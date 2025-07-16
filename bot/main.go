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

	"github.com/spf13/cobra"
	"gopkg.in/yaml.v2"
)

const (
	// Bot version
	Version = "1.0.0"
	
	// Bot configuration
	DefaultConfigPath = "./config/bot.yaml"
	DefaultLogLevel   = "info"
	
	// Default values
	DefaultCheckInterval = 5 * time.Minute
	DefaultSwapCooldown  = 30 * time.Minute
	DefaultPriceLimit    = "10.0"
	DefaultMaxSwapDaily  = "10000"
)

// BotConfig represents the bot configuration
type BotConfig struct {
	// Chain connection settings
	ChainRPC     string `yaml:"chain_rpc"`
	ChainGRPC    string `yaml:"chain_grpc"`
	ChainID      string `yaml:"chain_id"`
	
	// Bot settings
	LogLevel     string `yaml:"log_level"`
	CheckInterval time.Duration `yaml:"check_interval"`
	
	// IBC settings
	IBCEnabled   bool     `yaml:"ibc_enabled"`
	IBCChannels  []string `yaml:"ibc_channels"`
	
	// DEX settings
	MaxSwapDaily string `yaml:"max_swap_daily"`
	SwapCooldown time.Duration `yaml:"swap_cooldown"`
	PriceLimit   string `yaml:"price_limit"`
	
	// Telegram settings
	TelegramToken  string `yaml:"telegram_token"`
	TelegramChatID string `yaml:"telegram_chat_id"`
	
	// Safety settings
	EmergencyMode bool `yaml:"emergency_mode"`
}

// DefaultConfig returns a default configuration
func DefaultConfig() *BotConfig {
	return &BotConfig{
		ChainRPC:      "tcp://localhost:26657",
		ChainGRPC:     "localhost:9090",
		ChainID:       "gxr-1",
		LogLevel:      DefaultLogLevel,
		CheckInterval: DefaultCheckInterval,
		IBCEnabled:    false,
		IBCChannels:   []string{},
		MaxSwapDaily:  DefaultMaxSwapDaily,
		SwapCooldown:  DefaultSwapCooldown,
		PriceLimit:    DefaultPriceLimit,
		EmergencyMode: false,
	}
}

// Validate validates the bot configuration
func (c *BotConfig) Validate() error {
	if c.ChainRPC == "" {
		return fmt.Errorf("chain_rpc is required")
	}
	if c.ChainGRPC == "" {
		return fmt.Errorf("chain_grpc is required")
	}
	if c.ChainID == "" {
		return fmt.Errorf("chain_id is required")
	}
	if c.CheckInterval <= 0 {
		return fmt.Errorf("check_interval must be positive")
	}
	if c.SwapCooldown <= 0 {
		return fmt.Errorf("swap_cooldown must be positive")
	}
	
	// Validate IBC settings if enabled
	if c.IBCEnabled && len(c.IBCChannels) == 0 {
		return fmt.Errorf("ibc_channels must be specified when IBC is enabled")
	}
	
	// Validate Telegram settings if provided
	if c.TelegramToken != "" && c.TelegramChatID == "" {
		return fmt.Errorf("telegram_chat_id is required when telegram_token is provided")
	}
	
	return nil
}

// GXRBot represents the main bot instance
type GXRBot struct {
	config   *BotConfig
	ctx      context.Context
	cancel   context.CancelFunc
	wg       *sync.WaitGroup
	
	// Bot components
	ibcRelayer     *IBCRelayer
	rewardDistributor *RewardDistributor
	dexManager     *DEXManager
	rebalancer     *Rebalancer
	telegramAlert  *TelegramAlert
}

// NewGXRBot creates a new bot instance
func NewGXRBot(config *BotConfig) (*GXRBot, error) {
	if err := config.Validate(); err != nil {
		return nil, fmt.Errorf("invalid configuration: %w", err)
	}
	
	ctx, cancel := context.WithCancel(context.Background())
	
	return &GXRBot{
		config: config,
		ctx:    ctx,
		cancel: cancel,
		wg:     &sync.WaitGroup{},
	}, nil
}

// Initialize initializes all bot components
func (b *GXRBot) Initialize() error {
	log.Printf("Initializing GXR Bot v%s", Version)
	
	// Initialize Telegram Alert first (so other components can use it)
	if b.config.TelegramToken != "" {
		b.telegramAlert = NewTelegramAlert(b.config)
		if err := b.telegramAlert.Initialize(); err != nil {
			log.Printf("Warning: Failed to initialize Telegram alerts: %v", err)
			b.telegramAlert = nil // Disable if initialization fails
		} else {
			log.Println("✅ Telegram Alert initialized")
		}
	}
	
	// Initialize IBC Relayer
	if b.config.IBCEnabled {
		b.ibcRelayer = NewIBCRelayer(b.config)
		if err := b.ibcRelayer.Initialize(); err != nil {
			return fmt.Errorf("failed to initialize IBC relayer: %w", err)
		}
		log.Println("✅ IBC Relayer initialized")
	}
	
	// Initialize Reward Distributor
	b.rewardDistributor = NewRewardDistributor(b.config)
	if err := b.rewardDistributor.Initialize(); err != nil {
		return fmt.Errorf("failed to initialize reward distributor: %w", err)
	}
	log.Println("✅ Reward Distributor initialized")
	
	// Initialize DEX Manager
	b.dexManager = NewDEXManager(b.config)
	if err := b.dexManager.Initialize(); err != nil {
		return fmt.Errorf("failed to initialize DEX manager: %w", err)
	}
	log.Println("✅ DEX Manager initialized")
	
	// Initialize Rebalancer
	b.rebalancer = NewRebalancer(b.config)
	if err := b.rebalancer.Initialize(); err != nil {
		return fmt.Errorf("failed to initialize rebalancer: %w", err)
	}
	log.Println("✅ Rebalancer initialized")
	
	log.Println("🚀 All bot components initialized successfully")
	return nil
}

