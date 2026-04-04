// Package state implements persistent state and history for Dxrk Hell.
//
// The state package provides:
// - InstallState: Basic last-installed agents
// - History: Extended install history with timestamps and details
// - Memory: Integration with Engram for persistent agent memory
//
// Memory architecture:
// ┌─────────────┐     ┌─────────────┐
// │ state.json  │ ←→  │ history.json│
// │ (basic)    │     │ (extended)  │
// └──────┬──────┘     └──────┬──────┘
//
//	│                    │
//	└────────┬───────────┘
//	         ↓
//	┌───────────────┐
//	│   memory.go   │
//	│ (integration) │
//	└───────┬───────┘
//	        ↓
//	┌───────────────┐
//	│ Engram CLI    │
//	│ (persistent)  │
//	└───────────────┘
package state

import (
	"encoding/json"
	"os"
	"os/exec"
	"path/filepath"
	"time"
)

const (
	// Directory where state files are stored
	stateDirName = ".dxrk"

	// State files
	stateFileName   = "state.json"
	historyFileName = "history.json"
)

// HistoryEntry represents a single install event in history.
type HistoryEntry struct {
	Timestamp   time.Time `json:"timestamp"`
	Agents      []string  `json:"agents"`
	Preset      string    `json:"preset"`
	Components  []string  `json:"components"`
	Success     bool      `json:"success"`
	DurationSec int       `json:"duration_sec"`
}

// History stores the complete install history.
type History struct {
	Entries       []HistoryEntry `json:"entries"`
	LastInstall   *HistoryEntry  `json:"last_install,omitempty"`
	TotalInstalls int            `json:"total_installs"`
	TotalAgents   int            `json:"total_agents_installed"`
}

// HistoryPath returns the path to the history file.
func HistoryPath(homeDir string) string {
	return filepath.Join(homeDir, stateDirName, historyFileName)
}

// ReadHistory reads the history file. Returns empty history if not exists.
func ReadHistory(homeDir string) (*History, error) {
	path := HistoryPath(homeDir)
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return &History{Entries: []HistoryEntry{}}, nil
		}
		return nil, err
	}

	var h History
	if err := json.Unmarshal(data, &h); err != nil {
		return nil, err
	}
	return &h, nil
}

// WriteHistory saves the history to file.
func WriteHistory(homeDir string, history *History) error {
	dir := filepath.Join(homeDir, stateDirName)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}

	data, err := json.MarshalIndent(history, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(HistoryPath(homeDir), append(data, '\n'), 0644)
}

// AddInstall adds a new install entry to history.
func AddInstall(homeDir string, agents []string, preset string, components []string, success bool, durationSec int) error {
	entry := HistoryEntry{
		Timestamp:   time.Now(),
		Agents:      agents,
		Preset:      preset,
		Components:  components,
		Success:     success,
		DurationSec: durationSec,
	}

	history, err := ReadHistory(homeDir)
	if err != nil {
		return err
	}

	// Add to entries (most recent first)
	history.Entries = append([]HistoryEntry{entry}, history.Entries...)
	history.LastInstall = &entry
	history.TotalInstalls++
	history.TotalAgents += len(agents)

	// Keep only last 100 entries
	if len(history.Entries) > 100 {
		history.Entries = history.Entries[:100]
	}

	return WriteHistory(homeDir, history)
}

// GetLastInstall returns the most recent install entry.
func GetLastInstall(homeDir string) (*HistoryEntry, error) {
	history, err := ReadHistory(homeDir)
	if err != nil {
		return nil, err
	}
	return history.LastInstall, nil
}

// GetFrequentAgents returns agents that have been installed multiple times.
func GetFrequentAgents(homeDir string) ([]string, error) {
	history, err := ReadHistory(homeDir)
	if err != nil {
		return nil, err
	}

	// Count agent occurrences
	counts := make(map[string]int)
	for _, entry := range history.Entries {
		for _, agent := range entry.Agents {
			counts[agent]++
		}
	}

	// Filter agents installed more than once
	var frequent []string
	for agent, count := range counts {
		if count > 1 {
			frequent = append(frequent, agent)
		}
	}
	return frequent, nil
}

// EngramSync saves install history to Engram for persistent memory.
func EngramSync(homeDir string) error {
	history, err := ReadHistory(homeDir)
	if err != nil {
		return err
	}

	if history.LastInstall == nil {
		return nil // Nothing to sync
	}

	// Try to use engram CLI if available
	_, err = exec.LookPath("engram")
	if err != nil {
		return nil // Engram not installed, skip
	}

	// Create a memory entry for the install
	cmd := exec.Command("engram", "memory", "add",
		"--type", "install",
		"--content", formatInstallMemory(history.LastInstall),
	)

	return cmd.Run()
}

// formatInstallMemory formats an install entry for Engram.
func formatInstallMemory(entry *HistoryEntry) string {
	agents := ""
	if len(entry.Agents) > 0 {
		agents = entry.Agents[0]
		for i := 1; i < len(entry.Agents); i++ {
			agents += ", " + entry.Agents[i]
		}
	}
	return "Installed: " + agents + " with preset: " + entry.Preset
}

// GetStats returns memory statistics.
func GetStats(homeDir string) (map[string]any, error) {
	history, err := ReadHistory(homeDir)
	if err != nil {
		return nil, err
	}

	stats := map[string]any{
		"total_installs":      history.TotalInstalls,
		"total_agents":        history.TotalAgents,
		"unique_agents":       countUniqueAgents(history.Entries),
		"successful_installs": countSuccessfulInstalls(history.Entries),
		"last_install":        history.LastInstall,
	}

	return stats, nil
}

func countUniqueAgents(entries []HistoryEntry) int {
	set := make(map[string]bool)
	for _, entry := range entries {
		for _, agent := range entry.Agents {
			set[agent] = true
		}
	}
	return len(set)
}

func countSuccessfulInstalls(entries []HistoryEntry) int {
	count := 0
	for _, entry := range entries {
		if entry.Success {
			count++
		}
	}
	return count
}
