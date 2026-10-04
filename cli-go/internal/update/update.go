// Package update checks GitHub Releases for a newer grafana CLI and prints the
// update notice. The background check (Start) runs at most once a day and only
// in an interactive terminal; see NotifierDisabled for the opt-outs.
//
// The latest release is read from the redirect of
// https://github.com/<repo>/releases/latest, not from api.github.com, whose
// 60 requests/hour unauthenticated limit is shared by everyone behind one IP.
package update

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"runtime"
	"strconv"
	"strings"
	"time"
)

const (
	// Repo is the GitHub repository that publishes releases.
	Repo = "piyush-gambhir/grafana-cli"
	// BinaryName is the executable name used in notices.
	BinaryName = "grafana"
	// EnvPrefix is the CLI's environment variable prefix.
	EnvPrefix = "GRAFANA"

	// MethodSelf means `grafana update` replaces the binary in place.
	MethodSelf = "self"
	// MethodGo means the binary lives in a Go bin directory, where `make
	// install` from a source checkout put it; it is never self-replaced.
	MethodGo = "go"

	// SelfUpdateCommand updates a release binary.
	SelfUpdateCommand = BinaryName + " update"
	// SourceUpdateCommand updates a binary built from source. The main package
	// is named cli-go, so `go install .../cli-go@latest` would install a binary
	// named cli-go; the repo documents `make install` instead.
	SourceUpdateCommand = "git pull && make install (in your grafana-cli/cli-go checkout)"

	cacheDuration     = 24 * time.Hour
	cacheFileName     = "update-check.json"
	backgroundTimeout = 3 * time.Second
	explicitTimeout   = 15 * time.Second
)

// Test seams.
var (
	releaseBaseURL = "https://github.com"
	now            = time.Now
)

// releaseTagRE matches a plain release tag (vMAJOR.MINOR.PATCH). The version
// ends up in URLs and terminal output, so nothing else is accepted.
var releaseTagRE = regexp.MustCompile(`^v\d+\.\d+\.\d+$`)

// semverRE matches MAJOR.MINOR.PATCH with an optional pre-release or build
// suffix (git describe output such as v0.2.9-3-gabc1234-dirty qualifies).
var semverRE = regexp.MustCompile(`^v?(\d+)\.(\d+)\.(\d+)([-+][0-9A-Za-z.+-]+)?$`)

// UpdateInfo describes the latest release relative to the running version.
// Versions are stored without the leading "v".
type UpdateInfo struct {
	Available      bool
	CurrentVersion string
	LatestVersion  string
	ReleaseURL     string
}

// cacheEntry is the on-disk update-check.json.
type cacheEntry struct {
	LastChecked     string `json:"last_checked,omitempty"`
	LatestVersion   string `json:"latest_version,omitempty"`
	NotifiedVersion string `json:"notified_version,omitempty"`
	NotifiedAt      string `json:"notified_at,omitempty"`
}

// IsReleaseVersion reports whether v is a semver release version. "dev",
// empty, and commit-hash versions are development builds.
func IsReleaseVersion(v string) bool {
	return semverRE.MatchString(v)
}

// ReleaseURL returns the release notes page for version.
func ReleaseURL(version string) string {
	return fmt.Sprintf("https://github.com/%s/releases/tag/v%s", Repo, strings.TrimPrefix(version, "v"))
}

// NotifierDisabled reports whether the background check must be skipped: a
// non-terminal stderr, CI, an opt-out variable, quiet mode, or a dev build.
func NotifierDisabled(getenv func(string) string, version string, quiet, stderrIsTerminal bool) bool {
	if !stderrIsTerminal || quiet || !IsReleaseVersion(version) {
		return true
	}
	for _, name := range []string{"CI", EnvPrefix + "_NO_UPDATE_NOTIFIER", "NO_UPDATE_NOTIFIER"} {
		if getenv(name) != "" {
			return true
		}
	}
	return false
}

// Start begins the background check and returns a channel that receives the
// answer. A cache entry less than 24 hours old answers synchronously (no
// network), so the notice never depends on goroutine timing. Otherwise GitHub
// is asked in a goroutine (3-second timeout) and the outcome is cached; a
// failed check is cached too, so it is not retried for 24 hours.
func Start(currentVersion, configDir string) <-chan *UpdateInfo {
	ch := make(chan *UpdateInfo, 1)
	path := filepath.Join(configDir, cacheFileName)
	entry := readEntry(path)
	if fresh(entry) {
		ch <- infoFrom(currentVersion, entry)
		return ch
	}
	go func() {
		info, _ := fetchAndCache(currentVersion, path, backgroundTimeout)
		ch <- info
	}()
	return ch
}

