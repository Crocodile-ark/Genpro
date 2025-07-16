package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v2"
)

const (
	// Bot version
	Version = "2.0.0"
	
	// Bot configuration
	DefaultConfigPath = "./config/bot.yaml"
	DefaultLogLevel   = "info"
	
	// Default values
	DefaultCheckInterval = 5 * time.Minute
	DefaultSwapCooldown  = 1 * time.Hour
	DefaultPriceLimit    = "5.0"
	DefaultMaxSwapDaily  = "10000"
	
	// Health check interval
	HealthCheckInterval = 30 * time.Second
	
	// Shutdown timeout
	ShutdownTimeout = 30 * time.Second
)

// BotConfig represents the enhanced bot configuration
type BotConfig struct {
	// Chain connection settings
	ChainRPC     string `yaml:"chain_rpc"`
	ChainGRPC    string `yaml:"chain_grpc"`
	ChainID      string `yaml:"chain_id"`
	
	// Validator settings
	ValidatorAddress string `yaml:"validator_address"`
	ValidatorName    string `yaml:"validator_name"`
	ValidatorMnemonic string `yaml:"validator_mnemonic"`
	
	// Bot settings
	LogLevel     string        `yaml:"log_level"`
	CheckInterval time.Duration `yaml:"check_interval"`
	
	// Rebalancing settings
	SwapCooldown  time.Duration `yaml:"swap_cooldown"`
	PriceLimit    string        `yaml:"price_limit"`
	MaxSwapDaily  string        `yaml:"max_swap_daily"`
	
	// IBC settings
	IBCEnabled   bool     `yaml:"ibc_enabled"`
	IBCChannels  []string `yaml:"ibc_channels"`
	
	// DEX settings
	DEXEnabled bool     `yaml:"dex_enabled"`
	DEXPools   []string `yaml:"dex_pools"`
	
	// Telegram settings
	TelegramEnabled bool   `yaml:"telegram_enabled"`
	TelegramToken   string `yaml:"telegram_token"`
	TelegramChatID  string `yaml:"telegram_chat_id"`
	
	// Enhanced monitoring
	MonitoringEnabled     bool `yaml:"monitoring_enabled"`
	HealthCheckEnabled    bool `yaml:"health_check_enabled"`
	MetricsEnabled        bool `yaml:"metrics_enabled"`
	
	// Advanced settings
	RetryAttempts     int           `yaml:"retry_attempts"`
	RetryDelay        time.Duration `yaml:"retry_delay"`
	MaxConcurrentOps  int           `yaml:"max_concurrent_ops"`
	EnableProfiling   bool          `yaml:"enable_profiling"`
}

// BotService represents the main bot service
type BotService struct {
	config    *BotConfig
	clientCtx client.Context
	cdc       codec.Codec
	mu        sync.RWMutex
	
	// Core components
	rebalancer       *Rebalancer
	validatorMonitor *ValidatorMonitor
	ibcRelayer       *IBCRelayer
	dexManager       *DEXManager
	rewardDistributor *RewardDistributor
	telegramAlert    *TelegramAlert
	
	// State management
	running          bool
	startTime        time.Time
	lastHealthCheck  time.Time
	errorCount       int64
	successCount     int64
	
	// Health monitoring
	healthStatus     map[string]bool
	lastErrors       []ErrorRecord
	
	// Shutdown handling
	shutdownChan     chan struct{}
	shutdownComplete chan struct{}
}

// ErrorRecord represents an error record
type ErrorRecord struct {
	Timestamp time.Time
	Component string
	Error     string
}

// NewBotService creates a new enhanced bot service
func NewBotService(config *BotConfig) (*BotService, error) {
	bs := &BotService{
		config:           config,
		healthStatus:     make(map[string]bool),
		lastErrors:       make([]ErrorRecord, 0),
		shutdownChan:     make(chan struct{}),
		shutdownComplete: make(chan struct{}),
	}
	
	// Initialize components
	if err := bs.initializeComponents(); err != nil {
		return nil, fmt.Errorf("failed to initialize components: %w", err)
	}
	
	return bs, nil
}

