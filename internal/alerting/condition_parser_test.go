package alerting

import (
	"testing"
	"time"

	"pingpulse/internal/monitor"
)

func TestConditionParser_EvaluateCondition(t *testing.T) {
	parser := NewConditionParser()
	
	result := monitor.Result{
		Status:  500,
		Latency: 3 * time.Second,
		OK:      false,
	}

	tests := []struct {
		condition string
		expected  bool
		name      string
	}{
		// Simple status conditions
		{"status == 500", true, "status equals"},
		{"status != 200", true, "status not equals"},
		{"status > 400", true, "status greater than"},
		{"status >= 500", true, "status greater than or equal"},
		{"status < 600", true, "status less than"},
		{"status <= 500", true, "status less than or equal"},
		
		// Simple latency conditions
		{"latency > 2s", true, "latency greater than"},
		{"latency >= 3s", true, "latency greater than or equal"},
		{"latency < 5s", true, "latency less than"},
		{"latency <= 3s", true, "latency less than or equal"},
		{"latency == 3s", true, "latency equals"},
		{"latency != 2s", true, "latency not equals"},
		
		// Simple ok conditions
		{"ok == false", true, "ok equals false"},
		{"ok != true", true, "ok not equals true"},
		
		// Complex AND conditions
		{"status >= 500 AND latency > 2s", true, "AND condition true"},
		{"status == 500 AND ok == false", true, "AND condition with ok"},
		{"status == 200 AND latency > 2s", false, "AND condition false"},
		
		// Complex OR conditions
		{"status == 200 OR status == 500", true, "OR condition true"},
		{"status == 200 OR latency > 10s", false, "OR condition false"},
		{"ok == true OR status >= 500", true, "OR condition with ok"},
		
		// Backward compatibility
		{"status != 200", true, "legacy status condition"},
		{"latency > 2s", true, "legacy latency condition"},
		{"ok == false", true, "legacy ok condition"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actual := parser.EvaluateCondition(result, tt.condition)
			if actual != tt.expected {
				t.Errorf("EvaluateCondition(%q) = %v, want %v", tt.condition, actual, tt.expected)
			}
		})
	}
}

func TestConditionParser_ValidateCondition(t *testing.T) {
	parser := NewConditionParser()

	tests := []struct {
		condition string
		valid     bool
		name      string
	}{
		{"status == 200", true, "valid status condition"},
		{"latency > 2s", true, "valid latency condition"},
		{"ok == true", true, "valid ok condition"},
		{"status >= 400 AND latency > 1s", true, "valid AND condition"},
		{"status == 200 OR status == 500", true, "valid OR condition"},
		{"", false, "empty condition"},
		{"invalid_field == 200", false, "invalid field"},
		{"status", false, "missing operator"},
		{"status invalid_op 200", false, "invalid operator"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := parser.ValidateCondition(tt.condition)
			if tt.valid && err != nil {
				t.Errorf("ValidateCondition(%q) returned error: %v", tt.condition, err)
			}
			if !tt.valid && err == nil {
				t.Errorf("ValidateCondition(%q) should have returned error", tt.condition)
			}
		})
	}
}
