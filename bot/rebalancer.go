package main

import (
	"context"
	"fmt"
	"log"
	"math"
	"strconv"
	"sync"
	"time"
)

const (
	// RebalanceInterval is exactly 1 hour
	RebalanceInterval = 1 * time.Hour
	// PriceThreshold is exactly $5.00 USD
	PriceThreshold = 5.0
	// MonitorOnlyDuration is exactly 24 hours
	MonitorOnlyDuration = 24 * time.Hour
	// PriceUpdateInterval is 1 minute
	PriceUpdateInterval = 1 * time.Minute
	// MaxPriceHistory keeps last 60 price points
	MaxPriceHistory = 60
	// EmergencyStopThreshold is 500% above baseline
	EmergencyStopThreshold = 5.0
)

// RebalanceState represents the current state of the rebalancer
type RebalanceState int

const (
	StateActive RebalanceState = iota
	StateMonitorOnly
	StateEmergencyStop
	StateError
)

func (s RebalanceState) String() string {
	switch s {
	case StateActive:
		return "active"
	case StateMonitorOnly:
		return "monitor_only"
	case StateEmergencyStop:
		return "emergency_stop"
	case StateError:
		return "error"
	default:
		return "unknown"
	}
}

// Rebalancer handles automatic rebalancing with enhanced restrictions
type Rebalancer struct {
	config *BotConfig
	mu     sync.RWMutex
	
	// State management
	state               RebalanceState
	stateChangeTime     time.Time
	stateChangeReason   string
	
	// Price monitoring
	currentPrice        float64
	priceHistory        []float64
	lastPriceUpdate     time.Time
	priceUpdateErrors   int
	
	// Rebalancing state
	lastRebalance       time.Time
	rebalanceCount      int64
	nextRebalanceTime   time.Time
	totalRebalanceVolume float64
	
	// Monitor-only mode state
	monitorOnlyStart    time.Time
	monitorOnlyReason   string
	priceBreachTime     time.Time
	
	// Emergency state
	emergencyReason     string
	emergencyStartTime  time.Time
	
	// Alert integration
	telegramAlert       *TelegramAlert
	lastAlertTime       time.Time
	
	// Statistics
	dailyRebalanceCount int
	lastDailyReset      time.Time
	averagePrice        float64
	priceVolatility     float64
}

// NewRebalancer creates a new enhanced rebalancer instance
func NewRebalancer(config *BotConfig) *Rebalancer {
	return &Rebalancer{
		config:              config,
		state:               StateActive,
		stateChangeTime:     time.Now(),
		stateChangeReason:   "initialization",
		currentPrice:        3.0, // Default price below threshold
		priceHistory:        make([]float64, 0, MaxPriceHistory),
		lastRebalance:       time.Now(),
		nextRebalanceTime:   time.Now().Add(RebalanceInterval),
		lastDailyReset:      time.Now(),
		telegramAlert:       NewTelegramAlert(config),
	}
}

// Start starts the enhanced rebalancer with proper state management
func (r *Rebalancer) Start(ctx context.Context) error {
	log.Printf("Starting enhanced rebalancer with 1-hour intervals")
	
	// Send startup notification
	if err := r.sendStateChangeAlert("Rebalancer started", StateActive); err != nil {
		log.Printf("Failed to send startup alert: %v", err)
	}
	
	// Start price monitoring
	priceMonitorCtx, priceCancel := context.WithCancel(ctx)
	defer priceCancel()
	
	go r.monitorPrices(priceMonitorCtx)
	
	// Start daily reset routine
	go r.dailyResetRoutine(ctx)
	
	// Main rebalancing loop
	ticker := time.NewTicker(RebalanceInterval)
	defer ticker.Stop()
	
	for {
		select {
		case <-ctx.Done():
			log.Printf("Rebalancer stopping due to context cancellation")
			r.sendStateChangeAlert("Rebalancer stopped", StateError)
			return ctx.Err()
		case <-ticker.C:
			if err := r.processRebalanceCheck(ctx); err != nil {
				log.Printf("Error in rebalance check: %v", err)
				r.handleError(err)
			}
		}
	}
}