// initializeComponents initializes all bot components
func (bs *BotService) initializeComponents() error {
	log.Printf("Initializing bot components...")
	
	// Initialize telegram alert first
	if bs.config.TelegramEnabled {
		bs.telegramAlert = NewTelegramAlert(bs.config)
		if err := bs.telegramAlert.TestConnection(); err != nil {
			log.Printf("Warning: Telegram connection failed: %v", err)
		} else {
			bs.telegramAlert.SendTestAlert()
		}
	}
	
	// Initialize chain client context
	if err := bs.initializeChainClient(); err != nil {
		return fmt.Errorf("failed to initialize chain client: %w", err)
	}
	
	// Initialize rebalancer
	bs.rebalancer = NewRebalancer(bs.config)
	bs.healthStatus["rebalancer"] = true
	
	// Initialize validator monitor
	bs.validatorMonitor = NewValidatorMonitor(bs.config, bs.clientCtx, bs.cdc)
	bs.healthStatus["validator_monitor"] = true
	
	// Initialize IBC relayer if enabled
	if bs.config.IBCEnabled {
		bs.ibcRelayer = NewIBCRelayer(bs.config)
		bs.healthStatus["ibc_relayer"] = true
	}
	
	// Initialize DEX manager if enabled
	if bs.config.DEXEnabled {
		bs.dexManager = NewDEXManager(bs.config)
		bs.healthStatus["dex_manager"] = true
	}
	
	// Initialize reward distributor
	bs.rewardDistributor = NewRewardDistributor(bs.config)
	bs.healthStatus["reward_distributor"] = true
	
	log.Printf("All components initialized successfully")
	return nil
}

// initializeChainClient initializes the chain client
func (bs *BotService) initializeChainClient() error {
	log.Printf("Initializing chain client...")
	log.Printf("Chain ID: %s", bs.config.ChainID)
	log.Printf("Chain RPC: %s", bs.config.ChainRPC)
	log.Printf("Chain gRPC: %s", bs.config.ChainGRPC)
	
	// In a real implementation, this would create proper Cosmos SDK client
	// For now, we'll simulate the initialization
	time.Sleep(1 * time.Second)
	
	log.Printf("Chain client initialized successfully")
	return nil
}

// Start starts the bot service
func (bs *BotService) Start(ctx context.Context) error {
	bs.mu.Lock()
	bs.running = true
	bs.startTime = time.Now()
	bs.mu.Unlock()
	
	log.Printf("Starting GXR Bot Service v%s", Version)
	
	// Send startup notification
	if bs.telegramAlert != nil {
		bs.telegramAlert.SendBotAlert("GXR Bot", "started", "Bot service started successfully")
	}
	
	// Start all components
	if err := bs.startComponents(ctx); err != nil {
		return fmt.Errorf("failed to start components: %w", err)
	}
	
	// Start health monitoring
	if bs.config.HealthCheckEnabled {
		go bs.healthMonitor(ctx)
	}
	
	// Start heartbeat for validator monitoring
	go bs.sendHeartbeat(ctx)
	
	log.Printf("Bot service started successfully - All components running")
	return nil
}

// startComponents starts all bot components
func (bs *BotService) startComponents(ctx context.Context) error {
	var wg sync.WaitGroup
	errors := make(chan error, 10)
	
	// Start rebalancer
	wg.Add(1)
	go func() {
		defer wg.Done()
		if err := bs.rebalancer.Start(ctx); err != nil {
			errors <- fmt.Errorf("rebalancer failed: %w", err)
		}
	}()
	
	// Start validator monitor
	wg.Add(1)
	go func() {
		defer wg.Done()
		if err := bs.validatorMonitor.Start(ctx); err != nil {
			errors <- fmt.Errorf("validator monitor failed: %w", err)
		}
	}()
	
	// Start IBC relayer if enabled
	if bs.ibcRelayer != nil {
		wg.Add(1)
		go func() {
			defer wg.Done()
			if err := bs.ibcRelayer.Start(ctx); err != nil {
				errors <- fmt.Errorf("IBC relayer failed: %w", err)
			}
		}()
	}
	
	// Start DEX manager if enabled
	if bs.dexManager != nil {
		wg.Add(1)
		go func() {
			defer wg.Done()
			if err := bs.dexManager.Start(ctx); err != nil {
				errors <- fmt.Errorf("DEX manager failed: %w", err)
			}
		}()
	}
	
	// Start reward distributor
	wg.Add(1)
	go func() {
		defer wg.Done()
		if err := bs.rewardDistributor.Start(ctx); err != nil {
			errors <- fmt.Errorf("reward distributor failed: %w", err)
		}
	}()
	
	// Check for startup errors
	go func() {
		wg.Wait()
		close(errors)
	}()
	
	// Collect any startup errors
	for err := range errors {
		log.Printf("Component startup error: %v", err)
		bs.recordError("startup", err.Error())
		if bs.telegramAlert != nil {
			bs.telegramAlert.SendBotAlert("Startup", "error", err.Error())
		}
	}
	
	return nil
}

