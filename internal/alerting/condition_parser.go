package alerting

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"

	"pingpulse/internal/monitor"
)

// ConditionParser handles parsing and evaluation of alert conditions
type ConditionParser struct{}

// NewConditionParser creates a new condition parser
func NewConditionParser() *ConditionParser {
	return &ConditionParser{}
}

// EvaluateCondition evaluates a condition string against a monitor result
// Supports conditions like:
// - status == 200, status != 200, status > 400, status >= 500
// - latency > 2s, latency <= 1s, latency >= 500ms
// - ok == true, ok == false, ok != true
// - Complex conditions with AND/OR: status != 200 OR latency > 5s
func (cp *ConditionParser) EvaluateCondition(result monitor.Result, condition string) bool {
	if condition == "" {
		return false
	}

	// Handle OR conditions first (lower precedence)
	if strings.Contains(strings.ToUpper(condition), " OR ") {
		parts := cp.splitByOperator(condition, " OR ")
		for _, part := range parts {
			if cp.EvaluateCondition(result, strings.TrimSpace(part)) {
				return true
			}
		}
		return false
	}

	// Handle AND conditions (higher precedence)
	if strings.Contains(strings.ToUpper(condition), " AND ") {
		parts := cp.splitByOperator(condition, " AND ")
		for _, part := range parts {
			if !cp.EvaluateCondition(result, strings.TrimSpace(part)) {
				return false
			}
		}
		return true
	}

	// Handle single condition
	return cp.evaluateSingleCondition(result, condition)
}

// splitByOperator splits a string by an operator while respecting case insensitivity
func (cp *ConditionParser) splitByOperator(condition, operator string) []string {
	// Create a case-insensitive regex pattern
	pattern := "(?i)" + regexp.QuoteMeta(operator)
	re := regexp.MustCompile(pattern)
	
	// Find all matches to preserve original case in the split
	matches := re.FindAllStringIndex(condition, -1)
	if len(matches) == 0 {
		return []string{condition}
	}

	var parts []string
	start := 0
	for _, match := range matches {
		parts = append(parts, condition[start:match[0]])
		start = match[1]
	}
	parts = append(parts, condition[start:])
	
	return parts
}

// evaluateSingleCondition evaluates a single condition (no AND/OR)
func (cp *ConditionParser) evaluateSingleCondition(result monitor.Result, condition string) bool {
	condition = strings.TrimSpace(condition)
	
	// Status conditions
	if strings.HasPrefix(condition, "status") {
		return cp.evaluateStatusCondition(result.Status, condition)
	}
	
	// Latency conditions
	if strings.HasPrefix(condition, "latency") {
		return cp.evaluateLatencyCondition(result.Latency, condition)
	}
	
	// OK conditions
	if strings.HasPrefix(condition, "ok") {
		return cp.evaluateOkCondition(result.OK, condition)
	}

	// Backward compatibility with old format
	return cp.evaluateLegacyCondition(result, condition)
}

// evaluateStatusCondition evaluates status-based conditions
func (cp *ConditionParser) evaluateStatusCondition(status int, condition string) bool {
	// Extract operator and value
	operators := []string{">=", "<=", "!=", "==", ">", "<", "="}
	
	for _, op := range operators {
		if idx := strings.Index(condition, op); idx != -1 {
			valueStr := strings.TrimSpace(condition[idx+len(op):])
			expectedStatus, err := strconv.Atoi(valueStr)
			if err != nil {
				return false
			}
			
			switch op {
			case "==", "=":
				return status == expectedStatus
			case "!=":
				return status != expectedStatus
			case ">":
				return status > expectedStatus
			case "<":
				return status < expectedStatus
			case ">=":
				return status >= expectedStatus
			case "<=":
				return status <= expectedStatus
			}
		}
	}
	
	return false
}

// evaluateLatencyCondition evaluates latency-based conditions
func (cp *ConditionParser) evaluateLatencyCondition(latency time.Duration, condition string) bool {
	// Extract operator and value
	operators := []string{">=", "<=", "!=", "==", ">", "<", "="}
	
	for _, op := range operators {
		if idx := strings.Index(condition, op); idx != -1 {
			valueStr := strings.TrimSpace(condition[idx+len(op):])
			expectedLatency, err := time.ParseDuration(valueStr)
			if err != nil {
				return false
			}
			
			switch op {
			case "==", "=":
				return latency == expectedLatency
			case "!=":
				return latency != expectedLatency
			case ">":
				return latency > expectedLatency
			case "<":
				return latency < expectedLatency
			case ">=":
				return latency >= expectedLatency
			case "<=":
				return latency <= expectedLatency
			}
		}
	}
	
	return false
}

// evaluateOkCondition evaluates ok-based conditions
func (cp *ConditionParser) evaluateOkCondition(ok bool, condition string) bool {
	// Extract operator and value
	operators := []string{"!=", "==", "="}
	
	for _, op := range operators {
		if idx := strings.Index(condition, op); idx != -1 {
			valueStr := strings.TrimSpace(condition[idx+len(op):])
			
			var expectedOk bool
			switch strings.ToLower(valueStr) {
			case "true", "1":
				expectedOk = true
			case "false", "0":
				expectedOk = false
			default:
				return false
			}
			
			switch op {
			case "==", "=":
				return ok == expectedOk
			case "!=":
				return ok != expectedOk
			}
		}
	}
	
	return false
}

// evaluateLegacyCondition provides backward compatibility with old condition format
func (cp *ConditionParser) evaluateLegacyCondition(result monitor.Result, condition string) bool {
	switch condition {
	case "status != 200":
		return result.Status != 200
	case "latency > 2s":
		return result.Latency > 2*time.Second
	case "latency > 3s":
		return result.Latency > 3*time.Second
	case "latency > 5s":
		return result.Latency > 5*time.Second
	case "ok == false":
		return !result.OK
	default:
		return false
	}
}

// ValidateCondition validates if a condition string is syntactically correct
func (cp *ConditionParser) ValidateCondition(condition string) error {
	if condition == "" {
		return fmt.Errorf("condition cannot be empty")
	}

	// Split by OR first
	orParts := cp.splitByOperator(condition, " OR ")
	for _, orPart := range orParts {
		// Split by AND
		andParts := cp.splitByOperator(strings.TrimSpace(orPart), " AND ")
		for _, andPart := range andParts {
			if err := cp.validateSingleCondition(strings.TrimSpace(andPart)); err != nil {
				return err
			}
		}
	}

	return nil
}

// validateSingleCondition validates a single condition
func (cp *ConditionParser) validateSingleCondition(condition string) error {
	condition = strings.TrimSpace(condition)
	
	// Check if it starts with a valid field
	validFields := []string{"status", "latency", "ok"}
	validField := false
	for _, field := range validFields {
		if strings.HasPrefix(condition, field) {
			validField = true
			break
		}
	}
	
	if !validField {
		return fmt.Errorf("condition must start with one of: %s", strings.Join(validFields, ", "))
	}
	
	// Check for valid operators
	operators := []string{">=", "<=", "!=", "==", ">", "<", "="}
	hasOperator := false
	for _, op := range operators {
		if strings.Contains(condition, op) {
			hasOperator = true
			break
		}
	}
	
	if !hasOperator {
		return fmt.Errorf("condition must contain a valid operator: %s", strings.Join(operators, ", "))
	}
	
	return nil
}
