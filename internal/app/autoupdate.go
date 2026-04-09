package app

import (
	"archive/tar"
	"compress/gzip"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"syscall"
)

// AutoUpdateCheck checks for updates and returns true if a new version was
// successfully installed AND re-exec'd. It only runs on released builds.
//
// Guard order (mirrors selfUpdate to avoid double-triggering):
//  1. DXRK_SELF_UPDATE_DONE=1 → skip (loop guard set by selfUpdate or previous AutoUpdate)
//  2. version == "dev"         → skip
//  3. same version             → skip
func AutoUpdateCheck(stdout io.Writer) (updated bool, err error) {
	// Guard 1: loop prevention — selfUpdate (or a previous AutoUpdateCheck) already ran.
	if os.Getenv(envSelfUpdateDone) == "1" {
		return false, nil
	}

	// Guard 2: dev build — no meaningful version to compare.
	if Version == "dev" {
		return false, nil
	}

	// Get latest release tag.
	latestTag, err := getLatestTag()
	if err != nil {
		return false, nil // network failure is non-fatal
	}

	// Guard 3: already up-to-date.
	normalizedLatest := strings.TrimPrefix(latestTag, "v")
	normalizedCurrent := strings.TrimPrefix(Version, "v")
	if normalizedLatest == normalizedCurrent {
		return false, nil
	}

	_, _ = fmt.Fprintf(stdout, "🔄 Updating dxrk: %s → %s\n", normalizedCurrent, normalizedLatest)

	// Download and extract the new binary into a temp file.
	newBinaryPath, err := downloadAndExtractBinary(latestTag)
	if err != nil {
		_, _ = fmt.Fprintf(stdout, "⚠️  Update failed (download): %v\n", err)
		return false, err
	}
	defer os.Remove(newBinaryPath)

	// Replace the running binary on disk.
	installedPath, err := installBinary(newBinaryPath)
	if err != nil {
		_, _ = fmt.Fprintf(stdout, "⚠️  Update failed (install): %v\n", err)
		return false, err
	}

	_, _ = fmt.Fprintf(stdout, "✅ Updated to %s, restarting...\n", normalizedLatest)

	// Set loop guard BEFORE re-exec so the new process skips the update check.
	os.Setenv(envSelfUpdateDone, "1")

	// Re-exec with the new binary (Unix). On Windows just notify the user.
	if runtime.GOOS == "windows" {
		_, _ = fmt.Fprintf(stdout, "Updated to v%s — please restart.\n", normalizedLatest)
		return true, nil
	}

	// Use the installed path directly so we don't rely on PATH being refreshed.
	if err := syscall.Exec(installedPath, os.Args, os.Environ()); err != nil {
		// syscall.Exec only returns on error.
		_, _ = fmt.Fprintf(stdout, "⚠️  Re-exec failed: %v — please restart manually.\n", err)
		return true, nil // binary was updated; restart is up to the user
	}

	// Unreachable after a successful Exec.
	return true, nil
}

