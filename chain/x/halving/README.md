# Halving Module

This module governs the halving mechanism for reducing total supply every 5 years according to the updated GXR blockchain specification.

## 🎯 Overview

The new GXR halving system operates based on the following principles:

- **Cycle**: Every 5 years
- **Reduction**: 15% of the current total supply per cycle
- **Distribution**: Monthly over a 5-year period per cycle
- **Source**: The total blockchain supply, which continuously decreases
- **Minimum Threshold**: Halving stops automatically if the total supply falls below 1,000 GXR

## 📊 How the New System Works

### Supply Reduction Logic:

1. **Every 5 years**: The system calculates 15% of the current total supply
2. **Monthly distribution**: The 15% is split into 60 monthly distributions
3. **Burn & Mint**: Every month, tokens are burned from the supply, then re-minted for reward distribution
4. **Continuous cycle**: After 5 years, the remaining supply becomes the new 100% for the next cycle
5. **Auto-stop**: Halving ends automatically if the total supply < 1,000 GXR

### Example Illustration:

- **Cycle 1**: Total supply is 85,000,000 GXR

  - 15% = 12,750,000 GXR will be distributed over 5 years
  - Per month: 212,500 GXR
  - After 5 years: Remaining supply = 72,250,000 GXR

- **Cycle 2**: Total supply is 72,250,000 GXR (becomes the new 100%)

  - 15% = 10,837,500 GXR will be distributed over 5 years
  - Per month: 180,625 GXR
  - After 5 years: Remaining supply = 61,412,500 GXR

- **And so on...**

## 💰 Monthly Distribution

Each month, rewards are distributed evenly as follows:

- **70%** → Active validators (shared equally)
- **20%** → PoS Pool for delegators
- **10%** → DEX Pool (GXR/TON, GXR/POLYGON)

### Distribution Process:

1. **Burn**: Monthly reward tokens are burned from the total supply
2. **Mint**: New tokens are minted for distribution
3. **Distribute**: Rewards are allocated according to the specified percentages

## 🔧 Implementation

### Parameters

```go
type Params struct {
    HalvingCycleDuration time.Duration // 5 years
    ValidatorShare       sdk.Dec       // 0.70 (70%)
    DelegatorShare       sdk.Dec       // 0.20 (20%)
    DexShare             sdk.Dec       // 0.10 (10%)
}
```

### State

```go
type HalvingInfo struct {
    CurrentCycle       uint64   // Current cycle
    CycleStartTime     int64    // Cycle start time
    TotalFundsForCycle sdk.Coin // Total supply at the start of the cycle
    DistributedInCycle sdk.Coin // Amount distributed within the cycle
}
```

### Key Constants

```go
const (
    MinimumSupplyThreshold = 1000 * 1e8 // Minimum 1,000 GXR
    HalvingReductionRate   = "0.15"     // 15% reduction
    MainDenom              = "ugen"     // Main denomination
)
```

## 🚀 System Advantages

1. **Sustainable Deflation**: Supply decreases with each cycle
2. **Adaptive**: Rewards are based on the current supply
3. **Self-Regulating**: System halts automatically at the minimum threshold
4. **Transparent**: All calculations are based on actual supply
5. **Fair**: Equal monthly distribution over 5 years

## 📈 Monitoring

### Query Commands:

```bash
# Check current total supply
gxrchaind query bank total --denom ugen

# Check halving info
gxrchaind query halving info

# Check distribution records
gxrchaind query halving distributions
```

### Log Events:

- `Monthly rewards distributed`: Every monthly distribution
- `Advanced to next halving cycle`: Cycle progression
- `Halving stopped: total supply below minimum threshold`: Auto-stop event

## ⚠️ Important Notes

1. **Irreversible**: Every burn is permanent
2. **Automatic**: The system operates without manual intervention
3. **Predictable**: All calculations can be projected from current supply
4. **Secure**: Utilizes a secure burn/mint mechanism

## 🔍 Troubleshooting

### Common Issues:

- **Reward = 0**: Possibly the supply has dropped below the threshold
- **Distribution failed**: Check active validators and account module
- **Cycle not progressing**: Verify timing parameters and block time

### Solutions:

```bash
# Restart halving check
gxrchaind tx halving check-cycle --from validator

# Update parameters (if needed)
gxrchaind tx gov submit-proposal param-change proposal.json
```

