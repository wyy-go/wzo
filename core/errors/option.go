package errors

import "fmt"

type Option func(*Error)

// WithMessage modifies the message
func WithMessage(s string) Option {
	return func(e *Error) {
		e.Message = s
	}
}

// WithDetail modifies the detail
func WithDetail(d string) Option {
	return func(e *Error) {
		e.Detail = d
	}
}

// WithMessagef modifies the message
func WithMessagef(format string, args ...any) Option {
	return func(e *Error) {
		e.Message = fmt.Sprintf(format, args...)
	}
}

// WithMetadata add metadata to the error
func WithMetadata(k, v string) Option {
	return func(e *Error) {
		if k != "" && v != "" {
			if e.Metadata == nil {
				e.Metadata = make(map[string]string)
			}
			e.Metadata[k] = v
		}
	}
}
