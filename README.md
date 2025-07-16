# GXR (Gen X Raider) Blockchain
![banner](docs/static/img/banner.jpg)

> ✅ Anti-inflasi, berbasis PoS & IBC, tanpa smart contract
> ✅ Fokus pada efisiensi, distribusi adil, dan desentralisasi otomatis

## 🎯 Tentang GXR

GXR adalah blockchain Proof-of-Stake yang dirancang khusus untuk:
- **Fixed Supply**: 85,000,000 GXR (tanpa inflasi)
- **Auto Fee Distribution**: Pembagian fee otomatis sesuai spesifikasi
- **Halving System**: Distribusi reward 5 tahunan dengan pengurangan 15%
- **IBC Support**: Bridge otomatis dengan TON, Polygon, dan chain lain
- **Validator Bot**: Setiap validator wajib menjalankan bot otomatis

## 📁 Struktur Folder

```
gxrchaind/
├── chain/
|     ├── app/                   # Aplikasi blockchain utama
|     │   ├── app.go             # Konfigurasi aplikasi
|     │   ├── encoding.go        # Encoding setup
|     │   ├── genesis.go         # Genesis helpers
|     │   └── gxr_genesis.go     # Alokasi supply GXR
|     ├── cmd/gxrchaind/         # CLI daemon
|     │   └── cmd/
|     │       ├── root.go        # Root command
|     │       └── genaccounts.go # Genesis account command
|     ├── proto/                 # Protobuf definitions
|     |        └── gxr/
|     |            ├── halving/
|     |            └── feerouter/
|     └── x/                     # Modul kustom
|         ├── halving/           # Modul distribusi reward 5 tahunan
|         ├── feerouter/         # Modul routing fee otomatis
|         └── proto/             # Protobuf definitions
|              └── gxr/
|                   ├── halving/
|                   └── feerouter/
|
├── bot/                         # Bot validator (WAJIB)
│   ├── main.go                  # Bot utama
│   ├── ibc_relayer.go           # IBC relayer otomatis
│   ├── reward_distributor.go    # Distribusi reward
│   ├── dex_manager.go           # Auto refill DEX
│   ├── rebalancer.go            # Auto rebalancing
│   └── telegram_alert.go        # Alert Telegram
└── launcher/                    # Launcher untuk chain + bot
    └── main.go                  # Launcher utama

```

## 🔧 Build & Installation

### Prerequisites

```bash
# Install Go 1.21+
# Install make
# Install git
```

### Build Commands

```bash
# Clone repository
git clone https://github.com/Crocodile-ark/gxrchaind
cd gxrchaind

# Build blockchain daemon
make build

# Build bot
cd bot
go build -o gxr-bot .
cd ..

# Build launcher
cd launcher  
go build -o gxr-launcher .
cd ..
```

### Quick Start (Development)

```bash
# 1. Initialize node
./build/gxrchaind init mynode --chain-id gxr-1

# 2. Create validator key
./build/gxrchaind keys add validator

# 3. Add genesis account with initial tokens
./build/gxrchaind add-genesis-account validator 1000000000000000ugen

# 4. Create genesis transaction
./build/gxrchaind gentx validator 50000000000000ugen --chain-id gxr-1

# 5. Collect genesis transactions
./build/gxrchaind collect-gentxs

# 6. Start with launcher (recommended)
./launcher/gxr-launcher

# Atau start manual:
# Terminal 1: ./build/gxrchaind start
# Terminal 2: ./bot/gxr-bot
```

## 🚀 Production Deployment

### Validator Setup

```bash
# 1. Setup validator node
./build/gxrchaind init validator-01 --chain-id gxr-mainnet

# 2. Download genesis file
wget https://raw.githubusercontent.com/Crocodile-ark/gxrchaind/main/networks/mainnet/genesis.json
cp genesis.json ~/.gxrchaind/config/

# 3. Configure persistent peers
vim ~/.gxrchaind/config/config.toml
# Set persistent_peers = "..."

# 4. Configure bot
cp bot/config/bot.yaml.example bot/config/bot.yaml
vim bot/config/bot.yaml
# Set Telegram token, channels, etc.

# 5. Start with launcher
./launcher/gxr-launcher --auto-restart
```

### Systemd Service (Recommended)

```bash
# Create service file
sudo tee /etc/systemd/system/gxr.service > /dev/null <<EOF
[Unit]
Description=GXR Blockchain with Bot
After=network.target

[Service]
Type=simple
User=gxr
WorkingDirectory=/home/gxr/gxrchaind
ExecStart=/home/gxr/gxrchaind/launcher/gxr-launcher
Restart=always
RestartSec=3
LimitNOFILE=4096

[Install]
WantedBy=multi-user.target
EOF

# Enable and start
sudo systemctl enable gxr
sudo systemctl start gxr
sudo systemctl status gxr
```

