package main

import (
	"context"
	"fmt"
	"log"
	"time"
)

// RewardDistributor handles automatic reward distribution
type RewardDistributor struct {
	config *BotConfig
	
	// Chain client would be here in real implementation
	chainClient interface{}
	
	// Distribution state
	lastDistribution  time.Time
	distributionCount int64
	totalDistributed  string
	isConnected       bool
}

// NewRewardDistributor creates a new reward distributor instance
func NewRewardDistributor(config *BotConfig) *RewardDistributor {
	return &RewardDistributor{
		config: config,
	}
}

// Initialize initializes the reward distributor
func (rd *RewardDistributor) Initialize() error {
	log.Println("Initializing Reward Distributor...")
	
	// Initialize chain connection
	if err := rd.initializeChainClient(); err != nil {
		return fmt.Errorf("failed to initialize chain client: %w", err)
	}
	
	rd.lastDistribution = time.Now()
	rd.totalDistributed = "0ugen"
	rd.isConnected = true
	
	log.Println("Reward Distributor initialized successfully")
	return nil
}

// initializeChainClient initializes the blockchain client connection
func (rd *RewardDistributor) initializeChainClient() error {
	log.Printf("Connecting to chain: %s", rd.config.ChainID)
	log.Printf("Chain RPC: %s", rd.config.ChainRPC)
	log.Printf("Chain gRPC: %s", rd.config.ChainGRPC)
	
	// In a real implementation, this would create a Cosmos SDK client
	// For now, we'll simulate the connection
	if rd.config.ChainRPC == "" || rd.config.ChainGRPC == "" {
		return fmt.Errorf("chain RPC and gRPC endpoints are required")
	}
	
	// Simulate connection delay
	time.Sleep(1 * time.Second)
	
	log.Println("Chain client connected successfully")
	return nil
}

// Start starts the reward distributor service
func (rd *RewardDistributor) Start(ctx context.Context) error {
	log.Println("Starting Reward Distributor service...")
	
	// Check connection status
	if !rd.isConnected {
		return fmt.Errorf("reward distributor not connected to chain")
	}
	
	// Check every hour for monthly distributions
	ticker := time.NewTicker(1 * time.Hour)
	defer ticker.Stop()
	
	for {
		select {
		case <-ctx.Done():
			log.Println("Reward Distributor stopping...")
			return nil
			
		case <-ticker.C:
			if err := rd.checkAndDistribute(); err != nil {
				log.Printf("Reward Distributor error: %v", err)
			}
		}
	}
}

// checkAndDistribute checks if it's time to distribute rewards and does so
func (rd *RewardDistributor) checkAndDistribute() error {
	// Check if it's time for monthly distribution
	now := time.Now()
	if rd.shouldDistribute(now) {
		log.Println("Time for monthly reward distribution")
		
		// Distribute halving rewards
		if err := rd.distributeHalvingRewards(); err != nil {
			return fmt.Errorf("failed to distribute halving rewards: %w", err)
		}
		
		rd.lastDistribution = now
		rd.distributionCount++
		
		log.Printf("Monthly rewards distributed successfully (cycle %d)", rd.distributionCount)
	}
	
	return nil
}

// shouldDistribute determines if it's time for monthly distribution
func (rd *RewardDistributor) shouldDistribute(now time.Time) bool {
	// Check if 30 days have passed since last distribution
	return now.Sub(rd.lastDistribution) >= (30 * 24 * time.Hour)
}

// distributeHalvingRewards distributes rewards from the halving fund
func (rd *RewardDistributor) distributeHalvingRewards() error {
	log.Println("Distributing halving rewards...")
	
	// In a real implementation, this would:
	// 1. Create a transaction to call the halving module's distribute function
	// 2. Sign and broadcast the transaction
	// 3. Wait for confirmation
	
	// For now, we'll simulate the process
	if err := rd.simulateDistribution(); err != nil {
		return fmt.Errorf("distribution simulation failed: %w", err)
	}
	
	log.Println("- 70% distributed to active validators")
	log.Println("- 20% distributed to PoS pool (delegators)")
	log.Println("- 10% distributed to DEX pools")
	
	return nil
}

// simulateDistribution simulates the distribution process
func (rd *RewardDistributor) simulateDistribution() error {
	// Simulate transaction creation delay
	time.Sleep(2 * time.Second)
	
	// Simulate potential failures
	if rd.distributionCount > 0 && rd.distributionCount%10 == 0 {
		return fmt.Errorf("simulated network error")
	}
	
	// Update total distributed amount (this would come from the actual transaction)
	rd.totalDistributed = fmt.Sprintf("%dugen", (rd.distributionCount+1)*70833)
	
	return nil
}

// GetStatus returns the current reward distributor status
func (rd *RewardDistributor) GetStatus() map[string]interface{} {
	nextDistribution := rd.lastDistribution.Add(30 * 24 * time.Hour)
	timeUntilNext := nextDistribution.Sub(time.Now())
	
	return map[string]interface{}{
		"connected":          rd.isConnected,
		"last_distribution":  rd.lastDistribution,
		"distribution_count": rd.distributionCount,
		"total_distributed":  rd.totalDistributed,
		"next_distribution":  nextDistribution,
		"time_until_next":    timeUntilNext.String(),
		"chain_id":           rd.config.ChainID,
		"chain_rpc":          rd.config.ChainRPC,
		"chain_grpc":         rd.config.ChainGRPC,
	}
}

// ForceDistribution forces a manual distribution (for testing/emergency)
func (rd *RewardDistributor) ForceDistribution() error {
	if !rd.isConnected {
		return fmt.Errorf("not connected to chain")
	}
	
	log.Println("Forcing manual reward distribution...")
	
	if err := rd.distributeHalvingRewards(); err != nil {
		return fmt.Errorf("forced distribution failed: %w", err)
	}
	
	rd.lastDistribution = time.Now()
	rd.distributionCount++
	
	log.Println("Manual distribution completed successfully")
	return nil
}

// Reconnect attempts to reconnect to the chain
func (rd *RewardDistributor) Reconnect() error {
	log.Println("Attempting to reconnect to chain...")
	
	rd.isConnected = false
	
	if err := rd.initializeChainClient(); err != nil {
		return fmt.Errorf("reconnection failed: %w", err)
	}
	
	rd.isConnected = true
	log.Println("Reconnection successful")
	return nil
}