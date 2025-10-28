package report

import (
	"fmt"
	"log/slog"
	"strings"

	"github.com/mrkish/crawly/internal/model"
)

type Type string

const (
	JSON Type = "json"
	CSV  Type = "csv"
)

func VerifyOutput(output string) Type {
	switch strings.ToLower(output) {
	case "csv":
		return CSV
	default:
		slog.Error("invalid input, defaulting to JSON")
		return JSON
	}
}

// Out controls outputing the result into the output channel
func Out(results []model.Page, outputType Type) {
	for _, r := range results {
		fmt.Println(r)
	}
	switch outputType {
	case JSON:
		slog.Info("finished crawling",
			"pages", results,
		)
	case CSV:

	}
}
