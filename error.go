package web

import "fmt"

// Error describes an API error object
type Error struct {
	Code    int    `json:"code,omitempty"`
	Message string `json:"message,omitempty"`
}

// ValidationError conveience method to make a error object for validation errors
func ValidationError(format string, v ...interface{}) *Error {
	return &Error{
		Code:    400,
		Message: fmt.Sprintf(format, v...),
	}
}
