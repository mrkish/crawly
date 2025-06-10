package model

type Link struct {
	URL   string
	Depth int
}

func (l Link) IsInvalid(maxDepth int, hasFn func(string) bool) bool {
	return l.Depth > maxDepth || l.URL == "" || hasFn(l.URL)
}