// getLatestTag fetches the latest release tag from GitHub.
func getLatestTag() (string, error) {
	url := "https://api.github.com/repos/Dxrk777/Dxrk-AI/releases/latest"

	resp, err := http.Get(url) //nolint:gosec // URL is static
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("GitHub API returned %s", resp.Status)
	}

	var release struct {
		TagName string `json:"tag_name"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&release); err != nil {
		return "", err
	}
	if release.TagName == "" {
		return "", fmt.Errorf("tag_name is empty in GitHub response")
	}

	return release.TagName, nil
}

// getArchiveName returns the GoReleaser archive filename for the current platform.
// Format: dxrk_<version>_<os>_<arch>.tar.gz  (or .zip on Windows)
func getArchiveName(tag string) string {
	goos := runtime.GOOS
	goarch := runtime.GOARCH
	version := strings.TrimPrefix(tag, "v")

	ext := ".tar.gz"
	if goos == "windows" {
		ext = ".zip"
	}

	return fmt.Sprintf("dxrk_%s_%s_%s%s", version, goos, goarch, ext)
}

// downloadAndExtractBinary downloads the release archive for the current
// platform, extracts the "dxrk" (or "dxrk.exe") binary, writes it to a
// temporary file, and returns its path.
//
// FIX: the previous implementation wrote the raw .tar.gz bytes to the binary
// path, producing a corrupt executable. This version extracts the binary first.
func downloadAndExtractBinary(tag string) (string, error) {
	archiveName := getArchiveName(tag)
	downloadURL := fmt.Sprintf(
		"https://github.com/Dxrk777/Dxrk-AI/releases/download/%s/%s",
		tag, archiveName,
	)

	resp, err := http.Get(downloadURL) //nolint:gosec // URL is built from validated tag
	if err != nil {
		return "", fmt.Errorf("download %s: %w", downloadURL, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("download failed: %s", resp.Status)
	}

	// Create a temp file for the extracted binary.
	tmpFile, err := os.CreateTemp("", "dxrk-update-*")
	if err != nil {
		return "", fmt.Errorf("create temp file: %w", err)
	}
	tmpFile.Close()

	binaryName := "dxrk"
	if runtime.GOOS == "windows" {
		binaryName = "dxrk.exe"
	}

	if err := extractBinaryFromTarGz(resp.Body, binaryName, tmpFile.Name()); err != nil {
		os.Remove(tmpFile.Name())
		return "", fmt.Errorf("extract binary: %w", err)
	}

	return tmpFile.Name(), nil
}

// extractBinaryFromTarGz reads a .tar.gz stream, finds the named binary entry,
// and writes it (mode 0755) to destPath.
func extractBinaryFromTarGz(r io.Reader, binaryName, destPath string) error {
	gz, err := gzip.NewReader(r)
	if err != nil {
		return fmt.Errorf("gzip reader: %w", err)
	}
	defer gz.Close()

	tr := tar.NewReader(gz)
	for {
		hdr, err := tr.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return fmt.Errorf("tar next: %w", err)
		}

		// Match the bare filename (archives may include directory prefix).
		if filepath.Base(hdr.Name) != binaryName {
			continue
		}

		out, err := os.OpenFile(destPath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0o755)
		if err != nil {
			return fmt.Errorf("open dest: %w", err)
		}

		if _, err := io.Copy(out, tr); err != nil {
			out.Close()
			return fmt.Errorf("copy binary: %w", err)
		}
		return out.Close()
	}

	return fmt.Errorf("%q not found in archive", binaryName)
}

// installBinary replaces the running dxrk binary with newBinaryPath.
// Returns the path where it was installed.
func installBinary(newBinaryPath string) (string, error) {
	// Resolve current binary location.
	currentPath, err := resolveCurrentBinary()
	if err != nil {
		return "", err
	}

	data, err := os.ReadFile(newBinaryPath)
	if err != nil {
		return "", fmt.Errorf("read new binary: %w", err)
	}

	// Write atomically: write to a temp sibling, then rename.
	dir := filepath.Dir(currentPath)
	tmp, err := os.CreateTemp(dir, ".dxrk-update-*")
	if err != nil {
		// Fall back to direct write if temp in same dir fails (e.g. read-only fs).
		if err2 := os.WriteFile(currentPath, data, 0o755); err2 != nil {
			return "", fmt.Errorf("write binary: %w", err2)
		}
		return currentPath, nil
	}

	if _, err := tmp.Write(data); err != nil {
		tmp.Close()
		os.Remove(tmp.Name())
		return "", fmt.Errorf("write temp binary: %w", err)
	}
	if err := tmp.Close(); err != nil {
		os.Remove(tmp.Name())
		return "", fmt.Errorf("close temp binary: %w", err)
	}

	if err := os.Chmod(tmp.Name(), 0o755); err != nil {
		os.Remove(tmp.Name())
		return "", fmt.Errorf("chmod temp binary: %w", err)
	}

	if err := os.Rename(tmp.Name(), currentPath); err != nil {
		os.Remove(tmp.Name())
		return "", fmt.Errorf("replace binary (rename): %w", err)
	}

	return currentPath, nil
}

// resolveCurrentBinary returns the path to the currently running dxrk binary.
func resolveCurrentBinary() (string, error) {
	// Prefer os.Executable so we always replace the binary that is actually
	// running, not whichever one happens to be first in PATH.
	if exe, err := os.Executable(); err == nil && exe != "" {
		// Resolve symlinks so we write to the real file (important on Homebrew).
		if real, err := filepath.EvalSymlinks(exe); err == nil {
			return real, nil
		}
		return exe, nil
	}

	// Fallback: search well-known locations.
	candidates := []string{
		"/usr/local/bin/dxrk",
		"/opt/homebrew/bin/dxrk",
		"/usr/bin/dxrk",
	}
	for _, p := range candidates {
		if _, err := os.Stat(p); err == nil {
			return p, nil
		}
	}

	return "", fmt.Errorf("cannot locate dxrk binary")
}
