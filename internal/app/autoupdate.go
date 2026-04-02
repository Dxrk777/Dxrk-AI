package app

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"runtime"
	"strings"
)

// AutoUpdateCheck checks for updates and returns true if updated.
// It only downloads and installs if running as a released version (not "dev").
func AutoUpdateCheck(stdout io.Writer) (updated bool, err error) {
	// Don't auto-update in dev mode
	if Version == "dev" {
		return false, nil
	}

	// Get latest release
	latestTag, err := getLatestTag()
	if err != nil {
		return false, nil // Silently fail - don't bother user
	}

	// Check if update needed
	if latestTag == "v"+Version || latestTag == Version {
		return false, nil
	}

	// Update available!
	_, _ = fmt.Fprintf(stdout, "🔄 Updating dxrk: %s → %s\n", Version, strings.TrimPrefix(latestTag, "v"))

	// Download new binary
	filepath, err := downloadLatest(latestTag)
	if err != nil {
		_, _ = fmt.Fprintf(stdout, "⚠️  Update failed: %v\n", err)
		return false, err
	}
	defer os.Remove(filepath)

	// Install
	if err := install(filepath); err != nil {
		_, _ = fmt.Fprintf(stdout, "⚠️  Installation failed: %v\n", err)
		return false, err
	}

	_, _ = fmt.Fprintf(stdout, "✅ Updated to %s\n", strings.TrimPrefix(latestTag, "v"))
	return true, nil
}

func getLatestTag() (string, error) {
	url := "https://api.github.com/repos/Dxrk777/Dxrk-Hex/releases/latest"

	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	var release struct {
		TagName string `json:"tag_name"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&release); err != nil {
		return "", err
	}

	return release.TagName, nil
}

func downloadLatest(tag string) (string, error) {
	filename := getFilename(tag)
	url := fmt.Sprintf("https://github.com/Dxrk777/Dxrk-Hex/releases/download/%s/%s", tag, filename)

	// Create temp file
	tmpDir := os.TempDir()
	filepath := tmpDir + "/" + filename

	// Download
	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return "", fmt.Errorf("download failed: %s", resp.Status)
	}

	out, err := os.Create(filepath)
	if err != nil {
		return "", err
	}
	defer out.Close()

	if _, err := io.Copy(out, resp.Body); err != nil {
		return "", err
	}

	// Make executable
	os.Chmod(filepath, 0755)

	return filepath, nil
}

func getFilename(tag string) string {
	os := runtime.GOOS
	arch := runtime.GOARCH
	ext := ".tar.gz"

	if os == "windows" {
		ext = ".zip"
	}

	tag = strings.TrimPrefix(tag, "v")

	return fmt.Sprintf("dxrk_%s_%s_%s%s", tag, os, arch, ext)
}

func install(filepath string) error {
	// Find current binary
	currentPath, err := lookPath("dxrk")
	if err != nil {
		// Try common locations
		paths := []string{
			"/usr/local/bin/dxrk",
			"/opt/homebrew/bin/dxrk",
			"/usr/bin/dxrk",
		}
		for _, p := range paths {
			if _, err = os.Stat(p); err == nil {
				currentPath = p
				break
			}
		}
	}

	if currentPath == "" {
		return fmt.Errorf("cannot find dxrk in PATH")
	}

	// Check if we can write
	if _, err = os.OpenFile(currentPath, os.O_WRONLY, 0); err != nil {
		// Need sudo
		return fmt.Errorf("need sudo: chown +x %s", filepath)
	}

	// Simple copy (replace)
	data, err := os.ReadFile(filepath)
	if err != nil {
		return err
	}

	return os.WriteFile(currentPath, data, 0755)
}

func lookPath(name string) (string, error) {
	path := os.Getenv("PATH")
	for _, dir := range strings.Split(path, ":") {
		fullPath := dir + "/" + name
		if info, err := os.Stat(fullPath); err == nil && !info.IsDir() && info.Mode().IsRegular() {
			return fullPath, nil
		}
	}
	return "", fmt.Errorf("not found")
}
