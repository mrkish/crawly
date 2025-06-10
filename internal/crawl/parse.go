package crawl

import (
	"context"
	"crawly/pkg/cache"
	"crawly/pkg/log"
	"io"
	"log/slog"
	"net/url"
	"path/filepath"
	"strings"

	"golang.org/x/net/html"
)

const (
	HTML_BODY   string = "body"
	HREF_KEY    string = "href"
	REL_HASH    string = "#"
	QUERY_PARAM string = "?"
	SLASH       string = "/"
)

func parse(
	ctx context.Context,
	root *url.URL,
	page io.ReadCloser,
	currentDepth int,
	cache *cache.Cache[string, Link],
) ([]Link, error) {
	defer page.Close()

	var links []Link

	// use a map to de-dupe output
	foundLinks := make(map[string]Link)
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
				if cache.Has(href) {
					slog.Log(ctx, log.Trace, "seen link, skipping", slog.String("href", href))
					continue
				}
				if isRelativeLink(href) {
					cache.Add(href, Link{URL: href, Depth: currentDepth})
					slog.Log(ctx, log.Trace, "href is relative, skipping", slog.String("href", href))
					continue
				}
				if isExternalLink(root, href) {
					cache.Add(href, Link{URL: href, Depth: currentDepth})
					slog.Log(ctx, log.Trace, "href is external, skipping", slog.String("href", href))
					continue
				}
				if isRelativeLink(href) {
					cache.Add(href, Link{URL: href, Depth: currentDepth})
					slog.Log(ctx, log.Trace, "href is relative, skipping", slog.String("href", href))
					continue
				}
				if isMediaLink(href) {
					cache.Add(href, Link{URL: href, Depth: currentDepth})
					slog.Log(ctx, log.Trace, "href is for media, skipping", slog.String("href", href))
					continue
				}
				link := Link{URL: href, Depth: currentDepth + 1}
				foundLinks[href] = link
				slog.Log(ctx, log.Trace, "found link", slog.Any("link", link))
			}
		case html.ErrorToken:
			// end of page

			for _, l := range foundLinks {
				links = append(links, l)
			}

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
	return len(href) >= 1 && string(href[0]) == REL_HASH
}

func isMediaLink(href string) bool {
	return filepath.Ext(href) != ""
}

func getHref(attr []html.Attribute) string {
	for _, a := range attr {
		if a.Key == HREF_KEY {
			return strings.TrimRight(a.Val, SLASH)
		}
	}
	return ""
}
