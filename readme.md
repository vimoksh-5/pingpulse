# PingPulse 🚀

A lightweight, configurable monitoring service written in Go that tracks the health and performance of your web services with real-time alerts.

## Features ✨

- **HTTP/HTTPS Monitoring** - Monitor any web endpoint with customizable intervals
- **Multi-Channel Alerts** - Email, Slack, and Webhook notifications
- **Flexible Authentication** - Support for Bearer tokens, Basic auth, API keys, and OAuth2
- **Advanced Monitoring** - Response time tracking, status code validation, and JSON response validation
- **Project-Based Organization** - Group monitors by projects with individual API keys
- **Zero Dependencies** - Single binary deployment with no external database required
- **Free Alerting Options** - Built-in support for Gmail SMTP, Slack webhooks, and custom webhooks

## Quick Start 🚀

### 1. Installation

```bash
# Clone the repository
git clone https://github.com/vimoksh-5/pingpulse.git
cd pingpulse

# Build the binary
go build -o pingpulse ./cmd/pingpulse

# Or run directly
go run ./cmd/pingpulse
```

### 2. Configuration

Create a `config.yaml` file:

```yaml
default:
  interval: 60s
  timeout: 5s
  maxRetries: 3
  jitter: 5s
  userAgent: "PingPulse/1.0"

global_alerts:
  email:
    smtp_host: "smtp.gmail.com"
    smtp_port: 587
    username: "${GMAIL_USERNAME}"
    password: "${GMAIL_APP_PASSWORD}"
    from: "${GMAIL_USERNAME}"
    use_tls: true
  slack:
    webhook_url: "${SLACK_WEBHOOK_URL}"
    channel: "#alerts"
    username: "PingPulse"

projects:
  - name: "My Project"
    apiKey: "project-123"
    monitors:
      - name: "Website Health"
        url: "https://example.com"
        expected_status: 200
        expected_latency: 2s
        alerts:
          - name: "Website Down"
            condition: "status != 200"
            channels: ["email", "slack"]
            recipients: ["admin@example.com"]
            cooldown: 5m
```

### 3. Environment Variables

Create a `.env` file:

```bash
# Gmail SMTP Configuration
GMAIL_USERNAME=your-email@gmail.com
GMAIL_APP_PASSWORD=your-app-password

# Slack Webhook Configuration
SLACK_WEBHOOK_URL=https://hooks.slack.com/services/YOUR/SLACK/WEBHOOK

# Optional: Custom Webhook
WEBHOOK_URL=https://webhook.site/your-unique-id
```

### 4. Run PingPulse

```bash
./pingpulse
```

## Configuration Reference 📖

### Monitor Configuration

```yaml
monitors:
  - name: "API Monitor"
    url: "https://api.example.com/health"
    interval: 30s          # How often to check
    timeout: 10s           # Request timeout
    maxRetries: 3          # Retry attempts on failure
    method: "GET"          # HTTP method
    expected_status: 200   # Expected HTTP status
    expected_latency: 5s   # Maximum acceptable latency
    
    # Custom headers
    headers:
      Accept: "application/json"
      User-Agent: "PingPulse Monitor"
    
    # Request body (for POST/PUT)
    body:
      key: "value"
    
    # Authentication
    auth:
      type: "bearer"       # bearer, basic, oauth2, apikey
      token: "${API_TOKEN}"
    
    # SSL and redirect options
    ssl_verify: true
    follow_redirects: true
    
    # JSON response validation
    body_validation:
      path: "$.status"     # JSONPath expression
      expected: "healthy"  # Expected value
```

### Alert Conditions

PingPulse supports various alert conditions:

- `status != 200` - HTTP status code checks
- `latency > 2s` - Response time thresholds
- `latency > 3s` - Custom latency limits
- `latency > 5s` - High latency alerts
- `ok == false` - General health checks

### Authentication Types

#### Bearer Token
```yaml
auth:
  type: "bearer"
  token: "${BEARER_TOKEN}"
```

#### Basic Authentication
```yaml
auth:
  type: "basic"
  username: "${USERNAME}"
  password: "${PASSWORD}"
```

#### API Key
```yaml
auth:
  type: "apikey"
  api_key: "${API_KEY}"
  api_key_header: "X-API-Key"  # or "Authorization"
```

