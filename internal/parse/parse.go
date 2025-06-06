package parse

import (
	"crawly/internal/cache"
	"crawly/internal/model"
	"fmt"
	"net/url"

	"golang.org/x/net/html"
)

const (
	HTML_BODY        = "body"
	HREF_KEY         = "href"
	REL_HASH  string = "#"
)

func GetBody(doc *html.Node) *html.Node {
	for n := range doc.Descendants() {
		if n.Type == html.ElementNode && n.Data == HTML_BODY {
			return n
		}
	}
	return nil
}

func Tokens(fetched model.Fetched, cache *cache.Cache) {
	for {
		t := fetched.Tokens.Next()
		switch t {
		case html.ErrorToken:
			return
		case html.StartTagToken:
			token := fetched.Tokens.Token()
			if token.Data == "a" {
				href := getHref(token.Attr)
				fmt.Printf("found link: %s", href)
				cache.Add(&model.Link{
					URL: href,
				})
			}
		}
	}
}

// func Node(node *html.Node) []*html.Node {
// 	var parseableChildren []*html.Node
// 	for node := range node.Descendants() {
// 		if node.Type == html.ElementNode {
// 			switch node.Data {
// 			case "h1", "h2", "h3", "h4", "h5", "div", "span", "p", "li":
// 				if nodes := Node(node); len(nodes) > 0 {
// 					parseableChildren = append(parseableChildren, nodes...)
// 				}
// 			case "a":
// 				href := getHref(node.Attr)
// 				uri, err := url.Parse(href)
// 				if err != nil {
// 					continue
// 				}
// 				cache.Add(&model.Link{
// 					URL:   fetched.URL,
// 					Level: fetched.Level,
// 					LinksTo: map[string]*model.Link{
// 						fetched.URL: {URL: url},
// 					},
// 				})
//
// 				// Add link FROM
// 				cache.Add(&model.Link{
// 					URL:   url,
// 					Level: fetched.Level + 1,
// 					LinksFrom: map[string]*model.Link{
// 						fetched.URL: {URL: fetched.URL},
// 					},
// 				})
//
// 			default:
// 			}
// 		}
// 	}
// 	return parseableChildren
// }

func HTML(
	fetched model.Fetched,
	parseQueue chan<- model.Fetched,
	fetchQueue <-chan model.Parsed,
	cache *cache.Cache,
	depth int,
) {
	for node := range fetched.Doc.Descendants() {
		if !isAnchor(node) {
			switch node.Type {
			case html.ElementNode:
				parseQueue <- model.Fetched{
					Doc:   node,
					Level: fetched.Level,
					URL:   fetched.URL,
					Root:  fetched.Root,
				}
			}
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
		url := uri.Host + uri.Path
		// Add link TO
		cache.Add(&model.Link{
			URL:   fetched.URL,
			Level: fetched.Level,
			LinksTo: map[string]*model.Link{
				fetched.URL: {URL: url},
			},
		})

		// Add link FROM
		cache.Add(&model.Link{
			URL:   url,
			Level: fetched.Level + 1,
			LinksFrom: map[string]*model.Link{
				fetched.URL: {URL: fetched.URL},
			},
		})
		// if cache.Has(uri.Host + uri.Path) {
		// 	fmt.Printf("html node: %v\n", node)
		// 	// Update
		// } else {
		// 	cache.Add(&model.Link{
		// 		URL:   uri.Host + uri.Path,
		// 		Level: fetched.Level + 1,
		// 	})
		fmt.Printf("html node: %v\n", node)
		// }
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
