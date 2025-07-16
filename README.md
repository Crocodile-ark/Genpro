# GXR (Gen X Raider) Blockchain

> ✅ **Blockchain anti-inflasi berbasis PoS dengan sistem halving dinamis**
> 
> Blockchain Cosmos SDK yang dirancang untuk efisiensi, distribusi adil, dan desentralisasi otomatis dengan sistem halving yang mengurangi total supply secara bertahap.

---

## 🌟 Fitur Utama

### 🔥 Sistem Halving Revolusioner
- **Pengurangan Supply**: Setiap 5 tahun, total supply berkurang 15%
- **Distribusi Bulanan**: Reward terdistribusi merata selama 5 tahun per siklus
- **Auto-Stop**: Halving berhenti otomatis jika total supply < 1.000 GXR
- **Deflasi Berkelanjutan**: Supply terus berkurang, menciptakan scarcity alami

### 🤖 Bot Validator Otomatis
- **IBC Relayer**: Sinkronisasi otomatis antar-chain
- **Auto Rebalancing**: Rebalancing harga otomatis antar pool
- **Reward Distribution**: Distribusi reward bulanan otomatis
- **Telegram Alerts**: Monitoring real-time via Telegram

### 💎 Tokenomics Anti-Inflasi
- **Total Supply**: 85,000,000 GXR (fixed, tidak bertambah)
- **Konsensus**: Proof-of-Stake (PoS) dengan 85 validator maksimal
- **IBC Support**: Koneksi cross-chain ke TON, Polygon, dan lainnya
- **No Smart Contract**: Sederhana, ringan, aman

---

## 📊 Sistem Halving Terbaru

### Cara Kerja:
1. **Setiap 5 tahun**: Sistem menghitung 15% dari total supply saat ini
2. **Distribusi bulanan**: 15% tersebut dibagi menjadi 60 distribusi bulanan
3. **Burn & Mint**: Monthly reward di-burn dari supply, kemudian di-mint untuk distribusi
4. **Siklus berkelanjutan**: Setelah 5 tahun, sisa supply menjadi 100% baru untuk siklus berikutnya

### Contoh Proyeksi:
```
Siklus 1: 85,000,000 GXR → 72,250,000 GXR (15% terdistribusi)
Siklus 2: 72,250,000 GXR → 61,412,500 GXR (15% terdistribusi)
Siklus 3: 61,412,500 GXR → 52,200,625 GXR (15% terdistribusi)
...dan seterusnya hingga supply < 1,000 GXR
```

### Distribusi Reward:
- **70%** → Validator aktif (dibagi rata)
- **20%** → PoS Pool untuk delegator
- **10%** → DEX Pool (likuiditas GXR/TON, GXR/POLYGON)

---

## 🚀 Quick Start

### 1. Setup Node
```bash
# Clone repository
git clone https://github.com/your-org/gxr-blockchain.git
cd gxr-blockchain

# Build chain
cd chain/
make build

# Initialize node
./build/gxrchaind init your-node-name --chain-id gxr-1

# Start node
./build/gxrchaind start
```

### 2. Setup Bot Validator
```bash
# Build bot
cd bot/
go build -o gxr-bot

# Create config file
cat > config/bot.yaml << EOF
chain_rpc: "tcp://localhost:26657"
chain_grpc: "localhost:9090"
chain_id: "gxr-1"
telegram_token: "YOUR_TELEGRAM_BOT_TOKEN"
telegram_chat_id: "YOUR_CHAT_ID"
ibc_enabled: true
ibc_channels: ["channel-0", "channel-1"]
EOF

# Start bot
./gxr-bot --config config/bot.yaml
```

### 3. Setup Validator
```bash
# Create validator
gxrchaind tx staking create-validator \
  --amount=1000000ugen \
  --pubkey=$(gxrchaind tendermint show-validator) \
  --moniker="your-validator" \
  --chain-id=gxr-1 \
  --commission-rate="0.05" \
  --from=your-wallet

# Start validator bot (wajib untuk validator)
./gxr-bot --config config/validator.yaml
```

---

## 🔧 Struktur Proyek

