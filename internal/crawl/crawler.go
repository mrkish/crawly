package crawl

import (
	"context"
	"log/slog"
	"net/url"
	"strings"
	"sync/atomic"

	"github.com/mrkish/crawly/pkg/cache"
	"github.com/mrkish/crawly/pkg/log"
	"github.com/mrkish/crawly/pkg/semaphore"
)

func FromRoot(ctx context.Context, root string, workers, maxDepth int) ([]Page, error) {
	rootURL, err := url.Parse(root)
	if err != nil {
		return nil, err
	}

	rootPage := NewPage(root, 1)
	cache := cache.NewSet(func(key string) string { return strings.TrimRight(key, slash) })

	linkQueue := make(chan Link)
	defer close(linkQueue)

	parsedPages, err := startCrawl(ctx, rootURL, rootPage, workers, maxDepth, cache, linkQueue)
	if err != nil {
		return nil, err
	}

	var pending uint32 = 1
	var pages []Page

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

			slog.Debug("finished processing page result", slog.Any("pending", pending))

			if pending == 0 {
				slog.Info("crawling completed")
				return pages, nil
			}
		}
	}
}

func startCrawl(
	ctx context.Context,
	rootURL *url.URL,
	rootPage Page,
	workers, maxDepth int,
	cache *cache.Set[string],
	linkQueue <-chan Link,
) (<-chan Page, error) {
	var err error
	parsedPages := make(chan Page)
	sem := semaphore.New(workers)
	crawlPage := buildCrawler(rootURL)

	rootPage.Links, err = crawlPage(ctx, rootPage.Link, sem, cache)
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
			defer slog.Info("crawling goroutine exited")

			for {
				select {
				case <-ctx.Done():
					return
				case l := <-linkQueue:
					go func(link Link) {
						defer slog.Debug("crawling goroutine exited",
							slog.String("url", link.URL),
						)

						slog.Debug("crawling url",
							slog.String("url", link.URL),
						)

						links, err := crawlPage(ctx, link, sem, cache)
						if err != nil {
							slog.Error("error crawling",
								slog.String("url", link.URL),
								log.Err(err),
							)
						}

						parsedPages <- NewPage(link.URL, link.Depth, links...)
					}(l)
				}
			}
		}()
	}

	return parsedPages, nil
}

type crawlerFunc func(context.Context, Link, *semaphore.Weighted, *cache.Set[string]) ([]Link, error)

func buildCrawler(rootURL *url.URL) crawlerFunc {
	return func(
		ctx context.Context,
		link Link,
		sem *semaphore.Weighted,
		cache *cache.Set[string],
	) ([]Link, error) {
		defer sem.Free()
		sem.Acquire()
		cache.Add(link.URL)

		data, err := fetch(ctx, link.URL)
		if err != nil {
			return nil, err
		}

		return parse(ctx, data, rootURL, link.Depth, cache)
	}
}
