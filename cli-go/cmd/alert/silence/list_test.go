package silence

import (
	"testing"

	"github.com/piyush-gambhir/grafana-cli/cli-go/internal/client"
)

func TestFormatMatcher(t *testing.T) {
	for _, tc := range []struct {
		m    client.Matcher
		want string
	}{
		{client.Matcher{Name: "a", Value: "b", IsEqual: true}, "a=b"},
		{client.Matcher{Name: "a", Value: "b"}, "a!=b"},
		{client.Matcher{Name: "a", Value: "b.*", IsRegex: true, IsEqual: true}, "a=~b.*"},
		{client.Matcher{Name: "a", Value: "b.*", IsRegex: true}, "a!~b.*"},
	} {
		if got := formatMatcher(tc.m); got != tc.want {
			t.Errorf("formatMatcher(%+v) = %q, want %q", tc.m, got, tc.want)
		}
	}
}
