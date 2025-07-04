package app

import (
	"flag"
	"time"
)

type flags struct {
	Depth    int
	Workers  int
	Timeout  time.Duration
	LogLevel string
	Output   string
	URL      string
}

func readFlags() flags {
	depth := flag.Int("d", 3, "Defines how many levels below the root URL should be crawled. Default: 3.")
	timeout := flag.Int("t", 300, "Defines the maximum time in seconds to allow the crawling to continue. Default: 300.")
	logLevel := flag.String("l", "info", "Defines logging output level. Default: info.")
	workers := flag.Int("w", 5, "Defines the number of workers (concurrent) to use. Default: 5.")
	output := flag.String("o", "stdio", "Defines the output format. Default: stdio.")
	url := flag.String("u", "https://www.scrapingcourse.com/ecommerce/", "Defines the root URL to crawl. Required value.")
	flag.Parse()
	return flags{
		Depth:    *depth,
		Timeout:  time.Second * time.Duration(*timeout),
		Workers:  *workers,
		LogLevel: *logLevel,
		Output:   *output,
		URL:      *url,
	}
}
