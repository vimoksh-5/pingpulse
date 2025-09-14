# PingPulse 🚀

A comprehensive monitoring solution with a Go backend and React TypeScript dashboard for tracking the health and performance of your web services with real-time alerts and beautiful visualizations.

<img width="1428" height="942" alt="image" src="https://github.com/user-attachments/assets/85df21be-f099-4328-9eb8-2cf88ca7fc4c" />


## Features ✨

### 🔧 **Backend Monitoring**
- **HTTP/HTTPS Monitoring** - Monitor any web endpoint with customizable intervals
- **Multi-Channel Alerts** - Email, Slack, and Webhook notifications with advanced conditions
- **Flexible Authentication** - Support for Bearer tokens, Basic auth, API keys, and OAuth2
- **Advanced Monitoring** - Response time tracking, status code validation, and JSON response validation
- **Project-Based Organization** - Group monitors by projects with individual API keys
- **Zero Dependencies** - Single binary deployment with no external database required
- **Enhanced Alert Conditions** - Complex conditions with AND/OR logic, comparison operators

### 📊 **Dashboard & API**
- **Modern React Dashboard** - Beautiful, responsive TypeScript interface
- **Real-time Updates** - Auto-refresh every 30 seconds with manual refresh option
- **Comprehensive REST API** - Full API access for programmatic integration
- **Advanced Filtering** - Filter by project, search monitors, status overview
- **Mobile-Friendly** - Responsive design that works on all devices
- **CORS Enabled** - Web dashboard integration ready

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

## REST API Documentation 🔌

### **Base URL**
All API endpoints: `http://localhost:8080`

### **Standardized Response Format**
All v1 API endpoints return responses in a consistent format:

```json
{
  "data": {}, // The actual response data
  "status_code": 1, // 1 for success, 0 for failure
  "status_message": "Success message"
}
```

### **Core Endpoints**

#### **GET /healthz** - Health Check
```bash
curl http://localhost:8080/healthz
```
```json
{
  "data": {"status": "healthy"},
  "status_code": 1,
  "status_message": "Service is healthy"
}
```

#### **GET /api/v1/statistics** - System Statistics
```bash
curl http://localhost:8080/api/v1/statistics
```
```json
{
  "data": {
    "total_projects": 3,
    "total_monitors": 8,
    "healthy_monitors": 6,
    "unhealthy_monitors": 2,
    "project_stats": {
      "Production APIs": {
        "monitor_count": 4,
        "healthy_count": 3,
        "unhealthy_count": 1
      }
    },
    "status_codes": {"200": 6, "404": 1, "500": 1},
    "average_latency": 245000000
  },
  "status_code": 1,
  "status_message": "Statistics retrieved successfully"
}
```

#### **GET /api/v1/projects** - List Projects
```bash
curl http://localhost:8080/api/v1/projects
```
```json
{
  "data": {
    "projects": ["Production APIs", "External Services"],
    "count": 2
  },
  "status_code": 1,
  "status_message": "Projects retrieved successfully"
}
```

#### **GET /api/v1/monitors** - List All Monitors
```bash
curl http://localhost:8080/api/v1/monitors
```
```json
{
  "data": {
    "monitors": ["User Service Health", "Payment API"],
    "count": 2
  },
  "status_code": 1,
  "status_message": "Monitors retrieved successfully"
}
```

#### **GET /api/v1/projects/{project}/monitors** - Project Monitors
```bash
curl "http://localhost:8080/api/v1/projects/Production%20APIs/monitors"
```
```json
{
  "data": {
    "project": "Production APIs",
    "monitors": ["User Service Health", "Payment API"],
    "count": 2
  },
  "status_code": 1,
  "status_message": "Project monitors retrieved successfully"
}
```

#### **GET /api/v1/results** - Query Results
Advanced results query with multiple filter options:

**Query Parameters:**
- `project` (string): Filter by project name
- `monitor` (string): Filter by monitor name  
- `url` (string): Filter by specific URL
- `status` (int): Filter by HTTP status code
- `ok` (bool): Filter by OK status (true/false)
- `since` (RFC3339): Filter results after timestamp
- `until` (RFC3339): Filter results before timestamp
- `limit` (int): Limit number of results

