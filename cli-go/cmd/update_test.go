package cmd

import (
	"archive/tar"
	"archive/zip"
	"bytes"
	"compress/gzip"
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
	"time"

	"github.com/spf13/cobra"

	"github.com/piyush-gambhir/grafana-cli/cli-go/internal/build"
	"github.com/piyush-gambhir/grafana-cli/cli-go/internal/update"
)

const testReleaseURL = "https://github.com/piyush-gambhir/grafana-cli/releases/tag/v0.2.10"

var availableUpdate = update.UpdateInfo{
	Available:      true,
	CurrentVersion: "0.2.9",
	LatestVersion:  "0.2.10",
	ReleaseURL:     testReleaseURL,
}

var noUpdate = update.UpdateInfo{
	CurrentVersion: "0.2.9",
	LatestVersion:  "0.2.9",
	ReleaseURL:     "https://github.com/piyush-gambhir/grafana-cli/releases/tag/v0.2.9",
}

// testEnv isolates a test from the developer's environment and the package
// seams: release build 0.2.9, an empty config dir, no opt-out variables, an
// interactive stderr, and a GitHub check that fails the test if called.
func testEnv(t *testing.T) {
	t.Helper()
	t.Setenv("XDG_CONFIG_HOME", t.TempDir())
	for _, name := range []string{"CI", "GRAFANA_NO_UPDATE_NOTIFIER", "NO_UPDATE_NOTIFIER", "GRAFANA_QUIET", "GRAFANA_NO_INPUT", "GRAFANA_READ_ONLY", "GRAFANA_URL", "GRAFANA_TOKEN"} {
		t.Setenv(name, "")
	}
	oldVersion := build.Version
	oldSeams := []any{goos, goarch, checkForUpdateFresh, executablePath, installMethodFor, stdinIsTerminal, releaseDownloadBase, downloadClient, stderrIsTerminal, startUpdateCheck, detectInstallMethod}
	t.Cleanup(func() {
		build.Version = oldVersion
		goos = oldSeams[0].(string)
		goarch = oldSeams[1].(string)
		checkForUpdateFresh = oldSeams[2].(func(string, string) (*update.UpdateInfo, error))
		executablePath = oldSeams[3].(func() (string, error))
		installMethodFor = oldSeams[4].(func(string) string)
		stdinIsTerminal = oldSeams[5].(func() bool)
		releaseDownloadBase = oldSeams[6].(string)
		downloadClient = oldSeams[7].(*http.Client)
		stderrIsTerminal = oldSeams[8].(func() bool)
		startUpdateCheck = oldSeams[9].(func(string, string) <-chan *update.UpdateInfo)
		detectInstallMethod = oldSeams[10].(func() string)
	})
	build.Version = "0.2.9"
	stderrIsTerminal = func() bool { return true }
	stdinIsTerminal = func() bool { return false }
	installMethodFor = func(string) string { return update.MethodSelf }
	detectInstallMethod = func() string { return update.MethodSelf }
	executablePath = func() (string, error) { return "", fmt.Errorf("no executable in tests") }
	checkForUpdateFresh = func(string, string) (*update.UpdateInfo, error) {
		t.Error("unexpected GitHub release check")
		return nil, fmt.Errorf("no network in tests")
	}
	startUpdateCheck = func(string, string) <-chan *update.UpdateInfo {
		t.Error("unexpected background update check")
		return nil
	}
}

// stubFreshCheck answers `grafana update` release lookups with info.
func stubFreshCheck(info update.UpdateInfo) {
	checkForUpdateFresh = func(string, string) (*update.UpdateInfo, error) { return &info, nil }
}

// runRoot runs the root command (plus a no-op "noop" command) with args.
func runRoot(t *testing.T, stdin string, args ...string) (stdout, stderr string, err error) {
	t.Helper()
	root := newRootCmd()
	root.AddCommand(&cobra.Command{Use: "noop", RunE: func(*cobra.Command, []string) error { return nil }})
	var out, errOut bytes.Buffer
	root.SetOut(&out)
	root.SetErr(&errOut)
	root.SetIn(strings.NewReader(stdin))
	root.SetArgs(args)
	err = root.Execute()
	return out.String(), errOut.String(), err
}

