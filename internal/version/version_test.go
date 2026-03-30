package version

import (
	"path/filepath"
	"testing"
)

func TestVersionString(t *testing.T) {
	v := Version{
		Percentage: 0.01,
	}
	if got := v.String(); got != "0000.01%" {
		t.Errorf("String() = %v, want 0000.01%%", got)
	}

	v.Percentage = 50.5
	if got := v.String(); got != "0050.50%" {
		t.Errorf("String() = %v, want 0050.50%%", got)
	}

	v.Percentage = 100
	if got := v.String(); got != "0100.00%" {
		t.Errorf("String() = %v, want 0100.00%%", got)
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

	if m["Initial Release"] != 0.01 {
		t.Errorf("Initial Release = %v, want 0.01", m["Initial Release"])
	}

	if m["MVP Achieved"] != 100.0 {
		t.Errorf("MVP Achieved = %v, want 100", m["MVP Achieved"])
	}
}

func TestVersionBar(t *testing.T) {
	v := Version{Percentage: 50.0}
	bar := v.Bar()

	if len(bar) == 0 {
		t.Error("Bar() returned empty string")
	}

	// Check it contains the percentage
	expected := "50.00"
	if !contains(bar, expected) {
		t.Errorf("Bar() = %v, should contain %v", bar, expected)
	}
}

func TestProgress(t *testing.T) {
	v := Version{
		Percentage: 75.0,
		Milestone:  "Security Audit",
	}

	current, total, label := v.Progress()
	if current != 75.0 {
		t.Errorf("current = %v, want 75", current)
	}
	if total != 100.0 {
		t.Errorf("total = %v, want 100", total)
	}
	if label != "Security Audit" {
		t.Errorf("label = %v, want Security Audit", label)
	}
}

func TestDefault(t *testing.T) {
	v := Default()
	if v.Percentage != 0.01 {
		t.Errorf("Default.Percentage = %v, want 0.01", v.Percentage)
	}
	if v.Milestone != "Initial Release" {
		t.Errorf("Default.Milestone = %v, want Initial Release", v.Milestone)
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

	// Test increment
	if err := m.Increment(10.0, "Core Installer"); err != nil {
		t.Fatalf("Increment() error = %v", err)
	}

	if m.data.Percentage != 10.01 {
		t.Errorf("After Increment, Percentage = %v, want 10.01", m.data.Percentage)
	}

	if m.data.Milestone != "Core Installer" {
		t.Errorf("After Increment, Milestone = %v, want Core Installer", m.data.Milestone)
	}

	// Test cap at 100
	if err := m.Increment(100.0, "Overflow"); err != nil {
		t.Fatalf("Increment() error = %v", err)
	}
	if m.data.Percentage != 100 {
		t.Errorf("After large Increment, Percentage = %v, want 100", m.data.Percentage)
	}

	// Test invalid increment
	if err := m.Increment(0, "Invalid"); err == nil {
		t.Error("Increment(0) should return error")
	}

	if err := m.Increment(-1, "Invalid"); err == nil {
		t.Error("Increment(-1) should return error")
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

	if m.data.Percentage != 0.01 {
		t.Errorf("After Reset, Percentage = %v, want 0.01", m.data.Percentage)
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
