# GXR Blockchain Technical Specification

## Version: 2.0.0
**Last Updated**: December 2024

---

## Table of Contents

1. [System Overview](#system-overview)
2. [Halving Mechanism](#halving-mechanism)
3. [Validator Requirements](#validator-requirements)
4. [Bot Infrastructure](#bot-infrastructure)
5. [Immutable Architecture](#immutable-architecture)
6. [Technical Implementation](#technical-implementation)
7. [Security Features](#security-features)
8. [Performance Specifications](#performance-specifications)
9. [API Reference](#api-reference)
10. [Deployment Guide](#deployment-guide)

---

## System Overview

### Core Principles

GXR is an immutable blockchain built on Cosmos SDK with the following core principles:

- **Immutability**: No governance, no upgrades, no parameter changes
- **Deflation**: Actual supply reduction through burn-and-mint mechanism
- **Validator Accountability**: Strict requirements with automated enforcement
- **Predictability**: All operations follow deterministic rules
- **Single Token**: Only GXR exists in the ecosystem

### Architecture

```
┌─────────────────────────────────────────────────────────────┐
│                    GXR Blockchain                           │
├─────────────────────────────────────────────────────────────┤
│  Cosmos SDK Base Layer                                      │
│  ├── Tendermint Consensus                                   │
│  ├── Cosmos SDK Framework                                   │
│  └── Custom Modules                                         │
│      ├── x/halving (Custom)                                 │
│      ├── x/bank (Modified)                                  │
│      ├── x/staking (Modified)                               │
│      └── x/slashing (Modified)                              │
├─────────────────────────────────────────────────────────────┤
│  Bot Infrastructure                                         │
│  ├── Validator Monitor                                      │
│  ├── Rebalancer                                             │
│  ├── IBC Relayer                                            │
│  ├── DEX Manager                                            │
│  ├── Reward Distributor                                     │
│  └── Telegram Alert System                                  │
└─────────────────────────────────────────────────────────────┘
```

---

## Halving Mechanism

### Timeline Structure

```
Halving Cycle (5 years)
├── Distribution Period (2 years)
│   ├── Monthly Distribution (30 days)
│   │   ├── Burn Phase
│   │   └── Mint & Distribute Phase
│   └── Repeat for 24 months
└── Pause Period (3 years)
    └── No Distribution
```

### Mathematical Formula

#### Halving Fund Calculation
```
HalvingFund = CurrentSupply × 0.15
```

#### Monthly Distribution
```
MonthlyAmount = HalvingFund ÷ 24
```

#### Burn-and-Mint Process
```
1. Burn: MonthlyAmount from CurrentSupply
2. Mint: MonthlyAmount for distribution
3. Net Effect: CurrentSupply decreases by MonthlyAmount
```

### Distribution Allocation

| Recipient | Percentage | Years Active | Conditions |
|-----------|------------|--------------|------------|
| Validators | 70% | 1-2 | >20 days uptime/month |
| Delegators | 20% | 1-2 | Via fee pool |
| DEX Pools | 10% | 1-2 only | Bot managed |

### Supply Trajectory

```
Cycle 1: 1,000,000 GXR
├── HalvingFund: 150,000 GXR
├── Monthly: 6,250 GXR
└── Remaining: 850,000 GXR

Cycle 2: 850,000 GXR
├── HalvingFund: 127,500 GXR
├── Monthly: 5,312.5 GXR
└── Remaining: 722,500 GXR

Cycle 3: 722,500 GXR
├── HalvingFund: 108,375 GXR
├── Monthly: 4,515.625 GXR
└── Remaining: 614,125 GXR

...continues until supply < 1,000 GXR
```

### Implementation Details

#### Halving Module (`x/halving`)

**Key Functions:**
- `CheckAndAdvanceHalvingCycle()`: Monitors 5-year cycles
- `ShouldDistribute()`: Checks 30-day distribution timing
- `DistributeHalvingRewards()`: Executes burn-and-mint process
- `CheckAndUpdateDistributionStatus()`: Manages 2-year/3-year periods

**State Management:**
```go
type HalvingInfo struct {
    CurrentCycle       uint64
    CycleStartTime     int64
    TotalSupply        types.Coin
    HalvingFund        types.Coin
    DistributionActive bool
    DistributionStart  int64
    DistributedAmount  types.Coin
    PauseStart         int64
    LastMonthlyDistrib int64
}
```

---

## Validator Requirements

### Activity Metrics

#### Uptime Calculation
```
Monthly Uptime = (Active Days / 30) × 100%
Reward Eligibility = Monthly Uptime > 66.67% (>20 days)
```

#### Tracking Implementation
```go
type ValidatorUptime struct {
    ValidatorAddress string
    CurrentMonth     uint64
    InactiveDays     uint64
    LastActiveTime   time.Time
    LastCheck        time.Time
    MissedBlocks     uint64
    BotRunning       bool
    LastBotHeartbeat time.Time
    RewardEligible   bool
}
```

### Slashing Conditions

| Condition | Penalty | Recovery |
|-----------|---------|----------|
| Bot Not Running | Slashing | Restart bot |
| >10 Days Inactive | Reward Forfeiture | Increase uptime |
| Extended Downtime | Jail + Slash | Manual recovery |

### Bot Requirements

#### Mandatory Components
- **Validator Monitor**: Uptime tracking
- **Rebalancer**: Hourly operations
- **IBC Relayer**: Cross-chain operations
- **DEX Manager**: Pool management
- **Reward Distributor**: Reward handling
- **Telegram Alert**: Notifications

#### Heartbeat Protocol
```
Interval: 1 minute
Timeout: 5 minutes
Grace Period: 10 minutes
Slashing Delay: 1 hour
```

---

## Bot Infrastructure

### Rebalancer Specifications

#### Operational Constraints
- **Frequency**: Exactly 1 hour intervals
- **Price Threshold**: $5.00 USD
- **Monitor Mode**: 24-hour suspension
- **Recovery**: Automatic after 24 hours

#### State Machine
```
Active → (Price ≥ $5) → Monitor Only → (24h + Price < $5) → Active
Active → (Error) → Error State → (1h timeout) → Active
Active → (Emergency) → Emergency Stop → (Manual) → Active
```

### Validator Monitor

#### Tracking Metrics
- **Uptime**: Daily activity status
- **Bot Health**: Heartbeat monitoring
- **Reward Eligibility**: Monthly calculations
- **Alert Triggers**: Inactivity thresholds

#### Monthly Reset Process
```go
func (vm *ValidatorMonitor) performMonthlyReset() {
    // Store previous month statistics
    vm.monthlyStats[oldMonth] = &MonthlyStats{
        TotalValidators:    vm.totalValidators,
        ActiveValidators:   vm.activeValidators,
        InactiveValidators: vm.totalInactiveValidators,
        ForfeitedRewards:   vm.totalForfeitedRewards,
        AverageUptime:      vm.calculateAverageUptime(),
    }
    
    // Reset all validator counters
    for _, status := range vm.validators {
        status.CurrentMonth = vm.currentMonth
        status.InactiveDays = 0
        status.RewardEligible = true
    }
}
```

### Telegram Alert System

#### Alert Categories
```go
const (
    AlertTypeInfo AlertType = iota
    AlertTypeWarning
    AlertTypeError
    AlertTypeCritical
    AlertTypeSuccess
)
```

#### Rate Limiting
- **Max Rate**: 10 alerts per minute
- **Queue Size**: 100 pending alerts
- **Retry Logic**: 3 attempts with 5-second delays
- **Emergency Bypass**: Critical alerts skip rate limits

---

## Immutable Architecture

### Genesis Parameters

**Chain Parameters (FIXED):**
```json
{
  "chain_id": "gxr-1",
  "initial_supply": "85000000000000ugen",
  "min_supply_threshold": "1000000000ugen",
  "halving_cycle_duration": "157680000s",
  "distribution_period": "63072000s",
  "pause_period": "94608000s",
  "halving_reduction_rate": "0.15",
  "validator_share": "0.70",
  "delegator_share": "0.20",
  "dex_share": "0.10"
}
```

**Disabled Features:**
- Governance module
- Smart contract execution
- NFT support
- New token creation
- Parameter changes
- Software upgrades

### Consensus Parameters

```json
{
  "block": {
    "max_bytes": "1048576",
    "max_gas": "100000000",
    "time_iota_ms": "1000"
  },
  "evidence": {
    "max_age_num_blocks": "100000",
    "max_age_duration": "172800000000000",
    "max_bytes": "1048576"
  },
  "validator": {
    "pub_key_types": ["ed25519"]
  }
}
```

---

## Technical Implementation

### Halving Module Structure

```
x/halving/
├── keeper/
│   ├── keeper.go           # Core halving logic
│   ├── distribution.go     # Reward distribution
│   ├── validator.go        # Validator tracking
│   └── query.go           # Query handlers
├── types/
│   ├── halving.pb.go      # Protobuf definitions
│   ├── keys.go            # Storage keys
│   ├── params.go          # Parameters
│   └── query.pb.go        # Query types
├── client/
│   └── cli/
│       ├── query.go       # CLI queries
│       └── tx.go          # CLI transactions
└── module.go              # Module definition
```

### Key Storage Design

```go
var (
    CurrentHalvingKey     = []byte("current_halving")
    LastDistributionKey   = []byte("last_distribution")
    ValidatorUptimeKey    = []byte("validator_uptime")
    DistributionRecordKey = []byte("distribution_record")
)
```

### State Transitions

#### Halving Cycle Advancement
```go
func (k Keeper) CheckAndAdvanceHalvingCycle(ctx sdk.Context) error {
    info, found := k.GetHalvingInfo(ctx)
    cycleStart := time.Unix(info.CycleStartTime, 0)
    
    if ctx.BlockTime().Sub(cycleStart) >= HalvingCycleDuration {
        return k.advanceToNextCycle(ctx, info)
    }
    
    return nil
}
```

#### Monthly Distribution
```go
func (k Keeper) DistributeHalvingRewards(ctx sdk.Context) error {
    if !k.ShouldDistribute(ctx) {
        return nil
    }
    
    monthlyAmount := k.calculateMonthlyDistribution(ctx, info)
    
    // Burn from total supply
    k.bankKeeper.BurnCoins(ctx, types.ModuleName, sdk.NewCoins(monthlyAmount))
    
    // Mint for distribution
    k.bankKeeper.MintCoins(ctx, types.ModuleName, sdk.NewCoins(monthlyAmount))
    
    // Distribute to recipients
    return k.distributeRewards(ctx, monthlyAmount, info)
}
```

---

## Security Features

### Immutability Guarantees

#### Genesis Lock
- All parameters fixed at genesis
- No governance module compiled
- No upgrade handler registered
- No parameter change proposals

#### Code Verification
```bash
# Verify immutability
gxrchaind query gov params voting    # Should fail
gxrchaind query upgrade plan         # Should fail
gxrchaind query params subspace      # Should be empty
```

### Slashing Protection

#### Progressive Penalties
```go
type SlashingSchedule struct {
    FirstOffense:  sdk.NewDecWithPrec(1, 2)  // 1%
    SecondOffense: sdk.NewDecWithPrec(5, 2)  // 5%
    ThirdOffense:  sdk.NewDecWithPrec(10, 2) // 10%
    Jailing:      true
}
```

#### Bot Enforcement
- Heartbeat monitoring
- Automatic slashing for non-compliance
- Grace period for temporary failures
- Manual recovery procedures

---

## Performance Specifications

### Throughput Metrics

| Metric | Value | Unit |
|--------|-------|------|
| Block Time | 6 | seconds |
| TPS | 1,000 | transactions/second |
| Block Size | 1 | MB |
| Gas Limit | 100,000,000 | gas units |

### Network Requirements

#### Validator Hardware
- **CPU**: 8+ cores
- **RAM**: 32+ GB
- **Storage**: 1+ TB NVMe SSD
- **Network**: 1 Gbps connection

#### Bot Hardware
- **CPU**: 2+ cores
- **RAM**: 4+ GB
- **Storage**: 100+ GB SSD
- **Network**: 100 Mbps connection

### Monitoring Thresholds

```yaml
alerts:
  high_memory_usage: 80%
  high_cpu_usage: 85%
  disk_space_low: 90%
  network_latency_high: 500ms
  block_time_slow: 10s
  peer_count_low: 5
```

---

## API Reference

### Chain Queries

#### Halving Information
```bash
# Get current halving info
gxrchaind query halving info

# Get validator uptime
gxrchaind query halving uptime [validator-address]

# Get distribution records
gxrchaind query halving distributions --limit 100

# Get monthly statistics
gxrchaind query halving monthly-stats --month 12
```

#### Bank Queries
```bash
# Get total supply
gxrchaind query bank total --denom ugen

# Get account balance
gxrchaind query bank balance [address] ugen

# Get supply information
gxrchaind query bank supply ugen
```

### Bot API Endpoints

#### Health Check
```bash
GET /health
Response: {"status": "healthy", "timestamp": "2024-12-01T00:00:00Z"}
```

#### Component Status
```bash
GET /status
Response: {
  "version": "2.0.0",
  "running": true,
  "uptime": "24h30m15s",
  "components": {
    "rebalancer": {"state": "active"},
    "validator_monitor": {"active_validators": 150},
    "telegram_alert": {"alerts_sent": 1250}
  }
}
```

#### Metrics
```bash
GET /metrics
Response: Prometheus-formatted metrics
```

---

## Deployment Guide

### Prerequisites

#### System Requirements
```bash
# Ubuntu 22.04 LTS
sudo apt update
sudo apt install -y build-essential git curl jq

# Go 1.21+
wget https://go.dev/dl/go1.21.0.linux-amd64.tar.gz
sudo tar -C /usr/local -xzf go1.21.0.linux-amd64.tar.gz
echo 'export PATH=$PATH:/usr/local/go/bin' >> ~/.bashrc
```

#### Network Configuration
```bash
# Open required ports
sudo ufw allow 26656/tcp  # P2P
sudo ufw allow 26657/tcp  # RPC
sudo ufw allow 1317/tcp   # REST API
sudo ufw allow 9090/tcp   # gRPC
```

### Chain Deployment

#### 1. Build Binary
```bash
git clone https://github.com/your-org/gxr-chain
cd gxr-chain/chain
make build
sudo cp build/gxrchaind /usr/local/bin/
```

#### 2. Initialize Node
```bash
gxrchaind init [moniker] --chain-id gxr-1
gxrchaind add-genesis-account [address] 1000000000000ugen
gxrchaind gentx [key-name] 1000000ugen --chain-id gxr-1
gxrchaind collect-gentxs
```

#### 3. Configure Node
```bash
# Edit config.toml
sed -i 's/timeout_commit = "5s"/timeout_commit = "6s"/g' ~/.gxrchaind/config/config.toml
sed -i 's/prometheus = false/prometheus = true/g' ~/.gxrchaind/config/config.toml

# Edit app.toml
sed -i 's/minimum-gas-prices = ""/minimum-gas-prices = "0.001ugen"/g' ~/.gxrchaind/config/app.toml
```

#### 4. Start Node
```bash
# Create systemd service
sudo tee /etc/systemd/system/gxrchaind.service > /dev/null <<EOF
[Unit]
Description=GXR Chain Node
After=network-online.target

[Service]
User=$(whoami)
ExecStart=/usr/local/bin/gxrchaind start
Restart=always
RestartSec=3
LimitNOFILE=65535

[Install]
WantedBy=multi-user.target
EOF

sudo systemctl daemon-reload
sudo systemctl enable gxrchaind
sudo systemctl start gxrchaind
```

### Bot Deployment

#### 1. Build Bot
```bash
cd gxr-chain/bot
go build -o gxr-bot
sudo cp gxr-bot /usr/local/bin/
```

#### 2. Configure Bot
```bash
mkdir -p ~/.gxr-bot/config
cat > ~/.gxr-bot/config/bot.yaml << EOF
chain_rpc: "tcp://localhost:26657"
chain_grpc: "localhost:9090"
chain_id: "gxr-1"
validator_address: "gxrvaloper1..."
validator_name: "MyValidator"
telegram_enabled: true
telegram_token: "YOUR_BOT_TOKEN"
telegram_chat_id: "YOUR_CHAT_ID"
monitoring_enabled: true
health_check_enabled: true
EOF
```

#### 3. Start Bot
```bash
# Create systemd service
sudo tee /etc/systemd/system/gxr-bot.service > /dev/null <<EOF
[Unit]
Description=GXR Bot Service
After=network-online.target gxrchaind.service

[Service]
User=$(whoami)
ExecStart=/usr/local/bin/gxr-bot --config ~/.gxr-bot/config/bot.yaml
Restart=always
RestartSec=3
LimitNOFILE=65535

[Install]
WantedBy=multi-user.target
EOF

sudo systemctl daemon-reload
sudo systemctl enable gxr-bot
sudo systemctl start gxr-bot
```

### Monitoring Setup

#### 1. Install Prometheus
```bash
wget https://github.com/prometheus/prometheus/releases/download/v2.40.0/prometheus-2.40.0.linux-amd64.tar.gz
tar xzf prometheus-2.40.0.linux-amd64.tar.gz
sudo mv prometheus-2.40.0.linux-amd64/prometheus /usr/local/bin/
```

#### 2. Configure Prometheus
```yaml
# prometheus.yml
global:
  scrape_interval: 15s

scrape_configs:
  - job_name: 'gxr-chain'
    static_configs:
      - targets: ['localhost:26660']
  
  - job_name: 'gxr-bot'
    static_configs:
      - targets: ['localhost:8080']
```

#### 3. Install Grafana
```bash
sudo apt-get install -y software-properties-common
sudo add-apt-repository "deb https://packages.grafana.com/oss/deb stable main"
wget -q -O - https://packages.grafana.com/gpg.key | sudo apt-key add -
sudo apt-get update
sudo apt-get install grafana
```

---

## Troubleshooting

### Common Issues

#### Chain Won't Start
```bash
# Check logs
journalctl -u gxrchaind -f

# Common fixes
gxrchaind unsafe-reset-all
gxrchaind start --pruning=nothing
```

#### Bot Connection Issues
```bash
# Test chain connection
gxr-bot test

# Check bot logs
journalctl -u gxr-bot -f

# Verify configuration
gxr-bot status
```

#### Telegram Alerts Not Working
```bash
# Test token
curl -s "https://api.telegram.org/bot$TOKEN/getMe"

# Test chat ID
curl -s "https://api.telegram.org/bot$TOKEN/sendMessage?chat_id=$CHAT_ID&text=test"
```

### Performance Tuning

#### Chain Optimization
```bash
# Increase file limits
echo "* soft nofile 65536" >> /etc/security/limits.conf
echo "* hard nofile 65536" >> /etc/security/limits.conf

# Optimize kernel parameters
echo "net.core.rmem_max = 134217728" >> /etc/sysctl.conf
echo "net.core.wmem_max = 134217728" >> /etc/sysctl.conf
sysctl -p
```

#### Bot Optimization
```yaml
# bot.yaml
max_concurrent_ops: 20
retry_attempts: 5
retry_delay: "3s"
health_check_interval: "15s"
```

---

## Conclusion

The GXR blockchain represents a new paradigm in immutable blockchain design, combining:

- **Predictable Deflation**: Guaranteed supply reduction through mathematical precision
- **Validator Accountability**: Automated enforcement of participation requirements
- **Operational Excellence**: Comprehensive monitoring and alerting systems
- **Immutable Guarantee**: No governance or upgrade mechanisms

This specification provides the technical foundation for a truly immutable, deflationary blockchain ecosystem that operates with mathematical certainty and transparent accountability.

---

**Document Version**: 2.0.0  
**Last Updated**: December 2024  
**Next Review**: N/A (Immutable Specification)