// monitorPrices continuously monitors GXR price with enhanced tracking
func (r *Rebalancer) monitorPrices(ctx context.Context) {
	ticker := time.NewTicker(PriceUpdateInterval)
	defer ticker.Stop()
	
	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			if err := r.updatePrice(ctx); err != nil {
				log.Printf("Error updating price: %v", err)
				r.priceUpdateErrors++
				if r.priceUpdateErrors >= 5 {
					r.handlePriceError("Too many price update failures")
				}
			} else {
				r.priceUpdateErrors = 0
			}
		}
	}
}

// updatePrice updates the current GXR price and checks thresholds
func (r *Rebalancer) updatePrice(ctx context.Context) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	
	// Simulate price fetching with realistic variation
	// In production, this would fetch from actual price sources
	basePrice := 3.0
	variation := 0.1 * (2.0*math.Sin(float64(time.Now().Unix())/3600) + 1.0)
	newPrice := basePrice + variation
	
	// Add some randomness
	if time.Now().UnixNano()%7 == 0 {
		newPrice += 0.5 * (float64(time.Now().UnixNano()%100) / 100.0)
	}
	
	r.currentPrice = newPrice
	r.lastPriceUpdate = time.Now()
	
	// Update price history
	r.priceHistory = append(r.priceHistory, newPrice)
	if len(r.priceHistory) > MaxPriceHistory {
		r.priceHistory = r.priceHistory[1:]
	}
	
	// Calculate statistics
	r.calculatePriceStatistics()
	
	// Check for price threshold breach
	if newPrice >= PriceThreshold && r.state == StateActive {
		r.enterMonitorOnlyMode(fmt.Sprintf("Price threshold breach: $%.2f >= $%.2f", newPrice, PriceThreshold))
	}
	
	// Check for emergency conditions
	if newPrice >= EmergencyStopThreshold && r.state != StateEmergencyStop {
		r.enterEmergencyStop(fmt.Sprintf("Emergency price threshold: $%.2f", newPrice))
	}
	
	return nil
}

// calculatePriceStatistics calculates average price and volatility
func (r *Rebalancer) calculatePriceStatistics() {
	if len(r.priceHistory) == 0 {
		return
	}
	
	// Calculate average
	sum := 0.0
	for _, price := range r.priceHistory {
		sum += price
	}
	r.averagePrice = sum / float64(len(r.priceHistory))
	
	// Calculate volatility (standard deviation)
	varianceSum := 0.0
	for _, price := range r.priceHistory {
		diff := price - r.averagePrice
		varianceSum += diff * diff
	}
	r.priceVolatility = math.Sqrt(varianceSum / float64(len(r.priceHistory)))
}

// processRebalanceCheck processes the hourly rebalance check
func (r *Rebalancer) processRebalanceCheck(ctx context.Context) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	
	now := time.Now()
	
	// Check if it's time to rebalance (exactly 1 hour)
	if now.Before(r.nextRebalanceTime) {
		return nil // Not time yet
	}
	
	// Update next rebalance time
	r.nextRebalanceTime = now.Add(RebalanceInterval)
	
	// Check current state
	switch r.state {
	case StateActive:
		return r.performRebalance(ctx)
	case StateMonitorOnly:
		return r.handleMonitorOnlyMode(ctx)
	case StateEmergencyStop:
		return r.handleEmergencyStop(ctx)
	case StateError:
		return r.handleErrorState(ctx)
	default:
		return fmt.Errorf("unknown rebalancer state: %v", r.state)
	}
}

// performRebalance performs the actual rebalancing when in active state
func (r *Rebalancer) performRebalance(ctx context.Context) error {
	log.Printf("Performing hourly rebalance - Price: $%.2f", r.currentPrice)
	
	// Check if we're still in acceptable price range
	if r.currentPrice >= PriceThreshold {
		return r.enterMonitorOnlyMode(fmt.Sprintf("Price threshold reached during rebalance: $%.2f", r.currentPrice))
	}
	
	// Perform rebalancing logic
	rebalanceVolume := r.calculateRebalanceVolume()
	
	// Execute rebalance
	if err := r.executeRebalance(ctx, rebalanceVolume); err != nil {
		return fmt.Errorf("rebalance execution failed: %w", err)
	}
	
	// Update statistics
	r.lastRebalance = time.Now()
	r.rebalanceCount++
	r.dailyRebalanceCount++
	r.totalRebalanceVolume += rebalanceVolume
	
	log.Printf("Rebalance completed - Volume: %.2f GXR, Total: %d", rebalanceVolume, r.rebalanceCount)
	
	return nil
}

