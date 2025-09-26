package main

import ( 
	"os"
	"fmt"
	"strings"

	"github.com/non-political/crawler/internal" 
)

func main() {
	seedListBytes, err := os.ReadFile("seeds.txt")
	if err != nil {
		fmt.Fprintf(os.Stderr, "[ERROR]: Failed to load the seeds list: %v\n", err)
		os.Exit(-1)
	}

	seedList := string(seedListBytes)
	for seed := range strings.Lines(seedList) {
		urls, err := internal.ScrapePage(strings.TrimSpace(seed))
		if err != nil {
			fmt.Fprintf(os.Stderr, "[ERROR]: %v\n", err)
		}

		for _, url := range urls {
			fmt.Printf("Found: %s\n", url)
		}
	}
}
