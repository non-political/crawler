package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/non-political/crawler/internal"
)

const workers = 50
const maxPages = 10000
const queueSize = 100

func main() {
	seedListBytes, err := os.ReadFile("seeds.txt")
	if err != nil {
		fmt.Fprintf(os.Stderr, "[ERROR]: Failed to load the seeds list: %v\n", err)
		os.Exit(-1)
	}

	seeds := string(seedListBytes)
	var seedList []string
	for seed := range strings.Lines(seeds) {
		seedList = append(seedList, strings.TrimSpace(seed))
	}

	crawl := internal.CrawlInfo{
		Workers:   workers,
		MaxPages:  maxPages,
		QueueSize: queueSize,
		Seeds:     seedList,
	}
	numPages := crawl.StartCrawl()
	fmt.Printf("Found %d pages\n", numPages)
}
