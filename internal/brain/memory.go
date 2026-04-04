package brain

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"
)

// MemoryEntry represents a single memory entry.
type MemoryEntry struct {
	ID        string                 `json:"id"`
	Type      string                 `json:"type"` // "command", "install", "error", "info"
	Content   string                 `json:"content"`
	Metadata  map[string]interface{} `json:"metadata,omitempty"`
	Timestamp time.Time              `json:"timestamp"`
}

// Memory manages persistent memory for Dxrk Hell.
// Tracks commands, installations, errors, and other events.
type Memory struct {
	mu      sync.RWMutex
	dir     string
	entries []MemoryEntry
}

// NewMemory creates a new Memory instance.
func NewMemory(dataDir string) (*Memory, error) {
	mem := &Memory{
		dir:     filepath.Join(dataDir, ".dxrk", "memory"),
		entries: make([]MemoryEntry, 0),
	}

	// Ensure directory exists
	if err := os.MkdirAll(mem.dir, 0755); err != nil {
		return nil, err
	}

	// Load existing entries
	if err := mem.load(); err != nil {
		// Log but don't fail - start fresh
	}

	return mem, nil
}

// Remember adds a new memory entry.
func (m *Memory) Remember(entry MemoryEntry) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if entry.ID == "" {
		entry.ID = generateID()
	}
	if entry.Timestamp.IsZero() {
		entry.Timestamp = time.Now()
	}

	m.entries = append(m.entries, entry)
	return m.save()
}

// Query searches memory for entries matching the query.
func (m *Memory) Query(query string) []MemoryEntry {
	m.mu.RLock()
	defer m.mu.RUnlock()

	results := make([]MemoryEntry, 0)
	for _, entry := range m.entries {
		if containsIgnoreCase(entry.Content, query) {
			results = append(results, entry)
		}
	}
	return results
}

// Recent returns the most recent entries.
func (m *Memory) Recent(limit int) []MemoryEntry {
	m.mu.RLock()
	defer m.mu.RUnlock()

	if limit > len(m.entries) {
		limit = len(m.entries)
	}

	// Return last 'limit' entries
	start := len(m.entries) - limit
	if start < 0 {
		start = 0
	}

	result := make([]MemoryEntry, limit)
	copy(result, m.entries[start:])
	return result
}

// History returns all memory entries (for backward compatibility).
func (m *Memory) History() []MemoryEntry {
	m.mu.RLock()
	defer m.mu.RUnlock()

	result := make([]MemoryEntry, len(m.entries))
	copy(result, m.entries)
	return result
}

// Clear removes all memory entries.
func (m *Memory) Clear() error {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.entries = make([]MemoryEntry, 0)
	return m.save()
}

// Len returns the number of memory entries.
func (m *Memory) Len() int {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return len(m.entries)
}

// load reads memory from disk.
func (m *Memory) load() error {
	file := filepath.Join(m.dir, "memory.json")

	data, err := os.ReadFile(file)
	if err != nil {
		if os.IsNotExist(err) {
			return nil // No file yet, that's ok
		}
		return err
	}

	return json.Unmarshal(data, &m.entries)
}

// save writes memory to disk.
func (m *Memory) save() error {
	file := filepath.Join(m.dir, "memory.json")

	data, err := json.MarshalIndent(m.entries, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(file, data, 0644)
}

// generateID creates a unique ID for memory entries.
func generateID() string {
	return fmt.Sprintf("%d", time.Now().UnixNano())
}

// containsIgnoreCase checks if s contains substr (case-insensitive).
func containsIgnoreCase(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr ||
		containsLower(toLower(s), toLower(substr)))
}

func toLower(s string) string {
	result := make([]byte, len(s))
	for i := 0; i < len(s); i++ {
		c := s[i]
		if c >= 'A' && c <= 'Z' {
			c += 'a' - 'A'
		}
		result[i] = c
	}
	return string(result)
}

func containsLower(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}

var _ = filepath.Join // suppress unused import warning