#### OAuth2
```yaml
auth:
  type: "oauth2"
  client_id: "${CLIENT_ID}"
  client_secret: "${CLIENT_SECRET}"
  token_url: "https://auth.example.com/oauth/token"
```

## Alert Channels 📢

### Email Alerts

Configure Gmail SMTP for email notifications:

```yaml
global_alerts:
  email:
    smtp_host: "smtp.gmail.com"
    smtp_port: 587
    username: "${GMAIL_USERNAME}"
    password: "${GMAIL_APP_PASSWORD}"  # Use App Password, not regular password
    from: "${GMAIL_USERNAME}"
    use_tls: true
```

**Setup Gmail App Password:**
1. Enable 2-Factor Authentication on your Gmail account
2. Generate an App Password: Google Account → Security → App passwords
3. Use the App Password in your `.env` file

### Slack Alerts

Set up Slack webhook notifications:

```yaml
global_alerts:
  slack:
    webhook_url: "${SLACK_WEBHOOK_URL}"
    channel: "#alerts"
    username: "PingPulse"
```

**Setup Slack Webhook:**
1. Go to your Slack workspace settings
2. Create a new Incoming Webhook
3. Copy the webhook URL to your `.env` file

### Custom Webhooks

Send alerts to any HTTP endpoint:

```yaml
global_alerts:
  webhook:
    url: "${WEBHOOK_URL}"
    headers:
      Content-Type: "application/json"
      Authorization: "Bearer ${WEBHOOK_TOKEN}"
```

## Examples 💡

### Monitor with Multiple Alerts

```yaml
monitors:
  - name: "E-commerce API"
    url: "https://api.shop.com/products"
    interval: 15s
    expected_status: 200
    expected_latency: 1s
    headers:
      Authorization: "Bearer ${SHOP_API_TOKEN}"
    alerts:
      - name: "API Down"
        condition: "status != 200"
        channels: ["email", "slack"]
        recipients: ["devops@shop.com"]
        cooldown: 5m
      - name: "API Slow"
        condition: "latency > 2s"
        channels: ["slack"]
        cooldown: 10m
```

### Database Health Check

```yaml
monitors:
  - name: "Database Health"
    url: "https://api.example.com/health/db"
    interval: 30s
    body_validation:
      path: "$.database.status"
      expected: "connected"
    alerts:
      - name: "Database Connection Lost"
        condition: "ok == false"
        channels: ["email", "slack"]
        recipients: ["dba@example.com"]
        cooldown: 2m
```

### Multi-Environment Monitoring

```yaml
projects:
  - name: "Production"
    apiKey: "prod-key-123"
    monitors:
      - name: "Prod API"
        url: "https://api.prod.com"
        expected_latency: 500ms
  
  - name: "Staging"
    apiKey: "staging-key-456"
    monitors:
      - name: "Staging API"
        url: "https://api.staging.com"
        expected_latency: 1s
```

## Advanced Features 🔧

### JSON Response Validation

Validate specific fields in JSON responses:

```yaml
body_validation:
  path: "$.data.items[0].status"
  expected: "active"
```

### Custom Retry Logic

Configure retry behavior for failed requests:

```yaml
maxRetries: 5
jitter: 10s  # Random delay between retries
```

### SSL Certificate Validation

Control SSL certificate verification:

```yaml
ssl_verify: false  # Disable for self-signed certificates
```

## Monitoring Dashboard 📊

PingPulse provides console output with real-time monitoring status:

```
🚀 PingPulse Monitor Started
📊 Monitoring 5 endpoints across 2 projects
✅ Website Health: 200 OK (125ms)
❌ API Service: 500 Internal Server Error (2.3s)
📧 Email alert sent for API Service
💬 Slack alert sent for API Service
```

## Contributing 🤝

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## License 📄

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## Support 💬

- 🐛 **Issues**: [GitHub Issues](https://github.com/vimoksh-5/pingpulse/issues)
- 📖 **Documentation**: [Wiki](https://github.com/vimoksh-5/pingpulse/wiki)
- 💬 **Discussions**: [GitHub Discussions](https://github.com/vimoksh-5/pingpulse/discussions)

---

**Made with ❤️ by [Vimoksh]**

*Keep your services healthy with PingPulse!*