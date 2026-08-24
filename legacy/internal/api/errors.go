package api

import (
	"errors"
	"fmt"
	"net/http"
)

// Error codes and exit-code mapping for the CLI.
const (
	ExitOK                = 0
	ExitUsage             = 2
	ExitAuth              = 10
	ExitInsufficientScope = 11
	ExitNotFound          = 12
	ExitRateLimit         = 13
	ExitValidation        = 14
	ExitAPI               = 15
	ExitNetwork           = 16
)

// APIError is a structured Public API error.
type APIError struct {
	Status  int
	Code    string
	Message string
	Body    []byte
}

func (e *APIError) Error() string {
	if e.Code != "" {
		return fmt.Sprintf("%s: %s (HTTP %d)", e.Code, e.Message, e.Status)
	}
	return fmt.Sprintf("%s (HTTP %d)", e.Message, e.Status)
}

// UsageError is a CLI usage / prompt / validation problem (exit 2).
type UsageError struct {
	Msg string
}

func (e *UsageError) Error() string { return e.Msg }

// AbortError is returned when the user declines a confirmation (exit 0).
type AbortError struct{}

func (e *AbortError) Error() string { return "aborted" }

// ExitCode maps the error to a process exit code.
func ExitCode(err error) int {
	if err == nil {
		return ExitOK
	}
	var abort *AbortError
	if errors.As(err, &abort) {
		return ExitOK
	}
	var usage *UsageError
	if errors.As(err, &usage) {
		return ExitUsage
	}
	var ae *APIError
	if errors.As(err, &ae) {
		switch {
		case ae.Status == http.StatusUnauthorized || ae.Code == "MISSING_API_KEY":
			return ExitAuth
		case ae.Status == http.StatusForbidden:
			return ExitInsufficientScope
		case ae.Status == http.StatusNotFound:
			return ExitNotFound
		case ae.Status == http.StatusTooManyRequests:
			return ExitRateLimit
		case ae.Status >= 400 && ae.Status < 500:
			return ExitValidation
		default:
			return ExitAPI
		}
	}
	return ExitNetwork
}

func mapStatusError(status int, code, message string, body []byte) *APIError {
	if message == "" {
		message = http.StatusText(status)
	}
	if code == "" {
		switch status {
		case http.StatusUnauthorized:
			code = "AUTH_ERROR"
		case http.StatusForbidden:
			code = "INSUFFICIENT_SCOPE"
		case http.StatusNotFound:
			code = "NOT_FOUND"
		case http.StatusTooManyRequests:
			code = "RATE_LIMIT"
		default:
			if status >= 400 && status < 500 {
				code = "VALIDATION_ERROR"
			} else {
				code = "API_ERROR"
			}
		}
	}
	return &APIError{Status: status, Code: code, Message: message, Body: body}
}
