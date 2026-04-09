package brain_test

import (
	"testing"

	"github.com/Dxrk777/Dxrk-AI/internal/brain"
)

func TestMemoryNew(t *testing.T) {
	mem, err := brain.NewMemory(t.TempDir())
	if err != nil {
		t.Fatalf("NewMemory failed: %v", err)
	}
	if mem == nil {
		t.Fatal("Memory should not be nil")
	}
	if mem.Len() != 0 {
		t.Error("New memory should be empty")
	}
}

func TestMemoryRemember(t *testing.T) {
	mem, _ := brain.NewMemory(t.TempDir())

	entry := brain.MemoryEntry{
		Type:    "test",
		Content: "Test entry",
	}

	if err := mem.Remember(entry); err != nil {
		t.Fatalf("Remember failed: %v", err)
	}

	if mem.Len() != 1 {
		t.Errorf("Memory should have 1 entry, got %d", mem.Len())
	}
}

func TestMemoryQuery(t *testing.T) {
	mem, _ := brain.NewMemory(t.TempDir())

	// Add entries
	mem.Remember(brain.MemoryEntry{Type: "command", Content: "run ls -la"})
	mem.Remember(brain.MemoryEntry{Type: "install", Content: "installed opencode"})
	mem.Remember(brain.MemoryEntry{Type: "command", Content: "run git status"})

	// Query
	results := mem.Query("run")
	if len(results) != 2 {
		t.Errorf("Query should find 2 entries, got %d", len(results))
	}

	// Query with no results
	results = mem.Query("nonexistent")
	if len(results) != 0 {
		t.Errorf("Query should find 0 entries, got %d", len(results))
	}
}

func TestMemoryRecent(t *testing.T) {
	mem, _ := brain.NewMemory(t.TempDir())

	// Add 5 entries
	for i := 0; i < 5; i++ {
		mem.Remember(brain.MemoryEntry{
			Type:    "test",
			Content: "Entry " + string(rune('0'+i)),
		})
	}

	// Get recent 3
	recent := mem.Recent(3)
	if len(recent) != 3 {
		t.Errorf("Recent should return 3 entries, got %d", len(recent))
	}
}

func TestMemoryHistory(t *testing.T) {
	mem, _ := brain.NewMemory(t.TempDir())

	mem.Remember(brain.MemoryEntry{Type: "a", Content: "1"})
	mem.Remember(brain.MemoryEntry{Type: "b", Content: "2"})
	mem.Remember(brain.MemoryEntry{Type: "c", Content: "3"})

	history := mem.History()
	if len(history) != 3 {
		t.Errorf("History should return 3 entries, got %d", len(history))
	}
}

func TestMemoryClear(t *testing.T) {
	mem, _ := brain.NewMemory(t.TempDir())

	mem.Remember(brain.MemoryEntry{Type: "test", Content: "test"})
	if mem.Len() != 1 {
		t.Fatal("Should have 1 entry before clear")
	}

	if err := mem.Clear(); err != nil {
		t.Fatalf("Clear failed: %v", err)
	}

	if mem.Len() != 0 {
		t.Errorf("Memory should be empty after clear, got %d", mem.Len())
	}
}
