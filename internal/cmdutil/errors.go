package cmdutil

import (
	"fmt"
	"io"

	"github.com/piyush-gambhir/grafana-cli/internal/client"
)

// FlagErrorf creates a formatted error for flag-related issues.
func FlagErrorf(format string, args ...interface{}) error {
	return fmt.Errorf(format, args...)
}

// PrintError writes a formatted error message to the given writer.
func PrintError(w io.Writer, err error) {
	if apiErr, ok := err.(*client.APIError); ok {
		fmt.Fprintf(w, "Error: %s (HTTP %d)\n", apiErr.Message, apiErr.StatusCode)
		return
	}
	fmt.Fprintf(w, "Error: %s\n", err)
}

// ConfirmAction prompts the user for confirmation if not auto-confirmed.
func ConfirmAction(in io.Reader, out io.Writer, message string, confirmed bool) (bool, error) {
	if confirmed {
		return true, nil
	}

	fmt.Fprintf(out, "%s [y/N]: ", message)
	var response string
	_, err := fmt.Fscan(in, &response)
	if err != nil {
		return false, nil
	}

	return response == "y" || response == "Y" || response == "yes" || response == "Yes", nil
}
