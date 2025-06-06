package crawl

import (
	"context"
	"crawly/internal/cache"
	"crawly/pkg/log"
	"log/slog"

	// "crawly/internal/semaphore"
	"net/url"

)

func FromRoot(ctx context.Context, root string, maxDepth int) ([]Page, error) {
	rootURL, err := url.Parse(root)
	if err != nil {
		return nil, err
	}

	cache := cache.New[string, Link](mergeMap)

	pageBuilder := pageFactory(rootURL)
	homepage := pageBuilder(root, 1)

	type crawledPage struct {
		links []Link
		err   error
	}
	pageChan := make(chan crawledPage)
	defer close(pageChan)

	// for range maxDepth {
	go func() {
		links, err := crawl(ctx, homepage, cache)
		// pages = append(pages, linkToPage(pageBuilder, iterLinks...)...)
		pageChan <- crawledPage{links: links, err: err}
	}()
	// }

	var pages []Page
	for result := range pageChan {
		if result.err != nil {
			slog.Error("error crawling page", log.Err(result.err))
			continue
		}
		pages = append(pages, linkToPage(pageBuilder, result.links...)...)
	}

	return pages, nil
}

func linkToPage(builder PageBuilder, links ...Link) []Page {
	var out []Page
	for _, l := range links {
		out = append(out, builder(l.URL, l.Depth))
	}
	return out
}

// type crawlResult struct {
// 	Links []Link
// 	err   error
// }
//
// func fanOutSem(
// 	ctx context.Context,
// 	root *url.URL,
// 	currentDepth,
// 	workers int,
// 	links []Link,
// 	cache *cache.Cache[string, Link],
// ) {
// 	sem := semaphore.New(workers)
// 	out := make(chan crawlResult, len(links))
//
// 	for _, link := range links {
// 		sem.Acquire()
// 		go func(l Link) {
// 			defer sem.Free()
// 			li, err := crawl(ctx, root, link.URL, currentDepth, cache)
// 			out <- crawlResult{li, err}
// 		}(link)
// 	}
// }

func crawl(
	ctx context.Context,
	page Page,
	cache *cache.Cache[string, Link],
) ([]Link, error) {
	pageData, err := fetch(ctx, page.URL)
	if err != nil {
		return nil, err
	}
	return parse(ctx, page.Root, pageData, page.Depth, cache)
}
