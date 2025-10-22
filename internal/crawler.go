package internal

import (
	"context"
	"fmt"
	"net/http"
	"sync"

	"golang.org/x/net/html"
)

func GetPageHTML(ctx context.Context, url string) (pageDom *html.Node, errReturned error) {
	request, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		errReturned = err
		return
	}

	response, err := http.DefaultClient.Do(request)
	if err != nil {
		return nil, err
	}

	defer response.Body.Close()

	pageDom, errReturned = html.Parse(response.Body)
	return
}

func GetPageURLs(page *html.Node) []string {
	urls := make([]string, 0)

	for element := range page.Descendants() {
		for _, attr := range element.Attr {
			if attr.Key == "href" {
				urls = append(urls, attr.Val)
			}
		}
	}

	return urls
}

func ScrapePage(id int, ctx context.Context, scrapeQueue chan string, wg *sync.WaitGroup) {
	defer wg.Done()
	for {
		select {
		case <-ctx.Done():
			fmt.Printf("worker %d: stopping\n", id)
			return
		case currentURL, ok := <-scrapeQueue:
			// channel closed
			if !ok {
				fmt.Printf("worker %d: stopping (channel closed)\n", id)
				return
			}

			// check that the url is valid and the page can be accessed
			page, err := GetPageHTML(ctx, currentURL)
			if err != nil {
				continue
			}

			fmt.Printf("worker %d: scraping %s\n", id, currentURL)

			// URLs on the current page
			nextURLs := GetPageURLs(page)

			for _, nextURL := range nextURLs {
				// This is to prevent us from getting into a loop
				if Set().Contains(nextURL) {
					continue
				}
				fmt.Printf("worker %d: found %s\n", id, nextURL)

				Set().Add(nextURL)

				select {
				case scrapeQueue <- nextURL: // send if queue isn't full, discarded if queue is full
				case <-ctx.Done():
					fmt.Printf("worker %d: stopping\n", id)
					return
				default:
				}
			}
		}
	}
}
