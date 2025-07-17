# GXR Validator Bot

**WAJIB**: Setiap validator GXR harus menjalankan bot ini!

## 🎯 Overview

GXR Validator Bot adalah komponen wajib yang menyediakan:
- 🔗 **IBC Relayer**: Auto relay GXR/TON, GXR/POLYGON  
- 💰 **Reward Distributor**: Distribusi halving bulanan otomatis
- 🌊 **DEX Manager**: Auto refill liquidity pools
- ⚖️ **Rebalancer**: Auto rebalancing antar chain
- 📱 **Telegram Alert**: Monitoring dan notifikasi

## 🤖 Bot Components

### 1. IBC Relayer
Menjalankan IBC relay untuk:
- TON blockchain integration
- Polygon blockchain integration  
- Cosmos ecosystem chains
- Auto packet relaying setiap 30 detik

### 2. Reward Distributor
Memantau dan memicu:
- Distribusi halving bulanan (otomatis)
- Distribusi reward validator
- Distribusi reward delegator
- Distribusi ke DEX pools

### 3. DEX Manager
Mengelola:
- Auto refill GXR/TON pool
- Auto refill GXR/POLYGON pool
- LP community pool monitoring
- Balance threshold management

### 4. Rebalancer
Melakukan:
- Inter-chain rebalancing
- Price monitoring (emergency mode)
- Daily swap limits (10,000 GXR)
- Cooldown periods (30 menit)

### 5. Telegram Alert
Mengirim notifikasi:
- Bot startup/shutdown
- Emergency mode alerts
- Distribution success/failure
- Pool imbalance warnings

## 🔧 Installation

### Build from Source

```bash
cd bot
go mod tidy
go build -o gxr-bot .
```

### Download Binary

```bash
# Download latest release
wget https://github.com/Crocodile-ark/gxrchaind/releases/latest/gxr-bot
chmod +x gxr-bot
```

## ⚙️ Configuration

### Create Config File

```bash
cp config/bot.yaml.example config/bot.yaml
```

### Configuration Options

```yaml
# Chain connection
chain_rpc: "tcp://localhost:26657"
chain_grpc: "localhost:9090"
chain_id: "gxr-1"

# Bot settings
log_level: "info"
check_interval: "30s"

# IBC settings
ibc_enabled: true
ibc_channels:
  - "channel-0"  # TON channel
  - "channel-1"  # Polygon channel

# DEX settings
max_swap_daily: "10000ugen"  # 10,000 GXR
swap_cooldown: "30m"
price_limit: "10000000"      # $10 emergency threshold

# Telegram settings
telegram_token: "YOUR_BOT_TOKEN"
telegram_chat_id: "YOUR_CHAT_ID"

# Safety
emergency_mode: false
```

## 🚀 Running the Bot

### Standalone Mode

```bash
./gxr-bot --config config/bot.yaml
```

### With Launcher (Recommended)

```bash
# Start both chain and bot together
../launcher/gxr-launcher
```

### Systemd Service

```bash
# Create service file
sudo tee /etc/systemd/system/gxr-bot.service > /dev/null <<EOF
[Unit]
Description=GXR Validator Bot
After=network.target
Requires=gxr-chain.service

[Service]
Type=simple
User=gxr
WorkingDirectory=/home/gxr/gxrchaind/bot
ExecStart=/home/gxr/gxrchaind/bot/gxr-bot --config config/bot.yaml
Restart=always
RestartSec=5
LimitNOFILE=4096

[Install]
WantedBy=multi-user.target
EOF

sudo systemctl enable gxr-bot
sudo systemctl start gxr-bot
```

## 📊 Monitoring

### Status Checks

```bash
# Check bot status
curl localhost:8080/status

# Check component health
curl localhost:8080/health
```

### Log Monitoring

```bash
# Follow bot logs
tail -f ~/.gxrchaind/logs/bot.log

# Check specific component
grep "IBC Relayer" ~/.gxrchaind/logs/bot.log
grep "DEX Manager" ~/.gxrchaind/logs/bot.log
```

