package cmdutil

import "github.com/spf13/cobra"

// AddOutputFlag adds the --output/-o flag.
func AddOutputFlag(cmd *cobra.Command, output *string) {
	cmd.Flags().StringVarP(output, "output", "o", "", "Output format: table, json, yaml")
}

// AddFileFlag adds the --file/-f flag.
func AddFileFlag(cmd *cobra.Command, file *string) {
	cmd.Flags().StringVarP(file, "file", "f", "", "Path to JSON or YAML file (use - for stdin)")
}

// AddConfirmFlag adds the --confirm flag.
func AddConfirmFlag(cmd *cobra.Command, confirm *bool) {
	cmd.Flags().BoolVar(confirm, "confirm", false, "Skip confirmation prompt")
}

// AddPaginationFlags adds --page and --limit flags.
func AddPaginationFlags(cmd *cobra.Command, page, limit *int) {
	cmd.Flags().IntVar(page, "page", 1, "Page number")
	cmd.Flags().IntVar(limit, "limit", 100, "Number of results per page")
}
