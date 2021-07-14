package grpcerr

import (
	"testing"

	"github.com/tobbstr/testa/assert"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func TestAddDebugInfo(t *testing.T) {
	validGRPCErr, err := NewUnimplemented("dummy-err-msg")
	if err != nil {
		t.Fatal(err)
	}
	debugInfo := &errdetails.DebugInfo{
		StackEntries: []string{"dummy-stack-entry"},
		Detail:       "dummy-detail",
	}

	gRPCErrWithDebugInfo := status.New(codes.Unimplemented, "dummy-err-msg")
	di := errdetails.DebugInfo{
		StackEntries: []string{"dummy-stack-entry"},
		Detail:       "dummy-detail",
	}
	gRPCErrWithDebugInfo, err = gRPCErrWithDebugInfo.WithDetails(&di)
	if err != nil {
		t.Fatal(err)
	}

	type args struct {
		gRPCErr   *status.Status
		debugInfo *errdetails.DebugInfo
	}
	tests := []struct {
		name    string
		args    args
		want    *status.Status
		wantErr bool
	}{
		{
			name: "should return gRPC error with debugInfo for valid arguments",
			args: args{
				gRPCErr:   validGRPCErr,
				debugInfo: debugInfo,
			},
			want:    gRPCErrWithDebugInfo,
			wantErr: false,
		},
		{
			name: "should return error when get nil gRPCErr argument",
			args: args{
				gRPCErr:   nil,
				debugInfo: debugInfo,
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Given
			assert := assert.New(t)

			// When
			got, err := AddDebugInfo(tt.args.gRPCErr, tt.args.debugInfo)

			// Then
			assert(err).IsWantedError(tt.wantErr)
			assert(got).Equals(tt.want)
		})
	}
}

func TestAddRequestInfo(t *testing.T) {
	validGRPCErr, err := NewUnimplemented("dummy-err-msg")
	if err != nil {
		t.Fatal(err)
	}
	requestID := "dummy-request-id"
	servingData := "dummy-serving-data"

	gRPCErrWithRequestInfo := status.New(codes.Unimplemented, "dummy-err-msg")
	di := errdetails.RequestInfo{
		RequestId:   requestID,
		ServingData: servingData,
	}
	gRPCErrWithRequestInfo, err = gRPCErrWithRequestInfo.WithDetails(&di)
	if err != nil {
		t.Fatal(err)
	}

	type args struct {
		gRPCErr     *status.Status
		requestID   string
		servingData string
	}
	tests := []struct {
		name    string
		args    args
		want    *status.Status
		wantErr bool
	}{
		{
			name: "should return gRPC error with Request Info for valid arguments",
			args: args{
				gRPCErr:     validGRPCErr,
				requestID:   "dummy-request-id",
				servingData: "dummy-serving-data",
			},
			want:    gRPCErrWithRequestInfo,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Given
			assert := assert.New(t)

			// When
			got, err := AddRequestInfo(tt.args.gRPCErr, tt.args.requestID, tt.args.servingData)

			// Then
			assert(err).IsWantedError(tt.wantErr)
			assert(got).Equals(tt.want)
		})
	}
}

func TestAddHelp(t *testing.T) {
	unimplemented, err := NewUnimplemented("dummy-msg")
	if err != nil {
		t.Fatal(err)
	}

	link := &errdetails.Help_Link{Description: "dummy-description", Url: "dummy-url"}

	links := []*errdetails.Help_Link{link}
	helpDetails := errdetails.Help{Links: links}

	unimplementedWithHelpDetails, err := unimplemented.WithDetails(&helpDetails)
	if err != nil {
		t.Fatal(err)
	}

	type args struct {
		gRPCErr *status.Status
		links   []*errdetails.Help_Link
	}
	tests := []struct {
		name    string
		args    args
		want    *status.Status
		wantErr bool
	}{
		{
			name: "should return passed gRPCErr arg for empty links slice",
			args: args{
				gRPCErr: unimplemented,
				links:   []*errdetails.Help_Link{},
			},
			want:    unimplemented,
			wantErr: false,
		},
		{
			name: "should return gRPC error with help links when get valid arguments",
			args: args{
				gRPCErr: unimplemented,
				links:   links,
			},
			want:    unimplementedWithHelpDetails,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Given
			assert := assert.New(t)

			// When
			got, err := AddHelp(tt.args.gRPCErr, tt.args.links)

			// Then
			assert(err).IsWantedError(tt.wantErr)
			assert(got).Equals(tt.want)
		})
	}
}