// healthMonitor monitors the health of all components
func (bs *BotService) healthMonitor(ctx context.Context) {
	ticker := time.NewTicker(HealthCheckInterval)
	defer ticker.Stop()
	
	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			bs.performHealthCheck()
		}
	}
}

// performHealthCheck checks the health of all components
func (bs *BotService) performHealthCheck() {
	bs.mu.Lock()
	defer bs.mu.Unlock()
	
	bs.lastHealthCheck = time.Now()
	
	// Check rebalancer health
	if bs.rebalancer != nil {
		status := bs.rebalancer.GetStatus()
		bs.healthStatus["rebalancer"] = status["state"] != "error"
	}
	
	// Check validator monitor health
	if bs.validatorMonitor != nil {
		status := bs.validatorMonitor.GetStatus()
		bs.healthStatus["validator_monitor"] = status["total_validators"].(int) > 0
	}
	
	// Check IBC relayer health
	if bs.ibcRelayer != nil {
		status := bs.ibcRelayer.GetStatus()
		bs.healthStatus["ibc_relayer"] = status["connected"].(bool)
	}
	
	// Check DEX manager health
	if bs.dexManager != nil {
		status := bs.dexManager.GetStatus()
		bs.healthStatus["dex_manager"] = status["pools_active"].(int) > 0
	}
	
	// Check reward distributor health
	if bs.rewardDistributor != nil {
		status := bs.rewardDistributor.GetStatus()
		bs.healthStatus["reward_distributor"] = status["connected"].(bool)
	}
	
	// Check telegram alert health
	if bs.telegramAlert != nil {
		bs.healthStatus["telegram_alert"] = bs.telegramAlert.IsRunning()
	}
	
	// Count unhealthy components
	unhealthyCount := 0
	for component, healthy := range bs.healthStatus {
		if !healthy {
			unhealthyCount++
			log.Printf("Health check failed for component: %s", component)
		}
	}
	
	// Send alert if too many components are unhealthy
	if unhealthyCount > 2 && bs.telegramAlert != nil {
		bs.telegramAlert.SendEmergencyAlert("Multiple Component Failures", 
			fmt.Sprintf("%d components are unhealthy", unhealthyCount), 
			map[string]interface{}{"unhealthy_count": unhealthyCount})
	}
}

// sendHeartbeat sends periodic heartbeat to validator monitor
func (bs *BotService) sendHeartbeat(ctx context.Context) {
	ticker := time.NewTicker(1 * time.Minute)
	defer ticker.Stop()
	
	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			if bs.validatorMonitor != nil && bs.config.ValidatorAddress != "" {
				bs.validatorMonitor.RegisterBotHeartbeat(bs.config.ValidatorAddress, Version)
			}
		}
	}
}

// recordError records an error in the bot service
func (bs *BotService) recordError(component, errorMsg string) {
	bs.errorCount++
	
	record := ErrorRecord{
		Timestamp: time.Now(),
		Component: component,
		Error:     errorMsg,
	}
	
	bs.lastErrors = append(bs.lastErrors, record)
	
	// Keep only last 50 errors
	if len(bs.lastErrors) > 50 {
		bs.lastErrors = bs.lastErrors[1:]
	}
}

