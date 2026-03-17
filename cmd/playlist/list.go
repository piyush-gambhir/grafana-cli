package playlist

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"

	"github.com/piyush-gambhir/grafana-cli/internal/client"
	"github.com/piyush-gambhir/grafana-cli/internal/cmdutil"
	"github.com/piyush-gambhir/grafana-cli/internal/output"
)

func newCmdPlaylistList(f *cmdutil.Factory) *cobra.Command {
	var query string

	cmd := &cobra.Command{
		Use:     "list",
		Short:   "List playlists",
		Aliases: []string{"ls"},
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := f.Client()
			if err != nil {
				return err
			}

			results, err := c.ListPlaylists(context.Background(), query)
			if err != nil {
				return err
			}

			if len(results) == 0 {
				fmt.Fprintln(f.IOStreams.Out, "No playlists found.")
				return nil
			}

			return output.Print(f.IOStreams.Out, f.Resolved.Output, results, &output.TableDef{
				Headers: []string{"UID", "Name", "Interval"},
				RowFunc: func(item interface{}) []string {
					p := item.(client.Playlist)
					return []string{p.UID, p.Name, p.Interval}
				},
			})
		},
	}

	cmd.Flags().StringVarP(&query, "query", "q", "", "Search query")

	return cmd
}
