package errors

import (
	"fmt"
	"net/http"
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
	fmt.Println(opts)
	return &HttpError{Code: http.StatusBadRequest, Message: opts.Message, Err: opts.Err}
}

func NewNotFound(opts SimpleErrorFuncOptions) *HttpError {
	fmt.Println(opts)
	return &HttpError{Code: http.StatusNotFound, Message: opts.Message}
}

func NewUnauthorized(opts SimpleErrorFuncOptions) *HttpError {
	fmt.Println(opts)
	return &HttpError{Code: http.StatusUnauthorized, Message: opts.Message}
}

func NewInternal(opts ErrorFuncOptions) *HttpError {
	return &HttpError{Code: http.StatusInternalServerError, Message: opts.Message, Err: opts.Err}
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