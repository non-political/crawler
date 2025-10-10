package internal

import (
	"net/http"

	"golang.org/x/net/html"
)

func GetPageHTML(url string) (pageDom *html.Node, errReturned error) {
	response, err := http.Get(url)
	if err != nil {
		errReturned = err
		return
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

func ScrapePage(pageURL string, foundChannel chan string) {
	page, err := GetPageHTML(pageURL)
	if err != nil {
		return
	}

	// Ensure that the page is actually a page before we say we found it
	foundChannel <- pageURL
	currentURLs := GetPageURLs(page)

	// Also, right now the scraping is done recursively which may cause some problems...
	for _, url := range currentURLs {
		// This is to prevent us from getting into a loop
		if Set().Contains(url) {
			continue
		}

		Set().Add(url)
		go ScrapePage(url, foundChannel)
	}
}
