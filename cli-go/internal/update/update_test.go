package update

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"sync/atomic"
	"testing"
	"time"
)

// fakeGitHub serves /<repo>/releases/latest like github.com: a 302 to the tag
// page when status is 302 and tag is set (no Location when tag is empty), or
// a bare status otherwise. It counts requests and fails the test if the
// client follows the redirect.
func fakeGitHub(t *testing.T, status int, tag string) *atomic.Int32 {
	t.Helper()
	var hits atomic.Int32
	var srv *httptest.Server
	srv = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/"+Repo+"/releases/latest" {
			t.Errorf("unexpected request %s (the redirect must not be followed)", r.URL.Path)
			http.NotFound(w, r)
			return
		}
		hits.Add(1)
		if status == http.StatusFound && tag != "" {
			w.Header().Set("Location", tagURL(srv.URL, tag))
		}
		w.WriteHeader(status)
	}))
	t.Cleanup(srv.Close)
	setReleaseBase(t, srv.URL)
	return &hits
}

// tagURL is the release page URL GitHub redirects to, escaped like a real
// Location header.
func tagURL(base, tag string) string {
	u, err := url.Parse(base)
	if err != nil {
		panic(err)
	}
	u.Path = "/" + Repo + "/releases/tag/" + tag
	return u.String()
}

func setReleaseBase(t *testing.T, base string) {
	t.Helper()
	old := releaseBaseURL
	releaseBaseURL = base
	t.Cleanup(func() { releaseBaseURL = old })
}

func setNow(t *testing.T, at time.Time) {
	t.Helper()
	old := now
	now = func() time.Time { return at }
	t.Cleanup(func() { now = old })
}

func receive(t *testing.T, ch <-chan *UpdateInfo) *UpdateInfo {
	t.Helper()
	select {
	case info := <-ch:
		return info
	case <-time.After(5 * time.Second):
		t.Fatal("background check did not answer")
		return nil
	}
}

func TestNotifierDisabled(t *testing.T) {
	env := func(vars map[string]string) func(string) string {
		return func(k string) string { return vars[k] }
	}
	for _, tc := range []struct {
		name    string
		vars    map[string]string
		version string
		quiet   bool
		tty     bool
		want    bool
	}{
		{"interactive release build", nil, "0.2.9", false, true, false},
		{"stderr not a terminal", nil, "0.2.9", false, false, true},
		{"CI set", map[string]string{"CI": "true"}, "0.2.9", false, true, true},
		{"CLI opt-out", map[string]string{"GRAFANA_NO_UPDATE_NOTIFIER": "1"}, "0.2.9", false, true, true},
		{"generic opt-out", map[string]string{"NO_UPDATE_NOTIFIER": "yes"}, "0.2.9", false, true, true},
		{"quiet", nil, "0.2.9", true, true, true},
		{"dev build", nil, "dev", false, true, true},
		{"empty version", nil, "", false, true, true},
		{"commit hash version", nil, "e120dad", false, true, true},
		{"git describe build", nil, "v0.2.9-3-ge120dad-dirty", false, true, false},
	} {
		t.Run(tc.name, func(t *testing.T) {
			if got := NotifierDisabled(env(tc.vars), tc.version, tc.quiet, tc.tty); got != tc.want {
				t.Errorf("NotifierDisabled = %v, want %v", got, tc.want)
			}
		})
	}
}

func TestStartAnswersFromFreshCacheWithoutNetwork(t *testing.T) {
	hits := fakeGitHub(t, http.StatusFound, "v0.2.10")
	dir := t.TempDir()
	setNow(t, time.Date(2026, 10, 4, 12, 0, 0, 0, time.UTC))
	writeEntry(filepath.Join(dir, cacheFileName), cacheEntry{
		LastChecked:   "2026-10-04T01:00:00Z",
		LatestVersion: "0.2.10",
	})

	// A fresh cache answers synchronously: the value is already there.
	select {
	case info := <-Start("0.2.9", dir):
		if !info.Available || info.LatestVersion != "0.2.10" {
			t.Fatalf("info = %+v", info)
		}
		if want := "https://github.com/" + Repo + "/releases/tag/v0.2.10"; info.ReleaseURL != want {
			t.Errorf("ReleaseURL = %q, want %q", info.ReleaseURL, want)
		}
	default:
		t.Fatal("fresh cache did not answer synchronously")
	}
	if hits.Load() != 0 {
		t.Errorf("GitHub called %d times with a fresh cache", hits.Load())
	}
}

