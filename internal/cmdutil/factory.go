package cmdutil

import (
	"io"
	"os"

	"github.com/piyush-gambhir/grafana-cli/internal/client"
	"github.com/piyush-gambhir/grafana-cli/internal/config"
)

// Factory provides shared dependencies to all commands.
type Factory struct {
	Config   func() (*config.Config, error)
	Client   func() (*client.Client, error)
	IOStreams IOStreams
	// Resolved holds the resolved config after PersistentPreRunE.
	Resolved *config.ResolvedConfig
}

// IOStreams holds standard I/O streams.
type IOStreams struct {
	In     io.Reader
	Out    io.Writer
	ErrOut io.Writer
}

// DefaultIOStreams returns IOStreams connected to os.Stdin, os.Stdout, os.Stderr.
func DefaultIOStreams() IOStreams {
	return IOStreams{
		In:     os.Stdin,
		Out:    os.Stdout,
		ErrOut: os.Stderr,
	}
}
