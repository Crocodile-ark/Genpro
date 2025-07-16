# GXR (Gen X Raider) Blockchain Specification — UPDATED (2024)

> ✅ **Sistem Halving Dinamis dengan Pengurangan Total Supply**
> ✅ Dirancang tanpa smart contract, anti-inflasi, berbasis PoS & IBC
> ✅ Fokus pada efisiensi, distribusi adil, dan desentralisasi otomatis

---

## 1. IDENTITAS DASAR

| Komponen           | Rincian                               |
| ------------------ | ------------------------------------- |
| **Nama Chain**     | Gen X Raider (GXR)                    |
| **Ticker Token**   | GXR                                   |
| **Denom Terkecil** | `ugen` (1 GXR = 100,000,000 ugen)     |
| **Total Supply**   | 85,000,000 GXR *(awal, akan berkurang)* |
| **Desimal**        | 8                                     |
| **Smart Contract** | Tidak digunakan                       |
| **Konsensus**      | Proof-of-Stake (PoS)                  |
| **IBC Support**    | Aktif (GXR/TON, GXR/POLYGON, dll)     |
| **Max Validator**  | 85 node                               |
| **Waktu blok**     | 15 detik per blok                     |
---

## 2. TOKENOMICS GXR

### Total Supply Awal: 85,000,000 GXR

| Alokasi             | Jumlah (GXR) | Persentase | Keterangan                                              |
| ------------------- | ------------ | ---------- | ------------------------------------------------------- |
| Airdrop & Farming   | 17,000,000   | 20%        | Distribusi awal via Telegram bot farming                |
| Developer Core      | 5,950,000    | 7%         | Vesting keras 5 tahun, 10% unlock tiap 6 bulan          |
| Tim Inti (3 orang)  | 5,950,000    | 7%         | 3% / 2% / 2%, soft vesting 3 tahun                      |
| LP & Market         | 8,500,000    | 10%        | Likuiditas awal (GXR/TON, GXR/POLYGON, dll)             |
| Grant (3–7 pihak)   | 8,500,000    | 10%        | Hibah proyek dan mitra kolaborasi                       |
| Pool Staking (PoS)  | 8,500,000    | 10%        | Reward untuk delegator aktif                            |
| **Halving Reserve** | 21,250,000   | 25%        | **Tidak digunakan dalam sistem baru**                   |
| Cadangan/Ekspansi   | 8,500,000    | 10%        | Dana darurat dan pengembangan ekosistem                 |
| Validator Awal (30) | 850,000      | 1%         | 0.5% tahun 1 dan 0.5% tahun 2 jika aktif >20 hari/bulan |

**⚠️ Catatan Penting**: Dalam sistem halving yang baru, tidak ada lagi "Halving Fund" yang terpisah. Sistem halving bekerja langsung dengan total supply yang ada.

---

## 3. SISTEM HALVING DINAMIS BARU

### 🔥 Revolusi Sistem Halving

Sistem halving GXR yang baru **mengubah cara kerja fundamental** dari sistem reward tradisional:

- **Bukan lagi** distribusi dari pool tetap
- **Sekarang** pengurangan total supply secara langsung
- **Setiap 5 tahun**: 15% dari total supply saat ini dikurangi
- **Deflasi berkelanjutan**: Total supply terus berkurang selamanya

### Cara Kerja Sistem Baru:

1. **Perhitungan Bulanan**: 
   - Sistem menghitung 15% dari total supply saat ini
   - Dibagi menjadi 60 distribusi bulanan (5 tahun)
   - Setiap bulan: `(Current Supply × 0.15) ÷ 60`

2. **Proses Burn & Mint**:
   - Sejumlah reward bulanan di-**burn** dari total supply
   - Jumlah yang sama di-**mint** untuk distribusi reward
   - **Net effect**: Total supply berkurang setiap bulan

3. **Siklus Berkelanjutan**:
   - Setelah 5 tahun, sisa supply menjadi **100% baru** untuk siklus berikutnya
   - Tidak ada batasan jumlah siklus
   - Sistem berhenti otomatis jika supply < 1.000 GXR

### Proyeksi Siklus Halving:

