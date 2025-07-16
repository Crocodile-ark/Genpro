# Halving Module

Modul halving mengatur sistem pengurangan total supply setiap 5 tahun sesuai spesifikasi GXR blockchain yang telah diperbarui.

## 🎯 Overview

Sistem halving GXR yang baru bekerja dengan prinsip pengurangan total supply:
- **Siklus**: 5 tahun per halving
- **Pengurangan**: 15% dari total supply saat ini per siklus
- **Distribusi**: Bulanan selama 5 tahun untuk setiap siklus
- **Sumber**: Total supply blockchain yang terus berkurang
- **Batas Minimum**: Halving berhenti otomatis jika total supply < 1.000 GXR

## 📊 Cara Kerja Sistem Baru

### Logika Pengurangan Supply:
1. **Setiap 5 tahun**: Sistem menghitung 15% dari total supply saat ini
2. **Distribusi bulanan**: 15% tersebut dibagi menjadi 60 distribusi bulanan
3. **Burn & Mint**: Setiap bulan, sejumlah token di-burn dari supply, kemudian di-mint ulang untuk distribusi reward
4. **Siklus berkelanjutan**: Setelah 5 tahun, sisa supply menjadi 100% baru untuk siklus berikutnya
5. **Auto-stop**: Halving berhenti otomatis jika total supply < 1.000 GXR

### Contoh Ilustrasi:
- **Siklus 1**: Total supply 85,000,000 GXR
  - 15% = 12,750,000 GXR akan terdistribusi dalam 5 tahun
  - Per bulan: 212,500 GXR
  - Setelah 5 tahun: Sisa supply = 72,250,000 GXR

- **Siklus 2**: Total supply 72,250,000 GXR (menjadi 100% baru)
  - 15% = 10,837,500 GXR akan terdistribusi dalam 5 tahun
  - Per bulan: 180,625 GXR
  - Setelah 5 tahun: Sisa supply = 61,412,500 GXR

- **Dan seterusnya...**

## 💰 Distribusi Bulanan

Setiap bulan, reward didistribusikan dengan pembagian yang sama:

- **70%** → Validator aktif (dibagi rata)
- **20%** → PoS Pool untuk delegator
- **10%** → DEX Pool (GXR/TON, GXR/POLYGON)

### Proses Distribusi:
1. **Burn**: Sejumlah token (reward bulanan) di-burn dari total supply
2. **Mint**: Token baru di-mint untuk distribusi reward
3. **Distribute**: Token reward didistribusikan sesuai persentase

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
    CurrentCycle       uint64   // Siklus saat ini
    CycleStartTime     int64    // Waktu mulai siklus
    TotalFundsForCycle sdk.Coin // Total supply di awal siklus
    DistributedInCycle sdk.Coin // Sudah terdistribusi dalam siklus ini
}
```

### Konstanta Penting

```go
const (
    MinimumSupplyThreshold = 1000 * 1e8 // 1,000 GXR minimum
    HalvingReductionRate   = "0.15"     // 15% pengurangan
    MainDenom             = "ugen"      // Denominasi utama
)
```

## 🚀 Keunggulan Sistem Baru

1. **Deflasi Berkelanjutan**: Total supply terus berkurang setiap siklus
2. **Adaptif**: Reward disesuaikan dengan supply saat ini
3. **Sustainable**: Sistem berhenti otomatis di batas minimum
4. **Transparent**: Semua perhitungan berdasarkan supply aktual
5. **Fair**: Distribusi merata dalam periode 5 tahun

## 📈 Monitoring

### Query Commands:
```bash
# Cek total supply saat ini
gxrchaind query bank total --denom ugen

# Cek info halving
gxrchaind query halving info

# Cek distribusi record
gxrchaind query halving distributions
```

### Log Events:
- `Monthly rewards distributed`: Setiap distribusi bulanan
- `Advanced to next halving cycle`: Perpindahan siklus
- `Halving stopped: total supply below minimum threshold`: Auto-stop

## ⚠️ Catatan Penting

1. **Irreversible**: Setiap burn bersifat permanen
2. **Automatic**: Sistem berjalan otomatis tanpa intervensi manual
3. **Predictable**: Semua perhitungan dapat diprediksi dari supply saat ini
4. **Secure**: Menggunakan burn/mint mechanism yang aman

## 🔍 Troubleshooting

### Masalah Umum:
- **Reward = 0**: Kemungkinan total supply sudah di bawah threshold
- **Distribusi gagal**: Periksa validator aktif dan account module
- **Siklus tidak maju**: Periksa parameter waktu dan block time

### Solusi:
```bash
# Restart halving check
gxrchaind tx halving check-cycle --from validator

# Update parameters (jika diperlukan)
gxrchaind tx gov submit-proposal param-change proposal.json
```