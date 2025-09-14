# PingPulse Deployment Guide

This guide covers production deployment options for PingPulse.

## 🚀 Production Deployment

### Backend Deployment

#### Option 1: Binary Deployment

```bash
# Build optimized binary
go build -ldflags="-s -w" -o pingpulse ./cmd/pingpulse

# Copy to production server
scp pingpulse config.yaml user@server:/opt/pingpulse/

# Run with systemd (recommended)
sudo systemctl enable pingpulse
sudo systemctl start pingpulse
```

#### Option 2: Docker Deployment

```dockerfile
# Dockerfile
FROM golang:1.21-alpine AS builder
WORKDIR /app
COPY . .
RUN go build -ldflags="-s -w" -o pingpulse ./cmd/pingpulse

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=builder /app/pingpulse .
COPY config.yaml .
EXPOSE 8080
CMD ["./pingpulse"]
```

```bash
# Build and run
docker build -t pingpulse .
docker run -d -p 8080:8080 -v $(pwd)/config.yaml:/root/config.yaml pingpulse
```

#### Option 3: Docker Compose

```yaml
# docker-compose.yml
version: "3.8"
services:
  pingpulse:
    build: .
    ports:
      - "8080:8080"
    volumes:
      - ./config.yaml:/root/config.yaml
      - ./data:/root/data
    environment:
      - GMAIL_USERNAME=${GMAIL_USERNAME}
      - GMAIL_APP_PASSWORD=${GMAIL_APP_PASSWORD}
      - SLACK_WEBHOOK_URL=${SLACK_WEBHOOK_URL}
    restart: unless-stopped
```

### Dashboard Deployment

#### Option 1: Static Hosting

```bash
cd dashboard
npm install
npm run build

# Deploy to static hosting providers:
# - Netlify: drag and drop build/ folder
# - Vercel: vercel --prod
# - AWS S3: aws s3 sync build/ s3://your-bucket
# - GitHub Pages: push build/ to gh-pages branch
```

#### Option 2: Nginx Reverse Proxy

```nginx
# /etc/nginx/sites-available/pingpulse
server {
    listen 80;
    server_name your-domain.com;

    # Dashboard
    location / {
        root /var/www/pingpulse-dashboard;
        try_files $uri $uri/ /index.html;
    }

    # API proxy
    location /api/ {
        proxy_pass http://localhost:8080;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
    }

    location /healthz {
        proxy_pass http://localhost:8080;
    }
}
```

#### Option 3: Docker with Nginx

```dockerfile
# dashboard/Dockerfile
FROM node:18-alpine AS builder
WORKDIR /app
COPY package*.json ./
RUN npm ci
COPY . .
RUN npm run build

FROM nginx:alpine
COPY --from=builder /app/build /usr/share/nginx/html
COPY nginx.conf /etc/nginx/nginx.conf
EXPOSE 80
```

## 🔒 Security Configuration

### SSL/TLS Setup

```bash
# Using Let's Encrypt with Certbot
sudo apt install certbot python3-certbot-nginx
sudo certbot --nginx -d your-domain.com
```

### Firewall Configuration

```bash
# UFW (Ubuntu)
sudo ufw allow 22/tcp    # SSH
sudo ufw allow 80/tcp    # HTTP
sudo ufw allow 443/tcp   # HTTPS
sudo ufw enable

# Block direct access to backend port
sudo ufw deny 8080/tcp
```

### Environment Variables Security

```bash
# Use systemd environment files
sudo mkdir -p /etc/pingpulse
sudo chmod 700 /etc/pingpulse

# Create secure env file
sudo tee /etc/pingpulse/environment << EOF
GMAIL_USERNAME=your-email@gmail.com
GMAIL_APP_PASSWORD=your-app-password
SLACK_WEBHOOK_URL=https://hooks.slack.com/services/YOUR/SLACK/WEBHOOK
API_TOKEN=your-secure-token
EOF

sudo chmod 600 /etc/pingpulse/environment
```

## 📊 Monitoring & Logging

### Systemd Service Configuration

```ini
# /etc/systemd/system/pingpulse.service
[Unit]
Description=PingPulse Monitoring Service
After=network.target

[Service]
Type=simple
User=pingpulse
Group=pingpulse
WorkingDirectory=/opt/pingpulse
ExecStart=/opt/pingpulse/pingpulse
EnvironmentFile=/etc/pingpulse/environment
Restart=always
RestartSec=5
StandardOutput=journal
StandardError=journal

[Install]
WantedBy=multi-user.target
```

### Log Management

```bash
# Configure log rotation
sudo tee /etc/logrotate.d/pingpulse << EOF
/var/log/pingpulse/*.log {
    daily
    rotate 30
    compress
    delaycompress
    missingok
    notifempty
    create 644 pingpulse pingpulse
    postrotate
        systemctl reload pingpulse
    endscript
}
EOF
```

