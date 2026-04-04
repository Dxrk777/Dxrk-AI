// Package brain implements the central orchestrator for Dxrk Hell.
//
// The Brain coordinates all modules in Dxrk Hell:
// - Memory (persistent storage with history)
// - Vault (encryption for sensitive data)
// - Connector (remote control via Telegram, Discord, WhatsApp)
// - Commands (execute shell commands)
// - Email (send email notifications)
//
// Architecture:
// - Unified command processing
// - Memory integration with Engram sync
// - Secure command execution with timeout
// - Email notifications via SMTP
//
// This module is lightweight and does NOT include:
// - AI inference engines (Dxrk Hell uses external agents)
// - Vision/Imaging (not applicable to agent installer)
// - RAG/Embeddings (handled by Engram)
package brain

import (
	"fmt"
	"sync"
	"time"
)

// Config holds brain configuration.
type Config struct {
	// Memory settings
	MemoryDir string

	// Vault settings
	VaultEnabled bool

	// Connector settings
	ConnectorPort int

	// Email settings
	EmailEnabled bool
	SMTPHost     string
	SMTPPort     int
	SMTPUser     string
	SMTPPassword string
	EmailFrom    string

	// Command execution
	CommandTimeout time.Duration
}

// Brain is the central orchestrator for Dxrk Hell.
// Coordinates all modules and provides unified command processing.
type Brain struct {
	mu      sync.RWMutex
	config  *Config
	startAt time.Time
	running bool

	// Sub-systems (initialized lazily)
	memory *Memory
	vault  interface {
		Encrypt([]byte) ([]byte, error)
		Decrypt([]byte) ([]byte, error)
	}
	connector interface {
		Start() error
		Stop()
		Status() map[string]interface{}
	}
}

// New creates a new Brain instance.
func New(cfg *Config) *Brain {
	if cfg.CommandTimeout == 0 {
		cfg.CommandTimeout = 30 * time.Second
	}

	return &Brain{
		config:  cfg,
		startAt: time.Now(),
		running: true,
	}
}

// Status returns the current brain status.
func (b *Brain) Status() map[string]interface{} {
	b.mu.RLock()
	defer b.mu.RUnlock()

	status := map[string]interface{}{
		"running": b.running,
		"uptime":  time.Since(b.startAt).String(),
		"memory":  b.memory != nil,
		"vault":   b.config.VaultEnabled,
		"email":   b.config.EmailEnabled,
	}

	if b.connector != nil {
		connStatus := b.connector.Status()
		for k, v := range connStatus {
			status["connector_"+k] = v
		}
	}

	return status
}

// Uptime returns how long the brain has been running.
func (b *Brain) Uptime() time.Duration {
	return time.Since(b.startAt)
}

// IsRunning returns whether the brain is running.
func (b *Brain) IsRunning() bool {
	b.mu.RLock()
	defer b.mu.RUnlock()
	return b.running
}

// Shutdown gracefully shuts down the brain.
func (b *Brain) Shutdown() error {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.running = false
	return nil
}

// String returns a string representation of the brain.
func (b *Brain) String() string {
	return fmt.Sprintf("Dxrk Hell Brain (running for %s)", b.Uptime())
}
