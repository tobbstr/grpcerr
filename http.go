package grpcerr

import (
	"fmt"
	"net/http"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// ResponseWriterOption is an option function used to modify its http.ResponseWriter argument.
// For example, to set additional header values or to modify existing ones.
type ResponseWriterOption func(w http.ResponseWriter)

type httpResponseEncoder struct {
	st   *status.Status
	w    http.ResponseWriter
	opts []ResponseWriterOption
}

// AsJSON encodes the gRPC error as JSON and writes it to the http.ResponseWriter.
// If an error occurs it is returned, otherwise it returns nil.
func (f *httpResponseEncoder) AsJSON() error {
	if f.st == nil {
		f.w.WriteHeader(http.StatusInternalServerError)
		f.w.Write(nil)
		return fmt.Errorf("invalid argument: status was nil")
	}
	json, err := jsonBytesFromGrpcStatus(f.st)
	if err != nil {
		f.w.WriteHeader(http.StatusInternalServerError)
		f.w.Write(nil)
		return fmt.Errorf("could not get JSON as bytes from gRPC status: %w", err)
	}

	// Sets sane defaults
	f.w.Header().Set("Content-Type", "application/json")

	// Sets the passed options, which must be set between the Content-Type assignment and f.w.WriteHeader().
	// Otherwhise it's not possible to change the Content-Type header using the below options.
	for _, opt := range f.opts {
		opt(f.w)
	}

	// Sets sane defaults
	f.w.WriteHeader(httpStatusCodeFrom(f.st))

	f.w.Write(json)

	return nil
}

// NewHttpResponseEncodeWriter returns a function which is used to write a gRPC error to a http.ResponseWriter
// using an encoding such as JSON.
func NewHttpResponseEncodeWriter(w http.ResponseWriter, opts ...ResponseWriterOption) func(*status.Status) *httpResponseEncoder {
	return func(st *status.Status) *httpResponseEncoder {
		return &httpResponseEncoder{
			st:   st,
			w:    w,
			opts: opts,
		}
	}
}

func httpStatusCodeFrom(st *status.Status) int {
	switch st.Code() {
	case codes.Aborted, codes.AlreadyExists:
		return http.StatusConflict
	case codes.DataLoss, codes.Unknown, codes.Internal:
		return http.StatusInternalServerError
	case codes.InvalidArgument, codes.FailedPrecondition, codes.OutOfRange:
		return http.StatusBadRequest
	case codes.Unauthenticated:
		return http.StatusUnauthorized
	case codes.PermissionDenied:
		return http.StatusForbidden
	case codes.NotFound:
		return http.StatusNotFound
	case codes.ResourceExhausted:
		return http.StatusTooManyRequests
	case codes.Canceled:
		return 499
	case codes.Unimplemented:
		return http.StatusNotImplemented
	case codes.Unavailable:
		return http.StatusServiceUnavailable
	case codes.DeadlineExceeded:
		return http.StatusGatewayTimeout
	}

	// This error code should never be returned
	return http.StatusInternalServerError
}
