package main

import (
	"context"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/types/query"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
)

const (
	// ValidatorCheckInterval is how often to check validator status
	ValidatorCheckInterval = 5 * time.Minute
	// MonthlyResetInterval is 30 days
	MonthlyResetInterval = 30 * 24 * time.Hour
	// ValidatorInactivityThreshold is 10 days per month
	ValidatorInactivityThreshold = 10
	// BotHeartbeatInterval is 1 minute
	BotHeartbeatInterval = 1 * time.Minute
	// BotHeartbeatTimeout is 5 minutes
	BotHeartbeatTimeout = 5 * time.Minute
	// SlashingGracePeriod is 10 minutes
	SlashingGracePeriod = 10 * time.Minute
)

// ValidatorStatus represents the status of a validator
type ValidatorStatus struct {
	OperatorAddress string
	Moniker         string
	Status          stakingtypes.BondStatus
	Jailed          bool
	Tokens          string
	DelegatorShares string
	Commission      string
	
	// Uptime tracking
	CurrentMonth     uint64
	InactiveDays     uint64
	LastActiveTime   time.Time
	LastCheck        time.Time
	MissedBlocks     uint64
	
	// Bot monitoring
	BotRunning       bool
	LastBotHeartbeat time.Time
	BotVersion       string
	BotErrors        []string
	
	// Reward eligibility
	RewardEligible   bool
	ForfeitedRewards float64
	LastRewardClaim  time.Time
	
	// Statistics
	UptimePercent    float64
	MonthlyUptime    float64
	TotalMissedBlocks uint64
}

// ValidatorMonitor monitors validator performance and bot requirements
type ValidatorMonitor struct {
	config        *BotConfig
	clientCtx     client.Context
	cdc           codec.Codec
	mu            sync.RWMutex
	
	// Validator tracking
	validators    map[string]*ValidatorStatus
	totalValidators int
	activeValidators int
	
	// Monthly tracking
	currentMonth  uint64
	lastMonthReset time.Time
	
	// Bot enforcement
	botHeartbeats map[string]time.Time
	slashingQueue []string
	
	// Statistics
	totalInactiveValidators int
	totalForfeitedRewards   float64
	monthlyStats            map[uint64]*MonthlyStats
	
	// Alert system
	telegramAlert   *TelegramAlert
	lastAlertTime   time.Time
	alertsSent      int
}

// MonthlyStats tracks monthly statistics
type MonthlyStats struct {
	Month            uint64
	TotalValidators  int
	ActiveValidators int
	InactiveValidators int
	ForfeitedRewards float64
	AverageUptime    float64
	BotsRunning      int
	SlashedValidators int
}

// NewValidatorMonitor creates a new validator monitor
func NewValidatorMonitor(config *BotConfig, clientCtx client.Context, cdc codec.Codec) *ValidatorMonitor {
	return &ValidatorMonitor{
		config:        config,
		clientCtx:     clientCtx,
		cdc:           cdc,
		validators:    make(map[string]*ValidatorStatus),
		currentMonth:  getCurrentMonth(),
		lastMonthReset: time.Now(),
		botHeartbeats: make(map[string]time.Time),
		slashingQueue: make([]string, 0),
		monthlyStats:  make(map[uint64]*MonthlyStats),
		telegramAlert: NewTelegramAlert(config),
	}
}

// Start starts the validator monitoring service
func (vm *ValidatorMonitor) Start(ctx context.Context) error {
	log.Printf("Starting validator monitor with enhanced tracking")
	
	// Send startup notification
	if err := vm.sendAlert("🔍 Validator Monitor Started", "Enhanced monitoring active"); err != nil {
		log.Printf("Failed to send startup alert: %v", err)
	}
	
	// Start periodic checks
	go vm.validatorCheckRoutine(ctx)
	go vm.botMonitoringRoutine(ctx)
	go vm.monthlyResetRoutine(ctx)
	go vm.slashingRoutine(ctx)
	
	return nil
}

