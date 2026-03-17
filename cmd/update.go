package cmd

import (
	"archive/tar"
	"compress/gzip"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"github.com/spf13/cobra"

	"github.com/piyush-gambhir/grafana-cli/internal/build"
	"github.com/piyush-gambhir/grafana-cli/internal/config"
	"github.com/piyush-gambhir/grafana-cli/internal/update"
)

const updateRepo = "piyush-gambhir/grafana-cli"

func newUpdateCmd() *cobra.Command {
	var checkOnly bool

	cmd := &cobra.Command{
		Use:   "update",
		Short: "Update grafana to the latest version",
		Long:  "Check for and install the latest version of the grafana CLI from GitHub Releases.",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			configDir := config.ConfigDir()
			currentVersion := build.Version

			if currentVersion == "dev" {
				fmt.Fprintln(cmd.OutOrStdout(), "Update checking is not available for development builds.")
				fmt.Fprintln(cmd.OutOrStdout(), "Build from source or install a release to enable updates.")
				return nil
			}

			fmt.Fprintln(cmd.OutOrStdout(), "Checking for updates...")
			info, err := update.CheckForUpdateFresh(currentVersion, updateRepo, configDir)
			if err != nil {
				return fmt.Errorf("checking for updates: %w", err)
			}

			if checkOnly {
				if info.Available {
					update.PrintUpdateNotice(cmd.OutOrStdout(), info)
				} else {
					fmt.Fprintf(cmd.OutOrStdout(), "Already up to date (%s)\n", formatVer(currentVersion))
				}
				return nil
			}

			if !info.Available {
				fmt.Fprintf(cmd.OutOrStdout(), "Already up to date (%s)\n", formatVer(currentVersion))
				return nil
			}

			fmt.Fprintf(cmd.OutOrStdout(), "\nUpdate available: %s → %s\n",
				formatVer(info.CurrentVersion), formatVer(info.LatestVersion))
			if info.PublishedAt != "" {
				if t, err := time.Parse(time.RFC3339, info.PublishedAt); err == nil {
					fmt.Fprintf(cmd.OutOrStdout(), "Published: %s\n", t.Format("January 2, 2006"))
				}
			}
			fmt.Fprintf(cmd.OutOrStdout(), "Release:   %s\n\n", info.ReleaseURL)

			fmt.Fprint(cmd.OutOrStdout(), "Do you want to update? [y/N] ")
			var answer string
			fmt.Fscanln(os.Stdin, &answer)
			answer = strings.TrimSpace(strings.ToLower(answer))
			if answer != "y" && answer != "yes" {
				fmt.Fprintln(cmd.OutOrStdout(), "Update cancelled.")
				return nil
			}

			return performUpdate(cmd.OutOrStdout(), info.LatestVersion)
		},
	}

	cmd.Flags().BoolVar(&checkOnly, "check", false, "Only check if an update is available, don't install")

	return cmd
}

func performUpdate(w io.Writer, version string) error {
	osName := runtime.GOOS
	archName := runtime.GOARCH

	// Build the download URL matching the release asset naming convention.
	downloadURL := fmt.Sprintf(
		"https://github.com/%s/releases/download/v%s/grafana-cli_%s_%s.tar.gz",
		updateRepo, version, osName, archName,
	)

	fmt.Fprintf(w, "Downloading %s...\n", downloadURL)

	// Download to a temp directory.
	tmpDir, err := os.MkdirTemp("", "grafana-cli-update-*")
	if err != nil {
		return fmt.Errorf("creating temp directory: %w", err)
	}
	defer os.RemoveAll(tmpDir)

	archivePath := filepath.Join(tmpDir, "grafana-cli.tar.gz")
	if err := downloadFile(archivePath, downloadURL); err != nil {
		return fmt.Errorf("downloading update: %w", err)
	}

	fmt.Fprintf(w, "Extracting...\n")

	// Extract the binary from the tarball.
	binaryPath, err := extractBinary(archivePath, tmpDir)
	if err != nil {
		return fmt.Errorf("extracting update: %w", err)
	}

	// Get path to the currently running executable.
	execPath, err := os.Executable()
	if err != nil {
		return fmt.Errorf("finding current executable: %w", err)
	}
	execPath, err = filepath.EvalSymlinks(execPath)
	if err != nil {
		return fmt.Errorf("resolving executable path: %w", err)
	}

	fmt.Fprintf(w, "Replacing %s...\n", execPath)

	// Atomically replace: copy to a temp file next to the target, then rename.
	if err := atomicReplace(binaryPath, execPath); err != nil {
		return fmt.Errorf("replacing binary: %w", err)
	}

	fmt.Fprintf(w, "Successfully updated to %s!\n", formatVer(version))
	return nil
}

