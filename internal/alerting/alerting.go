package alerting

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/smtp"
	"strings"
	"time"

	"pingpulse/internal/config"
	"pingpulse/internal/monitor"
)

type AlertManager struct {
	config          *config.Config
	client          *http.Client
	conditionParser *ConditionParser
}

func NewAlertManager(cfg *config.Config) *AlertManager {
	return &AlertManager{
		config:          cfg,
		client:          &http.Client{Timeout: 10 * time.Second},
		conditionParser: NewConditionParser(),
	}
}

func (am *AlertManager) SendAlert(result monitor.Result, alert monitor.Alert) {
	// Check if alert condition is met
	if !am.conditionParser.EvaluateCondition(result, alert.Condition) {
		return
	}

	// Send to configured channels
	for _, channel := range alert.Channels {
		switch channel {
		case "email":
			am.sendEmailAlert(result, alert)
		case "slack":
			am.sendSlackAlert(result, alert)
		case "webhook":
			am.sendWebhookAlert(result, alert)
		}
	}
}

// ValidateCondition validates an alert condition syntax
func (am *AlertManager) ValidateCondition(condition string) error {
	return am.conditionParser.ValidateCondition(condition)
}

func (am *AlertManager) sendEmailAlert(result monitor.Result, alert monitor.Alert) {
	if am.config.GlobalAlerts.Email.SMTPHost == "" {
		return
	}

	subject := fmt.Sprintf("🚨 PingPulse Alert: %s", result.Name)
	body := am.formatEmailBody(result)

	// Create message
	msg := fmt.Sprintf("From: %s\r\nTo: %s\r\nSubject: %s\r\n\r\n%s",
		am.config.GlobalAlerts.Email.From,
		strings.Join(alert.Recipients, ", "),
		subject,
		body,
	)

	// Send email
	addr := fmt.Sprintf("%s:%d", am.config.GlobalAlerts.Email.SMTPHost, am.config.GlobalAlerts.Email.SMTPPort)
	auth := smtp.PlainAuth("", am.config.GlobalAlerts.Email.Username, am.config.GlobalAlerts.Email.Password, am.config.GlobalAlerts.Email.SMTPHost)

	err := smtp.SendMail(addr, auth, am.config.GlobalAlerts.Email.From, alert.Recipients, []byte(msg))
	if err != nil {
		fmt.Printf("Failed to send email alert: %v\n", err)
	} else {
		fmt.Printf("📧 Email alert sent for %s\n", result.Name)
	}
}

func (am *AlertManager) sendSlackAlert(result monitor.Result, alert monitor.Alert) {
	if am.config.GlobalAlerts.Slack.WebhookURL == "" {
		return
	}

	// Create Slack message
	slackMsg := map[string]interface{}{
		"text": fmt.Sprintf("🚨 *PingPulse Alert: %s*", result.Name),
		"attachments": []map[string]interface{}{
			{
				"color": am.getAlertColor(result),
				"fields": []map[string]interface{}{
					{"title": "Project", "value": result.Project, "short": true},
					{"title": "Monitor", "value": result.Name, "short": true},
					{"title": "URL", "value": result.URL, "short": false},
					{"title": "Status", "value": fmt.Sprintf("%d", result.Status), "short": true},
					{"title": "Latency", "value": result.Latency.String(), "short": true},
					{"title": "Time", "value": time.Unix(int64(result.CheckedAt), 0).Format("2006-01-02 15:04:05"), "short": true},
				},
			},
		},
	}

	jsonData, _ := json.Marshal(slackMsg)
	resp, err := am.client.Post(am.config.GlobalAlerts.Slack.WebhookURL, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		fmt.Printf("Failed to send Slack alert: %v\n", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode == 200 {
		fmt.Printf("💬 Slack alert sent for %s\n", result.Name)
	} else {
		fmt.Printf("Failed to send Slack alert: HTTP %d\n", resp.StatusCode)
	}
}

func (am *AlertManager) sendWebhookAlert(result monitor.Result, alert monitor.Alert) {
	webhookURL := alert.WebhookURL
	if webhookURL == "" {
		webhookURL = am.config.GlobalAlerts.Webhook.URL
	}
	if webhookURL == "" {
		return
	}

	// Create webhook payload
	payload := map[string]interface{}{
		"alert":   alert.Name,
		"project": result.Project,
		"monitor": result.Name,
		"url":     result.URL,
		"status":  result.Status,
		"ok":      result.OK,
		"latency": result.Latency.String(),
		"error":   result.Error,
		"time":    time.Unix(int64(result.CheckedAt), 0).Format(time.RFC3339),
	}

	jsonData, _ := json.Marshal(payload)
	req, _ := http.NewRequest("POST", webhookURL, bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")

	// Add custom headers
	for key, value := range am.config.GlobalAlerts.Webhook.Headers {
		req.Header.Set(key, value)
	}

	resp, err := am.client.Do(req)
	if err != nil {
		fmt.Printf("Failed to send webhook alert: %v\n", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 200 && resp.StatusCode < 300 {
		fmt.Printf("🔗 Webhook alert sent for %s\n", result.Name)
	} else {
		fmt.Printf("Failed to send webhook alert: HTTP %d\n", resp.StatusCode)
	}
}

func (am *AlertManager) formatEmailBody(result monitor.Result) string {
	status := "✅ OK"
	if !result.OK {
		status = "❌ FAILED"
	}

	return fmt.Sprintf(`
PingPulse Monitoring Alert

Project: %s
Monitor: %s
URL: %s
Status: %s (%d)
Latency: %s
Time: %s
Error: %s

Automated alert from PingPulse monitoring.`,
		result.Project,
		result.Name,
		result.URL,
		status,
		result.Status,
		result.Latency.String(),
		time.Unix(int64(result.CheckedAt), 0).Format("2006-01-02 15:04:05"),
		result.Error)
}

func (am *AlertManager) getAlertColor(result monitor.Result) string {
	if result.OK {
		return "good"
	}
	if result.Status >= 500 {
		return "danger"
	}
	return "warning"
}
