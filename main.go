package main

import (
	"os"

	"github.com/piyush-gambhir/grafana-cli/cmd"
	"github.com/piyush-gambhir/grafana-cli/internal/client"
	"github.com/piyush-gambhir/grafana-cli/internal/output"
)

func main() {
	if err := cmd.Execute(); err != nil {
		// Determine output format from flags/env.
		format := os.Getenv("GRAFANA_OUTPUT")
		for i, arg := range os.Args {
			if (arg == "-o" || arg == "--output") && i+1 < len(os.Args) {
				format = os.Args[i+1]
				break
			}
		}

		statusCode := 0
		if apiErr, ok := err.(*client.APIError); ok {
			statusCode = apiErr.StatusCode
		}

		output.WriteError(os.Stderr, format, err, statusCode)
		os.Exit(1)
	}
}