func TestAddLocalizedMessage(t *testing.T) {
	unimplemented, err := NewUnimplemented("dummy-msg")
	if err != nil {
		t.Fatal(err)
	}

	locale := "dummy-locale"
	msg := "dummy-localized-message"
	localizedMessageDetails := errdetails.LocalizedMessage{
		Locale:  locale,
		Message: msg,
	}

	unimplementedWithLocalizedMessageDetails, err := unimplemented.WithDetails(&localizedMessageDetails)
	if err != nil {
		t.Fatal(err)
	}

	type args struct {
		gRPCErr *status.Status
		locale  string
		msg     string
	}
	tests := []struct {
		name    string
		args    args
		want    *status.Status
		wantErr bool
	}{
		{
			name: "should return gRPC error with localized message when get valid arguments",
			args: args{
				gRPCErr: unimplemented,
				locale:  locale,
				msg:     msg,
			},
			want:    unimplementedWithLocalizedMessageDetails,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Given
			assert := assert.New(t)

			// When
			got, err := AddLocalizedMessage(tt.args.gRPCErr, tt.args.locale, tt.args.msg)

			// Then
			assert(err).IsWantedError(tt.wantErr)
			assert(got).Equals(tt.want)
		})
	}
}

func TestNewInvalidArgument(t *testing.T) {
	invalidArgumentWithDefaultMsg := status.New(codes.InvalidArgument, defaultInvalidArgumentErrMsg)
	invalidArgumentWithoutDetails := status.New(codes.InvalidArgument, "dummy-message")
	type args struct {
		errMsg          string
		fieldViolations []*errdetails.BadRequest_FieldViolation
	}
	tests := []struct {
		name    string
		args    args
		want    *status.Status
		wantErr bool
	}{
		{
			name: "should return gRPC error with default errMsg when get empty errMsg arg",
			args: args{
				errMsg:          "",
				fieldViolations: nil,
			},
			want:    invalidArgumentWithDefaultMsg,
			wantErr: false,
		},
		{
			name: "should return gRPC error without details when fieldViolations arg is nil",
			args: args{
				errMsg:          "dummy-message",
				fieldViolations: nil,
			},
			want:    invalidArgumentWithoutDetails,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Given
			assert := assert.New(t)

			// When
			got, err := NewInvalidArgument(tt.args.errMsg, tt.args.fieldViolations)

			// Then
			assert(err).IsWantedError(tt.wantErr)
			assert(got).Equals(tt.want)
		})
	}
}

func Test_newGRPCErrorWithErrorInfo(t *testing.T) {
	unauthenticatedWithZeroErrorInfo := status.New(codes.Unauthenticated, defaultUnauthenticatedErrMsg)
	unauthenticatedWithZeroErrorInfo, err := unauthenticatedWithZeroErrorInfo.WithDetails(&errdetails.ErrorInfo{})
	if err != nil {
		t.Fatal(err)
	}
	permissionDeniedWithZeroErrorInfo := status.New(codes.PermissionDenied, defaultPermissionDeniedErrMsg)
	permissionDeniedWithZeroErrorInfo, err = permissionDeniedWithZeroErrorInfo.WithDetails(&errdetails.ErrorInfo{})
	if err != nil {
		t.Fatal(err)
	}
	abortedWithZeroErrorInfo := status.New(codes.Aborted, defaultAbortedErrMsg)
	abortedWithZeroErrorInfo, err = abortedWithZeroErrorInfo.WithDetails(&errdetails.ErrorInfo{})
	if err != nil {
		t.Fatal(err)
	}

	type args struct {
		code      codes.Code
		errMsg    string
		errorInfo *errdetails.ErrorInfo
	}
	tests := []struct {
		name    string
		args    args
		want    *status.Status
		wantErr bool
	}{
		{
			name: "should return Unauthenticated gRPC error when get code codes.Unauthenticated",
			args: args{
				code:      codes.Unauthenticated,
				errMsg:    "",
				errorInfo: &errdetails.ErrorInfo{},
			},
			want:    unauthenticatedWithZeroErrorInfo,
			wantErr: false,
		},
		{
			name: "should return PermissionDenied gRPC error when get code codes.PermissionDenied",
			args: args{
				code:      codes.PermissionDenied,
				errMsg:    "",
				errorInfo: &errdetails.ErrorInfo{},
			},
			want:    permissionDeniedWithZeroErrorInfo,
			wantErr: false,
		},
		{
			name: "should return Aborted gRPC error when get code codes.Aborted",
			args: args{
				code:      codes.Aborted,
				errMsg:    "",
				errorInfo: &errdetails.ErrorInfo{},
			},
			want:    abortedWithZeroErrorInfo,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Given
			assert := assert.New(t)

			// When
			got, err := newGRPCErrorWithErrorInfo(tt.args.code, tt.args.errMsg, tt.args.errorInfo)

			// Then
			assert(err).IsWantedError(tt.wantErr)
			assert(got).Equals(tt.want)
		})
	}
}

