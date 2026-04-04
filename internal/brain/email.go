package brain

import (
	"crypto/tls"
	"fmt"
	"net"
	"net/smtp"
	"strings"
	"time"
)

// EmailConfig holds email configuration.
type EmailConfig struct {
	Host     string
	Port     int
	User     string
	Password string
	From     string
	UseTLS   bool
}

// Email represents an email message.
type Email struct {
	To      []string
	Subject string
	Body    string
	HTML    bool
}

// Emailer handles sending emails.
type Emailer struct {
	config EmailConfig
}

// NewEmailer creates a new Emailer instance.
func NewEmailer(cfg EmailConfig) *Emailer {
	if cfg.Port == 0 {
		cfg.Port = 587 // default SMTP port
	}
	return &Emailer{config: cfg}
}

// Send sends an email.
func (e *Emailer) Send(email Email) error {
	if len(email.To) == 0 {
		return fmt.Errorf("no recipients specified")
	}

	// Build email message
	var msg strings.Builder

	// Headers
	msg.WriteString(fmt.Sprintf("From: %s\r\n", e.config.From))
	msg.WriteString(fmt.Sprintf("To: %s\r\n", strings.Join(email.To, ", ")))
	msg.WriteString(fmt.Sprintf("Subject: %s\r\n", email.Subject))
	msg.WriteString(fmt.Sprintf("Date: %s\r\n", time.Now().Format(time.RFC1123Z)))

	if email.HTML {
		msg.WriteString("MIME-Version: 1.0\r\n")
		msg.WriteString("Content-Type: text/html; charset=\"UTF-8\"\r\n")
	} else {
		msg.WriteString("Content-Type: text/plain; charset=\"UTF-8\"\r\n")
	}

	msg.WriteString("\r\n")
	msg.WriteString(email.Body)

	// Send email
	return e.send(strings.Join(email.To, ", "), msg.String())
}

// SendHTML sends an HTML email.
func (e *Emailer) SendHTML(to []string, subject, html string) error {
	return e.Send(Email{
		To:      to,
		Subject: subject,
		Body:    html,
		HTML:    true,
	})
}

// send delivers the email via SMTP.
func (e *Emailer) send(to, msg string) error {
	addr := fmt.Sprintf("%s:%d", e.config.Host, e.config.Port)

	// Choose authentication method
	var auth smtp.Auth
	if e.config.User != "" && e.config.Password != "" {
		auth = smtp.PlainAuth("", e.config.User, e.config.Password, e.config.Host)
	}

	if e.config.UseTLS {
		return e.sendTLS(addr, auth, to, msg)
	}

	// Use dial timeout to avoid hanging
	return e.sendWithTimeout(addr, auth, to, msg)
}

// sendWithTimeout sends email with a 5-second connection timeout.
func (e *Emailer) sendWithTimeout(addr string, auth smtp.Auth, to, msg string) error {
	conn, err := net.DialTimeout("tcp", addr, 5*time.Second)
	if err != nil {
		return fmt.Errorf("connection timeout: %w", err)
	}
	defer conn.Close()

	client, err := smtp.NewClient(conn, e.config.Host)
	if err != nil {
		return fmt.Errorf("failed to create client: %w", err)
	}
	defer client.Close()

	if auth != nil {
		if err = client.Auth(auth); err != nil {
			return fmt.Errorf("authentication failed: %w", err)
		}
	}

	if err = client.Mail(e.config.From); err != nil {
		return fmt.Errorf("failed to set sender: %w", err)
	}
	if err = client.Rcpt(to); err != nil {
		return fmt.Errorf("failed to set recipient: %w", err)
	}

	w, err := client.Data()
	if err != nil {
		return fmt.Errorf("failed to start data: %w", err)
	}
	defer w.Close()

	_, err = w.Write([]byte(msg))
	if err != nil {
		return fmt.Errorf("failed to write message: %w", err)
	}

	return client.Quit()
}

// sendTLS sends email using TLS.
func (e *Emailer) sendTLS(addr string, auth smtp.Auth, to, msg string) error {
	tlsConfig := &tls.Config{
		ServerName: e.config.Host,
	}

	conn, err := tls.Dial("tcp", addr, tlsConfig)
	if err != nil {
		return fmt.Errorf("failed to connect: %w", err)
	}
	defer conn.Close()

	client, err := smtp.NewClient(conn, e.config.Host)
	if err != nil {
		return fmt.Errorf("failed to create client: %w", err)
	}
	defer client.Close()

	// Authenticate if needed
	if auth != nil {
		if err = client.Auth(auth); err != nil {
			return fmt.Errorf("authentication failed: %w", err)
		}
	}

	// Set sender and recipient
	if err = client.Mail(e.config.From); err != nil {
		return fmt.Errorf("failed to set sender: %w", err)
	}
	if err = client.Rcpt(to); err != nil {
		return fmt.Errorf("failed to set recipient: %w", err)
	}

	// Send message body
	writer, err := client.Data()
	if err != nil {
		return fmt.Errorf("failed to start data: %w", err)
	}
	defer writer.Close()

	_, err = writer.Write([]byte(msg))
	if err != nil {
		return fmt.Errorf("failed to write message: %w", err)
	}

	return nil
}

// TestConnection tests the email configuration.
func (e *Emailer) TestConnection() error {
	testEmail := Email{
		To:      []string{e.config.User}, // send to self
		Subject: "Dxrk Hell - Test Email",
		Body:    "This is a test email from Dxrk Hell.",
	}
	return e.Send(testEmail)
}

// IsConfigured returns true if email is properly configured.
func (e *Emailer) IsConfigured() bool {
	return e.config.Host != "" && e.config.User != ""
}

// String returns a string representation of the Emailer.
func (e *Emailer) String() string {
	return fmt.Sprintf("Emailer{host: %s, port: %d, user: %s}", e.config.Host, e.config.Port, e.config.User)
}
