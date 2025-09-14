package monitor

import (
	"bytes"
	"context"
	"crypto/tls"
	"encoding/json"
	"io"
	"log"
	"net"
	"net/http"
	"time"
)

type Auth struct {
	Type         string
	Token        string
	Username     string
	Password     string
	ClientID     string
	ClientSecret string
	TokenURL     string
	APIKey       string
	APIKeyHeader string
}

type BodyValidation struct {
	Path     string
	Expected interface{}
}

type Alert struct {
	Name       string
	Condition  string
	Channels   []string
	Recipients []string
	WebhookURL string
	Cooldown   time.Duration
}

type Target struct {
	Project         string            // which project this monitor belongs to
	Name            string            // a friendly name for the monitor
	URL             string
	Every           time.Duration     // interval between probes
	Timeout         time.Duration     // HTTP timeout per attempt
	Method          string            // GET, POST, etc.; default to GET
	Headers         map[string]string // optional HTTP headers (e.g., auth)
	Body            interface{}       // optional request body (for POST/PUT)
	MaxRetries      int               // number of retries on failure
	Auth            *Auth             // authentication configuration
	ExpectedStatus  *int              // expected HTTP status code
	ExpectedLatency *time.Duration    // expected response time threshold
	SSLVerify       *bool             // whether to verify SSL certificates
	FollowRedirects *bool             // whether to follow redirects
	BodyValidation  *BodyValidation   // response body validation
	Alerts          []Alert           // alerting configuration
	ConfigOrder     int               // Order from config.yaml for consistent sorting
}

// EpochTime represents Unix timestamp that marshals as integer
type EpochTime int64

// MarshalJSON ensures EpochTime is serialized as integer
func (e EpochTime) MarshalJSON() ([]byte, error) {
	return json.Marshal(int64(e))
}

type Result struct {
	Project     string        `json:"project"`
	Name        string        `json:"name"`
	URL         string        `json:"url"`
	OK          bool          `json:"ok"`
	Status      int           `json:"status"`
	Latency     time.Duration `json:"latency"`
	Error       string        `json:"error,omitempty"`
	CheckedAt   EpochTime     `json:"checked_at"`   // Unix timestamp in seconds
	ConfigOrder int           `json:"config_order"` // Order from config.yaml
}

type ResultSink interface {
	Append(Result)
	Latest(url string, n int) []Result
	All() map[string][]Result
}

type Manager struct {
	sink         ResultSink
	http         *http.Client
	alertManager interface {
		SendAlert(Result, Alert)
	}
}

func NewManager(sink ResultSink) *Manager {
	transport := &http.Transport{
		ResponseHeaderTimeout: 5 * time.Second,
		TLSClientConfig:       &tls.Config{MinVersion: tls.VersionTLS12},
		DialContext: (&net.Dialer{
			Timeout: 5 * time.Second,
		}).DialContext,
		Proxy: http.ProxyFromEnvironment,
	}
	return &Manager{
		sink: sink,
		http: &http.Client{Transport: transport},
	}
}

func (m *Manager) SetAlertManager(alertManager interface {
	SendAlert(Result, Alert)
}) {
	m.alertManager = alertManager
}

func (m *Manager) Start(ctx context.Context, t Target) {
	ticker := time.NewTicker(t.Every)
	go func() {
		defer ticker.Stop()
		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				m.checkOnce(ctx, t)
			}
		}
	}()
	// run immediately on start
	m.checkOnce(ctx, t)
}

// applyAuth applies authentication to the HTTP request
func (m *Manager) applyAuth(req *http.Request, auth *Auth) {
	if auth == nil {
		return
	}

	switch auth.Type {
	case "bearer":
		req.Header.Set("Authorization", "Bearer "+auth.Token)
	case "basic":
		req.SetBasicAuth(auth.Username, auth.Password)
	case "apikey":
		header := auth.APIKeyHeader
		if header == "" {
			header = "X-API-Key"
		}
		req.Header.Set(header, auth.APIKey)
	case "oauth2":
		// For OAuth2, you would typically get a token first
		// This is a simplified version - in production you'd want to cache tokens
		if auth.Token != "" {
			req.Header.Set("Authorization", "Bearer "+auth.Token)
		}
	}
}

// validateResponseBody validates the response body against the expected value
func (m *Manager) validateResponseBody(body []byte, validation *BodyValidation) bool {
	if validation == nil {
		return true
	}

	// Simple JSON path validation (you might want to use a proper JSONPath library)
	var data map[string]interface{}
	if err := json.Unmarshal(body, &data); err != nil {
		return false
	}

	// Simple path validation for common cases
	// In production, use a proper JSONPath library like github.com/PaesslerAG/jsonpath
	switch validation.Path {
	case "$.status":
		if status, ok := data["status"].(string); ok {
			return status == validation.Expected
		}
	case "$.healthy":
		if healthy, ok := data["healthy"].(bool); ok {
			return healthy == validation.Expected
		}
	}

	return true
}

func (m *Manager) checkOnce(ctx context.Context, t Target) {
	// Default method is GET
	method := t.Method
	if method == "" {
		method = http.MethodGet
	}

	var body io.Reader
	if t.Body != nil {
		// marshal body to JSON (you could support other formats too)
		b, _ := json.Marshal(t.Body)
		body = bytes.NewReader(b)
	}

	// simple retry loop
	var lastErr error
	for attempt := 0; attempt <= t.MaxRetries; attempt++ {
		reqCtx, cancel := context.WithTimeout(ctx, t.Timeout)
		start := time.Now()
		req, _ := http.NewRequestWithContext(reqCtx, method, t.URL, body)

		// apply authentication
		m.applyAuth(req, t.Auth)

		// apply headers (Content-Type if body present)
		for k, v := range t.Headers {
			req.Header.Set(k, v)
		}
		if t.Body != nil {
			req.Header.Set("Content-Type", "application/json")
		}

		// send the request
		res, err := m.http.Do(req)
		cancel()
		r := Result{
			Project:     t.Project,
			Name:        t.Name,
			URL:         t.URL,
			CheckedAt:   EpochTime(time.Now().Unix()),
			Latency:     time.Since(start),
			ConfigOrder: t.ConfigOrder,
		}

		if err != nil {
			lastErr = err
			if attempt == t.MaxRetries {
				r.OK = false
				r.Error = lastErr.Error()
				m.sink.Append(r)
				log.Printf("Monitor %s failed: %v", t.Name, lastErr.Error())
				return
			}
			continue // retry
		}

		// success
		defer res.Body.Close()
		r.Status = res.StatusCode

		// Check expected status
		if t.ExpectedStatus != nil {
			r.OK = (res.StatusCode == *t.ExpectedStatus)
		} else {
			r.OK = (res.StatusCode >= 200 && res.StatusCode < 400)
		}

		// Check expected latency
		if t.ExpectedLatency != nil && r.Latency > *t.ExpectedLatency {
			r.OK = false
			r.Error = "Response time exceeded threshold"
		}

		// Validate response body if configured
		if t.BodyValidation != nil && r.OK {
			bodyBytes, _ := io.ReadAll(res.Body)
			if !m.validateResponseBody(bodyBytes, t.BodyValidation) {
				r.OK = false
				r.Error = "Response body validation failed"
			}
		}

		m.sink.Append(r)
		log.Printf("Monitor %s: status=%d, latency=%v, ok=%v", t.Name, res.StatusCode, r.Latency, r.OK)

		// Send alerts if configured
		if m.alertManager != nil {
			for _, alert := range t.Alerts {
				m.alertManager.SendAlert(r, alert)
			}
		}

		return
	}
}
