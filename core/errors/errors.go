package errors

import (
	"errors"
	"fmt"
	"github.com/bytedance/sonic"
	"io"
	"net/http"

	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/status"
)

const (
	defaultStatusCode = 500
	defaultErrorCode  = 500
)

var _ error = (*Error)(nil)

// Error define the error type
type Error struct {
	StatusCode int32             `json:"status_code,omitempty"`
	Code       int32             `json:"code,omitempty"`
	Message    string            `json:"message,omitempty"`
	Detail     string            `json:"detail,omitempty"`
	Metadata   map[string]string `json:"metadata,omitempty"`
}

// Error implement `Error() string` interface.
func (e *Error) Error() string {
	s, _ := sonic.MarshalString(e)
	return s
}

func (e *Error) clone() *Error {
	if e == nil {
		return nil
	}

	metadata := make(map[string]string, len(e.Metadata))
	for k, v := range e.Metadata {
		metadata[k] = v
	}
	return &Error{
		StatusCode: e.StatusCode,
		Code:       e.Code,
		Message:    e.Message,
		Detail:     e.Detail,
		Metadata:   metadata,
	}
}

func (e *Error) Format(s fmt.State, verb rune) {
	copied := e.clone()
	delete(copied.Metadata, "_wzo_error_stack")
	msg := fmt.Sprintf("error: code = %d message = %s detail = %s metadata = %v", e.Code, e.Message, e.Detail, copied.Metadata)

	switch verb {
	case 's', 'v':
		switch {
		case s.Flag('+'):
			_, _ = io.WriteString(s, msg+"\n"+e.Stack())

		default:
			_, _ = io.WriteString(s, msg)
		}
	}
}

// TakeOption custom options
func (e *Error) TakeOption(opts ...Option) *Error {
	if e == nil {
		return nil
	}
	for _, f := range opts {
		f(e)
	}
	return e
}

func Code(err error) int {
	if err == nil {
		return 0
	}
	return int(FromError(err).Code)
}

func Parse(err string) *Error {
	e := new(Error)
	errr := sonic.Unmarshal([]byte(err), e)
	if errr != nil {
		e.StatusCode = defaultStatusCode
		e.Code = defaultErrorCode
		e.Message = "Internal Server Error"
		e.Detail = err
	}
	if e.Code == 0 {
		e.StatusCode = defaultStatusCode
		e.Code = defaultErrorCode
	}
	return e
}

func FromError(err error) *Error {
	if err == nil {
		return nil
	}

	if e := new(Error); errors.As(err, &e) {
		return e
	}

	return Parse(err.Error())
}

func WrapRpcError(err error) *Error {
	if err == nil {
		return nil
	}

	e := FromError(err)
	if e.Metadata["_wzo_error_stack"] != "" {
		e.Metadata["_wzo_error_stack"] = stacktrace() + "\n" + e.Metadata["_wzo_error_stack"]
	}

	return e
}

func New(code int32, message string, opts ...Option) *Error {
	e := &Error{
		StatusCode: defaultStatusCode,
		Code:       code,
		Message:    message,
		Metadata:   map[string]string{},
	}

	e.Metadata["_wzo_error_stack"] = stacktrace()
	return e.TakeOption(opts...)
}

func Newf(code int32, message, format string, a ...any) *Error {
	e := &Error{
		StatusCode: defaultStatusCode,
		Code:       code,
		Message:    message,
		Detail:     fmt.Sprintf(format, a...),
		Metadata:   map[string]string{},
	}

	e.Metadata["_wzo_error_stack"] = stacktrace()
	return e
}

func NewWithStatusCode(statusCode, code int32, message, detail string) *Error {
	e := &Error{
		StatusCode: statusCode,
		Code:       code,
		Message:    message,
		Detail:     detail,
		Metadata:   map[string]string{},
	}

	e.Metadata["_wzo_error_stack"] = stacktrace()
	return e
}

func NewfWithStatusCode(statusCode, code int32, message, format string, a ...any) *Error {
	e := &Error{
		StatusCode: statusCode,
		Code:       code,
		Message:    message,
		Detail:     fmt.Sprintf(format, a...),
		Metadata:   map[string]string{},
	}

	e.Metadata["_wzo_error_stack"] = stacktrace()
	return e
}

func Errorf(code int32, message, format string, a ...any) error {
	e := &Error{
		StatusCode: defaultStatusCode,
		Code:       code,
		Message:    message,
		Detail:     fmt.Sprintf(format, a...),
		Metadata:   map[string]string{},
	}

	e.Metadata["_wzo_error_stack"] = stacktrace()
	return e
}

func ErrorfWithStatusCode(statusCode, code int32, message, format string, a ...any) error {
	e := &Error{
		StatusCode: statusCode,
		Code:       code,
		Message:    message,
		Detail:     fmt.Sprintf(format, a...),
		Metadata:   map[string]string{},
	}

	e.Metadata["_wzo_error_stack"] = stacktrace()
	return e
}

// ParseError parser error is `Error`, if not `Error`, new `Error` with code 500 and warp the err.
// err == nil: return nil
// err is not Error: return NewInternalServer
// err is Error(as te):
//
//	te == nil:  return nil
//	te != nil:  return te
func ParseError(err error) *Error {
	if err == nil {
		return nil
	}
	if te := new(Error); errors.As(err, &te) {
		return te
	}
	return NewInternalServer(WithDetail(err.Error()))
}

// EqualCode return true if error underlying code equal target code.
// err == nil: code = 200
// err is not Error: code = 500
// err is Error(as te):
//
//	te == nil:  code = 200
//	te != nil:  te.code
func EqualCode(err error, targetCode int32) bool {
	if err == nil {
		return http.StatusOK == targetCode
	}
	if te := new(Error); errors.As(err, &te) {
		if te == nil {
			return http.StatusOK == targetCode
		} else {
			return te != nil && te.Code == targetCode
		}
	}
	return http.StatusInternalServerError == targetCode
}

// GRPCStatus returns the Status represented by se.
func (e *Error) GRPCStatus() *status.Status {
	s, _ := status.New(ToGRPCCode(int(e.Code)), e.Message).
		WithDetails(&errdetails.ErrorInfo{
			Reason:   e.Error(),
			Metadata: e.Metadata,
		})
	return s
}
