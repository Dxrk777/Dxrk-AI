package brain

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"
)

// EngramConfig holds Engram MCP server configuration.
type EngramConfig struct {
	// URL of the Engram MCP server (default: http://localhost:8080)
	URL string
	// APIKey for authentication (optional)
	APIKey string
	// Timeout for requests
	Timeout time.Duration
}

// DefaultEngramConfig returns default Engram configuration.
func DefaultEngramConfig() *EngramConfig {
	return &EngramConfig{
		URL:     "http://localhost:8080",
		Timeout: 10 * time.Second,
	}
}

// EngramClient is a client for the Engram MCP server.
type EngramClient struct {
	config *EngramConfig
	client *http.Client
}

// NewEngramClient creates a new Engram client.
func NewEngramClient(cfg *EngramConfig) *EngramClient {
	if cfg == nil {
		cfg = DefaultEngramConfig()
	}
	if cfg.Timeout == 0 {
		cfg.Timeout = 10 * time.Second
	}

	return &EngramClient{
		config: cfg,
		client: &http.Client{
			Timeout: cfg.Timeout,
		},
	}
}

// IsAvailable checks if the Engram server is running.
func (e *EngramClient) IsAvailable() bool {
	resp, err := e.client.Get(e.config.URL + "/health")
	if err != nil {
		return false
	}
	defer resp.Body.Close()
	return resp.StatusCode == http.StatusOK
}

// MemoryEntry represents an entry in Engram.
type EngramMemoryEntry struct {
	ID        string                 `json:"id,omitempty"`
	Content   string                 `json:"content"`
	Type      string                 `json:"type,omitempty"`
	Metadata  map[string]interface{} `json:"metadata,omitempty"`
	Timestamp time.Time              `json:"timestamp,omitempty"`
}

// Remember saves a memory entry to Engram.
func (e *EngramClient) Remember(entry MemoryEntry) error {
	engramEntry := EngramMemoryEntry{
		Content: entry.Content,
		Type:    entry.Type,
		Metadata: map[string]interface{}{
			"original_type": entry.Type,
		},
	}

	body, err := json.Marshal(engramEntry)
	if err != nil {
		return fmt.Errorf("failed to marshal entry: %w", err)
	}

	req, err := http.NewRequest("POST", e.config.URL+"/memory", strings.NewReader(string(body)))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	if e.config.APIKey != "" {
		req.Header.Set("Authorization", "Bearer "+e.config.APIKey)
	}

	resp, err := e.client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to connect to Engram: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		return fmt.Errorf("Engram returned status %d", resp.StatusCode)
	}

	return nil
}

// Query searches for memory entries in Engram.
func (e *EngramClient) Query(query string) ([]MemoryEntry, error) {
	url := fmt.Sprintf("%s/memory/search?q=%s", e.config.URL, query)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	if e.config.APIKey != "" {
		req.Header.Set("Authorization", "Bearer "+e.config.APIKey)
	}

	resp, err := e.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to Engram: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("Engram returned status %d", resp.StatusCode)
	}

	var results struct {
		Entries []EngramMemoryEntry `json:"entries"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&results); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	entries := make([]MemoryEntry, len(results.Entries))
	for i, e := range results.Entries {
		entries[i] = MemoryEntry{
			Content:  e.Content,
			Type:     e.Type,
			Metadata: e.Metadata,
		}
	}

	return entries, nil
}

// Recent returns recent memory entries from Engram.
func (e *EngramClient) Recent(limit int) ([]MemoryEntry, error) {
	url := fmt.Sprintf("%s/memory/recent?limit=%d", e.config.URL, limit)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	if e.config.APIKey != "" {
		req.Header.Set("Authorization", "Bearer "+e.config.APIKey)
	}

	resp, err := e.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to Engram: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("Engram returned status %d", resp.StatusCode)
	}

	var results struct {
		Entries []EngramMemoryEntry `json:"entries"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&results); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	entries := make([]MemoryEntry, len(results.Entries))
	for i, e := range results.Entries {
		entries[i] = MemoryEntry{
			Content:  e.Content,
			Type:     e.Type,
			Metadata: e.Metadata,
		}
	}

	return entries, nil
}

// String returns a string representation of the Engram client.
func (e *EngramClient) String() string {
	return fmt.Sprintf("EngramClient{URL: %s}", e.config.URL)
}