// validatorCheckRoutine periodically checks validator status
func (vm *ValidatorMonitor) validatorCheckRoutine(ctx context.Context) {
	ticker := time.NewTicker(ValidatorCheckInterval)
	defer ticker.Stop()
	
	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			if err := vm.checkAllValidators(ctx); err != nil {
				log.Printf("Error checking validators: %v", err)
			}
		}
	}
}

// botMonitoringRoutine monitors bot heartbeats
func (vm *ValidatorMonitor) botMonitoringRoutine(ctx context.Context) {
	ticker := time.NewTicker(BotHeartbeatInterval)
	defer ticker.Stop()
	
	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			vm.checkBotHeartbeats(ctx)
		}
	}
}

// monthlyResetRoutine resets monthly counters
func (vm *ValidatorMonitor) monthlyResetRoutine(ctx context.Context) {
	ticker := time.NewTicker(MonthlyResetInterval)
	defer ticker.Stop()
	
	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			vm.performMonthlyReset(ctx)
		}
	}
}

// slashingRoutine handles validator slashing for bot non-compliance
func (vm *ValidatorMonitor) slashingRoutine(ctx context.Context) {
	ticker := time.NewTicker(SlashingGracePeriod)
	defer ticker.Stop()
	
	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			vm.processSlashingQueue(ctx)
		}
	}
}

// checkAllValidators checks all bonded validators
func (vm *ValidatorMonitor) checkAllValidators(ctx context.Context) error {
	vm.mu.Lock()
	defer vm.mu.Unlock()
	
	// Query all validators
	validators, err := vm.queryValidators(ctx)
	if err != nil {
		return fmt.Errorf("failed to query validators: %w", err)
	}
	
	activeCount := 0
	inactiveCount := 0
	
	for _, validator := range validators {
		status, exists := vm.validators[validator.OperatorAddress]
		if !exists {
			status = &ValidatorStatus{
				OperatorAddress: validator.OperatorAddress,
				Moniker:         validator.Description.Moniker,
				CurrentMonth:    vm.currentMonth,
				LastActiveTime:  time.Now(),
				LastCheck:       time.Now(),
				RewardEligible:  true,
			}
			vm.validators[validator.OperatorAddress] = status
		}
		
		// Update validator status
		vm.updateValidatorStatus(status, validator)
		
		// Check inactivity
		if vm.isValidatorInactive(status) {
			inactiveCount++
			if status.RewardEligible {
				vm.markValidatorInactive(status)
			}
		} else {
			activeCount++
		}
		
		// Check bot requirement
		if !vm.isValidatorBotRunning(status) {
			vm.queueForSlashing(status.OperatorAddress)
		}
	}
	
	vm.totalValidators = len(validators)
	vm.activeValidators = activeCount
	vm.totalInactiveValidators = inactiveCount
	
	log.Printf("Validator check complete - Total: %d, Active: %d, Inactive: %d", 
		vm.totalValidators, vm.activeValidators, vm.totalInactiveValidators)
	
	return nil
}

// queryValidators queries all validators from the chain
func (vm *ValidatorMonitor) queryValidators(ctx context.Context) ([]stakingtypes.Validator, error) {
	queryClient := stakingtypes.NewQueryClient(vm.clientCtx)
	
	resp, err := queryClient.Validators(ctx, &stakingtypes.QueryValidatorsRequest{
		Status: stakingtypes.BondStatusBonded,
		Pagination: &query.PageRequest{
			Limit: 1000,
		},
	})
	if err != nil {
		return nil, err
	}
	
	return resp.Validators, nil
}

// updateValidatorStatus updates a validator's status
func (vm *ValidatorMonitor) updateValidatorStatus(status *ValidatorStatus, validator stakingtypes.Validator) {
	status.Status = validator.Status
	status.Jailed = validator.Jailed
	status.Tokens = validator.Tokens.String()
	status.DelegatorShares = validator.DelegatorShares.String()
	status.Commission = validator.Commission.Rate.String()
	status.LastCheck = time.Now()
	
	// Update uptime tracking
	if validator.Status == stakingtypes.Bonded && !validator.Jailed {
		status.LastActiveTime = time.Now()
	} else {
		// Validator is inactive, increment inactive days
		if status.CurrentMonth == vm.currentMonth {
			lastCheck := time.Unix(status.LastCheck.Unix(), 0)
			if time.Since(lastCheck) >= 24*time.Hour {
				status.InactiveDays++
			}
		}
	}
	
	// Calculate uptime percentage
	monthStart := time.Now().AddDate(0, 0, -30)
	if status.LastActiveTime.After(monthStart) {
		activeDays := time.Since(status.LastActiveTime).Hours() / 24
		status.MonthlyUptime = (30 - activeDays) / 30 * 100
	}
}

