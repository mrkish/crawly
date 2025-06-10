package crawl

import (
	"context"
	"io"
	"log/slog"
	"net/url"
	"path/filepath"
	"strings"

	"github.com/mrkish/crawly/pkg/cache"
	"github.com/mrkish/crawly/pkg/log"

	"golang.org/x/net/html"
)

func parse(
	ctx context.Context,
	page io.ReadCloser,
	rootURL *url.URL,
	currentDepth int,
	cache *cache.Set[string],
) ([]Link, error) {
	defer page.Close()

	var links []Link

	foundLinks := make(map[string]struct{})
	tokens := html.NewTokenizer(page)
	isLinkInvalid := linkValidator(rootURL, foundLinks, cache)
	nextDepth := currentDepth + 1

	for {
		if ctx.Err() != nil {
			return links, ctx.Err()
		}

		tokenType := tokens.Next()
		token := tokens.Token()

		switch tokenType {
		case html.StartTagToken:
			if token.Data == a_tag {
				href := getHref(token.Attr)
				if isLinkInvalid(href) {
					continue
				}
				foundLinks[href] = struct{}{}
				link := Link{URL: href, Depth: nextDepth}
				links = append(links, link)
				slog.Debug("found link", slog.Any("link", link))
			}
		case html.ErrorToken:
			// end of page
			return links, nil
		}
	}
}

func linkValidator(rootURL *url.URL, found map[string]struct{}, cache *cache.Set[string]) func(string) bool {
	return func(href string) bool {
		if href == "" || isRelativeLink(href) || isExternalLink(rootURL, href) || isMediaLink(href) {
			return true
		}
		if _, ok := found[href]; ok || cache.Has(href) {
			log.Trace("seen link, skipping", slog.String(href_key, href))
			return true
		}
		return false
	}
}

func isExternalLink(root *url.URL, href string) bool {
	uri, err := url.Parse(href)
	if err != nil {
		return false
	}
	isExternal := uri.Host != root.Host
	if isExternal {
		log.Trace("href is external, skipping", slog.String(href_key, href))
	}
	return isExternal
}

func isRelativeLink(href string) bool {
	isRelLink := len(href) >= 1 && string(href[0]) == rel_hash
	if isRelLink {
		log.Trace("href is relative, skipping", slog.String(href_key, href))
	}
	return isRelLink
}

func isMediaLink(href string) bool {
	isFileLink := filepath.Ext(href) != ""
	if isFileLink {
		log.Trace("href is for media, skipping", slog.String(href_key, href))
	}
	return isFileLink
}

func getHref(attr []html.Attribute) string {
	for _, a := range attr {
		if a.Key == href_key {
			return strings.TrimRight(a.Val, slash)
		}
	}
	return ""
}
