package model

import "sync"

type link int

const (
	linkTo link = iota
	linkFrom
)

type Link struct {
	URL   string
	Level int

	mu        sync.Mutex
	LinksTo   map[string]*Link
	LinksFrom map[string]*Link
}

func New(url string, level int) *Link {
	return &Link{
		URL:       url,
		Level:     level,
		mu:        sync.Mutex{},
		LinksTo:   make(map[string]*Link),
		LinksFrom: make(map[string]*Link),
	}
}

func (l *Link) AddLinkTo(link *Link) {
	l.add(link, linkFrom)
}

func (l *Link) AddLinkedFrom(link *Link) {
	l.add(link, linkFrom)
}

func (l *Link) add(link *Link, kind link) {
	l.mu.Lock()
	defer l.mu.Unlock()
	switch kind {
	case linkTo:
		l.LinksTo[link.URL] = link
	case linkFrom:
		l.LinksFrom[link.URL] = link
	}
}
