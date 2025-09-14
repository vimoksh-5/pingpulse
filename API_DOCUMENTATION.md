# PingPulse API Documentation

PingPulse provides both legacy and enhanced API endpoints for querying monitoring results and statistics.

## Base URL

All API endpoints are available at: `http://localhost:8080`


1. **`GET /healthz`** - Health check with structured response
2. **`GET /api/v1/results`** - Query results with enhanced error handling
3. **`GET /api/v1/projects`** - Project listing with consistent format
4. **`GET /api/v1/monitors`** - Monitor listing with consistent format
5. **`GET /api/v1/projects/{project}/monitors`** - Project monitors with validation
6. **`GET /api/v1/statistics`** - Statistics with structured response
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

## Enhanced API Endpoints (v1)

### GET /api/v1/results

Advanced results query endpoint with multiple filter options.

**Query Parameters:**

- `project` (string): Filter by project name
- `monitor` (string): Filter by monitor name
- `url` (string): Filter by specific URL
- `status` (int): Filter by HTTP status code
- `ok` (bool): Filter by OK status (true/false)
- `since` (RFC3339): Filter results after this timestamp
- `until` (RFC3339): Filter results before this timestamp
- `limit` (int): Limit number of results returned

**Examples:**

```bash
# Get results for specific project
curl "http://localhost:8080/api/v1/results?project=Production%20APIs"

# Get failed results
curl "http://localhost:8080/api/v1/results?ok=false"

# Get results with specific status code
curl "http://localhost:8080/api/v1/results?status=500"

# Get results for specific monitor
curl "http://localhost:8080/api/v1/results?monitor=User%20Service%20Health"

# Get recent results (last hour)
curl "http://localhost:8080/api/v1/results?since=2024-01-01T10:00:00Z&limit=50"

# Complex query: Failed results for specific project in time range
curl "http://localhost:8080/api/v1/results?project=Production%20APIs&ok=false&since=2024-01-01T09:00:00Z&until=2024-01-01T10:00:00Z"
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
        "checked_at": "2024-01-01T10:30:00Z"
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

Returns all available projects.

**Example:**

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

Returns all available monitors.

**Example:**

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

Returns all monitors for a specific project.

**Example:**

```bash
curl "http://localhost:8080/api/v1/projects/Production%20APIs/monitors"
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

Returns comprehensive monitoring statistics.

**Example:**

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

## Enhanced Alert Conditions

PingPulse now supports flexible alert conditions with multiple operators and logical combinations.

### Supported Fields

- `status`: HTTP status code
- `latency`: Response time (supports duration formats like `2s`, `500ms`)
- `ok`: Boolean status of the check

### Supported Operators

- `==` or `=`: Equal to
- `!=`: Not equal to
- `>`: Greater than
- `<`: Less than
- `>=`: Greater than or equal to
- `<=`: Less than or equal to

### Logical Operators

- `AND`: Both conditions must be true
- `OR`: Either condition must be true

### Example Conditions

```yaml
# Simple conditions
condition: "status != 200"
condition: "latency > 2s"
condition: "ok == false"

# Comparison conditions
condition: "status >= 500"
condition: "status < 400"
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

### Backward Compatibility

All existing condition formats are still supported:

- `status != 200`
- `latency > 2s`
- `latency > 3s`
- `latency > 5s`
- `ok == false`

## Error Handling

All API endpoints return appropriate HTTP status codes:

- `200 OK`: Successful request
- `400 Bad Request`: Invalid parameters or malformed request
- `404 Not Found`: Resource not found
- `500 Internal Server Error`: Server error

Error responses follow the standardized format:

```json
{
  "data": null,
  "status_code": 0,
  "status_message": "Invalid path. Use /api/v1/projects/{project}/monitors"
}
```

## Rate Limiting

Currently, no rate limiting is implemented. Consider implementing rate limiting for production deployments.

## Authentication

Currently, no authentication is required for API access. Consider implementing API key authentication for production deployments.

## CORS

CORS headers are not currently set. Add CORS middleware if accessing from web applications.