func TestNotifierSuppressed(t *testing.T) {
	for _, tc := range []struct {
		name  string
		setup func(t *testing.T)
		args  []string
	}{
		{"stderr not a terminal", func(t *testing.T) { stderrIsTerminal = func() bool { return false } }, []string{"noop"}},
		{"CI set", func(t *testing.T) { t.Setenv("CI", "true") }, []string{"noop"}},
		{"GRAFANA_NO_UPDATE_NOTIFIER", func(t *testing.T) { t.Setenv("GRAFANA_NO_UPDATE_NOTIFIER", "1") }, []string{"noop"}},
		{"NO_UPDATE_NOTIFIER", func(t *testing.T) { t.Setenv("NO_UPDATE_NOTIFIER", "1") }, []string{"noop"}},
		{"--quiet", func(t *testing.T) {}, []string{"noop", "--quiet"}},
		{"-q", func(t *testing.T) {}, []string{"noop", "-q"}},
		{"GRAFANA_QUIET", func(t *testing.T) { t.Setenv("GRAFANA_QUIET", "1") }, []string{"noop"}},
		{"dev build", func(t *testing.T) { build.Version = "dev" }, []string{"noop"}},
		{"empty version", func(t *testing.T) { build.Version = "" }, []string{"noop"}},
		{"non-semver version", func(t *testing.T) { build.Version = "e120dad" }, []string{"noop"}},
		{"version command", func(t *testing.T) {}, []string{"version"}},
		{"completion command", discardStdout, []string{"completion", "bash"}},
		{"help command", func(t *testing.T) {}, []string{"help"}},
		{"__complete", func(t *testing.T) {}, []string{"__complete", "dash"}},
		{"update command", func(t *testing.T) { stubFreshCheck(noUpdate) }, []string{"update", "--check"}},
	} {
		t.Run(tc.name, func(t *testing.T) {
			testEnv(t) // startUpdateCheck fails the test if it is called
			tc.setup(t)
			_, stderr, err := runRoot(t, "", tc.args...)
			if err != nil {
				t.Fatalf("command failed: %v", err)
			}
			if strings.Contains(stderr, "new version") {
				t.Errorf("notice printed: %q", stderr)
			}
		})
	}
}

func TestNotifierShowsNoticeOncePerVersion(t *testing.T) {
	testEnv(t)
	calls := 0
	startUpdateCheck = func(version, _ string) <-chan *update.UpdateInfo {
		calls++
		if version != "0.2.9" {
			t.Errorf("checked version %q", version)
		}
		ch := make(chan *update.UpdateInfo, 1)
		info := availableUpdate
		ch <- &info
		return ch
	}

	_, stderr, err := runRoot(t, "", "noop")
	if err != nil {
		t.Fatal(err)
	}
	want := "\nA new version of grafana is available: v0.2.9 -> v0.2.10\nUpdate with: grafana update\nRelease notes: " + testReleaseURL + "\n"
	if stderr != want {
		t.Errorf("stderr = %q, want %q", stderr, want)
	}

	_, stderr, err = runRoot(t, "", "noop")
	if err != nil {
		t.Fatal(err)
	}
	if stderr != "" {
		t.Errorf("notice repeated within 24h: %q", stderr)
	}
	if calls != 2 {
		t.Errorf("background check started %d times, want 2", calls)
	}
}

func TestNotifierGoInstallUpdateLine(t *testing.T) {
	testEnv(t)
	detectInstallMethod = func() string { return update.MethodGo }
	startUpdateCheck = func(string, string) <-chan *update.UpdateInfo {
		ch := make(chan *update.UpdateInfo, 1)
		info := availableUpdate
		ch <- &info
		return ch
	}
	_, stderr, err := runRoot(t, "", "noop")
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(stderr, "Update with: "+update.SourceUpdateCommand+"\n") || strings.Contains(stderr, "grafana update") {
		t.Errorf("stderr = %q", stderr)
	}
}

// The notice must never wait for a slow GitHub answer.
func TestNotifierDoesNotWaitForPendingCheck(t *testing.T) {
	testEnv(t)
	startUpdateCheck = func(string, string) <-chan *update.UpdateInfo { return make(chan *update.UpdateInfo) }
	start := time.Now()
	_, stderr, err := runRoot(t, "", "noop")
	if err != nil {
		t.Fatal(err)
	}
	if stderr != "" || time.Since(start) > time.Second {
		t.Errorf("stderr = %q after %s", stderr, time.Since(start))
	}
}