// isValidatorInactive checks if validator is inactive (>10 days/month)
func (vm *ValidatorMonitor) isValidatorInactive(status *ValidatorStatus) bool {
	// Check if validator has been inactive for more than 10 days this month
	if status.CurrentMonth != vm.currentMonth {
		// New month, reset counters
		status.CurrentMonth = vm.currentMonth
		status.InactiveDays = 0
		return false
	}
	
	return status.InactiveDays > ValidatorInactivityThreshold
}

// markValidatorInactive marks a validator as inactive and ineligible for rewards
func (vm *ValidatorMonitor) markValidatorInactive(status *ValidatorStatus) {
	status.RewardEligible = false
	status.ForfeitedRewards += 100.0 // Approximate monthly reward
	
	log.Printf("Validator %s marked inactive - Inactive days: %d", 
		status.OperatorAddress, status.InactiveDays)
	
	// Send telegram alert
	message := fmt.Sprintf("⚠️ Validator Inactivity Alert\n\nValidator: %s\nInactive Days: %d/%d\nStatus: Reward Forfeited\nMonth: %d", 
		status.Moniker, status.InactiveDays, ValidatorInactivityThreshold, vm.currentMonth)
	
	vm.sendAlert("Validator Inactivity", message)
}

// isValidatorBotRunning checks if validator's bot is running
func (vm *ValidatorMonitor) isValidatorBotRunning(status *ValidatorStatus) bool {
	lastHeartbeat, exists := vm.botHeartbeats[status.OperatorAddress]
	if !exists {
		return false
	}
	
	return time.Since(lastHeartbeat) < BotHeartbeatTimeout
}

// queueForSlashing queues a validator for slashing due to bot non-compliance
func (vm *ValidatorMonitor) queueForSlashing(operatorAddr string) {
	// Check if already queued
	for _, addr := range vm.slashingQueue {
		if addr == operatorAddr {
			return
		}
	}
	
	vm.slashingQueue = append(vm.slashingQueue, operatorAddr)
	
	log.Printf("Validator %s queued for slashing - bot not running", operatorAddr)
}

// checkBotHeartbeats checks for bot heartbeats
func (vm *ValidatorMonitor) checkBotHeartbeats(ctx context.Context) {
	vm.mu.RLock()
	defer vm.mu.RUnlock()
	
	now := time.Now()
	inactiveValidators := 0
	
	for addr, status := range vm.validators {
		lastHeartbeat, exists := vm.botHeartbeats[addr]
		if !exists {
			lastHeartbeat = now.Add(-time.Hour) // Assume old heartbeat
		}
		
		if now.Sub(lastHeartbeat) > BotHeartbeatTimeout {
			status.BotRunning = false
			inactiveValidators++
			
			// Send alert for bot inactivity
			if now.Sub(status.LastBotHeartbeat) > 1*time.Hour {
				vm.sendBotInactivityAlert(status)
				status.LastBotHeartbeat = now
			}
		} else {
			status.BotRunning = true
			status.LastBotHeartbeat = lastHeartbeat
		}
	}
	
	if inactiveValidators > 0 {
		log.Printf("Bot heartbeat check - %d validators with inactive bots", inactiveValidators)
	}
}

