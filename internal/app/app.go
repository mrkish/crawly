package app

import (
	"context"
	"errors"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/mrkish/crawly/internal/crawl"
	"github.com/mrkish/crawly/pkg/log"
)

var quit = func() chan os.Signal {
	q := make(chan os.Signal, 1)
	signal.Notify(q,
		os.Interrupt,
		syscall.SIGTERM,
		syscall.SIGHUP,
		syscall.SIGINT,
		syscall.SIGQUIT)
	return q
}()

func Run(info BuildInfo) error {
	flags := readFlags()
	if flags.URL == "" {
		return errors.New("no root URL specified")
	}

	log.Init(flags.LogLevel, info.Commit, info.Version)
	slog.Info("initiating crawling",
		slog.Any("flags", flags),
	)

	ctx, cancel := context.WithTimeout(context.Background(), flags.Timeout)
	defer cancel()

	type result struct {
		pages []crawl.Page
		err   error
	}
	resultChan := make(chan result)

	start := time.Now()

	go func() {
		defer close(resultChan)
		links, err := crawl.FromRoot(ctx, flags.URL, flags.Workers, flags.Depth)
		resultChan <- result{links, err}
	}()

	for {
		select {
		case <-ctx.Done():
			return errors.New("timed out attempting to crawl")
		case <-quit:
			slog.Error("received request to cancel, stopping")
			return nil
		case result := <-resultChan:
			slog.Info("finished crawling",
				log.Duration(start),
				"pages", result.pages,
			)
			if result.err != nil {
				slog.Error("error from crawlin", log.Err(result.err))
			}
			return nil
		}
	}
}
