package main

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"strconv"
	"time"
)

// Rebalancer handles automatic rebalancing between chains
type Rebalancer struct {
	config *BotConfig
	
	// Price monitoring
	currentPrice      float64
	priceHistory      []float64
	lastPriceUpdate   time.Time
	
	// Rebalancer state
	emergencyMode     bool
	lastRebalance     time.Time
	rebalanceCount    int64
	dailySwapCount    int64
	lastSwapTime      time.Time
	emergencyStartTime time.Time
	
	// Telegram alert integration
	telegramAlert     *TelegramAlert
}

// NewRebalancer creates a new rebalancer instance
func NewRebalancer(config *BotConfig) *Rebalancer {
	return &Rebalancer{
		config: config,
		currentPrice: 5.0, // Default price
		priceHistory: make([]float64, 0),
	}
}

// Initialize initializes the rebalancer
func (r *Rebalancer) Initialize() error {
	log.Println("Initializing Rebalancer...")
	
	// Validate configuration
	if err := r.validateConfig(); err != nil {
		return fmt.Errorf("invalid rebalancer configuration: %w", err)
	}
	
	r.emergencyMode = r.config.EmergencyMode
	r.lastRebalance = time.Now()
	r.lastPriceUpdate = time.Now()
	
	// Initialize price monitoring
	r.startPriceMonitoring()
	
	log.Printf("Rebalancer initialized (Emergency mode: %v)", r.emergencyMode)
	return nil
}

// validateConfig validates the rebalancer configuration
func (r *Rebalancer) validateConfig() error {
	if r.config.MaxSwapDaily == "" {
		return fmt.Errorf("max_swap_daily is required")
	}
	
	if r.config.PriceLimit == "" {
		return fmt.Errorf("price_limit is required")
	}
	
	if r.config.SwapCooldown <= 0 {
		return fmt.Errorf("swap_cooldown must be positive")
	}
	
	return nil
}

// startPriceMonitoring starts the price monitoring goroutine
func (r *Rebalancer) startPriceMonitoring() {
	// In a real implementation, this would connect to DEX APIs or oracles
	// For now, we'll simulate price fluctuations
	log.Println("Starting price monitoring...")
}

// SetTelegramAlert sets the telegram alert instance
func (r *Rebalancer) SetTelegramAlert(alert *TelegramAlert) {
	r.telegramAlert = alert
}

// Start starts the rebalancer service
func (r *Rebalancer) Start(ctx context.Context) error {
	log.Println("Starting Rebalancer service...")
	
	ticker := time.NewTicker(r.config.CheckInterval)
	defer ticker.Stop()
	
	for {
		select {
		case <-ctx.Done():
			log.Println("Rebalancer stopping...")
			return nil
			
		case <-ticker.C:
			if err := r.checkAndRebalance(); err != nil {
				log.Printf("Rebalancer error: %v", err)
			}
		}
	}
}

// checkAndRebalance checks prices and rebalances if needed
func (r *Rebalancer) checkAndRebalance() error {
	// Update current price
	r.updateCurrentPrice()
	
	// Check if we should enter emergency mode
	if r.shouldEnterEmergencyMode() {
		r.enterEmergencyMode()
		return nil
	}
	
	// Exit emergency mode if conditions are met
	if r.emergencyMode && r.shouldExitEmergencyMode() {
		r.exitEmergencyMode()
	}
	
	// Skip rebalancing if in emergency mode
	if r.emergencyMode {
		return nil
	}
	
	// Check daily limits
	if r.hasReachedDailyLimit() {
		log.Println("Daily swap limit reached, skipping rebalance")
		return nil
	}
	
	// Check cooldown period
	if time.Since(r.lastSwapTime) < r.config.SwapCooldown {
		return nil // Still in cooldown
	}
	
	// Perform rebalancing
	return r.performRebalancing()
}

// updateCurrentPrice updates the current GXR price
func (r *Rebalancer) updateCurrentPrice() {
	// Simulate price fluctuations for demo purposes
	// In a real implementation, this would fetch from DEX APIs or oracles
	
	// Add some random price movement
	change := (rand.Float64() - 0.5) * 0.1 // ±5% change
	r.currentPrice = r.currentPrice * (1 + change)
	
	// Ensure price stays within reasonable bounds
	if r.currentPrice < 1.0 {
		r.currentPrice = 1.0
	} else if r.currentPrice > 20.0 {
		r.currentPrice = 20.0
	}
	
	// Update price history
	r.priceHistory = append(r.priceHistory, r.currentPrice)
	if len(r.priceHistory) > 100 {
		r.priceHistory = r.priceHistory[1:]
	}
	
	r.lastPriceUpdate = time.Now()
}

// shouldEnterEmergencyMode checks if we should enter emergency mode
func (r *Rebalancer) shouldEnterEmergencyMode() bool {
	priceLimit, err := strconv.ParseFloat(r.config.PriceLimit, 64)
	if err != nil {
		log.Printf("Error parsing price limit: %v", err)
		return false
	}
	
	return !r.emergencyMode && r.currentPrice > priceLimit
}

// shouldExitEmergencyMode checks if we should exit emergency mode
func (r *Rebalancer) shouldExitEmergencyMode() bool {
	// Check if 24 hours have passed since entering emergency mode
	if time.Since(r.emergencyStartTime) < (24 * time.Hour) {
		return false
	}
	
	priceLimit, err := strconv.ParseFloat(r.config.PriceLimit, 64)
	if err != nil {
		log.Printf("Error parsing price limit: %v", err)
		return false
	}
	
	return r.currentPrice <= priceLimit
}

