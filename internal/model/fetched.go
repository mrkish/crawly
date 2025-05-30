package model

import "golang.org/x/net/html"

type Fetched struct {
	Doc *html.Node
	Level int
	URL string
}
