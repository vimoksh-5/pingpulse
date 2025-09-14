# PingPulse Configuration Guide

Complete guide for configuring PingPulse monitoring system.

## 🎯 What PingPulse Does

PingPulse is a continuous API monitoring system that:

- 🔄 **Monitors APIs** - Continuously checks HTTP/HTTPS endpoints
- ⏱️ **Tracks Performance** - Measures response times and status codes
- 🚨 **Sends Alerts** - Notifies via email, Slack, or webhooks when issues occur
- 📊 **Stores Results** - Keeps monitoring history for analysis
- 🌐 **Provides API** - REST API for programmatic access to monitoring data
- 📱 **Web Dashboard** - Beautiful React interface for real-time monitoring

## 🚀 Quick Start

### Minimal Configuration

Create a `config.yaml` file to get started:

```yaml
projects:
  - name: "My First Project"
    monitors:
      - name: "Google Health Check"
        url: "https://www.google.com"
        interval: 60s
```

This minimal config will:

- Check Google every 60 seconds
- Store results in memory
- Provide API access at `http://localhost:8080`
- Show results in dashboard at `http://localhost:3000`

### Complete Example

```yaml
default:
  interval: 60s
  timeout: 5s

global_alerts:
  email:
    smtp_host: "smtp.gmail.com"
    smtp_port: 587
    username: "${GMAIL_USERNAME}"
    password: "${GMAIL_APP_PASSWORD}"
    from: "${GMAIL_USERNAME}"
    use_tls: true

projects:
  - name: "Production APIs"
    monitors:
      - name: "User API"
        url: "https://api.example.com/users"
        interval: 30s
        alerts:
          - name: "API Down"
            condition: "status != 200"
            channels: ["email"]
            recipients: ["ops@company.com"]
```

## ⚙️ Configuration Structure

### 1. Default Settings

```yaml
default:
  interval: 60s # How often to check each API (default)
  timeout: 5s # Maximum time to wait for response
  maxRetries: 3 # Number of retry attempts on failure
  jitter: 5s # Random delay to avoid thundering herd
  userAgent: "PingPulse/1.0" # HTTP User-Agent header
```

**Use Cases:**

- **interval**: Check critical APIs every 15s, less important ones every 5 minutes
- **timeout**: Set based on expected response time (1s for fast APIs, 30s for slow ones)
- **maxRetries**: Retry 3 times for transient failures, 1 time for critical alerts
- **jitter**: Prevents all monitors from hitting at exactly the same time

### 2. Global Alert Configuration

Configure notification channels that can be used by all monitors:

```yaml
global_alerts:
  email:
    smtp_host: "smtp.gmail.com"
    smtp_port: 587
    username: "${GMAIL_USERNAME}"
    password: "${GMAIL_APP_PASSWORD}" # Use App Password, not regular password
    from: "${GMAIL_USERNAME}"
    use_tls: true

  slack:
    webhook_url: "${SLACK_WEBHOOK_URL}"
    channel: "#alerts"
    username: "PingPulse Bot"

  webhook:
    url: "${WEBHOOK_URL}"
    headers:
      Authorization: "Bearer ${WEBHOOK_TOKEN}"
      Content-Type: "application/json"
```

**Channel Use Cases:**

| Channel     | Best For                                 | Example                                |
| ----------- | ---------------------------------------- | -------------------------------------- |
| **Email**   | Formal notifications, management reports | `ops@company.com`, `cto@company.com`   |
| **Slack**   | Real-time team notifications             | `#devops`, `#alerts`, `#critical`      |
| **Webhook** | Integration with other systems           | PagerDuty, OpsGenie, custom dashboards |

### **3. Projects (Multi-tenant Support)**

```yaml
projects:
  - name: "Production APIs" # Project name
    apiKey: "prod_abc123" # Project identifier
    monitors: [...] # List of monitors for this project
  - name: "Development APIs" # Separate project
    apiKey: "dev_xyz789" # Different API key
    monitors: [...] # Different monitors
```

**Use Cases:**

- **Separate environments**: Production, staging, development
- **Different teams**: Frontend team, backend team, data team
- **Different clients**: Customer A, Customer B, internal tools

### **4. Monitor Configuration**

```yaml
monitors:
  - name: "User Service Health" # Friendly name
    url: "https://api.company.com/users/health" # API endpoint
    interval: 30s # Override default interval
    timeout: 10s # Override default timeout
    maxRetries: 2 # Override default retries
    method: "GET" # HTTP method (GET, POST, PUT, etc.)
    expected_status: 200 # Expected HTTP status code
    expected_latency: 2s # Maximum acceptable response time
    ssl_verify: true # Verify SSL certificates
    follow_redirects: true # Follow HTTP redirects
```