### Health Monitoring

```bash
# Add to crontab for health checks
*/5 * * * * curl -f http://localhost:8080/healthz || systemctl restart pingpulse
```

## 🌐 Load Balancing & High Availability

### Multiple Instances with HAProxy

```
# /etc/haproxy/haproxy.cfg
global
    daemon

defaults
    mode http
    timeout connect 5000ms
    timeout client 50000ms
    timeout server 50000ms

frontend pingpulse_frontend
    bind *:80
    default_backend pingpulse_backend

backend pingpulse_backend
    balance roundrobin
    server pingpulse1 10.0.1.10:8080 check
    server pingpulse2 10.0.1.11:8080 check
    server pingpulse3 10.0.1.12:8080 check
```

### Kubernetes Deployment

```yaml
# k8s/deployment.yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: pingpulse
spec:
  replicas: 3
  selector:
    matchLabels:
      app: pingpulse
  template:
    metadata:
      labels:
        app: pingpulse
    spec:
      containers:
        - name: pingpulse
          image: pingpulse:latest
          ports:
            - containerPort: 8080
          env:
            - name: GMAIL_USERNAME
              valueFrom:
                secretKeyRef:
                  name: pingpulse-secrets
                  key: gmail-username
          volumeMounts:
            - name: config
              mountPath: /root/config.yaml
              subPath: config.yaml
      volumes:
        - name: config
          configMap:
            name: pingpulse-config

---
apiVersion: v1
kind: Service
metadata:
  name: pingpulse-service
spec:
  selector:
    app: pingpulse
  ports:
    - port: 80
      targetPort: 8080
  type: LoadBalancer
```

## 🔧 Performance Optimization

### Backend Optimization

```yaml
# config.yaml - Production settings
default:
  interval: 60s
  timeout: 10s
  max_retries: 3
  jitter: 10s

# Optimize for high-load scenarios
server:
  read_timeout: 30s
  write_timeout: 30s
  max_connections: 1000
```

### Database Considerations

```go
// For large deployments, consider external storage
// PostgreSQL, MySQL, or Redis for results storage
// This requires custom implementation or future updates
```

## 📈 Scaling Recommendations

### Small Deployment (< 100 monitors)

- Single instance deployment
- In-memory storage
- Basic alerting

### Medium Deployment (100-1000 monitors)

- Load balanced instances
- Shared configuration
- Enhanced monitoring

### Large Deployment (> 1000 monitors)

- Kubernetes cluster
- External database
- Message queue for alerts
- Distributed monitoring

## 🚨 Alert Channel Setup

### Email Configuration (Gmail)

1. Enable 2-Factor Authentication
2. Generate App Password: Google Account → Security → App passwords
3. Use App Password in environment variables

### Slack Configuration

1. Go to Slack workspace settings
2. Create new Incoming Webhook
3. Copy webhook URL to environment variables

### Custom Webhooks

Configure any HTTP endpoint for custom integrations:

```yaml
global_alerts:
  webhook:
    url: "https://your-webhook-endpoint.com/alerts"
    headers:
      Authorization: "Bearer ${WEBHOOK_TOKEN}"
      Content-Type: "application/json"
```

## 🔍 Troubleshooting Production Issues

### Common Deployment Issues

1. **Port conflicts**: Ensure port 8080 is available
2. **Permission issues**: Run as non-root user
3. **Config file location**: Verify config.yaml path
4. **Environment variables**: Check all required vars are set

### Performance Issues

1. **High memory usage**: Reduce monitoring frequency
2. **Slow responses**: Optimize timeout settings
3. **Alert spam**: Increase cooldown periods

### Monitoring the Monitor

Set up external monitoring for PingPulse itself:

```bash
# External health check
curl -f https://your-pingpulse-domain.com/healthz
```

## 📋 Production Checklist

- [ ] SSL/TLS certificates configured
- [ ] Firewall rules in place
- [ ] Environment variables secured
- [ ] Log rotation configured
- [ ] Backup strategy implemented
- [ ] Monitoring alerts tested
- [ ] Health checks configured
- [ ] Documentation updated
- [ ] Team access configured
- [ ] Incident response plan ready

## 🆘 Emergency Procedures

### Service Recovery

```bash
# Quick restart
sudo systemctl restart pingpulse

# Check logs
sudo journalctl -u pingpulse -f

# Rollback deployment
sudo systemctl stop pingpulse
sudo cp /opt/pingpulse/backup/pingpulse /opt/pingpulse/
sudo systemctl start pingpulse
```

### Data Recovery

```bash
# Backup monitoring data (if using file storage)
sudo cp -r /opt/pingpulse/data /backup/pingpulse-$(date +%Y%m%d)
```