// Only the top-level commands skip the check; `dashboard update` is an
// ordinary command.
func TestSkipsUpdateCheckOnlyForTopLevelCommands(t *testing.T) {
	root := newRootCmd()
	for _, tc := range []struct {
		path []string
		want bool
	}{
		{[]string{"update"}, true},
		{[]string{"version"}, true},
		{[]string{"dashboard", "update"}, false},
		{[]string{"alert", "rule", "update"}, false},
		{[]string{"dashboard", "list"}, false},
	} {
		cmd, _, err := root.Find(tc.path)
		if err != nil {
			t.Fatal(err)
		}
		if got := skipsUpdateCheck(cmd); got != tc.want {
			t.Errorf("skipsUpdateCheck(%v) = %v, want %v", tc.path, got, tc.want)
		}
	}
}

// `dashboard update` shares its name with `grafana update`; it must still get
// the client setup and the read-only check (it used to skip both and panic).
func TestNestedUpdateCommandIsReadOnlyChecked(t *testing.T) {
	testEnv(t)
	stderrIsTerminal = func() bool { return false }
	file := filepath.Join(t.TempDir(), "dashboard.json")
	if err := os.WriteFile(file, []byte(`{"dashboard":{"title":"x"}}`), 0o600); err != nil {
		t.Fatal(err)
	}
	_, _, err := runRoot(t, "", "dashboard", "update", "-f", file, "--read-only", "--url", "http://127.0.0.1:1")
	if err == nil || !strings.Contains(err.Error(), "read-only") {
		t.Fatalf("err = %v, want the read-only block", err)
	}
}

func TestUpdateCheckJSON(t *testing.T) {
	for _, method := range []string{update.MethodSelf, update.MethodGo} {
		t.Run(method, func(t *testing.T) {
			testEnv(t)
			stubFreshCheck(availableUpdate)
			executablePath = func() (string, error) { return "/somewhere/grafana", nil }
			installMethodFor = func(string) string { return method }

			stdout, _, err := runRoot(t, "", "update", "--check", "-o", "json")
			if err != nil {
				t.Fatal(err)
			}
			var got map[string]any
			if err := json.Unmarshal([]byte(stdout), &got); err != nil {
				t.Fatalf("stdout is not JSON: %v\n%s", err, stdout)
			}
			want := map[string]any{
				"current_version":  "0.2.9",
				"latest_version":   "0.2.10",
				"update_available": true,
				"release_url":      testReleaseURL,
				"install_method":   method,
			}
			if fmt.Sprint(got) != fmt.Sprint(want) {
				t.Errorf("JSON = %v, want %v", got, want)
			}
		})
	}
}

func TestUpdateCheckText(t *testing.T) {
	testEnv(t)
	stubFreshCheck(availableUpdate)
	stdout, _, err := runRoot(t, "", "update", "--check")
	if err != nil {
		t.Fatal(err)
	}
	for _, want := range []string{"Current version: v0.2.9", "Latest version:  v0.2.10", "Update available: yes", "Release notes: " + testReleaseURL, "Update with: grafana update"} {
		if !strings.Contains(stdout, want) {
			t.Errorf("output missing %q:\n%s", want, stdout)
		}
	}
}

func TestUpdateAlreadyLatest(t *testing.T) {
	testEnv(t)
	stubFreshCheck(noUpdate)
	executablePath = func() (string, error) { return "/somewhere/grafana", nil }
	stdout, _, err := runRoot(t, "", "update", "--no-input")
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(stdout, "grafana v0.2.9 is already the latest version.") {
		t.Errorf("output = %q", stdout)
	}
}

func TestUpdateNeedsYesWithoutPrompt(t *testing.T) {
	for _, tc := range []struct {
		name string
		args []string
		env  string
		tty  bool
	}{
		{"--no-input", []string{"update", "--no-input"}, "", true},
		{"GRAFANA_NO_INPUT", []string{"update"}, "1", true},
		{"stdin not a terminal", []string{"update"}, "", false},
	} {
		t.Run(tc.name, func(t *testing.T) {
			testEnv(t)
			t.Setenv("GRAFANA_NO_INPUT", tc.env)
			stdinIsTerminal = func() bool { return tc.tty }
			stubFreshCheck(availableUpdate)
			exe := writeFakeExe(t, t.TempDir(), "grafana")
			executablePath = func() (string, error) { return exe, nil }

			_, _, err := runRoot(t, "", tc.args...)
			if err == nil || !strings.Contains(err.Error(), "--yes") {
				t.Fatalf("err = %v, want a message saying to pass --yes", err)
			}
			assertFile(t, exe, "old binary")
		})
	}
}