**Use Cases:**

- **Health checks**: Monitor `/health` endpoints
- **API endpoints**: Monitor critical business APIs
- **Database connections**: Monitor database health
- **External services**: Monitor third-party APIs

### **5. Authentication Configuration**

```yaml
auth:
  type: "bearer"                  # Authentication type
  token: "${API_TOKEN}"           # Bearer token
  # OR
  type: "basic"                   # Basic authentication
  username: "${API_USERNAME}"     # Username
  password: "${API_PASSWORD}"     # Password
  # OR
  type: "apikey"                  # API key authentication
  api_key: "${API_KEY}"           # API key
  api_key_header: "X-API-Key"     # Header name (default: X-API-Key)
  # OR
  type: "oauth2"                  # OAuth2 authentication
  client_id: "${OAUTH_CLIENT_ID}" # OAuth2 client ID
  client_secret: "${OAUTH_CLIENT_SECRET}"  # OAuth2 client secret
  token_url: "https://auth.example.com/oauth/token"  # Token endpoint
  token: "${OAUTH_ACCESS_TOKEN}"  # Access token
```

**Use Cases:**

- **Bearer tokens**: JWT tokens, OAuth2 access tokens
- **Basic auth**: Legacy APIs, internal services
- **API keys**: Third-party services, rate-limited APIs
- **OAuth2**: Modern APIs with proper authentication

### **6. Request Configuration**

```yaml
headers:                          # Custom HTTP headers
  Authorization: "Bearer ${TOKEN}"
  Content-Type: "application/json"
  X-Request-ID: "{{.UUID}}"
  X-Client-Version: "1.0"

body:                             # Request body (for POST/PUT)
  {"action": "health_check", "timestamp": "{{.Now}}"}
  # OR
  {"query": "SELECT 1", "timeout": 5}
```

**Use Cases:**

- **Custom headers**: Authentication, tracking, versioning
- **POST requests**: Health checks that require data
- **GraphQL queries**: Complex queries for health checks
- **Database queries**: Direct database health checks

### **7. Response Validation**

```yaml
body_validation:                  # Validate response body
  path: "$.status"                # JSONPath to check
  expected: "healthy"             # Expected value
  # OR
  path: "$.healthy"               # Check boolean field
  expected: true                  # Expected boolean value
```

**Use Cases:**

- **Health check responses**: Verify status is "healthy"
- **Database responses**: Check if query succeeded
- **API responses**: Validate specific fields
- **Complex validations**: Check multiple fields

### **8. Alert Configuration**

```yaml
alerts:
  - name: "API Down Alert" # Alert name
    condition: "status != 200" # When to trigger
    channels: ["email", "slack"] # Where to send
    recipients: ["ops@company.com"] # Email recipients
    cooldown: 5m # Wait between alerts
  - name: "Slow Response Alert" # Another alert
    condition: "latency > 5s" # Different condition
    channels: ["slack"] # Different channel
    cooldown: 10m # Different cooldown
```

**Use Cases:**

- **Critical alerts**: Immediate notification for down services
- **Performance alerts**: Notify when response times are high
- **Escalation**: Different alerts for different severity levels
- **Cooldown**: Prevent alert spam during outages

## 🏢 Real-World Configuration Examples

### E-commerce Platform

Monitor critical e-commerce services with different alert levels:

```yaml
projects:
  - name: "E-commerce Production"
    monitors:
      # Critical: Payment processing
      - name: "Payment Gateway"
        url: "https://payments.company.com/health"
        interval: 15s
        expected_latency: 1s
        auth:
          type: "bearer"
          token: "${PAYMENT_API_TOKEN}"
        alerts:
          - name: "Payment Critical"
            condition: "status != 200 OR latency > 3s"
            channels: ["email", "slack"]
            recipients: ["finance@company.com", "ops@company.com"]
            cooldown: 2m

      # Important: Product catalog
      - name: "Product API"
        url: "https://api.company.com/products/health"
        interval: 30s
        alerts:
          - name: "Product API Issues"
            condition: "status >= 500"
            channels: ["slack"]
            cooldown: 5m

      # User-facing: Website
      - name: "Website Health"
        url: "https://www.company.com/health"
        interval: 60s
        alerts:
          - name: "Website Down"
            condition: "status != 200"
            channels: ["email", "slack"]
            recipients: ["marketing@company.com"]
```

### Microservices Architecture

Monitor internal services with service mesh health checks:

