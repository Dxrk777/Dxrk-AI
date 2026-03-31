package state

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestReadWriteHistory(t *testing.T) {
	// Create temp dir
	tmpDir, err := os.MkdirTemp("", "dxrk-test")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	// Test empty history
	history, err := ReadHistory(tmpDir)
	if err != nil {
		t.Fatalf("ReadHistory() error = %v", err)
	}
	if len(history.Entries) != 0 {
		t.Errorf("Expected empty history, got %d entries", len(history.Entries))
	}

	// Add an install
	err = AddInstall(tmpDir, []string{"claude", "opencode"}, "full", []string{"engram", "sdd"}, true, 120)
	if err != nil {
		t.Fatalf("AddInstall() error = %v", err)
	}

	// Read again
	history, err = ReadHistory(tmpDir)
	if err != nil {
		t.Fatalf("ReadHistory() error = %v", err)
	}
	if len(history.Entries) != 1 {
		t.Errorf("Expected 1 entry, got %d", len(history.Entries))
	}
	if history.TotalInstalls != 1 {
		t.Errorf("Expected TotalInstalls=1, got %d", history.TotalInstalls)
	}
	if history.TotalAgents != 2 {
		t.Errorf("Expected TotalAgents=2, got %d", history.TotalAgents)
	}
}

func TestGetLastInstall(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "dxrk-test")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	// No install yet
	last, err := GetLastInstall(tmpDir)
	if err != nil {
		t.Fatalf("GetLastInstall() error = %v", err)
	}
	if last != nil {
		t.Error("Expected nil for no installs")
	}

	// Add install
	AddInstall(tmpDir, []string{"claude"}, "minimal", []string{"engram"}, true, 60)

	last, err = GetLastInstall(tmpDir)
	if err != nil {
		t.Fatalf("GetLastInstall() error = %v", err)
	}
	if last == nil {
		t.Fatal("Expected last install to exist")
	}
	if len(last.Agents) != 1 || last.Agents[0] != "claude" {
		t.Errorf("Expected agent=claude, got %v", last.Agents)
	}
}

func TestGetFrequentAgents(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "dxrk-test")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	// Install clau   de twice
	AddInstall(tmpDir, []string{"claude"}, "full", nil, true, 60)
	AddInstall(tmpDir, []string{"opencode"}, "full", nil, true, 60)
	AddInstall(tmpDir, []string{"claude"}, "full", nil, true, 60)

	frequent, err := GetFrequentAgents(tmpDir)
	if err != nil {
		t.Fatalf("GetFrequentAgents() error = %v", err)
	}

	if len(frequent) != 1 {
		t.Errorf("Expected 1 frequent agent, got %d", len(frequent))
	}
	if len(frequent) > 0 && frequent[0] != "claude" {
		t.Errorf("Expected cla   de, got %s", frequent[0])
	}
}

func TestGetStats(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "dxrk-test")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	// Add installs
	AddInstall(tmpDir, []string{"claude", "opencode"}, "full", nil, true, 120)
	AddInstall(tmpDir, []string{"cursor"}, "minimal", nil, false, 60) // failed

	stats, err := GetStats(tmpDir)
	if err != nil {
		t.Fatalf("GetStats() error = %v", err)
	}

	if stats["total_installs"].(int) != 2 {
		t.Errorf("Expected total_installs=2, got %v", stats["total_installs"])
	}
	if stats["total_agents"].(int) != 3 {
		t.Errorf("Expected total_agents=3, got %v", stats["total_agents"])
	}
	if stats["unique_agents"].(int) != 3 {
		t.Errorf("Expected unique_agents=3, got %v", stats["unique_agents"])
	}
	if stats["successful_installs"].(int) != 1 {
		t.Errorf("Expected successful_installs=1, got %v", stats["successful_installs"])
	}
}

func TestHistoryPath(t *testing.T) {
	path := HistoryPath("/home/user")
	expected := filepath.Join("/home/user", stateDirName, historyFileName)
	if path != expected {
		t.Errorf("Expected %s, got %s", expected, path)
	}
}

func TestMultipleInstalls(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "dxrk-test")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	// Add 5 installs
	for i := 0; i < 5; i++ {
		AddInstall(tmpDir, []string{"agent"}, "full", nil, true, 60)
	}

	history, _ := ReadHistory(tmpDir)
	if len(history.Entries) != 5 {
		t.Errorf("Expected 5 entries, got %d", len(history.Entries))
	}
	if history.Entries[0].Timestamp.Before(history.Entries[4].Timestamp) {
		t.Error("Expected most recent first")
	}
}

func TestEngramSync_NoPanic(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "dxrk-test")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	// Should not panic even if engram is not installed
	err = EngramSync(tmpDir)
	if err != nil {
		t.Logf("EngramSync error (expected if engram not installed): %v", err)
	}
}

func TestHistoryEntryTiming(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "dxrk-test")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	before := time.Now()
	AddInstall(tmpDir, []string{"test"}, "full", nil, true, 30)
	after := time.Now()

	history, _ := ReadHistory(tmpDir)
	last := history.LastInstall

	if last.Timestamp.Before(before) || last.Timestamp.After(after) {
		t.Error("Timestamp not within expected range")
	}
}
