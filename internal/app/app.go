package app

import (
	"context"
	"crawly/internal/cache"
	"crawly/internal/fetch"
	"crawly/internal/model"
	"errors"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"
)

type BuildInfo struct {
	BuildTime time.Time
	Commit    string
	Version   string
}

type Flags struct {
	Depth    int
	Workers  int
	Timeout  int
	LogLevel string
	Output   string
	URL      string
}

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

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*time.Duration(flags.Timeout))
	defer cancel()

	parseQueue := make(chan model.Fetched)
	fetchQueue := make(chan model.Parsed)

	cache := cache.New()

	fmt.Printf("got flags: %v\n", flags)
	go func() {
		_, _ = fetch.Start(ctx, flags.URL, flags.Workers, flags.Depth, fetchQueue, parseQueue, cache)
	}()

	<- quit
	cancel()
	

	// report.Output(result)

	return nil
}

// Exit sends an interrupt signal to end the program
func Exit() {
	quit <- os.Interrupt
}

func readFlags() Flags {
	depth := flag.Int("d", 3, "Defines how many levels below the root URL should be crawled. Default: 3.")
	timeout := flag.Int("t", 300, "Defines the maximum time in seconds to allow the crawling to continue. Default: 300.")
	logLevel := flag.String("l", "debug", "Defines logging output level. Default: error.")
	workers := flag.Int("w", 5, "Defines the number of workers (concurrent) to use. Default: 5.")
	output := flag.String("o", "json", "Defines the output format. Default: JSON.")
	url := flag.String("u", "https://www.scrapingcourse.com/ecommerce/", "Defines the root URL to crawl. Required value.")
	flag.Parse()
	return Flags{
		Depth:    *depth,
		Timeout:  *timeout,
		Workers:  *workers,
		LogLevel: *logLevel,
		Output:   *output,
		URL:      *url,
	}
}
