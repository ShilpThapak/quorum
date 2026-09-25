//////////////////////////////////////////////////////////////////////
//
// Rate-limit grader: at most one Fetch per ~1s while Crawl stays
// concurrent. Written to be `go vet` clean (no FailNow in a
// non-test goroutine).
//

package main

import (
	"sync"
	"testing"
	"time"
)

func TestRateLimitedCrawler(t *testing.T) {
	fetchSig := fetchSignalInstance()

	var (
		mu      sync.Mutex
		tooFast bool
		fetches int
	)
	start := time.Now()

	go func() {
		for range fetchSig {
			mu.Lock()
			fetches++
			if fetches > 1 && time.Since(start).Nanoseconds() < 950000000 {
				tooFast = true
			}
			start = time.Now()
			mu.Unlock()
		}
	}()

	main()

	mu.Lock()
	defer mu.Unlock()
	if fetches < 2 {
		t.Errorf("expected at least 2 fetches, got %d", fetches)
	}
	if tooFast {
		t.Errorf("two crawls executed less than 1 second apart; solution is incorrect")
	}
}
