package cmd

import (
	"archive/tar"
	"bytes"
	"compress/gzip"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/piyush-gambhir/grafana-cli/cli-go/internal/build"
	"github.com/piyush-gambhir/grafana-cli/cli-go/internal/update"
)

// An entry named with a traversal path must still land inside destDir.
func TestExtractBinaryStaysInDestDir(t *testing.T) {
	dir := t.TempDir()
	archive := filepath.Join(dir, "release.tar.gz")
	f, err := os.Create(archive)
	if err != nil {
		t.Fatal(err)
	}
	gz := gzip.NewWriter(f)
	tw := tar.NewWriter(gz)
	payload := []byte("new binary")
	if err := tw.WriteHeader(&tar.Header{Name: "../../grafana-cli", Mode: 0o755, Size: int64(len(payload)), Typeflag: tar.TypeReg}); err != nil {
		t.Fatal(err)
	}
	if _, err := tw.Write(payload); err != nil {
		t.Fatal(err)
	}
	for _, c := range []interface{ Close() error }{tw, gz, f} {
		if err := c.Close(); err != nil {
			t.Fatal(err)
		}
	}

	destDir := filepath.Join(dir, "out")
	if err := os.Mkdir(destDir, 0o755); err != nil {
		t.Fatal(err)
	}
	got, err := extractBinary(archive, destDir)
	if err != nil {
		t.Fatal(err)
	}
	if want := filepath.Join(destDir, "grafana"); got != want {
		t.Fatalf("extracted to %s, want %s", got, want)
	}
	data, err := os.ReadFile(got)
	if err != nil || string(data) != string(payload) {
		t.Fatalf("extracted payload = %q, %v", data, err)
	}
}

// stubUpdateCheck pretends to run a release build on osName with info as the
// latest-release answer, without calling the GitHub API.
func stubUpdateCheck(t *testing.T, osName string, info *update.UpdateInfo) {
	t.Helper()
	oldGOOS, oldCheck, oldVersion := goos, checkForUpdateFresh, build.Version
	t.Cleanup(func() { goos, checkForUpdateFresh, build.Version = oldGOOS, oldCheck, oldVersion })
	goos = osName
	build.Version = info.CurrentVersion
	checkForUpdateFresh = func(string, string, string) (*update.UpdateInfo, error) { return info, nil }
}

func runUpdateCmd(args ...string) (string, error) {
	var out bytes.Buffer
	cmd := newUpdateCmd()
	cmd.SetOut(&out)
	cmd.SetErr(&out)
	cmd.SetArgs(args)
	err := cmd.Execute()
	return out.String(), err
}

var windowsUpdate = &update.UpdateInfo{
	Available:      true,
	CurrentVersion: "0.2.7",
	LatestVersion:  "0.2.8",
	ReleaseURL:     "https://github.com/piyush-gambhir/grafana-cli/releases/tag/v0.2.8",
}

// Windows releases ship a .zip with grafana.exe, so update must refuse to
// install (before prompting) and point at the release page.
func TestUpdateRefusesInstallOnWindows(t *testing.T) {
	stubUpdateCheck(t, "windows", windowsUpdate)
	out, err := runUpdateCmd()
	if err == nil {
		t.Fatal("expected update to refuse on Windows")
	}
	for _, want := range []string{"grafana.exe", windowsUpdate.ReleaseURL} {
		if !strings.Contains(err.Error(), want) {
			t.Errorf("error %q does not mention %q", err, want)
		}
	}
	if strings.Contains(out, "Do you want to update?") {
		t.Errorf("prompted before refusing: %q", out)
	}
}

func TestUpdateCheckWorksOnWindows(t *testing.T) {
	stubUpdateCheck(t, "windows", windowsUpdate)
	out, err := runUpdateCmd("--check")
	if err != nil {
		t.Fatalf("update --check on Windows: %v", err)
	}
	if !strings.Contains(out, windowsUpdate.ReleaseURL) {
		t.Errorf("output %q does not link the release", out)
	}
}
