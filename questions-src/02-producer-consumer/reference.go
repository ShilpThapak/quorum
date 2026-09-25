//////////////////////////////////////////////////////////////////////
//
// Given is a producer-consumer scenario, where a producer reads in
// tweets from a mockstream and a consumer is processing the
// data. Your task is to change the code so that the producer as well
// as the consumer can run concurrently.
//
// Reference (correct) solution: the producer and consumer run
// concurrently. The producer sends every tweet over a channel *as soon
// as it has read it* from the stream, and a separate goroutine drives
// it. The consumer drains the same channel. Because the channel
// decouples the two, the consumer starts processing the first tweet
// while the producer is still reading the remaining ones.
//
// Compare with `starter.go` which runs producer then consumer
// sequentially.
//

package main

import (
	"fmt"
	"time"
)

func producer(stream Stream) <-chan *Tweet {
	tweets := make(chan *Tweet, 5)

	go func() {
		defer close(tweets)
		for {
			tweet, err := stream.Next()
			if err == ErrEOF {
				return
			}
			tweets <- tweet
		}
	}()

	return tweets
}

func consumer(tweets <-chan *Tweet) {
	for tweet := range tweets {
		if tweet.IsTalkingAboutGo() {
			fmt.Println(tweet.Username, "\ttweets about golang")
		} else {
			fmt.Println(tweet.Username, "\tdoes not tweet about golang")
		}
	}
}

func main() {
	start := time.Now()

	consumer(producer(GetMockStream()))

	fmt.Printf("Process took %s\n", time.Since(start))
}
