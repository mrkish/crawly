package fetch

import (
	"bytes"
	"context"
	"crawly/internal/cache"
	"crawly/internal/model"
	"crawly/internal/parse"
	"crawly/internal/semaphore"
	"io"
	"net/http"

	"golang.org/x/net/html"
)

func Start(
	ctx context.Context,
	root string,
	workers, depth int,
	fetchQueue <-chan model.Parsed,
	parseQueue chan model.Fetched,
	cache *cache.Cache,
) (any, error) {
	fetchSem := semaphore.New(workers)
	parseSem := semaphore.New(workers)
	defer fetchSem.Close()
	defer parseSem.Close()

	go fetch(ctx, root, depth, 1, parseQueue)

	for {
		select {
		case <-ctx.Done():
			return nil, nil
		case parseable := <-parseQueue:
			parseSem.Acquire()
			go func() {
				defer parseSem.Free()
				parse.HTML(parseable, parseQueue, fetchQueue, cache, depth)
			}()
		case req := <-fetchQueue:
			fetchSem.Acquire()
			go func() {
				defer fetchSem.Free()
				fetch(ctx, req.URL, depth, req.Level+1, parseQueue)
			}()
		}
	}
}

func fetch(ctx context.Context, url string, depth, currentLevel int, out chan<- model.Fetched) {
	select {
	case <-ctx.Done():
		return
	default:
	}
	if currentLevel > depth {
		return
	}
	res, err := http.Get(url)
	if err != nil {
		return
	}
	defer res.Body.Close()
	b, err := io.ReadAll(res.Body)
	if err != nil {
		return
	}
	doc, err := html.Parse(bytes.NewBuffer(b))
	if err != nil {
		return
	}
	tok := html.NewTokenizer(bytes.NewBuffer(b))
	// fmt.Printf("html node: %v\n", fetched)
	body := parse.GetBody(doc)
	if body == nil {
		return
	}
	out <- model.Fetched{
		URL:    url,
		Level:  currentLevel,
		Doc:    body,
		Tokens: tok,
	}
}

// func extract(res *http.Response) *html.Node {
// 	defer res.Body.Close()
// 	doc, err := html.Parse(res.Body)
// 	if err != nil {
// 		return nil
// 	}
// 	return doc
// }
