package brain_test

import (
	"testing"

	"github.com/Dxrk777/Dxrk-Hex/internal/brain"
)

// NOTE: Email sending tests are integration tests that require
// a real SMTP server. We only test configuration and validation here.

func TestEmailerNew(t *testing.T) {
	cfg := brain.EmailConfig{
		Host:     "smtp.example.com",
		Port:     587,
		User:     "test@example.com",
		Password: "password123",
		From:     "test@example.com",
	}

	e := brain.NewEmailer(cfg)
	if e == nil {
		t.Fatal("Emailer should not be nil")
	}
}

func TestEmailerNewDefaultPort(t *testing.T) {
	cfg := brain.EmailConfig{
		Host: "smtp.example.com",
		User: "test@example.com",
	}

	e := brain.NewEmailer(cfg)
	if e == nil {
		t.Fatal("Emailer should not be nil")
	}
}

func TestEmailerIsConfigured(t *testing.T) {
	tests := []struct {
		name string
		cfg  brain.EmailConfig
		want bool
	}{
		{
			name: "empty config",
			cfg:  brain.EmailConfig{},
			want: false,
		},
		{
			name: "only host",
			cfg:  brain.EmailConfig{Host: "smtp.example.com"},
			want: false,
		},
		{
			name: "host and user",
			cfg:  brain.EmailConfig{Host: "smtp.example.com", User: "test@example.com"},
			want: true,
		},
		{
			name: "with tls",
			cfg:  brain.EmailConfig{Host: "smtp.gmail.com", User: "user@gmail.com", UseTLS: true},
			want: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := brain.NewEmailer(tt.cfg)
			if got := e.IsConfigured(); got != tt.want {
				t.Errorf("IsConfigured() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestEmailerSendNoRecipients(t *testing.T) {
	e := brain.NewEmailer(brain.EmailConfig{
		Host: "smtp.example.com",
		User: "test@example.com",
	})

	err := e.Send(brain.Email{
		To:      []string{},
		Subject: "Test",
		Body:    "Body",
	})

	if err == nil {
		t.Error("Send should fail with no recipients")
	}
}

func TestEmailStruct(t *testing.T) {
	email := brain.Email{
		To:      []string{"recipient@example.com"},
		Subject: "Test Subject",
		Body:    "Test Body",
		HTML:    true,
	}

	if len(email.To) != 1 {
		t.Error("Email should have 1 recipient")
	}
	if email.To[0] != "recipient@example.com" {
		t.Errorf("Recipient should be 'recipient@example.com', got %s", email.To[0])
	}
	if email.Subject != "Test Subject" {
		t.Errorf("Subject should be 'Test Subject', got %s", email.Subject)
	}
	if email.Body != "Test Body" {
		t.Errorf("Body should be 'Test Body', got %s", email.Body)
	}
	if !email.HTML {
		t.Error("Email should be HTML")
	}
}

func TestEmailerSendMultipleRecipients(t *testing.T) {
	e := brain.NewEmailer(brain.EmailConfig{
		Host: "smtp.example.com",
		User: "test@example.com",
	})

	err := e.Send(brain.Email{
		To:      []string{"user1@example.com", "user2@example.com"},
		Subject: "Test",
		Body:    "Body",
	})

	// Will fail due to no SMTP, but tests validation
	if err == nil {
		t.Error("Should fail due to no SMTP server")
	}
}

func TestEmailerSendEmptySubject(t *testing.T) {
	e := brain.NewEmailer(brain.EmailConfig{
		Host: "smtp.example.com",
		User: "test@example.com",
	})

	err := e.Send(brain.Email{
		To:      []string{"test@example.com"},
		Subject: "",
		Body:    "Body",
	})

	// Empty subject is allowed, should just fail on SMTP
	if err == nil {
		t.Error("Should fail due to no SMTP server")
	}
}

func TestEmailerSendEmptyBody(t *testing.T) {
	e := brain.NewEmailer(brain.EmailConfig{
		Host: "smtp.example.com",
		User: "test@example.com",
	})

	err := e.Send(brain.Email{
		To:      []string{"test@example.com"},
		Subject: "Test",
		Body:    "",
	})

	// Empty body is allowed, should just fail on SMTP
	if err == nil {
		t.Error("Should fail due to no SMTP server")
	}
}

func TestEmailConfig(t *testing.T) {
	cfg := brain.EmailConfig{
		Host:     "smtp.gmail.com",
		Port:     465,
		User:     "user@gmail.com",
		Password: "secret",
		From:     "user@gmail.com",
		UseTLS:   true,
	}

	if cfg.Host != "smtp.gmail.com" {
		t.Errorf("Host should be 'smtp.gmail.com', got %s", cfg.Host)
	}
	if cfg.Port != 465 {
		t.Errorf("Port should be 465, got %d", cfg.Port)
	}
	if !cfg.UseTLS {
		t.Error("UseTLS should be true")
	}
}

func TestEmailConfigDefault(t *testing.T) {
	cfg := brain.EmailConfig{}

	if cfg.Host != "" {
		t.Error("Default host should be empty")
	}
	if cfg.Port != 0 {
		t.Error("Default port should be 0")
	}
	if cfg.UseTLS {
		t.Error("Default UseTLS should be false")
	}
}

func TestEmailResult(t *testing.T) {
	result := &brain.EmailResult{
		Sent:    true,
		To:      []string{"test@example.com"},
		Subject: "Test",
	}

	if !result.Sent {
		t.Error("Email should be marked as sent")
	}
	if len(result.To) != 1 {
		t.Error("Should have 1 recipient")
	}
	if result.To[0] != "test@example.com" {
		t.Error("Recipient mismatch")
	}
}

func TestEmailResultWithError(t *testing.T) {
	result := &brain.EmailResult{
		Sent:    false,
		To:      []string{"test@example.com"},
		Subject: "Test",
		Error:   "connection refused",
	}

	if result.Sent {
		t.Error("Email should not be marked as sent")
	}
	if result.Error == "" {
		t.Error("Error should be set")
	}
}
