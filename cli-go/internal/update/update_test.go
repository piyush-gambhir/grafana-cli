package update

import (
	"bytes"
	"strings"
	"testing"
)

// `grafana update` cannot install on Windows, so the notice must not suggest it there.
func TestPrintUpdateNoticeByOS(t *testing.T) {
	info := &UpdateInfo{
		Available:      true,
		CurrentVersion: "0.2.7",
		LatestVersion:  "0.2.8",
		ReleaseURL:     "https://github.com/piyush-gambhir/grafana-cli/releases/tag/v0.2.8",
	}
	for _, tc := range []struct {
		goos        string
		wantCommand bool
	}{{"linux", true}, {"darwin", true}, {"windows", false}} {
		t.Run(tc.goos, func(t *testing.T) {
			old := goos
			t.Cleanup(func() { goos = old })
			goos = tc.goos

			var out bytes.Buffer
			PrintUpdateNotice(&out, info)
			if !strings.Contains(out.String(), info.ReleaseURL) {
				t.Errorf("notice %q does not link the release", out.String())
			}
			if got := strings.Contains(out.String(), "grafana update"); got != tc.wantCommand {
				t.Errorf("notice mentions `grafana update` = %v, want %v: %q", got, tc.wantCommand, out.String())
			}
		})
	}
}
