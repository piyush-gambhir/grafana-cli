package policy

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"

	"github.com/piyush-gambhir/grafana-cli/internal/cmdutil"
	"github.com/piyush-gambhir/grafana-cli/internal/output"
)

func newCmdPolicyGet(f *cmdutil.Factory) *cobra.Command {
	return &cobra.Command{
		Use:   "get",
		Short: "Get the notification policy tree",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := f.Client()
			if err != nil {
				return err
			}

			policy, err := c.GetNotificationPolicy(context.Background())
			if err != nil {
				return err
			}

			if f.Resolved.Output == "table" {
				fmt.Fprintf(f.IOStreams.Out, "Receiver:        %s\n", policy.Receiver)
				fmt.Fprintf(f.IOStreams.Out, "Group By:        %v\n", policy.GroupBy)
				fmt.Fprintf(f.IOStreams.Out, "Group Wait:      %s\n", policy.GroupWait)
				fmt.Fprintf(f.IOStreams.Out, "Group Interval:  %s\n", policy.GroupInterval)
				fmt.Fprintf(f.IOStreams.Out, "Repeat Interval: %s\n", policy.RepeatInterval)
				fmt.Fprintf(f.IOStreams.Out, "Routes:          %d\n", len(policy.Routes))
				return nil
			}

			return output.Print(f.IOStreams.Out, f.Resolved.Output, policy, nil)
		},
	}
}
