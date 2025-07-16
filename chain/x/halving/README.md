# Halving Module

Modul halving mengatur distribusi reward 5 tahunan sesuai spesifikasi GXR blockchain.

## 🎯 Overview

Sistem halving GXR berbeda dengan blockchain lain:
- **Siklus**: 5 tahun per halving (bukan 4 tahun seperti Bitcoin)
- **Pengurangan**: 15% per siklus (bukan 50%)
- **Distribusi**: Bulanan selama 5 tahun (bukan sekaligus)
- **Sumber**: Fixed pool 21.25M GXR (25% total supply)

## 📊 Halving Schedule

| Halving | Periode      | Dana GXR    | Pengurangan | Per Bulan |
|---------|--------------|-------------|-------------|-----------|
| 1       | Tahun 1–5    | 4,250,000   | —           | 70,833    |
| 2       | Tahun 6–10   | 3,612,500   | -15%        | 60,208    |
| 3       | Tahun 11–15  | 3,070,625   | -15%        | 51,177    |
| 4       | Tahun 16–20  | 2,610,032   | -15%        | 43,500    |
| 5       | Tahun 21–25  | 2,218,528   | -15%        | 36,975    |

**Total**: 21,250,000 GXR akan terdistribusi selama 25 tahun.

## 💰 Distribusi Bulanan

Setiap bulan, reward didistribusikan dengan pembagian:

- **70%** → Validator aktif (dibagi rata)
- **20%** → PoS Pool untuk delegator
- **10%** → DEX Pool (GXR/TON, GXR/POLYGON)

### Contoh Bulan Pertama:
- Total reward: 70,833 GXR
- Validator (70%): 49,583 GXR
- Delegator (20%): 14,167 GXR  
- DEX Pool (10%): 7,083 GXR

## 🔧 Implementasi

### Parameters

```go
type Params struct {
    HalvingCycleDuration time.Duration // 5 tahun
    ValidatorShare       sdk.Dec       // 0.70 (70%)
    DelegatorShare       sdk.Dec       // 0.20 (20%)
    DexShare            sdk.Dec       // 0.10 (10%)
}
```

### State

```go
type HalvingInfo struct {
    CurrentCycle       uint64   // Siklus saat ini (1-5)
    CycleStartTime     int64    // Waktu mulai siklus
    TotalFundsForCycle sdk.Coin // Total dana siklus ini
    DistributedInCycle sdk.Coin // Sudah terdistribusi
}
```

### Messages

Modul halving tidak memiliki message dari user. Semua distribusi otomatis via BeginBlocker.

### Queries

```bash
# Query parameters
gxrchaind q halving params

# Query halving info
gxrchaind q halving halving-info

# Query distribution history
gxrchaind q halving distribution-history
```

## 🤖 Automasi

### BeginBlocker

Setiap blok, modul mengecek:
1. Apakah sudah 30 hari sejak distribusi terakhir?
2. Jika ya, lakukan distribusi bulanan
3. Apakah sudah 5 tahun sejak siklus dimulai?
4. Jika ya, maju ke siklus berikutnya

### Bot Integration

Validator bot memantau dan membantu:
- Memastikan distribusi berjalan tepat waktu
- Mengirim alert Telegram saat distribusi
- Monitoring kesehatan halving fund

## 📝 Events

```go
// Monthly distribution event
EventTypeDistribution = "halving_distribution"
AttributeKeyAmount    = "amount"
AttributeKeyCycle     = "cycle"
AttributeKeyMonth     = "month"

// Cycle advancement event  
EventTypeCycleAdvance = "halving_cycle_advance"
AttributeKeyNewCycle  = "new_cycle"
AttributeKeyNewFunds  = "new_funds"
```

## 🚨 Error Handling

### Common Errors:
- `ErrInvalidCycle`: Siklus halving tidak valid (harus 1-5)
- `ErrInsufficientFunds`: Dana halving tidak mencukupi
- `ErrDistributionTooEarly`: Belum waktunya distribusi bulanan

### Recovery:
- Jika distribusi gagal, akan dicoba ulang di blok berikutnya
- Semua error dicatat dalam event log
- Bot akan mengirim alert jika ada masalah

## 🧪 Testing

```bash
# Unit tests
go test ./x/halving/keeper/...

# Integration tests
go test ./x/halving/...

# Simulation tests
go test ./x/halving/simulation/...
```

### Test Scenarios:
- ✅ Distribusi bulanan normal
- ✅ Perpindahan siklus halving
- ✅ Validasi parameter
- ✅ Error handling
- ✅ Genesis import/export

## 🔍 Monitoring

### Metrics:
- Total terdistribusi per siklus
- Waktu distribusi terakhir
- Sisa dana di halving fund
- Jumlah validator aktif yang menerima reward

### Alerts:
- 🔔 Distribusi bulanan berhasil
- ⚠️ Distribusi tertunda atau gagal
- 🚨 Siklus halving berganti
- 💰 Halving fund rendah

## 📚 References

- [GXR Specification](../../blockchain_gxr_spec.md)
- [Cosmos SDK Modules](https://docs.cosmos.network/main/modules/)
- [BeginBlocker/EndBlocker](https://docs.cosmos.network/main/building-modules/beginblock-endblock)