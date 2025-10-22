package internal

import (
	"context"
	"sync"
	"time"
)

type CrawlInfo struct {
	Workers   int
	MaxPages  int
	QueueSize int
	Seeds     []string
}

func (info *CrawlInfo) StartCrawl() int {
	scrapeQueue := make(chan string, info.QueueSize)
	for _, seed := range info.Seeds {
		scrapeQueue <- seed
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	var wg sync.WaitGroup

	// start workers
	for i := 1; i <= info.Workers; i++ {
		wg.Add(1)
		go ScrapePage(i, ctx, scrapeQueue, &wg)
	}

	// check every 100ms whether we have reached target, might overshoot
	go func() {
		for {
			time.Sleep(100 * time.Millisecond)
			if Set().Size() >= info.MaxPages {
				cancel()
				return
			}
		}
	}()

	// wait until all workers done (ensure resources released cleanly)
	wg.Wait()

	close(scrapeQueue)

	return Set().Size()
}
