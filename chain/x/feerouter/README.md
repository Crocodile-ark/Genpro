# FeeRouter Module

Modul feerouter mengatur distribusi fee otomatis sesuai spesifikasi GXR blockchain.

## 🎯 Overview

FeeRouter adalah sistem revolusioner yang:
- **Auto distributes fees**: Pembagian fee otomatis per transaksi
- **Dual schemes**: Transaksi umum vs LP farming berbeda pembagian
- **Real-time**: Fee langsung terdistribusi, bukan dikumpulkan dulu
- **LP rewards**: Komunitas LP farming mendapat reward khusus

## 💸 Fee Distribution Schemes

### 🔄 Transaksi Umum (40/30/30)

Untuk transaksi biasa (transfer, staking, dll):

| Penerima     | Persentase | Tujuan                    |
|--------------|------------|---------------------------|
| Validator    | 40%        | Reward validator langsung |
| DEX Pool     | 30%        | Auto refill liquidity     |
| PoS Pool     | 30%        | Reward delegator          |

### 🚜 LP Community Farming (30/25/25/20)

Untuk transaksi LP farming community:

| Penerima     | Persentase | Tujuan                        |
|--------------|------------|-------------------------------|
| Validator    | 30%        | Reward validator              |
| DEX Pool     | 25%        | Auto refill liquidity         |
| LP Community | 25%        | **Reward LP farmers khusus**  |
| PoS Pool     | 20%        | Reward delegator              |

## 🔧 Implementasi

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
    TotalCollected   sdk.Coins // Total fee terkumpul
    TotalToValidators sdk.Coins // Total ke validator
    TotalToDex       sdk.Coins // Total ke DEX pool
    TotalToPos       sdk.Coins // Total ke PoS pool
    TotalToLPRewards sdk.Coins // Total ke LP rewards
}

type LPPool struct {
    Name         string    // Nama pool
    Address      string    // Address pool
    Active       bool      // Status aktif
    TotalRewards sdk.Coin  // Total reward diterima
}
```

### Fee Processing

```go
// ProcessTransactionFees dipanggil untuk setiap transaksi
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

## 🤖 Automasi

### Ante Handler Integration

FeeRouter terintegrasi dengan ante handler untuk:
1. Mengidentifikasi jenis transaksi (umum vs farming)
2. Menghitung pembagian fee sesuai skema
3. Mendistribusikan ke alamat yang tepat
4. Mencatat statistik

### Bot Functions

Validator bot membantu:
- **DEX Auto Refill**: Menggunakan fee yang terkumpul untuk refill pool
- **LP Pool Management**: Mengelola daftar LP pool aktif
- **Fee Monitoring**: Memantau dan alert jika ada anomali

## 🌊 DEX Pool Auto Refill

### Cara Kerja:
1. Fee untuk DEX pool terkumpul di fee collector
2. Bot validator memantau balance pool
3. Jika pool imbalance atau low liquidity detected
4. Bot otomatis refill dari fee collector
5. Refill dilakukan proporsional untuk semua pool aktif

### Supported Pools:
- `GXR/TON` - Pool utama
- `GXR/POLYGON` - Pool sekunder  
- LP Community pools (whitelist otomatis)

## 🏆 LP Community Rewards

### Kriteria LP Farming:
- Transaksi harus dari/ke LP pool address
- Pool harus terdaftar dan aktif
- Volume minimum untuk kualifikasi
- Anti-spam protection

### Distribution:
- 25% fee dari LP farming transaksi
- Dibagi rata ke semua LP pool aktif
- Real-time distribution (per transaksi)
- Auto whitelist pool baru via bot

## 📝 Events

```go
// Fee distribution event
EventTypeFeeDistribution = "fee_distribution"
AttributeKeyFees         = "fees"
AttributeKeyScheme       = "scheme" // "general" atau "farming"
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
- Validator vs LP vs PoS ratios
- Pool performance tracking

### Statistics:
```bash
# Contoh output fee-stats
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
- `ErrInvalidFeeAmount`: Jumlah fee tidak valid
- `ErrInsufficientBalance`: Balance tidak cukup untuk distribusi  
- `ErrInactiveLPPool`: LP pool tidak aktif
- `ErrDistributionFailed`: Gagal mendistribusi ke alamat tujuan

### Safety Mechanisms:
- Fee validation sebelum distribusi
- Fallback jika target address tidak valid
- Rate limiting untuk LP pool registration
- Anti-spam protection

## 🧪 Testing

```bash
# Unit tests
go test ./x/feerouter/keeper/...

# Integration tests dengan halving module
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

### Dengan Modul Lain:
- **Halving**: Distribusi reward bulanan
- **Bank**: Transfer fee ke alamat tujuan
- **Staking**: Reward validator dan delegator
- **Distribution**: PoS pool rewards

### Dengan Bot:
- DEX pool monitoring & refill
- LP pool whitelist management
- Fee anomaly detection
- Telegram alerts untuk events

## 📊 Monitoring

### Metrics:
- Fee volume per scheme (general vs farming)
- Distribution efficiency
- DEX pool health
- LP community activity

### Alerts:
- 💰 Fee distribution berhasil
- ⚠️ Distribusi gagal atau tertunda
- 🌊 DEX pool refill diperlukan
- 🏆 LP pool reward terdistribusi

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