| Siklus | Periode      | Supply Awal    | 15% Reward     | Supply Akhir   | Bulanan (±)  |
|--------|--------------|----------------|----------------|----------------|--------------|
| 1      | Tahun 1–5    | 85,000,000     | 12,750,000     | 72,250,000     | 212,500      |
| 2      | Tahun 6–10   | 72,250,000     | 10,837,500     | 61,412,500     | 180,625      |
| 3      | Tahun 11–15  | 61,412,500     | 9,211,875      | 52,200,625     | 153,531      |
| 4      | Tahun 16–20  | 52,200,625     | 7,830,094      | 44,370,531     | 130,501      |
| 5      | Tahun 21–25  | 44,370,531     | 6,655,580      | 37,714,951     | 110,926      |
| ...    | ...          | ...            | ...            | ...            | ...          |
| ∞      | Tahun ∞      | < 1,000 GXR    | **AUTO STOP**  | < 1,000 GXR    | 0            |

### Distribusi Halving per Siklus:

- **70%** → Validator aktif (dibagi rata)
- **20%** → PoS Pool (delegator)
- **10%** → DEX Pool (likuiditas GXR/TON, GXR/POLYGON, dll)

### Keunggulan Sistem Baru:

1. **Deflasi Berkelanjutan**: Supply terus berkurang selamanya
2. **Scarcity Alami**: Semakin lama semakin langka
3. **Sustainable**: Tidak pernah habis sampai threshold minimum
4. **Predictable**: Semua perhitungan deterministik
5. **Self-Regulating**: Otomatis berhenti di batas minimum

---

## 4. VALIDATOR & DELEGATOR

### Validator:

- **Max Node**: 85 validator
- **Komisi Awal**: 5%–10% (dinamis)
- **Reward**: Bulanan dari halving + fee transaksi
- **Unstake Fee Delegator**: 0.5% (masuk ke validator)
- **Nonaktif >10 hari/bulan**: Tidak terima reward bulan itu
- **Bot Wajib**: Setiap validator harus menjalankan bot otomatis

### Delegator:

- **Reward**: Dari PoS Pool (20% dari halving bulanan)
- **Unstake Fee**: 0.5% ke validator
- **Minimum Stake**: 1 GXR
- **Pilih validator**: Berdasarkan performa dan komisii

---

## 5. SISTEM FEE TRANSAKSI

### A. Fee Transaksi Umum

Digunakan untuk semua transaksi biasa:

| Komponen               | Persentase |
| ---------------------- | ---------- |
| Validator              | 40%        |
| DEX Pool (Auto Refill) | 30%        |
| PoS Pool (Delegator)   | 30%        |

### B. Fee dari Aktivitas LP Komunitas (Farming)

Untuk pool komunitas buatan pengguna:

| Komponen               | Persentase |
| ---------------------- | ---------- |
| Validator              | 30%        |
| DEX Pool (Auto Refill) | 25%        |
| LP Komunitas (reward)  | 25%        |
| PoS Pool (Delegator)   | 20%        |

---

## 6. SISTEM LP: TIM & KOMUNITAS

### A. LP Resmi

- **Dibuat oleh tim GXR**: Pool utama
- **Likuiditas awal**: Dari alokasi 10% LP & Market
- **Dikelola oleh bot validator**: Auto refill, auto balancing
- **Menggunakan aturan fee default**: 40/30/30

### B. LP Komunitas

- **Siapa pun boleh buat LP**: GXR/TOKEN lain
- **Insentif berbasis fee farming**: 30/25/25/20
- **Dideteksi oleh bot validator**: Diberi reward otomatis
- **Tidak butuh smart contract**: Cukup whitelist address LP aktif

---

## 7. BOT VALIDATOR (WAJIB)

### Fungsi Bot:

- **Auto IBC Relayer**: Sinkronisasi antar-chain
- **Auto Rebalancing**: Rebalancing harga antar pool
- **Auto Reward Distribution**: Distribusi reward bulanan otomatis
- **Auto Refill DEX Pool**: Dari fee transaksi
- **Telegram Alert**: Uptime, pool imbalance, emergency
- **Health Monitoring**: Connection health, packet relaying

### Proteksi Bot:

- **Max swap harian**: 10.000 GXR/hari
- **Cooldown swap**: 30 menit jika lonjakan ekstrem
- **Emergency mode**: Jika harga GXR > $10
- **Auto recovery**: Reconnection otomatis pada network issues
- **Rate limiting**: Telegram alerts (max 10/menit)

### Komponen Bot:

1. **Reward Distributor**: Trigger halving distribution
2. **Rebalancer**: Cross-chain price balancing
3. **DEX Manager**: Pool management & refill
4. **IBC Relayer**: Cross-chain packet relaying
5. **Telegram Alert**: Real-time monitoring & alerts

---