// Start starts all bot services
func (b *GXRBot) Start() error {
	log.Println("Starting GXR Bot services...")
	
	// Start IBC Relayer
	if b.ibcRelayer != nil {
		b.wg.Add(1)
		go func() {
			defer b.wg.Done()
			if err := b.ibcRelayer.Start(b.ctx); err != nil {
				log.Printf("IBC Relayer error: %v", err)
			}
		}()
	}
	
	// Start Reward Distributor
	b.wg.Add(1)
	go func() {
		defer b.wg.Done()
		if err := b.rewardDistributor.Start(b.ctx); err != nil {
			log.Printf("Reward Distributor error: %v", err)
		}
	}()
	
	// Start DEX Manager
	b.wg.Add(1)
	go func() {
		defer b.wg.Done()
		if err := b.dexManager.Start(b.ctx); err != nil {
			log.Printf("DEX Manager error: %v", err)
		}
	}()
	
	// Start Rebalancer
	b.wg.Add(1)
	go func() {
		defer b.wg.Done()
		if err := b.rebalancer.Start(b.ctx); err != nil {
			log.Printf("Rebalancer error: %v", err)
		}
	}()
	
	// Send startup alert
	if b.telegramAlert != nil {
		b.telegramAlert.SendAlert("🚀 GXR Bot started successfully\n\nAll services are now running.")
	}
	
	log.Println("🎉 All GXR Bot services started successfully")
	return nil
}

// Stop gracefully stops all bot services
func (b *GXRBot) Stop() {
	log.Println("Stopping GXR Bot...")
	
	// Send shutdown alert
	if b.telegramAlert != nil {
		b.telegramAlert.SendAlert("🛑 GXR Bot shutting down\n\nAll services will be stopped.")
	}
	
	// Cancel context to stop all goroutines
	b.cancel()
	
	// Wait for all goroutines to finish
	b.wg.Wait()
	
	log.Println("👋 GXR Bot stopped successfully")
}

// GetStatus returns the status of all bot components
func (b *GXRBot) GetStatus() map[string]interface{} {
	status := map[string]interface{}{
		"version": Version,
		"config":  b.config,
	}
	
	if b.ibcRelayer != nil {
		status["ibc_relayer"] = b.ibcRelayer.GetStatus()
	}
	
	if b.rewardDistributor != nil {
		status["reward_distributor"] = b.rewardDistributor.GetStatus()
	}
	
	if b.dexManager != nil {
		status["dex_manager"] = b.dexManager.GetStatus()
	}
	
	if b.rebalancer != nil {
		status["rebalancer"] = b.rebalancer.GetStatus()
	}
	
	if b.telegramAlert != nil {
		status["telegram_alert"] = b.telegramAlert.GetStatus()
	}
	
	return status
}

// loadConfig loads configuration from file
func loadConfig(path string) (*BotConfig, error) {
	if path == "" {
		path = DefaultConfigPath
	}
	
	// Check if file exists
	if _, err := os.Stat(path); os.IsNotExist(err) {
		log.Printf("Config file not found at %s, using defaults", path)
		return DefaultConfig(), nil
	}
	
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}
	
	config := DefaultConfig()
	if err := yaml.Unmarshal(data, config); err != nil {
		return nil, fmt.Errorf("failed to parse config file: %w", err)
	}
	
	return config, nil
}

// main function and CLI setup
func main() {
	var configPath string
	
	rootCmd := &cobra.Command{
		Use:   "gxr-bot",
		Short: "GXR Blockchain Validator Bot",
		Long:  `GXR Bot handles automatic IBC relaying, reward distribution, DEX management, and rebalancing for GXR blockchain validators.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			// Load configuration
			config, err := loadConfig(configPath)
			if err != nil {
				return fmt.Errorf("failed to load configuration: %w", err)
			}
			
			// Create bot instance
			bot, err := NewGXRBot(config)
			if err != nil {
				return fmt.Errorf("failed to create bot: %w", err)
			}
			
			// Initialize bot
			if err := bot.Initialize(); err != nil {
				return fmt.Errorf("failed to initialize bot: %w", err)
			}
			
			// Start bot services
			if err := bot.Start(); err != nil {
				return fmt.Errorf("failed to start bot: %w", err)
			}
			
			// Set up signal handling for graceful shutdown
			c := make(chan os.Signal, 1)
			signal.Notify(c, os.Interrupt, syscall.SIGTERM)
			
			// Wait for shutdown signal
			<-c
			log.Println("Received shutdown signal")
			
			// Gracefully stop the bot
			bot.Stop()
			
			return nil
		},
	}
	
	rootCmd.Flags().StringVarP(&configPath, "config", "c", DefaultConfigPath, "Path to configuration file")
	
	if err := rootCmd.Execute(); err != nil {
		log.Fatal(err)
	}
}