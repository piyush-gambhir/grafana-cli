package cmd

import (
	"bytes"
	"strings"
	"testing"

	"github.com/piyush-gambhir/grafana-cli/internal/cmdutil"
)

func TestConfirmAction_NoInputBlocks(t *testing.T) {
	var out bytes.Buffer
	in := strings.NewReader("")

	// With noInput=true and confirmed=false, should return error.
	ok, err := cmdutil.ConfirmAction(in, &out, "Delete?", false, true)
	if err == nil {
		t.Fatal("expected error when no-input is set and not confirmed")
	}
	if ok {
		t.Fatal("expected ok=false when no-input blocks")
	}
	if !strings.Contains(err.Error(), "--no-input") {
		t.Errorf("error should mention --no-input, got: %v", err)
	}
}

func TestConfirmAction_NoInputAllowsConfirm(t *testing.T) {
	var out bytes.Buffer
	in := strings.NewReader("")

	// With noInput=true but confirmed=true, should succeed.
	ok, err := cmdutil.ConfirmAction(in, &out, "Delete?", true, true)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !ok {
		t.Fatal("expected ok=true when confirmed=true even with no-input")
	}
}

func TestConfirmAction_NoInputNotSet(t *testing.T) {
	var out bytes.Buffer
	in := strings.NewReader("y\n")

	// Without noInput, should prompt and accept "y".
	ok, err := cmdutil.ConfirmAction(in, &out, "Delete?", false)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !ok {
		t.Fatal("expected ok=true after user types 'y'")
	}
}

func TestLoginCmd_NoInput(t *testing.T) {
	f := &cmdutil.Factory{
		IOStreams: cmdutil.DefaultIOStreams(),
		NoInput:  true,
	}
	cmd := newLoginCmd(f)
	err := cmd.RunE(cmd, nil)
	if err == nil {
		t.Fatal("expected error for login with --no-input")
	}
	if !strings.Contains(err.Error(), "--no-input") {
		t.Errorf("error should mention --no-input, got: %v", err)
	}
}
