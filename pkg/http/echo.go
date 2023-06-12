package http

import (
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"

	"github.com/karta0898098/iam/pkg/errors"
)

type Config struct {
	Mode string `mapstructure:"mode"`
	Port string `mapstructure:"port"`
}

// NewEcho http handler
func NewEcho(config Config) *echo.Echo {
	echo.NotFoundHandler = EchoNotFoundHandler

	e := echo.New()
	e.Validator = NewEchoValidator()

	if config.Mode == "release" {
		e.Debug = false
		e.HideBanner = true
		e.HidePort = true

	} else {
		e.Debug = true
		e.HideBanner = false
		e.HidePort = false
	}

	e.HTTPErrorHandler = EchoErrorHandler
	return e
}

// EchoErrorHandler error handle for echo
func EchoErrorHandler(err error, c echo.Context) {
	if err == nil {
		return
	}

	echoError, ok := err.(*echo.HTTPError)
	if ok {
		_ = c.JSON(echoError.Code, echoError)
		return
	}

	e := errors.TryConvert(err)
	if e != nil {

		code, resp := e.ToViewModel()
		_ = c.JSON(code, resp)
	} else {
		code, resp := errors.New(err.Error()).ToViewModel()
		_ = c.JSON(code, resp)
	}
}

// EchoNotFoundHandler responds not found response.
func EchoNotFoundHandler(c echo.Context) error {
	return errors.ErrPageNotFound
}

// EchoValidator fot echo default validator
type EchoValidator struct {
	validator *validator.Validate
}

// NewEchoValidator new echo validator
func NewEchoValidator() *EchoValidator {
	return &EchoValidator{validator: validator.New()}
}

// Validate for echo validator interface
func (e *EchoValidator) Validate(i interface{}) error {
	return e.validator.Struct(i)
}