**Examples:**
```bash
# Get results for specific project
curl "http://localhost:8080/api/v1/results?project=Production%20APIs"

# Get failed results
curl "http://localhost:8080/api/v1/results?ok=false"

# Get results with specific status code
curl "http://localhost:8080/api/v1/results?status=500"

# Complex query: Failed results for project in time range
curl "http://localhost:8080/api/v1/results?project=Production%20APIs&ok=false&since=2024-01-01T09:00:00Z&limit=50"
```

**Response:**
```json
{
  "data": {
    "results": [
      {
        "project": "Production APIs",
        "name": "User Service Health",
        "url": "https://api.company.com/users/health",
        "ok": true,
        "status": 200,
        "latency": 150000000,
        "checked_at": "2024-01-01T10:30:00Z"
      }
    ],
    "count": 1,
    "filter": {"project": "Production APIs", "limit": 50}
  },
  "status_code": 1,
  "status_message": "Results retrieved successfully"
}
```

### **Legacy Endpoints (Backward Compatibility)**
- `GET /results` - Returns all monitoring results or results for specific URL
- `GET /results?url=https://example.com` - Filter by URL

## Configuration Reference 📖

### **Basic Configuration**

Create a `config.yaml` file:

```yaml
default:
  interval: 60s
  timeout: 5s
  max_retries: 3
  jitter: 5s
  user_agent: "PingPulse/1.0"

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
  - name: "Production APIs"
    api_key: "prod_abc123"
    monitors:
      - name: "User Service Health"
        url: "https://api.company.com/users/health"
        interval: 30s
        expected_status: 200
        expected_latency: 2s
        alerts:
          - name: "User Service Down"
            condition: "status != 200 OR ok == false"
            channels: ["email", "slack"]
            recipients: ["ops@company.com"]
            cooldown: 5m
```

### **Enhanced Alert Conditions**

PingPulse supports flexible alert conditions with multiple operators:

#### **Supported Fields**
- `status`: HTTP status code
- `latency`: Response time (supports duration formats like `2s`, `500ms`)
- `ok`: Boolean status of the check

#### **Supported Operators**
- `==` or `=`: Equal to
- `!=`: Not equal to
- `>`: Greater than
- `<`: Less than
- `>=`: Greater than or equal to
- `<=`: Less than or equal to

#### **Logical Operators**
- `AND`: Both conditions must be true
- `OR`: Either condition must be true

#### **Example Conditions**
```yaml
# Simple conditions
condition: "status != 200"
condition: "latency > 2s"
condition: "ok == false"

# Comparison conditions
condition: "status >= 500"
condition: "latency <= 1s"

# Complex conditions with AND/OR
condition: "status != 200 OR ok == false"
condition: "status >= 500 AND latency > 5s"
condition: "status >= 400 AND status < 500"
condition: "latency > 1s AND latency <= 3s"

# Multiple conditions
condition: "status == 401 OR status == 403"
condition: "status >= 500 OR latency > 10s OR ok == false"
```

### **Authentication Types**

#### **Bearer Token**
```yaml
auth:
  type: "bearer"
  token: "${API_TOKEN}"
```

#### **Basic Authentication**
```yaml
auth:
  type: "basic"
  username: "${USERNAME}"
  password: "${PASSWORD}"
```

#### **API Key**
```yaml
auth:
  type: "apikey"
  api_key: "${API_KEY}"
  api_key_header: "X-API-Key"
```

#### **OAuth2**
```yaml
auth:
  type: "oauth2"
  client_id: "${CLIENT_ID}"
  client_secret: "${CLIENT_SECRET}"
  token_url: "https://auth.example.com/oauth/token"
  token: "${OAUTH_ACCESS_TOKEN}"
```

### **Advanced Monitor Configuration**

