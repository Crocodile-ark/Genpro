# GXR Launcher

Launcher to run the GXR blockchain and validator bot together.

## 🎯 Overview

GXR Launcher offers:

- **Unified Management**: Run chain + bot with one command
- **Auto Restart**: Automatically restarts if a process crashes
- **Proper Sequencing**: Starts the chain first, waits until ready, then starts the bot
- **Graceful Shutdown**: Stops the bot first, then the chain
- **Centralized Logging**: Logs from both processes in a single stream

## ⚙️ Installation

### Build from Source

```bash
cd launcher
go mod tidy
go build -o gxr-launcher .
```

### Download Binary

```bash
wget https://github.com/Crocodile-ark/gxrchaind/releases/latest/gxr-launcher
chmod +x gxr-launcher
```

## 🚀 Quick Start

### Basic Usage

```bash
# Start with default settings
./gxr-launcher

# Start with custom binaries
./gxr-launcher --chain-binary ./build/gxrchaind --bot-binary ./bot/gxr-bot

# Start with custom configs
./gxr-launcher --chain-config ~/.gxrchaind/config/config.toml --bot-config ./bot/config/bot.yaml
```

### Advanced Usage

```bash
# Production mode with auto-restart
./gxr-launcher --auto-restart --chain-home /opt/gxr/data

# Development mode (no auto-restart)
./gxr-launcher --auto-restart=false

# Custom restart delay
./gxr-launcher --restart-delay 10s
```

## 📋 Command Line Options

```bash
Usage:
  gxr-launcher [flags]

Flags:
      --auto-restart           Automatically restart failed processes (default true)
      --bot-binary string      Path to gxr-bot binary
      --bot-config string      Bot configuration file
      --chain-binary string    Path to gxrchaind binary
      --chain-config string    Chain configuration file
      --chain-home string      Chain home directory
  -h, --help                   help for gxr-launcher
  -v, --version                version for gxr-launcher

Commands:
  status      Show status of chain and bot processes
  help        Help about any command
```

## 🔧 Configuration

### Default Paths

```bash
Chain Binary:  ./build/gxrchaind
Bot Binary:    ./bot/gxr-bot
Chain Home:    $HOME/.gxrchaind
Auto Restart:  true
Restart Delay: 5 seconds
```

### Environment Variables

```bash
# Override defaults via environment
export GXR_CHAIN_BINARY="/usr/local/bin/gxrchaind"
export GXR_BOT_BINARY="/usr/local/bin/gxr-bot"
export GXR_CHAIN_HOME="/opt/gxr/data"
```

## 🔄 Process Management

### Startup Sequence

1. **Chain Start**: Run `gxrchaind start`
2. **Wait**: Wait 10 seconds for chain initialization
3. **Bot Start**: Run `gxr-bot` with config
4. **Monitor**: Monitor both processes in real time

### Shutdown Sequence

1. **Bot Stop**: Send SIGTERM to bot (gracefully)
2. **Wait**: Wait for bot to clean up
3. **Chain Stop**: Send SIGTERM to chain (gracefully)
4. **Wait**: Wait for all processes to stop

### Auto Restart

If a process crashes:

1. Log the error and reason
2. Wait for the restart delay (default 5 seconds)
3. Restart the failed process
4. Continue monitoring

## 📊 Monitoring

### Status Check

```bash
# Check running processes
./gxr-launcher status

# Output example:
Chain: ✅ Running (PID: 12345)
Bot:   ✅ Running (PID: 12346)
Auto-restart: ✅ Enabled
```

### Log Output

```bash
# Launcher prefix log lines:
[CHAIN] 2024-01-01T10:00:00Z INF Chain started successfully
[BOT]   2024-01-01T10:00:10Z INF Bot initialized successfully

# Follow logs in real time
./gxr-launcher | tee launcher.log
```

## 🛡️ Production Deployment

### Systemd Service

```bash
sudo tee /etc/systemd/system/gxr.service > /dev/null <<EOF
[Unit]
Description=GXR Blockchain Launcher
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

sudo systemctl enable gxr
sudo systemctl start gxr
sudo systemctl status gxr
```

### Docker Compose

```yaml
version: '3.8'
services:
  gxr:
    image: gxr/blockchain:latest
    command: ./launcher/gxr-launcher
    ports:
      - "26657:26657"
      - "1317:1317"
      - "9090:9090"
    volumes:
      - gxr_data:/home/gxr/.gxrchaind
    restart: unless-stopped

volumes:
  gxr_data:
```

