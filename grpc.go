package grpcerr

import (
	"fmt"

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

// Describes additional debugging info.
//
// Source: https://pkg.go.dev/google.golang.org/genproto/googleapis/rpc/errdetails
type DebugInfo struct {
	// The stack trace entries indicating where the error occurred.
	StackEntries []string
	// Additional debugging information provided by the server.
	Detail string
}

// AddDebugInfo adds additional debug info to a gRPC error. For example useful when the server
// wants to include a stack trace.
//
// Source: https://github.com/grpc/grpc-go/blob/master/codes/codes.go
func AddDebugInfo(gRPCErr error, debugInfo *DebugInfo) (error, error) {
	if debugInfo == nil {
		return gRPCErr, nil
	}

	status, ok := status.FromError(gRPCErr)
	if !ok {
		return nil, fmt.Errorf("invalid argument: gRPCErr must hold a status.Error struct")
	}

	errDetailsDebugInfo := &errdetails.DebugInfo{
		StackEntries: debugInfo.StackEntries,
		Detail:       debugInfo.Detail,
	}

	statusWithDebugInfo, err := status.WithDetails(errDetailsDebugInfo)
	if err != nil {
		return nil, err
	}

	return statusWithDebugInfo.Err(), nil
}

// DebugInfoFrom returns the DebugInfo from a gRPC error. If there isn't any,
// the zero value of DebugInfo is returned.
func DebugInfoFrom(gRPCErr error) DebugInfo {
	st := status.Convert(gRPCErr)

	for _, detail := range st.Details() {
		if debugInfo, ok := detail.(*errdetails.DebugInfo); ok {
			return DebugInfo{
				StackEntries: debugInfo.StackEntries,
				Detail:       debugInfo.Detail,
			}
		}
	}

	return DebugInfo{}
}

// Contains metadata about the request that clients can attach when filing a bug
// or providing other forms of feedback.
//
// Source: https://pkg.go.dev/google.golang.org/genproto/googleapis/rpc/errdetails
type RequestInfo struct {
	// An opaque string that should only be interpreted by the service generating
	// it. For example, it can be used to identify requests in the service's logs.
	RequestID string
	// Any data that was used to serve this request. For example, an encrypted
	// stack trace that can be sent back to the service provider for debugging.
	ServingData string
}

// AddRequestInfo adds metadata to a gRPC error about the request that clients
// can attach when filing a bug or providing other forms of feedback.
//
// Source: https://github.com/grpc/grpc-go/blob/master/codes/codes.go
func AddRequestInfo(gRPCErr error, requestInfo *RequestInfo) (error, error) {
	if requestInfo == nil {
		return gRPCErr, nil
	}

	status, ok := status.FromError(gRPCErr)
	if !ok {
		return nil, fmt.Errorf("invalid argument: gRPCErr must hold a status.Error struct")
	}

	requestInfoDetails := errdetails.RequestInfo{
		RequestId:   requestInfo.RequestID,
		ServingData: requestInfo.ServingData,
	}

	statusWithInfoDetails, err := status.WithDetails(&requestInfoDetails)
	if err != nil {
		return nil, err
	}

	return statusWithInfoDetails.Err(), nil
}

// RequestInfoFrom returns the RequestInfo from a gRPC error. If there's no
// RequestInfo details,the zero value of RequestInfo is returned.
func RequestInfoFrom(gRPCErr error) RequestInfo {
	st := status.Convert(gRPCErr)

	for _, detail := range st.Details() {
		if requestInfo, ok := detail.(*errdetails.RequestInfo); ok {
			return RequestInfo{
				RequestID:   requestInfo.RequestId,
				ServingData: requestInfo.ServingData,
			}
		}
	}

	return RequestInfo{}
}

// Provides a link to documentation or for performing an out of band action.
//
// For example, if a quota check failed with an error indicating the calling
// project hasn't enabled the accessed service, this can contain a URL pointing
// directly to the right place in the developer console to flip the bit.
//
// Source: https://pkg.go.dev/google.golang.org/genproto/googleapis/rpc/errdetails
type HelpLink struct {
	// Describes what the link offers.
	Description string
	// URL pointing to additional information on handling the current error.
	URL string
}

// AddHelp adds links to a gRPC error to documentation or for performing an out of band action.
//
// For example, if a quota check failed with an error indicating the calling
// project hasn't enabled the accessed service, this can contain a URL pointing
// directly to the right place in the developer console to flip the bit.
//
// Source: https://github.com/grpc/grpc-go/blob/master/codes/codes.go
func AddHelp(gRPCErr error, links []HelpLink) (error, error) {
	if len(links) == 0 {
		return gRPCErr, nil
	}

	status, ok := status.FromError(gRPCErr)
	if !ok {
		return nil, fmt.Errorf("invalid argument: gRPCErr must hold a status.Error struct")
	}

	helpDetails := errdetails.Help{}
	for _, link := range links {
		l := &errdetails.Help_Link{
			Description: link.Description,
			Url:         link.URL,
		}
		helpDetails.Links = append(helpDetails.Links, l)
	}

	statusWithHelpDetails, err := status.WithDetails(&helpDetails)
	if err != nil {
		return nil, err
	}

	return statusWithHelpDetails.Err(), nil
}

// HelpLinksFrom returns the slice of HelpLinks from a gRPC error. If there isn't any,
// an empty slice is returned.
func HelpLinksFrom(gRPCErr error) []HelpLink {
	st := status.Convert(gRPCErr)

	for _, detail := range st.Details() {
		if help, ok := detail.(*errdetails.Help); ok {
			helpLinks := make([]HelpLink, 0, len(help.Links))
			for _, link := range help.GetLinks() {
				l := HelpLink{
					Description: link.Description,
					URL:         link.Url,
				}
				helpLinks = append(helpLinks, l)
			}
			return helpLinks
		}
	}

	return []HelpLink{}
}

// Provides a localized error message that is safe to return to the user
// which can be attached to an RPC error.
//
// Source: https://pkg.go.dev/google.golang.org/genproto/googleapis/rpc/errdetails
type LocalizedMessage struct {
	// The locale used following the specification defined at
	// http://www.rfc-editor.org/rfc/bcp/bcp47.txt.
	// Examples are: "en-US", "fr-CH", "es-MX"
	Locale string
	// The localized error message in the above locale.
	Message string
}

// AddLocalizedMessage adds a localized error message to a gRPC error
func AddLocalizedMessage(gRPCErr error, localizedMsg *LocalizedMessage) (error, error) {
	if localizedMsg == nil {
		return gRPCErr, nil
	}

	status, ok := status.FromError(gRPCErr)
	if !ok {
		return nil, fmt.Errorf("invalid argument: gRPCErr must hold a status.Error struct")
	}

	localizedMessageDetails := errdetails.LocalizedMessage{
		Locale:  localizedMsg.Locale,
		Message: localizedMsg.Message,
	}

	statusWithLocalizedMsgDetails, err := status.WithDetails(&localizedMessageDetails)
	if err != nil {
		return nil, err
	}

	return statusWithLocalizedMsgDetails.Err(), nil
}

// LocalizedMessageFrom returns the LocalizedMessage from a gRPC error. If there isn't any,
// the zero value of LocalizedMessage is returned.
func LocalizedMessageFrom(gRPCErr error) LocalizedMessage {
	st := status.Convert(gRPCErr)

	for _, detail := range st.Details() {
		if localizedMsg, ok := detail.(*errdetails.LocalizedMessage); ok {
			return LocalizedMessage{
				Locale:  localizedMsg.Locale,
				Message: localizedMsg.Message,
			}
		}
	}

	return LocalizedMessage{}
}

// A message type used to describe a single bad request field.
//
// Source: https://pkg.go.dev/google.golang.org/genproto/googleapis/rpc/errdetails
type FieldViolation struct {
	// A path leading to a field in the request body. The value will be a
	// sequence of dot-separated identifiers that identify a protocol buffer
	// field. E.g., "field_violations.field" would identify this field.
	Field string
	// A description of why the request element is bad.
	Description string
}

// NewInvalidArgument constructs a gRPC error that indicates the client specified an invalid argument.
// Note that this differs from FailedPrecondition. It indicates arguments that are problematic regardless
// of the state of the system (e.g., a malformed file name).
//
// This error code will not be generated by the gRPC framework.
//
// Source: https://github.com/grpc/grpc-go/blob/master/codes/codes.go
func NewInvalidArgument(errMsg string, fieldViolations []FieldViolation) (error, error) {
	st, err := newStatusWithBadRequestDetails(codes.InvalidArgument, errMsg, fieldViolations)
	if err != nil {
		return nil, err
	}
	return st.Err(), nil
}

func newStatusWithBadRequestDetails(code codes.Code, errMsg string, fieldViolations []FieldViolation) (*status.Status, error) {
	var st *status.Status
	if errMsg == "" {
		switch code {
		case codes.InvalidArgument:
			st = status.New(code, defaultInvalidArgumentErrMsg)
		case codes.OutOfRange:
			st = status.New(code, defaultOutOfRangeErrMsg)
		}
	} else {
		st = status.New(code, errMsg)
	}

	if len(fieldViolations) == 0 {
		return st, nil
	}

	badRequestDetails := errdetails.BadRequest{}
	for _, violation := range fieldViolations {
		fv := &errdetails.BadRequest_FieldViolation{
			Field:       violation.Field,
			Description: violation.Description,
		}
		badRequestDetails.FieldViolations = append(badRequestDetails.FieldViolations, fv)
	}

	return st.WithDetails(&badRequestDetails)
}

// FieldViolationsFrom returns the slice of FieldViolations from a gRPC error. If there isn't any,
// an empty slice is returned.
func FieldViolationsFrom(gRPCErr error) []FieldViolation {
	st := status.Convert(gRPCErr)

	for _, detail := range st.Details() {
		if badReq, ok := detail.(*errdetails.BadRequest); ok {
			fieldViolations := make([]FieldViolation, 0, len(badReq.FieldViolations))
			for _, violation := range badReq.FieldViolations {
				fv := FieldViolation{
					Field:       violation.Field,
					Description: violation.Description,
				}
				fieldViolations = append(fieldViolations, fv)
			}
			return fieldViolations
		}
	}

	return []FieldViolation{}
}

func jsonBytesFromGrpcStatus(status *status.Status) ([]byte, error) {
	data, err := protojson.Marshal(status.Proto())
	if err != nil {
		return nil, err
	}

	return data, nil
}

// NewOutOfRange constructs a gRPC error that means the operation was
// attempted past the valid range.
// E.g., seeking or reading past end of file.
//
// Unlike InvalidArgument, this error indicates a problem that may
// be fixed if the system state changes. For example, a 32-bit file
// system will generate InvalidArgument if asked to read at an
// offset that is not in the range [0,2^32-1], but it will generate
// OutOfRange if asked to read from an offset past the current
// file size.
//
// There is a fair bit of overlap between FailedPrecondition and
// OutOfRange. We recommend using OutOfRange (the more specific
// error) when it applies so that callers who are iterating through
// a space can easily look for an OutOfRange error to detect when
// they are done.
//
// This error code will not be generated by the gRPC framework.
//
// Source: https://github.com/grpc/grpc-go/blob/master/codes/codes.go
func NewOutOfRange(errMsg string, fieldViolations []FieldViolation) (error, error) {
	st, err := newStatusWithBadRequestDetails(codes.OutOfRange, errMsg, fieldViolations)
	if err != nil {
		return nil, err
	}
	return st.Err(), nil
}

// A message type used to describe a single precondition failure.
//
// Source: https://pkg.go.dev/google.golang.org/genproto/googleapis/rpc/errdetails
type PreconditionFailure struct {
	// The type of PreconditionFailure. We recommend using a service-specific
	// enum type to define the supported precondition violation subjects. For
	// example, "TOS" for "Terms of Service violation".
	Type string
	// The subject, relative to the type, that failed.
	// For example, "google.com/cloud" relative to the "TOS" type would indicate
	// which terms of service is being referenced.
	Subject string
	// A description of how the precondition failed. Developers can use this
	// description to understand how to fix the failure.
	//
	// For example: "Terms of service not accepted".
	Description string
}

// NewFailedPrecondition constructs a gRPC error that indicates operation was rejected because the
// system is not in a state required for the operation's execution.
// For example, directory to be deleted may be non-empty, an rmdir
// operation is applied to a non-directory, etc.
//
// A litmus test that may help a service implementor in deciding
// between FailedPrecondition, Aborted, and Unavailable:
//  (a) Use Unavailable if the client can retry just the failing call.
//  (b) Use Aborted if the client should retry at a higher-level
//      (e.g., restarting a read-modify-write sequence).
//  (c) Use FailedPrecondition if the client should not retry until
//      the system state has been explicitly fixed. E.g., if an "rmdir"
//      fails because the directory is non-empty, FailedPrecondition
//      should be returned since the client should not retry unless
//      they have first fixed up the directory by deleting files from it.
//  (d) Use FailedPrecondition if the client performs conditional
//      REST Get/Update/Delete on a resource and the resource on the
//      server does not match the condition. E.g., conflicting
//      read-modify-write on the same resource.
//
// This error code will not be generated by the gRPC framework.
//
// Source: https://github.com/grpc/grpc-go/blob/master/codes/codes.go
func NewFailedPrecondition(errMsg string, failures []PreconditionFailure) (error, error) {
	st, err := newStatusWithFailedPreconditionDetails(codes.FailedPrecondition, errMsg, failures)
	if err != nil {
		return nil, err
	}
	return st.Err(), nil
}

func newStatusWithFailedPreconditionDetails(code codes.Code, errMsg string, failures []PreconditionFailure) (*status.Status, error) {
	var st *status.Status
	if errMsg == "" {
		st = status.New(code, defaultFailedPreconditionErrMsg)
	} else {
		st = status.New(code, errMsg)
	}

	if len(failures) == 0 {
		return st, nil
	}

	preconditionFailureDetails := errdetails.PreconditionFailure{}
	for _, failure := range failures {
		v := &errdetails.PreconditionFailure_Violation{
			Type:        failure.Type,
			Subject:     failure.Subject,
			Description: failure.Description,
		}
		preconditionFailureDetails.Violations = append(preconditionFailureDetails.Violations, v)
	}

	return st.WithDetails(&preconditionFailureDetails)
}

// PreconditionFailuresFrom returns the slice of PreconditionFailures from a gRPC error. If there isn't any,
// an empty slice is returned.
func PreconditionFailuresFrom(gRPCErr error) []PreconditionFailure {
	st := status.Convert(gRPCErr)

	for _, detail := range st.Details() {
		if precondFailure, ok := detail.(*errdetails.PreconditionFailure); ok {
			failures := make([]PreconditionFailure, 0, len(precondFailure.Violations))
			for _, violation := range precondFailure.Violations {
				failure := PreconditionFailure{
					Type:        violation.Type,
					Subject:     violation.Subject,
					Description: violation.Description,
				}
				failures = append(failures, failure)
			}
			return failures
		}
	}

	return []PreconditionFailure{}
}

// Describes the cause of the error with structured details.
//
// Example of an error when contacting the "pubsub.googleapis.com" API when it
// is not enabled:
//
//     { "reason": "API_DISABLED"
//       "domain": "googleapis.com"
//       "metadata": {
//         "resource": "projects/123",
//         "service": "pubsub.googleapis.com"
//       }
//     }
//
// This response indicates that the pubsub.googleapis.com API is not enabled.
//
// Example of an error that is returned when attempting to create a Spanner
// instance in a region that is out of stock:
//
//     { "reason": "STOCKOUT"
//       "domain": "spanner.googleapis.com",
//       "metadata": {
//         "availableRegions": "us-central1,us-east2"
//       }
//     }
//
// Source: https://pkg.go.dev/google.golang.org/genproto/googleapis/rpc/errdetails
type ErrorInfo struct {
	// The reason of the error. This is a constant value that identifies the
	// proximate cause of the error. Error reasons are unique within a particular
	// domain of errors. This should be at most 63 characters and match
	// /[A-Z0-9_]+/.
	Reason string
	// The logical grouping to which the "reason" belongs. The error domain
	// is typically the registered service name of the tool or product that
	// generates the error. Example: "pubsub.googleapis.com". If the error is
	// generated by some common infrastructure, the error domain must be a
	// globally unique value that identifies the infrastructure. For Google API
	// infrastructure, the error domain is "googleapis.com".
	Domain string
	// Additional structured details about this error.
	//
	// Keys should match /[a-zA-Z0-9-_]/ and be limited to 64 characters in
	// length. When identifying the current value of an exceeded limit, the units
	// should be contained in the key, not the value.  For example, rather than
	// {"instanceLimit": "100/request"}, should be returned as,
	// {"instanceLimitPerRequest": "100"}, if the client exceeds the number of
	// instances that can be created in a single (batch) request.
	Metadata map[string]string
}

// NewUnauthenticated constructs a gRPC error that indicates the request does not have valid
// authentication credentials for the operation.
//
// The gRPC framework will generate this error code when the
// authentication metadata is invalid or a Credentials callback fails,
// but also expect authentication middleware to generate it.
//
// Source: https://github.com/grpc/grpc-go/blob/master/codes/codes.go
func NewUnauthenticated(errMsg string, errorInfo *ErrorInfo) (error, error) {
	st, err := newStatusWithErrorInfo(codes.Unauthenticated, errMsg, errorInfo)
	if err != nil {
		return nil, err
	}
	return st.Err(), nil
}

func newStatusWithErrorInfo(code codes.Code, errMsg string, errorInfo *ErrorInfo) (*status.Status, error) {
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

	if errorInfo == nil {
		return st, nil
	}

	errorInfoDetails := errdetails.ErrorInfo{
		Reason:   errorInfo.Reason,
		Domain:   errorInfo.Domain,
		Metadata: errorInfo.Metadata,
	}

	return st.WithDetails(&errorInfoDetails)
}

// ErrorInfoFrom returns the ErrorInfo from a gRPC error. If there isn't any,
// the zero value of ErrorInfo is returned.
func ErrorInfoFrom(gRPCErr error) ErrorInfo {
	st := status.Convert(gRPCErr)

	for _, detail := range st.Details() {
		if errorInfo, ok := detail.(*errdetails.ErrorInfo); ok {
			return ErrorInfo{
				Reason:   errorInfo.Reason,
				Domain:   errorInfo.Domain,
				Metadata: errorInfo.Metadata,
			}
		}
	}

	return ErrorInfo{}
}

// NewPermissionDenied constructs a gRPC error that indicates the caller does not have permission to
// execute the specified operation. It must not be used for rejections
// caused by exhausting some resource (use ResourceExhausted
// instead for those errors). It must not be
// used if the caller cannot be identified (use Unauthenticated
// instead for those errors).
//
// This error code will not be generated by the gRPC core framework,
// but expect authentication middleware to use it.
//
// Source: https://github.com/grpc/grpc-go/blob/master/codes/codes.go
func NewPermissionDenied(errMsg string, errorInfo *ErrorInfo) (error, error) {
	st, err := newStatusWithErrorInfo(codes.PermissionDenied, errMsg, errorInfo)
	if err != nil {
		return nil, err
	}
	return st.Err(), nil
}

// NewAborted constructs a gRPC error that indicates the operation was aborted, typically due to a
// concurrency issue like sequencer check failures, transaction aborts,
// etc.
//
// See litmus test above for deciding between FailedPrecondition,
// Aborted, and Unavailable.
//
// This error code will not be generated by the gRPC framework.
//
// Source: https://github.com/grpc/grpc-go/blob/master/codes/codes.go
func NewAborted(errMsg string, errorInfo *ErrorInfo) (error, error) {
	st, err := newStatusWithErrorInfo(codes.Aborted, errMsg, errorInfo)
	if err != nil {
		return nil, err
	}
	return st.Err(), nil
}

// Describes the resource that is being accessed.
//
// Source: https://pkg.go.dev/google.golang.org/genproto/googleapis/rpc/errdetails
type ResourceInfo struct {
	// A name for the type of resource being accessed, e.g. "sql table",
	// "cloud storage bucket", "file", "Google calendar"; or the type URL
	// of the resource: e.g. "type.googleapis.com/google.pubsub.v1.Topic".
	ResourceType string
	// The name of the resource being accessed.  For example, a shared calendar
	// name: "example.com_4fghdhgsrgh@group.calendar.google.com", if the current
	// error is [google.rpc.Code.PERMISSION_DENIED][google.rpc.Code.PERMISSION_DENIED].
	ResourceName string
	// The owner of the resource (optional).
	// For example, "user:<owner email>" or "project:<Google developer project
	// id>".
	Owner string
	// Describes what error is encountered when accessing this resource.
	// For example, updating a cloud project may require the `writer` permission
	// on the developer console project.
	Description string
}

// NewNotFound constructs a gRPC error that means some requested entity (e.g., file or directory) was
// not found.
//
// This error code will not be generated by the gRPC framework.
//
// Source: https://github.com/grpc/grpc-go/blob/master/codes/codes.go
func NewNotFound(errMsg string, resourceInfo *ResourceInfo) (error, error) {
	st, err := newStatusWithResourceInfo(codes.NotFound, errMsg, resourceInfo)
	if err != nil {
		return nil, err
	}
	return st.Err(), nil
}

// NewAlreadyExists constructs a gRPC error that means an attempt to create an entity failed because one
// already exists.
//
// This error code will not be generated by the gRPC framework.
//
// Source: https://github.com/grpc/grpc-go/blob/master/codes/codes.go
func NewAlreadyExists(errMsg string, resourceInfo *ResourceInfo) (error, error) {
	st, err := newStatusWithResourceInfo(codes.AlreadyExists, errMsg, resourceInfo)
	if err != nil {
		return nil, err
	}
	return st.Err(), nil
}

func newStatusWithResourceInfo(code codes.Code, errMsg string, resourceInfo *ResourceInfo) (*status.Status, error) {
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

	if resourceInfo == nil {
		return st, nil
	}

	resourceInfoDetails := errdetails.ResourceInfo{
		ResourceType: resourceInfo.ResourceType,
		ResourceName: resourceInfo.ResourceName,
		Owner:        resourceInfo.Owner,
		Description:  resourceInfo.Description,
	}

	return st.WithDetails(&resourceInfoDetails)
}

// ResourceInfoFrom returns the ResourceInfo from a gRPC error. If there isn't any,
// the zero value of ResourceInfo is returned.
func ResourceInfoFrom(gRPCErr error) ResourceInfo {
	st := status.Convert(gRPCErr)

	for _, detail := range st.Details() {
		if resourceInfo, ok := detail.(*errdetails.ResourceInfo); ok {
			return ResourceInfo{
				ResourceType: resourceInfo.ResourceType,
				ResourceName: resourceInfo.ResourceName,
				Owner:        resourceInfo.Owner,
				Description:  resourceInfo.Description,
			}
		}
	}

	return ResourceInfo{}
}

// A message type used to describe a single quota violation.  For example, a
// daily quota or a custom quota that was exceeded.
//
// Source: https://pkg.go.dev/google.golang.org/genproto/googleapis/rpc/errdetails
type QuotaViolation struct {
	// The subject on which the quota check failed.
	// For example, "clientip:<ip address of client>" or "project:<Google
	// developer project id>".
	Subject string
	// A description of how the quota check failed. Clients can use this
	// description to find more about the quota configuration in the service's
	// public documentation, or find the relevant quota limit to adjust through
	// developer console.
	//
	// For example: "Service disabled" or "Daily Limit for read operations
	// exceeded".
	Description string
}

// NewResourceExhausted constructs a gRPC error that indicates some resource has been exhausted, perhaps
// a per-user quota, or perhaps the entire file system is out of space.
//
// This error code will be generated by the gRPC framework in
// out-of-memory and server overload situations, or when a message is
// larger than the configured maximum size.
//
// Source: https://github.com/grpc/grpc-go/blob/master/codes/codes.go
func NewResourceExhausted(errMsg string, quotaViolations []QuotaViolation) (error, error) {
	st, err := newStatusWithQuotaFailure(codes.ResourceExhausted, errMsg, quotaViolations)
	if err != nil {
		return nil, err
	}
	return st.Err(), nil
}

func newStatusWithQuotaFailure(code codes.Code, errMsg string, violations []QuotaViolation) (*status.Status, error) {
	var st *status.Status
	if errMsg == "" {
		st = status.New(code, defaultResourceExhaustedErrMsg)
	} else {
		st = status.New(code, errMsg)
	}

	if len(violations) == 0 {
		return st, nil
	}

	quotaFailureDetails := errdetails.QuotaFailure{}
	for _, violation := range violations {
		v := &errdetails.QuotaFailure_Violation{
			Subject:     violation.Subject,
			Description: violation.Description,
		}
		quotaFailureDetails.Violations = append(quotaFailureDetails.Violations, v)
	}

	return st.WithDetails(&quotaFailureDetails)
}

// QuotaViolationsFrom returns the slice of QuotaViolations from a gRPC error. If there isn't any,
// an empty slice is returned.
func QuotaViolationsFrom(gRPCErr error) []QuotaViolation {
	st := status.Convert(gRPCErr)

	for _, detail := range st.Details() {
		if quotaFailure, ok := detail.(*errdetails.QuotaFailure); ok {
			violations := make([]QuotaViolation, 0, len(quotaFailure.Violations))
			for _, violation := range quotaFailure.Violations {
				v := QuotaViolation{
					Subject:     violation.Subject,
					Description: violation.Description,
				}
				violations = append(violations, v)
			}
			return violations
		}
	}

	return []QuotaViolation{}
}

// NewCancelled constructs a gRPC error that indicates the operation was canceled (typically by the caller).
//
// The gRPC framework will generate this error code when cancellation
// is requested.
//
// Source: https://github.com/grpc/grpc-go/blob/master/codes/codes.go
func NewCancelled(errMsg string) error {
	var st *status.Status
	if errMsg == "" {
		st = status.New(codes.Canceled, defaultCanceledErrMsg)
	} else {
		st = status.New(codes.Canceled, errMsg)
	}

	return st.Err()
}

// NewDataLoss constructs a gRPC error that indicates unrecoverable data loss or corruption.
//
// This error code will not be generated by the gRPC framework.
//
// Source: https://github.com/grpc/grpc-go/blob/master/codes/codes.go
func NewDataLoss(errMsg string, debugInfo *DebugInfo) (error, error) {
	var st *status.Status
	if errMsg == "" {
		st = status.New(codes.DataLoss, defaultCanceledErrMsg)
	} else {
		st = status.New(codes.DataLoss, errMsg)
	}

	if debugInfo == nil {
		return st.Err(), nil
	}

	debugInfoDetails := errdetails.DebugInfo{
		StackEntries: debugInfo.StackEntries,
		Detail:       debugInfo.Detail,
	}

	statusWithDetails, err := st.WithDetails(&debugInfoDetails)
	if err != nil {
		return nil, err
	}

	return statusWithDetails.Err(), nil
}

func newStatusWithDebugInfo(code codes.Code, errMsg string, debugInfo *DebugInfo) (*status.Status, error) {
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

	if debugInfo == nil {
		return st, nil
	}

	debugInfoDetails := errdetails.DebugInfo{
		StackEntries: debugInfo.StackEntries,
		Detail:       debugInfo.Detail,
	}

	return st.WithDetails(&debugInfoDetails)
}

// NewUnknown constructs a gRPC error that means an unknown error has occured.
// An example of where this error may be returned is
// if a Status value received from another address space belongs to
// an error-space that is not known in this address space. Also
// errors raised by APIs that do not return enough error information
// may be converted to this error.
//
// The gRPC framework will generate this error code in the above two
// mentioned cases.
//
// Source: https://github.com/grpc/grpc-go/blob/master/codes/codes.go
func NewUnknown(errMsg string, debugInfo *DebugInfo) (error, error) {
	st, err := newStatusWithDebugInfo(codes.Unknown, errMsg, debugInfo)
	if err != nil {
		return nil, err
	}
	return st.Err(), nil
}

// NewInternal construct a gRPC error that means some invariants expected by underlying
// system has been broken. If you see one of these errors,
// something is very broken.
//
// This error code will be generated by the gRPC framework in several
// internal error conditions.
//
// Source: https://github.com/grpc/grpc-go/blob/master/codes/codes.go
func NewInternal(errMsg string, debugInfo *DebugInfo) (error, error) {
	st, err := newStatusWithDebugInfo(codes.Internal, errMsg, debugInfo)
	if err != nil {
		return nil, err
	}
	return st.Err(), nil
}

// NewUnimplemented constructs a gRPC error that indicates operation is not implemented or not
// supported/enabled in this service.
//
// This error code will be generated by the gRPC framework. Most
// commonly, you will see this error code when a method implementation
// is missing on the server. It can also be generated for unknown
// compression algorithms or a disagreement as to whether an RPC should
// be streaming.
//
// Source: https://github.com/grpc/grpc-go/blob/master/codes/codes.go
func NewUnimplemented(errMsg string) error {
	var st *status.Status
	if errMsg == "" {
		st = status.New(codes.Unimplemented, defaultUnimplementedErrMsg)
	} else {
		st = status.New(codes.Unimplemented, errMsg)
	}

	return st.Err()
}

// NewUnavailable constructs a gRPC error that indicates the service is currently unavailable.
// This is a most likely a transient condition and may be corrected
// by retrying with a backoff. Note that it is not always safe to retry
// non-idempotent operations.
//
// See litmus test above for deciding between FailedPrecondition,
// Aborted, and Unavailable.
//
// This error code will be generated by the gRPC framework during
// abrupt shutdown of a server process or network connection.
//
// Source: https://github.com/grpc/grpc-go/blob/master/codes/codes.go
func NewUnavailable(errMsg string, debugInfo *DebugInfo) (error, error) {
	st, err := newStatusWithDebugInfo(codes.Unavailable, errMsg, debugInfo)
	if err != nil {
		return nil, err
	}
	return st.Err(), nil
}

// NewDeadlineExceeded constructs a gRPC error that means operation expired before completion.
// For operations that change the state of the system, this error may be
// returned even if the operation has completed successfully. For
// example, a successful response from a server could have been delayed
// long enough for the deadline to expire.
//
// The gRPC framework will generate this error code when the deadline is
// exceeded.
//
// Source: https://github.com/grpc/grpc-go/blob/master/codes/codes.go
func NewDeadlineExceeded(errMsg string, debugInfo *DebugInfo) (error, error) {
	st, err := newStatusWithDebugInfo(codes.DeadlineExceeded, errMsg, debugInfo)
	if err != nil {
		return nil, err
	}
	return st.Err(), nil
}

func Code(gRPCErr error) codes.Code {
	st := status.Convert(gRPCErr)
	return st.Code()
}

func Message(gRPCErr error) string {
	st := status.Convert(gRPCErr)
	return st.Message()
}