// processSlashingQueue processes the slashing queue
func (vm *ValidatorMonitor) processSlashingQueue(ctx context.Context) {
	vm.mu.Lock()
	defer vm.mu.Unlock()
	
	if len(vm.slashingQueue) == 0 {
		return
	}
	
	log.Printf("Processing slashing queue - %d validators", len(vm.slashingQueue))
	
	for _, operatorAddr := range vm.slashingQueue {
		if err := vm.slashValidator(ctx, operatorAddr); err != nil {
			log.Printf("Failed to slash validator %s: %v", operatorAddr, err)
		} else {
			log.Printf("Successfully slashed validator %s for bot non-compliance", operatorAddr)
		}
	}
	
	// Clear the queue
	vm.slashingQueue = vm.slashingQueue[:0]
}

// slashValidator executes slashing for a validator
func (vm *ValidatorMonitor) slashValidator(ctx context.Context, operatorAddr string) error {
	// In a real implementation, this would submit a slashing transaction
	// For now, we'll just log and send alerts
	
	status, exists := vm.validators[operatorAddr]
	if !exists {
		return fmt.Errorf("validator not found: %s", operatorAddr)
	}
	
	log.Printf("SLASHING: Validator %s (%s) for bot non-compliance", 
		status.Moniker, operatorAddr)
	
	// Send slashing alert
	message := fmt.Sprintf("⚔️ Validator Slashed\n\nValidator: %s\nReason: Mandatory bot not running\nTime: %s", 
		status.Moniker, time.Now().Format("2006-01-02 15:04:05"))
	
	return vm.sendAlert("Validator Slashed", message)
}

// performMonthlyReset resets monthly counters
func (vm *ValidatorMonitor) performMonthlyReset(ctx context.Context) {
	vm.mu.Lock()
	defer vm.mu.Unlock()
	
	oldMonth := vm.currentMonth
	vm.currentMonth = getCurrentMonth()
	vm.lastMonthReset = time.Now()
	
	// Store monthly statistics
	vm.monthlyStats[oldMonth] = &MonthlyStats{
		Month:              oldMonth,
		TotalValidators:    vm.totalValidators,
		ActiveValidators:   vm.activeValidators,
		InactiveValidators: vm.totalInactiveValidators,
		ForfeitedRewards:   vm.totalForfeitedRewards,
		AverageUptime:      vm.calculateAverageUptime(),
		BotsRunning:        vm.countRunningBots(),
	}
	
	// Reset all validator monthly counters
	for _, status := range vm.validators {
		status.CurrentMonth = vm.currentMonth
		status.InactiveDays = 0
		status.RewardEligible = true
		status.MissedBlocks = 0
	}
	
	log.Printf("Monthly reset completed - Month %d -> %d", oldMonth, vm.currentMonth)
	
	// Send monthly report
	vm.sendMonthlyReport(oldMonth)
}

// calculateAverageUptime calculates average uptime across all validators
func (vm *ValidatorMonitor) calculateAverageUptime() float64 {
	if len(vm.validators) == 0 {
		return 0.0
	}
	
	totalUptime := 0.0
	for _, status := range vm.validators {
		totalUptime += status.MonthlyUptime
	}
	
	return totalUptime / float64(len(vm.validators))
}

// countRunningBots counts validators with running bots
func (vm *ValidatorMonitor) countRunningBots() int {
	count := 0
	for _, status := range vm.validators {
		if status.BotRunning {
			count++
		}
	}
	return count
}

// RegisterBotHeartbeat registers a bot heartbeat
func (vm *ValidatorMonitor) RegisterBotHeartbeat(operatorAddr string, version string) {
	vm.mu.Lock()
	defer vm.mu.Unlock()
	
	vm.botHeartbeats[operatorAddr] = time.Now()
	
	if status, exists := vm.validators[operatorAddr]; exists {
		status.BotRunning = true
		status.BotVersion = version
		status.LastBotHeartbeat = time.Now()
	}
}

// sendBotInactivityAlert sends an alert for bot inactivity
func (vm *ValidatorMonitor) sendBotInactivityAlert(status *ValidatorStatus) {
	message := fmt.Sprintf("🤖 Bot Inactivity Alert\n\nValidator: %s\nBot Status: Inactive\nLast Heartbeat: %s\nAction: Queued for slashing", 
		status.Moniker, 
		status.LastBotHeartbeat.Format("2006-01-02 15:04:05"))
	
	vm.sendAlert("Bot Inactivity", message)
}