```yaml
monitors:
  - name: "Complex API Monitor"
    url: "https://api.example.com/health"
    interval: 45s
    timeout: 10s
    max_retries: 2
    method: "POST"
    headers:
      Accept: "application/json"
      X-Client-Version: "1.0"
    body: {"action": "health_check", "timestamp": "{{.Now}}"}
    auth:
      type: "bearer"
      token: "${API_TOKEN}"
    expected_status: 200
    expected_latency: 3s
    ssl_verify: true
    follow_redirects: true
    body_validation:
      path: "$.healthy"
      expected: true
    alerts:
      - name: "API Critical Issues"
        condition: "status != 200 OR latency > 3s"
        channels: ["email", "slack", "webhook"]
        recipients: ["alerts@company.com"]
        cooldown: 15m
      - name: "API Performance Warning"
        condition: "latency > 1s AND latency <= 3s"
        channels: ["slack"]
        cooldown: 30m
```

## Environment Variables 🔧

Create a `.env` file:

```bash
# Gmail SMTP Configuration
GMAIL_USERNAME=your-email@gmail.com
GMAIL_APP_PASSWORD=your-app-password

# Slack Webhook Configuration  
SLACK_WEBHOOK_URL=https://hooks.slack.com/services/YOUR/SLACK/WEBHOOK

# API Tokens
API_TOKEN=your-bearer-token
OAUTH_CLIENT_ID=your-oauth-client-id
OAUTH_CLIENT_SECRET=your-oauth-client-secret

# Custom Webhook
WEBHOOK_URL=https://webhook.site/your-unique-id
WEBHOOK_TOKEN=your-webhook-auth-token
```

## Alert Channels 📢

### **Email Alerts (Gmail)**
1. Enable 2-Factor Authentication on Gmail
2. Generate App Password: Google Account → Security → App passwords
3. Use App Password in `.env` file

### **Slack Alerts**
1. Go to Slack workspace settings
2. Create new Incoming Webhook
3. Copy webhook URL to `.env` file

### **Custom Webhooks**
Send alerts to any HTTP endpoint with custom headers and authentication.

## Dashboard Development 💻

### **Technology Stack**
- **Frontend**: React 18 + TypeScript
- **Styling**: CSS Modules with CSS Custom Properties
- **HTTP Client**: Axios with comprehensive error handling
- **Build Tool**: Create React App
- **Icons**: Custom SVG icons

### **Development Setup**
```bash
cd dashboard
npm install
npm start  # Starts on http://localhost:3000
```

### **Dashboard Features**
- **Zero Configuration**: Works out of the box with CSS Modules
- **Type Safety**: Full TypeScript coverage
- **Responsive Design**: Mobile-friendly interface
- **Real-time Updates**: Auto-refresh with connection status
- **Error Handling**: Graceful degradation when backend unavailable

## Production Deployment 🚀

### **Backend Deployment**
```bash
# Build optimized binary
go build -ldflags="-s -w" -o pingpulse ./cmd/pingpulse

# Run with systemd
sudo systemctl enable pingpulse
sudo systemctl start pingpulse
```

### **Dashboard Deployment**
```bash
cd dashboard
npm run build

# Deploy to static hosting (Netlify, Vercel, S3, etc.)
# Set REACT_APP_API_URL=https://your-api-domain.com
```

### **Docker Deployment**
```dockerfile
# Backend
FROM golang:1.21-alpine AS builder
WORKDIR /app
COPY . .
RUN go build -o pingpulse ./cmd/pingpulse

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=builder /app/pingpulse .
COPY config.yaml .
CMD ["./pingpulse"]
```

## Examples & Use Cases 💡

### **E-commerce Monitoring**
```yaml
projects:
  - name: "E-commerce Platform"
    monitors:
      - name: "Product API"
        url: "https://api.shop.com/products"
        condition: "status != 200 OR latency > 1s"
      - name: "Payment Gateway"
        url: "https://payments.shop.com/health"
        condition: "status >= 500 OR ok == false"
      - name: "User Authentication"
        url: "https://auth.shop.com/health"
        condition: "status != 200"
```

### **Microservices Health Checks**
```yaml
projects:
  - name: "Microservices"
    monitors:
      - name: "User Service"
        url: "http://user-service:8080/health"
        body_validation:
          path: "$.status"
          expected: "healthy"
      - name: "Order Service"
        url: "http://order-service:8080/health"
      - name: "Inventory Service"
        url: "http://inventory-service:8080/health"
```

