package crawl

import (
	"context"
	"crawly/pkg/log"
	"io"
	"log/slog"
	"net/url"

	"golang.org/x/net/html"
)

const (
	HTML_BODY   string = "body"
	HREF_KEY    string = "href"
	REL_HASH    string = "#"
	QUERY_PARAM string = "?"
)

func parse(ctx context.Context, root *url.URL, page io.ReadCloser, maxDepth, currentDepth int, seen map[string]Link) ([]Link, error) {
	defer page.Close()
	if maxDepth == currentDepth {
		return nil, ErrMaxDepthStop
	}

	if seen == nil {
		seen = make(map[string]Link)
	}

	var links []Link
	tokens := html.NewTokenizer(page)

	for {
		if ctx.Err() != nil {
			return links, ctx.Err()
		}

		tokenType := tokens.Next()
		token := tokens.Token()

		switch tokenType {
		case html.StartTagToken:
			if token.Data == "a" {
				href := getHref(token.Attr)
				if href == "" {
					continue
				}
				if _, ok := seen[href]; ok {
					slog.Log(ctx, log.Trace, "seen link, skipping", slog.String("href", href))
					continue
				}
				if isRelativeLink(href) {
					slog.Log(ctx, log.Trace, "href is relative, skipping", slog.String("href", href))
					continue
				}
				if isExternalLink(root, href) {
					slog.Log(ctx, log.Trace, "href is external, skipping", slog.String("href", href))
					continue
				}
				// is a link we can parse next
				link := Link{URL: href, Depth: currentDepth}
				seen[href] = link
				slog.Log(ctx, log.Trace, "found link", slog.Any("link", link))
				links = append(links, link)
			}
		case html.ErrorToken:
			// end of page
			return links, nil
		}
	}
}

func isExternalLink(root *url.URL, href string) bool {
	uri, err := url.Parse(href)
	if err != nil {
		return false
	}
	return uri.Host != root.Host
}

func isRelativeLink(href string) bool {
	return len(href) > 1 && (string(href[0]) == REL_HASH || string(href[0]) == QUERY_PARAM)
}

func getHref(attr []html.Attribute) string {
	// fmt.Printf("parsing tag attributes: %v\n", attr)
	for _, a := range attr {
		// fmt.Printf("tag attr: %s\n", a.Key)
		if a.Key == HREF_KEY {
			// fmt.Printf("found href value: %s\n", a.Val)
			return a.Val
		}
	}
	return ""
}
