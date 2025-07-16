package main

import (
	"context"
	"log"
	"time"
)

// RewardDistributor handles automatic reward distribution
type RewardDistributor struct {
	config *BotConfig
	
	// Distribution state
	lastDistribution  time.Time
	distributionCount int64
	totalDistributed  string
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
	
	// Connect to chain
	log.Printf("Connecting to chain: %s", rd.config.ChainID)
	// TODO: Initialize chain client
	
	rd.lastDistribution = time.Now()
	rd.totalDistributed = "0ugen"
	
	log.Println("Reward Distributor initialized")
	return nil
}

// Start starts the reward distributor service
func (rd *RewardDistributor) Start(ctx context.Context) error {
	log.Println("Starting Reward Distributor service...")
	
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
			return err
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
	// TODO: Call halving module to distribute monthly rewards
	// This would involve creating and broadcasting a transaction
	// to trigger the monthly distribution in the halving module
	
	log.Println("Distributing halving rewards...")
	log.Println("- 70% to active validators")
	log.Println("- 20% to PoS pool (delegators)")
	log.Println("- 10% to DEX pools")
	
	return nil
}

// GetStatus returns the current reward distributor status
func (rd *RewardDistributor) GetStatus() map[string]interface{} {
	return map[string]interface{}{
		"last_distribution":  rd.lastDistribution,
		"distribution_count": rd.distributionCount,
		"total_distributed":  rd.totalDistributed,
		"next_distribution":  rd.lastDistribution.Add(30 * 24 * time.Hour),
	}
}