### **Multi-Environment Monitoring**
```yaml
projects:
  - name: "Production"
    monitors:
      - name: "Prod API"
        url: "https://api.prod.com"
        expected_latency: 500ms
        alerts:
          - condition: "latency > 1s OR status != 200"
            channels: ["email", "slack"]
            recipients: ["oncall@company.com"]
  
  - name: "Staging"  
    monitors:
      - name: "Staging API"
        url: "https://api.staging.com"
        expected_latency: 2s
        alerts:
          - condition: "status != 200"
            channels: ["slack"]
```

## Troubleshooting 🔍

### **Common Issues**

#### **Dashboard Won't Load**
1. Ensure backend is running on port 8080
2. Check CORS is enabled (automatic in recent versions)
3. Verify network connectivity

#### **No Monitoring Data**
1. Check `config.yaml` syntax
2. Verify environment variables are set
3. Check monitor URLs are accessible

#### **Alerts Not Working**
1. Verify SMTP/Slack credentials
2. Check alert conditions syntax
3. Ensure cooldown periods aren't preventing alerts

### **Debug Mode**
```bash
# Run with verbose logging
./pingpulse --debug

# Check API responses
curl -v http://localhost:8080/api/v1/statistics
```

## API Rate Limiting & Security 🔒

### **Current Status**
- No rate limiting implemented
- No authentication required for API access
- CORS enabled for web dashboard

### **Production Recommendations**
- Implement API key authentication
- Add rate limiting middleware
- Use HTTPS in production
- Configure firewall rules
- Set up monitoring for the monitoring system

## Contributing 🤝

### **Development Setup**
```bash
# Backend development
go mod tidy
go run ./cmd/pingpulse

# Frontend development  
cd dashboard
npm install
npm start

# Run tests
go test ./...
cd dashboard && npm test
```

### **Contributing Guidelines**
1. Fork the repository
2. Create feature branch (`git checkout -b feature/amazing-feature`)
3. Follow existing code patterns
4. Add tests for new features
5. Update documentation
6. Submit pull request

## Performance & Scalability 📈

### **Current Capabilities**
- **Concurrent Monitoring**: Handles hundreds of endpoints
- **Memory Efficient**: In-memory storage with configurable retention
- **Low Latency**: Sub-second API response times
- **Lightweight**: Single binary < 10MB

### **Scaling Recommendations**
- Use external database for large deployments
- Implement horizontal scaling with load balancer
- Add caching layer for high-traffic APIs
- Consider message queue for alert processing

## License 📄

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## Support & Community 💬

- 🐛 **Issues**: [GitHub Issues](https://github.com/vimoksh-5/pingpulse/issues)
- 📖 **Documentation**: [Wiki](https://github.com/vimoksh-5/pingpulse/wiki)  
- 💬 **Discussions**: [GitHub Discussions](https://github.com/vimoksh-5/pingpulse/discussions)
- 📧 **Email**: support@pingpulse.dev

## Roadmap 🗺️

### **Upcoming Features**
- [ ] Historical data persistence (PostgreSQL/MySQL)
- [ ] Advanced dashboard charts and analytics
- [ ] User authentication and role-based access
- [ ] Mobile app for iOS/Android
- [ ] Kubernetes operator
- [ ] Prometheus metrics export
- [ ] Custom dashboard themes
- [ ] Alert escalation policies

---

**Made with ❤️ by [Vimoksh](https://github.com/vimoksh-5)**

*Keep your services healthy with PingPulse!* 🚀

---

## Screenshots 📸

### Dashboard Overview
![PingPulse Dashboard](dashboard-screenshot.png)

The dashboard provides:
- **Real-time Statistics**: Overview of all projects and monitors
- **Status Indicators**: Color-coded health status for quick assessment  
- **Detailed Metrics**: Individual monitor cards with comprehensive information
- **Filtering & Search**: Easy navigation through large numbers of monitors
- **Responsive Design**: Works perfectly on desktop, tablet, and mobile devices
