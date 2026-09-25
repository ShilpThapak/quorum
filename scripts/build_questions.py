#!/usr/bin/env python3
"""Rebuild every question in questions-src/ from the canonical upstream
clone at /tmp/gce (github.com/loong/go-concurrency-exercises).

Layout per question NN-slug/:
  meta.json       - credits + description (hand-authored, keep)
  support/*.go    - DO-NOT-EDIT helpers, verbatim upstream
  starter.go      - buggy starting point, verbatim upstream main.go
  test.go         - grader (check_test.go upstream, or authored when
                    upstream ships no grader)
  reference.go    - the correct fix (authored here)
"""
import json, os, shutil, sys

ROOT = os.path.dirname(os.path.dirname(os.path.abspath(__file__)))
SRC = os.path.join(ROOT, "questions-src")
UP = "/tmp/gce"

QMAP = {
    "01-rate-limited-crawler": dict(up="0-limit-crawler",
        support=["mockfetcher.go"], starter="main.go",
        test_kind="check_test.go", grader=None),
    "02-producer-consumer": dict(up="1-producer-consumer",
        support=["mockstream.go"], starter="main.go",
        test_kind="none", grader="authored"),
    "03-thread-safe-cache": dict(up="2-race-in-cache",
        support=["mockdb.go", "mockserver.go"], starter="main.go",
        test_kind="check_test.go", grader=None),
    "04-session-cleaner": dict(up="5-session-cleaner",
        support=["helper.go", "errors.go"], starter="main.go",
        test_kind="check_test.go", grader=None),
}

CACHE_SUPPORT = """// DO NOT EDIT - the mock DB and server that drive the cache exercise.

package main

import (
	"strconv"
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

type MockDB struct {
	Calls int32
}

// Get simulates a slow DB read, returning the key as the value.
func (db *MockDB) Get(key string) (string, error) {
	time.Sleep(20 * time.Millisecond)
	atomic.AddInt32(&db.Calls, 1)
	return key, nil
}

// GetMockDB returns a fresh mock DB.
func GetMockDB() *MockDB { return &MockDB{} }

// Loader loads from the mock DB.
type Loader struct{ DB *MockDB }

// Load fetches a key's value from the DB.
func (l *Loader) Load(key string) string {
	val, err := l.DB.Get(key)
	if err != nil {
		panic(err)
	}
	return val
}

const (
	cycles        = 15
	callsPerCycle = 100
)

// RunMockServer hammers the cache from many goroutines.
func RunMockServer(cache *KeyStoreCache, t *testing.T) {
	var wg sync.WaitGroup
	for c := 0; c < cycles; c++ {
		wg.Add(1)
		go func() {
			for i := 0; i < callsPerCycle; i++ {
				wg.Add(1)
				go func(i int) {
					if v := cache.Get("Test" + strconv.Itoa(i)); v != "Test"+strconv.Itoa(i) {
						t.Errorf("incorrect db response %v", v)
					}
					wg.Done()
				}(i)
			}
			wg.Done()
		}()
	}
	wg.Wait()
}
"""

# Upstream expects these symbol names from the student's main.go;
# keep them identical to upstream so starter <-> grader line up.
for slug, cfg in QMAP.items():
    d = os.path.join(SRC, slug)
    updir = os.path.join(UP, cfg["up"])
    os.makedirs(os.path.join(d, "support"), exist_ok=True)

    # support: DO NOT EDIT helpers
    for s in cfg["support"]:
        shutil.copy(os.path.join(updir, s), os.path.join(d, "support", s))

    # starter: buggy main.go verbatim
    shutil.copy(os.path.join(updir, cfg["starter"]), os.path.join(d, "starter.go"))

print("copied support + starter for all questions")
