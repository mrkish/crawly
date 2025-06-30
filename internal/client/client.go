package client

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/url"

	"github.com/mrkish/crawly/pkg/cache"
	"github.com/mrkish/crawly/pkg/limiter"
)

const (
	robotsPath = "/robots.txt"
)

type Client struct {
	client http.Client

	disallowed *cache.Hashed[string, map[string]struct{}]
	limit      *limiter.Request
}

func New(connections ...int) *Client {
	var limit *limiter.Request
	if len(connections) > 0 {
		limit = limiter.New(connections[0])
	}
	return &Client{
		client:     http.Client{},
		disallowed: cache.NewHashed[string, map[string]struct{}](nil, nil),
		limit:      limit,
	}
}

func (c *Client) Do(ctx context.Context, method string, url string, payload any) (resp *http.Response, err error) {
	var body io.Reader
	if payload != nil {
		var b []byte
		b, err = json.Marshal(payload)
		if err != nil {
			return
		}
		body = bytes.NewBuffer(b)
	}

	if list := c.getDisallowed(url); list != nil {
		if _, ok := list[url]; ok {
			return nil, errors.New("provided URL is in site's robots disallow list")
		}
	}

	var req *http.Request
	req, err = http.NewRequestWithContext(ctx, method, url, body)
	if err != nil {
		return
	}

	if c.limit != nil {
		c.limit.Request(url, func() {
			resp, err = c.client.Do(req)
		})
	} else {
		resp, err = c.client.Do(req)
	}

	return
}

func (c *Client) getDisallowed(path string) map[string]struct{} {
	if list := c.disallowed.Get(path); list != nil {
		return list
	}

	disallowed := make(map[string]struct{})

	uri, err := url.Parse(path)
	if err != nil {
		return disallowed
	}

	robotsURL := uri.Host + robotsPath

	req, err := http.NewRequestWithContext(context.Background(), http.MethodGet, robotsURL, nil)
	if err != nil {
		return disallowed
	}

	var resp *http.Response
	resp, err = c.client.Do(req)
	if err != nil {
		return disallowed
	}

	switch resp.StatusCode {
	case http.StatusOK:
		c.disallowed.Add(path, disallowed)
	default:
	}

	return disallowed
}

func parseRobotsTxt(body io.ReadCloser, disallowed map[string]struct{}) {
	defer body.Close()

}