```
gxr-blockchain/
├── chain/                    # Blockchain core
│   ├── x/halving/           # Halving module (sistem utama)
│   ├── x/bank/              # Bank module
│   ├── x/staking/           # Staking module
│   ├── app/                 # Application setup
│   └── cmd/                 # CLI commands
├── bot/                     # Validator bot
│   ├── main.go             # Bot utama
│   ├── reward_distributor.go # Distribusi reward
│   ├── rebalancer.go       # Rebalancing otomatis
│   ├── ibc_relayer.go      # IBC relaying
│   └── telegram_alert.go   # Telegram alerts
├── launcher/               # Node launcher
└── docs/                   # Dokumentasi
```

---

## 📈 Monitoring & Status

### Query Commands:
```bash
# Cek total supply saat ini
gxrchaind query bank total --denom ugen

# Cek info halving
gxrchaind query halving info

# Cek status validator
gxrchaind query staking validators

# Cek distribusi history
gxrchaind query halving distributions
```

### Bot Status:
```bash
# Cek status bot
curl http://localhost:8080/status

# Cek status IBC
curl http://localhost:8080/ibc/status

# Cek status rebalancer
curl http://localhost:8080/rebalancer/status
```

---

## 🔐 Keamanan & Compliance

### Fitur Keamanan:
- **Immutable Genesis**: Parameter dikunci dari awal
- **No Governance**: Tidak ada voting yang bisa mengubah aturan
- **Anti-Manipulation**: Bot protection untuk extreme price movements
- **Rate Limiting**: Telegram alerts dengan rate limiting
- **Auto Recovery**: Automatic reconnection pada network issues

### Compliance:
- **Audit Trail**: Semua transaksi tercatat on-chain
- **Transparent**: Open source dan dapat diverifikasi
- **Predictable**: Semua perhitungan deterministik
- **Decentralized**: Tidak ada central authority

---

## 🛠️ Development

### Build Requirements:
- Go 1.21+
- Make
- Git

### Development Setup:
```bash
# Clone dan setup
git clone https://github.com/your-org/gxr-blockchain.git
cd gxr-blockchain

# Install dependencies
cd chain && go mod tidy
cd ../bot && go mod tidy

# Run tests
make test

# Run linter
make lint
```

### Testing:
```bash
# Unit tests
make test-unit

# Integration tests
make test-integration

# E2E tests
make test-e2e

# Coverage
make coverage
```

---

## 🤝 Contributing

1. Fork repository
2. Create feature branch: `git checkout -b feature/new-feature`
3. Commit changes: `git commit -m 'Add new feature'`
4. Push branch: `git push origin feature/new-feature`
5. Create Pull Request

### Code Style:
- Follow Go conventions
- Add tests for new features
- Update documentation
- Run linter before commit

---

## 📞 Support

### Telegram:
- **Official Group**: [t.me/gxr_blockchain](https://t.me/gxr_blockchain)
- **Developer Chat**: [t.me/gxr_devs](https://t.me/gxr_devs)
- **Validator Support**: [t.me/gxr_validators](https://t.me/gxr_validators)

### Resources:
- **Documentation**: [docs.gxr.blockchain](https://docs.gxr.blockchain)
- **Explorer**: [explorer.gxr.blockchain](https://explorer.gxr.blockchain)
- **API Reference**: [api.gxr.blockchain](https://api.gxr.blockchain)

---

## 📄 License

MIT License - see [LICENSE](LICENSE) file for details.

---

## 🎯 Roadmap

### Phase 1: Foundation ✅
- [x] Core blockchain dengan halving module
- [x] Bot validator otomatis
- [x] IBC integration
- [x] Telegram monitoring

### Phase 2: Expansion 🚧
- [ ] DEX integration (GXR/TON, GXR/POLYGON)
- [ ] Web dashboard untuk monitoring
- [ ] Mobile app untuk delegators
- [ ] Advanced analytics

### Phase 3: Ecosystem 🔮
- [ ] Cross-chain bridges
- [ ] DeFi protocols integration
- [ ] NFT marketplace
- [ ] Gaming integrations

---

**⚡ GXR: The Future of Deflationary Blockchain**

*Dirancang untuk masa depan, dibangun untuk bertahan.*

