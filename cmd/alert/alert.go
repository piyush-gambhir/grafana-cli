package alert

import (
	"github.com/spf13/cobra"

	"github.com/piyush-gambhir/grafana-cli/cmd/alert/contactpoint"
	"github.com/piyush-gambhir/grafana-cli/cmd/alert/mutetiming"
	"github.com/piyush-gambhir/grafana-cli/cmd/alert/policy"
	"github.com/piyush-gambhir/grafana-cli/cmd/alert/rule"
	"github.com/piyush-gambhir/grafana-cli/cmd/alert/silence"
	"github.com/piyush-gambhir/grafana-cli/cmd/alert/template"
	"github.com/piyush-gambhir/grafana-cli/internal/cmdutil"
)

// NewCmdAlert returns the alert parent command.
func NewCmdAlert(f *cmdutil.Factory) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "alert",
		Short: "Manage alerting resources",
		Long:  "Manage Grafana alerting rules, contact points, policies, mute timings, templates, and silences.",
	}

	cmd.AddCommand(rule.NewCmdRule(f))
	cmd.AddCommand(contactpoint.NewCmdContactPoint(f))
	cmd.AddCommand(policy.NewCmdPolicy(f))
	cmd.AddCommand(mutetiming.NewCmdMuteTiming(f))
	cmd.AddCommand(template.NewCmdTemplate(f))
	cmd.AddCommand(silence.NewCmdSilence(f))

	return cmd
}