func TestStartFetchesWhenStaleAndCaches(t *testing.T) {
	hits := fakeGitHub(t, http.StatusFound, "v0.2.10")
	dir := t.TempDir()
	setNow(t, time.Date(2026, 10, 4, 12, 0, 0, 0, time.UTC))
	writeEntry(filepath.Join(dir, cacheFileName), cacheEntry{
		LastChecked:     "2026-10-03T11:00:00Z", // 25h ago
		LatestVersion:   "0.2.9",
		NotifiedVersion: "0.2.9",
		NotifiedAt:      "2026-10-03T11:00:00Z",
	})

	info := receive(t, Start("0.2.9", dir))
	if info == nil || !info.Available || info.LatestVersion != "0.2.10" {
		t.Fatalf("info = %+v", info)
	}
	entry := readEntry(filepath.Join(dir, cacheFileName))
	if entry.LatestVersion != "0.2.10" || entry.LastChecked != "2026-10-04T12:00:00Z" {
		t.Errorf("cache not refreshed: %+v", entry)
	}
	if entry.NotifiedVersion != "0.2.9" {
		t.Errorf("refresh dropped notified_version: %+v", entry)
	}

	receive(t, Start("0.2.9", dir))
	if hits.Load() != 1 {
		t.Errorf("GitHub called %d times, want 1 (second check should use the cache)", hits.Load())
	}
}

// A notice recorded by another process while the GitHub request is in flight
// must survive the cache refresh, or the next command repeats the notice.
func TestRefreshKeepsNoticeRecordedDuringFetch(t *testing.T) {
	release := make(chan struct{})
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		<-release
		http.Redirect(w, r, "/"+Repo+"/releases/tag/v0.2.10", http.StatusFound)
	}))
	t.Cleanup(srv.Close)
	setReleaseBase(t, srv.URL)
	dir := t.TempDir()
	setNow(t, time.Date(2026, 10, 4, 12, 0, 0, 0, time.UTC))

	ch := Start("0.2.9", dir) // empty cache: fetch in the background
	info := &UpdateInfo{Available: true, CurrentVersion: "0.2.9", LatestVersion: "0.2.10", ReleaseURL: ReleaseURL("0.2.10")}
	if !Notify(&bytes.Buffer{}, info, dir, MethodSelf) {
		t.Fatal("notice not shown")
	}
	close(release)
	receive(t, ch)

	if entry := readEntry(filepath.Join(dir, cacheFileName)); entry.NotifiedVersion != "0.2.10" || entry.LatestVersion != "0.2.10" {
		t.Errorf("refresh lost the recorded notice: %+v", entry)
	}
}

func TestFailedCheckIsCached(t *testing.T) {
	for _, tc := range []struct {
		name   string
		status int
		tag    string
	}{
		{"server error", http.StatusInternalServerError, ""},
		{"rate limited", http.StatusForbidden, ""},
		{"redirect without Location", http.StatusFound, ""},
		{"non-semver tag", http.StatusFound, "v0.2.10-rc.1"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			hits := fakeGitHub(t, tc.status, tc.tag)
			dir := t.TempDir()
			setNow(t, time.Date(2026, 10, 4, 12, 0, 0, 0, time.UTC))

			if info := receive(t, Start("0.2.9", dir)); info != nil {
				t.Errorf("failed check returned %+v", info)
			}
			if info := receive(t, Start("0.2.9", dir)); info == nil || info.Available {
				t.Errorf("cached failure answer = %+v, want no update", info)
			}
			if hits.Load() != 1 {
				t.Errorf("GitHub called %d times, want 1 (a failed check is cached)", hits.Load())
			}
		})
	}
}

