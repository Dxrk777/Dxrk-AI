package screens

import (
	"fmt"
	"strings"

	"github.com/Dxrk777/Dxrk-Hex/internal/connector"
	"github.com/Dxrk777/Dxrk-Hex/internal/tui/styles"
)

// RemoteOptions returns the remote/connector menu options.
func RemoteOptions() []string {
	return []string{
		"Configure Telegram",
		"Configure Discord",
		"Configure WhatsApp",
		"Start Server",
		"Stop Server",
		"Status",
		"Back",
	}
}

// RemoteState holds the state for remote connector configuration.
type RemoteState struct {
	TelegramToken  string
	TelegramChatID string
	DiscordWebhook string
	WhatsAppSID    string
	WhatsAppToken  string
	WhatsAppFrom   string
	ServerRunning  bool
	Config         *connector.Config
}

// RenderRemote renders the remote connector configuration screen.
func RenderRemote(cursor int, state RemoteState) string {
	var b strings.Builder

	b.WriteString(styles.RenderLogo())
	b.WriteString("\n\n")
	b.WriteString(styles.HeadingStyle.Render("🔗 Remote Control"))
	b.WriteString("\n\n")

	b.WriteString(styles.SubtextStyle.Render("Connect Dxrk Hell to messaging platforms:"))
	b.WriteString("\n\n")

	// Telegram status
	if state.TelegramToken != "" {
		b.WriteString("✅ Telegram: Configured\n")
	} else {
		b.WriteString("❌ Telegram: Not configured\n")
	}

	// Discord status
	if state.DiscordWebhook != "" {
		b.WriteString("✅ Discord: Configured\n")
	} else {
		b.WriteString("❌ Discord: Not configured\n")
	}

	// WhatsApp status
	if state.WhatsAppSID != "" && state.WhatsAppToken != "" {
		b.WriteString("✅ WhatsApp: Configured\n")
	} else {
		b.WriteString("❌ WhatsApp: Not configured\n")
	}

	// Server status
	b.WriteString("\n")
	if state.ServerRunning {
		b.WriteString("🟢 Server: Running (port 8081)\n")
	} else {
		b.WriteString("⚫ Server: Stopped\n")
	}

	b.WriteString("\n")
	b.WriteString(renderOptions(RemoteOptions(), cursor))
	b.WriteString("\n")
	b.WriteString(styles.HelpStyle.Render("j/k: navigate • enter: select • esc: back"))

	return styles.FrameStyle.Render(b.String())
}

// RenderRemoteHelp shows configuration instructions.
func RenderRemoteHelp(platform string) string {
	var b strings.Builder

	b.WriteString(styles.RenderLogo())
	b.WriteString("\n\n")
	b.WriteString(styles.HeadingStyle.Render(fmt.Sprintf("📱 Configure %s", platform)))
	b.WriteString("\n\n")

	switch platform {
	case "telegram":
		b.WriteString(styles.SubtextStyle.Render("Telegram Setup:"))
		b.WriteString("\n\n")
		b.WriteString("1. Open Telegram and message @BotFather\n")
		b.WriteString("2. Send /newbot to create a new bot\n")
		b.WriteString("3. Copy the Bot Token\n")
		b.WriteString("4. Start a chat with your bot and send /start\n")
		b.WriteString("5. Visit https://api.telegram.org/bot<TOKEN>/getUpdates\n")
		b.WriteString("6. Copy your Chat ID\n")

	case "discord":
		b.WriteString(styles.SubtextStyle.Render("Discord Setup:"))
		b.WriteString("\n\n")
		b.WriteString("1. Open Discord Server Settings\n")
		b.WriteString("2. Go to Integrations > Webhooks\n")
		b.WriteString("3. Create a new webhook or copy existing\n")
		b.WriteString("4. Paste the Webhook URL\n")

	case "whatsapp":
		b.WriteString(styles.SubtextStyle.Render("WhatsApp Setup (Twilio):"))
		b.WriteString("\n\n")
		b.WriteString("1. Create account at twilio.com\n")
		b.WriteString("2. Get Account SID and Auth Token\n")
		b.WriteString("3. Purchase a WhatsApp-enabled number\n")
		b.WriteString("4. Configure Twilio sandbox or verify number\n")
		b.WriteString("5. Enter credentials here\n")
	}

	b.WriteString("\n")
	b.WriteString(styles.HelpStyle.Render("esc: back"))

	return styles.FrameStyle.Render(b.String())
}

// RenderRemoteStatus shows remote server status.
func RenderRemoteStatus(status map[string]interface{}) string {
	var b strings.Builder

	b.WriteString(styles.RenderLogo())
	b.WriteString("\n\n")
	b.WriteString(styles.HeadingStyle.Render("📊 Remote Server Status"))
	b.WriteString("\n\n")

	b.WriteString(fmt.Sprintf("Enabled: %v\n", status["enabled"]))
	b.WriteString(fmt.Sprintf("Running: %v\n", status["running"]))
	b.WriteString(fmt.Sprintf("Port: %v\n", status["port"]))
	b.WriteString(fmt.Sprintf("Uptime: %v\n", status["uptime"]))
	b.WriteString("\n")
	b.WriteString("Platforms:\n")
	b.WriteString(fmt.Sprintf("  Telegram: %v\n", status["telegram"]))
	b.WriteString(fmt.Sprintf("  Discord: %v\n", status["discord"]))
	b.WriteString(fmt.Sprintf("  WhatsApp: %v\n", status["whatsapp"]))

	b.WriteString("\n")
	b.WriteString(styles.HelpStyle.Render("esc: back"))

	return styles.FrameStyle.Render(b.String())
}
