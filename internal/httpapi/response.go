package httpapi

import (
	"encoding/json"
	"net/http"
)

// StandardResponse represents the standard API response format
type StandardResponse struct {
	Data          interface{} `json:"data"`
	StatusCode    int         `json:"status_code"`    // 1 for success, 0 for failure
	StatusMessage string      `json:"status_message"` // Human-readable status message
}

// ResponseHelper provides methods for creating standardized responses
type ResponseHelper struct{}

// NewResponseHelper creates a new response helper
func NewResponseHelper() *ResponseHelper {
	return &ResponseHelper{}
}

// Success creates a successful response with status code 1
func (rh *ResponseHelper) Success(w http.ResponseWriter, data interface{}, message string) {
	if message == "" {
		message = "Success"
	}
	
	response := StandardResponse{
		Data:          data,
		StatusCode:    1,
		StatusMessage: message,
	}
	
	rh.writeJSONResponse(w, http.StatusOK, response)
}

// Error creates an error response with status code 0
func (rh *ResponseHelper) Error(w http.ResponseWriter, httpStatus int, message string, data interface{}) {
	if message == "" {
		message = "An error occurred"
	}
	
	response := StandardResponse{
		Data:          data,
		StatusCode:    0,
		StatusMessage: message,
	}
	
	rh.writeJSONResponse(w, httpStatus, response)
}

// BadRequest creates a bad request error response
func (rh *ResponseHelper) BadRequest(w http.ResponseWriter, message string) {
	rh.Error(w, http.StatusBadRequest, message, nil)
}

// NotFound creates a not found error response
func (rh *ResponseHelper) NotFound(w http.ResponseWriter, message string) {
	rh.Error(w, http.StatusNotFound, message, nil)
}

// InternalError creates an internal server error response
func (rh *ResponseHelper) InternalError(w http.ResponseWriter, message string) {
	rh.Error(w, http.StatusInternalServerError, message, nil)
}

// writeJSONResponse writes a JSON response with proper headers
func (rh *ResponseHelper) writeJSONResponse(w http.ResponseWriter, httpStatus int, response StandardResponse) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(httpStatus)
	
	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	_ = enc.Encode(response)
}

// HealthResponse represents a simple health check response
func (rh *ResponseHelper) Health(w http.ResponseWriter) {
	rh.Success(w, map[string]string{"status": "healthy"}, "Service is healthy")
}
