# PingPulse 🚀

A lightweight monitoring solution that continuously checks your APIs and services, sending alerts when things go wrong. Built with Go backend and React TypeScript dashboard.

<img width="1428" height="942" alt="PingPulse Dashboard" src="https://github.com/user-attachments/assets/85df21be-f099-4328-9eb8-2cf88ca7fc4c" />

## ✨ Key Features

- 🔄 **Continuous Monitoring** - Check HTTP/HTTPS endpoints at custom intervals
- 🚨 **Smart Alerts** - Email, Slack, and webhook notifications with flexible conditions
- 🔐 **Authentication Support** - Bearer tokens, Basic auth, API keys, and OAuth2
- 📊 **Beautiful Dashboard** - Real-time React TypeScript interface
- 🏗️ **Project Organization** - Group monitors by projects with separate API keys
- 📦 **Zero Dependencies** - Single binary, no database required
- 🌐 **REST API** - Full programmatic access to all monitoring data

## Quick Start 🚀

### Option 1: Start Everything (Recommended)

```bash
# Clone and build
git clone https://github.com/vimoksh-5/pingpulse.git
cd pingpulse
go build -o pingpulse ./cmd/pingpulse

# Start both backend and dashboard
./start-dashboard.sh
```

### Option 2: Manual Setup

```bash
# Terminal 1: Start backend
./pingpulse

# Terminal 2: Start dashboard
cd dashboard && npm install && npm start
```

### 🌐 **Access Your Services**

- **📊 Dashboard**: http://localhost:3000 (Beautiful web interface)
- **🔧 API Backend**: http://localhost:8080 (REST API endpoints)
- **📋 Health Check**: http://localhost:8080/healthz

## Dashboard Features 📊

### **Real-time Monitoring Interface**

- **Statistics Overview**: Total projects, monitors, uptime percentage, average latency
- **Status Indicators**: Color-coded health status (🟢 Healthy, 🟡 Warning, 🔴 Critical)
- **Monitor Cards**: Individual service status with detailed metrics
- **Project Filtering**: Filter monitors by specific projects
- **Search Functionality**: Find monitors by name, URL, or project
- **Connection Status**: Real-time backend connection indicator

### **Monitor Information Display**

Each monitor shows:

- Service name and project
- Current status with color coding
- HTTP status code
- Response latency with performance indicators
- Uptime percentage
- Last check timestamp
- Service URL

### **Status System**

- 🟢 **Healthy**: 200-299 status codes, low latency
- 🟡 **Warning**: High latency (1-2s) or 4xx errors
- 🔴 **Critical**: 5xx errors, timeouts, or very high latency (>2s)
- ⚪ **Unknown**: Status cannot be determined

## 🔌 API Access

PingPulse provides a comprehensive REST API for programmatic access:

- **Base URL**: `http://localhost:8080`
- **Health Check**: `GET /healthz`
- **Statistics**: `GET /api/v1/statistics`
- **Projects**: `GET /api/v1/projects`
- **Monitors**: `GET /api/v1/monitors`
- **Results**: `GET /api/v1/results` (with advanced filtering)

📖 **Full API Documentation**: See [API_DOCUMENTATION.md](API_DOCUMENTATION.md) for complete endpoint details, examples, and response formats.

## ⚙️ Configuration

Create a `config.yaml` file to define your monitoring setup:

```yaml
default:
  interval: 60s
  timeout: 5s

projects:
  - name: "Production APIs"
    monitors:
      - name: "User Service Health"
        url: "https://api.company.com/users/health"
        interval: 30s
        alerts:
          - name: "Service Down"
            condition: "status != 200"
            channels: ["email", "slack"]
            recipients: ["ops@company.com"]
```

📖 **Complete Configuration Guide**: See [CONFIGURATION_GUIDE.md](CONFIGURATION_GUIDE.md) for detailed configuration options, authentication, alerts, and real-world examples.

## 🔐 Environment Variables

Set up your credentials in a `.env` file:

```bash
# Email alerts
GMAIL_USERNAME=your-email@gmail.com
GMAIL_APP_PASSWORD=your-app-password

# Slack alerts
SLACK_WEBHOOK_URL=https://hooks.slack.com/services/YOUR/SLACK/WEBHOOK

# API authentication
API_TOKEN=your-bearer-token
```

## 💻 Development

### Dashboard Development

```bash
cd dashboard
npm install
npm start  # Starts on http://localhost:3000
```

**Tech Stack**: React 18 + TypeScript, CSS Modules, Axios

## 📚 Documentation

- 📖 **[API Documentation](API_DOCUMENTATION.md)** - Complete REST API reference
- ⚙️ **[Configuration Guide](CONFIGURATION_GUIDE.md)** - Detailed setup and examples
- 🚀 **[Deployment Guide](DEPLOYMENT.md)** - Production deployment instructions
- 🔍 **[Troubleshooting](TROUBLESHOOTING.md)** - Common issues and solutions

## 🤝 Contributing

1. Fork the repository
2. Create feature branch (`git checkout -b feature/amazing-feature`)
3. Add tests for new features
4. Submit pull request

```bash
# Development setup
go mod tidy
go run ./cmd/pingpulse
```

## 📄 License

MIT License - see [LICENSE](LICENSE) file for details.

## 🌟 Support

- 🐛 **Issues**: [GitHub Issues](https://github.com/vimoksh-5/pingpulse/issues)
- 💬 **Discussions**: [GitHub Discussions](https://github.com/vimoksh-5/pingpulse/discussions)

---

**Made with ❤️ by [Vimoksh](https://github.com/vimoksh-5)**

_Keep your services healthy with PingPulse!_ 🚀
