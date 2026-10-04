package cmd

import (
	"archive/tar"
	"archive/zip"
	"bufio"
	"compress/gzip"
	"context"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"github.com/spf13/cobra"
	"golang.org/x/term"

	"github.com/piyush-gambhir/grafana-cli/cli-go/internal/build"
	"github.com/piyush-gambhir/grafana-cli/cli-go/internal/config"
	"github.com/piyush-gambhir/grafana-cli/cli-go/internal/output"
	"github.com/piyush-gambhir/grafana-cli/cli-go/internal/update"
)

const (
	installScriptURL = "https://raw.githubusercontent.com/" + update.Repo + "/main/install.sh"
	maxChecksumsSize = 1 << 20
)

// Test seams: tests point these at temp dirs, fake platforms, and httptest
// servers so nothing touches the real network or the test binary.
var (
	goos                = runtime.GOOS
	goarch              = runtime.GOARCH
	checkForUpdateFresh = update.CheckForUpdateFresh
	executablePath      = currentExecutable
	installMethodFor    = update.DetectInstallMethod
	stdinIsTerminal     = func() bool { return term.IsTerminal(int(os.Stdin.Fd())) }
	releaseDownloadBase = "https://github.com"
	downloadClient      = &http.Client{Timeout: 120 * time.Second}
)

// updateCheckResult is the `update --check` report.
type updateCheckResult struct {
	CurrentVersion  string `json:"current_version" yaml:"current_version"`
	LatestVersion   string `json:"latest_version" yaml:"latest_version"`
	UpdateAvailable bool   `json:"update_available" yaml:"update_available"`
	ReleaseURL      string `json:"release_url" yaml:"release_url"`
	InstallMethod   string `json:"install_method" yaml:"install_method"`
}