// The latest version comes from the Location of the /releases/latest
// redirect; anything but a same-host vX.Y.Z tag page is a failed check.
func TestFetchLatestReadsRedirect(t *testing.T) {
	var location string
	status := http.StatusFound
	var followed atomic.Int32
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/"+Repo+"/releases/latest" {
			followed.Add(1)
			_, _ = w.Write([]byte("followed"))
			return
		}
		if location != "" {
			w.Header().Set("Location", location)
		}
		w.WriteHeader(status)
	}))
	t.Cleanup(srv.Close)
	setReleaseBase(t, srv.URL)
	foreign, _ := url.Parse(srv.URL)
	foreign.Host = "github.com.evil.example"

	for _, tc := range []struct {
		name     string
		status   int
		location string
		want     string // empty means the check must fail
	}{
		{"good tag", http.StatusFound, tagURL(srv.URL, "v0.2.10"), "0.2.10"},
		{"relative Location", http.StatusFound, "/" + Repo + "/releases/tag/v1.2.3", "1.2.3"},
		{"missing Location", http.StatusFound, "", ""},
		{"foreign host", http.StatusFound, tagURL(foreign.String(), "v0.2.10"), ""},
		{"other repo", http.StatusFound, srv.URL + "/someone/else/releases/tag/v0.2.10", ""},
		{"no releases yet", http.StatusFound, srv.URL + "/" + Repo + "/releases", ""},
		{"pre-release tag", http.StatusFound, tagURL(srv.URL, "v0.2.10-rc.1"), ""},
		{"tag without v", http.StatusFound, tagURL(srv.URL, "0.2.10"), ""},
		{"terminal escape in tag", http.StatusFound, tagURL(srv.URL, "v1.0.0\x1b[31m"), ""},
		{"C1 control in foreign query", http.StatusFound, "https://evil.example/x?\x9b2J", ""},
		{"not a redirect", http.StatusOK, "", ""},
		{"rate limited", http.StatusForbidden, "", ""},
	} {
		t.Run(tc.name, func(t *testing.T) {
			status, location = tc.status, tc.location
			got, err := fetchLatest(explicitTimeout)
			if tc.want == "" {
				if err == nil {
					t.Fatalf("fetchLatest = %q, want an error", got)
				}
				// The error reaches the terminal via `grafana update`.
				for _, b := range []byte(err.Error()) {
					if b < 0x20 || b >= 0x7f {
						t.Errorf("error contains raw byte %#x: %q", b, err)
						break
					}
				}
				return
			}
			if err != nil || got != tc.want {
				t.Fatalf("fetchLatest = %q, %v, want %q", got, err, tc.want)
			}
		})
	}
	if n := followed.Load(); n != 0 {
		t.Errorf("client followed the redirect %d times", n)
	}
}

func TestCheckRejectsUnexpectedTag(t *testing.T) {
	fakeGitHub(t, http.StatusFound, "v1.0.0\x1b[31m")
	_, err := CheckForUpdateFresh("0.2.9", t.TempDir())
	if err == nil || !strings.Contains(err.Error(), "unexpected release tag") {
		t.Fatalf("err = %v, want an unexpected release tag error", err)
	}
}

func TestCheckForUpdateFreshBypassesCache(t *testing.T) {
	hits := fakeGitHub(t, http.StatusFound, "v0.2.10")
	dir := t.TempDir()
	setNow(t, time.Date(2026, 10, 4, 12, 0, 0, 0, time.UTC))
	writeEntry(filepath.Join(dir, cacheFileName), cacheEntry{LastChecked: "2026-10-04T11:00:00Z", LatestVersion: "0.2.9"})

	info, err := CheckForUpdateFresh("0.2.9", dir)
	if err != nil {
		t.Fatal(err)
	}
	if !info.Available || info.LatestVersion != "0.2.10" || hits.Load() != 1 {
		t.Errorf("info = %+v, hits = %d", info, hits.Load())
	}
	// `update --check` stores its answer, so the notifier agrees with it.
	if cached := receive(t, Start("0.2.9", dir)); cached == nil || cached.LatestVersion != "0.2.10" || hits.Load() != 1 {
		t.Errorf("notifier after a fresh check = %+v, hits = %d", cached, hits.Load())
	}
}

// The attempt is recorded before the request, so a process that exits before
// the answer arrives does not make every later command ask GitHub again.
func TestStartRecordsAttemptBeforeRequest(t *testing.T) {
	release := make(chan struct{})
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		<-release
		w.WriteHeader(http.StatusBadGateway)
	}))
	t.Cleanup(srv.Close)
	setReleaseBase(t, srv.URL)
	dir := t.TempDir()
	setNow(t, time.Date(2026, 10, 4, 12, 0, 0, 0, time.UTC))

	// Nobody reads this answer during the test, like a command that exited.
	pending := Start("0.2.9", dir)
	t.Cleanup(func() { // runs before the seams are restored
		close(release)
		<-pending
	})
	if entry := readEntry(filepath.Join(dir, cacheFileName)); entry.LastChecked != "2026-10-04T12:00:00Z" {
		t.Fatalf("attempt not recorded before the request: %+v", entry)
	}
	select {
	case <-Start("0.2.9", dir):
	default:
		t.Fatal("second check did not answer from the cache")
	}
}

