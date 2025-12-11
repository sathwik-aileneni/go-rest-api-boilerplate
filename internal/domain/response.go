package domain

// StandardResponse is the standard API response format
// All API responses should follow this structure for consistency
type StandardResponse struct {
	APIID  string        `json:"api_id"`            // Unique identifier for request tracing
	Errors []ErrorDetail `json:"errors,omitempty"`  // List of errors (empty on success)
	Data   interface{}   `json:"data,omitempty"`    // Response data (structured based on endpoint)
}

// ErrorDetail represents a single error in the response
type ErrorDetail struct {
	Code    string `json:"code"`              // Error code (e.g., "VALIDATION_ERROR", "NOT_FOUND")
	Message string `json:"message"`           // Human-readable error message
	Field   string `json:"field,omitempty"`   // Field name (for validation errors)
}

// Deprecated: Use StandardResponse instead
type Response struct {
	Success bool        `json:"success"`
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"data,omitempty"`
	Error   string      `json:"error,omitempty"`
}

type HealthResponse struct {
	Status    string `json:"status"`
	Timestamp string `json:"timestamp"`
}