func TestUpdateReadOnlyBlocksInstallButNotCheck(t *testing.T) {
	testEnv(t)
	stubFreshCheck(availableUpdate)
	exe := writeFakeExe(t, t.TempDir(), "grafana")
	executablePath = func() (string, error) { return exe, nil }

	if _, _, err := runRoot(t, "", "update", "--check", "--read-only"); err != nil {
		t.Fatalf("update --check --read-only: %v", err)
	}
	_, _, err := runRoot(t, "", "update", "--yes", "--read-only")
	if err == nil || !strings.Contains(err.Error(), "read-only") {
		t.Fatalf("err = %v, want the read-only block", err)
	}
	t.Setenv("GRAFANA_READ_ONLY", "1")
	if _, _, err := runRoot(t, "", "update", "--yes"); err == nil || !strings.Contains(err.Error(), "read-only") {
		t.Fatalf("GRAFANA_READ_ONLY: err = %v, want the read-only block", err)
	}
	assertFile(t, exe, "old binary")
}

func TestUpdateGoInstallPrintsCommandInsteadOfReplacing(t *testing.T) {
	testEnv(t)
	stubFreshCheck(availableUpdate)
	exe := writeFakeExe(t, t.TempDir(), "grafana")
	executablePath = func() (string, error) { return exe, nil }
	installMethodFor = func(string) string { return update.MethodGo }
	releaseDownloadBase = "http://127.0.0.1:1" // any download attempt fails

	stdout, _, err := runRoot(t, "", "update", "--yes")
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(stdout, "Update with: "+update.SourceUpdateCommand) {
		t.Errorf("output = %q", stdout)
	}
	assertFile(t, exe, "old binary")
}

func TestUpdateInstallsRelease(t *testing.T) {
	for _, tc := range []struct {
		goos, goarch, exeName string
	}{
		{"linux", "amd64", "grafana"},
		{"darwin", "arm64", "grafana"},
		{"windows", "amd64", "grafana.exe"},
	} {
		t.Run(tc.goos, func(t *testing.T) {
			testEnv(t)
			stubFreshCheck(availableUpdate)
			goos, goarch = tc.goos, tc.goarch
			dir := t.TempDir()
			exe := writeFakeExe(t, dir, tc.exeName)
			executablePath = func() (string, error) { return exe, nil }
			stdinIsTerminal = func() bool { return true }

			archive := releaseArchive(t, tc.goos, []archiveEntry{
				{name: "README.md", body: "readme"},
				{name: binaryName(tc.goos), body: "new binary"},
			})
			name := archiveName(tc.goos, tc.goarch)
			serveRelease(t, name, archive, checksumLine(archive, name))
			cacheFile := filepath.Join(os.Getenv("XDG_CONFIG_HOME"), "grafana-cli", "update-check.json")
			if err := os.MkdirAll(filepath.Dir(cacheFile), 0o700); err != nil {
				t.Fatal(err)
			}
			if err := os.WriteFile(cacheFile, []byte(`{"latest_version":"0.2.10"}`), 0o600); err != nil {
				t.Fatal(err)
			}

			// Pressing Enter at the [Y/n] prompt means yes.
			stdout, _, err := runRoot(t, "\n", "update")
			if err != nil {
				t.Fatalf("update: %v\n%s", err, stdout)
			}
			for _, want := range []string{"Update now? [Y/n]", "Updated grafana v0.2.9 -> v0.2.10\nRelease notes: " + testReleaseURL} {
				if !strings.Contains(stdout, want) {
					t.Errorf("output missing %q:\n%s", want, stdout)
				}
			}
			assertFile(t, exe, "new binary")
			if tc.goos == "windows" {
				assertFile(t, exe+".old", "old binary")
			}
			if _, err := os.Stat(cacheFile); !os.IsNotExist(err) {
				t.Errorf("update cache not cleared: %v", err)
			}
			assertOnlyFiles(t, dir, tc.goos == "windows", tc.exeName)
		})
	}
}

