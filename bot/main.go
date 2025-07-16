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
)

const (
	// Bot version
	Version = "1.0.0"
	
	// Bot configuration
	DefaultConfigPath = "./config/bot.yaml"
	DefaultLogLevel   = "info"
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
func NewGXRBot(config *BotConfig) *GXRBot {
	ctx, cancel := context.WithCancel(context.Background())
	
	return &GXRBot{
		config: config,
		ctx:    ctx,
		cancel: cancel,
		wg:     &sync.WaitGroup{},
	}
}

// Initialize initializes all bot components
func (b *GXRBot) Initialize() error {
	log.Printf("Initializing GXR Bot v%s", Version)
	
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
	
	// Initialize Telegram Alert
	if b.config.TelegramToken != "" {
		b.telegramAlert = NewTelegramAlert(b.config)
		if err := b.telegramAlert.Initialize(); err != nil {
			log.Printf("⚠️  Telegram alerts disabled: %v", err)
		} else {
			log.Println("✅ Telegram Alert initialized")
		}
	}
	
	log.Println("🚀 GXR Bot initialization complete")
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
	
	// Send startup notification
	if b.telegramAlert != nil {
		b.telegramAlert.SendAlert("🚀 GXR Bot started successfully")
	}
	
	log.Println("✅ All GXR Bot services started")
	return nil
}

// Stop gracefully stops all bot services
func (b *GXRBot) Stop() {
	log.Println("Stopping GXR Bot services...")
	
	// Cancel context to signal all goroutines to stop
	b.cancel()
	
	// Wait for all goroutines to finish
	b.wg.Wait()
	
	// Send shutdown notification
	if b.telegramAlert != nil {
		b.telegramAlert.SendAlert("🛑 GXR Bot stopped")
	}
	
	log.Println("✅ GXR Bot stopped gracefully")
}

// DefaultConfig returns the default bot configuration
func DefaultConfig() *BotConfig {
	return &BotConfig{
		ChainRPC:      "tcp://localhost:26657",
		ChainGRPC:     "localhost:9090",
		ChainID:       "gxr-1",
		LogLevel:      DefaultLogLevel,
		CheckInterval: 30 * time.Second,
		IBCEnabled:    true,
		IBCChannels:   []string{"channel-0", "channel-1"},
		MaxSwapDaily:  "10000ugen", // 10,000 GXR daily limit
		SwapCooldown:  30 * time.Minute,
		PriceLimit:    "10000000", // $10 limit for emergency mode
		EmergencyMode: false,
	}
}

// Main CLI command
func main() {
	var configPath string
	
	rootCmd := &cobra.Command{
		Use:   "gxr-bot",
		Short: "GXR Validator Bot",
		Long: `GXR Validator Bot provides automated services for GXR blockchain:
- IBC Relayer for cross-chain synchronization
- Reward Distribution automation
- DEX Pool auto-refill and rebalancing
- Telegram alerts for monitoring`,
		Version: Version,
		RunE: func(cmd *cobra.Command, args []string) error {
			// Load configuration
			config := DefaultConfig()
			if configPath != "" {
				// TODO: Load from file
				log.Printf("Loading config from: %s", configPath)
			}
			
			// Create and initialize bot
			bot := NewGXRBot(config)
			if err := bot.Initialize(); err != nil {
				return fmt.Errorf("failed to initialize bot: %w", err)
			}
			
			// Start bot services
			if err := bot.Start(); err != nil {
				return fmt.Errorf("failed to start bot: %w", err)
			}
			
			// Wait for interrupt signal
			sigChan := make(chan os.Signal, 1)
			signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
			
			log.Println("GXR Bot is running. Press Ctrl+C to stop.")
			<-sigChan
			
			// Graceful shutdown
			bot.Stop()
			return nil
		},
	}
	
	rootCmd.Flags().StringVarP(&configPath, "config", "c", DefaultConfigPath, "Path to configuration file")
	
	if err := rootCmd.Execute(); err != nil {
		log.Fatal(err)
	}
}