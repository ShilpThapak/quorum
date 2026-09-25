# AGENTS.md — Quorum

Quorum is a local Go-concurrency lab: 4 exercises (rate-limited crawler,
producer-consumer, thread-safe cache, session cleaner) where a hidden grader
`go vet`s the user's solution and runs hidden tests with `go test -race`.

## Structure
- `questions-src/<slug>/` — one folder per exercise: `starter.go` (user's
  starting point), `reference.go` (correct), `meta.json`. Private grading
  contract (never shown to the user): hidden `check_test.go` + grader-side
  answers.
- `server/` — Go grader (`quorum-server`). Reads env `QUORUM_QUESTIONS` (default
  `../questions-src`), grep-ensures hidden files aren't served, builds a scratch
  module per attempt, runs `go vet` then `go test -race` as subprocesses (a
  working Go toolchain must be present at runtime). Routes: `/api/health`,
  `/api/questions`, `/api/source/{slug}/{kind}`, `/api/grade`. Port from `PORT`
  (default 8090). Report always emits `cases` as `[]` (never `null`).
- `web/` — Vite + React + Monaco editor. Dev proxy `/api` -> localhost:8090
  via `vite.config.ts`. Reads reference/starter through `loadSource` (must
  unwrap `.code`, never `res.text()`).

## Editing rules
- DO NOT modify anything under `questions-src/<slug>/reference.go` or hidden
  grading files — that's the grading contract. Touching it breaks the lab.
- Keep files small; small-diff edits. Avoid giant `write` tool calls to root
  paths (harness path quirk: the root is seen as non-existent; a 0-byte
  AGENTS.md was the result). Prefer writing temp files + `mv`, or edits inside
  `server/` / `web/src/` which work fine.
- No comments in code unless asked. No new files unless required.

## The #1 rule — verify after EVERY change
After every change (server OR web), prove it works before stopping:
1. Build server: `cd server && go build ./... && go vet ./...`
2. Build web: `cd web && npm run build` (type-checks too)
3. Start `go run .` in server/, then run the full grader matrix via the web
   proxy: reference -> PASS and starter -> FAIL for ALL 4 questions, and the
   server must not crash or log errors (watch for a blank Results pane or a
   React `Cannot read properties of null (reading 'map')` — those mean a stale
   server binary or a null `cases`; fix by rebuilding the server AND null-
   guarding `report.cases ?? []` in App.tsx).
4. Restart the dev server from the FRESH build (a stale binary serves old
   routes and broken `cases`, which is exactly the reported blank page).

## Frontend never-blank rules
- `report.cases` may be `null` from a stale/old server -> render with a guard
  (`(report.cases ?? []).map(...)`); it already is, do not regress it.
- Report fields (`vetNotes`, `output`, `cases`) must all render null-tolerant.
- Editor text must be real source (unwrap `.code`), never an escaped JSON
  envelope.
