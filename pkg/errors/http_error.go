package errors

import (
	"fmt"
	"net/http"
	"github.com/aragornz325/piloto-api/pkg/logger"
	"go.uber.org/zap"
)

// Error implements the error interface for HttpError.
// It returns the error message, optionally including the wrapped error if present.
func (e *HttpError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("%s: %v", e.Message, e.Err)
	}
	return e.Message
}

// Unwrap returns the underlying error wrapped by the HttpError, enabling error unwrapping with errors.Unwrap and error inspection.
func (e *HttpError) Unwrap() error {
	return e.Err
}

func NewBadRequest(opts ErrorFuncOptions ) *HttpError {
	logger.Log.Error("Bad request", zap.String("message", opts.Message), zap.Error(opts.Err))
	return &HttpError{Code: http.StatusBadRequest, Message: opts.Message, Err: opts.Err}
}

func NewNotFound(opts SimpleErrorFuncOptions) *HttpError {
	logger.Log.Error("Not found", zap.String("message", opts.Message))
	return &HttpError{Code: http.StatusNotFound, Message: opts.Message}
}

func NewUnauthorized(opts SimpleErrorFuncOptions) *HttpError {
	logger.Log.Error("Unauthorized", zap.String("message", opts.Message))
	return &HttpError{Code: http.StatusUnauthorized, Message: opts.Message}
}

func NewInternal(opts ErrorFuncOptions) *HttpError {
	logger.Log.Error("Internal server error", zap.String("message", opts.Message), zap.Error(opts.Err))
	return &HttpError{Code: http.StatusInternalServerError, Message: opts.Message, Err: opts.Err}
}

func NewForbidden(opts ErrorFuncOptions) *HttpError {
	logger.Log.Error("Forbidden", zap.String("message", opts.Message), zap.Error(opts.Err))
	return &HttpError{Code: http.StatusForbidden, Message: opts.Message, Err: opts.Err}
}


//-------structs--------///

type ErrorFuncOptions struct {
	Message string
	Err     error
}

type SimpleErrorFuncOptions struct {
	Message string
}

type HttpError struct {
	Code    int    // Código HTTP
	Message string // Mensaje amigable
	Err     error  // Error original (opcional)
}