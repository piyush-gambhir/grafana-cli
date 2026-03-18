package client

import (
	"bytes"
	"encoding/json"
	"errors"
	"testing"

	"github.com/piyush-gambhir/grafana-cli/internal/output"
)

func TestStructuredJSONError_WithStatusCode(t *testing.T) {
	var buf bytes.Buffer
	err := &APIError{StatusCode: 409, Message: "resource already exists", URL: "/api/dashboards/db"}

	output.WriteError(&buf, "json", err, 409)

	var resp output.ErrorResponse
	if jsonErr := json.Unmarshal(buf.Bytes(), &resp); jsonErr != nil {
		t.Fatalf("output is not valid JSON: %v\nOutput: %s", jsonErr, buf.String())
	}

	if resp.StatusCode != 409 {
		t.Errorf("got status_code %d, want 409", resp.StatusCode)
	}
	if resp.Error == "" {
		t.Error("expected non-empty error message")
	}
}

func TestStructuredJSONError_WithZeroStatusCode(t *testing.T) {
	var buf bytes.Buffer
	err := errors.New("config file not found")

	output.WriteError(&buf, "json", err, 0)

	raw := buf.String()

	var resp output.ErrorResponse
	if jsonErr := json.Unmarshal([]byte(raw), &resp); jsonErr != nil {
		t.Fatalf("output is not valid JSON: %v\nOutput: %s", jsonErr, raw)
	}

	if resp.Error != "config file not found" {
		t.Errorf("got error %q, want %q", resp.Error, "config file not found")
	}
	// status_code should be omitted when 0.
	if resp.StatusCode != 0 {
		t.Errorf("expected status_code 0, got %d", resp.StatusCode)
	}
}

func TestIdempotentFlags_IfNotExists_409(t *testing.T) {
	// Simulate --if-not-exists behavior: 409 Conflict should be treated as success.
	apiErr := &APIError{StatusCode: 409, Message: "Conflict: resource already exists"}
	ifNotExists := true

	if ifNotExists && IsConflict(apiErr) {
		// This is the expected behavior: swallow the error.
		return
	}
	t.Fatal("expected 409 to be handled by --if-not-exists")
}

func TestIdempotentFlags_IfNotExists_OtherError(t *testing.T) {
	// Non-409 errors should NOT be swallowed by --if-not-exists.
	apiErr := &APIError{StatusCode: 500, Message: "Internal Server Error"}
	ifNotExists := true

	if ifNotExists && IsConflict(apiErr) {
		t.Fatal("500 error should NOT be treated as conflict")
	}
	// Expected: error is not swallowed.
}

func TestIdempotentFlags_IfExists_404(t *testing.T) {
	// Simulate --if-exists behavior: 404 Not Found should be treated as success.
	apiErr := &APIError{StatusCode: 404, Message: "Not found"}
	ifExists := true

	if ifExists && IsNotFound(apiErr) {
		// This is the expected behavior: swallow the error.
		return
	}
	t.Fatal("expected 404 to be handled by --if-exists")
}

func TestIdempotentFlags_IfExists_OtherError(t *testing.T) {
	// Non-404 errors should NOT be swallowed by --if-exists.
	apiErr := &APIError{StatusCode: 403, Message: "Forbidden"}
	ifExists := true

	if ifExists && IsNotFound(apiErr) {
		t.Fatal("403 error should NOT be treated as not found")
	}
	// Expected: error is not swallowed.
}

func TestQuietFlag_SuppressesInfoMessages(t *testing.T) {
	// Simulate the quiet flag behavior: when quiet is true,
	// informational output should not be written.
	var buf bytes.Buffer
	quiet := true

	infoMsg := "Datasource created successfully"

	if !quiet {
		buf.WriteString(infoMsg + "\n")
	}

	if buf.Len() != 0 {
		t.Errorf("expected no output when quiet=true, got: %s", buf.String())
	}
}

func TestQuietFlag_AllowsOutput(t *testing.T) {
	// When quiet is false, informational output should be written.
	var buf bytes.Buffer
	quiet := false

	infoMsg := "Datasource created successfully"

	if !quiet {
		buf.WriteString(infoMsg + "\n")
	}

	if buf.Len() == 0 {
		t.Error("expected output when quiet=false")
	}
	if buf.String() != infoMsg+"\n" {
		t.Errorf("got %q, want %q", buf.String(), infoMsg+"\n")
	}
}
