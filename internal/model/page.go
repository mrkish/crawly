package model

type Page struct {
	Link  `json:"page"`
	Links []Link `json:"links,omitempty"`
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
