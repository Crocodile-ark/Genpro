package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"
)

const (
	// TelegramAPIBaseURL is the base URL for Telegram Bot API
	TelegramAPIBaseURL = "https://api.telegram.org/bot"
	// MaxAlertsPerMinute is the rate limit for alerts
	MaxAlertsPerMinute = 10
	// AlertQueueSize is the maximum number of queued alerts
	AlertQueueSize = 100
	// RetryAttempts is the number of retry attempts for failed alerts
	RetryAttempts = 3
	// RetryDelay is the delay between retry attempts
	RetryDelay = 5 * time.Second
	// MessageSizeLimit is the maximum message size for Telegram
	MessageSizeLimit = 4096
	// AlertPriorityHigh is for high priority alerts
	AlertPriorityHigh = 1
	// AlertPriorityMedium is for medium priority alerts
	AlertPriorityMedium = 2
	// AlertPriorityLow is for low priority alerts
	AlertPriorityLow = 3
)

// AlertType represents different types of alerts
type AlertType int

const (
	AlertTypeInfo AlertType = iota
	AlertTypeWarning
	AlertTypeError
	AlertTypeCritical
	AlertTypeSuccess
)

func (at AlertType) String() string {
	switch at {
	case AlertTypeInfo:
		return "INFO"
	case AlertTypeWarning:
		return "WARNING"
	case AlertTypeError:
		return "ERROR"
	case AlertTypeCritical:
		return "CRITICAL"
	case AlertTypeSuccess:
		return "SUCCESS"
	default:
		return "UNKNOWN"
	}
}

func (at AlertType) Emoji() string {
	switch at {
	case AlertTypeInfo:
		return "ℹ️"
	case AlertTypeWarning:
		return "⚠️"
	case AlertTypeError:
		return "❌"
	case AlertTypeCritical:
		return "🚨"
	case AlertTypeSuccess:
		return "✅"
	default:
		return "📱"
	}
}

// TelegramAlert represents a structured alert message
type TelegramAlert struct {
	config    *BotConfig
	client    *http.Client
	mu        sync.RWMutex
	
	// Rate limiting
	alertTimes       []time.Time
	alertQueue       chan *Alert
	rateLimitEnabled bool
	
	// Statistics
	totalAlerts        int64
	successfulAlerts   int64
	failedAlerts       int64
	rateLimitedAlerts  int64
	lastAlertTime      time.Time
	
	// Alert categorization
	alertCounts  map[AlertType]int64
	alertHistory []AlertRecord
	
	// Configuration
	botToken    string
	chatID      string
	apiURL      string
	maxRetries  int
	retryDelay  time.Duration
	
	// Control
	running    bool
	stopChan   chan struct{}
}

// Alert represents an individual alert
type Alert struct {
	ID          string
	Type        AlertType
	Priority    int
	Title       string
	Message     string
	Timestamp   time.Time
	Metadata    map[string]interface{}
	Retries     int
	LastAttempt time.Time
}

// AlertRecord represents a historical alert record
type AlertRecord struct {
	ID        string
	Type      AlertType
	Title     string
	Message   string
	Timestamp time.Time
	Success   bool
	Attempts  int
}

// TelegramMessage represents a Telegram API message
type TelegramMessage struct {
	ChatID    string `json:"chat_id"`
	Text      string `json:"text"`
	ParseMode string `json:"parse_mode,omitempty"`
}

// TelegramResponse represents a Telegram API response
type TelegramResponse struct {
	OK          bool   `json:"ok"`
	Result      interface{} `json:"result,omitempty"`
	ErrorCode   int    `json:"error_code,omitempty"`
	Description string `json:"description,omitempty"`
}

