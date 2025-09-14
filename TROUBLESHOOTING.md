# PingPulse Troubleshooting Guide

This guide helps you diagnose and fix common issues with PingPulse.

## 🚨 Quick Diagnostics

### Health Check

```bash
# Check if PingPulse is running
curl http://localhost:8080/healthz

# Expected response:
{
  "data": {"status": "healthy"},
  "status_code": 1,
  "status_message": "Service is healthy"
}
```

### Service Status

```bash
# Check service status (systemd)
sudo systemctl status pingpulse

# Check logs
sudo journalctl -u pingpulse -f

# Check process
ps aux | grep pingpulse
```

## 🔧 Common Issues

### 1. Dashboard Won't Load

#### Symptoms

- Dashboard shows "Unable to connect to backend"
- White screen or loading spinner
- Network errors in browser console

#### Diagnosis

```bash
# Check if backend is running
curl http://localhost:8080/healthz

# Check if port 8080 is open
netstat -tlnp | grep 8080

# Check CORS headers
curl -H "Origin: http://localhost:3000" -v http://localhost:8080/api/v1/statistics
```

#### Solutions

```bash
# 1. Start the backend
./pingpulse

# 2. Check config file exists
ls -la config.yaml

# 3. Verify network connectivity
ping localhost

# 4. Check firewall rules
sudo ufw status
```

### 2. No Monitoring Data

#### Symptoms

- API returns empty results
- Statistics show 0 monitors
- Dashboard shows "No monitors found"

#### Diagnosis

```bash
# Check config file syntax
./pingpulse --validate-config

# Check if monitors are configured
grep -A 5 "monitors:" config.yaml

# Check API response
curl http://localhost:8080/api/v1/monitors
```

#### Solutions

```yaml
# 1. Verify config.yaml structure
projects:
  - name: "Test Project"
    monitors:
      - name: "Test Monitor"
        url: "https://httpbin.org/status/200"
        interval: 30s

# 2. Check URL accessibility
curl -v https://httpbin.org/status/200

# 3. Restart service
sudo systemctl restart pingpulse
```

### 3. Alerts Not Working

#### Symptoms

- No email/Slack notifications
- Alert conditions met but no alerts sent
- Error messages in logs about SMTP/Slack

#### Diagnosis

```bash
# Check alert configuration
grep -A 10 "global_alerts:" config.yaml

# Check environment variables
echo $GMAIL_USERNAME
echo $SLACK_WEBHOOK_URL

# Test SMTP connection
telnet smtp.gmail.com 587
```

#### Solutions

##### Email Alerts (Gmail)

```bash
# 1. Enable 2FA on Gmail account
# 2. Generate App Password (not regular password)
# 3. Set environment variables
export GMAIL_USERNAME="your-email@gmail.com"
export GMAIL_APP_PASSWORD="your-16-char-app-password"

# 4. Test configuration
curl -X POST http://localhost:8080/test/email
```

##### Slack Alerts

```bash
# 1. Create Slack webhook
# 2. Set webhook URL
export SLACK_WEBHOOK_URL="https://hooks.slack.com/services/YOUR/SLACK/WEBHOOK"

# 3. Test webhook
curl -X POST -H 'Content-type: application/json' \
  --data '{"text":"Test from PingPulse"}' \
  $SLACK_WEBHOOK_URL
```

### 4. High Memory Usage

#### Symptoms

- System running out of memory
- PingPulse process killed by OOM killer
- Slow response times

#### Diagnosis

```bash
# Check memory usage
ps aux | grep pingpulse
free -h
top -p $(pgrep pingpulse)

# Check number of monitors
curl http://localhost:8080/api/v1/statistics | jq '.data.total_monitors'
```

#### Solutions

```yaml
# 1. Reduce monitoring frequency
default:
  interval: 300s # 5 minutes instead of 60s

# 2. Limit result retention (if implemented)
storage:
  max_results: 1000
  cleanup_interval: 24h

# 3. Reduce concurrent monitors
monitors:
  - name: "High Frequency Monitor"
    interval: 60s # Only for critical services
```

### 5. SSL/TLS Certificate Issues

#### Symptoms

- "Certificate verification failed" errors
- Monitors failing with SSL errors
- HTTPS endpoints not working

#### Diagnosis

```bash
# Test SSL certificate
openssl s_client -connect api.example.com:443 -servername api.example.com

# Check certificate expiry
echo | openssl s_client -connect api.example.com:443 2>/dev/null | openssl x509 -noout -dates
```

#### Solutions

```yaml
# 1. Disable SSL verification (not recommended for production)
monitors:
  - name: "Test Monitor"
    url: "https://self-signed.example.com"
    ssl_verify: false

# 2. Add custom CA certificates
# Copy certificates to /etc/ssl/certs/ and update ca-certificates

# 3. Use HTTP instead of HTTPS for internal services
monitors:
  - name: "Internal Service"
    url: "http://internal-api.local/health"
```

### 6. Authentication Issues

#### Symptoms

- 401/403 errors from monitored APIs
- "Authentication failed" messages
- OAuth token expired errors

#### Diagnosis

```bash
# Test authentication manually
curl -H "Authorization: Bearer $API_TOKEN" https://api.example.com/health

# Check token expiry (JWT)
echo $API_TOKEN | cut -d. -f2 | base64 -d | jq .exp
```

#### Solutions

```yaml
# 1. Verify credentials
auth:
  type: "bearer"
  token: "${API_TOKEN}"  # Ensure env var is set

# 2. Use correct auth type
auth:
  type: "basic"
  username: "${API_USERNAME}"
  password: "${API_PASSWORD}"

# 3. Refresh OAuth tokens (manual process)
# Implement token refresh logic or update tokens regularly
```

