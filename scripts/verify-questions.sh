#!/bin/bash
# Verifies every question exactly the way the grader will: assemble a
# fresh module in a temp dir per (question, variant) and run
# `go vet` + `go test -race`. Reference must PASS; starter must FAIL.
#
# EXPECTED OUTCOME:
#   reference (correct solution) -> vet OK, test OK (race-clean)
#   starter   (buggy)            -> vet OK, test FAIL (or race report)
set -euo pipefail

ROOT="$(cd "$(dirname "$0")/.." && pwd)"
SRC="$ROOT/questions-src"
TMP="$(mktemp -d)"
trap 'rm -rf "$TMP"' EXIT

for qdir in "$SRC"/[0-9][0-9]-*/; do
  q="$(basename "$qdir")"
  for kind in reference starter; do
    file="$qdir/$kind.go"
    [ -f "$file" ] || continue

    work="$TMP/$q-$kind/main"
    mkdir -p "$work"
    cp "$qdir"/support/*.go "$work/"
    cp "$file" "$work/main.go"
    cp "$qdir"/test.go "$work/test.go"
    printf 'module grader\n\ngo 1.24\n' > "$TMP/$q-$kind/go.mod"

    vet=FAIL; tst=FAIL; race=""
    if ( cd "$work" && go vet ./... ) >"$work/vet.log" 2>&1; then
      vet=ok
      if ( cd "$work" && go test -race -count=1 -timeout 120s ./... ) >"$work/test.log" 2>&1; then
        tst=PASS
      else
        grep -qi "DATA RACE" "$work/test.log" && race="[race]"
      fi
    fi
    printf "%-22s %-9s vet=%-4s test=%-6s %s\n" "$q" "$kind" "$vet" "$tst" "$race"
  done
done
