package errors

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc/codes"
)

// DetailData is metadata for client debug
type DetailData map[string]interface{}

type Detail struct {
	Type     string                 `json:"@type,omitempty"`
	Reason   string                 `json:"reason,omitempty"`
	Domain   string                 `json:"domain,omitempty"`
	Metadata map[string]interface{} `json:"metadata,omitempty"`
}

// Exception define custom error for tracmo cloud
type Exception struct {
	Code     int        `json:"code"`              // error code for client define how to handle error
	Status   int        `json:"status"`            // status is http status
	Message  string     `json:"message"`           // error message for client
	GRPCCode codes.Code `json:"grpc_code"`         //  grpc error code
	Details  []Detail   `json:"details,omitempty"` // details is metadata for client debug
}

// New server internal error with message
func New(message string) *Exception {
	return &Exception{
		Code:    ErrInternal.Code,
		Status:  ErrInternal.Status,
		Message: message,
	}
}

// Is Check input is same
func Is(err error, target error) bool {
	causeTargetErr, ok := errors.Cause(target).(*Exception)
	if !ok {
		return errors.Is(err, target)
	}

	causeErr, ok := errors.Cause(err).(*Exception)
	if !ok {
		return errors.Is(err, target)
	}

	return causeErr.Code == causeTargetErr.Code
}

func As(err error, target interface{}) bool {
	return errors.As(err, target)
}

func TryConvert(target error) *Exception {
	err, ok := errors.Cause(target).(*Exception)
	if !ok {
		return nil
	}
	return err
}

func Cause(target error) error {
	return errors.Cause(target)
}

// Wrap returns an error annotating err with a stack trace
// at the point Wrap is called, and the supplied message.
// If err is nil, Wrap returns nil.
func Wrap(err error, msg string) error {
	return errors.Wrap(err, msg)
}

// Wrapf returns an error annotating err with a stack trace
// at the point Wrapf is called, and the format specifier.
// If err is nil, Wrapf returns nil.
func Wrapf(err error, format string, args ...interface{}) error {
	return errors.Wrapf(err, format, args...)
}

// WithMessage annotates err with a new message.
// If err is nil, WithMessage returns nil.
func WithMessage(err error, message string) error {
	return errors.WithMessage(err, message)
}

// WithMessagef annotates err with the format specifier.
// If err is nil, WithMessagef returns nil.
func WithMessagef(err error, format string, args ...interface{}) error {
	if err == nil {
		return nil
	}
	return errors.WithMessagef(err, format, args...)
}

// Error implement golang error
func (e *Exception) Error() string {
	var (
		b strings.Builder
	)
	_, _ = b.WriteRune('[')
	_, _ = b.WriteString(strconv.Itoa(e.Code))
	_, _ = b.WriteRune(']')
	_, _ = b.WriteRune(' ')
	_, _ = b.WriteString(e.Message)
	return b.String()
}

// WithDetails set detail error message
func (e *Exception) WithDetails(details ...Detail) *Exception {
	newErr := *e
	newErr.Details = append(newErr.Details, details...)
	return &newErr
}

type View struct {
	Code    int      `json:"code"`
	Info    string   `json:"info"`
	Details []Detail `json:"details,omitempty"`
}

// ToViewModel to restful view
func (e *Exception) ToViewModel() (int, *View) {
	if len(e.Details) == 0 {
		e.Details = make([]Detail, 0)
	}

	return e.Status, &View{
		Code:    e.Code,
		Info:    e.Message,
		Details: e.Details,
	}
}

// ErrorResponse error response for go-kit
func ErrorResponse(ctx context.Context, err error, w http.ResponseWriter) {
	var (
		code     int
		response *View
	)

	switch val := Cause(err).(type) {
	case *Exception:
		code, response = val.ToViewModel()
	default:
		code, response = New(err.Error()).ToViewModel()
	}

	w.Header().Set(echo.HeaderContentType, echo.MIMEApplicationJSONCharsetUTF8)
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(response)
}

type LoggerErrorHandle struct {
}

func NewLoggingErrorHandle() *LoggerErrorHandle {
	return &LoggerErrorHandle{}
}

func (h *LoggerErrorHandle) Handle(ctx context.Context, err error) {
	logger := log.Ctx(ctx)

	switch val := errors.Cause(err).(type) {
	case *Exception:
		var (
			event *zerolog.Event
		)

		if val.Code >= http.StatusBadRequest &&
			val.Code < http.StatusInternalServerError {
			event = logger.Warn().Stack()
		} else {
			event = logger.Error().Stack()
		}

		event.Err(err).Msgf("server occur error")
	default:
		logger.
			Error().
			Stack().
			Err(err).
			Msg("failed to handle")
	}
}
