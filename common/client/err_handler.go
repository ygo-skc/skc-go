package client

import (
	"net/http"

	"github.com/ygo-skc/skc-go/common/v2/model"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// RPCErrorToAPIError translates a gRPC error into a *model.APIError, mapping the
// gRPC status code to the closest HTTP status code. For not-found the service's
// own status message is surfaced so the caller sees which resource was missing;
// every other code gets a generic message to avoid leaking internal detail.
// Logging is intentionally left to the caller so each RPC method can record its
// own context.
func RPCErrorToAPIError(err error) *model.APIError {
	st := status.Convert(err)
	httpStatus := httpStatusFromGRPCCode(st.Code())

	message := http.StatusText(httpStatus)
	if httpStatus == http.StatusNotFound {
		message = st.Message()
	}
	return &model.APIError{Message: message, StatusCode: httpStatus}
}

// httpStatusFromGRPCCode maps a gRPC status code to the closest HTTP status
// code, following the same conventions as grpc-gateway's HTTPStatusFromCode.
func httpStatusFromGRPCCode(code codes.Code) int {
	switch code {
	case codes.OK:
		return http.StatusOK
	case codes.Canceled:
		return 499 // client closed request (nginx convention; no stdlib constant)
	case codes.InvalidArgument, codes.FailedPrecondition, codes.OutOfRange:
		return http.StatusBadRequest
	case codes.DeadlineExceeded:
		return http.StatusGatewayTimeout
	case codes.NotFound:
		return http.StatusNotFound
	case codes.AlreadyExists, codes.Aborted:
		return http.StatusConflict
	case codes.PermissionDenied:
		return http.StatusForbidden
	case codes.Unauthenticated:
		return http.StatusUnauthorized
	case codes.ResourceExhausted:
		return http.StatusTooManyRequests
	case codes.Unimplemented:
		return http.StatusNotImplemented
	case codes.Unavailable:
		return http.StatusServiceUnavailable
	case codes.Unknown, codes.Internal, codes.DataLoss:
		return http.StatusInternalServerError
	default:
		return http.StatusInternalServerError
	}
}