func newUpdateCmd() *cobra.Command {
	var checkOnly, yes bool

	cmd := &cobra.Command{
		Use:         "update",
		Annotations: map[string]string{"mutates": "true"},
		Short:       "Update grafana to the latest version",
		Long: `Check for and install the latest grafana release from GitHub Releases.

grafana update downloads the release archive for this OS and architecture,
verifies its SHA-256 against the release's checksums.txt, and replaces the
running executable. It works on macOS, Linux, and Windows; on Windows the old
binary is moved aside to grafana.exe.old and deleted on a later start. If the
executable's directory is not writable, re-run with sudo or reinstall with the
install script into a writable directory. A grafana built from source into a
Go bin directory ($GOBIN, $GOPATH/bin, ~/go/bin) is not replaced: update it
with "git pull && make install" in your grafana-cli/cli-go checkout.

It asks "Update now? [Y/n]" in a terminal. Pass --yes to skip the prompt;
with --no-input or without a terminal, --yes is required. --read-only blocks
installing, while --check is always allowed. --check -o json reports
current_version, latest_version, update_available, release_url, and
install_method (self or go).

Update notice: in an interactive terminal, grafana checks GitHub at most once
a day and, when a newer release exists, prints a three-line notice on stderr
after the command's output (once a day per version). The check is skipped
when stderr is not a terminal, CI is set, --quiet or GRAFANA_QUIET is on,
GRAFANA_NO_UPDATE_NOTIFIER or NO_UPDATE_NOTIFIER is set to any value, or the
build is a development build.

Examples:
  grafana update                  # prompt, then install the latest release
  grafana update --yes            # install without a prompt
  grafana update --check          # report current and latest versions
  grafana update --check -o json  # the same, as JSON`,
		Args: cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			out := cmd.OutOrStdout()
			OutputFormat = flagOutput
			current := strings.TrimPrefix(build.Version, "v")
			devBuild := !update.IsReleaseVersion(build.Version)

			if devBuild && !checkOnly {
				fmt.Fprintln(out, "Self-update is not available for development builds.")
				fmt.Fprintln(out, "Install a release (see https://github.com/"+update.Repo+"/releases) to enable updates.")
				return nil
			}

			info, err := checkForUpdateFresh(build.Version, config.ConfigDir())
			if err != nil {
				return fmt.Errorf("checking for updates: %w", err)
			}

			exe, exeErr := executablePath()
			method := update.MethodSelf
			if exeErr == nil {
				method = installMethodFor(exe)
			}

			if checkOnly {
				return printUpdateCheck(out, info, current, method, devBuild)
			}

			if !info.Available {
				fmt.Fprintf(out, "grafana v%s is already the latest version.\n", current)
				return nil
			}

			fmt.Fprintf(out, "Update available: v%s -> v%s\n", current, info.LatestVersion)
			fmt.Fprintf(out, "Release notes: %s\n", info.ReleaseURL)

			if method == update.MethodGo {
				fmt.Fprintf(out, "This grafana was built from source into a Go bin directory (%s), so it is not replaced.\n", filepath.Dir(exe))
				fmt.Fprintf(out, "Update with: %s\n", update.SourceUpdateCommand)
				return nil
			}
			if exeErr != nil {
				return exeErr
			}

			// update skips the root's config setup, so enforce --read-only here.
			resolved, _, err := loadAndResolveConfig(cmd)
			if err != nil {
				return err
			}
			if err := checkPermissions(cmd, resolved); err != nil {
				return err
			}

			if !yes {
				if flagNoInput {
					return errors.New("update needs confirmation: pass --yes to install without a prompt")
				}
				if !stdinIsTerminal() {
					return errors.New("update needs confirmation and stdin is not a terminal: pass --yes to install without a prompt")
				}
				if !confirmUpdate(cmd.InOrStdin(), out) {
					fmt.Fprintln(out, "Update cancelled.")
					return nil
				}
			}

			if err := installRelease(cmd.Context(), out, info.LatestVersion, exe); err != nil {
				return err
			}
			fmt.Fprintf(out, "Updated grafana v%s -> v%s\n", current, info.LatestVersion)
			fmt.Fprintf(out, "Release notes: %s\n", info.ReleaseURL)
			return nil
		},
	}

	cmd.Flags().BoolVar(&checkOnly, "check", false, "Only report whether an update is available (always queries GitHub)")
	cmd.Flags().BoolVarP(&yes, "yes", "y", false, "Install without asking for confirmation")

	return cmd
}

func printUpdateCheck(w io.Writer, info *update.UpdateInfo, current, method string, devBuild bool) error {
	result := updateCheckResult{
		CurrentVersion:  current,
		LatestVersion:   info.LatestVersion,
		UpdateAvailable: info.Available && !devBuild,
		ReleaseURL:      info.ReleaseURL,
		InstallMethod:   method,
	}
	if flagOutput != "" && flagOutput != string(output.FormatTable) {
		return output.Print(w, flagOutput, result, nil)
	}

	display := current
	if !devBuild {
		display = "v" + current
	}
	fmt.Fprintf(w, "Current version: %s\n", display)
	fmt.Fprintf(w, "Latest version:  v%s\n", result.LatestVersion)
	switch {
	case devBuild:
		fmt.Fprintln(w, "Update available: unknown (development build)")
	case result.UpdateAvailable:
		fmt.Fprintln(w, "Update available: yes")
	default:
		fmt.Fprintln(w, "Update available: no")
	}
	fmt.Fprintf(w, "Release notes: %s\n", result.ReleaseURL)
	if result.UpdateAvailable {
		fmt.Fprintf(w, "Update with: %s\n", update.UpdateCommand(method))
	}
	return nil
}

// confirmUpdate asks "Update now? [Y/n]"; an empty answer means yes.
func confirmUpdate(in io.Reader, w io.Writer) bool {
	fmt.Fprint(w, "Update now? [Y/n] ")
	line, err := bufio.NewReader(in).ReadString('\n')
	if err != nil && line == "" {
		return false
	}
	switch strings.ToLower(strings.TrimSpace(line)) {
	case "", "y", "yes":
		return true
	default:
		return false
	}
}