// NewTelegramAlert creates a new enhanced Telegram alert system
func NewTelegramAlert(config *BotConfig) *TelegramAlert {
	ta := &TelegramAlert{
		config:           config,
		client:           &http.Client{Timeout: 30 * time.Second},
		alertTimes:       make([]time.Time, 0),
		alertQueue:       make(chan *Alert, AlertQueueSize),
		rateLimitEnabled: true,
		alertCounts:      make(map[AlertType]int64),
		alertHistory:     make([]AlertRecord, 0),
		maxRetries:       RetryAttempts,
		retryDelay:       RetryDelay,
		stopChan:         make(chan struct{}),
	}
	
	// Validate and set configuration
	if err := ta.validateConfig(); err != nil {
		log.Printf("Telegram alert configuration error: %v", err)
		return ta
	}
	
	// Start alert processing
	go ta.processAlerts()
	
	return ta
}

// validateConfig validates the Telegram configuration
func (ta *TelegramAlert) validateConfig() error {
	if ta.config.TelegramToken == "" {
		return fmt.Errorf("telegram_token is required")
	}
	
	if ta.config.TelegramChatID == "" {
		return fmt.Errorf("telegram_chat_id is required")
	}
	
	ta.botToken = ta.config.TelegramToken
	ta.chatID = ta.config.TelegramChatID
	ta.apiURL = fmt.Sprintf("%s%s", TelegramAPIBaseURL, ta.botToken)
	
	// Validate bot token format
	if !strings.Contains(ta.botToken, ":") {
		return fmt.Errorf("invalid bot token format")
	}
	
	// Validate chat ID format
	if _, err := strconv.ParseInt(ta.chatID, 10, 64); err != nil {
		if !strings.HasPrefix(ta.chatID, "@") {
			return fmt.Errorf("invalid chat ID format")
		}
	}
	
	ta.running = true
	log.Printf("Telegram alert system initialized - Chat: %s", ta.chatID)
	
	return nil
}

// processAlerts processes the alert queue
func (ta *TelegramAlert) processAlerts() {
	for {
		select {
		case alert := <-ta.alertQueue:
			ta.handleAlert(alert)
		case <-ta.stopChan:
			log.Printf("Stopping Telegram alert processor")
			return
		}
	}
}

// handleAlert handles an individual alert
func (ta *TelegramAlert) handleAlert(alert *Alert) {
	ta.mu.Lock()
	defer ta.mu.Unlock()
	
	// Check rate limiting
	if ta.rateLimitEnabled && !ta.canSendAlert() {
		ta.rateLimitedAlerts++
		log.Printf("Alert rate limited: %s", alert.Title)
		return
	}
	
	// Format message
	message := ta.formatAlert(alert)
	
	// Send with retries
	success := ta.sendWithRetries(message, alert)
	
	// Update statistics
	ta.totalAlerts++
	ta.lastAlertTime = time.Now()
	ta.alertCounts[alert.Type]++
	
	if success {
		ta.successfulAlerts++
	} else {
		ta.failedAlerts++
	}
	
	// Add to history
	ta.addToHistory(alert, success)
	
	// Update rate limiting
	if ta.rateLimitEnabled {
		ta.alertTimes = append(ta.alertTimes, time.Now())
		ta.cleanupOldAlerts()
	}
}

// canSendAlert checks if we can send an alert based on rate limiting
func (ta *TelegramAlert) canSendAlert() bool {
	ta.cleanupOldAlerts()
	return len(ta.alertTimes) < MaxAlertsPerMinute
}

// cleanupOldAlerts removes old alert timestamps for rate limiting
func (ta *TelegramAlert) cleanupOldAlerts() {
	cutoff := time.Now().Add(-1 * time.Minute)
	newTimes := make([]time.Time, 0)
	
	for _, alertTime := range ta.alertTimes {
		if alertTime.After(cutoff) {
			newTimes = append(newTimes, alertTime)
		}
	}
	
	ta.alertTimes = newTimes
}