// CheckForUpdateFresh always asks GitHub, bypassing (and refreshing) the cache.
func CheckForUpdateFresh(currentVersion, configDir string) (*UpdateInfo, error) {
	return fetchAndCache(currentVersion, filepath.Join(configDir, cacheFileName), explicitTimeout)
}

// Cached returns the cached answer without touching the network, or nil when
// there is no successful check from the last 24 hours.
func Cached(currentVersion, configDir string) *UpdateInfo {
	entry := readEntry(filepath.Join(configDir, cacheFileName))
	if !fresh(entry) || entry.LatestVersion == "" {
		return nil
	}
	return infoFrom(currentVersion, entry)
}

// ClearCache removes the cached check, for example after a successful update.
func ClearCache(configDir string) error {
	err := os.Remove(filepath.Join(configDir, cacheFileName))
	if os.IsNotExist(err) {
		return nil
	}
	return err
}

// Notify prints the notice unless it was already shown for this latest
// version in the last 24 hours, then records that it was shown. It reports
// whether the notice was printed.
func Notify(w io.Writer, info *UpdateInfo, configDir, method string) bool {
	if info == nil || !info.Available {
		return false
	}
	path := filepath.Join(configDir, cacheFileName)
	entry := readEntry(path)
	if entry.NotifiedVersion == info.LatestVersion {
		if at, err := time.Parse(time.RFC3339, entry.NotifiedAt); err == nil {
			if age := now().Sub(at); age >= 0 && age < cacheDuration {
				return false
			}
		}
	}
	PrintNotice(w, info, method)
	entry.NotifiedVersion = info.LatestVersion
	entry.NotifiedAt = now().UTC().Format(time.RFC3339)
	writeEntry(path, entry)
	return true
}

// PrintNotice writes the update notice, preceded by a blank line.
func PrintNotice(w io.Writer, info *UpdateInfo, method string) {
	fmt.Fprintf(w, "\nA new version of %s is available: v%s -> v%s\n", BinaryName, info.CurrentVersion, info.LatestVersion)
	fmt.Fprintf(w, "Update with: %s\n", UpdateCommand(method))
	fmt.Fprintf(w, "Release notes: %s\n", info.ReleaseURL)
}

// UpdateCommand returns the command that updates a binary installed by method.
func UpdateCommand(method string) string {
	if method == MethodGo {
		return SourceUpdateCommand
	}
	return SelfUpdateCommand
}

// DetectInstallMethod reports MethodGo when exePath lives in a Go bin
// directory ($GOBIN, $GOPATH/bin, or ~/go/bin), else MethodSelf.
func DetectInstallMethod(exePath string) string {
	home, _ := os.UserHomeDir()
	return InstallMethodFor(exePath, os.Getenv, home)
}

// InstallMethodFor is DetectInstallMethod with the environment and home
// directory supplied by the caller.
func InstallMethodFor(exePath string, getenv func(string) string, home string) string {
	var dirs []string
	if gobin := getenv("GOBIN"); gobin != "" {
		dirs = append(dirs, gobin)
	}
	for _, gopath := range filepath.SplitList(getenv("GOPATH")) {
		if gopath != "" {
			dirs = append(dirs, filepath.Join(gopath, "bin"))
		}
	}
	if home != "" {
		dirs = append(dirs, filepath.Join(home, "go", "bin"))
	}
	exeDir := canonicalPath(filepath.Dir(exePath))
	for _, dir := range dirs {
		if samePath(exeDir, canonicalPath(dir)) {
			return MethodGo
		}
	}
	return MethodSelf
}

func canonicalPath(p string) string {
	if resolved, err := filepath.EvalSymlinks(p); err == nil {
		p = resolved
	}
	if abs, err := filepath.Abs(p); err == nil {
		p = abs
	}
	return filepath.Clean(p)
}

func samePath(a, b string) bool {
	if runtime.GOOS == "windows" {
		return strings.EqualFold(a, b)
	}
	return a == b
}

func fresh(entry cacheEntry) bool {
	checked, err := time.Parse(time.RFC3339, entry.LastChecked)
	if err != nil {
		return false
	}
	age := now().Sub(checked)
	return age >= 0 && age < cacheDuration
}

func infoFrom(currentVersion string, entry cacheEntry) *UpdateInfo {
	current := strings.TrimPrefix(currentVersion, "v")
	info := &UpdateInfo{CurrentVersion: current}
	if entry.LatestVersion == "" {
		return info
	}
	info.LatestVersion = entry.LatestVersion
	info.ReleaseURL = ReleaseURL(entry.LatestVersion)
	info.Available, _ = isNewer(entry.LatestVersion, current)
	return info
}