func currentExecutable() (string, error) {
	exe, err := os.Executable()
	if err != nil {
		return "", fmt.Errorf("finding the current executable: %w", err)
	}
	if resolved, err := filepath.EvalSymlinks(exe); err == nil {
		exe = resolved
	}
	return exe, nil
}

// cleanupOldExecutable deletes the grafana.exe.old a Windows update leaves
// behind (it cannot be deleted while it is still running). Best-effort.
func cleanupOldExecutable(osName, exe string) {
	if osName == "windows" && exe != "" {
		_ = os.Remove(exe + ".old")
	}
}

// archiveName is the GoReleaser archive for a platform (see .goreleaser.yaml:
// "{{ .ProjectName }}_{{ .Os }}_{{ .Arch }}", zip on Windows).
func archiveName(osName, arch string) string {
	ext := ".tar.gz"
	if osName == "windows" {
		ext = ".zip"
	}
	return fmt.Sprintf("grafana-cli_%s_%s%s", osName, arch, ext)
}

// binaryName is the executable inside the release archive.
func binaryName(osName string) string {
	if osName == "windows" {
		return "grafana.exe"
	}
	return "grafana"
}

// installRelease downloads, verifies, and installs version over exe. Any
// failure leaves the current executable in place.
func installRelease(ctx context.Context, w io.Writer, version, exe string) error {
	if err := checkWritableDir(filepath.Dir(exe)); err != nil {
		return err
	}

	tmpDir, err := os.MkdirTemp("", "grafana-cli-update-*")
	if err != nil {
		return fmt.Errorf("creating temp directory: %w", err)
	}
	defer os.RemoveAll(tmpDir)

	archive := archiveName(goos, goarch)
	baseURL := fmt.Sprintf("%s/%s/releases/download/v%s/", releaseDownloadBase, update.Repo, version)
	archivePath := filepath.Join(tmpDir, "release-archive")

	fmt.Fprintf(w, "Downloading %s...\n", archive)
	if err := downloadFile(ctx, archivePath, baseURL+archive); err != nil {
		return fmt.Errorf("downloading %s: %w", archive, err)
	}

	fmt.Fprintln(w, "Verifying checksum...")
	if err := verifyChecksum(ctx, archivePath, baseURL+"checksums.txt", archive); err != nil {
		return fmt.Errorf("checksum verification failed: %w", err)
	}

	binaryPath, err := extractBinary(archivePath, tmpDir, goos)
	if err != nil {
		return fmt.Errorf("extracting update: %w", err)
	}

	fmt.Fprintf(w, "Replacing %s...\n", exe)
	return replaceExecutable(goos, binaryPath, exe)
}

// checkWritableDir fails early, before any download, when the new binary
// could not be written next to the current one.
func checkWritableDir(dir string) error {
	f, err := os.CreateTemp(dir, ".grafana-update-*")
	if err != nil {
		return notWritableError(dir, err)
	}
	name := f.Name()
	f.Close()
	return os.Remove(name)
}

func notWritableError(dir string, err error) error {
	return fmt.Errorf("cannot write to %s (%v): the current grafana was left unchanged; re-run with sudo (or as Administrator on Windows), or reinstall into a writable directory with the install script: curl -sSfL %s | INSTALL_DIR=~/.local/bin sh", dir, err, installScriptURL)
}

func newRequest(ctx context.Context, url string) (*http.Request, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", "grafana-cli")
	return req, nil
}

func downloadFile(ctx context.Context, dst, url string) error {
	req, err := newRequest(ctx, url)
	if err != nil {
		return err
	}
	resp, err := downloadClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("download failed with status %d", resp.StatusCode)
	}

	out, err := os.OpenFile(dst, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0o600)
	if err != nil {
		return err
	}
	if err := copyUpdatePayload(out, resp.Body); err != nil {
		out.Close()
		return err
	}
	return out.Close()
}

