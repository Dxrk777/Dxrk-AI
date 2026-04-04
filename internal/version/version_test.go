package version

import (
	"path/filepath"
	"testing"
)

func TestVersionString(t *testing.T) {
	// Test v1.0 format
	v := Version{
		Major:      1.0,
		Percentage: 100.0,
	}
	if got := v.String(); got != "v1.0.0" {
		t.Errorf("String() = %v, want v1.0.0", got)
	}

	// Test percentage format for older versions
	v.Major = 0
	v.Percentage = 50.5
	if got := v.String(); got != "0050.50%" {
		t.Errorf("String() = %v, want 0050.50%%", got)
	}
}

func TestParseVersion(t *testing.T) {
	tests := []struct {
		input    string
		expected float64
		wantErr  bool
	}{
		{"000.01%", 0.01, false},
		{"50.00%", 50.0, false},
		{"100.00%", 100.0, false},
		{"0.01%", 0.01, false},
		{"invalid", 0, true},
	}

	for _, tt := range tests {
		got, err := ParseVersion(tt.input)
		if (err != nil) != tt.wantErr {
			t.Errorf("ParseVersion(%q) error = %v, wantErr %v", tt.input, err, tt.wantErr)
			continue
		}
		if !tt.wantErr && got != tt.expected {
			t.Errorf("ParseVersion(%q) = %v, want %v", tt.input, got, tt.expected)
		}
	}
}

func TestMilestones(t *testing.T) {
	m := Milestones()

	if len(m) == 0 {
		t.Error("Milestones() returned empty map")
	}

	if m["v1.0.0 — Initial Release"] != 100.0 {
		t.Errorf("v1.0.0 — Initial Release = %v, want 100.0", m["v1.0.0 — Initial Release"])
	}
}

func TestVersionBar(t *testing.T) {
	// Test v1.0 bar
	v := Version{Major: 1.0, Milestone: "v1.0.0 — Initial Release"}
	bar := v.Bar()

	if len(bar) == 0 {
		t.Error("Bar() returned empty string")
	}

	// Check it contains v1.0
	if !contains(bar, "v1.0") {
		t.Errorf("Bar() = %v, should contain v1.0", bar)
	}
}

func TestProgress(t *testing.T) {
	v := Version{
		Percentage: 100.0,
		Milestone:  "v1.0.0 — Initial Release",
	}

	current, total, label := v.Progress()
	if current != 100.0 {
		t.Errorf("current = %v, want 100", current)
	}
	if total != 100.0 {
		t.Errorf("total = %v, want 100", total)
	}
	if label != "v1.0.0 — Initial Release" {
		t.Errorf("label = %v, want v1.0.0 — Initial Release", label)
	}
}

func TestDefault(t *testing.T) {
	v := Default()
	if v.Percentage != 100.0 {
		t.Errorf("Default.Percentage = %v, want 100.0", v.Percentage)
	}
	if v.Major != 1.0 {
		t.Errorf("Default.Major = %v, want 1.0", v.Major)
	}
	if v.Milestone != "v1.0.0 — Initial Release" {
		t.Errorf("Default.Milestone = %v, want v1.0.0 — Initial Release", v.Milestone)
	}
}

func TestManager(t *testing.T) {
	// Create temp dir
	tmpDir := t.TempDir()

	m := &Manager{
		filePath: filepath.Join(tmpDir, "version.json"),
		data:     Default(),
	}

	// Test save/load
	if err := m.save(); err != nil {
		t.Fatalf("save() error = %v", err)
	}

	// Load should work
	m2 := &Manager{filePath: m.filePath}
	if err := m2.load(); err != nil {
		t.Fatalf("load() error = %v", err)
	}

	if m2.data.Percentage != m.data.Percentage {
		t.Errorf("loaded Percentage = %v, want %v", m2.data.Percentage, m.data.Percentage)
	}
}

func TestReset(t *testing.T) {
	tmpDir := t.TempDir()

	m := &Manager{
		filePath: filepath.Join(tmpDir, "version.json"),
		data: Version{
			Percentage: 50.0,
			Milestone:  "Some Milestone",
		},
	}

	if err := m.save(); err != nil {
		t.Fatalf("save() error = %v", err)
	}

	if err := m.Reset(); err != nil {
		t.Fatalf("Reset() error = %v", err)
	}

	if m.data.Percentage != 100.0 {
		t.Errorf("After Reset, Percentage = %v, want 100.0", m.data.Percentage)
	}
}

func TestNewManager(t *testing.T) {
	// This will create in home dir, which may or may not exist
	// Just verify it doesn't crash
	m, err := NewManager()
	if err != nil {
		t.Skipf("NewManager() error = %v (home dir may not exist)", err)
	}

	v := m.Get()
	if v.Percentage <= 0 {
		t.Errorf("Get() returned invalid version: %v", v)
	}
}

func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(s) > 0 && containsHelper(s, substr))
}

func containsHelper(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