// calculateRebalanceVolume calculates the volume to rebalance
func (r *Rebalancer) calculateRebalanceVolume() float64 {
	// Simple volume calculation based on price volatility
	baseVolume := 1000.0 // 1000 GXR base volume
	volatilityMultiplier := 1.0 + r.priceVolatility
	
	return baseVolume * volatilityMultiplier
}

// executeRebalance executes the actual rebalancing operation
func (r *Rebalancer) executeRebalance(ctx context.Context, volume float64) error {
	// Simulate rebalancing - in production this would interact with DEX
	log.Printf("Executing rebalance of %.2f GXR", volume)
	
	// Simulate processing time
	time.Sleep(100 * time.Millisecond)
	
	// Simulate potential errors
	if time.Now().UnixNano()%100 == 0 {
		return fmt.Errorf("simulated rebalance error")
	}
	
	return nil
}

// handleMonitorOnlyMode handles the bot when in monitor-only mode
func (r *Rebalancer) handleMonitorOnlyMode(ctx context.Context) error {
	elapsed := time.Since(r.monitorOnlyStart)
	
	log.Printf("Monitor-only mode - Elapsed: %v, Price: $%.2f", elapsed, r.currentPrice)
	
	// Check if 24 hours have passed
	if elapsed >= MonitorOnlyDuration {
		// Check if price is back below threshold
		if r.currentPrice < PriceThreshold {
			return r.exitMonitorOnlyMode("24-hour period elapsed and price below threshold")
		} else {
			// Extend monitor-only period
			r.monitorOnlyStart = time.Now()
			r.sendStateChangeAlert(fmt.Sprintf("Monitor-only mode extended - Price: $%.2f", r.currentPrice), StateMonitorOnly)
		}
	}
	
	return nil
}

// handleEmergencyStop handles emergency stop conditions
func (r *Rebalancer) handleEmergencyStop(ctx context.Context) error {
	log.Printf("Emergency stop active - Price: $%.2f", r.currentPrice)
	
	// Check if conditions have normalized
	if r.currentPrice < PriceThreshold {
		return r.exitEmergencyStop("Price returned to normal levels")
	}
	
	return nil
}

// handleErrorState handles error state recovery
func (r *Rebalancer) handleErrorState(ctx context.Context) error {
	log.Printf("Error state active - attempting recovery")
	
	// Simple recovery logic - reset to active after 1 hour
	if time.Since(r.stateChangeTime) >= time.Hour {
		return r.recoverFromError("Auto-recovery after 1 hour")
	}
	
	return nil
}

// State transition methods

// enterMonitorOnlyMode transitions to monitor-only mode
func (r *Rebalancer) enterMonitorOnlyMode(reason string) error {
	r.state = StateMonitorOnly
	r.stateChangeTime = time.Now()
	r.stateChangeReason = reason
	r.monitorOnlyStart = time.Now()
	r.monitorOnlyReason = reason
	r.priceBreachTime = time.Now()
	
	log.Printf("Entering monitor-only mode: %s", reason)
	return r.sendStateChangeAlert(reason, StateMonitorOnly)
}

// exitMonitorOnlyMode transitions out of monitor-only mode
func (r *Rebalancer) exitMonitorOnlyMode(reason string) error {
	r.state = StateActive
	r.stateChangeTime = time.Now()
	r.stateChangeReason = reason
	
	log.Printf("Exiting monitor-only mode: %s", reason)
	return r.sendStateChangeAlert(reason, StateActive)
}

// enterEmergencyStop transitions to emergency stop
func (r *Rebalancer) enterEmergencyStop(reason string) error {
	r.state = StateEmergencyStop
	r.stateChangeTime = time.Now()
	r.stateChangeReason = reason
	r.emergencyReason = reason
	r.emergencyStartTime = time.Now()
	
	log.Printf("EMERGENCY STOP: %s", reason)
	return r.sendStateChangeAlert(fmt.Sprintf("EMERGENCY: %s", reason), StateEmergencyStop)
}

// exitEmergencyStop transitions out of emergency stop
func (r *Rebalancer) exitEmergencyStop(reason string) error {
	r.state = StateActive
	r.stateChangeTime = time.Now()
	r.stateChangeReason = reason
	
	log.Printf("Exiting emergency stop: %s", reason)
	return r.sendStateChangeAlert(fmt.Sprintf("Recovery: %s", reason), StateActive)
}

