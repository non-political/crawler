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