func TestUpdatePromptCanBeDeclined(t *testing.T) {
	testEnv(t)
	stubFreshCheck(availableUpdate)
	exe := writeFakeExe(t, t.TempDir(), "grafana")
	executablePath = func() (string, error) { return exe, nil }
	stdinIsTerminal = func() bool { return true }
	releaseDownloadBase = "http://127.0.0.1:1"

	stdout, _, err := runRoot(t, "n\n", "update")
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(stdout, "Update cancelled.") {
		t.Errorf("output = %q", stdout)
	}
	assertFile(t, exe, "old binary")
}

func TestInstallRefusesBadChecksum(t *testing.T) {
	for _, tc := range []struct {
		name      string
		checksums func(archive []byte, name string) string
		wantErr   string
	}{
		{"mismatch", func(_ []byte, name string) string { return checksumLine([]byte("other bytes"), name) }, "mismatch"},
		{"missing entry", func(archive []byte, _ string) string { return checksumLine(archive, "grafana-cli_plan9_amd64.tar.gz") }, "no checksum"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			testEnv(t)
			goos, goarch = "linux", "amd64"
			exe := writeFakeExe(t, t.TempDir(), "grafana")
			archive := releaseArchive(t, "linux", []archiveEntry{{name: "grafana", body: "new binary"}})
			name := archiveName("linux", "amd64")
			serveRelease(t, name, archive, tc.checksums(archive, name))

			err := installRelease(context.Background(), io.Discard, "0.2.10", exe)
			if err == nil || !strings.Contains(err.Error(), tc.wantErr) {
				t.Fatalf("err = %v, want %q", err, tc.wantErr)
			}
			assertFile(t, exe, "old binary")
		})
	}
}

func TestInstallUnwritableDirKeepsOldBinary(t *testing.T) {
	if runtime.GOOS == "windows" || os.Geteuid() == 0 {
		t.Skip("needs POSIX permissions and a non-root user")
	}
	testEnv(t)
	goos, goarch = "linux", "amd64"
	dir := t.TempDir()
	exe := writeFakeExe(t, dir, "grafana")
	archive := releaseArchive(t, "linux", []archiveEntry{{name: "grafana", body: "new binary"}})
	name := archiveName("linux", "amd64")
	serveRelease(t, name, archive, checksumLine(archive, name))
	if err := os.Chmod(dir, 0o555); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = os.Chmod(dir, 0o755) })

	err := installRelease(context.Background(), io.Discard, "0.2.10", exe)
	if err == nil || !strings.Contains(err.Error(), "sudo") || !strings.Contains(err.Error(), "install script") {
		t.Fatalf("err = %v, want advice to use sudo or the install script", err)
	}
	assertFile(t, exe, "old binary")

	// The replace step itself refuses the same way.
	src := filepath.Join(t.TempDir(), "grafana")
	if err := os.WriteFile(src, []byte("new binary"), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := replaceExecutable("linux", src, exe); err == nil || !strings.Contains(err.Error(), "sudo") {
		t.Fatalf("replaceExecutable err = %v", err)
	}
	assertFile(t, exe, "old binary")
}

func TestReplaceExecutableUnix(t *testing.T) {
	dir := t.TempDir()
	exe := writeFakeExe(t, dir, "grafana")
	if err := os.Chmod(exe, 0o700); err != nil {
		t.Fatal(err)
	}
	src := filepath.Join(t.TempDir(), "grafana")
	if err := os.WriteFile(src, []byte("new binary"), 0o600); err != nil {
		t.Fatal(err)
	}
	if err := replaceExecutable("linux", src, exe); err != nil {
		t.Fatal(err)
	}
	assertFile(t, exe, "new binary")
	if runtime.GOOS != "windows" {
		if info, err := os.Stat(exe); err != nil || info.Mode().Perm() != 0o755 {
			t.Errorf("mode = %v, %v; want 0755", info.Mode().Perm(), err)
		}
	}
	assertOnlyFiles(t, dir, false, "grafana")
}

