package model

type Link struct {
	URL   string `json:"url"`
	Depth int    `json:"depth"`
}

func (l Link) IsInvalid(maxDepth int, hasFn func(string) bool) bool {
	return l.Depth > maxDepth || l.URL == "" || hasFn(l.URL)
}
