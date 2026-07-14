package output

import (
	"bytes"
	"encoding/json"
	"errors"
	"strings"
	"testing"
)

func TestWriteError_JSON(t *testing.T) {
	var buf bytes.Buffer
	err := errors.New("something went wrong")

	WriteError(&buf, "json", err, 500)

	output := buf.String()

	var resp ErrorResponse
	if err := json.Unmarshal([]byte(output), &resp); err != nil {
		t.Fatalf("output is not valid JSON: %v\nOutput: %s", err, output)
	}

	if resp.Error != "something went wrong" {
		t.Errorf("got error %q, want %q", resp.Error, "something went wrong")
	}
	if resp.StatusCode != 500 {
		t.Errorf("got status_code %d, want %d", resp.StatusCode, 500)
	}
}

func TestWriteError_JSON_NoStatusCode(t *testing.T) {
	var buf bytes.Buffer
	err := errors.New("config error")

	WriteError(&buf, "json", err, 0)

	output := buf.String()

	var resp ErrorResponse
	if err := json.Unmarshal([]byte(output), &resp); err != nil {
		t.Fatalf("output is not valid JSON: %v\nOutput: %s", err, output)
	}

	if resp.Error != "config error" {
		t.Errorf("got error %q, want %q", resp.Error, "config error")
	}
	// status_code should be omitted (zero value) from JSON.
	if strings.Contains(output, "status_code") {
		t.Errorf("expected status_code to be omitted when 0, got: %s", output)
	}
}

func TestWriteError_PlainText(t *testing.T) {
	var buf bytes.Buffer
	err := errors.New("something went wrong")

	WriteError(&buf, "table", err, 404)

	output := buf.String()
	if !strings.Contains(output, "Error: something went wrong") {
		t.Errorf("expected plain text error, got: %s", output)
	}
	// Should NOT contain JSON.
	if strings.HasPrefix(strings.TrimSpace(output), "{") {
		t.Errorf("expected plain text, not JSON: %s", output)
	}
}

func TestWriteError_EmptyFormat(t *testing.T) {
	var buf bytes.Buffer
	err := errors.New("unknown error")

	// Empty format should use plain text.
	WriteError(&buf, "", err, 0)

	output := buf.String()
	if !strings.Contains(output, "Error: unknown error") {
		t.Errorf("expected plain text error for empty format, got: %s", output)
	}
}
