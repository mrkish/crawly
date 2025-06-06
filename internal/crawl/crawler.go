package crawl

import (
	"context"
	"crawly/internal/cache"
	"errors"
	"net/url"
)

var ErrMaxDepthStop = errors.New("max depth reached")

func FromRoot(ctx context.Context, root string, maxDepth int) ([]Link, error) {
	rootURL, err := url.Parse(root)
	if err != nil {
		return nil, err
	}

	seen := make(map[string]Link)
	cache := cache.New()

	pageData, err := fetch(ctx, root)
	if err != nil {
		return nil, err
	}

	return parse(ctx, rootURL, pageData, maxDepth, 1, seen)
}

func crawl(ctx context.Context, root *url.URL, pageURL string, maxDepth, currentDepth int) ([]Link, error) {
	if maxDepth == currentDepth {
		return nil, ErrMaxDepthStop
	}

	seen := make(map[string]Link)

	pageData, err := fetch(ctx, pageURL)
	if err != nil {
		return nil, err
	}

	return parse(ctx, root, pageData, maxDepth, currentDepth, seen)
}
