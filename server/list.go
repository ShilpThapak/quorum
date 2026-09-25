package main

import (
	"encoding/json"
	"os"
	"path/filepath"
	"sort"
)

type QuestionMeta struct {
	ID         string   `json:"id"`
	Title      string   `json:"title"`
	Slug       string   `json:"slug"`
	Difficulty int      `json:"difficulty"`
	Tags       []string `json:"tags"`
	Task       string   `json:"task"`
	TimeoutSec int      `json:"timeoutSec"`
}

// List returns meta for every questions-src/NN-slug dir.
func (g *Grader) List() ([]QuestionMeta, error) {
	dirs, _ := filepath.Glob(filepath.Join(g.Base, "[0-9][0-9]-*"))
	sort.Strings(dirs)
	var out []QuestionMeta
	for _, d := range dirs {
		b, err := os.ReadFile(filepath.Join(d, "meta.json"))
		if err != nil {
			continue
		}
		var m QuestionMeta
		if err := json.Unmarshal(b, &m); err != nil {
			continue
		}
		m.Slug = filepath.Base(d)
		out = append(out, m)
	}
	return out, nil
}