// recoverFromError recovers from error state
func (r *Rebalancer) recoverFromError(reason string) error {
	r.state = StateActive
	r.stateChangeTime = time.Now()
	r.stateChangeReason = reason
	
	log.Printf("Recovering from error: %s", reason)
	return r.sendStateChangeAlert(fmt.Sprintf("Recovery: %s", reason), StateActive)
}

// handleError handles general errors
func (r *Rebalancer) handleError(err error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	
	r.state = StateError
	r.stateChangeTime = time.Now()
	r.stateChangeReason = err.Error()
	
	log.Printf("Rebalancer error: %v", err)
	r.sendStateChangeAlert(fmt.Sprintf("Error: %v", err), StateError)
}

// handlePriceError handles price-related errors
func (r *Rebalancer) handlePriceError(reason string) {
	r.mu.Lock()
	defer r.mu.Unlock()
	
	r.state = StateError
	r.stateChangeTime = time.Now()
	r.stateChangeReason = reason
	
	log.Printf("Price error: %s", reason)
	r.sendStateChangeAlert(fmt.Sprintf("Price Error: %s", reason), StateError)
}

// sendStateChangeAlert sends telegram alert for state changes
func (r *Rebalancer) sendStateChangeAlert(message string, newState RebalanceState) error {
	if r.telegramAlert == nil {
		return nil
	}
	
	// Rate limiting - don't send alerts too frequently
	if time.Since(r.lastAlertTime) < 5*time.Minute {
		return nil
	}
	
	fullMessage := fmt.Sprintf("🔄 Rebalancer State Change\n\nState: %s\nReason: %s\nPrice: $%.2f\nTime: %s",
		newState.String(),
		message,
		r.currentPrice,
		time.Now().Format("2006-01-02 15:04:05"),
	)
	
	if err := r.telegramAlert.SendAlert(fullMessage); err != nil {
		log.Printf("Failed to send state change alert: %v", err)
		return err
	}
	
	r.lastAlertTime = time.Now()
	return nil
}

// dailyResetRoutine resets daily counters
func (r *Rebalancer) dailyResetRoutine(ctx context.Context) {
	ticker := time.NewTicker(24 * time.Hour)
	defer ticker.Stop()
	
	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			r.mu.Lock()
			log.Printf("Daily reset - Rebalances today: %d", r.dailyRebalanceCount)
			r.dailyRebalanceCount = 0
			r.lastDailyReset = time.Now()
			r.mu.Unlock()
		}
	}
}

// GetStatus returns current rebalancer status
func (r *Rebalancer) GetStatus() map[string]interface{} {
	r.mu.RLock()
	defer r.mu.RUnlock()
	
	return map[string]interface{}{
		"state":                 r.state.String(),
		"state_change_time":     r.stateChangeTime.Format(time.RFC3339),
		"state_change_reason":   r.stateChangeReason,
		"current_price":         r.currentPrice,
		"last_price_update":     r.lastPriceUpdate.Format(time.RFC3339),
		"price_history_count":   len(r.priceHistory),
		"average_price":         r.averagePrice,
		"price_volatility":      r.priceVolatility,
		"last_rebalance":        r.lastRebalance.Format(time.RFC3339),
		"next_rebalance":        r.nextRebalanceTime.Format(time.RFC3339),
		"rebalance_count":       r.rebalanceCount,
		"daily_rebalance_count": r.dailyRebalanceCount,
		"total_volume":          r.totalRebalanceVolume,
		"monitor_only_start":    r.monitorOnlyStart.Format(time.RFC3339),
		"monitor_only_reason":   r.monitorOnlyReason,
		"emergency_reason":      r.emergencyReason,
		"emergency_start":       r.emergencyStartTime.Format(time.RFC3339),
	}
}

// Stop gracefully stops the rebalancer
func (r *Rebalancer) Stop() {
	r.mu.Lock()
	defer r.mu.Unlock()
	
	log.Printf("Stopping rebalancer - Final stats: %d rebalances, $%.2f total volume", 
		r.rebalanceCount, r.totalRebalanceVolume)
	
	r.sendStateChangeAlert("Rebalancer stopped", StateError)
}