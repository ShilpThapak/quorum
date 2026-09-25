package main

import (
	"testing"
	"time"
)

// TestConcurrentProcessing verifies the producer and the consumer run
// concurrently. mockstream.Next() sleeps 320ms and IsTalkingAboutGo()
// sleeps 330ms; with 5 tweets sequential processing takes ~3.25s,
// while concurrent processing only takes ~1.7s.
func TestConcurrentProcessing(t *testing.T) {
	start := time.Now()
	main()
	elapsed := time.Since(start)

	if elapsed > 2200*time.Millisecond {
		t.Errorf("processing took %v; producer and consumer should run concurrently", elapsed)
	}
}
