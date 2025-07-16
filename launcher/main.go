package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/exec"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/spf13/cobra"
)

const (
	LauncherVersion = "1.0.0"
)

// LauncherConfig holds the launcher configuration
type LauncherConfig struct {
	ChainBinary    string
	BotBinary      string
	ChainHome      string
	ChainConfig    string
	BotConfig      string
	LogLevel       string
	AutoRestart    bool
	RestartDelay   time.Duration
}

// GXRLauncher manages both chain and bot processes
type GXRLauncher struct {
	config     *LauncherConfig
	ctx        context.Context
	cancel     context.CancelFunc
	wg         *sync.WaitGroup
	
	chainCmd   *exec.Cmd
	botCmd     *exec.Cmd
	
	chainRunning bool
	botRunning   bool
}

// NewGXRLauncher creates a new launcher instance
func NewGXRLauncher(config *LauncherConfig) *GXRLauncher {
	ctx, cancel := context.WithCancel(context.Background())
	
	return &GXRLauncher{
		config: config,
		ctx:    ctx,
		cancel: cancel,
		wg:     &sync.WaitGroup{},
	}
}

// Start starts both chain and bot processes
func (l *GXRLauncher) Start() error {
	log.Printf("🚀 Starting GXR Launcher v%s", LauncherVersion)
	
	// Start chain first
	if err := l.startChain(); err != nil {
		return fmt.Errorf("failed to start chain: %w", err)
	}
	
	// Wait a bit for chain to initialize
	log.Println("⏳ Waiting for chain initialization...")
	time.Sleep(10 * time.Second)
	
	// Start bot
	if err := l.startBot(); err != nil {
		log.Printf("⚠️  Failed to start bot: %v", err)
		log.Println("📄 Chain will continue running without bot")
	}
	
	log.Println("✅ GXR Launcher started successfully")
	log.Println("   📦 Chain: Running")
	if l.botRunning {
		log.Println("   🤖 Bot: Running")
	} else {
		log.Println("   🤖 Bot: Failed to start")
	}
	
	return nil
}

// startChain starts the GXR blockchain daemon
func (l *GXRLauncher) startChain() error {
	log.Println("🔗 Starting GXR Chain...")
	
	// Build chain command
	l.chainCmd = exec.CommandContext(l.ctx, l.config.ChainBinary, "start")
	
	// Set environment variables
	if l.config.ChainHome != "" {
		l.chainCmd.Env = append(os.Environ(), fmt.Sprintf("HOME=%s", l.config.ChainHome))
	}
	
	// Set up logging
	l.chainCmd.Stdout = &PrefixedWriter{prefix: "[CHAIN]", writer: os.Stdout}
	l.chainCmd.Stderr = &PrefixedWriter{prefix: "[CHAIN]", writer: os.Stderr}
	
	// Start chain process
	if err := l.chainCmd.Start(); err != nil {
		return fmt.Errorf("failed to start chain process: %w", err)
	}
	
	l.chainRunning = true
	
	// Monitor chain process
	l.wg.Add(1)
	go func() {
		defer l.wg.Done()
		defer func() { l.chainRunning = false }()
		
		if err := l.chainCmd.Wait(); err != nil {
			log.Printf("❌ Chain process exited with error: %v", err)
		} else {
			log.Println("🔗 Chain process exited normally")
		}
		
		// Auto-restart if enabled
		if l.config.AutoRestart && l.ctx.Err() == nil {
			log.Printf("🔄 Restarting chain in %v...", l.config.RestartDelay)
			time.Sleep(l.config.RestartDelay)
			if err := l.startChain(); err != nil {
				log.Printf("❌ Failed to restart chain: %v", err)
			}
		}
	}()
	
	return nil
}

// startBot starts the GXR bot
func (l *GXRLauncher) startBot() error {
	log.Println("🤖 Starting GXR Bot...")
	
	// Build bot command
	args := []string{}
	if l.config.BotConfig != "" {
		args = append(args, "--config", l.config.BotConfig)
	}
	
	l.botCmd = exec.CommandContext(l.ctx, l.config.BotBinary, args...)
	
	// Set up logging
	l.botCmd.Stdout = &PrefixedWriter{prefix: "[BOT] ", writer: os.Stdout}
	l.botCmd.Stderr = &PrefixedWriter{prefix: "[BOT] ", writer: os.Stderr}
	
	// Start bot process
	if err := l.botCmd.Start(); err != nil {
		return fmt.Errorf("failed to start bot process: %w", err)
	}
	
	l.botRunning = true
	
	// Monitor bot process
	l.wg.Add(1)
	go func() {
		defer l.wg.Done()
		defer func() { l.botRunning = false }()
		
		if err := l.botCmd.Wait(); err != nil {
			log.Printf("❌ Bot process exited with error: %v", err)
		} else {
			log.Println("🤖 Bot process exited normally")
		}
		
		// Auto-restart if enabled
		if l.config.AutoRestart && l.ctx.Err() == nil {
			log.Printf("🔄 Restarting bot in %v...", l.config.RestartDelay)
			time.Sleep(l.config.RestartDelay)
			if err := l.startBot(); err != nil {
				log.Printf("❌ Failed to restart bot: %v", err)
			}
		}
	}()
	
	return nil
}

