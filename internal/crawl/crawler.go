package crawl

import (
	"context"
	"crawly/pkg/cache"
	"crawly/pkg/log"
	"log/slog"
	"net/url"
	"strings"
	"sync"
	"sync/atomic"
)

const slash = "/"

type crawledPage struct {
	links []Link
	err   error
}

func FromRoot(ctx context.Context, root string, workers, maxDepth int) ([]Page, error) {
	rootURL, err := url.Parse(root)
	if err != nil {
		return nil, err
	}
	pageBuilder := pageFactory(rootURL)
	rootPage := pageBuilder(root, 1)

	pageChan, err := crawl(ctx, rootPage, workers, maxDepth, pageBuilder)
	if err != nil {
		return nil, err
	}

	var pages []Page
	for page := range pageChan {
		page.Root = rootURL
		pages = append(pages, page)
	}

	return pages, nil
}

func crawl(ctx context.Context, rootPage Page, workers, maxDepth int, pageBuilder PageBuilder) (<-chan Page, error) {
	var err error
	out := make(chan Page)
	cache := cache.New[string, Link](cacheKeyCleaner, mergeMap)

	// the initial page does not need to be parsed concurently
	rootPage.Links, err = crawlPage(ctx, rootPage, cache)
	if err != nil {
		return nil, err
	}

	foundLinks := len(rootPage.Links) > 0

	if foundLinks {
		// send result from crawling root page
		defer func() {
			out <- rootPage
			// close channel so FromRoot will exit when only crawling root
			if maxDepth == 1 {
				close(out)
			}
		}()
	}

	if maxDepth > 1 && foundLinks {
		linkQueue := make(chan []Link)
		// launch routine that will manage concurrent crawling
		go func(queue chan []Link) {
			defer close(out)
			defer slog.Info("crawling function exited")

			var pending uint32
			// sem := semaphore.New(workers)
			defer close(queue)

			for {
				select {
				case <-ctx.Done():
					return
				case links := <-queue:
					go func() {
						for {
						select {
						case links := <- fanIn(ctx, fanOut(
								ctx,
								links,
								out,
								&pending,
								maxDepth,
								pageBuilder,
								cache)...),
							}
						}
					}()
					if pending == 0 {
						slog.Info("queue cleared, done crawling")
						return
					}
				}
			}
		}(linkQueue)

		// seed the processing queue
		linkQueue <- rootPage.Links
	}

	return out, nil
}

func fanOut(
	ctx context.Context,
	links []Link,
	pageResult chan<- Page,
	pending *uint32,
	maxDepth int,
	pageBuilder PageBuilder,
	// sem *semaphore.Weighted,
	cache *cache.Cache[string, Link],
) []chan []Link {
	// add the new links to the counter and then start crawling them
	var channels []chan []Link
	for _, l := range links {
		if l.Depth > maxDepth {
			slog.Debug("page is beyond max depth", slog.Int("url", l.Depth), slog.Int("maxDepth", maxDepth))
			continue
		}
		// Do not re-crawl
		if cache.Has(l.URL) {
			slog.Debug("page already crawled, skipping", slog.String("url", l.URL))
			continue
		}
		// sem.Acquire()
		atomic.SwapUint32(pending, *pending+1)
		linkChan := make(chan []Link)
		channels = append(channels, linkChan)
		go func(link Link) {
			// defer sem.Free()
			defer close(linkChan)
			defer func() { atomic.SwapUint32(pending, *pending-1) }()
			foundLinks, crawlErr := crawlPage(ctx, pageBuilder(link.URL, link.Depth), cache)
			if crawlErr != nil {
				slog.Error("error crawling", log.Err(crawlErr))
			}
			pageResult <- Page{Link: link, Links: links}
			linkChan <- foundLinks
		}(l)
	}
	return channels
}

func fanIn(ctx context.Context, channels ...<-chan []Link) <-chan []Link {
	var wg sync.WaitGroup
	out := make(chan []Link)

	multiplex := func(c <-chan []Link) {
		defer wg.Done()
		for link := range c {
			select {
			case <-ctx.Done():
				return
			case out <- link:
			}
		}
	}

	wg.Add(len(channels))
	for _, c := range channels {
		go multiplex(c)
	}

	go func() {
		wg.Wait()
		close(out)
	}()

	return out
}

func crawlPage(
	ctx context.Context,
	page Page,
	cache *cache.Cache[string, Link],
) ([]Link, error) {
	pageData, err := fetch(ctx, page.URL)
	if err != nil {
		return nil, err
	}
	cache.Add(page.URL, Link{URL: page.URL, Depth: page.Depth})
	return parse(ctx, page.Root, pageData, page.Depth, cache)
}

func cacheKeyCleaner(key string) string {
	return strings.TrimRight(key, slash)
}
