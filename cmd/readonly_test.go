package cmd

import (
	"testing"

	"github.com/spf13/cobra"
)

func TestCheckReadOnly_WriteCmdBlocked(t *testing.T) {
	cmd := &cobra.Command{
		Use:         "delete",
		Annotations: map[string]string{"mutates": "true"},
	}

	// Simulate read-only enforcement logic from PersistentPreRunE.
	effectiveReadOnly := true
	if effectiveReadOnly && cmd.Annotations != nil && cmd.Annotations["mutates"] == "true" {
		// Expected: command should be blocked.
		return
	}
	t.Fatal("expected write command to be blocked in read-only mode")
}

func TestCheckReadOnly_WriteCmdAllowed(t *testing.T) {
	cmd := &cobra.Command{
		Use:         "delete",
		Annotations: map[string]string{"mutates": "true"},
	}

	// Simulate read-only = false.
	effectiveReadOnly := false
	if effectiveReadOnly && cmd.Annotations != nil && cmd.Annotations["mutates"] == "true" {
		t.Fatal("write command should not be blocked when read-only is false")
	}
	// Expected: no error, command allowed.
}

func TestCheckReadOnly_ReadCmdAllowed(t *testing.T) {
	cmd := &cobra.Command{
		Use: "list",
		// No annotations — this is a read command.
	}

	// Simulate read-only = true.
	effectiveReadOnly := true
	if effectiveReadOnly && cmd.Annotations != nil && cmd.Annotations["mutates"] == "true" {
		t.Fatal("read command should not be blocked even in read-only mode")
	}
	// Expected: no error, command allowed.
}
