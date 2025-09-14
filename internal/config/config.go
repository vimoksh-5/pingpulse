package config

import (
	"os"
	"strings"
	"time"

	"gopkg.in/yaml.v3"
)

// Duration wraps time.Duration to handle YAML durations like "15s".
type Duration struct{ time.Duration }

func (d *Duration) UnmarshalYAML(value *yaml.Node) error {
	var s string
	if err := value.Decode(&s); err != nil {
		return err
	}
	dur, err := time.ParseDuration(s)
	if err != nil {
		return err
	}
	d.Duration = dur
	return nil
}

type Defaults struct {
	Interval   Duration `yaml:"interval"`
	Timeout    Duration `yaml:"timeout"`
	MaxRetries int      `yaml:"max_retries"`
	Jitter     Duration `yaml:"jitter"`
	UserAgent  string   `yaml:"user_agent"`
}

type Auth struct {
	Type         string `yaml:"type,omitempty"` // bearer, basic, oauth2, apikey
	Token        string `yaml:"token,omitempty"`
	Username     string `yaml:"username,omitempty"`
	Password     string `yaml:"password,omitempty"`
	ClientID     string `yaml:"client_id,omitempty"`
	ClientSecret string `yaml:"client_secret,omitempty"`
	TokenURL     string `yaml:"token_url,omitempty"`
	APIKey       string `yaml:"api_key,omitempty"`
	APIKeyHeader string `yaml:"api_key_header,omitempty"` // X-API-Key, Authorization, etc.
}

type Alert struct {
	Name       string   `yaml:"name"`
	Condition  string   `yaml:"condition"` // latency > 5s, status != 200, etc.
	Channels   []string `yaml:"channels"`  // email, slack, webhook
	Recipients []string `yaml:"recipients,omitempty"`
	WebhookURL string   `yaml:"webhook_url,omitempty"`
	Cooldown   Duration `yaml:"cooldown,omitempty"`
}

type Monitor struct {
	Name            string            `yaml:"name"`
	URL             string            `yaml:"url"`
	Interval        *Duration         `yaml:"interval,omitempty"`
	Timeout         *Duration         `yaml:"timeout,omitempty"`
	MaxRetries      *int              `yaml:"max_retries,omitempty"`
	Method          string            `yaml:"method,omitempty"`
	Headers         map[string]string `yaml:"headers,omitempty"`
	Body            interface{}       `yaml:"body,omitempty"`
	Auth            *Auth             `yaml:"auth,omitempty"`
	ExpectedStatus  *int              `yaml:"expected_status,omitempty"`
	ExpectedLatency *Duration         `yaml:"expected_latency,omitempty"`
	SSLVerify       *bool             `yaml:"ssl_verify,omitempty"`
	FollowRedirects *bool             `yaml:"follow_redirects,omitempty"`
	BodyValidation  *BodyValidation   `yaml:"body_validation,omitempty"`
	Alerts          []Alert           `yaml:"alerts,omitempty"`
}

type BodyValidation struct {
	Path     string      `yaml:"path"`     // JSONPath like "$.status"
	Expected interface{} `yaml:"expected"` // Expected value
}

type Project struct {
	Name     string    `yaml:"name"`
	ApiKey   string    `yaml:"api_key"`
	Monitors []Monitor `yaml:"monitors"`
}

type GlobalAlerts struct {
	Email   EmailConfig   `yaml:"email,omitempty"`
	Slack   SlackConfig   `yaml:"slack,omitempty"`
	Webhook WebhookConfig `yaml:"webhook,omitempty"`
}

type EmailConfig struct {
	SMTPHost string `yaml:"smtp_host"`
	SMTPPort int    `yaml:"smtp_port"`
	Username string `yaml:"username"`
	Password string `yaml:"password"`
	From     string `yaml:"from"`
	UseTLS   bool   `yaml:"use_tls"`
}

type SlackConfig struct {
	WebhookURL string `yaml:"webhook_url"`
	Channel    string `yaml:"channel,omitempty"`
	Username   string `yaml:"username,omitempty"`
}

type WebhookConfig struct {
	URL     string            `yaml:"url"`
	Headers map[string]string `yaml:"headers,omitempty"`
}

type Config struct {
	Default      Defaults     `yaml:"default"`
	Projects     []Project    `yaml:"projects"`
	GlobalAlerts GlobalAlerts `yaml:"global_alerts,omitempty"`
}

func Load(path string) (*Config, error) {
	b, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var cfg Config
	if err := yaml.Unmarshal(b, &cfg); err != nil {
		return nil, err
	}
	expandEnvInConfig(&cfg)
	return &cfg, nil
}

// expandEnvVars replaces ${VAR} with os.Getenv("VAR") in strings
func expandEnvVars(s string) string {
	if !strings.Contains(s, "${") {
		return s
	}
	start := strings.Index(s, "${")
	end := strings.Index(s[start:], "}")
	if end == -1 {
		return s
	}
	key := s[start+2 : start+end]
	return strings.Replace(s, "${"+key+"}", os.Getenv(key), 1)
}

// expandEnvInConfig replaces ${VAR} with os.Getenv("VAR") throughout the config
func expandEnvInConfig(cfg *Config) {
	// Expand in headers
	for i := range cfg.Projects {
		for j := range cfg.Projects[i].Monitors {
			monitor := &cfg.Projects[i].Monitors[j]

			// Expand headers
			for k, v := range monitor.Headers {
				monitor.Headers[k] = expandEnvVars(v)
			}

			// Expand auth fields
			if monitor.Auth != nil {
				monitor.Auth.Token = expandEnvVars(monitor.Auth.Token)
				monitor.Auth.Username = expandEnvVars(monitor.Auth.Username)
				monitor.Auth.Password = expandEnvVars(monitor.Auth.Password)
				monitor.Auth.ClientID = expandEnvVars(monitor.Auth.ClientID)
				monitor.Auth.ClientSecret = expandEnvVars(monitor.Auth.ClientSecret)
				monitor.Auth.TokenURL = expandEnvVars(monitor.Auth.TokenURL)
				monitor.Auth.APIKey = expandEnvVars(monitor.Auth.APIKey)
			}

			// Expand alert webhook URLs
			for k := range monitor.Alerts {
				monitor.Alerts[k].WebhookURL = expandEnvVars(monitor.Alerts[k].WebhookURL)
			}
		}
	}

	// Expand global alert configs
	cfg.GlobalAlerts.Email.Username = expandEnvVars(cfg.GlobalAlerts.Email.Username)
	cfg.GlobalAlerts.Email.Password = expandEnvVars(cfg.GlobalAlerts.Email.Password)
	cfg.GlobalAlerts.Email.From = expandEnvVars(cfg.GlobalAlerts.Email.From)
	cfg.GlobalAlerts.Email.SMTPHost = expandEnvVars(cfg.GlobalAlerts.Email.SMTPHost)
	cfg.GlobalAlerts.Email.UseTLS = (cfg.GlobalAlerts.Email.UseTLS)
	cfg.GlobalAlerts.Slack.WebhookURL = expandEnvVars(cfg.GlobalAlerts.Slack.WebhookURL)
	cfg.GlobalAlerts.Webhook.URL = expandEnvVars(cfg.GlobalAlerts.Webhook.URL)
}
