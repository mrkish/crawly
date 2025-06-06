package app

import "time"

type BuildInfo struct {
	BuildTime time.Time
	Commit    string
	Version   string
}
