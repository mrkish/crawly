commit := $(shell git rev-parse --short HEAD)
build:
	go build -o crawly -ldflags "-X 'main.version=1.0.0' -X 'main.commit=$(commit)'" cmd/main.go
