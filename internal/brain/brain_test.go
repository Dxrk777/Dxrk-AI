package brain_test

import (
	"testing"
	"time"

	"github.com/Dxrk777/Dxrk-AI/internal/brain"
)

func TestBrainNew(t *testing.T) {
	cfg := &brain.Config{
		MemoryDir:      "/tmp/test-dxrk",
		VaultEnabled:   true,
		EmailEnabled:   false,
		CommandTimeout: 10 * time.Second,
	}

	b := brain.New(cfg)
	if b == nil {
		t.Fatal("Brain should not be nil")
	}
	if !b.IsRunning() {
		t.Error("Brain should be running")
	}
}

func TestBrainStatus(t *testing.T) {
	cfg := &brain.Config{}
	b := brain.New(cfg)

	status := b.Status()
	if status["running"] != true {
		t.Error("Status should show running")
	}
	if status["memory"] != false {
		t.Error("Memory should not be initialized yet")
	}
}

func TestBrainUptime(t *testing.T) {
	cfg := &brain.Config{}
	b := brain.New(cfg)

	time.Sleep(10 * time.Millisecond)
	uptime := b.Uptime()
	if uptime < 10*time.Millisecond {
		t.Error("Uptime should be at least 10ms")
	}
}

func TestBrainShutdown(t *testing.T) {
	cfg := &brain.Config{}
	b := brain.New(cfg)

	err := b.Shutdown()
	if err != nil {
		t.Errorf("Shutdown should not return error: %v", err)
	}
	if b.IsRunning() {
		t.Error("Brain should not be running after shutdown")
	}
}

func TestBrainString(t *testing.T) {
	cfg := &brain.Config{}
	b := brain.New(cfg)

	s := b.String()
	if s == "" {
		t.Error("String should not be empty")
	}
}
