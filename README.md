# Quorum

Quorum is a small, honest coding lab for **Go concurrency**. It ships four
exercises (all WTFPL-licensed, adapted from
[go-concurrency-exercises](https://github.com/loong/go-concurrency-exercises)),
a Go learning/grader backend, and a browser IDE with Go language support.

Each exercise gives you a **buggy starter** (`starter.go`) in the editor.
Fix it until all hidden tests pass. The grader runs your code through
`go vet` and `go test -race` in an isolated scratch module and reports the
verdict — vet result, data-race detection, and PASS/FAIL per hidden test.

## The four exercises

| # | Exercise | Difficulty | What to fix |
|---|----------|-----------|-------------|
| 01 | Rate-limited crawler | ★★ | Add a shared 1s rate limiter so fetches leave at least ~1s apart |
| 02 | Producer-consumer | ★★ | Run the producer and consumer concurrently so nothing runs twice as slow |
| 03 | Thread-safe cache | ★★★ | Guard the LRU cache with a mutex so `-race` stops complaining |
| 04 | Session cleaner | ★★★★ | Fix the background cleaner so sessions outlive their expiry |

Correct reference solutions live in each `reference.go`; the grader itself is
verified to keep the matrix *reference → PASS, starter → FAIL*, vet-clean, for
all four.

## Layout

```
questions-src/NN-slug/      exercise: starter.go + reference.go + support/ + check_test.go + meta.json
server/                     Go grader + HTTP API  (:8090)
web/                        Vite + React + TypeScript + Monaco frontend (:5173)
scripts/                    build / verify helpers (Python + shell)
```

## Run it (dev)

Backend (needs Go ≥ 1.24):

```sh
cd server
QUORUM_QUESTIONS=../questions-src go run .        # listens on :8090
```

Frontend (needs Node ≥ 20):

```sh
cd web
npm install
npm run dev                                      # :5173, proxies /api → :8090
```

Open http://localhost:5173, pick an exercise, edit, hit **Run tests**.

### Production build

```sh
cd web && npm run build            # static site in web/dist
cd server && go build -o quorum .  # single static binary
QUORUM_QUESTIONS=../questions-src ./quorum       # serves web/dist too
```

## How grading works (server/grader.go)

1. Copy `support/*.go` (the do-not-edit harness) plus a hidden `check_test.go`
   into a fresh scratch module.
2. `go vet ./...` — the candidate must be vet-clean or it fails immediately.
3. `go test -race ./...` — hidden tests run under the race detector.
4. Output is parsed into per-test PASS/FAIL cases and a JSON `Report`.

The test files are named `check_test.go` deliberately: Go only compiles
`*_test.go` for the test binary, so the harness and hidden tests never
leak into the candidate's source.

## Verify the matrix

```sh
scripts/verify-questions.sh        # expect: all reference PASS, all starter FAIL
```

## License

Adapted exercises: WTFPL. Everything else here: MIT.
