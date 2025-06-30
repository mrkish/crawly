package crawl

import (
	"context"
	"log/slog"
	"net/url"
	"strings"
	"sync/atomic"

	"github.com/mrkish/crawly/internal/client"
	"github.com/mrkish/crawly/internal/constants"
	"github.com/mrkish/crawly/internal/fetch"
	"github.com/mrkish/crawly/internal/model"
	"github.com/mrkish/crawly/internal/parse"
	"github.com/mrkish/crawly/pkg/cache"
	"github.com/mrkish/crawly/pkg/log"
	"github.com/mrkish/crawly/pkg/semaphore"
)

func FromRoot(ctx context.Context, root string, workers, maxDepth int) ([]model.Page, error) {
	rootURL, err := url.Parse(root)
	if err != nil {
		return nil, err
	}

	cache := cache.NewSet(func(key string) string { return strings.TrimRight(key, constants.SLASH) })

	linkQueue := make(chan model.Link)
	defer close(linkQueue)

	parsedPages, err := startCrawl(ctx, rootURL, workers, maxDepth, cache, linkQueue)
	if err != nil {
		return nil, err
	}

	var pending uint32 = 1
	var pages []model.Page

	for {
		select {
		case <-ctx.Done():
			return pages, nil

		case page := <-parsedPages:
			log.Trace("received page result", slog.Any("pending", pending))
			atomic.SwapUint32(&pending, pending-1)

			pages = append(pages, page)

			if len(page.Links) > 0 {
				for _, link := range page.Links {
					if link.IsInvalid(maxDepth, cache.Has) {
						continue
					}
					atomic.SwapUint32(&pending, pending+1)
					linkQueue <- link
				}
			}

			slog.Debug("finished processing page result", slog.Uint64("pending", uint64(pending)))

			if pending == 0 {
				return pages, nil
			}
		}
	}
}

func startCrawl(
	ctx context.Context,
	rootURL *url.URL,
	workers, maxDepth int,
	cache *cache.Set[string],
	linkQueue <-chan model.Link,
) (<-chan model.Page, error) {
	var err error
	parsedPages := make(chan model.Page)
	sem := semaphore.New(workers)

	crawlPage := buildCrawler(rootURL, workers, sem, cache)

	rootPage := model.Page{Link: model.Link{URL: rootURL.String(), Depth: 1}}
	rootPage.Links, err = crawlPage(ctx, rootPage.Link)
	if err != nil {
		return nil, err
	}

	foundLinks := len(rootPage.Links) > 0

	if foundLinks {
		defer func() {
			go func() {
				parsedPages <- rootPage
				if maxDepth == 1 {
					close(parsedPages)
				}
			}()
		}()
	}

	if maxDepth > 1 && foundLinks {
		go func() {
			defer close(parsedPages)
			defer slog.Debug("crawling goroutine exited")

			for {
				select {
				case <-ctx.Done():
					return
				case l := <-linkQueue:
					go func(link model.Link) {
						defer slog.Debug("crawling goroutine exited",
							slog.String("url", link.URL),
						)

						slog.Debug("crawling url",
							slog.String("url", link.URL),
						)

						links, err := crawlPage(ctx, link)
						if err != nil {
							slog.Error("error crawling",
								slog.String("url", link.URL),
								log.Err(err),
							)
						}

						parsedPages <- model.NewPage(link.URL, link.Depth, links...)
					}(l)
				}
			}
		}()
	}

	return parsedPages, nil
}

type crawlerFunc func(context.Context, model.Link) ([]model.Link, error)

func buildCrawler(rootURL *url.URL, workers int, sem *semaphore.Weighted, cache *cache.Set[string]) crawlerFunc {
	c := client.New(workers)
	return func(
		ctx context.Context,
		link model.Link,
	) ([]model.Link, error) {
		defer sem.Free()
		sem.Acquire()
		cache.Add(link.URL)

		data, err := fetch.Page(ctx, c, link.URL)
		if err != nil {
			return nil, err
		}

		return parse.Page(ctx, data, rootURL, link.Depth, cache)
	}
}
