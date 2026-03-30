package version

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
)

// Version represents the Dxrk AI version as a percentage
type Version struct {
	Major      float64 `json:"major"`      // 0-100
	Minor      float64 `json:"minor"`      // sub-percent
	Percentage float64 `json:"percentage"` // calculated percentage
	Milestone  string  `json:"milestone"`  // current milestone
	UpdatedBy  string  `json:"updated_by"` // who/what updated
}

// Default returns the default starting version
func Default() Version {
	return Version{
		Major:      0.01,
		Minor:      0,
		Percentage: 0.01,
		Milestone:  "Initial Release",
		UpdatedBy:  "Dxrk AI",
	}
}

// String returns the version string in format "000.XX%"
func (v Version) String() string {
	return fmt.Sprintf("%07.2f%%", v.Percentage)
}

// Short returns short version like "000.01%"
func (v Version) Short() string {
	return fmt.Sprintf("%.2f%%", v.Percentage)
}

// Manager handles version persistence and updates
type Manager struct {
	mu       sync.RWMutex
	filePath string
	data     Version
}

// NewManager creates a new version manager
func NewManager() (*Manager, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return nil, fmt.Errorf("get home dir: %w", err)
	}

	dxrkDir := filepath.Join(homeDir, ".config", "dxrk")
	if err := os.MkdirAll(dxrkDir, 0755); err != nil {
		return nil, fmt.Errorf("create dxrk dir: %w", err)
	}

	m := &Manager{
		filePath: filepath.Join(dxrkDir, "version.json"),
	}

	if err := m.load(); err != nil {
		// If no version file, use default
		m.data = Default()
		if err := m.save(); err != nil {
			return nil, fmt.Errorf("save default version: %w", err)
		}
	}

	return m, nil
}

// Get returns the current version
func (m *Manager) Get() Version {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.data
}

// Increment increases the version by the given percentage
func (m *Manager) Increment(percent float64, milestone string) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if percent <= 0 || percent > 100 {
		return fmt.Errorf("invalid percentage: %f (must be 0.01-100)", percent)
	}

	m.data.Percentage += percent
	if m.data.Percentage > 100 {
		m.data.Percentage = 100
	}
	m.data.Major = m.data.Percentage
	m.data.Milestone = milestone
	m.data.UpdatedBy = "Dxrk AI"

	return m.save()
}

// SetMilestone updates the current milestone
func (m *Manager) SetMilestone(milestone string) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.data.Milestone = milestone
	return m.save()
}

// Reset resets version to initial state
func (m *Manager) Reset() error {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.data = Default()
	return m.save()
}

// Milestones returns predefined milestones with their percentages
func Milestones() map[string]float64 {
	return map[string]float64{
		"Initial Release":    0.01,
		"Core Installer":     1.0,
		"Agent Support":      5.0,
		"Skills System":      10.0,
		"MCP Integration":    15.0,
		"Memory System":      20.0,
		"SDD Workflow":       25.0,
		"Persona System":     30.0,
		"Theme System":       35.0,
		"Backup/Restore":     40.0,
		"Multi-Platform":     50.0,
		"E2E Testing":        55.0,
		"Documentation":      60.0,
		"Community Features": 65.0,
		"Performance":        70.0,
		"Security Audit":     75.0,
		"Stable Release":     80.0,
		"Production Ready":   90.0,
		"Feature Complete":   95.0,
		"MVP Achieved":       100.0,
	}
}

// Progress returns progress info
func (v Version) Progress() (current, total float64, label string) {
	total = 100.0
	current = v.Percentage
	label = v.Milestone
	return
}

// Bar returns an ASCII progress bar
func (v Version) Bar() string {
	progress := int(v.Percentage / 2) // 50 chars = 100%
	filled := strings.Repeat("█", progress)
	empty := strings.Repeat("░", 50-progress)
	return fmt.Sprintf("[%s%s] %.2f%%", filled, empty, v.Percentage)
}

func (m *Manager) load() error {
	data, err := os.ReadFile(m.filePath)
	if err != nil {
		if os.IsNotExist(err) {
			return fmt.Errorf("version file not found")
		}
		return err
	}

	return json.Unmarshal(data, &m.data)
}

func (m *Manager) save() error {
	data, err := json.MarshalIndent(m.data, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(m.filePath, data, 0644)
}

// ParseVersion parses a version string like "000.01%"
func ParseVersion(s string) (float64, error) {
	s = strings.TrimSpace(s)
	s = strings.TrimSuffix(s, "%")
	s = strings.TrimPrefix(s, "0")
	if s == "" {
		s = "0"
	}
	return strconv.ParseFloat(s, 64)
}

// FormatPercent formats a float as percentage string
func FormatPercent(f float64) string {
	return fmt.Sprintf("%.2f%%", f)
}
