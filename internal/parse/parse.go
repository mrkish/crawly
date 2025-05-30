package parse

import (
	"crawly/internal/cache"
	"crawly/internal/model"

	"golang.org/x/net/html"
)

func HTML(fetched model.Fetched, fetchQueue <-chan model.Parsed, cache *cache.Cache) {
	for node := range fetched.Doc.Descendants() {
		if !isAnchor(node) {
			continue
		}
		if cache.Has(node.Data) {
			// Update
		} else {
			// Add
		}
	}
}

func isAnchor(node *html.Node) bool {
	return node.Type == html.ElementNode && node.Data == "a"
}
