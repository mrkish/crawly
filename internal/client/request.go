package client

import (
	"io"
	"log/slog"
	"net/http"

	"github.com/mrkish/crawly/pkg/log"
)

type Response struct {
	Body       []byte
	Header     http.Header
	Request    *http.Request
	Status     string
	StatusCode int
}

func toResponse(res *http.Response) Response {
	var body []byte
	if res.Body != nil {
		defer res.Body.Close()
		var err error
		body, err = io.ReadAll(res.Body)
		if err != nil {
			slog.Error("error reading request body", log.Err(err))
		}
	}
	return Response{
		Body:       body,
		Header:     res.Header,
		Request:    res.Request,
		Status:     res.Status,
		StatusCode: res.StatusCode,
	}
}
