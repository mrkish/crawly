package crawl

import (
	"context"
	"errors"
	"io"
	"log/slog"
	"net/http"
	"time"

	"github.com/mrkish/crawly/pkg/log"
)

func fetch(ctx context.Context, url string) (io.ReadCloser, error) {
	slog.Debug("fetching URL", slog.String("url", url))
	if ctx.Err() != nil {
		return nil, ctx.Err()
	}
	var res *http.Response
	start := time.Now()
	res, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	slog.Debug("response status",
		log.Duration(start),
		slog.String("url", url),
		slog.Int("response code", res.StatusCode),
	)
	switch res.StatusCode {
	case http.StatusOK:
		return res.Body, err
	default:
		return nil, errors.New("failed to pull page")
	}
}