// GetStatus returns the current status of the bot service
func (bs *BotService) GetStatus() map[string]interface{} {
	bs.mu.RLock()
	defer bs.mu.RUnlock()
	
	status := map[string]interface{}{
		"version":           Version,
		"running":           bs.running,
		"start_time":        bs.startTime.Format(time.RFC3339),
		"uptime":            time.Since(bs.startTime).String(),
		"last_health_check": bs.lastHealthCheck.Format(time.RFC3339),
		"error_count":       bs.errorCount,
		"success_count":     bs.successCount,
		"health_status":     bs.healthStatus,
		"config": map[string]interface{}{
			"chain_id":           bs.config.ChainID,
			"validator_address":  bs.config.ValidatorAddress,
			"validator_name":     bs.config.ValidatorName,
			"telegram_enabled":   bs.config.TelegramEnabled,
			"ibc_enabled":        bs.config.IBCEnabled,
			"dex_enabled":        bs.config.DEXEnabled,
			"monitoring_enabled": bs.config.MonitoringEnabled,
		},
	}
	
	// Add component statuses
	componentStatuses := make(map[string]interface{})
	
	if bs.rebalancer != nil {
		componentStatuses["rebalancer"] = bs.rebalancer.GetStatus()
	}
	
	if bs.validatorMonitor != nil {
		componentStatuses["validator_monitor"] = bs.validatorMonitor.GetStatus()
	}
	
	if bs.ibcRelayer != nil {
		componentStatuses["ibc_relayer"] = bs.ibcRelayer.GetStatus()
	}
	
	if bs.dexManager != nil {
		componentStatuses["dex_manager"] = bs.dexManager.GetStatus()
	}
	
	if bs.rewardDistributor != nil {
		componentStatuses["reward_distributor"] = bs.rewardDistributor.GetStatus()
	}
	
	if bs.telegramAlert != nil {
		componentStatuses["telegram_alert"] = bs.telegramAlert.GetStatistics()
	}
	
	status["components"] = componentStatuses
	
	return status
}

// Stop gracefully stops the bot service
func (bs *BotService) Stop() error {
	bs.mu.Lock()
	if !bs.running {
		bs.mu.Unlock()
		return nil
	}
	bs.running = false
	bs.mu.Unlock()
	
	log.Printf("Stopping bot service...")
	
	// Signal shutdown
	close(bs.shutdownChan)
	
	// Stop all components
	if bs.rebalancer != nil {
		bs.rebalancer.Stop()
	}
	
	if bs.validatorMonitor != nil {
		bs.validatorMonitor.Stop()
	}
	
	if bs.ibcRelayer != nil {
		bs.ibcRelayer.Stop()
	}
	
	if bs.dexManager != nil {
		bs.dexManager.Stop()
	}
	
	if bs.rewardDistributor != nil {
		bs.rewardDistributor.Stop()
	}
	
	// Send shutdown notification
	if bs.telegramAlert != nil {
		bs.telegramAlert.SendBotAlert("GXR Bot", "stopped", "Bot service stopped")
		bs.telegramAlert.Stop()
	}
	
	// Wait for graceful shutdown or timeout
	select {
	case <-bs.shutdownComplete:
		log.Printf("Bot service stopped gracefully")
	case <-time.After(ShutdownTimeout):
		log.Printf("Bot service shutdown timeout")
	}
	
	return nil
}

// LoadConfig loads the bot configuration
func LoadConfig(configPath string) (*BotConfig, error) {
	if configPath == "" {
		configPath = DefaultConfigPath
	}
	
	// Set default values
	config := &BotConfig{
		LogLevel:      DefaultLogLevel,
		CheckInterval: DefaultCheckInterval,
		SwapCooldown:  DefaultSwapCooldown,
		PriceLimit:    DefaultPriceLimit,
		MaxSwapDaily:  DefaultMaxSwapDaily,
		RetryAttempts: 3,
		RetryDelay:    5 * time.Second,
		MaxConcurrentOps: 10,
		HealthCheckEnabled: true,
		MonitoringEnabled: true,
	}
	
	// Try to load from file
	if _, err := os.Stat(configPath); err == nil {
		data, err := os.ReadFile(configPath)
		if err != nil {
			return nil, fmt.Errorf("failed to read config file: %w", err)
		}
		
		if err := yaml.Unmarshal(data, config); err != nil {
			return nil, fmt.Errorf("failed to parse config file: %w", err)
		}
		
		log.Printf("Configuration loaded from: %s", configPath)
	} else {
		log.Printf("Config file not found, using defaults: %s", configPath)
	}
	
	// Validate configuration
	if err := ValidateConfig(config); err != nil {
		return nil, fmt.Errorf("configuration validation failed: %w", err)
	}
	
	return config, nil
}

