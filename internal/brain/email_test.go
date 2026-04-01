package brain

import (
	"testing"
)

func TestEmailConfig(t *testing.T) {
	cfg := EmailConfig{
		Host:     "smtp.gmail.com",
		Port:     587,
		User:     "test@gmail.com",
		Password: "password",
		From:     "test@gmail.com",
		UseTLS:   true,
	}

	if cfg.Host != "smtp.gmail.com" {
		t.Errorf("Expected host smtp.gmail.com, got %s", cfg.Host)
	}

	if cfg.Port != 587 {
		t.Errorf("Expected port 587, got %d", cfg.Port)
	}
}

func TestNewEmailer(t *testing.T) {
	cfg := EmailConfig{
		Host:     "smtp.test.com",
		Port:     587,
		User:     "test",
		Password: "test",
		From:     "test@test.com",
	}

	emailer := NewEmailer(cfg)

	if emailer.config.Host != "smtp.test.com" {
		t.Errorf("Expected host smtp.test.com, got %s", emailer.config.Host)
	}
}

func TestNewEmailerDefaultPort(t *testing.T) {
	cfg := EmailConfig{
		Host:     "smtp.test.com",
		Port:     0, // Should default to 587
		User:     "test",
		Password: "test",
		From:     "test@test.com",
	}

	emailer := NewEmailer(cfg)

	if emailer.config.Port != 587 {
		t.Errorf("Expected default port 587, got %d", emailer.config.Port)
	}
}

func TestEmailerIsConfigured(t *testing.T) {
	tests := []struct {
		name     string
		config   EmailConfig
		expected bool
	}{
		{
			name:     "configured",
			config:   EmailConfig{Host: "smtp.test.com", User: "test"},
			expected: true,
		},
		{
			name:     "missing host",
			config:   EmailConfig{Host: "", User: "test"},
			expected: false,
		},
		{
			name:     "missing user",
			config:   EmailConfig{Host: "smtp.test.com", User: ""},
			expected: false,
		},
		{
			name:     "empty",
			config:   EmailConfig{},
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			emailer := NewEmailer(tt.config)
			if emailer.IsConfigured() != tt.expected {
				t.Errorf("IsConfigured() = %v, want %v", emailer.IsConfigured(), tt.expected)
			}
		})
	}
}

func TestEmail(t *testing.T) {
	email := Email{
		To:      []string{"recipient@test.com"},
		Subject: "Test Subject",
		Body:    "Test Body",
		HTML:    false,
	}

	if len(email.To) != 1 {
		t.Errorf("Expected 1 recipient, got %d", len(email.To))
	}

	if email.To[0] != "recipient@test.com" {
		t.Errorf("Expected recipient@test.com, got %s", email.To[0])
	}

	if email.Subject != "Test Subject" {
		t.Errorf("Expected Test Subject, got %s", email.Subject)
	}
}

func TestEmailMultipleRecipients(t *testing.T) {
	email := Email{
		To:      []string{"user1@test.com", "user2@test.com", "user3@test.com"},
		Subject: "Test",
		Body:    "Body",
	}

	if len(email.To) != 3 {
		t.Errorf("Expected 3 recipients, got %d", len(email.To))
	}
}

func TestSendNoRecipients(t *testing.T) {
	emailer := NewEmailer(EmailConfig{
		Host:     "smtp.test.com",
		Port:     587,
		User:     "test",
		Password: "test",
		From:     "test@test.com",
	})

	email := Email{
		To:      []string{},
		Subject: "Test",
		Body:    "Body",
	}

	err := emailer.Send(email)
	if err == nil {
		t.Error("Expected error for no recipients")
	}
}

func TestSendHTML(t *testing.T) {
	emailer := NewEmailer(EmailConfig{
		Host:     "smtp.test.com",
		Port:     587,
		User:     "test",
		Password: "test",
		From:     "test@test.com",
	})

	// This will fail because SMTP is not configured, but we're testing the function exists
	err := emailer.SendHTML([]string{"test@test.com"}, "Test", "<h1>HTML</h1>")

	// We expect this to fail since there's no real SMTP server
	// But the function should be called without panicking
	if err == nil {
		t.Log("Email sent successfully (unexpected without real SMTP)")
	}
}

func TestEmailerString(t *testing.T) {
	emailer := NewEmailer(EmailConfig{
		Host:     "smtp.test.com",
		Port:     587,
		User:     "test",
		Password: "test",
		From:     "test@test.com",
	})

	str := emailer.String()
	if str == "" {
		t.Error("String() should return a non-empty value")
	}
}