func Test_newGRPCErrorWithResourceInfo(t *testing.T) {
	notFoundWithZeroResourceInfo := status.New(codes.NotFound, defaultNotFoundErrMsg)
	notFoundWithZeroResourceInfo, err := notFoundWithZeroResourceInfo.WithDetails(&errdetails.ResourceInfo{})
	if err != nil {
		t.Fatal(err)
	}
	alreadyExistsWithZeroResourceInfo := status.New(codes.AlreadyExists, defaultAlreadyExistsErrMsg)
	alreadyExistsWithZeroResourceInfo, err = alreadyExistsWithZeroResourceInfo.WithDetails(&errdetails.ResourceInfo{})
	if err != nil {
		t.Fatal(err)
	}

	type args struct {
		code         codes.Code
		errMsg       string
		resourceInfo *errdetails.ResourceInfo
	}
	tests := []struct {
		name    string
		args    args
		want    *status.Status
		wantErr bool
	}{
		{
			name: "should return NotFound gRPC error when get code codes.NotFound",
			args: args{
				code:         codes.NotFound,
				errMsg:       "",
				resourceInfo: &errdetails.ResourceInfo{},
			},
			want:    notFoundWithZeroResourceInfo,
			wantErr: false,
		},
		{
			name: "should return AlreadyExists gRPC error when get code codes.AlreadyExists",
			args: args{
				code:         codes.AlreadyExists,
				errMsg:       "",
				resourceInfo: &errdetails.ResourceInfo{},
			},
			want:    alreadyExistsWithZeroResourceInfo,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Given
			assert := assert.New(t)

			// When
			got, err := newGRPCErrorWithResourceInfo(tt.args.code, tt.args.errMsg, tt.args.resourceInfo)

			// Then
			assert(err).IsWantedError(tt.wantErr)
			assert(got).Equals(tt.want)
		})
	}
}