## 💰 Tokenomics

| Alokasi                | Jumlah GXR  | Persentase | Vesting              |
|------------------------|-------------|------------|----------------------|
| Airdrop & Farming      | 17,000,000  | 20%        | Telegram bot farming |
| Developer Core         | 5,950,000   | 7%         | 5 tahun (hard)       |
| Tim Inti (3 orang)     | 5,950,000   | 7%         | 3 tahun (soft)       |
| LP & Market            | 8,500,000   | 10%        | Likuiditas awal      |
| Grant Kolaborasi       | 8,500,000   | 10%        | Hibah proyek         |
| Pool Staking (PoS)     | 8,500,000   | 10%        | Reward delegator     |
| **Halving Fund**       | 21,250,000  | 25%        | **Reward 5 tahunan** |
| Cadangan/Ekspansi      | 8,500,000   | 10%        | Dana darurat         |
| Validator Awal (30)    | 850,000     | 1%         | Bonus 2 tahun        |

**Total Supply: 85,000,000 GXR (Fixed, Anti-Inflasi)**

## 🔄 Sistem Halving

### Contoh Siklus 5 Tahunan

| Halving | Periode      | Dana GXR    | Pengurangan | Bulanan |
|---------|--------------|-------------|-------------|---------|
| 1       | Tahun 1–5    | 4,250,000   | —           | 70,833  |
| 2       | Tahun 6–10   | 3,612,500   | -15%        | 60,208  |
| 3       | Tahun 11–15  | 3,070,625   | -15%        | 51,177  |
| 4       | Tahun 16–20  | 2,610,032   | -15%        | 43,500  |
| 5       | Tahun 21–25  | 2,218,528   | -15%        | 36,975  |

### Distribusi Bulanan

- **70%** → Validator aktif
- **20%** → PoS Pool (delegator)  
- **10%** → DEX Pool (GXR/TON, GXR/POLYGON)

## 💸 Sistem Fee

### Transaksi Umum (40/30/30)
- **40%** → Validator
- **30%** → DEX Pool (auto refill)
- **30%** → PoS Pool (delegator)

### LP Farming (30/25/25/20)
- **30%** → Validator
- **25%** → DEX Pool
- **25%** → LP Komunitas (reward)
- **20%** → PoS Pool

## 🤖 Bot Validator

**WAJIB**: Setiap validator harus menjalankan bot!

### Fitur Bot:
- ✅ **Auto IBC Relayer** (GXR/TON, GXR/POLYGON)
- ✅ **Auto Reward Distributor** (bulanan)
- ✅ **Auto DEX Refill** (dari fee)
- ✅ **Auto Rebalancing** (antar pool)
- ✅ **Telegram Alert** (monitoring)

### Proteksi Bot:
- Max swap 10,000 GXR/hari
- Cooldown 30 menit per swap
- Emergency mode jika GXR > $10
- Auto restart jika crash

## 🌉 IBC & Bridge

### Supported Chains:
- **TON** (Telegram Open Network)
- **Polygon** (MATIC)
- **Cosmos Hub** (ATOM)
- *Lebih banyak akan ditambahkan*

### LP Resmi:
- `GXR/TON` - Pool utama
- `GXR/ATOM` - Pool utama
- `GXR/POLYGON` - Pool sekunder

### LP Komunitas:
- Siapa saja bisa buat pool
- Auto reward via fee farming
- Whitelist otomatis oleh bot

## 📊 Monitoring

### CLI Commands:

```bash
# Status chain
./build/gxrchaind status

# Query halving info
./build/gxrchaind q halving halving-info

# Query fee stats  
./build/gxrchaind q feerouter fee-stats

# Query LP pools
./build/gxrchaind q feerouter lp-pools
```

### Bot Status:

```bash
# Bot logs
tail -f ~/.gxrchaind/logs/bot.log

# Launcher status
./launcher/gxr-launcher status
```

## 🔒 Security

### Validator Security:
- Run on dedicated server
- Firewall only necessary ports
- Regular backups of validator keys
- Monitor via Telegram alerts

### Bot Security:
- Bot keys separate from validator
- Rate limiting & cooldowns
- Emergency stop mechanisms
- Secure Telegram integration

## 🆘 Support

### Documentation:
- [Halving Module README](x/halving/README.md)
- [FeeRouter Module README](x/feerouter/README.md)
- [Bot README](bot/README.md)
- [Launcher README](launcher/README.md)

### Community:
- **Telegram**: @GXRBlockchain
- **GitHub**: https://github.com/Crocodile-ark/gxrchaind
- **Explorer**: https://explorer.gxr.network

---

## ⚖️ License

Apache License 2.0

---

**🎯 GXR: Fixed Supply, Auto Distribution, IBC Ready** 🚀