// verifyChecksum checks filePath against the archive's entry in the release's
// checksums.txt and refuses a mismatch or a missing entry.
func verifyChecksum(ctx context.Context, filePath, checksumURL, archive string) error {
	req, err := newRequest(ctx, checksumURL)
	if err != nil {
		return err
	}
	resp, err := downloadClient.Do(req)
	if err != nil {
		return fmt.Errorf("downloading checksums: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("checksums download failed with status %d", resp.StatusCode)
	}
	body, err := io.ReadAll(io.LimitReader(resp.Body, maxChecksumsSize+1))
	if err != nil {
		return fmt.Errorf("reading checksums: %w", err)
	}
	if len(body) > maxChecksumsSize {
		return errors.New("checksums.txt is too large")
	}

	// Format: "<sha256>  <filename>" per line.
	var expected string
	for _, line := range strings.Split(string(body), "\n") {
		fields := strings.Fields(line)
		if len(fields) == 2 && strings.TrimPrefix(fields[1], "*") == archive {
			expected = strings.ToLower(fields[0])
			break
		}
	}
	if expected == "" {
		return fmt.Errorf("no checksum for %s in checksums.txt", archive)
	}
	if decoded, err := hex.DecodeString(expected); err != nil || len(decoded) != sha256.Size {
		return fmt.Errorf("malformed checksum for %s in checksums.txt", archive)
	}

	f, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer f.Close()
	h := sha256.New()
	if err := copyUpdatePayload(h, f); err != nil {
		return fmt.Errorf("computing checksum: %w", err)
	}
	if actual := hex.EncodeToString(h.Sum(nil)); actual != expected {
		return fmt.Errorf("SHA-256 mismatch for %s: expected %s, got %s", archive, expected, actual)
	}
	return nil
}

// extractBinary writes the release binary from the archive to destDir under
// the fixed name "grafana", so no part of an archive entry name ever reaches
// the file system (CodeQL go/zipslip). It refuses entries with absolute or
// parent-relative paths, a binary entry that is not a regular file, and a
// binary larger than the release size limit.
func extractBinary(archivePath, destDir, osName string) (string, error) {
	outPath := filepath.Join(destDir, "grafana")
	want := binaryName(osName)
	var found bool
	var err error
	if osName == "windows" {
		found, err = extractFromZip(archivePath, outPath, want)
	} else {
		found, err = extractFromTarGz(archivePath, outPath, want)
	}
	if err != nil {
		return "", err
	}
	if !found {
		return "", fmt.Errorf("%s not found in the release archive", want)
	}
	return outPath, nil
}

func extractFromTarGz(archivePath, outPath, want string) (bool, error) {
	f, err := os.Open(archivePath)
	if err != nil {
		return false, err
	}
	defer f.Close()
	gz, err := gzip.NewReader(f)
	if err != nil {
		return false, fmt.Errorf("opening gzip: %w", err)
	}
	defer gz.Close()

	tr := tar.NewReader(gz)
	found := false
	for {
		hdr, err := tr.Next()
		if err == io.EOF {
			return found, nil
		}
		if err != nil {
			return false, fmt.Errorf("reading tar: %w", err)
		}
		name, err := archiveEntryName(hdr.Name)
		if err != nil {
			return false, err
		}
		if name != want {
			continue
		}
		if hdr.Typeflag != tar.TypeReg {
			return false, fmt.Errorf("archive entry %q is not a regular file", hdr.Name)
		}
		if found {
			return false, fmt.Errorf("archive has more than one %s", want)
		}
		if hdr.Size > maxReleaseArtifactBytes {
			return false, fmt.Errorf("%s exceeds the %d MiB limit", want, maxReleaseArtifactBytes>>20)
		}
		if err := writeExtracted(outPath, tr); err != nil {
			return false, err
		}
		found = true
	}
}

func extractFromZip(archivePath, outPath, want string) (bool, error) {
	r, err := zip.OpenReader(archivePath)
	if err != nil {
		return false, fmt.Errorf("opening zip: %w", err)
	}
	defer r.Close()

	found := false
	for _, entry := range r.File {
		name, err := archiveEntryName(entry.Name)
		if err != nil {
			return false, err
		}
		if name != want {
			continue
		}
		if !entry.Mode().IsRegular() {
			return false, fmt.Errorf("archive entry %q is not a regular file", entry.Name)
		}
		if found {
			return false, fmt.Errorf("archive has more than one %s", want)
		}
		if entry.UncompressedSize64 > uint64(maxReleaseArtifactBytes) {
			return false, fmt.Errorf("%s exceeds the %d MiB limit", want, maxReleaseArtifactBytes>>20)
		}
		rc, err := entry.Open()
		if err != nil {
			return false, err
		}
		err = writeExtracted(outPath, rc)
		rc.Close()
		if err != nil {
			return false, err
		}
		found = true
	}
	return found, nil
}

// archiveEntryName returns name without a leading "./" or trailing "/",
// refusing absolute paths, drive letters, backslashes, and ".." components.
func archiveEntryName(name string) (string, error) {
	if path.IsAbs(name) || strings.ContainsAny(name, `\:`) {
		return "", fmt.Errorf("unsafe path %q in the release archive", name)
	}
	for _, part := range strings.Split(name, "/") {
		if part == ".." {
			return "", fmt.Errorf("unsafe path %q in the release archive", name)
		}
	}
	return strings.TrimSuffix(strings.TrimPrefix(name, "./"), "/"), nil
}

func writeExtracted(outPath string, r io.Reader) error {
	out, err := os.OpenFile(outPath, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0o755)
	if err != nil {
		return err
	}
	if err := copyUpdatePayload(out, r); err != nil {
		out.Close()
		return err
	}
	return out.Close()
}

// replaceExecutable installs the binary at src over exe. The new file is
// written next to exe first, so a failure never leaves exe half-written.
// macOS/Linux rename it over exe. Windows cannot overwrite a running .exe but
// can rename it, so exe moves aside to exe.old (deleted on a later start by
// cleanupOldExecutable) and is restored if the final rename fails.
func replaceExecutable(osName, src, exe string) error {
	dir := filepath.Dir(exe)
	tmp, err := os.CreateTemp(dir, ".grafana-update-*")
	if err != nil {
		return notWritableError(dir, err)
	}
	tmpPath := tmp.Name()
	defer func() {
		if tmpPath != "" {
			_ = os.Remove(tmpPath)
		}
	}()

	in, err := os.Open(src)
	if err != nil {
		tmp.Close()
		return err
	}
	err = copyUpdatePayload(tmp, in)
	in.Close()
	if err == nil {
		err = tmp.Chmod(0o755)
	}
	if cerr := tmp.Close(); err == nil {
		err = cerr
	}
	if err != nil {
		return fmt.Errorf("writing the new executable: %w", err)
	}

	if osName != "windows" {
		if err := os.Rename(tmpPath, exe); err != nil {
			return fmt.Errorf("replacing %s: %w", exe, err)
		}
		tmpPath = ""
		return nil
	}

	old := exe + ".old"
	_ = os.Remove(old) // left over from an earlier update
	if err := os.Rename(exe, old); err != nil {
		return fmt.Errorf("moving %s aside: %w", exe, err)
	}
	if err := os.Rename(tmpPath, exe); err != nil {
		if rerr := os.Rename(old, exe); rerr != nil {
			return fmt.Errorf("installing the new executable: %w (restoring the old one also failed: %v; it is at %s)", err, rerr, old)
		}
		return fmt.Errorf("installing the new executable: %w", err)
	}
	tmpPath = ""
	return nil
}
