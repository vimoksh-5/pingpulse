package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"pingpulse/internal/alerting"
	"pingpulse/internal/config"
	"pingpulse/internal/httpapi"
	"pingpulse/internal/monitor"
	"pingpulse/internal/store"

	"github.com/joho/godotenv"
)

func main() {
	// Load .env file
	if err := godotenv.Load(); err != nil {
		log.Printf("Warning: .env file not found: %v", err)
	}

	cfg, err := config.Load("config.yaml")
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	mem := store.NewMemoryStore(50)
	mgr := monitor.NewManager(mem)

	// Initialize alert manager
	alertMgr := alerting.NewAlertManager(cfg)
	mgr.SetAlertManager(alertMgr)

	// Loop over projects and their monitors
	for _, proj := range cfg.Projects {
		for _, mon := range proj.Monitors {
			// Merge monitor settings with defaults
			interval := cfg.Default.Interval.Duration
			if mon.Interval != nil {
				interval = mon.Interval.Duration
			}
			timeout := cfg.Default.Timeout.Duration
			if mon.Timeout != nil {
				timeout = mon.Timeout.Duration
			}
			maxRetries := cfg.Default.MaxRetries
			if mon.MaxRetries != nil {
				maxRetries = *mon.MaxRetries
			}

			// Convert config types to monitor types
			var auth *monitor.Auth
			if mon.Auth != nil {
				auth = &monitor.Auth{
					Type:         mon.Auth.Type,
					Token:        mon.Auth.Token,
					Username:     mon.Auth.Username,
					Password:     mon.Auth.Password,
					ClientID:     mon.Auth.ClientID,
					ClientSecret: mon.Auth.ClientSecret,
					TokenURL:     mon.Auth.TokenURL,
					APIKey:       mon.Auth.APIKey,
					APIKeyHeader: mon.Auth.APIKeyHeader,
				}
			}

			var bodyValidation *monitor.BodyValidation
			if mon.BodyValidation != nil {
				bodyValidation = &monitor.BodyValidation{
					Path:     mon.BodyValidation.Path,
					Expected: mon.BodyValidation.Expected,
				}
			}

			var alerts []monitor.Alert
			for _, alert := range mon.Alerts {
				alerts = append(alerts, monitor.Alert{
					Name:       alert.Name,
					Condition:  alert.Condition,
					Channels:   alert.Channels,
					Recipients: alert.Recipients,
					WebhookURL: alert.WebhookURL,
					Cooldown:   alert.Cooldown.Duration,
				})
			}

			target := monitor.Target{
				Project:        proj.Name,
				Name:           mon.Name,
				URL:            mon.URL,
				Method:         mon.Method,
				Headers:        mon.Headers,
				Body:           mon.Body,
				Every:          interval,
				Timeout:        timeout,
				MaxRetries:     maxRetries,
				Auth:           auth,
				ExpectedStatus: mon.ExpectedStatus,
				ExpectedLatency: func() *time.Duration {
					if mon.ExpectedLatency != nil {
						d := time.Duration(mon.ExpectedLatency.Duration)
						return &d
					}
					return nil
				}(),
				SSLVerify:       mon.SSLVerify,
				FollowRedirects: mon.FollowRedirects,
				BodyValidation:  bodyValidation,
				Alerts:          alerts,
			}
			mgr.Start(ctx, target)
		}
	}

	srv := httpapi.NewServer(mem, ":8080")
	go func() {
		if err := srv.Start(); err != nil {
			log.Printf("http server stopped: %v", err)
			stop()
		}
	}()

	<-ctx.Done()
	_ = srv.Shutdown(ctx)
}