// formatAlert formats an alert message for Telegram
func (ta *TelegramAlert) formatAlert(alert *Alert) string {
	timestamp := alert.Timestamp.Format("2006-01-02 15:04:05")
	
	var parts []string
	
	// Add header with emoji and type
	header := fmt.Sprintf("%s *%s*", alert.Type.Emoji(), alert.Type.String())
	parts = append(parts, header)
	
	// Add title
	if alert.Title != "" {
		parts = append(parts, fmt.Sprintf("*%s*", alert.Title))
	}
	
	// Add message
	if alert.Message != "" {
		parts = append(parts, alert.Message)
	}
	
	// Add timestamp
	parts = append(parts, fmt.Sprintf("📅 %s", timestamp))
	
	// Add metadata if present
	if len(alert.Metadata) > 0 {
		parts = append(parts, "")
		parts = append(parts, "*Details:*")
		for key, value := range alert.Metadata {
			parts = append(parts, fmt.Sprintf("• %s: %v", key, value))
		}
	}
	
	message := strings.Join(parts, "\n")
	
	// Truncate if too long
	if len(message) > MessageSizeLimit {
		message = message[:MessageSizeLimit-3] + "..."
	}
	
	return message
}

// sendWithRetries sends a message with retry logic
func (ta *TelegramAlert) sendWithRetries(message string, alert *Alert) bool {
	for attempt := 0; attempt < ta.maxRetries; attempt++ {
		if attempt > 0 {
			time.Sleep(ta.retryDelay)
		}
		
		if ta.sendMessage(message) {
			return true
		}
		
		alert.Retries++
		alert.LastAttempt = time.Now()
		
		log.Printf("Alert retry %d/%d failed: %s", attempt+1, ta.maxRetries, alert.Title)
	}
	
	return false
}

// sendMessage sends a message to Telegram
func (ta *TelegramAlert) sendMessage(message string) bool {
	if !ta.running {
		return false
	}
	
	telegramMsg := TelegramMessage{
		ChatID:    ta.chatID,
		Text:      message,
		ParseMode: "Markdown",
	}
	
	jsonData, err := json.Marshal(telegramMsg)
	if err != nil {
		log.Printf("Failed to marshal Telegram message: %v", err)
		return false
	}
	
	url := fmt.Sprintf("%s/sendMessage", ta.apiURL)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		log.Printf("Failed to create Telegram request: %v", err)
		return false
	}
	
	req.Header.Set("Content-Type", "application/json")
	
	resp, err := ta.client.Do(req)
	if err != nil {
		log.Printf("Failed to send Telegram message: %v", err)
		return false
	}
	defer resp.Body.Close()
	
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("Failed to read Telegram response: %v", err)
		return false
	}
	
	var telegramResp TelegramResponse
	if err := json.Unmarshal(body, &telegramResp); err != nil {
		log.Printf("Failed to parse Telegram response: %v", err)
		return false
	}
	
	if !telegramResp.OK {
		log.Printf("Telegram API error: %d - %s", telegramResp.ErrorCode, telegramResp.Description)
		return false
	}
	
	return true
}

// addToHistory adds an alert to the history
func (ta *TelegramAlert) addToHistory(alert *Alert, success bool) {
	record := AlertRecord{
		ID:        alert.ID,
		Type:      alert.Type,
		Title:     alert.Title,
		Message:   alert.Message,
		Timestamp: alert.Timestamp,
		Success:   success,
		Attempts:  alert.Retries + 1,
	}
	
	ta.alertHistory = append(ta.alertHistory, record)
	
	// Keep only last 100 records
	if len(ta.alertHistory) > 100 {
		ta.alertHistory = ta.alertHistory[1:]
	}
}

// SendAlert sends a basic alert (backward compatibility)
func (ta *TelegramAlert) SendAlert(message string) error {
	return ta.SendAlertWithType(AlertTypeInfo, "Alert", message)
}

