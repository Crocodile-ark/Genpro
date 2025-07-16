package main

import (
	"context"
	"log"
	"strconv"
	"time"
)

// Rebalancer handles automatic rebalancing between chains
type Rebalancer struct {
	config *BotConfig
	
	// Rebalancer state
	emergencyMode     bool
	lastRebalance     time.Time
	rebalanceCount    int64
	dailySwapCount    int64
	lastSwapTime      time.Time
}

// NewRebalancer creates a new rebalancer instance
func NewRebalancer(config *BotConfig) *Rebalancer {
	return &Rebalancer{
		config: config,
	}
}

// Initialize initializes the rebalancer
func (r *Rebalancer) Initialize() error {
	log.Println("Initializing Rebalancer...")
	
	r.emergencyMode = r.config.EmergencyMode
	r.lastRebalance = time.Now()
	
	log.Printf("Rebalancer initialized (Emergency mode: %v)", r.emergencyMode)
	return nil
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
	// Check if we're in emergency mode
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

// shouldEnterEmergencyMode checks if we should enter emergency mode
func (r *Rebalancer) shouldEnterEmergencyMode() bool {
	// Check GXR price against limit
	currentPrice := r.getCurrentGXRPrice()
	priceLimit, _ := strconv.ParseFloat(r.config.PriceLimit, 64)
	
	return currentPrice > priceLimit
}

// shouldExitEmergencyMode checks if we should exit emergency mode
func (r *Rebalancer) shouldExitEmergencyMode() bool {
	// Check if 24 hours have passed and price is stable
	if time.Since(r.lastRebalance) < (24 * time.Hour) {
		return false
	}
	
	currentPrice := r.getCurrentGXRPrice()
	priceLimit, _ := strconv.ParseFloat(r.config.PriceLimit, 64)
	
	return currentPrice <= priceLimit
}

// enterEmergencyMode activates emergency mode
func (r *Rebalancer) enterEmergencyMode() {
	if !r.emergencyMode {
		log.Println("🚨 Entering emergency mode - price limit exceeded")
		r.emergencyMode = true
		// TODO: Send Telegram alert
	}
}

// exitEmergencyMode deactivates emergency mode
func (r *Rebalancer) exitEmergencyMode() {
	log.Println("✅ Exiting emergency mode - conditions normalized")
	r.emergencyMode = false
	
	// Perform immediate rebalancing after emergency mode
	r.performRebalancing()
}

// hasReachedDailyLimit checks if daily swap limit is reached
func (r *Rebalancer) hasReachedDailyLimit() bool {
	// Reset daily count if new day
	now := time.Now()
	if now.Day() != r.lastSwapTime.Day() {
		r.dailySwapCount = 0
	}
	
	// Parse max daily swap amount
	// TODO: Parse actual amount from config.MaxSwapDaily
	maxDaily := int64(10000) // 10,000 GXR limit
	
	return r.dailySwapCount >= maxDaily
}

// getCurrentGXRPrice gets the current GXR price (placeholder)
func (r *Rebalancer) getCurrentGXRPrice() float64 {
	// TODO: Implement actual price fetching from DEX or oracle
	return 5.0 // Placeholder price
}

// performRebalancing performs the actual rebalancing
func (r *Rebalancer) performRebalancing() error {
	log.Println("Performing inter-chain rebalancing...")
	
	// TODO: Implement actual rebalancing logic
	// This would involve:
	// 1. Checking pool imbalances across chains
	// 2. Calculating optimal rebalancing amounts
	// 3. Executing cross-chain swaps
	// 4. Updating pool states
	
	r.lastRebalance = time.Now()
	r.lastSwapTime = time.Now()
	r.rebalanceCount++
	r.dailySwapCount++
	
	log.Printf("Rebalancing completed (cycle %d)", r.rebalanceCount)
	return nil
}

// GetStatus returns the current rebalancer status
func (r *Rebalancer) GetStatus() map[string]interface{} {
	return map[string]interface{}{
		"emergency_mode":    r.emergencyMode,
		"last_rebalance":    r.lastRebalance,
		"rebalance_count":   r.rebalanceCount,
		"daily_swap_count":  r.dailySwapCount,
		"last_swap_time":    r.lastSwapTime,
		"price_limit":       r.config.PriceLimit,
		"max_daily_swap":    r.config.MaxSwapDaily,
		"swap_cooldown":     r.config.SwapCooldown,
	}
}