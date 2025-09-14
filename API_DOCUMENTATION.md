# PingPulse API Documentation

Complete REST API reference for PingPulse monitoring system.

## 🌐 Base URL

All API endpoints are available at: `http://localhost:8080`

## 📋 Quick Reference

| Endpoint                              | Method | Description                           |
| ------------------------------------- | ------ | ------------------------------------- |
| `/healthz`                            | GET    | Service health check                  |
| `/api/v1/statistics`                  | GET    | System statistics and metrics         |
| `/api/v1/projects`                    | GET    | List all projects                     |
| `/api/v1/monitors`                    | GET    | List all monitors                     |
| `/api/v1/projects/{project}/monitors` | GET    | Monitors for specific project         |
| `/api/v1/results`                     | GET    | Query monitoring results with filters |
| `/results`                            | GET    | Legacy results endpoint               |

## Standardized Response Format

All enhanced API endpoints (v1) return responses in a standardized format:

```json
{
  "data": {}, // The actual response data
  "status_code": 1, // 1 for success, 0 for failure
  "status_message": "Success message" // Human-readable status message
}
```

### Status Codes

- `1`: Success - The request was processed successfully
- `0`: Failure - The request failed due to an error

### Success Response Example

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

### Error Response Example

```json
{
  "data": null,
  "status_code": 0,
  "status_message": "Invalid status parameter: must be a valid integer"
}
```

## Legacy Endpoints (Backward Compatibility)

### GET /results

Returns all monitoring results or results for a specific URL.

**Query Parameters:**

- `url` (optional): Filter results by specific URL

**Example:**

```bash
# Get all results
curl "http://localhost:8080/results"

# Get results for specific URL
curl "http://localhost:8080/results?url=https://api.example.com/health"
```

### GET /healthz

Health check endpoint.

**Response:**

```json
{
  "data": {
    "status": "healthy"
  },
  "status_code": 1,
  "status_message": "Service is healthy"
}
```

## 🔍 Enhanced API Endpoints (v1)

### GET /api/v1/results

Query monitoring results with advanced filtering capabilities.

**Query Parameters:**

| Parameter | Type   | Description                | Example                      |
| --------- | ------ | -------------------------- | ---------------------------- |
| `project` | string | Filter by project name     | `Production APIs`            |
| `monitor` | string | Filter by monitor name     | `User Service Health`        |
| `url`     | string | Filter by specific URL     | `https://api.example.com`    |
| `status`  | int    | Filter by HTTP status code | `200`, `404`, `500`          |
| `ok`      | bool   | Filter by success status   | `true`, `false`              |
| `since`   | string | Results after timestamp    | `1757849600` (epoch seconds) |
| `until`   | string | Results before timestamp   | `1757849700` (epoch seconds) |
| `limit`   | int    | Maximum results to return  | `50`, `100`                  |

**Usage Examples:**

```bash
# Get all results
curl "http://localhost:8080/api/v1/results"

# Filter by project
curl "http://localhost:8080/api/v1/results?project=Production%20APIs"

# Get only failed results
curl "http://localhost:8080/api/v1/results?ok=false"

# Filter by status code
curl "http://localhost:8080/api/v1/results?status=500"

# Get results for specific monitor
curl "http://localhost:8080/api/v1/results?monitor=User%20Service%20Health"

# Time-based filtering (epoch timestamp)
curl "http://localhost:8080/api/v1/results?since=1757849600&limit=50"

# Complex filtering with time range
curl "http://localhost:8080/api/v1/results?project=Production%20APIs&ok=false&since=1757849600&until=1757850000&limit=100"
```

**Response Format:**

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
        "latency": "150ms",
        "checked_at": "1757851919"
      }
    ],
    "count": 1,
    "filter": {
      "project": "Production APIs",
      "limit": 50
    }
  },
  "status_code": 1,
  "status_message": "Results retrieved successfully"
}
```

### GET /api/v1/projects

List all configured projects in the monitoring system.

**Usage:**

```bash
curl "http://localhost:8080/api/v1/projects"
```

**Response:**

```json
{
  "data": {
    "projects": ["Production APIs", "External Services", "Development APIs"],
    "count": 3
  },
  "status_code": 1,
  "status_message": "Projects retrieved successfully"
}
```

### GET /api/v1/monitors

List all configured monitors across all projects.

**Usage:**

```bash
curl "http://localhost:8080/api/v1/monitors"
```

**Response:**

```json
{
  "data": {
    "monitors": ["User Service Health", "Payment API", "Internal Service"],
    "count": 3
  },
  "status_code": 1,
  "status_message": "Monitors retrieved successfully"
}
```

### GET /api/v1/projects/{project}/monitors

Get all monitors for a specific project.

**Path Parameters:**

- `project` (string): Project name (URL encoded)

**Usage:**

```bash
# Get monitors for "Production APIs" project
curl "http://localhost:8080/api/v1/projects/Production%20APIs/monitors"

