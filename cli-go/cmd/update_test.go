package cmd

import (
	"archive/tar"
	"compress/gzip"
	"os"
	"path/filepath"
	"testing"
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
