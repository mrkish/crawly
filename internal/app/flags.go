package app

import (
	"flag"
	"time"
)

type Flags struct {
	Depth    int
	Workers  int
	Timeout  time.Duration
	LogLevel string
	Output   string
	URL      string
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
		Timeout:  time.Millisecond * time.Duration(*timeout),
		Workers:  *workers,
		LogLevel: *logLevel,
		Output:   *output,
		URL:      *url,
	}
}