```yaml
projects:
  - name: "Microservices Cluster"
    monitors:
      # Core services
      - name: "User Service"
        url: "http://user-service.internal:8080/health"
        interval: 30s
        body_validation:
          path: "$.status"
          expected: "healthy"
        alerts:
          - name: "User Service Down"
            condition: "ok == false"
            channels: ["slack"]

      - name: "Order Service"
        url: "http://order-service.internal:8080/health"
        interval: 30s
        alerts:
          - name: "Order Processing Issues"
            condition: "status != 200 OR latency > 2s"
            channels: ["slack"]

      # Database connections
      - name: "Database Connectivity"
        url: "http://db-proxy.internal:8080/ping"
        interval: 60s
        alerts:
          - name: "Database Connection Issues"
            condition: "status != 200"
            channels: ["email", "slack"]
            recipients: ["dba@company.com"]
```

### Third-Party API Monitoring

Monitor external dependencies with appropriate intervals:

```yaml
projects:
  - name: "External Dependencies"
    monitors:
      # Payment processor
      - name: "Stripe API"
        url: "https://api.stripe.com/v1/charges"
        interval: 120s # Less frequent for external APIs
        method: "GET"
        auth:
          type: "bearer"
          token: "${STRIPE_SECRET_KEY}"
        headers:
          Stripe-Version: "2023-10-16"
        alerts:
          - name: "Stripe API Issues"
            condition: "status != 200"
            channels: ["email"]
            recipients: ["payments@company.com"]

      # Communication service
      - name: "Twilio API"
        url: "https://api.twilio.com/2010-04-01/Accounts/${TWILIO_ACCOUNT_SID}.json"
        interval: 300s # 5 minutes for less critical
        auth:
          type: "basic"
          username: "${TWILIO_ACCOUNT_SID}"
          password: "${TWILIO_AUTH_TOKEN}"
        alerts:
          - name: "SMS Service Issues"
            condition: "status >= 400"
            channels: ["slack"]

      # CDN monitoring
      - name: "CloudFlare CDN"
        url: "https://cdn.company.com/health"
        interval: 180s
        alerts:
          - name: "CDN Performance"
            condition: "latency > 5s OR status != 200"
            channels: ["slack"]
```

### Multi-Environment Setup

Monitor different environments with environment-specific settings:

```yaml
projects:
  # Production environment
  - name: "Production"
    monitors:
      - name: "Prod API Gateway"
        url: "https://api.company.com/health"
        interval: 30s
        expected_latency: 500ms
        alerts:
          - name: "Production Critical"
            condition: "status != 200 OR latency > 2s"
            channels: ["email", "slack"]
            recipients: ["oncall@company.com"]
            cooldown: 1m

  # Staging environment
  - name: "Staging"
    monitors:
      - name: "Staging API Gateway"
        url: "https://staging-api.company.com/health"
        interval: 120s
        expected_latency: 2s
        alerts:
          - name: "Staging Issues"
            condition: "status != 200"
            channels: ["slack"]
            cooldown: 10m

  # Development environment
  - name: "Development"
    monitors:
      - name: "Dev API Gateway"
        url: "https://dev-api.company.com/health"
        interval: 300s # Less frequent monitoring
        alerts:
          - name: "Dev Environment Down"
            condition: "ok == false"
            channels: ["slack"]
            cooldown: 30m
```

## 📊 Monitoring Best Practices

### Interval Guidelines

Choose monitoring intervals based on service criticality:

| Service Type              | Interval      | Reasoning                                       |
| ------------------------- | ------------- | ----------------------------------------------- |
| **Payment/Critical APIs** | 15-30s        | Immediate detection of revenue-impacting issues |
| **User-facing APIs**      | 30-60s        | Quick detection of user experience issues       |
| **Internal Services**     | 1-2 minutes   | Balance between detection speed and load        |
| **Background Services**   | 2-5 minutes   | Less critical, reduce monitoring overhead       |
| **External APIs**         | 2-5 minutes   | Respect rate limits, less control over uptime   |
| **Database Health**       | 30s-2 minutes | Critical but stable, moderate monitoring        |

### Alert Condition Patterns

Common alert patterns for different scenarios:

```yaml
# Basic availability
condition: "status != 200"

# Performance degradation
condition: "latency > 2s"

# Critical failure (any issue)
condition: "ok == false"

# Server errors only
condition: "status >= 500"

# Authentication/authorization issues
condition: "status == 401 OR status == 403"

# Combined availability and performance
condition: "status != 200 OR latency > 5s"

# Moderate performance warning
condition: "latency > 1s AND latency <= 3s"
```

### Configuration Best Practices

#### 1. Start Simple, Iterate

```yaml
# Begin with basic monitoring
monitors:
  - name: "API Health"
    url: "https://api.example.com/health"
    interval: 60s

# Add complexity as needed
monitors:
  - name: "API Health"
    url: "https://api.example.com/health"
    interval: 30s
    expected_latency: 2s
    alerts:
      - name: "API Issues"
        condition: "status != 200 OR latency > 2s"
        channels: ["slack"]
```

