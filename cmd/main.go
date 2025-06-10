package main

import (
	"fmt"
	"os"
	"time"

	"github.com/mrkish/crawly/internal/app"
)

var (
	buildTime time.Time
	commit    string
	version   string
)

func main() {
	if err := app.Run(app.BuildInfo{
		BuildTime: buildTime,
		Commit:    commit,
		Version:   version,
	}); err != nil {
		fmt.Printf("error occurred: %s\n", err.Error())
		os.Exit(1)
	}
}
