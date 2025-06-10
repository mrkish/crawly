# Crawly
This is a toy project to build a simple web crawler.

## Usage
The CLI app accepts a few flags to govern its behavior:

  -u
    Root URL to crawl. Required.
    Default: https://www.scrapingcourse.com/ecommerce/

  -o
    Output format. TODO!

  -t
    Timeout in seconds.
    Default: 300

  -w
    Number of workers (concurrent requests) to allow.
    Default: 5

  -d
    Depth -- how far from the root URL to crawl.
    Default: 3

  -l
    Log level. Levels are:
      - trace
      - debug
      - info
      - warn
      - error
    Default: info
