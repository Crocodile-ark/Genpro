# GXR Blockchain - Enhanced Immutable Deflationary System

## Overview

GXR is an immutable blockchain built on Cosmos SDK with a revolutionary deflationary halving mechanism. The system features actual supply reduction through a burn-and-mint process, strict validator requirements, and comprehensive bot infrastructure.

## Key Features

### 🔥 Revolutionary Halving System
- **15% Supply Reduction**: Every 5 years, 15% of total supply is allocated to HalvingFund
- **Actual Deflation**: Monthly burn-and-mint process reduces total supply permanently
- **Distribution Period**: 2 years distribution followed by 3 years pause
- **DEX Rewards**: Only years 1-2 of each distribution period
- **Auto-Stop**: Halving stops when supply < 1,000 GXR

### 👥 Validator Requirements
- **Inactivity Penalty**: >10 days inactive per month = reward forfeiture
- **Mandatory Bot**: Validators must run bot or face slashing
- **Uptime Tracking**: Comprehensive monitoring and alerts
- **Reward Eligibility**: Active validators only receive halving rewards

### 🤖 Enhanced Bot System
- **Hourly Rebalancing**: Exactly 1-hour intervals, not per block
- **Price Monitoring**: Monitor-only mode when GXR ≥ $5 for 24 hours
- **Telegram Alerts**: Real-time notifications for all state changes
- **Health Monitoring**: Comprehensive component health checks
- **Validator Tracking**: Heartbeat monitoring and inactivity detection

### 🏛️ Immutable Architecture
- **No Governance**: No parameter changes or upgrades
- **No Smart Contracts**: Pure blockchain functionality
- **No NFTs**: Single token ecosystem
- **One Token**: Only GXR exists

## Quick Start

### Prerequisites
- Go 1.21+
- Cosmos SDK
- Node.js (for launcher)

### Installation

```bash
# Clone repository
git clone https://github.com/your-org/gxr-chain
cd gxr-chain

# Build chain
cd chain
make build

# Build bot
cd ../bot
go build -o gxr-bot

# Build launcher
cd ../launcher
npm install
npm run build
```

### Running the Chain

```bash
# Initialize chain
./build/gxrchaind init mynode --chain-id gxr-1

# Start chain
./build/gxrchaind start
```

### Running the Bot

```bash
# Create bot configuration
cp config/bot.example.yaml config/bot.yaml

# Edit configuration
nano config/bot.yaml

# Start bot
./gxr-bot --config config/bot.yaml
```

## Architecture

### Chain Structure
```
chain/
├── x/halving/          # Halving module
│   ├── keeper/         # Core halving logic
│   ├── types/          # Data structures
│   └── client/         # CLI commands
├── app/                # Application setup
├── cmd/                # Chain binary
└── proto/              # Protocol buffers
```

### Bot Structure
```
bot/
├── main.go             # Main bot service
├── rebalancer.go       # Hourly rebalancing
├── validator_monitor.go # Validator tracking
├── telegram_alert.go   # Alert system
├── dex_manager.go      # DEX integration
├── ibc_relayer.go      # IBC operations
└── reward_distributor.go # Reward distribution
```

## Halving System Details

### Timeline
- **Cycle Duration**: 5 years
- **Distribution**: 2 years active
- **Pause**: 3 years inactive
- **Monthly Distribution**: 30-day intervals

### Distribution Formula
```
Monthly Burn = (Total Supply × 0.15) ÷ 24 months
```

### Reward Allocation
- **70%**: Active validators
- **20%**: Delegators (PoS staking)
- **10%**: DEX pools (years 1-2 only)

### Example Scenario
```
Initial Supply: 1,000,000 GXR
Cycle 1: 850,000 GXR (15% burned)
Cycle 2: 722,500 GXR (15% burned)
Cycle 3: 614,125 GXR (15% burned)
...continues until < 1,000 GXR
```

## Validator Requirements

### Activity Rules
- **Monthly Uptime**: >20 days active required
- **Inactivity Penalty**: >10 days inactive = no rewards
- **Status Tracking**: Real-time monitoring
- **Reward Forfeiture**: Inactive validators lose monthly rewards

### Bot Requirements
- **Mandatory**: All validators must run bot
- **Heartbeat**: 1-minute intervals
- **Slashing**: No bot = validator slashing
- **Monitoring**: Comprehensive health checks

## Bot Configuration

### Example Configuration
```yaml
# Chain connection
chain_rpc: "tcp://localhost:26657"
chain_grpc: "localhost:9090"
chain_id: "gxr-1"

# Validator settings
validator_address: "gxrvaloper1..."
validator_name: "MyValidator"

# Rebalancing (1-hour intervals)
swap_cooldown: "1h"
price_limit: "5.0"
max_swap_daily: "10000"

# Telegram alerts
telegram_enabled: true
telegram_token: "YOUR_BOT_TOKEN"
telegram_chat_id: "YOUR_CHAT_ID"

# Enhanced monitoring
monitoring_enabled: true
health_check_enabled: true
```