# Get monitors for "External Services" project
curl "http://localhost:8080/api/v1/projects/External%20Services/monitors"
```

**Response:**

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

### GET /api/v1/statistics

Get comprehensive system statistics and health metrics.

**Usage:**

```bash
curl "http://localhost:8080/api/v1/statistics"
```

**Response:**

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
      },
      "External Services": {
        "monitor_count": 2,
        "healthy_count": 2,
        "unhealthy_count": 0
      }
    },
    "status_codes": {
      "200": 6,
      "404": 1,
      "500": 1
    },
    "average_latency": "245ms"
  },
  "status_code": 1,
  "status_message": "Statistics retrieved successfully"
}
```

### GET /api/v1/monitors/{monitor}/history

Get historical monitoring data for a specific monitor to create uptime graphs and analyze trends.

**Path Parameters:**

| Parameter | Type   | Description                | Example    |
| --------- | ------ | -------------------------- | ---------- |
| `monitor` | string | Monitor name (URL encoded) | `My%20API` |

**Query Parameters:**

| Parameter | Type   | Description                | Example      |
| --------- | ------ | -------------------------- | ------------ |
| `since`   | string | Start time (epoch seconds) | `1757849600` |
| `until`   | string | End time (epoch seconds)   | `1757849700` |
| `limit`   | int    | Maximum results to return  | `1000`       |

**Usage:**

```bash
# Get last 30 days of data for a monitor
curl "http://localhost:8080/api/v1/monitors/My%20API/history?since=1757849600&limit=1000"

# Get specific time range
curl "http://localhost:8080/api/v1/monitors/Payment%20Service/history?since=1757849600&until=1757850000"
```

**Response:**

```json
{
  "data": {
    "monitor": "My API",
    "results": [
      {
        "project": "Production APIs",
        "name": "My API",
        "url": "https://api.example.com/health",
        "ok": true,
        "status": 200,
        "latency": "150ms",
        "checked_at": "1757851919"
      }
    ],
    "count": 1,
    "filter": {
      "monitor": "My API",
      "since": "1757851919",
      "limit": 1000
    }
  },
  "status_code": 1,
  "status_message": "Monitor history retrieved successfully"
}
```

**Use Cases:**

- Creating uptime status graphs (like Atlassian Status)
- Analyzing performance trends over time
- Generating historical reports
- Building custom dashboards

## 🚨 Alert Conditions Reference

PingPulse supports flexible alert conditions with multiple operators and logical combinations.

### Available Fields

| Field     | Type     | Description          | Examples              |
| --------- | -------- | -------------------- | --------------------- |
| `status`  | int      | HTTP status code     | `200`, `404`, `500`   |
| `latency` | duration | Response time        | `2s`, `500ms`, `1.5s` |
| `ok`      | bool     | Overall check status | `true`, `false`       |

### Comparison Operators

| Operator    | Description           | Example         |
| ----------- | --------------------- | --------------- |
| `==` or `=` | Equal to              | `status == 200` |
| `!=`        | Not equal to          | `status != 200` |
| `>`         | Greater than          | `latency > 2s`  |
| `<`         | Less than             | `latency < 1s`  |
| `>=`        | Greater than or equal | `status >= 500` |
| `<=`        | Less than or equal    | `latency <= 1s` |

### Logical Operators

| Operator | Description                   | Example                          |
| -------- | ----------------------------- | -------------------------------- |
| `AND`    | Both conditions must be true  | `status >= 500 AND latency > 5s` |
| `OR`     | Either condition must be true | `status != 200 OR ok == false`   |

### Example Conditions

```yaml
# Basic conditions
condition: "status != 200"           # Alert on non-200 status
condition: "latency > 2s"            # Alert on slow responses
condition: "ok == false"             # Alert on any failure

# Server error conditions
condition: "status >= 500"           # Alert on 5xx errors
condition: "status >= 400"           # Alert on 4xx and 5xx errors

# Performance conditions
condition: "latency > 1s"            # Alert on responses > 1 second
condition: "latency <= 100ms"        # Alert if too fast (unusual)

# Complex conditions
condition: "status != 200 OR ok == false"                    # Any failure
condition: "status >= 500 AND latency > 5s"                  # Slow server errors
condition: "status >= 400 AND status < 500"                  # Client errors only
condition: "latency > 1s AND latency <= 3s"                  # Moderate slowness

# Multiple OR conditions
condition: "status == 401 OR status == 403"                  # Auth errors
condition: "status >= 500 OR latency > 10s OR ok == false"   # Critical issues
```

### Real-World Examples

```yaml
# E-commerce monitoring
- name: "Payment Gateway Critical"
  condition: "status != 200 OR latency > 3s"

- name: "API Performance Warning"
  condition: "latency > 1s AND latency <= 3s"

- name: "Authentication Issues"
  condition: "status == 401 OR status == 403"

# Infrastructure monitoring
- name: "Service Completely Down"
  condition: "ok == false"

- name: "High Latency Warning"
  condition: "latency > 2s AND status == 200"

- name: "Server Errors"
  condition: "status >= 500"
```

## ❌ Error Handling

### HTTP Status Codes

