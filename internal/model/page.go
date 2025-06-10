package model

type Page struct {
	Link
	Links []Link
}

func NewPage(url string, depth int, links ...Link) Page {
	return Page{
		Link: Link{
			URL:   url,
			Depth: depth,
		},
		Links: links,
	}
}
