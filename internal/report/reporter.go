package report

import (
	"log/slog"
	"strings"

	"github.com/mrkish/crawly/internal/model"
)

type Output string

const (
	JSON Output = "json"
	CSV  Output = "csv"
)

func VerifyOutput(output string) Output {
	switch strings.ToLower(output) {
	case "csv":
		return CSV
	default:
		slog.Error("invalid input, defaulting to JSON")
		return JSON
	}
}

func Out(results []model.Page, output Output) {
	switch output {
	case JSON:
		slog.Info("finished crawling",
			"pages", results,
		)
	case CSV:

	}
}
