package fetch

import (
	"context"
	"crawly/internal/cache"
	"crawly/internal/model"
	"crawly/internal/parse"
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
	fetch(root, depth, 1, parseQueue)
	select {
	case <-ctx.Done():
		return nil, nil
	case parseable := <- parseQueue:
		go parse.HTML(parseable, fetchQueue, cache)
	case req := <- fetchQueue:
		// TODO: Add semaphore here
		go fetch(req.URL, depth, req.Level + 1, parseQueue)
	}
	return nil, nil
}

func fetch(url string, depth, currentLevel int, out chan<- model.Fetched) {
	if currentLevel > depth {
		return
	}
	res, err := http.Get(url)
	if err != nil {
		return
	}
	fetched := extract(res)
	out <- model.Fetched{
		URL: url,
		Level: currentLevel,
		Doc: fetched,
	}
}

func extract(res *http.Response) *html.Node {
	defer res.Body.Close()
	doc, err := html.Parse(res.Body)
	if err != nil {
		return nil
	}
	return doc
}
