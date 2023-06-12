package errors

import (
	"net/http"

	"google.golang.org/grpc/codes"
)

var (
	ErrInvalidInput       = &Exception{Code: 400001, Message: "One of the request inputs is not valid.", Status: http.StatusBadRequest, GRPCCode: codes.InvalidArgument}
	ErrInvalidHeaderValue = &Exception{Code: 400003, Message: "The value provided for one of the HTTP headers was not in the correct format.", Status: http.StatusBadRequest, GRPCCode: codes.InvalidArgument}
	ErrUnauthorized       = &Exception{Code: 401001, Message: "The request unauthorized", Status: http.StatusUnauthorized, GRPCCode: codes.PermissionDenied}
	ErrForbidden          = &Exception{Code: 403001, Message: "Forbidden.", Status: http.StatusForbidden, GRPCCode: codes.PermissionDenied}
	ErrPageNotFound       = &Exception{Code: 404001, Message: "Page not found.", Status: http.StatusNotFound, GRPCCode: codes.NotFound}
	ErrResourceNotFound   = &Exception{Code: 404002, Message: "The specified resource does not exist.", Status: http.StatusNotFound}
	ErrConflict           = &Exception{Code: 409001, Message: "The request conflict.", Status: http.StatusConflict, GRPCCode: codes.AlreadyExists}
	ErrTooManyRequests    = &Exception{Code: 429001, Message: "Too Many Requests", Status: http.StatusTooManyRequests, GRPCCode: codes.PermissionDenied}
	ErrInternal           = &Exception{Code: 500001, Message: "Serve occur error.", Status: http.StatusInternalServerError, GRPCCode: codes.Internal}
)
