# GXR (Gen X Raider) Blockchain Specification — FINAL (Juli 2025)

> ✅ Dirancang tanpa smart contract, anti-inflasi, berbasis PoS & IBC
> ✅ Fokus pada efisiensi, distribusi adil, dan desentralisasi otomatis

---

## 1. IDENTITAS DASAR

| Komponen           | Rincian                               |
| ------------------ | ------------------------------------- |
| **Nama Chain**     | Gen X Raider (GXR)                    |
| **Ticker Token**   | GXR                                   |
| **Denom Terkecil** | `gen` (1 GXR = 100,000,000 gen)       |
| **Total Supply**   | 85,000,000 GXR *(fixed, non-inflasi)* |
| **Desimal**        | 8                                     |
| **Smart Contract** | Tidak digunakan                       |
| **Konsensus**      | Proof-of-Stake (PoS)                  |
| **IBC Support**    | Aktif (GXR/TON, GXR/POLYGON, dll)     |
| **Max Validator**  | 85 node                               |
| **Waktu blok**     | 15 detik per blok                     |
---

## 2. TOKENOMICS GXR

### Total Supply: 85,000,000 GXR

| Alokasi             | Jumlah (GXR) | Persentase | Keterangan                                              |
| ------------------- | ------------ | ---------- | ------------------------------------------------------- |
| Airdrop & Farming   | 17,000,000   | 20%        | Distribusi awal via Telegram bot farming                |
| Developer Core      | 5,950,000    | 7%         | Vesting keras 5 tahun, 10% unlock tiap 6 bulan          |
| Tim Inti (3 orang)  | 5,950,000    | 7%         | 3% / 2% / 2%, soft vesting 3 tahun                      |
| LP & Market         | 8,500,000    | 10%        | Likuiditas awal (GXR/TON, GXR/POLYGON, dll)             |
| Grant (3–7 pihak)   | 8,500,000    | 10%        | Hibah proyek dan mitra kolaborasi                       |
| Pool Staking (PoS)  | 8,500,000    | 10%        | Reward untuk delegator aktif                            |
| Halving Fund        | 21,250,000   | 25%        | Untuk reward per 5 tahun                                |
| Cadangan/Ekspansi   | 8,500,000    | 10%        | Dana darurat dan pengembangan ekosistem                 |
| Validator Awal (30) | 850,000      | 1%         | 0.5% tahun 1 dan 0.5% tahun 2 jika aktif >20 hari/bulan |

---

## 3. SISTEM HALVING DINAMIS

### Siklus Halving: Tiap 5 Tahun Sekali

- Semua dana diambil dari Halving Fund (21,250,000 GXR)
- Tidak ada minting baru (fixed supply)

### Mekanisme Penurunan Bertahap:

| Halving Ke | Periode     | Dana (GXR) | Pengurangan (%) | Bulanan (±) |
| ---------- | ----------- | ---------- | --------------- | ----------- |
| 1          | Tahun 1–5   | 4,250,000  | —               | 70,833      |
| 2          | Tahun 6–10  | 3,612,500  | -15%            | 60,208      |
| 3          | Tahun 11–15 | 3,070,625  | -15%            | 51,177      |
| 4          | Tahun 16–20 | 2,610,032  | -15%            | 43,500      |
| 5          | Tahun 21–25 | 2,218,528  | -15%            | 36,975      |

### Distribusi Halving per Siklus:

- 70% → Validator aktif
- 20% → PoS Pool (delegator)
- 10% → DEX Pool (likuiditas GXR/TON, GXR/POLYGON, dll)
  - Tahun 1: 5% dibagi rata
  - Tahun 2: 5% dibagi habis, dinamis berdasarkan volume.
  - Tahun 3–5: tidak dibagikan

---

## 4. VALIDATOR & DELEGATOR

### Validator:

- Max Node: 85
- Komisi Awal: 5%–10% (dinamis)
- Reward: Bulanan dari halving + fee transaksi
- Unstake Fee Delegator: 0.5% (masuk ke validator)
- Nonaktif >10 hari/bulan → tidak terima reward bulan itu

### Delegator:

- Mendapat reward dari PoS Pool
- Pilih validator aktif

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

- Dibuat oleh tim GXR
- Likuiditas awal dari alokasi 10% LP & Market
- Dikelola oleh bot validator (auto refill, auto balancing)
- Menggunakan aturan fee default (40/30/30)

### B. LP Komunitas

- Siapa pun boleh buat LP GXR/TOKEN lain
- Insentif berbasis fee farming (30/25/25/20)
- Dideteksi oleh bot validator & diberi reward otomatis
- Tidak butuh smart contract, cukup whitelist address LP aktif

---

## 7. BOT VALIDATOR (WAJIB)

### Fungsi:

- Auto IBC Relayer (sinkronisasi antar-chain)
- Auto Rebalancing harga antar pool
- Auto distribusi reward
- Auto refill DEX Pool dari fee transaksi
- Telegram alert untuk uptime, pool imbalance, dsb
- Sekali klik deploy node + bot aktif otomatis. blokchain immutable tapi bos bisa aku update.

### Proteksi Bot:

- Max swap harian (misal 10.000 GXR/hari)
- Cooldown swap 30 menit jika lonjakan ekstrem terdeteksi
- Jika harga GXR > \$5–\$10: **Bot masuk mode nonaktif sementara 24 Jam (monitor-only)**
- Setelah cooldown, bot kembali aktif → langsung **jalankan rebalancing antar pool**
- Distribusi reward tetap berjalan normal secara terpisah dari rebalancing
- Tujuan: mencegah manipulasi harga dan proteksi sistem saat market ekstrem

---

## 8. CIRI KHAS GXR

| Fitur                   | Penjelasan                                                                                                            |
| ----------------------- | --------------------------------------------------------------------------------------------------------------------- |
| No Smart Contract       | Sederhana, ringan, aman                                                                                               |
| Fixed Supply            | Anti inflasi                                                                                                          |
| Auto Fee Refill         | DEX Pool terisi otomatis                                                                                              |
| Bot Desentral           | Validator jalankan node+otomatir bot relayer sekali klik semua otomatis aktif dan tidak bisa di rubah, misalnya ingin jalankan node saja, tidak bisa.                                                                               |
| Anti Dump Halving       | Reward bulanan selama 2 tahun, tidak langsung                                                                                        |
| Insentif Validator Awal | 1% bonus ke 30 node pertama di bagi selama 2 tahun (dengan uptime minimum)                                                                   |
| LP Komunitas            | Komunitas bebas buat LP, reward otomatis via fee farming                                                              |
| **Immutable Genesis**   | Semua parameter dikunci dari awal, tidak ada upgrade SC atau governance, hanya update via hard fork jika mutlak perlu |

---

