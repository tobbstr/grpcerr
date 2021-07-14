package grpcerr

import (
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/encoding/protojson"
)

const (
	defaultInvalidArgumentErrMsg    = "Client specified an invalid argument. Check error details for more information."
	defaultOutOfRangeErrMsg         = "Client specified an invalid range."
	defaultFailedPreconditionErrMsg = "Request can not be executed in the current system state, such as deleting a non-empty directory."
	defaultUnauthenticatedErrMsg    = "Request not authenticated due to missing, invalid, or expired security credentials."
	defaultPermissionDeniedErrMsg   = "Client does not have sufficient permission. This can happen because the client doesn't have permission, or the API has not been enabled."
	defaultAbortedErrMsg            = "Concurrency conflict, such as read-modify-write conflict."
	defaultNotFoundErrMsg           = "A specified resource is not found."
	defaultAlreadyExistsErrMsg      = "Resource a client tried to create already exists."
	defaultResourceExhaustedErrMsg  = "Either out of resource quota or reaching rate limiting. The client should look for google.rpc.QuotaFailure error detail for more information."
	defaultCanceledErrMsg           = "Request cancelled by the client."
	defaultDataLossErrMsg           = "Unrecoverable data loss or data corruption. The client should report the error to the user."
	defaultUnknownErrMsg            = "Unknown server error. Typically a server bug."
	defaultInternalErrMsg           = "Internal server error. Typically a server bug."
	defaultUnimplementedErrMsg      = "API method not implemented by the server."
	defaultNotAvailableErrMsg       = "Network error occurred before reaching the server. Typically a network outage or misconfiguration."
	defaultUnavailableErrMsg        = "Service unavailable. Typically the server is down."
	defaultDeadlineExceededErrMsg   = "Request deadline exceeded. This will happen only if the caller sets a deadline that is shorter than the method's default deadline (i.e. requested deadline is not enough for the server to process the request) and the request did not finish within the deadline."
)

func AddDebugInfo(gRPCErr *status.Status, debugInfo *errdetails.DebugInfo) (*status.Status, error) {
	return gRPCErr.WithDetails(debugInfo)
}

func AddRequestInfo(gRPCErr *status.Status, requestID, servingData string) (*status.Status, error) {
	requestInfoDetails := errdetails.RequestInfo{
		RequestId:   requestID,
		ServingData: servingData,
	}

	return gRPCErr.WithDetails(&requestInfoDetails)
}

func AddHelp(gRPCErr *status.Status, links []*errdetails.Help_Link) (*status.Status, error) {
	if len(links) == 0 {
		return gRPCErr, nil
	}

	helpDetails := errdetails.Help{Links: links}

	return gRPCErr.WithDetails(&helpDetails)
}

func AddLocalizedMessage(gRPCErr *status.Status, locale, msg string) (*status.Status, error) {
	localizedMessageDetails := errdetails.LocalizedMessage{
		Locale:  locale,
		Message: msg,
	}

	return gRPCErr.WithDetails(&localizedMessageDetails)
}

func NewInvalidArgument(errMsg string, fieldViolations []*errdetails.BadRequest_FieldViolation) (*status.Status, error) {
	return newGRPCErrorWithBadRequestDetails(codes.InvalidArgument, errMsg, fieldViolations)
}

func newGRPCErrorWithBadRequestDetails(code codes.Code, errMsg string, fieldViolations []*errdetails.BadRequest_FieldViolation) (*status.Status, error) {
	var st *status.Status
	if errMsg == "" {
		st = status.New(code, defaultInvalidArgumentErrMsg)
	} else {
		st = status.New(code, errMsg)
	}

	if len(fieldViolations) == 0 {
		return st, nil
	}

	badRequestDetails := errdetails.BadRequest{FieldViolations: fieldViolations}

	return st.WithDetails(&badRequestDetails)
}

func jsonBytesFromGrpcStatus(status *status.Status) ([]byte, error) {
	data, err := protojson.Marshal(status.Proto())
	if err != nil {
		return nil, err
	}

	return data, nil
}

func NewOutOfRange(errMsg string, fieldViolations []*errdetails.BadRequest_FieldViolation) (*status.Status, error) {
	return newGRPCErrorWithBadRequestDetails(codes.OutOfRange, errMsg, fieldViolations)
}

func NewFailedPrecondition(errMsg string, violations []*errdetails.PreconditionFailure_Violation) (*status.Status, error) {
	return newGRPCErrorWithFailedPreconditionDetails(codes.FailedPrecondition, errMsg, violations)
}

func newGRPCErrorWithFailedPreconditionDetails(code codes.Code, errMsg string, violations []*errdetails.PreconditionFailure_Violation) (*status.Status, error) {
	var st *status.Status
	if errMsg == "" {
		st = status.New(code, defaultFailedPreconditionErrMsg)
	} else {
		st = status.New(code, errMsg)
	}

	if len(violations) == 0 {
		return st, nil
	}

	failedPreconditionsViolationDetails := errdetails.PreconditionFailure{Violations: violations}

	return st.WithDetails(&failedPreconditionsViolationDetails)
}

func NewUnauthenticated(errMsg string, errorInfo *errdetails.ErrorInfo) (*status.Status, error) {
	return newGRPCErrorWithErrorInfo(codes.Unauthenticated, errMsg, errorInfo)
}

