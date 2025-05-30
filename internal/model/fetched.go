package model

import (
	"net/url"

	"golang.org/x/net/html"
)

type Fetched struct {
	Doc *html.Node
	Level int
	URL string
	Root url.URL
}
