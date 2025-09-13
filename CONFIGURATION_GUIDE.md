# PingPulse Configuration Guide

## **What PingPulse Does:**

PingPulse is a continuous API monitoring system that:

- 🔄 **Continuously pings APIs** at specified intervals
- ⏱️ **Measures response times** and status codes
- 🚨 **Sends alerts** when APIs fail or are slow
- 📊 **Stores results** for analysis and reporting
- 🌐 **Provides web API** to view monitoring data

## **Configuration Structure:**

### **1. Default Settings**

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

### **2. Global Alerting Configuration**

```yaml
global_alerts:
  email:
    smtp_host: "smtp.gmail.com" # Email server
    smtp_port: 587 # SMTP port (587 for TLS)
    username: "${GMAIL_USERNAME}" # Your email address
    password: "${GMAIL_APP_PASSWORD}" # App password (not regular password)
    from: "${GMAIL_USERNAME}" # Sender email
    use_tls: true # Use TLS encryption
  slack:
    webhook_url: "${SLACK_WEBHOOK_URL}" # Slack webhook URL
    channel: "#alerts" # Slack channel name
    username: "PingPulse" # Bot username
  webhook:
    url: "${WEBHOOK_URL}" # Custom webhook URL
    headers: # Custom headers
      Content-Type: "application/json"
```

**Use Cases:**

- **Email**: Send alerts to ops team, management, stakeholders
- **Slack**: Real-time notifications to development teams
- **Webhook**: Integrate with PagerDuty, OpsGenie, custom systems

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

## **Real-World Use Cases:**

### **1. E-commerce Website**

```yaml
projects:
  - name: "E-commerce Production"
    monitors:
      - name: "Payment API"
        url: "https://payments.company.com/health"
        interval: 15s
        expected_latency: 1s
        auth:
          type: "bearer"
          token: "${PAYMENT_TOKEN}"
        alerts:
          - name: "Payment Down"
            condition: "status != 200"
            channels: ["email", "slack"]
            recipients: ["finance@company.com"]
```

### **2. Microservices Architecture**

```yaml
projects:
  - name: "Microservices"
    monitors:
      - name: "User Service"
        url: "https://users.company.com/health"
        interval: 30s
      - name: "Order Service"
        url: "https://orders.company.com/health"
        interval: 30s
      - name: "Inventory Service"
        url: "https://inventory.company.com/health"
        interval: 30s
```

### **3. Third-party Integrations**

```yaml
projects:
  - name: "External APIs"
    monitors:
      - name: "Stripe API"
        url: "https://api.stripe.com/v1/charges"
        method: "GET"
        auth:
          type: "apikey"
          api_key: "${STRIPE_SECRET_KEY}"
        expected_status: 200
      - name: "Twilio API"
        url: "https://api.twilio.com/2010-04-01/Accounts"
        auth:
          type: "basic"
          username: "${TWILIO_ACCOUNT_SID}"
          password: "${TWILIO_AUTH_TOKEN}"
```

### **4. Database Monitoring**

```yaml
projects:
  - name: "Database Health"
    monitors:
      - name: "PostgreSQL"
        url: "https://db.company.com/health"
        method: "POST"
        body: { "query": "SELECT 1", "timeout": 5 }
        expected_status: 200
        body_validation:
          path: "$.status"
          expected: "healthy"
```

## **Monitoring Intervals Guide:**

- **Critical APIs**: 15-30 seconds
- **Important APIs**: 1-2 minutes
- **Less critical APIs**: 5-10 minutes
- **External APIs**: 2-5 minutes
- **Database health**: 30 seconds - 2 minutes

## **Alert Conditions:**

- `status != 200` - API is down
- `latency > 2s` - API is slow
- `latency > 5s` - API is very slow
- `ok == false` - Any failure condition
- `status >= 500` - Server errors
- `status == 401` - Authentication issues
- `status == 403` - Authorization issues

## **Best Practices:**

1. **Start with longer intervals** (2-5 minutes) and reduce as needed
2. **Set appropriate timeouts** based on expected response times
3. **Use cooldowns** to prevent alert spam
4. **Test your alerts** before going to production
5. **Monitor from multiple locations** if possible
6. **Use different alert channels** for different severity levels
7. **Keep credentials secure** in environment variables
8. **Document your monitoring strategy** for your team