func TestCachedNeverFetches(t *testing.T) {
	hits := fakeGitHub(t, http.StatusFound, "v0.2.10")
	dir := t.TempDir()
	setNow(t, time.Date(2026, 10, 4, 12, 0, 0, 0, time.UTC))
	if info := Cached("0.2.9", dir); info != nil {
		t.Errorf("Cached with no cache = %+v", info)
	}
	writeEntry(filepath.Join(dir, cacheFileName), cacheEntry{LastChecked: "2026-10-04T11:00:00Z", LatestVersion: "0.2.10"})
	if info := Cached("0.2.9", dir); info == nil || !info.Available {
		t.Errorf("Cached = %+v", info)
	}
	if hits.Load() != 0 {
		t.Errorf("Cached called GitHub %d times", hits.Load())
	}
}

func TestNoticeFormat(t *testing.T) {
	info := &UpdateInfo{Available: true, CurrentVersion: "0.2.9", LatestVersion: "0.2.10", ReleaseURL: ReleaseURL("0.2.10")}
	var self, gobin bytes.Buffer
	PrintNotice(&self, info, MethodSelf)
	PrintNotice(&gobin, info, MethodGo)

	want := "\nA new version of grafana is available: v0.2.9 -> v0.2.10\n" +
		"Update with: grafana update\n" +
		"Release notes: https://github.com/piyush-gambhir/grafana-cli/releases/tag/v0.2.10\n"
	if self.String() != want {
		t.Errorf("notice =\n%q\nwant\n%q", self.String(), want)
	}
	if !strings.Contains(gobin.String(), "Update with: git pull && make install") || strings.Contains(gobin.String(), "grafana update") {
		t.Errorf("Go bin notice = %q", gobin.String())
	}
}

func TestNotifyOncePerVersionPerDay(t *testing.T) {
	dir := t.TempDir()
	start := time.Date(2026, 10, 4, 12, 0, 0, 0, time.UTC)
	info := &UpdateInfo{Available: true, CurrentVersion: "0.2.9", LatestVersion: "0.2.10", ReleaseURL: ReleaseURL("0.2.10")}
	notify := func(at time.Time, info *UpdateInfo) bool {
		setNow(t, at)
		var out bytes.Buffer
		shown := Notify(&out, info, dir, MethodSelf)
		if shown != (out.Len() > 0) {
			t.Fatalf("Notify returned %v but wrote %q", shown, out.String())
		}
		return shown
	}

	if !notify(start, info) {
		t.Fatal("first notice not shown")
	}
	if notify(start.Add(time.Minute), info) {
		t.Error("notice shown twice for the same version within 24h")
	}
	if notify(start.Add(23*time.Hour), info) {
		t.Error("notice shown again within 24h")
	}
	newer := *info
	newer.LatestVersion, newer.ReleaseURL = "0.2.11", ReleaseURL("0.2.11")
	if !notify(start.Add(23*time.Hour+time.Minute), &newer) {
		t.Error("notice for a newer version suppressed")
	}
	if !notify(start.Add(48*time.Hour), &newer) {
		t.Error("notice not shown again after 24h")
	}
}

func TestInstallMethodFor(t *testing.T) {
	home := t.TempDir()
	gobin := t.TempDir()
	gopath := t.TempDir()
	other := t.TempDir()
	for _, d := range []string{filepath.Join(home, "go", "bin"), filepath.Join(gopath, "bin")} {
		if err := os.MkdirAll(d, 0o755); err != nil {
			t.Fatal(err)
		}
	}
	env := func(k string) string {
		return map[string]string{"GOBIN": gobin, "GOPATH": gopath}[k]
	}
	for _, tc := range []struct {
		exe  string
		want string
	}{
		{filepath.Join(gobin, "grafana"), MethodGo},
		{filepath.Join(gopath, "bin", "grafana"), MethodGo},
		{filepath.Join(home, "go", "bin", "grafana"), MethodGo},
		{filepath.Join(other, "grafana"), MethodSelf},
	} {
		if got := InstallMethodFor(tc.exe, env, home); got != tc.want {
			t.Errorf("InstallMethodFor(%s) = %q, want %q", tc.exe, got, tc.want)
		}
	}
}

func TestIsNewer(t *testing.T) {
	for _, tc := range []struct {
		latest, current string
		want            bool
	}{
		{"0.2.10", "0.2.9", true},
		{"0.2.9", "0.2.9", false},
		{"0.2.9", "0.2.10", false},
		{"1.0.0", "0.9.9", true},
		{"0.2.10", "v0.2.9-3-gabc", true},
	} {
		if got, _ := isNewer(tc.latest, tc.current); got != tc.want {
			t.Errorf("isNewer(%s, %s) = %v, want %v", tc.latest, tc.current, got, tc.want)
		}
	}
}