func Test_newGRPCErrorWithDebugInfo(t *testing.T) {
	dataLossWithZeroDebugInfo := status.New(codes.DataLoss, defaultDataLossErrMsg)
	dataLossWithZeroDebugInfo, err := dataLossWithZeroDebugInfo.WithDetails(&errdetails.DebugInfo{})
	if err != nil {
		t.Fatal(err)
	}
	unknownWithZeroDebugInfo := status.New(codes.Unknown, defaultUnknownErrMsg)
	unknownWithZeroDebugInfo, err = unknownWithZeroDebugInfo.WithDetails(&errdetails.DebugInfo{})
	if err != nil {
		t.Fatal(err)
	}
	internalWithZeroDebugInfo := status.New(codes.Internal, defaultInternalErrMsg)
	internalWithZeroDebugInfo, err = internalWithZeroDebugInfo.WithDetails(&errdetails.DebugInfo{})
	if err != nil {
		t.Fatal(err)
	}
	unavailableWithZeroDebugInfo := status.New(codes.Unavailable, defaultUnavailableErrMsg)
	unavailableWithZeroDebugInfo, err = unavailableWithZeroDebugInfo.WithDetails(&errdetails.DebugInfo{})
	if err != nil {
		t.Fatal(err)
	}
	deadlineExceededWithZeroDebugInfo := status.New(codes.DeadlineExceeded, defaultDeadlineExceededErrMsg)
	deadlineExceededWithZeroDebugInfo, err = deadlineExceededWithZeroDebugInfo.WithDetails(&errdetails.DebugInfo{})
	if err != nil {
		t.Fatal(err)
	}

	type args struct {
		code      codes.Code
		errMsg    string
		debugInfo *errdetails.DebugInfo
	}
	tests := []struct {
		name    string
		args    args
		want    *status.Status
		wantErr bool
	}{
		{
			name: "should return DataLoss gRPC error when get code codes.DataLoss",
			args: args{
				code:      codes.DataLoss,
				errMsg:    "",
				debugInfo: &errdetails.DebugInfo{},
			},
			want:    dataLossWithZeroDebugInfo,
			wantErr: false,
		},
		{
			name: "should return Unknown gRPC error when get code codes.Unknown",
			args: args{
				code:      codes.Unknown,
				errMsg:    "",
				debugInfo: &errdetails.DebugInfo{},
			},
			want:    unknownWithZeroDebugInfo,
			wantErr: false,
		},
		{
			name: "should return Internal gRPC error when get code codes.Internal",
			args: args{
				code:      codes.Internal,
				errMsg:    "",
				debugInfo: &errdetails.DebugInfo{},
			},
			want:    internalWithZeroDebugInfo,
			wantErr: false,
		},
		{
			name: "should return Unavailable gRPC error when get code codes.Unavailable",
			args: args{
				code:      codes.Unavailable,
				errMsg:    "",
				debugInfo: &errdetails.DebugInfo{},
			},
			want:    unavailableWithZeroDebugInfo,
			wantErr: false,
		},
		{
			name: "should return DeadlineExceeded gRPC error when get code codes.DeadlineExceeded",
			args: args{
				code:      codes.DeadlineExceeded,
				errMsg:    "",
				debugInfo: &errdetails.DebugInfo{},
			},
			want:    deadlineExceededWithZeroDebugInfo,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Given
			assert := assert.New(t)

			// When
			got, err := newGRPCErrorWithDebugInfo(tt.args.code, tt.args.errMsg, tt.args.debugInfo)

			// Then
			assert(err).IsWantedError(tt.wantErr)
			assert(got).Equals(tt.want)
		})
	}
}

func Test_newGRPCErrorWithQuotaFailure(t *testing.T) {
	resourceExhaustedWithDefaultErrMsg := status.New(codes.ResourceExhausted, defaultResourceExhaustedErrMsg)
	quotaFailureViolations := []*errdetails.QuotaFailure_Violation{
		{Subject: "dummy-subject", Description: "dummy-description"},
	}
	quotaFailureDetails := errdetails.QuotaFailure{
		Violations: quotaFailureViolations,
	}
	resourceExhaustedWithDefaultErrMsgAndQuotaFailureDetails, err := resourceExhaustedWithDefaultErrMsg.WithDetails(&quotaFailureDetails)
	if err != nil {
		t.Fatal(err)
	}

	type args struct {
		code       codes.Code
		errMsg     string
		violations []*errdetails.QuotaFailure_Violation
	}
	tests := []struct {
		name    string
		args    args
		want    *status.Status
		wantErr bool
	}{
		{
			name: "should return ResourceExhausted gRPC error with default errMsg when get empty errMsg arg",
			args: args{
				code:       codes.ResourceExhausted,
				errMsg:     "",
				violations: quotaFailureViolations,
			},
			want:    resourceExhaustedWithDefaultErrMsgAndQuotaFailureDetails,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Given
			assert := assert.New(t)

			// When
			got, err := newGRPCErrorWithQuotaFailure(tt.args.code, tt.args.errMsg, tt.args.violations)

			// Then
			assert(err).IsWantedError(tt.wantErr)
			assert(got).Equals(tt.want)
		})
	}
}