// fetchAndCache asks GitHub for the latest release and records the outcome.
// On failure it keeps the last known latest version but still stamps
// last_checked, so a broken network does not cause a request per command.
// The entry is re-read after the request so a notice another process recorded
// meanwhile is kept.
func fetchAndCache(currentVersion, path string, timeout time.Duration) (*UpdateInfo, error) {
	latest, err := fetchLatest(timeout)
	entry := readEntry(path)
	entry.LastChecked = now().UTC().Format(time.RFC3339)
	if err == nil {
		entry.LatestVersion = latest
	}
	writeEntry(path, entry)
	if err != nil {
		return nil, err
	}
	return infoFrom(currentVersion, entry), nil
}

// fetchLatest returns the latest release version (without the "v") from the
// redirect GitHub sends for /releases/latest. The redirect is not followed:
// its Location must be https://github.com/<repo>/releases/tag/vX.Y.Z on the
// same host, or the check fails.
func fetchLatest(timeout time.Duration) (string, error) {
	client := &http.Client{
		Timeout: timeout,
		CheckRedirect: func(*http.Request, []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}
	latestURL := fmt.Sprintf("%s/%s/releases/latest", releaseBaseURL, Repo)
	req, err := http.NewRequest(http.MethodGet, latestURL, nil)
	if err != nil {
		return "", fmt.Errorf("creating request: %w", err)
	}
	req.Header.Set("User-Agent", "grafana-cli")

	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("requesting %s: %w", latestURL, err)
	}
	_ = resp.Body.Close()
	if resp.StatusCode < 300 || resp.StatusCode > 399 {
		return "", fmt.Errorf("%s returned status %d, expected a redirect to the latest release", latestURL, resp.StatusCode)
	}
	loc, err := resp.Location()
	if err != nil {
		return "", fmt.Errorf("%s redirect has no usable Location header", latestURL)
	}
	tag, ok := strings.CutPrefix(loc.Path, "/"+Repo+"/releases/tag/")
	if loc.Scheme != req.URL.Scheme || loc.Host != req.URL.Host || !ok {
		return "", fmt.Errorf("%s redirected to %s, not a %s release tag", latestURL, loc.Redacted(), Repo)
	}
	if !releaseTagRE.MatchString(tag) {
		return "", fmt.Errorf("unexpected release tag %q", tag)
	}
	return strings.TrimPrefix(tag, "v"), nil
}

func readEntry(path string) cacheEntry {
	var entry cacheEntry
	data, err := os.ReadFile(path)
	if err != nil {
		return entry
	}
	if json.Unmarshal(data, &entry) != nil {
		return cacheEntry{}
	}
	if entry.LatestVersion != "" && (!IsReleaseVersion(entry.LatestVersion) || strings.ContainsAny(entry.LatestVersion, "-+")) {
		entry.LatestVersion = ""
	}
	return entry
}

// writeEntry saves the cache best-effort via a temp file and rename, so a
// concurrent reader never sees a partial file.
func writeEntry(path string, entry cacheEntry) {
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0o700); err != nil {
		return
	}
	data, err := json.MarshalIndent(entry, "", "  ")
	if err != nil {
		return
	}
	tmp, err := os.CreateTemp(dir, ".update-check-*")
	if err != nil {
		return
	}
	_, werr := tmp.Write(data)
	cerr := tmp.Close()
	if werr != nil || cerr != nil || os.Rename(tmp.Name(), path) != nil {
		_ = os.Remove(tmp.Name())
	}
}

// isNewer reports whether latest is a higher MAJOR.MINOR.PATCH than current.
func isNewer(latest, current string) (bool, error) {
	l, err := parseSemver(latest)
	if err != nil {
		return false, err
	}
	c, err := parseSemver(current)
	if err != nil {
		return false, err
	}
	for i := range l {
		if l[i] != c[i] {
			return l[i] > c[i], nil
		}
	}
	return false, nil
}

func parseSemver(v string) ([3]int, error) {
	var out [3]int
	m := semverRE.FindStringSubmatch(v)
	if m == nil {
		return out, fmt.Errorf("invalid semver: %s", v)
	}
	for i := range out {
		n, err := strconv.Atoi(m[i+1])
		if err != nil {
			return out, fmt.Errorf("invalid semver: %s", v)
		}
		out[i] = n
	}
	return out, nil
}
