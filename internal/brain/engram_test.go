package brain

import (
	"testing"
)

func TestDefaultEngramConfig(t *testing.T) {
	cfg := DefaultEngramConfig()

	if cfg.URL != "http://localhost:8080" {
		t.Errorf("Expected URL http://localhost:8080, got %s", cfg.URL)
	}

	if cfg.Timeout != 10*1000000000 {
		t.Errorf("Expected timeout 10s, got %v", cfg.Timeout)
	}
}

func TestNewEngramClient(t *testing.T) {
	cfg := &EngramConfig{
		URL:     "http://localhost:9090",
		Timeout: 5 * 1000000000,
	}

	client := NewEngramClient(cfg)

	if client.config.URL != "http://localhost:9090" {
		t.Errorf("Expected URL http://localhost:9090, got %s", client.config.URL)
	}
}

func TestNewEngramClientNilConfig(t *testing.T) {
	client := NewEngramClient(nil)

	if client.config.URL != "http://localhost:8080" {
		t.Errorf("Expected default URL, got %s", client.config.URL)
	}
}

func TestEngramClientString(t *testing.T) {
	client := NewEngramClient(&EngramConfig{
		URL: "http://test:1234",
	})

	str := client.String()
	if str != "EngramClient{URL: http://test:1234}" {
		t.Errorf("Unexpected string: %s", str)
	}
}

func TestEngramClientIsAvailable(t *testing.T) {
	client := NewEngramClient(&EngramConfig{
		URL: "http://localhost:99999", // Invalid port
	})

	// Should return false because server is not running
	available := client.IsAvailable()
	if available {
		t.Error("Expected false for unavailable server")
	}
}

func TestEngramMemoryEntry(t *testing.T) {
	entry := EngramMemoryEntry{
		ID:      "test-id",
		Content: "test content",
		Type:    "command",
		Metadata: map[string]interface{}{
			"key": "value",
		},
	}

	if entry.ID != "test-id" {
		t.Errorf("Expected ID test-id, got %s", entry.ID)
	}

	if entry.Content != "test content" {
		t.Errorf("Expected content, got %s", entry.Content)
	}

	if entry.Type != "command" {
		t.Errorf("Expected type command, got %s", entry.Type)
	}

	if entry.Metadata["key"] != "value" {
		t.Errorf("Expected metadata key=value")
	}
}

func TestEngramClientRemember(t *testing.T) {
	client := NewEngramClient(&EngramConfig{
		URL: "http://localhost:99999",
	})

	entry := MemoryEntry{
		Content: "test",
		Type:    "command",
	}

	err := client.Remember(entry)
	if err == nil {
		t.Error("Expected error for unreachable server")
	}
}

func TestEngramClientQuery(t *testing.T) {
	client := NewEngramClient(&EngramConfig{
		URL: "http://localhost:99999",
	})

	entries, err := client.Query("test")
	if err == nil && len(entries) == 0 {
		t.Log("Query returned empty (server might be unreachable)")
	}
}

func TestEngramClientRecent(t *testing.T) {
	client := NewEngramClient(&EngramConfig{
		URL: "http://localhost:99999",
	})

	entries, err := client.Recent(10)
	if err == nil && len(entries) == 0 {
		t.Log("Recent returned empty (expected for unreachable server)")
	}
}

func TestEngramClientRecentWithLimit(t *testing.T) {
	client := NewEngramClient(&EngramConfig{
		URL: "http://localhost:99999",
	})

	// Test with different limits
	limits := []int{1, 5, 10, 50, 100}
	for _, limit := range limits {
		entries, err := client.Recent(limit)
		if err != nil {
			t.Logf("Recent(%d) error: %v (expected for unreachable server)", limit, err)
		} else if len(entries) == 0 {
			t.Logf("Recent(%d) returned empty", limit)
		}
	}
}

func TestEngramClientQueryEmpty(t *testing.T) {
	client := NewEngramClient(&EngramConfig{
		URL: "http://localhost:99999",
	})

	entries, err := client.Query("")
	if err == nil {
		t.Log("Empty query might return all entries")
	}
	_ = entries // suppress unused warning
}

func TestEngramClientQuerySpecialChars(t *testing.T) {
	client := NewEngramClient(&EngramConfig{
		URL: "http://localhost:99999",
	})

	entries, err := client.Query("test with spaces and @#$%")
	if err == nil {
		t.Logf("Query with special chars returned %d entries", len(entries))
	}
}
