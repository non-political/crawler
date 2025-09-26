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

func ScrapePage(pageURL string) (foundURLs []string, err error) {
	page, err := GetPageHTML(pageURL)
	if err != nil {
		return
	}

	currentURLs := GetPageURLs(page)

	// Also, right now the scraping is done recursively which may cause some problems...
	for _, url := range currentURLs {
		foundURLs = append(foundURLs, url)

		// We should probably not ignore the error, but... not sure what to do
		// if it fails so...
		innerURLs, _ := ScrapePage(url)
		foundURLs = append(foundURLs, innerURLs...)
	}

	return 
}