// ValidateConfig validates the bot configuration
func ValidateConfig(config *BotConfig) error {
	if config.ChainID == "" {
		return fmt.Errorf("chain_id is required")
	}
	
	if config.ChainRPC == "" {
		return fmt.Errorf("chain_rpc is required")
	}
	
	if config.ChainGRPC == "" {
		return fmt.Errorf("chain_grpc is required")
	}
	
	if config.ValidatorAddress == "" {
		return fmt.Errorf("validator_address is required")
	}
	
	if config.TelegramEnabled {
		if config.TelegramToken == "" {
			return fmt.Errorf("telegram_token is required when telegram is enabled")
		}
		if config.TelegramChatID == "" {
			return fmt.Errorf("telegram_chat_id is required when telegram is enabled")
		}
	}
	
	if config.CheckInterval < 1*time.Minute {
		return fmt.Errorf("check_interval must be at least 1 minute")
	}
	
	if config.SwapCooldown < 1*time.Hour {
		return fmt.Errorf("swap_cooldown must be at least 1 hour")
	}
	
	if config.RetryAttempts < 1 || config.RetryAttempts > 10 {
		return fmt.Errorf("retry_attempts must be between 1 and 10")
	}
	
	if config.MaxConcurrentOps < 1 || config.MaxConcurrentOps > 100 {
		return fmt.Errorf("max_concurrent_ops must be between 1 and 100")
	}
	
	return nil
}

// CreateRootCmd creates the root command
func CreateRootCmd() *cobra.Command {
	var configPath string
	
	rootCmd := &cobra.Command{
		Use:   "gxr-bot",
		Short: "GXR Blockchain Bot Service",
		Long:  "Enhanced GXR blockchain bot with validator monitoring, rebalancing, and alert systems",
		RunE: func(cmd *cobra.Command, args []string) error {
			return runBot(configPath)
		},
	}
	
	rootCmd.PersistentFlags().StringVar(&configPath, "config", DefaultConfigPath, "Path to configuration file")
	
	// Add subcommands
	rootCmd.AddCommand(createStatusCmd())
	rootCmd.AddCommand(createTestCmd())
	rootCmd.AddCommand(createVersionCmd())
	
	return rootCmd
}

// runBot runs the main bot service
func runBot(configPath string) error {
	// Load configuration
	config, err := LoadConfig(configPath)
	if err != nil {
		return fmt.Errorf("failed to load configuration: %w", err)
	}
	
	// Create bot service
	botService, err := NewBotService(config)
	if err != nil {
		return fmt.Errorf("failed to create bot service: %w", err)
	}
	
	// Setup signal handling
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	
	// Start bot service
	go func() {
		if err := botService.Start(ctx); err != nil {
			log.Printf("Bot service error: %v", err)
			cancel()
		}
	}()
	
	// Wait for shutdown signal
	<-sigChan
	log.Printf("Received shutdown signal")
	
	// Graceful shutdown
	cancel()
	return botService.Stop()
}

// createStatusCmd creates the status command
func createStatusCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "status",
		Short: "Show bot status",
		RunE: func(cmd *cobra.Command, args []string) error {
			// In a real implementation, this would connect to a running bot instance
			fmt.Println("Bot Status: This would show the current bot status")
			return nil
		},
	}
}

// createTestCmd creates the test command
func createTestCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "test",
		Short: "Test bot configuration",
		RunE: func(cmd *cobra.Command, args []string) error {
			config, err := LoadConfig(DefaultConfigPath)
			if err != nil {
				return fmt.Errorf("configuration test failed: %w", err)
			}
			
			fmt.Printf("Configuration test passed for chain: %s\n", config.ChainID)
			return nil
		},
	}
}

// createVersionCmd creates the version command
func createVersionCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "version",
		Short: "Show bot version",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Printf("GXR Bot version %s\n", Version)
		},
	}
}

// main is the entry point
func main() {
	rootCmd := CreateRootCmd()
	
	if err := rootCmd.Execute(); err != nil {
		log.Fatalf("Command execution failed: %v", err)
	}
}