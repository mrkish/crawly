package report

import (
	"log/slog"
	"strings"

	"github.com/mrkish/crawly/internal/model"
)

type Output string

const (
	JSON  Output = "json"
	Stdio Output = "stdio"
	CSV   Output = "csv"
)

func VerifyOutput(output string) Output {
	switch strings.ToLower(output) {
	case "json":
		return JSON
	case "csv":
		return CSV
	case "stdio":
		return Stdio
	default:
		slog.Error("invalid input, defaulting to stdio")
		return Stdio
	}
}

func Out(results []model.Page, output Output) {
	switch output {
	case JSON:
	case Stdio:
		slog.Info("finished crawling",
			"pages", results,
		)
	case CSV:

	}
}