// SendAlertWithType sends an alert with a specific type
func (ta *TelegramAlert) SendAlertWithType(alertType AlertType, title, message string) error {
	alert := &Alert{
		ID:        fmt.Sprintf("alert-%d", time.Now().UnixNano()),
		Type:      alertType,
		Priority:  AlertPriorityMedium,
		Title:     title,
		Message:   message,
		Timestamp: time.Now(),
		Metadata:  make(map[string]interface{}),
	}
	
	return ta.QueueAlert(alert)
}

// SendRebalancerAlert sends a rebalancer state change alert
func (ta *TelegramAlert) SendRebalancerAlert(state, reason string, price float64) error {
	alert := &Alert{
		ID:        fmt.Sprintf("rebalancer-%d", time.Now().UnixNano()),
		Type:      AlertTypeWarning,
		Priority:  AlertPriorityHigh,
		Title:     "Rebalancer State Change",
		Message:   reason,
		Timestamp: time.Now(),
		Metadata: map[string]interface{}{
			"state": state,
			"price": fmt.Sprintf("$%.2f", price),
		},
	}
	
	return ta.QueueAlert(alert)
}

// SendValidatorAlert sends a validator-related alert
func (ta *TelegramAlert) SendValidatorAlert(validatorName, reason string, inactiveDays int) error {
	alertType := AlertTypeWarning
	if inactiveDays > 10 {
		alertType = AlertTypeCritical
	}
	
	alert := &Alert{
		ID:        fmt.Sprintf("validator-%d", time.Now().UnixNano()),
		Type:      alertType,
		Priority:  AlertPriorityHigh,
		Title:     "Validator Inactivity",
		Message:   reason,
		Timestamp: time.Now(),
		Metadata: map[string]interface{}{
			"validator":      validatorName,
			"inactive_days":  inactiveDays,
			"threshold":      10,
		},
	}
	
	return ta.QueueAlert(alert)
}

// SendBotAlert sends a bot-related alert
func (ta *TelegramAlert) SendBotAlert(botType, status, reason string) error {
	alertType := AlertTypeWarning
	if status == "error" || status == "stopped" {
		alertType = AlertTypeError
	}
	
	alert := &Alert{
		ID:        fmt.Sprintf("bot-%d", time.Now().UnixNano()),
		Type:      alertType,
		Priority:  AlertPriorityMedium,
		Title:     fmt.Sprintf("Bot Status: %s", botType),
		Message:   reason,
		Timestamp: time.Now(),
		Metadata: map[string]interface{}{
			"bot_type": botType,
			"status":   status,
		},
	}
	
	return ta.QueueAlert(alert)
}

// SendHalvingAlert sends a halving-related alert
func (ta *TelegramAlert) SendHalvingAlert(cycle uint64, event, details string) error {
	alert := &Alert{
		ID:        fmt.Sprintf("halving-%d", time.Now().UnixNano()),
		Type:      AlertTypeInfo,
		Priority:  AlertPriorityMedium,
		Title:     "Halving Event",
		Message:   fmt.Sprintf("Cycle %d: %s", cycle, event),
		Timestamp: time.Now(),
		Metadata: map[string]interface{}{
			"cycle":   cycle,
			"event":   event,
			"details": details,
		},
	}
	
	return ta.QueueAlert(alert)
}

// SendEmergencyAlert sends a high-priority emergency alert
func (ta *TelegramAlert) SendEmergencyAlert(title, message string, metadata map[string]interface{}) error {
	alert := &Alert{
		ID:        fmt.Sprintf("emergency-%d", time.Now().UnixNano()),
		Type:      AlertTypeCritical,
		Priority:  AlertPriorityHigh,
		Title:     title,
		Message:   message,
		Timestamp: time.Now(),
		Metadata:  metadata,
	}
	
	// Emergency alerts bypass rate limiting
	oldRateLimit := ta.rateLimitEnabled
	ta.rateLimitEnabled = false
	defer func() { ta.rateLimitEnabled = oldRateLimit }()
	
	return ta.QueueAlert(alert)
}

