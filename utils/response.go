// Standard JSON response helpers used by all API controllers.
// Every JSON response from /api/* routes uses these helpers for consistency.

package utils

// JSONResponse is the standard envelope for all successful API responses.
type JSONResponse struct {
	Status  string      `json:"status"`
	Data    interface{} `json:"data,omitempty"`
	Message string      `json:"message,omitempty"`
}

// JSONError is the standard envelope for all API error responses.
type JSONError struct {
	Status  string `json:"status"`
	Message string `json:"message"`
	Code    int    `json:"code"`
}

// OKResponse builds a standard success response envelope.
func OKResponse(data interface{}) JSONResponse {
	return JSONResponse{Status: "ok", Data: data}
}

// CreatedResponse builds a 201 success response envelope.
func CreatedResponse(data interface{}) JSONResponse {
	return JSONResponse{Status: "created", Data: data}
}

// ErrorResponse builds a standard error response envelope.
func ErrorResponse(message string, code int) JSONError {
	return JSONError{Status: "error", Message: message, Code: code}
}
