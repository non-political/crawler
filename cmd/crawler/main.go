package main

import (
	"context"
	"fmt"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/non-political/crawler/internal"
)

const workers = 50
const maxPages = 10000

func main() {
	seedListBytes, err := os.ReadFile("seeds.txt")
	if err != nil {
		fmt.Fprintf(os.Stderr, "[ERROR]: Failed to load the seeds list: %v\n", err)
		os.Exit(-1)
	}

	scrapeQueue := make(chan string, 10)
	seedList := string(seedListBytes)
	for seed := range strings.Lines(seedList) {
		scrapeQueue <- strings.TrimSpace(seed)
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	var wg sync.WaitGroup

	for i := 1; i <= workers; i++ {
		wg.Add(1)
		go internal.ScrapePage(i, ctx, scrapeQueue, &wg)
	}

	go func() {
		for {
			time.Sleep(100 * time.Millisecond)
			if internal.Set().Size() >= maxPages {
				cancel()
				return
			}
		}
	}()

	wg.Wait()
	close(scrapeQueue)
	fmt.Printf("Found %d pages\n", internal.Set().Size())
}