// QueueAlert adds an alert to the processing queue
func (ta *TelegramAlert) QueueAlert(alert *Alert) error {
	if !ta.running {
		return fmt.Errorf("telegram alert system is not running")
	}
	
	select {
	case ta.alertQueue <- alert:
		return nil
	case <-time.After(5 * time.Second):
		return fmt.Errorf("alert queue is full")
	}
}

// EnableRateLimit enables or disables rate limiting
func (ta *TelegramAlert) EnableRateLimit(enabled bool) {
	ta.mu.Lock()
	defer ta.mu.Unlock()
	
	ta.rateLimitEnabled = enabled
	log.Printf("Telegram rate limiting %s", map[bool]string{true: "enabled", false: "disabled"}[enabled])
}

// GetStatistics returns alert statistics
func (ta *TelegramAlert) GetStatistics() map[string]interface{} {
	ta.mu.RLock()
	defer ta.mu.RUnlock()
	
	stats := map[string]interface{}{
		"total_alerts":         ta.totalAlerts,
		"successful_alerts":    ta.successfulAlerts,
		"failed_alerts":        ta.failedAlerts,
		"rate_limited_alerts":  ta.rateLimitedAlerts,
		"last_alert_time":      ta.lastAlertTime.Format(time.RFC3339),
		"queue_size":           len(ta.alertQueue),
		"rate_limit_enabled":   ta.rateLimitEnabled,
		"current_rate_count":   len(ta.alertTimes),
		"max_rate_per_minute":  MaxAlertsPerMinute,
		"alert_history_size":   len(ta.alertHistory),
		"running":              ta.running,
	}
	
	// Add alert counts by type
	typeCounts := make(map[string]int64)
	for alertType, count := range ta.alertCounts {
		typeCounts[alertType.String()] = count
	}
	stats["alert_counts_by_type"] = typeCounts
	
	return stats
}

// GetHistory returns recent alert history
func (ta *TelegramAlert) GetHistory() []AlertRecord {
	ta.mu.RLock()
	defer ta.mu.RUnlock()
	
	// Return a copy to avoid race conditions
	history := make([]AlertRecord, len(ta.alertHistory))
	copy(history, ta.alertHistory)
	
	return history
}

// TestConnection tests the Telegram connection
func (ta *TelegramAlert) TestConnection() error {
	if !ta.running {
		return fmt.Errorf("telegram alert system is not running")
	}
	
	url := fmt.Sprintf("%s/getMe", ta.apiURL)
	resp, err := ta.client.Get(url)
	if err != nil {
		return fmt.Errorf("failed to connect to Telegram: %w", err)
	}
	defer resp.Body.Close()
	
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read response: %w", err)
	}
	
	var telegramResp TelegramResponse
	if err := json.Unmarshal(body, &telegramResp); err != nil {
		return fmt.Errorf("failed to parse response: %w", err)
	}
	
	if !telegramResp.OK {
		return fmt.Errorf("telegram API error: %d - %s", telegramResp.ErrorCode, telegramResp.Description)
	}
	
	return nil
}

// SendTestAlert sends a test alert
func (ta *TelegramAlert) SendTestAlert() error {
	return ta.SendAlertWithType(AlertTypeSuccess, "Test Alert", "Telegram alert system is working correctly")
}

// Stop gracefully stops the alert system
func (ta *TelegramAlert) Stop() {
	ta.mu.Lock()
	defer ta.mu.Unlock()
	
	if !ta.running {
		return
	}
	
	ta.running = false
	close(ta.stopChan)
	
	log.Printf("Telegram alert system stopped - Final stats: %d total alerts, %d successful, %d failed", 
		ta.totalAlerts, ta.successfulAlerts, ta.failedAlerts)
}

// IsRunning returns whether the alert system is running
func (ta *TelegramAlert) IsRunning() bool {
	ta.mu.RLock()
	defer ta.mu.RUnlock()
	return ta.running
}