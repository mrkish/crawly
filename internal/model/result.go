package model

type Result struct {
	URL string
	Level int
	Connections []*Result
}