#### 2. Environment-Specific Settings

```yaml
# Production: Strict monitoring
- name: "Production API"
  interval: 30s
  expected_latency: 1s
  alerts:
    - condition: "status != 200"
      channels: ["email", "slack"]
      cooldown: 1m

# Staging: Relaxed monitoring
- name: "Staging API"
  interval: 120s
  expected_latency: 3s
  alerts:
    - condition: "status >= 500"
      channels: ["slack"]
      cooldown: 10m
```

#### 3. Alert Fatigue Prevention

```yaml
alerts:
  # Critical: Immediate notification
  - name: "Service Down"
    condition: "ok == false"
    channels: ["email", "slack"]
    cooldown: 2m # Short cooldown for critical issues

  # Warning: Less frequent notification
  - name: "Performance Degradation"
    condition: "latency > 2s AND latency <= 5s"
    channels: ["slack"]
    cooldown: 15m # Longer cooldown for warnings
```

#### 4. Security Best Practices

```yaml
# Use environment variables for secrets
auth:
  type: "bearer"
  token: "${API_TOKEN}"  # Never hardcode tokens

# Separate credentials by environment
auth:
  type: "basic"
  username: "${PROD_API_USERNAME}"  # Environment-specific
  password: "${PROD_API_PASSWORD}"
```

## 🔧 Environment Variables Setup

### Required Variables

Create a `.env` file or set system environment variables:

```bash
# Email configuration (Gmail example)
GMAIL_USERNAME=monitoring@company.com
GMAIL_APP_PASSWORD=abcd-efgh-ijkl-mnop  # 16-character app password

# Slack configuration
SLACK_WEBHOOK_URL=https://hooks.slack.com/services/T00000000/B00000000/XXXXXXXXXXXXXXXXXXXXXXXX

# API authentication tokens
PROD_API_TOKEN=eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...
STAGING_API_TOKEN=eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...

# Third-party service credentials
STRIPE_SECRET_KEY=sk_live_...
TWILIO_ACCOUNT_SID=AC...
TWILIO_AUTH_TOKEN=...

# Custom webhook
WEBHOOK_URL=https://your-system.com/webhooks/pingpulse
WEBHOOK_TOKEN=your-secure-webhook-token
```

### Environment Variable Security

```bash
# Set proper permissions
chmod 600 .env

# For systemd services
sudo mkdir -p /etc/pingpulse
sudo tee /etc/pingpulse/environment << EOF
GMAIL_USERNAME=monitoring@company.com
GMAIL_APP_PASSWORD=your-app-password
EOF
sudo chmod 600 /etc/pingpulse/environment
```

## 📋 Configuration Validation

### Validate Your Configuration

```bash
# Test configuration syntax
./pingpulse --validate-config

# Test specific monitor
curl -v https://your-api.com/health

# Test alert channels
# Email test (manual)
echo "Test email" | mail -s "PingPulse Test" your-email@company.com

# Slack test
curl -X POST -H 'Content-type: application/json' \
  --data '{"text":"Test from PingPulse"}' \
  $SLACK_WEBHOOK_URL
```

### Common Configuration Errors

```yaml
# ❌ Wrong: Missing quotes for complex strings
condition: status != 200 OR latency > 2s

# ✅ Correct: Properly quoted
condition: "status != 200 OR latency > 2s"

# ❌ Wrong: Invalid duration format
interval: 60
timeout: 5

# ✅ Correct: Include time unit
interval: 60s
timeout: 5s

# ❌ Wrong: Hardcoded secrets
password: "my-secret-password"

# ✅ Correct: Environment variables
password: "${API_PASSWORD}"
```

## 🎯 Quick Reference

### Minimal Working Config

```yaml
projects:
  - name: "Test"
    monitors:
      - name: "Google"
        url: "https://www.google.com"
```

### Production-Ready Config

```yaml
default:
  interval: 60s
  timeout: 10s
  max_retries: 3

global_alerts:
  email:
    smtp_host: "smtp.gmail.com"
    smtp_port: 587
    username: "${GMAIL_USERNAME}"
    password: "${GMAIL_APP_PASSWORD}"
    use_tls: true
  slack:
    webhook_url: "${SLACK_WEBHOOK_URL}"
    channel: "#alerts"

projects:
  - name: "Production APIs"
    monitors:
      - name: "Main API"
        url: "https://api.company.com/health"
        interval: 30s
        auth:
          type: "bearer"
          token: "${API_TOKEN}"
        alerts:
          - name: "API Critical"
            condition: "status != 200 OR latency > 3s"
            channels: ["email", "slack"]
            recipients: ["ops@company.com"]
            cooldown: 5m
```
