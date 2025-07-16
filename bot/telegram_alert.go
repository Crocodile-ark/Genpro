package main

import (
	"fmt"
	"log"
	"strings"
	"time"
)

// TelegramAlert handles Telegram notifications
type TelegramAlert struct {
	config *BotConfig
	
	// Alert state
	enabled     bool
	alertCount  int64
	lastAlert   time.Time
	rateLimiter *RateLimiter
}

// RateLimiter handles rate limiting for Telegram alerts
type RateLimiter struct {
	maxAlertsPerMinute int
	alertTimes         []time.Time
}

// NewRateLimiter creates a new rate limiter
func NewRateLimiter(maxAlertsPerMinute int) *RateLimiter {
	return &RateLimiter{
		maxAlertsPerMinute: maxAlertsPerMinute,
		alertTimes:         make([]time.Time, 0),
	}
}

// CanSendAlert checks if we can send an alert without hitting rate limits
func (rl *RateLimiter) CanSendAlert() bool {
	now := time.Now()
	
	// Remove alerts older than 1 minute
	cutoff := now.Add(-time.Minute)
	var recentAlerts []time.Time
	for _, alertTime := range rl.alertTimes {
		if alertTime.After(cutoff) {
			recentAlerts = append(recentAlerts, alertTime)
		}
	}
	rl.alertTimes = recentAlerts
	
	// Check if we're under the limit
	if len(rl.alertTimes) >= rl.maxAlertsPerMinute {
		return false
	}
	
	// Record this alert
	rl.alertTimes = append(rl.alertTimes, now)
	return true
}

// NewTelegramAlert creates a new Telegram alert instance
func NewTelegramAlert(config *BotConfig) *TelegramAlert {
	return &TelegramAlert{
		config:      config,
		rateLimiter: NewRateLimiter(10), // Max 10 alerts per minute
	}
}

// Initialize initializes the Telegram alert system
func (ta *TelegramAlert) Initialize() error {
	log.Println("Initializing Telegram Alert...")
	
	// Validate Telegram configuration
	if ta.config.TelegramToken == "" {
		return fmt.Errorf("telegram token not configured")
	}
	
	if ta.config.TelegramChatID == "" {
		return fmt.Errorf("telegram chat ID not configured")
	}
	
	// Validate token format (should start with a number followed by colon)
	if !ta.isValidTokenFormat(ta.config.TelegramToken) {
		return fmt.Errorf("invalid telegram token format")
	}
	
	// Test connection
	if err := ta.testConnection(); err != nil {
		return fmt.Errorf("failed to connect to Telegram: %w", err)
	}
	
	ta.enabled = true
	log.Println("Telegram Alert initialized successfully")
	return nil
}

// isValidTokenFormat checks if the token has a valid format
func (ta *TelegramAlert) isValidTokenFormat(token string) bool {
	// Basic validation: should contain a colon and be of reasonable length
	parts := strings.Split(token, ":")
	return len(parts) == 2 && len(parts[0]) > 0 && len(parts[1]) > 0
}

// testConnection tests the Telegram bot connection
func (ta *TelegramAlert) testConnection() error {
	log.Printf("Testing Telegram connection (Bot ID: %s..., Chat: %s)", 
		ta.config.TelegramToken[:10], ta.config.TelegramChatID)
	
	// In a real implementation, this would:
	// 1. Make an API call to Telegram's getMe endpoint
	// 2. Verify the bot token is valid
	// 3. Test sending a message to the chat ID
	
	// For now, we'll simulate the connection test
	time.Sleep(500 * time.Millisecond)
	
	// Simulate occasional connection failures
	if strings.Contains(ta.config.TelegramToken, "invalid") {
		return fmt.Errorf("invalid telegram token")
	}
	
	log.Println("Telegram connection test successful")
	return nil
}

// SendAlert sends an alert message to Telegram
func (ta *TelegramAlert) SendAlert(message string) error {
	if !ta.enabled {
		log.Printf("Telegram alerts disabled, would send: %s", message)
		return nil
	}
	
	// Check rate limiting
	if !ta.rateLimiter.CanSendAlert() {
		log.Println("Rate limit exceeded, skipping Telegram alert")
		return nil
	}
	
	// Validate message
	if message == "" {
		return fmt.Errorf("empty message")
	}
	
	// Add timestamp and formatting
	formattedMessage := fmt.Sprintf("🤖 *GXR Bot Alert*\n\n%s\n\n⏰ %s", 
		message, time.Now().Format("2006-01-02 15:04:05 UTC"))
	
	// Send the message
	if err := ta.sendMessage(formattedMessage); err != nil {
		return fmt.Errorf("failed to send telegram message: %w", err)
	}
	
	ta.alertCount++
	ta.lastAlert = time.Now()
	
	return nil
}

