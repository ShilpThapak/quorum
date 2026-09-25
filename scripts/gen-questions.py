#!/usr/bin/env python3
"""Regenerate the four question bundles under questions-src/.

Each bundle NN-<slug>/ contains:
  meta.json      - exercise metadata + credits
  starter.go     - buggy starting point (upstream verbatim)
  reference.go   - correct solution (= upstream main + the 3-line fix)
  test.go        - grader test (race + vet clean)
  support/       - DO-NOT-EDIT helpers (upstream verbatim)
"""
import json
import os
import shutil

ROOT = os.path.dirname(os.path.dirname(os.path.abspath(__file__)))
SRC = os.path.join(ROOT, "questions-src")

UPSTREAM = {
    "01-rate-limited-crawler": {
        "up": "0-limit-crawler",
        "support": ["mockfetcher.go"],
        "starter": "main.go",
    },
    "02-producer-consumer": {
        "up": "1-producer-consumer",
        "support": ["mockstream.go"],
        "starter": "main.go",
    },
    "03-thread-safe-cache": {
        "up": "2-race-in-cache",
        "support": ["mockdb.go", "mockserver.go"],
        "starter": "main.go",
    },
    "04-session-cleaner": {
        "up": "5-session-cleaner",
        "support": ["helper.go", "errors_stub.go"],
        "starter": "main.go",
    },
}

def main():
    up_root = os.environ.get("UPSTREAM_DIR", "/tmp/gce")

    for slug, cfg in UPSTREAM.items():
        dst = os.path.join(SRC, slug)
        up = os.path.join(up_root, cfg["up"])
        os.makedirs(os.path.join(dst, "support"), exist_ok=True)

        # support: verbatim upstream
        for s in cfg["support"]:
            shutil.copyfile(os.path.join(up, s),
                            os.path.join(dst, "support", s))
        # starter: verbatim upstream main.go
        shutil.copyfile(os.path.join(up, cfg["starter"]),
                        os.path.join(dst, "starter.go"))
        print(f"[copied] {slug}: support={cfg['support']}")

if __name__ == "__main__":
    main()