| Status Code                 | Description        | When It Occurs                        |
| --------------------------- | ------------------ | ------------------------------------- |
| `200 OK`                    | Successful request | Request processed successfully        |
| `400 Bad Request`           | Invalid parameters | Malformed query parameters or request |
| `404 Not Found`             | Resource not found | Invalid endpoint or project not found |
| `500 Internal Server Error` | Server error       | Internal processing error             |

### Error Response Format

All errors follow the standardized response format:

```json
{
  "data": null,
  "status_code": 0,
  "status_message": "Descriptive error message"
}
```

### Common Error Examples

```json
// Invalid project name
{
  "data": null,
  "status_code": 0,
  "status_message": "Project 'NonExistent' not found"
}

// Invalid query parameter
{
  "data": null,
  "status_code": 0,
  "status_message": "Invalid status parameter: must be a valid integer"
}

// Invalid time format
{
  "data": null,
  "status_code": 0,
  "status_message": "Invalid time format for 'since' parameter"
}
```

## 🔒 Security & Production Considerations

### Current Security Status

⚠️ **Development Mode**: PingPulse currently runs in development mode with minimal security:

- ❌ No authentication required
- ❌ No rate limiting implemented
- ❌ CORS headers not configured
- ❌ No API key validation

### Production Recommendations

For production deployments, implement these security measures:

#### 1. API Authentication

```yaml
# Future implementation
security:
  api_key_required: true
  api_keys:
    - key: "prod_key_123"
      permissions: ["read", "write"]
    - key: "readonly_key_456"
      permissions: ["read"]
```

#### 2. Rate Limiting

```yaml
# Recommended limits
rate_limiting:
  requests_per_minute: 60
  burst_limit: 10
  per_ip_limit: 100
```

#### 3. CORS Configuration

```yaml
# For web dashboard access
cors:
  allowed_origins:
    - "https://your-dashboard.com"
    - "http://localhost:3000" # Development
  allowed_methods: ["GET", "POST"]
  allowed_headers: ["Content-Type", "Authorization"]
```

#### 4. HTTPS Only

```bash
# Use reverse proxy (nginx/apache) for HTTPS
# Never expose port 8080 directly to internet
# Use proper SSL certificates
```

### Network Security

```bash
# Firewall configuration
sudo ufw deny 8080/tcp          # Block direct access
sudo ufw allow 80/tcp           # Allow HTTP (redirect to HTTPS)
sudo ufw allow 443/tcp          # Allow HTTPS

# Use reverse proxy
location /api/ {
    proxy_pass http://localhost:8080;
    proxy_set_header Host $host;
    proxy_set_header X-Real-IP $remote_addr;
}
```

## 📊 API Usage Examples

### Monitoring Dashboard Integration

```javascript
// React/JavaScript example
const API_BASE = "http://localhost:8080";

// Get system statistics
async function getStatistics() {
  const response = await fetch(`${API_BASE}/api/v1/statistics`);
  const data = await response.json();

  if (data.status_code === 1) {
    return data.data;
  } else {
    throw new Error(data.status_message);
  }
}

// Get failed monitors
async function getFailedMonitors() {
  const response = await fetch(`${API_BASE}/api/v1/results?ok=false&limit=10`);
  const data = await response.json();
  return data.data.results;
}
```

### Python Integration

```python
import requests
import json

class PingPulseClient:
    def __init__(self, base_url="http://localhost:8080"):
        self.base_url = base_url

    def get_statistics(self):
        response = requests.get(f"{self.base_url}/api/v1/statistics")
        data = response.json()

        if data["status_code"] == 1:
            return data["data"]
        else:
            raise Exception(data["status_message"])

    def get_project_monitors(self, project_name):
        url = f"{self.base_url}/api/v1/projects/{project_name}/monitors"
        response = requests.get(url)
        return response.json()["data"]["monitors"]

# Usage
client = PingPulseClient()
stats = client.get_statistics()
print(f"Total monitors: {stats['total_monitors']}")
```

### Shell/Bash Integration

```bash
#!/bin/bash
# PingPulse monitoring script

API_BASE="http://localhost:8080"

# Check system health
check_health() {
    curl -s "$API_BASE/healthz" | jq -r '.status_message'
}

# Get failed monitors count
get_failed_count() {
    curl -s "$API_BASE/api/v1/results?ok=false" | jq '.data.count'
}

# Alert if too many failures
failed_count=$(get_failed_count)
if [ "$failed_count" -gt 5 ]; then
    echo "ALERT: $failed_count monitors are failing!"
    # Send notification
fi
```

## 🔄 Webhook Integration

### Custom Webhook Alerts

Configure PingPulse to send alerts to your systems:

```yaml
# config.yaml
global_alerts:
  webhook:
    url: "https://your-system.com/webhooks/pingpulse"
    headers:
      Authorization: "Bearer your-webhook-token"
      Content-Type: "application/json"
```

### Webhook Payload Format

```json
{
  "alert_name": "API Down Alert",
  "monitor_name": "User Service Health",
  "project": "Production APIs",
  "status": "CRITICAL",
  "condition": "status != 200",
  "current_status": 500,
  "current_latency": "5.2s",
  "url": "https://api.company.com/users/health",
  "timestamp": "1757851919",
  "message": "Monitor 'User Service Health' is failing: status != 200"
}
```
