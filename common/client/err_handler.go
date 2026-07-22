package client

import (
	"net/http"

	"github.com/ygo-skc/skc-go/common/v2/model"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func rpcErrorToAPIError(err error) *model.APIError {
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