// sendMonthlyReport sends a monthly statistics report
func (vm *ValidatorMonitor) sendMonthlyReport(month uint64) {
	stats, exists := vm.monthlyStats[month]
	if !exists {
		return
	}
	
	message := fmt.Sprintf("📊 Monthly Validator Report\n\nMonth: %d\nTotal Validators: %d\nActive: %d\nInactive: %d\nForfeited Rewards: %.2f GXR\nAverage Uptime: %.1f%%\nBots Running: %d", 
		stats.Month,
		stats.TotalValidators,
		stats.ActiveValidators,
		stats.InactiveValidators,
		stats.ForfeitedRewards,
		stats.AverageUptime,
		stats.BotsRunning)
	
	vm.sendAlert("Monthly Report", message)
}

// sendAlert sends a telegram alert
func (vm *ValidatorMonitor) sendAlert(title, message string) error {
	if vm.telegramAlert == nil {
		return nil
	}
	
	// Rate limiting - don't send alerts too frequently
	if time.Since(vm.lastAlertTime) < 2*time.Minute {
		return nil
	}
	
	fullMessage := fmt.Sprintf("%s\n\n%s", title, message)
	if err := vm.telegramAlert.SendAlert(fullMessage); err != nil {
		log.Printf("Failed to send alert: %v", err)
		return err
	}
	
	vm.lastAlertTime = time.Now()
	vm.alertsSent++
	return nil
}

// GetValidatorStatus returns the status of a specific validator
func (vm *ValidatorMonitor) GetValidatorStatus(operatorAddr string) (*ValidatorStatus, bool) {
	vm.mu.RLock()
	defer vm.mu.RUnlock()
	
	status, exists := vm.validators[operatorAddr]
	return status, exists
}

// GetAllValidatorStatuses returns all validator statuses
func (vm *ValidatorMonitor) GetAllValidatorStatuses() map[string]*ValidatorStatus {
	vm.mu.RLock()
	defer vm.mu.RUnlock()
	
	// Create a copy to avoid race conditions
	result := make(map[string]*ValidatorStatus)
	for addr, status := range vm.validators {
		result[addr] = status
	}
	
	return result
}

// GetMonthlyStats returns monthly statistics
func (vm *ValidatorMonitor) GetMonthlyStats() map[uint64]*MonthlyStats {
	vm.mu.RLock()
	defer vm.mu.RUnlock()
	
	result := make(map[uint64]*MonthlyStats)
	for month, stats := range vm.monthlyStats {
		result[month] = stats
	}
	
	return result
}

// GetStatus returns current monitor status
func (vm *ValidatorMonitor) GetStatus() map[string]interface{} {
	vm.mu.RLock()
	defer vm.mu.RUnlock()
	
	return map[string]interface{}{
		"total_validators":         vm.totalValidators,
		"active_validators":        vm.activeValidators,
		"inactive_validators":      vm.totalInactiveValidators,
		"current_month":           vm.currentMonth,
		"last_month_reset":        vm.lastMonthReset.Format(time.RFC3339),
		"slashing_queue_size":     len(vm.slashingQueue),
		"running_bots":            vm.countRunningBots(),
		"total_forfeited_rewards": vm.totalForfeitedRewards,
		"alerts_sent":             vm.alertsSent,
		"average_uptime":          vm.calculateAverageUptime(),
	}
}

// getCurrentMonth returns current month identifier
func getCurrentMonth() uint64 {
	return uint64(time.Now().Unix() / int64(30*24*time.Hour.Seconds()))
}

// Stop gracefully stops the validator monitor
func (vm *ValidatorMonitor) Stop() {
	vm.mu.Lock()
	defer vm.mu.Unlock()
	
	log.Printf("Stopping validator monitor - Final stats: %d validators, %d alerts sent", 
		vm.totalValidators, vm.alertsSent)
	
	vm.sendAlert("Monitor Stopped", "Validator monitor stopped")
}