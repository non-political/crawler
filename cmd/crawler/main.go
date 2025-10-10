package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/non-political/crawler/internal"
)

func main() {
	seedListBytes, err := os.ReadFile("seeds.txt")
	if err != nil {
		fmt.Fprintf(os.Stderr, "[ERROR]: Failed to load the seeds list: %v\n", err)
		os.Exit(-1)
	}

	foundPages := make(chan string, 10)

	seedList := string(seedListBytes)
	for seed := range strings.Lines(seedList) {
		go internal.ScrapePage(strings.TrimSpace(seed), foundPages)
	}

	for url := range foundPages {
		fmt.Printf("Found: %s\n", url)
	}
}
