package crawl

type Link struct {
	URL   string
	Depth int
}

func (l Link) IsInvalid(maxDepth int, hasFn func(string) bool) bool {
	return l.Depth > maxDepth || l.URL == "" || hasFn(l.URL)
}

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