// enterEmergencyMode activates emergency mode
func (r *Rebalancer) enterEmergencyMode() {
	if !r.emergencyMode {
		log.Println("🚨 Entering emergency mode - price limit exceeded")
		r.emergencyMode = true
		r.emergencyStartTime = time.Now()
		
		// Send Telegram alert
		if r.telegramAlert != nil {
			priceLimit, _ := strconv.ParseFloat(r.config.PriceLimit, 64)
			r.telegramAlert.SendPriceAlert(r.currentPrice, priceLimit)
		}
	}
}

// exitEmergencyMode deactivates emergency mode
func (r *Rebalancer) exitEmergencyMode() {
	log.Println("✅ Exiting emergency mode - conditions normalized")
	r.emergencyMode = false
	
	// Send Telegram alert
	if r.telegramAlert != nil {
		r.telegramAlert.SendAlert(fmt.Sprintf("✅ Emergency mode deactivated\n\nPrice: $%.2f\nTime in emergency: %v", 
			r.currentPrice, time.Since(r.emergencyStartTime)))
	}
	
	// Perform immediate rebalancing after emergency mode
	if err := r.performRebalancing(); err != nil {
		log.Printf("Error during post-emergency rebalancing: %v", err)
	}
}

// hasReachedDailyLimit checks if daily swap limit is reached
func (r *Rebalancer) hasReachedDailyLimit() bool {
	// Reset daily count if new day
	now := time.Now()
	if now.Day() != r.lastSwapTime.Day() {
		r.dailySwapCount = 0
	}
	
	// Parse max daily swap amount
	maxDailyStr := r.config.MaxSwapDaily
	if maxDailyStr == "" {
		maxDailyStr = "10000"
	}
	
	maxDaily, err := strconv.ParseInt(maxDailyStr, 10, 64)
	if err != nil {
		log.Printf("Error parsing max daily swap: %v", err)
		maxDaily = 10000 // Default to 10,000 GXR
	}
	
	return r.dailySwapCount >= maxDaily
}

// performRebalancing performs the actual rebalancing
func (r *Rebalancer) performRebalancing() error {
	log.Println("Performing inter-chain rebalancing...")
	
	// Simulate rebalancing logic
	if err := r.simulateRebalancing(); err != nil {
		return fmt.Errorf("rebalancing simulation failed: %w", err)
	}
	
	r.lastRebalance = time.Now()
	r.lastSwapTime = time.Now()
	r.rebalanceCount++
	r.dailySwapCount++
	
	log.Printf("Rebalancing completed (cycle %d)", r.rebalanceCount)
	
	// Send success alert
	if r.telegramAlert != nil {
		r.telegramAlert.SendAlert(fmt.Sprintf("⚖️ Rebalancing completed\n\nCycle: %d\nPrice: $%.2f\nDaily swaps: %d", 
			r.rebalanceCount, r.currentPrice, r.dailySwapCount))
	}
	
	return nil
}

// simulateRebalancing simulates the rebalancing process
func (r *Rebalancer) simulateRebalancing() error {
	// Simulate checking pool imbalances
	log.Println("Checking pool imbalances...")
	time.Sleep(1 * time.Second)
	
	// Simulate calculating optimal amounts
	log.Println("Calculating optimal rebalancing amounts...")
	time.Sleep(1 * time.Second)
	
	// Simulate executing swaps
	log.Println("Executing cross-chain swaps...")
	time.Sleep(2 * time.Second)
	
	// Simulate potential failures
	if r.rebalanceCount > 0 && r.rebalanceCount%20 == 0 {
		return fmt.Errorf("simulated rebalancing failure")
	}
	
	return nil
}

// GetStatus returns the current rebalancer status
func (r *Rebalancer) GetStatus() map[string]interface{} {
	priceLimit, _ := strconv.ParseFloat(r.config.PriceLimit, 64)
	maxDaily, _ := strconv.ParseInt(r.config.MaxSwapDaily, 10, 64)
	
	status := map[string]interface{}{
		"emergency_mode":     r.emergencyMode,
		"current_price":      r.currentPrice,
		"price_limit":        priceLimit,
		"last_rebalance":     r.lastRebalance,
		"rebalance_count":    r.rebalanceCount,
		"daily_swap_count":   r.dailySwapCount,
		"max_daily_swap":     maxDaily,
		"last_swap_time":     r.lastSwapTime,
		"swap_cooldown":      r.config.SwapCooldown,
		"last_price_update":  r.lastPriceUpdate,
		"price_history_length": len(r.priceHistory),
	}
	
	if r.emergencyMode {
		status["emergency_duration"] = time.Since(r.emergencyStartTime).String()
	}
	
	return status
}

// GetPriceHistory returns the recent price history
func (r *Rebalancer) GetPriceHistory() []float64 {
	return r.priceHistory
}

// ForceRebalance forces a manual rebalancing (for testing/emergency)
func (r *Rebalancer) ForceRebalance() error {
	if r.emergencyMode {
		return fmt.Errorf("cannot force rebalance in emergency mode")
	}
	
	log.Println("Forcing manual rebalancing...")
	
	if err := r.performRebalancing(); err != nil {
		return fmt.Errorf("forced rebalancing failed: %w", err)
	}
	
	log.Println("Manual rebalancing completed successfully")
	return nil
}