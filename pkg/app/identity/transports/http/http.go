package http

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	httptransport "github.com/go-kit/kit/transport/http"
	"github.com/labstack/echo/v4"

	"github.com/karta0898098/iam/pkg/app/identity/endpoints"
	"github.com/karta0898098/iam/pkg/errors"
)

// MakeSignin signin endpoint
func MakeSignin(endpoints endpoints.Endpoints) http.Handler {
	return httptransport.NewServer(
		endpoints.SigninEndpoint,
		decodeHTTPSigninRequest,
		encodeHTTPSigninResponse,
		httptransport.ServerErrorEncoder(errors.ErrorResponse),
		httptransport.ServerErrorHandler(errors.NewLoggingErrorHandle()),
	)
}

// decodeHTTPSigninRequest is a transport/http.DecodeRequestFunc that decodes a
// JSON-encoded request from the HTTP request body. Primarily useful in a server.
func decodeHTTPSigninRequest(ctx context.Context, r *http.Request) (interface{}, error) {

	fmt.Printf("%#v\n", r.URL.Path)
	var req endpoints.SigninRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		return nil, errors.Wrap(errors.ErrInvalidInput, "input request is not json")
	}
	return req, err
}

// encodeHTTPSigninResponse is a transport/http.EncodeResponseFunc that decodes a
// JSON-encode response from the HTTP response body. Primarily useful in a server.
func encodeHTTPSigninResponse(ctx context.Context, w http.ResponseWriter, response interface{}) error {
	w.Header().Set(echo.HeaderContentType, echo.MIMEApplicationJSONCharsetUTF8)
	if headers, ok := response.(httptransport.Headerer); ok {
		for k, values := range headers.Headers() {
			for _, v := range values {
				w.Header().Add(k, v)
			}
		}
	}
	code := http.StatusOK
	if sc, ok := response.(httptransport.StatusCoder); ok {
		code = sc.StatusCode()
	}
	w.WriteHeader(code)
	if code == http.StatusNoContent {
		return nil
	}

	return json.NewEncoder(w).Encode(response)
}

// MakeSignup make signup endpoint
func MakeSignup(endpoints endpoints.Endpoints) http.Handler {
	return httptransport.NewServer(
		endpoints.SignupEndpoint,
		decodeHTTPSignupRequest,
		encodeHTTPSignupResponse,
		httptransport.ServerErrorEncoder(errors.ErrorResponse),
		httptransport.ServerErrorHandler(errors.NewLoggingErrorHandle()),
	)
}

// decodeHTTPSignupRequest is a transport/http.DecodeRequestFunc that decodes a
// JSON-encoded request from the HTTP request body. Primarily useful in a server.
func decodeHTTPSignupRequest(ctx context.Context, r *http.Request) (interface{}, error) {
	var req endpoints.SignupRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		return nil, errors.Wrap(errors.ErrInvalidInput, "input request is not json")
	}
	return req, err
}

// encodeHTTPSigninResponse is a transport/http.EncodeResponseFunc that decodes a
// JSON-encode response from the HTTP response body. Primarily useful in a server.
func encodeHTTPSignupResponse(ctx context.Context, w http.ResponseWriter, response interface{}) error {
	w.Header().Set(echo.HeaderContentType, echo.MIMEApplicationJSONCharsetUTF8)
	if headers, ok := response.(httptransport.Headerer); ok {
		for k, values := range headers.Headers() {
			for _, v := range values {
				w.Header().Add(k, v)
			}
		}
	}
	code := http.StatusOK
	if sc, ok := response.(httptransport.StatusCoder); ok {
		code = sc.StatusCode()
	}
	w.WriteHeader(code)
	if code == http.StatusNoContent {
		return nil
	}

	return json.NewEncoder(w).Encode(response)
}
