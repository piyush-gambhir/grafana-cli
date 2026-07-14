package output

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"
)

func TestJSONFormatter_Output(t *testing.T) {
	data := map[string]interface{}{
		"id":   1,
		"name": "test",
	}

	var buf bytes.Buffer
	f := &JSONFormatter{}
	err := f.Format(&buf, data, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	output := strings.TrimSpace(buf.String())

	// Should be valid JSON.
	var parsed map[string]interface{}
	if err := json.Unmarshal([]byte(output), &parsed); err != nil {
		t.Fatalf("output is not valid JSON: %v\nOutput: %s", err, output)
	}

	if parsed["name"] != "test" {
		t.Errorf("got name %v, want %q", parsed["name"], "test")
	}
}

func TestJSONFormatter_PrettyPrint(t *testing.T) {
	data := map[string]string{"key": "value"}

	var buf bytes.Buffer
	f := &JSONFormatter{}
	err := f.Format(&buf, data, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	output := buf.String()
	// Pretty-printed JSON should contain newlines and indentation.
	if !strings.Contains(output, "\n") {
		t.Error("expected pretty-printed JSON with newlines")
	}
	if !strings.Contains(output, "  ") {
		t.Error("expected indentation in pretty-printed JSON")
	}
}

func TestYAMLFormatter_Output(t *testing.T) {
	data := map[string]interface{}{
		"name":  "test",
		"count": 42,
	}

	var buf bytes.Buffer
	f := &YAMLFormatter{}
	err := f.Format(&buf, data, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	output := buf.String()
	if !strings.Contains(output, "name: test") {
		t.Errorf("expected YAML key-value, got: %s", output)
	}
	if !strings.Contains(output, "count: 42") {
		t.Errorf("expected YAML count, got: %s", output)
	}
}

func TestTableFormatter_Headers(t *testing.T) {
	type Item struct {
		Name string
		Age  int
	}

	data := []Item{
		{Name: "Alice", Age: 30},
		{Name: "Bob", Age: 25},
	}

	tableDef := &TableDef{
		Headers: []string{"Name", "Age"},
		RowFunc: func(item interface{}) []string {
			i := item.(Item)
			return []string{i.Name, strings.TrimSpace(strings.Repeat(" ", 0) + string(rune('0'+i.Age/10)) + string(rune('0'+i.Age%10)))}
		},
	}

	var buf bytes.Buffer
	f := &TableFormatter{}
	err := f.Format(&buf, data, tableDef)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	output := buf.String()
	lines := strings.Split(strings.TrimSpace(output), "\n")
	if len(lines) < 3 {
		t.Fatalf("expected at least 3 lines (header + 2 rows), got %d:\n%s", len(lines), output)
	}

	// Headers should be uppercase.
	if !strings.Contains(lines[0], "NAME") {
		t.Errorf("expected NAME in header, got: %s", lines[0])
	}
	if !strings.Contains(lines[0], "AGE") {
		t.Errorf("expected AGE in header, got: %s", lines[0])
	}

	// Check data rows.
	if !strings.Contains(lines[1], "Alice") {
		t.Errorf("expected Alice in first row, got: %s", lines[1])
	}
	if !strings.Contains(lines[2], "Bob") {
		t.Errorf("expected Bob in second row, got: %s", lines[2])
	}
}

func TestNewFormatter_Default(t *testing.T) {
	tests := []struct {
		format  string
		wantErr bool
	}{
		{format: "table", wantErr: false},
		{format: "json", wantErr: false},
		{format: "yaml", wantErr: false},
		{format: "csv", wantErr: true},
		{format: "xml", wantErr: true},
		{format: "", wantErr: true},
	}

	for _, tt := range tests {
		t.Run(tt.format, func(t *testing.T) {
			f, err := NewFormatter(tt.format)
			if tt.wantErr {
				if err == nil {
					t.Errorf("expected error for format %q", tt.format)
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error for format %q: %v", tt.format, err)
			}
			if f == nil {
				t.Errorf("expected non-nil formatter for format %q", tt.format)
			}
		})
	}
}