func newGRPCErrorWithErrorInfo(code codes.Code, errMsg string, errorInfo *errdetails.ErrorInfo) (*status.Status, error) {
	var st *status.Status
	if errMsg == "" {
		switch code {
		case codes.Unauthenticated:
			st = status.New(code, defaultUnauthenticatedErrMsg)
		case codes.PermissionDenied:
			st = status.New(code, defaultPermissionDeniedErrMsg)
		case codes.Aborted:
			st = status.New(code, defaultAbortedErrMsg)
		}
	} else {
		st = status.New(code, errMsg)
	}

	return st.WithDetails(errorInfo)
}

func NewPermissionDenied(errMsg string, errorInfo *errdetails.ErrorInfo) (*status.Status, error) {
	return newGRPCErrorWithErrorInfo(codes.PermissionDenied, errMsg, errorInfo)
}

func NewAborted(errMsg string, errorInfo *errdetails.ErrorInfo) (*status.Status, error) {
	return newGRPCErrorWithErrorInfo(codes.Aborted, errMsg, errorInfo)
}

func NewNotFound(errMsg string, resourceInfo *errdetails.ResourceInfo) (*status.Status, error) {
	return newGRPCErrorWithResourceInfo(codes.NotFound, errMsg, resourceInfo)
}

func NewAlreadyExists(errMsg string, resourceInfo *errdetails.ResourceInfo) (*status.Status, error) {
	return newGRPCErrorWithResourceInfo(codes.AlreadyExists, errMsg, resourceInfo)
}

func newGRPCErrorWithResourceInfo(code codes.Code, errMsg string, resourceInfo *errdetails.ResourceInfo) (*status.Status, error) {
	var st *status.Status
	if errMsg == "" {
		switch code {
		case codes.NotFound:
			st = status.New(code, defaultNotFoundErrMsg)
		case codes.AlreadyExists:
			st = status.New(code, defaultAlreadyExistsErrMsg)
		}
	} else {
		st = status.New(code, errMsg)
	}

	return st.WithDetails(resourceInfo)
}

func NewResourceExhausted(errMsg string, quotaViolations []*errdetails.QuotaFailure_Violation) (*status.Status, error) {
	return newGRPCErrorWithQuotaFailure(codes.ResourceExhausted, errMsg, quotaViolations)
}

func newGRPCErrorWithQuotaFailure(code codes.Code, errMsg string, violations []*errdetails.QuotaFailure_Violation) (*status.Status, error) {
	var st *status.Status
	if errMsg == "" {
		st = status.New(code, defaultResourceExhaustedErrMsg)
	} else {
		st = status.New(code, errMsg)
	}

	quotaFailureDetails := errdetails.QuotaFailure{Violations: violations}

	return st.WithDetails(&quotaFailureDetails)
}

func NewCancelled(errMsg string) (*status.Status, error) {
	var st *status.Status
	if errMsg == "" {
		st = status.New(codes.Canceled, defaultCanceledErrMsg)
	} else {
		st = status.New(codes.Canceled, errMsg)
	}

	return st, nil
}

func NewDataLoss(errMsg string, debugInfo *errdetails.DebugInfo) (*status.Status, error) {
	var st *status.Status
	if errMsg == "" {
		st = status.New(codes.DataLoss, defaultCanceledErrMsg)
	} else {
		st = status.New(codes.DataLoss, errMsg)
	}

	return st.WithDetails(debugInfo)
}

func newGRPCErrorWithDebugInfo(code codes.Code, errMsg string, debugInfo *errdetails.DebugInfo) (*status.Status, error) {
	var st *status.Status
	if errMsg == "" {
		switch code {
		case codes.DataLoss:
			st = status.New(code, defaultDataLossErrMsg)
		case codes.Unknown:
			st = status.New(code, defaultUnknownErrMsg)
		case codes.Internal:
			st = status.New(code, defaultInternalErrMsg)
		case codes.Unavailable:
			st = status.New(code, defaultUnavailableErrMsg)
		case codes.DeadlineExceeded:
			st = status.New(code, defaultDeadlineExceededErrMsg)
		}
	} else {
		st = status.New(code, errMsg)
	}

	return st.WithDetails(debugInfo)
}

func NewUnknown(errMsg string, debugInfo *errdetails.DebugInfo) (*status.Status, error) {
	return newGRPCErrorWithDebugInfo(codes.Unknown, errMsg, debugInfo)
}

func NewInternal(errMsg string, debugInfo *errdetails.DebugInfo) (*status.Status, error) {
	return newGRPCErrorWithDebugInfo(codes.Internal, errMsg, debugInfo)
}

func NewUnimplemented(errMsg string) (*status.Status, error) {
	var st *status.Status
	if errMsg == "" {
		st = status.New(codes.Unimplemented, defaultUnimplementedErrMsg)
	} else {
		st = status.New(codes.Unimplemented, errMsg)
	}

	return st, nil
}

func NewUnavailable(errMsg string, debugInfo *errdetails.DebugInfo) (*status.Status, error) {
	return newGRPCErrorWithDebugInfo(codes.Unavailable, errMsg, debugInfo)
}

func NewDeadlineExceeded(errMsg string, debugInfo *errdetails.DebugInfo) (*status.Status, error) {
	return newGRPCErrorWithDebugInfo(codes.DeadlineExceeded, errMsg, debugInfo)
}