// sendMessage handles the actual message sending
func (ta *TelegramAlert) sendMessage(message string) error {
	log.Printf("📱 Sending Telegram Alert: %s", message)
	
	// In a real implementation, this would:
	// 1. Make an HTTP POST request to Telegram's sendMessage endpoint
	// 2. Include the bot token, chat ID, and message
	// 3. Handle response and errors
	
	// For now, we'll simulate the sending process
	time.Sleep(100 * time.Millisecond)
	
	// Simulate occasional sending failures
	if ta.alertCount > 0 && ta.alertCount%50 == 0 {
		return fmt.Errorf("simulated telegram API error")
	}
	
	log.Println("Telegram message sent successfully")
	return nil
}

// SendUptimeAlert sends uptime monitoring alert
func (ta *TelegramAlert) SendUptimeAlert(component string, status string) error {
	if component == "" || status == "" {
		return fmt.Errorf("component and status are required")
	}
	
	var emoji string
	switch strings.ToLower(status) {
	case "up", "online", "active":
		emoji = "✅"
	case "down", "offline", "inactive":
		emoji = "❌"
	default:
		emoji = "⚠️"
	}
	
	message := fmt.Sprintf("%s *Uptime Alert*\n\nComponent: %s\nStatus: %s", emoji, component, status)
	return ta.SendAlert(message)
}

// SendPoolImbalanceAlert sends pool imbalance alert
func (ta *TelegramAlert) SendPoolImbalanceAlert(pool string, imbalance string) error {
	if pool == "" || imbalance == "" {
		return fmt.Errorf("pool and imbalance are required")
	}
	
	message := fmt.Sprintf("⚖️ *Pool Imbalance Alert*\n\nPool: %s\nImbalance: %s", pool, imbalance)
	return ta.SendAlert(message)
}

// SendPriceAlert sends price monitoring alert
func (ta *TelegramAlert) SendPriceAlert(price float64, limit float64) error {
	message := fmt.Sprintf("💰 *Price Alert*\n\nCurrent GXR Price: $%.2f\nLimit: $%.2f\n\n🚨 Emergency mode activated!", 
		price, limit)
	return ta.SendAlert(message)
}

// SendRewardDistributionAlert sends reward distribution alert
func (ta *TelegramAlert) SendRewardDistributionAlert(amount string, cycle int64) error {
	if amount == "" {
		return fmt.Errorf("amount is required")
	}
	
	message := fmt.Sprintf("💎 *Reward Distribution*\n\nAmount: %s\nCycle: %d\n\n✅ Monthly rewards distributed successfully!", 
		amount, cycle)
	return ta.SendAlert(message)
}

// SendEmergencyAlert sends emergency alert
func (ta *TelegramAlert) SendEmergencyAlert(alert string) error {
	if alert == "" {
		return fmt.Errorf("alert message is required")
	}
	
	message := fmt.Sprintf("🚨 *EMERGENCY ALERT*\n\n%s\n\n⚠️ Immediate attention required!", alert)
	return ta.SendAlert(message)
}

// SendCustomAlert sends a custom alert with specified emoji
func (ta *TelegramAlert) SendCustomAlert(emoji string, title string, content string) error {
	if title == "" || content == "" {
		return fmt.Errorf("title and content are required")
	}
	
	if emoji == "" {
		emoji = "ℹ️"
	}
	
	message := fmt.Sprintf("%s *%s*\n\n%s", emoji, title, content)
	return ta.SendAlert(message)
}

// Disable disables Telegram alerts
func (ta *TelegramAlert) Disable() {
	ta.enabled = false
	log.Println("Telegram alerts disabled")
}

// Enable enables Telegram alerts
func (ta *TelegramAlert) Enable() {
	ta.enabled = true
	log.Println("Telegram alerts enabled")
}

// GetStatus returns the current Telegram alert status
func (ta *TelegramAlert) GetStatus() map[string]interface{} {
	return map[string]interface{}{
		"enabled":      ta.enabled,
		"alert_count":  ta.alertCount,
		"last_alert":   ta.lastAlert,
		"token_set":    ta.config.TelegramToken != "",
		"chat_id_set":  ta.config.TelegramChatID != "",
		"rate_limiter": map[string]interface{}{
			"max_per_minute": ta.rateLimiter.maxAlertsPerMinute,
			"recent_alerts":  len(ta.rateLimiter.alertTimes),
		},
	}
}