### 7. Performance Issues

#### Symptoms

- Slow dashboard loading
- API timeouts
- High CPU usage

#### Diagnosis

```bash
# Check API response times
time curl http://localhost:8080/api/v1/statistics

# Monitor CPU usage
top -p $(pgrep pingpulse)

# Check number of concurrent monitors
curl http://localhost:8080/api/v1/statistics | jq '.data'
```

#### Solutions

```yaml
# 1. Optimize monitoring intervals
default:
  interval: 120s # Increase from 60s
  timeout: 5s # Reduce from 10s

# 2. Reduce concurrent operations
monitors:
  - name: "Critical Service"
    interval: 30s # Keep short for critical
  - name: "Less Critical"
    interval: 300s # Longer for less critical

# 3. Implement caching (future enhancement)
```

## 🐛 Debug Mode

### Enable Verbose Logging

```bash
# Run with debug flag
./pingpulse --debug

# Or set environment variable
export PINGPULSE_DEBUG=true
./pingpulse
```

### Log Analysis

```bash
# Follow logs in real-time
sudo journalctl -u pingpulse -f

# Search for errors
sudo journalctl -u pingpulse | grep -i error

# Check specific time range
sudo journalctl -u pingpulse --since "2024-01-01 10:00:00" --until "2024-01-01 11:00:00"
```

## 🔍 Network Diagnostics

### Connectivity Tests

```bash
# Test DNS resolution
nslookup api.example.com

# Test connectivity
telnet api.example.com 443

# Test HTTP response
curl -v https://api.example.com/health

# Check routing
traceroute api.example.com
```

### Firewall Issues

```bash
# Check UFW status
sudo ufw status verbose

# Check iptables rules
sudo iptables -L

# Test port accessibility
nc -zv localhost 8080
```

## 📊 Performance Monitoring

### System Resources

```bash
# Memory usage
free -h
cat /proc/meminfo

# CPU usage
top
htop

# Disk usage
df -h
du -sh /opt/pingpulse/

# Network connections
netstat -tulpn | grep pingpulse
```

### Application Metrics

```bash
# Monitor count and status
curl http://localhost:8080/api/v1/statistics | jq '.'

# Check specific monitor results
curl "http://localhost:8080/api/v1/results?monitor=YourMonitorName&limit=10"

# Performance over time
while true; do
  curl -s http://localhost:8080/api/v1/statistics | jq '.data.average_latency'
  sleep 60
done
```

## 🚑 Emergency Recovery

### Service Recovery

```bash
# Quick restart
sudo systemctl restart pingpulse

# Force kill and restart
sudo pkill -9 pingpulse
sudo systemctl start pingpulse

# Rollback to previous version
sudo systemctl stop pingpulse
sudo cp /backup/pingpulse /opt/pingpulse/
sudo systemctl start pingpulse
```

### Configuration Recovery

```bash
# Restore from backup
sudo cp /backup/config.yaml /opt/pingpulse/

# Reset to minimal config
cat > config.yaml << EOF
default:
  interval: 60s
  timeout: 5s

projects:
  - name: "Emergency Test"
    monitors:
      - name: "Health Check"
        url: "https://httpbin.org/status/200"
EOF
```

### Data Recovery

```bash
# If using file-based storage
sudo cp -r /backup/data/ /opt/pingpulse/

# Clear corrupted data
sudo rm -rf /opt/pingpulse/data/*
sudo systemctl restart pingpulse
```

## 📋 Diagnostic Checklist

When reporting issues, include:

- [ ] PingPulse version: `./pingpulse --version`
- [ ] Operating system: `uname -a`
- [ ] Config file (sanitized): Remove sensitive data
- [ ] Error logs: `sudo journalctl -u pingpulse --since "1 hour ago"`
- [ ] System resources: `free -h && df -h`
- [ ] Network connectivity: Test URLs manually
- [ ] Environment variables: Verify all required vars are set

## 🆘 Getting Help

### Before Asking for Help

1. Check this troubleshooting guide
2. Search existing GitHub issues
3. Enable debug mode and collect logs
4. Test with minimal configuration
5. Verify system requirements

### Reporting Issues

Include in your issue report:

- Clear description of the problem
- Steps to reproduce
- Expected vs actual behavior
- System information
- Relevant logs and error messages
- Configuration file (sanitized)

### Community Support

- 🐛 **GitHub Issues**: [Report bugs and issues](https://github.com/vimoksh-5/pingpulse/issues)
- 💬 **Discussions**: [Ask questions and share ideas](https://github.com/vimoksh-5/pingpulse/discussions)

## 🔧 Advanced Debugging

### Memory Profiling

```bash
# Enable Go profiling (if compiled with debug info)
go tool pprof http://localhost:8080/debug/pprof/heap

# Monitor memory over time
while true; do
  ps -o pid,vsz,rss,comm -p $(pgrep pingpulse)
  sleep 10
done
```

### Network Debugging

```bash
# Capture network traffic
sudo tcpdump -i any -w pingpulse.pcap host api.example.com

# Analyze with Wireshark or tcpdump
tcpdump -r pingpulse.pcap -A
```

### Configuration Validation

```bash
# Validate YAML syntax
python3 -c "import yaml; yaml.safe_load(open('config.yaml'))"

# Check for common issues
grep -n "interval.*[0-9]s" config.yaml
grep -n "timeout.*[0-9]s" config.yaml
```