func TestReplaceExecutableWindowsMovesOldAside(t *testing.T) {
	dir := t.TempDir()
	exe := writeFakeExe(t, dir, "grafana.exe")
	if err := os.WriteFile(exe+".old", []byte("stale"), 0o755); err != nil {
		t.Fatal(err)
	}
	src := filepath.Join(t.TempDir(), "grafana")
	if err := os.WriteFile(src, []byte("new binary"), 0o755); err != nil {
		t.Fatal(err)
	}

	if err := replaceExecutable("windows", src, exe); err != nil {
		t.Fatal(err)
	}
	assertFile(t, exe, "new binary")
	assertFile(t, exe+".old", "old binary")
	assertOnlyFiles(t, dir, true, "grafana.exe")

	// A later start deletes the leftover, only on Windows.
	cleanupOldExecutable("linux", exe)
	assertFile(t, exe+".old", "old binary")
	cleanupOldExecutable("windows", exe)
	if _, err := os.Stat(exe + ".old"); !os.IsNotExist(err) {
		t.Errorf("grafana.exe.old not removed: %v", err)
	}
}

func TestExtractBinary(t *testing.T) {
	for _, osName := range []string{"linux", "windows"} {
		bin := binaryName(osName)
		for _, tc := range []struct {
			name    string
			entries []archiveEntry
			wantErr string
		}{
			{"binary at root", []archiveEntry{{name: "LICENSE", body: "x"}, {name: bin, body: "new binary"}}, ""},
			{"dot-slash prefix", []archiveEntry{{name: "./" + bin, body: "new binary"}}, ""},
			{"parent traversal", []archiveEntry{{name: "../../" + bin, body: "evil"}}, "unsafe path"},
			{"traversal in another entry", []archiveEntry{{name: bin, body: "new binary"}, {name: "docs/../../evil", body: "x"}}, "unsafe path"},
			{"absolute path", []archiveEntry{{name: "/tmp/" + bin, body: "evil"}}, "unsafe path"},
			{"backslash path", []archiveEntry{{name: `..\` + bin, body: "evil"}}, "unsafe path"},
			{"symlink binary", []archiveEntry{{name: bin, link: "/bin/sh"}}, "not a regular file"},
			{"directory binary", []archiveEntry{{name: bin, dir: true}}, "not a regular file"},
			{"duplicate binary", []archiveEntry{{name: bin, body: "a"}, {name: bin, body: "b"}}, "more than one"},
			{"binary missing", []archiveEntry{{name: "grafana-cli", body: "x"}}, "not found"},
		} {
			t.Run(osName+"/"+tc.name, func(t *testing.T) {
				dir := t.TempDir()
				archivePath := filepath.Join(dir, "release-archive")
				if err := os.WriteFile(archivePath, releaseArchive(t, osName, tc.entries), 0o600); err != nil {
					t.Fatal(err)
				}
				destDir := filepath.Join(dir, "out")
				if err := os.Mkdir(destDir, 0o755); err != nil {
					t.Fatal(err)
				}
				got, err := extractBinary(archivePath, destDir, osName)
				if tc.wantErr != "" {
					if err == nil || !strings.Contains(err.Error(), tc.wantErr) {
						t.Fatalf("err = %v, want %q", err, tc.wantErr)
					}
					if _, err := os.Stat(filepath.Join(dir, "evil")); !os.IsNotExist(err) {
						t.Errorf("traversal wrote outside destDir")
					}
					return
				}
				if err != nil {
					t.Fatal(err)
				}
				if want := filepath.Join(destDir, "grafana"); got != want {
					t.Fatalf("extracted to %s, want the fixed name %s", got, want)
				}
				assertFile(t, got, "new binary")
			})
		}
	}
}

func TestVersionShowsCachedLatestWithoutNetwork(t *testing.T) {
	testEnv(t)
	stdout, _, err := runRoot(t, "", "version")
	if err != nil {
		t.Fatal(err)
	}
	if strings.Contains(stdout, "latest") {
		t.Errorf("unknown latest printed: %q", stdout)
	}

	cacheFile := filepath.Join(os.Getenv("XDG_CONFIG_HOME"), "grafana-cli", "update-check.json")
	if err := os.MkdirAll(filepath.Dir(cacheFile), 0o700); err != nil {
		t.Fatal(err)
	}
	entry := fmt.Sprintf(`{"last_checked":%q,"latest_version":"0.2.10"}`, time.Now().UTC().Format(time.RFC3339))
	if err := os.WriteFile(cacheFile, []byte(entry), 0o600); err != nil {
		t.Fatal(err)
	}
	stdout, _, err = runRoot(t, "", "version")
	if err != nil {
		t.Fatal(err)
	}
	if !strings.HasPrefix(stdout, "grafana-cli version 0.2.9\n") || !strings.Contains(stdout, "  latest: 0.2.10\n  update_available: true\n") {
		t.Errorf("version output = %q", stdout)
	}
}

// --- helpers ---

// discardStdout silences commands that write to os.Stdout directly.
func discardStdout(t *testing.T) {
	devNull, err := os.OpenFile(os.DevNull, os.O_WRONLY, 0)
	if err != nil {
		t.Fatal(err)
	}
	old := os.Stdout
	os.Stdout = devNull
	t.Cleanup(func() { os.Stdout = old; devNull.Close() })
}

type archiveEntry struct {
	name, body, link string
	dir              bool
}

// releaseArchive builds a zip (Windows) or tar.gz archive in memory.
func releaseArchive(t *testing.T, osName string, entries []archiveEntry) []byte {
	t.Helper()
	var buf bytes.Buffer
	if osName == "windows" {
		zw := zip.NewWriter(&buf)
		for _, e := range entries {
			hdr := &zip.FileHeader{Name: e.name, Method: zip.Deflate}
			body := e.body
			switch {
			case e.link != "":
				hdr.SetMode(os.ModeSymlink | 0o777)
				body = e.link
			case e.dir:
				hdr.Name += "/"
				hdr.SetMode(os.ModeDir | 0o755)
			default:
				hdr.SetMode(0o755)
			}
			w, err := zw.CreateHeader(hdr)
			if err != nil {
				t.Fatal(err)
			}
			if _, err := w.Write([]byte(body)); err != nil {
				t.Fatal(err)
			}
		}
		if err := zw.Close(); err != nil {
			t.Fatal(err)
		}
		return buf.Bytes()
	}

	gz := gzip.NewWriter(&buf)
	tw := tar.NewWriter(gz)
	for _, e := range entries {
		hdr := &tar.Header{Name: e.name, Mode: 0o755, Size: int64(len(e.body)), Typeflag: tar.TypeReg}
		switch {
		case e.link != "":
			hdr.Typeflag, hdr.Linkname, hdr.Size = tar.TypeSymlink, e.link, 0
		case e.dir:
			hdr.Typeflag, hdr.Size = tar.TypeDir, 0
		}
		if err := tw.WriteHeader(hdr); err != nil {
			t.Fatal(err)
		}
		if hdr.Size > 0 {
			if _, err := tw.Write([]byte(e.body)); err != nil {
				t.Fatal(err)
			}
		}
	}
	if err := tw.Close(); err != nil {
		t.Fatal(err)
	}
	if err := gz.Close(); err != nil {
		t.Fatal(err)
	}
	return buf.Bytes()
}

func checksumLine(data []byte, name string) string {
	sum := sha256.Sum256(data)
	return hex.EncodeToString(sum[:]) + "  " + name + "\n"
}

// serveRelease serves the v0.2.10 release assets and points the installer at it.
func serveRelease(t *testing.T, name string, archive []byte, checksums string) {
	t.Helper()
	prefix := "/" + update.Repo + "/releases/download/v0.2.10/"
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case prefix + name:
			_, _ = w.Write(archive)
		case prefix + "checksums.txt":
			_, _ = io.WriteString(w, checksums)
		default:
			http.NotFound(w, r)
		}
	}))
	t.Cleanup(srv.Close)
	releaseDownloadBase = srv.URL
	downloadClient = srv.Client()
}

func writeFakeExe(t *testing.T, dir, name string) string {
	t.Helper()
	exe := filepath.Join(dir, name)
	if err := os.WriteFile(exe, []byte("old binary"), 0o755); err != nil {
		t.Fatal(err)
	}
	return exe
}

func assertFile(t *testing.T, path, want string) {
	t.Helper()
	data, err := os.ReadFile(path)
	if err != nil || string(data) != want {
		t.Fatalf("%s = %q, %v; want %q", filepath.Base(path), data, err, want)
	}
}

// assertOnlyFiles checks that no temp files were left next to the executable.
func assertOnlyFiles(t *testing.T, dir string, withOld bool, exeName string) {
	t.Helper()
	entries, err := os.ReadDir(dir)
	if err != nil {
		t.Fatal(err)
	}
	want := map[string]bool{exeName: true}
	if withOld {
		want[exeName+".old"] = true
	}
	for _, e := range entries {
		if !want[e.Name()] {
			t.Errorf("unexpected file %s left in the executable's directory", e.Name())
		}
	}
}
