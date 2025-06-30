package fetch

import (
	"context"
	"errors"
	"io"
	"log/slog"
	"net/http"
	"time"

	"github.com/mrkish/crawly/internal/client"
	"github.com/mrkish/crawly/pkg/log"
)

func Page(ctx context.Context, c *client.Client, url string) (io.ReadCloser, error) {
	if c == nil {
		return page(ctx, client.New(), url)
	}
	return page(ctx, c, url)
}

func page(ctx context.Context, c *client.Client, url string) (io.ReadCloser, error) {
	slog.Debug("fetching URL", slog.String("url", url))
	if ctx.Err() != nil {
		return nil, ctx.Err()
	}
	start := time.Now()
	res, err := c.Do(ctx, http.MethodGet, url, nil)
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