func downloadFile(dst, url string) error {
	client := &http.Client{Timeout: 120 * time.Second}
	resp, err := client.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("download failed with status %d", resp.StatusCode)
	}

	out, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, resp.Body)
	return err
}

func extractBinary(archivePath, destDir string) (string, error) {
	f, err := os.Open(archivePath)
	if err != nil {
		return "", err
	}
	defer f.Close()

	gz, err := gzip.NewReader(f)
	if err != nil {
		return "", fmt.Errorf("opening gzip: %w", err)
	}
	defer gz.Close()

	tr := tar.NewReader(gz)
	for {
		hdr, err := tr.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return "", fmt.Errorf("reading tar: %w", err)
		}

		// Look for the binary: could be "grafana", "grafana-cli", or at the
		// top level of the archive. Accept any executable-looking file.
		name := filepath.Base(hdr.Name)
		if name != "grafana" && name != "grafana-cli" {
			continue
		}
		if hdr.Typeflag != tar.TypeReg {
			continue
		}

		outPath := filepath.Join(destDir, name)
		out, err := os.OpenFile(outPath, os.O_CREATE|os.O_WRONLY, 0o755)
		if err != nil {
			return "", err
		}
		if _, err := io.Copy(out, tr); err != nil {
			out.Close()
			return "", err
		}
		out.Close()
		return outPath, nil
	}

	return "", fmt.Errorf("binary not found in archive")
}

func atomicReplace(src, dst string) error {
	// Preserve the permissions of the destination file.
	dstInfo, err := os.Stat(dst)
	if err != nil {
		return fmt.Errorf("stat destination: %w", err)
	}
	dstMode := dstInfo.Mode()

	// Create a temporary file in the same directory as the destination
	// so that os.Rename works (same filesystem).
	dstDir := filepath.Dir(dst)
	tmpFile, err := os.CreateTemp(dstDir, ".grafana-update-*")
	if err != nil {
		return fmt.Errorf("creating temp file: %w", err)
	}
	tmpPath := tmpFile.Name()

	// Clean up the temp file on error.
	defer func() {
		if tmpPath != "" {
			os.Remove(tmpPath)
		}
	}()

	// Copy the new binary to the temp file.
	srcFile, err := os.Open(src)
	if err != nil {
		tmpFile.Close()
		return fmt.Errorf("opening new binary: %w", err)
	}

	if _, err := io.Copy(tmpFile, srcFile); err != nil {
		srcFile.Close()
		tmpFile.Close()
		return fmt.Errorf("copying new binary: %w", err)
	}
	srcFile.Close()

	if err := tmpFile.Chmod(dstMode); err != nil {
		tmpFile.Close()
		return fmt.Errorf("setting permissions: %w", err)
	}

	if err := tmpFile.Close(); err != nil {
		return fmt.Errorf("closing temp file: %w", err)
	}

	// Atomic rename.
	if err := os.Rename(tmpPath, dst); err != nil {
		return fmt.Errorf("renaming: %w (you may need to run with sudo)", err)
	}

	// Clear tmpPath so the deferred cleanup doesn't remove the installed binary.
	tmpPath = ""
	return nil
}

func formatVer(v string) string {
	if strings.HasPrefix(v, "v") {
		return v
	}
	return "v" + v
}
