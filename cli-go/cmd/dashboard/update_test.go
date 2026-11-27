package dashboard

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"

	"github.com/spf13/cobra"

	"github.com/piyush-gambhir/grafana-cli/cli-go/internal/client"
	"github.com/piyush-gambhir/grafana-cli/cli-go/internal/cmdutil"
	"github.com/piyush-gambhir/grafana-cli/cli-go/internal/config"
)

func runDashboardUpdate(t *testing.T, args ...string) (client.DashboardCreateRequest, error) {
	t.Helper()
	return runDashboardWrite(t, newCmdDashboardUpdate, "abc123", http.StatusOK, "prod", args...)
}

func runDashboardWrite(t *testing.T, newCmd func(*cmdutil.Factory) *cobra.Command, uid string, lookupStatus int, existingFolder string, args ...string) (client.DashboardCreateRequest, error) {
	t.Helper()
	var posted client.DashboardCreateRequest
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch {
		case r.Method == http.MethodGet && r.URL.Path == "/api/dashboards/uid/"+uid:
			if lookupStatus == http.StatusNotFound {
				http.NotFound(w, r)
				return
			}
			json.NewEncoder(w).Encode(client.DashboardFullResponse{Meta: client.DashboardMeta{FolderUID: existingFolder}})
		case r.Method == http.MethodPost && r.URL.Path == "/api/dashboards/db":
			json.NewDecoder(r.Body).Decode(&posted)
			json.NewEncoder(w).Encode(client.DashboardCreateResponse{UID: uid})
		default:
			http.NotFound(w, r)
		}
	}))
	t.Cleanup(ts.Close)

	file := filepath.Join(t.TempDir(), "d.json")
	if err := os.WriteFile(file, []byte(`{"uid":"`+uid+`","title":"Prod"}`), 0600); err != nil {
		t.Fatal(err)
	}
	var out bytes.Buffer
	f := &cmdutil.Factory{
		IOStreams: cmdutil.IOStreams{Out: &out, ErrOut: &out},
		Resolved:  &config.ResolvedConfig{Output: "table"},
		Client: func() (*client.Client, error) {
			return &client.Client{BaseURL: ts.URL, HTTPClient: ts.Client(), UserAgent: "test"}, nil
		},
	}
	cmd := newCmd(f)
	cmd.SetArgs(append(args, "-f", file))
	cmd.SetOut(&out)
	cmd.SetErr(&out)
	err := cmd.Execute()
	return posted, err
}

func TestDashboardUpdate_KeepsExistingFolder(t *testing.T) {
	posted, err := runDashboardUpdate(t)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if posted.FolderUID != "prod" {
		t.Errorf("folderUid = %q, want existing folder %q", posted.FolderUID, "prod")
	}

	posted, err = runDashboardUpdate(t, "--folder", "other")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if posted.FolderUID != "other" {
		t.Errorf("folderUid = %q, want %q", posted.FolderUID, "other")
	}
}

func TestDashboardUpdate_RejectsPositionalUID(t *testing.T) {
	if _, err := runDashboardUpdate(t, "abc123"); err == nil {
		t.Fatal("expected an error for a positional argument")
	}
}

func TestDashboardUpdate_NewDashboardPostsWithoutFolder(t *testing.T) {
	posted, err := runDashboardWrite(t, newCmdDashboardUpdate, "abc123", http.StatusNotFound, "")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if posted.FolderUID != "" {
		t.Errorf("folderUid = %q, want empty for a new dashboard", posted.FolderUID)
	}
}

func TestDashboardUpdate_GeneralFolderStays(t *testing.T) {
	posted, err := runDashboardWrite(t, newCmdDashboardUpdate, "abc123", http.StatusOK, "")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if posted.FolderUID != "" {
		t.Errorf("folderUid = %q, want General (empty)", posted.FolderUID)
	}
}

func TestDashboardImport_OverwriteKeepsExistingFolder(t *testing.T) {
	posted, err := runDashboardWrite(t, newCmdDashboardImport, "abc123", http.StatusOK, "prod", "--overwrite")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if posted.FolderUID != "prod" {
		t.Errorf("folderUid = %q, want existing folder %q", posted.FolderUID, "prod")
	}

	posted, err = runDashboardWrite(t, newCmdDashboardImport, "abc123", http.StatusOK, "prod", "--overwrite", "--folder", "other")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if posted.FolderUID != "other" {
		t.Errorf("folderUid = %q, want %q", posted.FolderUID, "other")
	}
}

func TestDashboardImport_WithoutOverwriteSkipsLookup(t *testing.T) {
	posted, err := runDashboardWrite(t, newCmdDashboardImport, "abc123", http.StatusOK, "prod")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if posted.FolderUID != "" {
		t.Errorf("folderUid = %q, want empty without --overwrite", posted.FolderUID)
	}
}
