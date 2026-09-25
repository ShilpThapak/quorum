package main

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"runtime"
	"strings"
	"time"
)

// graderToolchainOK reports whether a working Go toolchain is reachable from
// this process. It is the cheap signal /api/health uses so the web UI can gate.
func graderToolchainOK() bool {
	goBin, err := exec.LookPath("go")
	if err != nil {
		return false
	}
	ctx, cancel := context.WithTimeout(context.Background(), 4*time.Second)
	defer cancel()
	out, err := exec.CommandContext(ctx, goBin, "env", "GOROOT").Output()
	return err == nil && len(bytes.TrimSpace(out)) > 0 && runtime.GOOS == "windows" == false
}



// Case is one hidden-test result.
type Case struct {
	Name    string  `json:"name"`
	Passed  bool    `json:"passed"`
	Elapsed float64 `json:"elapsed"`
}

// Report is the JSON returned for one grade attempt.
type Report struct {
	Question string `json:"question"`
	Kind     string `json:"kind"`
	VetOK    bool   `json:"vetOK"`
	TestOK   bool   `json:"testOK"`
	Race     bool   `json:"race"`
	Cases    []Case `json:"cases"`
	Output   string `json:"output,omitempty"`
	TimeMs   int64  `json:"timeMs"`
}

// Grader grades one candidate for one question.
type Grader struct {
	Base    string // path to questions-src
	Scratch string // temp base dir
	GoCmd   string
	Timeout time.Duration
}

func NewGrader(base string) *Grader {
	return &Grader{
		Base:    base,
		Scratch: os.TempDir(),
		GoCmd:   "go",
		Timeout: 5 * time.Minute,
	}
}

var (
	reCase = regexp.MustCompile(`(?m)^--- (PASS|FAIL): (\S+) \(([^)]*)\)`)
)

// copyIn copies every *.go in srcDir into dstDir.
func copyIn(srcDir, dstDir string) error {
	entries, err := os.ReadDir(srcDir)
	if err != nil {
		return err
	}
	for _, e := range entries {
		if e.IsDir() || !strings.HasSuffix(e.Name(), ".go") {
			continue
		}
		data, err := os.ReadFile(filepath.Join(srcDir, e.Name()))
		if err != nil {
			return err
		}
		if err := os.WriteFile(filepath.Join(dstDir, e.Name()), data, 0o600); err != nil {
			return err
		}
	}
	return nil
}

// run executes go <args...> in dir. Returns combined output + exit error.
func (g *Grader) run(ctx context.Context, dir string, args ...string) (string, error) {
	cctx, cancel := context.WithTimeout(ctx, g.Timeout)
	defer cancel()
	cmd := exec.CommandContext(cctx, g.GoCmd, args...)
	cmd.Dir = dir
	cmd.Env = append(os.Environ(), "GOWORK=off")
	var buf bytes.Buffer
	cmd.Stdout = &buf
	cmd.Stderr = &buf
	err := cmd.Run()
	return buf.String(), err
}

// Grade assembles a scratch module and grades it.
//
//	kind: "reference" | "starter" | "user"
//	code: only used when kind == "user"
func (g *Grader) Grade(ctx context.Context, slug, kind string, code []byte) (*Report, error) {
	start := time.Now()
	qdir := filepath.Join(g.Base, slug)
	if fi, err := os.Stat(qdir); err != nil || !fi.IsDir() {
		return nil, fmt.Errorf("question %q: not found", slug)
	}

	scratch, err := os.MkdirTemp(g.Scratch, "quorum-"+slug+"-")
	if err != nil {
		return nil, err
	}
	defer os.RemoveAll(scratch)

	// question-level support files
	if err := copyIn(filepath.Join(qdir, "support"), scratch); err != nil && !os.IsNotExist(err) {
		return nil, err
	}
	// candidate
	if kind == "user" {
		if err := os.WriteFile(filepath.Join(scratch, "main.go"), code, 0o600); err != nil {
			return nil, err
		}
	} else {
		body, err := os.ReadFile(filepath.Join(qdir, kind+".go"))
		if err != nil {
			return nil, err
		}
		if err := os.WriteFile(filepath.Join(scratch, "main.go"), body, 0o600); err != nil {
			return nil, err
		}
	}
	// hidden grader test
	if body, err := os.ReadFile(filepath.Join(qdir, "check_test.go")); err == nil {
		if err := os.WriteFile(filepath.Join(scratch, "check_test.go"), body, 0o600); err != nil {
			return nil, err
		}
	}

	// module file
	if err := os.WriteFile(filepath.Join(scratch, "go.mod"),
		[]byte("module grader\n\ngo 1.24\n"), 0o600); err != nil {
		return nil, err
	}

	r := &Report{Question: slug, Kind: kind}

	// vet first
	if out, err := g.run(ctx, scratch, "vet", "./..."); err != nil {
		r.VetOK = false
		r.Output = out
		r.TimeMs = time.Since(start).Milliseconds()
		return r, nil
	}
	r.VetOK = true

	// test -race
	args := []string{"test", "-count=1", "-race", "-timeout", "180s", "./..."}
	out, err := g.run(ctx, scratch, args...)
	r.Output = out
	r.Race = strings.Contains(out, "DATA RACE")
	for _, m := range reCase.FindAllStringSubmatch(out, -1) {
		r.Cases = append(r.Cases, Case{Name: m[2], Passed: m[1] == "PASS", Elapsed: 0})
	}
	r.TestOK = err == nil && !r.Race
	r.TimeMs = time.Since(start).Milliseconds()
	if r.Cases == nil {
		r.Cases = []Case{} // never null in JSON: the frontend maps over cases
	}
	return r, nil
}