### Telegram Setup

1. Create Telegram bot via @BotFather
2. Get bot token
3. Add bot to monitoring group
4. Get chat ID via @userinfobot
5. Update config dengan token dan chat ID

## 🛡️ Security Features

### Rate Limiting
- Max 10,000 GXR swap per hari
- Cooldown 30 menit per swap
- Emergency brake jika GXR > $10

### Access Control
- Bot menggunakan key terpisah dari validator
- Read-only access untuk monitoring
- Limited permissions untuk transactions

### Monitoring
- Real-time health checks
- Automatic restart on failure
- Telegram alerts untuk semua events

## 🚨 Emergency Mode

Bot otomatis masuk emergency mode jika:
- GXR price > $10 (configurable)
- Network congestion terdeteksi
- Validator tidak responsive

### Emergency Actions:
- Stop semua auto-swapping
- Send immediate Telegram alert
- Log semua activities
- Wait 24 jam sebelum resume

## 📈 Performance Tuning

### Resource Requirements

**Minimum:**
- CPU: 1 core
- RAM: 2GB  
- Storage: 10GB
- Network: 10 Mbps

**Recommended:**
- CPU: 2 cores
- RAM: 4GB
- Storage: 50GB  
- Network: 100 Mbps

### Optimization

```yaml
# High-performance config
check_interval: "10s"    # Faster checks
swap_cooldown: "15m"     # Shorter cooldown
max_swap_daily: "50000ugen"  # Higher limits
```

## 🧪 Testing

### Development Mode

```bash
# Run with test config
./gxr-bot --config config/test.yaml --log-level debug
```

### Component Testing

```bash
# Test IBC relayer only
./gxr-bot --components ibc

# Test DEX manager only  
./gxr-bot --components dex

# Test telegram alerts
./gxr-bot --test-telegram
```

## 📝 Troubleshooting

### Common Issues

**Bot tidak connect ke chain:**
```bash
# Check chain status
gxrchaind status

# Check RPC endpoint
curl localhost:26657/status
```

**IBC relayer error:**
```bash
# Check channel status
gxrchaind q ibc channel channels

# Check packet acknowledgments
gxrchaind q ibc channel packets
```

**Telegram alerts tidak kirim:**
```bash
# Test bot token
curl https://api.telegram.org/bot<TOKEN>/getMe

# Test chat ID
curl https://api.telegram.org/bot<TOKEN>/sendMessage?chat_id=<CHAT_ID>&text=test
```

### Error Codes

| Code | Error | Solution |
|------|-------|----------|
| 1001 | Chain connection failed | Check RPC endpoint |
| 1002 | IBC channel not found | Setup IBC channels |
| 1003 | Insufficient balance | Check bot wallet |
| 1004 | Telegram API error | Check token/chat ID |
| 1005 | Emergency mode active | Wait or manual override |

## 🔄 Updates

### Auto Updates

Bot checks for updates setiap 24 jam:
```bash
# Enable auto-update
./gxr-bot --auto-update
```

### Manual Updates

```bash
# Download latest
wget https://github.com/Crocodile-ark/gxrchaind/releases/latest/gxr-bot

# Backup current
mv gxr-bot gxr-bot.backup

# Replace and restart
chmod +x gxr-bot
systemctl restart gxr-bot
```

## 📋 Best Practices

### Validator Setup
1. Run bot on sama server dengan validator
2. Use firewall untuk proteksi
3. Monitor via Telegram 24/7
4. Regular backup bot config
5. Test emergency procedures

### Maintenance
1. Update bot minimum 1x per bulan
2. Monitor performance metrics
3. Check log files weekly
4. Test alert systems regularly
5. Keep backup server ready

## 🆘 Support

### Emergency Contact
- **Telegram**: @gxr_foundation
- **Email**: gxrfoundation@gmail.com

---

**⚠️ PENTING**: Bot ini WAJIB untuk semua validator GXR. Validator yang tidak menjalankan bot akan di-slash!
