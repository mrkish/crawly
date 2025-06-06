package app

import (
	"context"
	"crawly/internal/crawl"
	"crawly/pkg/log"
	"errors"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"
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
	slog.Debug("setting up for crawling", slog.Any("flags", flags))

	ctx, cancel := context.WithTimeout(context.Background(), flags.Timeout)
	defer cancel()

	// cache := cache.New()
	log.Init(flags.LogLevel)

	type result struct {
		pages []crawl.Page
		err   error
	}
	resultChan := make(chan result)

	start := time.Now()

	go func(out chan<- result) {
		defer close(resultChan)
		links, err := crawl.FromRoot(ctx, flags.URL, flags.Depth)
		out <- result{links, err}
	}(resultChan)

	for {
		select {
		case <-ctx.Done():
			cancel()
			return errors.New("timed out attempting to crawl")
		case <-quit:
			slog.Error("received request to cancel, stopping")
			cancel()
			return nil
		case result := <-resultChan:
			slog.Info("finished crawling",
				log.Duration(start),
				"pages", result.pages,
				log.Err(result.err),
			)
			return nil
		}
	}
}
