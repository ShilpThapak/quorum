//////////////////////////////////////////////////////////////////////
//
// Given is a mockfetcher that fetches URLs. The fetcher is slow, so
// we want to limit the crawler to at most one page per second while
// STILL crawling concurrently.
//
// Reference (correct) solution: a shared 1-second ticker. Each crawl
// goroutine waits for the tick before fetching, so no two fetches
// start within the same second, yet all crawls still run in parallel.
//

package main

import (
	"fmt"
	"sync"
	"time"
)

// rateTicker pulses once per second and is shared by every crawl
// goroutine: the rate limiter.
var rateTicker = time.NewTicker(time.Second)

// Crawl uses `fetcher` from the `mockfetcher.go` file to imitate a
// real crawler. It crawls until the maximum depth has reached.
func Crawl(url string, depth int, wg *sync.WaitGroup) {
	defer wg.Done()

	if depth <= 0 {
		return
	}

	// Wait for the next tick, enforcing "at most one page per second"
	// while crawls still run concurrently.
	<-rateTicker.C

	body, urls, err := fetcher.Fetch(url)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Printf("found: %s %q\n", url, body)

	wg.Add(len(urls))
	for _, u := range urls {
		// Do not remove the `go` keyword, as Crawl() must be
		// called concurrently
		go Crawl(u, depth-1, wg)
	}
}

func main() {
	var wg sync.WaitGroup

	wg.Add(1)
	Crawl("http://golang.org/", 4, &wg)
	wg.Wait()
}