// Stop gracefully stops both processes
func (l *GXRLauncher) Stop() {
	log.Println("🛑 Stopping GXR Launcher...")
	
	// Cancel context to signal all processes to stop
	l.cancel()
	
	// Stop bot first
	if l.botCmd != nil && l.botRunning {
		log.Println("🤖 Stopping bot...")
		if err := l.botCmd.Process.Signal(syscall.SIGTERM); err != nil {
			log.Printf("Error stopping bot: %v", err)
		}
	}
	
	// Stop chain
	if l.chainCmd != nil && l.chainRunning {
		log.Println("🔗 Stopping chain...")
		if err := l.chainCmd.Process.Signal(syscall.SIGTERM); err != nil {
			log.Printf("Error stopping chain: %v", err)
		}
	}
	
	// Wait for all processes to finish
	l.wg.Wait()
	
	log.Println("✅ GXR Launcher stopped gracefully")
}

// GetStatus returns the current status of both processes
func (l *GXRLauncher) GetStatus() map[string]interface{} {
	return map[string]interface{}{
		"chain_running": l.chainRunning,
		"bot_running":   l.botRunning,
		"auto_restart":  l.config.AutoRestart,
	}
}

// PrefixedWriter adds a prefix to log lines
type PrefixedWriter struct {
	prefix string
	writer *os.File
}

func (pw *PrefixedWriter) Write(p []byte) (n int, err error) {
	prefixed := fmt.Sprintf("%s %s", pw.prefix, string(p))
	return pw.writer.Write([]byte(prefixed))
}

// DefaultConfig returns the default launcher configuration
func DefaultConfig() *LauncherConfig {
	return &LauncherConfig{
		ChainBinary:  "./build/gxrchaind",
		BotBinary:    "./bot/gxr-bot",
		ChainHome:    os.ExpandEnv("$HOME/.gxrchaind"),
		LogLevel:     "info",
		AutoRestart:  true,
		RestartDelay: 5 * time.Second,
	}
}

// Main CLI command
func main() {
	var (
		chainBinary string
		botBinary   string
		chainHome   string
		chainConfig string
		botConfig   string
		autoRestart bool
	)
	
	rootCmd := &cobra.Command{
		Use:   "gxr-launcher",
		Short: "GXR Blockchain Launcher",
		Long: `GXR Launcher starts and manages both the GXR blockchain daemon and the validator bot.
		
According to GXR specification, validators must run both the node and bot together.
The launcher ensures both services start together and can be managed as a single unit.`,
		Version: LauncherVersion,
		RunE: func(cmd *cobra.Command, args []string) error {
			// Create configuration
			config := DefaultConfig()
			if chainBinary != "" {
				config.ChainBinary = chainBinary
			}
			if botBinary != "" {
				config.BotBinary = botBinary
			}
			if chainHome != "" {
				config.ChainHome = chainHome
			}
			config.ChainConfig = chainConfig
			config.BotConfig = botConfig
			config.AutoRestart = autoRestart
			
			// Create and start launcher
			launcher := NewGXRLauncher(config)
			if err := launcher.Start(); err != nil {
				return fmt.Errorf("failed to start launcher: %w", err)
			}
			
			// Wait for interrupt signal
			sigChan := make(chan os.Signal, 1)
			signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
			
			log.Println("🏃 GXR Launcher is running. Press Ctrl+C to stop.")
			<-sigChan
			
			// Graceful shutdown
			launcher.Stop()
			return nil
		},
	}
	
	// Add flags
	rootCmd.Flags().StringVar(&chainBinary, "chain-binary", "", "Path to gxrchaind binary")
	rootCmd.Flags().StringVar(&botBinary, "bot-binary", "", "Path to gxr-bot binary")
	rootCmd.Flags().StringVar(&chainHome, "chain-home", "", "Chain home directory")
	rootCmd.Flags().StringVar(&chainConfig, "chain-config", "", "Chain configuration file")
	rootCmd.Flags().StringVar(&botConfig, "bot-config", "", "Bot configuration file")
	rootCmd.Flags().BoolVar(&autoRestart, "auto-restart", true, "Automatically restart failed processes")
	
	// Add status command
	statusCmd := &cobra.Command{
		Use:   "status",
		Short: "Show status of chain and bot processes",
		RunE: func(cmd *cobra.Command, args []string) error {
			// TODO: Implement status checking for running processes
			fmt.Println("Status checking not implemented yet")
			return nil
		},
	}
	rootCmd.AddCommand(statusCmd)
	
	if err := rootCmd.Execute(); err != nil {
		log.Fatal(err)
	}
}