# FeeRouter Module

The FeeRouter module manages automatic fee distribution according to the GXR blockchain specification.

## 🎯 Overview

FeeRouter is a revolutionary system that:

- **Auto distributes fees**: Automatically splits transaction fees
- **Dual schemes**: General vs. LP farming transactions have different distribution rules
- **Real-time**: Fees are distributed immediately, not accumulated
- **LP rewards**: LP farming community receives special rewards

## 💸 Fee Distribution Schemes

### 🔄 General Transactions (40/30/30)

For typical transactions (transfer, staking, etc.):

| Recipient | Percentage | Purpose                 |
| --------- | ---------- | ----------------------- |
| Validator | 40%        | Direct validator reward |
| DEX Pool  | 30%        | Auto liquidity refill   |
| PoS Pool  | 30%        | Delegator reward        |

### 🚜 LP Community Farming (30/25/25/20)

For LP farming transactions:

| Recipient    | Percentage | Purpose                      |
| ------------ | ---------- | ---------------------------- |
| Validator    | 30%        | Validator reward             |
| DEX Pool     | 25%        | Auto liquidity refill        |
| LP Community | 25%        | **Special LP farmer reward** |
| PoS Pool     | 20%        | Delegator reward             |

## 🔧 Implementation

### Parameters

```go
type Params struct {
    // General transaction fees (40/30/30)
    GeneralValidatorShare sdk.Dec // 0.40
    GeneralDexShare       sdk.Dec // 0.30  
    GeneralPosShare       sdk.Dec // 0.30

    // LP farming transaction fees (30/25/25/20)
    FarmingValidatorShare sdk.Dec // 0.30
    FarmingDexShare       sdk.Dec // 0.25
    FarmingLPRewardShare  sdk.Dec // 0.25
    FarmingPosShare       sdk.Dec // 0.20
}
```

### State

```go
type FeeStats struct {
    TotalCollected   sdk.Coins // Total collected fees
    TotalToValidators sdk.Coins // Total sent to validators
    TotalToDex       sdk.Coins // Total sent to DEX pool
    TotalToPos       sdk.Coins // Total sent to PoS pool
    TotalToLPRewards sdk.Coins // Total sent to LP rewards
}

type LPPool struct {
    Name         string    // Pool name
    Address      string    // Pool address
    Active       bool      // Active status
    TotalRewards sdk.Coin  // Total rewards received
}
```

### Fee Processing

```go
// ProcessTransactionFees is called for every transaction
func (k Keeper) ProcessTransactionFees(
    ctx sdk.Context,
    fees sdk.Coins,
    isFarmingTransaction bool
) error
```

### Queries

```bash
# Query parameters
gxrchaind q feerouter params

# Query fee statistics
gxrchaind q feerouter fee-stats

# Query LP pools
gxrchaind q feerouter lp-pools
```

## 🤖 Automation

### Ante Handler Integration

FeeRouter integrates with the ante handler to:

1. Identify transaction type (general vs. farming)
2. Calculate fee distribution per scheme
3. Distribute to correct addresses
4. Record statistics

### Bot Functions

Validator bot helps with:

- **DEX Auto Refill**: Uses collected fees to refill pools
- **LP Pool Management**: Manage active LP pool list
- **Fee Monitoring**: Detect anomalies and alert

## 🌊 DEX Pool Auto Refill

### How It Works:

1. Fees for DEX pool accumulate in the fee collector
2. Validator bot monitors pool balance
3. If imbalance or low liquidity detected
4. Bot automatically refills from the fee collector
5. Refill is proportional to all active pools

### Supported Pools:

- `GXR/TON` - Main pool
- `GXR/POLYGON` - Secondary pool
- LP Community pools (auto-whitelisted)

## 🏆 LP Community Rewards

### LP Farming Criteria:

- Transaction must be from/to LP pool address
- Pool must be registered and active
- Minimum volume required to qualify
- Anti-spam protection in place

### Distribution:

- 25% of LP farming transaction fee
- Evenly distributed to all active LP pools
- Real-time distribution (per transaction)
- New pools auto-whitelisted by bot

## 📝 Events

```go
// Fee distribution event
EventTypeFeeDistribution = "fee_distribution"
AttributeKeyFees         = "fees"
AttributeKeyScheme       = "scheme" // "general" or "farming"
AttributeKeyValidator    = "validator_amount"
AttributeKeyDex          = "dex_amount"
AttributeKeyPos          = "pos_amount"
AttributeKeyLP           = "lp_amount"

// LP pool registration event
EventTypeLPPoolRegistered = "lp_pool_registered"
AttributeKeyPoolName      = "pool_name"
AttributeKeyPoolAddress   = "pool_address"
```

## 🔍 Fee Analysis

### Real-time Tracking:

- Total volume per scheme
- Fee efficiency metrics
- Validator vs. LP vs. PoS ratios
- Pool performance tracking

### Statistics:

```bash
# Example output of fee-stats
{
  "total_collected": "1000000000ugen",
  "total_to_validators": "350000000ugen",
  "total_to_dex": "275000000ugen",
  "total_to_pos": "275000000ugen",
  "total_to_lp_rewards": "100000000ugen"
}
```

## 🚨 Error Handling

### Common Errors:

- `ErrInvalidFeeAmount`: Invalid fee amount
- `ErrInsufficientBalance`: Not enough balance for distribution
- `ErrInactiveLPPool`: LP pool is inactive
- `ErrDistributionFailed`: Failed to distribute to destination

### Safety Mechanisms:

- Fee validation before distribution
- Fallback if target address is invalid
- Rate limiting for LP pool registration
- Anti-spam protection

## 🧪 Testing

```bash
# Unit tests
go test ./x/feerouter/keeper/...

# Integration tests with halving module
go test ./x/feerouter/...

# Fee distribution simulation
go test ./x/feerouter/simulation/...
```

### Test Scenarios:

- ✅ General transaction fee distribution
- ✅ LP farming fee distribution
- ✅ DEX pool auto refill
- ✅ LP pool registration/deregistration
- ✅ Error handling & edge cases

## 🔗 Integration

### With Other Modules:

- **Halving**: Monthly reward distribution
- **Bank**: Fee transfers to target addresses
- **Staking**: Validator & delegator rewards
- **Distribution**: PoS pool rewards

### With Bot:

- DEX pool monitoring & refill
- LP pool whitelist management
- Fee anomaly detection
- Telegram alerts for events

## 📊 Monitoring

### Metrics:

- Fee volume per scheme (general vs. farming)
- Distribution efficiency
- DEX pool health
- LP community activity

### Alerts:

- 💰 Fee distribution success
- ⚠️ Failed or delayed distribution
- 🌊 DEX pool refill needed
- 🏆 LP pool rewards distributed

## 🔧 Configuration

### Genesis Setup:

```json
{
  "feerouter": {
    "params": {
      "general_validator_share": "0.40",
      "general_dex_share": "0.30",
      "general_pos_share": "0.30",
      "farming_validator_share": "0.30",
      "farming_dex_share": "0.25",
      "farming_lp_reward_share": "0.25",
      "farming_pos_share": "0.20"
    },
    "fee_stats": {
      "total_collected": []
    },
    "lp_pools": []
  }
}
```

## 📚 References

- [GXR Fee Specification](../../blockchain_gxr_spec.md#fee-distribution)
- [Cosmos SDK Ante Handler](https://docs.cosmos.network/main/basics/gas-fees)
- [Fee Distribution Architecture](./docs/fee-distribution.md)