## Bot Operations

### Rebalancing Rules
- **Frequency**: Exactly every 1 hour
- **Price Threshold**: $5.00 USD
- **Monitor Mode**: 24-hour suspension when price ≥ $5
- **State Alerts**: All state changes via Telegram

### Alert Categories
- **Info**: General notifications
- **Warning**: Price thresholds, inactivity
- **Error**: Component failures
- **Critical**: Emergency situations
- **Success**: Successful operations

### Rate Limiting
- **Max Alerts**: 10 per minute
- **Queue Size**: 100 pending alerts
- **Retry Logic**: 3 attempts with 5s delay

## Monitoring & Alerts

### Health Checks
- **Rebalancer**: State monitoring
- **Validator Monitor**: Uptime tracking
- **IBC Relayer**: Connection status
- **DEX Manager**: Pool health
- **Reward Distributor**: Distribution status

### Telegram Integration
```bash
# Test connection
./gxr-bot test-telegram

# Send test alert
./gxr-bot send-alert "Test message"
```

### Dashboard Metrics
- Active validators
- Current halving cycle
- Supply statistics
- Bot health status
- Alert history

## Command Line Interface

### Chain Commands
```bash
# Query halving info
gxrchaind query halving info

# Query validator uptime
gxrchaind query halving uptime [validator]

# Query distribution records
gxrchaind query halving distributions
```

### Bot Commands
```bash
# Start bot
./gxr-bot --config config/bot.yaml

# Check status
./gxr-bot status

# Test configuration
./gxr-bot test

# Show version
./gxr-bot version
```

## Security Features

### Immutability
- **No Governance**: Parameters cannot be changed
- **No Upgrades**: Chain version is fixed
- **No Smart Contracts**: Code cannot be modified
- **Deterministic**: All operations are predictable

### Validator Security
- **Slashing**: Non-compliant validators are penalized
- **Monitoring**: Real-time tracking of all validators
- **Mandatory Bots**: Required for participation
- **Transparent**: All actions are publicly auditable

## Network Statistics

### Current Metrics
- **Total Supply**: Dynamic (decreasing)
- **Validator Count**: Active participants
- **Halving Cycle**: Current cycle number
- **Next Distribution**: Countdown timer
- **Bot Health**: System status

### Historical Data
- Supply reduction history
- Validator performance
- Distribution records
- Alert history

## Troubleshooting

### Common Issues

#### Bot Not Starting
```bash
# Check configuration
./gxr-bot test

# Verify chain connection
./gxr-bot status

# Check logs
tail -f logs/gxr-bot.log
```

#### Validator Inactivity
```bash
# Check validator status
gxrchaind query staking validator [validator-address]

# Check uptime record
gxrchaind query halving uptime [validator-address]

# Restart validator
systemctl restart gxrchaind
```

#### Telegram Alerts Not Working
```bash
# Test connection
./gxr-bot test-telegram

# Check configuration
grep -A 3 telegram config/bot.yaml

# Verify token and chat ID
curl -s "https://api.telegram.org/bot$TOKEN/getMe"
```

### Performance Optimization

#### Bot Performance
- **Concurrent Operations**: Configure max_concurrent_ops
- **Retry Logic**: Adjust retry_attempts and retry_delay
- **Health Checks**: Optimize health_check_interval
- **Alert Rate**: Configure telegram rate limiting

#### Chain Performance
- **Pruning**: Configure state pruning
- **Indexing**: Optimize transaction indexing
- **Caching**: Enable state caching
- **Logging**: Adjust log levels

## Development

### Building from Source
```bash
# Build chain
cd chain
make build

# Build bot
cd bot
go build -o gxr-bot

# Run tests
make test
```

### Running Tests
```bash
# Unit tests
go test ./...

# Integration tests
make test-integration

# End-to-end tests
make test-e2e
```

## API Reference

### Chain API
- **REST**: http://localhost:1317
- **gRPC**: localhost:9090
- **WebSocket**: ws://localhost:26657/websocket

### Bot API
- **Health**: GET /health
- **Status**: GET /status
- **Metrics**: GET /metrics

## Support

### Documentation
- **Technical Specs**: See `docs/blockchain_gxr_spec.md`
- **Bot Guide**: See `bot/README.md`
- **API Docs**: See `docs/api.md`

### Community
- **Discord**: [GXR Community](https://discord.gg/gxr)
- **Telegram**: [GXR Announcements](https://t.me/gxr_announcements)
- **GitHub**: [Issues & PRs](https://github.com/your-org/gxr-chain/issues)

## License

MIT License - see LICENSE file for details.

## Disclaimer

This is an immutable blockchain system. Once deployed, the parameters cannot be changed through governance or upgrades. Validators and users should fully understand the system before participating.

---

**GXR Blockchain** - Truly Immutable, Deflationary, and Decentralized

