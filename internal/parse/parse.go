package parse

import (
	"crawly/internal/cache"
	"crawly/internal/model"
	"fmt"
	"net/url"

	"golang.org/x/net/html"
)

const (
	HREF_KEY = "href"
	REL_HASH string = "#"
)

func HTML(fetched model.Fetched, fetchQueue <-chan model.Parsed, cache *cache.Cache) {
	for node := range fetched.Doc.Descendants() {
		if !isAnchor(node) {
			continue
		}
		href := getHref(node.Attr)
		uri, err := url.Parse(href)
		if err != nil {
			continue
		}
		// external link
		if uri.Host != fetched.Root.Host {
			continue
		}
		// relative link in page
		if isRelativeLink(href) {
			continue
		}
		if cache.Has(uri.Host + uri.Path) {
			fmt.Printf("html node: %v\n", node)
			// Update
		} else {
			fmt.Printf("html node: %v\n", node)
			// Add
		}
	}
}

func isRelativeLink(href string) bool {
	return len(href) > 1 && string(href[0]) == REL_HASH
}

func getHref(attr []html.Attribute) string {
	for _, a := range attr {
		if a.Key == HREF_KEY {
			return a.Val
		}
	}
	return ""
}

func isAnchor(node *html.Node) bool {
	return node.Type == html.ElementNode && node.Data == "a"
}