### Kubernetes Deployment

```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: gxr-blockchain
spec:
  replicas: 1
  selector:
    matchLabels:
      app: gxr-blockchain
  template:
    metadata:
      labels:
        app: gxr-blockchain
    spec:
      containers:
      - name: gxr
        image: gxr/blockchain:latest
        command: ["./launcher/gxr-launcher"]
        ports:
        - containerPort: 26657
        - containerPort: 1317
        - containerPort: 9090
        volumeMounts:
        - name: gxr-data
          mountPath: /home/gxr/.gxrchaind
      volumes:
      - name: gxr-data
        persistentVolumeClaim:
          claimName: gxr-data-pvc
```

## 🚨 Troubleshooting

### Common Issues

**Chain fails to start:**

```bash
# Check binary path
ls -la ./build/gxrchaind

# Check permissions
chmod +x ./build/gxrchaind

# Check home directory
ls -la ~/.gxrchaind/
```

**Bot fails to start:**

```bash
# Check bot binary
ls -la ./bot/gxr-bot

# Check bot config
cat ./bot/config/bot.yaml

# Test bot standalone
./bot/gxr-bot --config ./bot/config/bot.yaml
```

**Auto-restart not working:**

```bash
# Check launcher config
./gxr-launcher --help

# Test with longer delay
./gxr-launcher --restart-delay 30s

# Disable for debugging
./gxr-launcher --auto-restart=false
```

### Debug Mode

```bash
# Enable verbose logging
./gxr-launcher --log-level debug

# Check process status
ps aux | grep gxr

# Check system resources
top -p $(pgrep gxrchaind)
top -p $(pgrep gxr-bot)
```

## 🔧 Development

### Building

```bash
# Build launcher
go build -o gxr-launcher .

# Build with version info
go build -ldflags "-X main.Version=v1.0.0" -o gxr-launcher .

# Cross-compile for different platforms
GOOS=linux GOARCH=amd64 go build -o gxr-launcher-linux .
GOOS=darwin GOARCH=amd64 go build -o gxr-launcher-darwin .
```

### Testing

```bash
# Unit tests
go test ./...

# Integration tests with mock processes
go test -tags integration ./...

# Test with real binaries
./gxr-launcher --chain-binary ./test/mock-chain --bot-binary ./test/mock-bot
```

## 📈 Performance

### Resource Usage

**Launcher overhead:**

- CPU: <1%
- RAM: \~10MB
- Network: Negligible

**Total system requirements:**

- CPU: 4+ cores (chain + bot + launcher)
- RAM: 8GB+ (chain + bot + launcher)
- Storage: 100GB+ (chain data)
- Network: 100Mbps+ (P2P + IBC)

### Optimization

```bash
# Reduce restart delay for faster recovery
./gxr-launcher --restart-delay 1s

# Custom process priorities
nice -n -5 ./gxr-launcher  # Higher priority
```

## 🛠️ Configuration Examples

### Development

```bash
./gxr-launcher \
  --chain-binary ./build/gxrchaind \
  --bot-binary ./bot/gxr-bot \
  --chain-home ~/.gxrchaind-dev \
  --bot-config ./bot/config/dev.yaml \
  --auto-restart=false
```

### Testnet

```bash
./gxr-launcher \
  --chain-binary /usr/local/bin/gxrchaind \
  --bot-binary /usr/local/bin/gxr-bot \
  --chain-home /opt/gxr/testnet \
  --chain-config /opt/gxr/testnet/config/config.toml \
  --bot-config /opt/gxr/testnet/bot.yaml
```

### Mainnet

```bash
./gxr-launcher \
  --chain-binary /usr/local/bin/gxrchaind \
  --bot-binary /usr/local/bin/gxr-bot \
  --chain-home /opt/gxr/mainnet \
  --chain-config /opt/gxr/mainnet/config/config.toml \
  --bot-config /opt/gxr/mainnet/bot.yaml \
  --auto-restart=true \
  --restart-delay=10s
```

## 📚 References

- [GXR Chain README](../README.md)
- [Bot README](../bot/README.md)
- [Systemd Documentation](https://www.freedesktop.org/software/systemd/man/systemd.service.html)
- [Process Management Best Practices](https://12factor.net/processes)

---

**💡 Tip**: Use the launcher for production deployment, run manually only for development!