## 8. IMPLEMENTASI TEKNIS

### Halving Module:

```go
// Konstanta sistem
const (
    MinimumSupplyThreshold = 1000 * 1e8 // 1,000 GXR
    HalvingReductionRate   = "0.15"     // 15%
    HalvingCycleDuration   = 5 * 365 * 24 * time.Hour // 5 tahun
)

// Fungsi utama halving
func CalculateMonthlyReward(currentSupply sdk.Coin) sdk.Coin {
    reductionRate := sdk.MustNewDecFromStr("0.15")
    cycleReduction := currentSupply.Amount.ToDec().Mul(reductionRate)
    monthlyAmount := cycleReduction.QuoRaw(60) // 60 bulan
    return sdk.NewCoin("ugen", monthlyAmount.TruncateInt())
}

func DistributeMonthlyRewards(ctx sdk.Context) error {
    // 1. Burn monthly reward dari total supply
    // 2. Mint reward untuk distribusi
    // 3. Distribute: 70% validator, 20% delegator, 10% DEX
    // 4. Update halving info
    // 5. Record distribution
}
```

### Bot Architecture:

```go
// Bot utama
type GXRBot struct {
    ibcRelayer        *IBCRelayer
    rewardDistributor *RewardDistributor
    dexManager        *DEXManager
    rebalancer        *Rebalancer
    telegramAlert     *TelegramAlert
}

// Setiap komponen berjalan independen dengan error handling
```

---

## 9. KEAMANAN & COMPLIANCE

### Fitur Keamanan:

- **Immutable Genesis**: Semua parameter dikunci dari awal
- **No Governance**: Tidak ada voting yang bisa mengubah aturan
- **Anti-Manipulation**: Bot protection untuk extreme price movements
- **Rate Limiting**: Semua operasi bot dengan rate limiting
- **Auto Recovery**: Automatic reconnection pada network issues

### Compliance:

- **Audit Trail**: Semua transaksi tercatat on-chain
- **Transparent**: Open source dan dapat diverifikasi
- **Predictable**: Semua perhitungan deterministik
- **Decentralized**: Tidak ada central authority

---

## 10. CIRI KHAS GXR

| Fitur                   | Penjelasan                                                                                                            |
| ----------------------- | --------------------------------------------------------------------------------------------------------------------- |
| **Deflasi Berkelanjutan** | Total supply terus berkurang setiap bulan                                                                           |
| **No Smart Contract**   | Sederhana, ringan, aman                                                                                              |
| **Dynamic Supply**      | Supply berubah berdasarkan distribusi halving                                                                        |
| **Auto Fee Refill**     | DEX Pool terisi otomatis dari fee                                                                                    |
| **Bot Mandatory**       | Setiap validator wajib menjalankan bot otomatis                                                                       |
| **Sustainable Halving** | Sistem halving yang sustainable dengan auto-stop                                                                     |
| **LP Komunitas**        | Komunitas bebas buat LP, reward otomatis                                                                             |
| **Immutable Genesis**   | Semua parameter dikunci dari awal, tidak ada governance                                                              |

---

## 11. MONITORING & ANALYTICS

### Query Commands:

```bash
# Cek total supply saat ini
gxrchaind query bank total --denom ugen

# Cek info halving
gxrchaind query halving info

# Cek distribusi history
gxrchaind query halving distributions

# Cek status validator
gxrchaind query staking validators
```

### Metrics Penting:

- **Current Total Supply**: Supply saat ini
- **Monthly Burn Rate**: Jumlah yang di-burn per bulan
- **Halving Cycle Progress**: Progress siklus saat ini
- **Validator Uptime**: Uptime validator
- **Bot Health**: Status kesehatan bot

---

## 12. ROADMAP TEKNIS

### Phase 1: Core Implementation ✅
- [x] Halving module dengan supply reduction
- [x] Bot validator otomatis
- [x] IBC integration
- [x] Telegram monitoring

### Phase 2: Enhancement 🚧
- [ ] Web dashboard untuk monitoring
- [ ] Advanced analytics
- [ ] Mobile app untuk delegators
- [ ] API improvements

### Phase 3: Ecosystem 🔮
- [ ] Cross-chain DEX integration
- [ ] Advanced rebalancing algorithms
- [ ] DeFi protocols integration
- [ ] NFT marketplace

---

**⚡ GXR: The Future of Deflationary Blockchain**

*Sistem halving yang benar-benar revolusioner - mengurangi supply secara langsung, bukan sekedar distribusi dari pool tetap.*

