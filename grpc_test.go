package grpcerr

import (
	"fmt"
	"testing"

	"github.com/tobbstr/testa/assert"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func TestAddDebugInfo(t *testing.T) {
	validGRPCErr := NewUnimplemented("dummy-err-msg")

	debugInfo := &DebugInfo{
		StackEntries: []string{"dummy-stack-entry"},
		Detail:       "dummy-detail",
	}

	statusWithDebugInfo := status.New(codes.Unimplemented, "dummy-err-msg")
	di := errdetails.DebugInfo{
		StackEntries: []string{"dummy-stack-entry"},
		Detail:       "dummy-detail",
	}
	statusWithDebugInfo, err := statusWithDebugInfo.WithDetails(&di)
	if err != nil {
		t.Fatal(err)
	}

	gRPCErrWithDebugInfo := statusWithDebugInfo.Err()

	type args struct {
		gRPCErr   error
		debugInfo *DebugInfo
	}
	tests := []struct {
		name    string
		args    args
		want    error
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
		{
			name: "should return error when get gRPCErr which does not have a GRPCStatus() method",
			args: args{
				gRPCErr:   fmt.Errorf("dummy-error"),
				debugInfo: debugInfo,
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "should return same gRPCErr when get nil debugInfo",
			args: args{
				gRPCErr:   validGRPCErr,
				debugInfo: nil,
			},
			want:    validGRPCErr,
			wantErr: false,
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

func TestDebugInfoFrom(t *testing.T) {
	stackEntries := []string{"dummy-stack-entry-1", "dummy-stack-entry-2"}
	detail := "dummy-detail"
	debugInfoDetails := &errdetails.DebugInfo{
		StackEntries: stackEntries,
		Detail:       detail,
	}

	status := status.New(codes.Unimplemented, defaultUnimplementedErrMsg)
	gRPCErrWithoutDebugInfo := status.Err()

	zeroDebugInfo := DebugInfo{}

	statusWithDebugInfo, err := status.WithDetails(debugInfoDetails)
	if err != nil {
		t.Fatal(err)
	}
	gRPCErrWithDebugInfo := statusWithDebugInfo.Err()

	debugInfo := DebugInfo{
		StackEntries: stackEntries,
		Detail:       detail,
	}

	type args struct {
		gRPCErr error
	}
	tests := []struct {
		name string
		args args
		want DebugInfo
	}{
		{
			name: "Should return debugInfo when get gRPCErr with debugInfoDetails",
			args: args{
				gRPCErrWithDebugInfo,
			},
			want: debugInfo,
		},
		{
			name: "Should return zeroDebugInfo when get gRPCErr without debugInfoDetails",
			args: args{
				gRPCErrWithoutDebugInfo,
			},
			want: zeroDebugInfo,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Given
			assert := assert.New(t)

			// When
			got := DebugInfoFrom(tt.args.gRPCErr)

			// Then
			assert(got).Equals(tt.want)
		})
	}
}

func TestAddRequestInfo(t *testing.T) {
	validGRPCErr := NewUnimplemented("dummy-err-msg")

	requestID := "dummy-request-id"
	servingData := "dummy-serving-data"

	statusWithRequestInfo := status.New(codes.Unimplemented, "dummy-err-msg")
	di := errdetails.RequestInfo{
		RequestId:   requestID,
		ServingData: servingData,
	}
	statusWithRequestInfo, err := statusWithRequestInfo.WithDetails(&di)
	if err != nil {
		t.Fatal(err)
	}

	gRPCErrWithRequestInfo := statusWithRequestInfo.Err()

	type args struct {
		gRPCErr     error
		requestInfo *RequestInfo
	}
	tests := []struct {
		name    string
		args    args
		want    error
		wantErr bool
	}{
		{
			name: "should return gRPC error with Request Info for valid arguments",
			args: args{
				gRPCErr: validGRPCErr,
				requestInfo: &RequestInfo{
					RequestID:   "dummy-request-id",
					ServingData: "dummy-serving-data",
				},
			},
			want:    gRPCErrWithRequestInfo,
			wantErr: false,
		},
		{
			name: "should return error when get gRPCErr which does not have a GRPCStatus() method",
			args: args{
				gRPCErr: fmt.Errorf("dummy-error"),
				requestInfo: &RequestInfo{
					RequestID:   "dummy-request-id",
					ServingData: "dummy-serving-data",
				},
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "should return error when get nil gRPCErr",
			args: args{
				gRPCErr: nil,
				requestInfo: &RequestInfo{
					RequestID:   "dummy-request-id",
					ServingData: "dummy-serving-data",
				},
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "should return unmodified gRPCErr when get nil requestInfo",
			args: args{
				gRPCErr:     validGRPCErr,
				requestInfo: nil,
			},
			want:    validGRPCErr,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Given
			assert := assert.New(t)

			// When
			got, err := AddRequestInfo(tt.args.gRPCErr, tt.args.requestInfo)

			// Then
			assert(err).IsWantedError(tt.wantErr)
			assert(got).Equals(tt.want)
		})
	}
}

func TestRequestInfoFrom(t *testing.T) {
	requestID := "dummy-request-id"
	servingData := "dummy-servingdata"
	requestInfoDetails := &errdetails.RequestInfo{
		RequestId:   requestID,
		ServingData: servingData,
	}

	status := status.New(codes.Unimplemented, defaultUnimplementedErrMsg)
	gRPCErrWithoutRequestInfo := status.Err()

	zeroRequestInfo := RequestInfo{}

	statusWithRequestInfo, err := status.WithDetails(requestInfoDetails)
	if err != nil {
		t.Fatal(err)
	}
	gRPCErrWithRequestInfo := statusWithRequestInfo.Err()

	requestInfo := RequestInfo{
		RequestID:   requestID,
		ServingData: servingData,
	}

	type args struct {
		gRPCErr error
	}
	tests := []struct {
		name string
		args args
		want RequestInfo
	}{
		{
			name: "Should return RequestInfo when get gRPCErr with requestInfoDetails",
			args: args{
				gRPCErr: gRPCErrWithRequestInfo,
			},
			want: requestInfo,
		},
		{
			name: "Should return zeroRequestInfo when get gRPCErr without requestInfoDetails",
			args: args{
				gRPCErr: gRPCErrWithoutRequestInfo,
			},
			want: zeroRequestInfo,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Given
			assert := assert.New(t)

			// When
			got := RequestInfoFrom(tt.args.gRPCErr)

			// Then
			assert(got).Equals(tt.want)
		})
	}
}

func TestAddHelp(t *testing.T) {
	unimplemented := NewUnimplemented("dummy-msg")

	links := []HelpLink{
		{Description: "dummy-description", URL: "dummy-url"},
	}

	errDetailsLink := &errdetails.Help_Link{Description: "dummy-description", Url: "dummy-url"}

	errDetailsLinks := []*errdetails.Help_Link{errDetailsLink}
	helpDetails := errdetails.Help{Links: errDetailsLinks}

	statusWithHelpDetails := status.New(codes.Unimplemented, "dummy-msg")
	statusWithHelpDetails, err := statusWithHelpDetails.WithDetails(&helpDetails)
	if err != nil {
		t.Fatal(err)
	}

	gRPCErrWithHelpDetails := statusWithHelpDetails.Err()

	type args struct {
		gRPCErr error
		links   []HelpLink
	}
	tests := []struct {
		name    string
		args    args
		want    error
		wantErr bool
	}{
		{
			name: "should return passed gRPCErr arg for empty links slice",
			args: args{
				gRPCErr: unimplemented,
				links:   []HelpLink{},
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
			want:    gRPCErrWithHelpDetails,
			wantErr: false,
		},
		{
			name: "should return error when get nil gRPCErr",
			args: args{
				gRPCErr: nil,
				links:   links,
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "should return error when get gRPCErr which does not have a GRPCStatus() method",
			args: args{
				gRPCErr: fmt.Errorf("dummy-error"),
				links:   links,
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
			got, err := AddHelp(tt.args.gRPCErr, tt.args.links)

			// Then
			assert(err).IsWantedError(tt.wantErr)
			assert(got).Equals(tt.want)
		})
	}
}

func TestHelpLinksFrom(t *testing.T) {
	description1 := "dummy-description-1"
	url1 := "dummy-url-1"
	description2 := "dummy-description-2"
	url2 := "dummy-url-2"

	helpInfoLink1 := &errdetails.Help_Link{
		Description: description1,
		Url:         url1,
	}
	helpInfoLink2 := &errdetails.Help_Link{
		Description: description2,
		Url:         url2,
	}
	helpInfoDetails := &errdetails.Help{Links: []*errdetails.Help_Link{helpInfoLink1, helpInfoLink2}}

	status := status.New(codes.Unimplemented, defaultUnimplementedErrMsg)
	gRPCErrWithoutHelpInfo := status.Err()

	zeroHelpLinks := []HelpLink{}

	statusWithHelpInfo, err := status.WithDetails(helpInfoDetails)
	if err != nil {
		t.Fatal(err)
	}
	gRPCErrWithHelpInfo := statusWithHelpInfo.Err()

	helpLinks := []HelpLink{
		{Description: description1, URL: url1},
		{Description: description2, URL: url2},
	}

	type args struct {
		gRPCErr error
	}
	tests := []struct {
		name string
		args args
		want []HelpLink
	}{
		{
			name: "should return helpLinks when get gRPCErr with helpInfoDetails",
			args: args{
				gRPCErr: gRPCErrWithHelpInfo,
			},
			want: helpLinks,
		},
		{
			name: "should return zeroHelpLinks when get gRPCErr without helpInfoDetails",
			args: args{
				gRPCErr: gRPCErrWithoutHelpInfo,
			},
			want: zeroHelpLinks,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Given
			assert := assert.New(t)

			// When
			got := HelpLinksFrom(tt.args.gRPCErr)

			// Then
			assert(got).Equals(tt.want)
		})
	}
}

func TestAddLocalizedMessage(t *testing.T) {
	unimplemented := NewUnimplemented("dummy-msg")

	locale := "dummy-locale"
	msg := "dummy-localized-message"
	localizedMessageDetails := errdetails.LocalizedMessage{
		Locale:  locale,
		Message: msg,
	}

	unimplementedWithLocalizedMessageDetails := status.New(codes.Unimplemented, "dummy-msg")
	unimplementedWithLocalizedMessageDetails, err := unimplementedWithLocalizedMessageDetails.WithDetails(&localizedMessageDetails)
	if err != nil {
		t.Fatal(err)
	}

	gRPCErrWithLocalizedMessageDetails := unimplementedWithLocalizedMessageDetails.Err()

	type args struct {
		gRPCErr          error
		localizedMessage *LocalizedMessage
	}
	tests := []struct {
		name    string
		args    args
		want    error
		wantErr bool
	}{
		{
			name: "should return gRPC error with localized message when get valid arguments",
			args: args{
				gRPCErr: unimplemented,
				localizedMessage: &LocalizedMessage{
					Locale:  locale,
					Message: msg,
				},
			},
			want:    gRPCErrWithLocalizedMessageDetails,
			wantErr: false,
		},
		{
			name: "should return error get nil gRPCErr",
			args: args{
				gRPCErr: nil,
				localizedMessage: &LocalizedMessage{
					Locale:  locale,
					Message: msg,
				},
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "should return error get gRPCErr which does not have a GRPCStatus() method",
			args: args{
				gRPCErr: fmt.Errorf("dummy-error"),
				localizedMessage: &LocalizedMessage{
					Locale:  locale,
					Message: msg,
				},
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
			got, err := AddLocalizedMessage(tt.args.gRPCErr, tt.args.localizedMessage)

			// Then
			assert(err).IsWantedError(tt.wantErr)
			assert(got).Equals(tt.want)
		})
	}
}

func TestLocalizedMessageFrom(t *testing.T) {
	locale := "dummy-locale"
	msg := "dummy-message"
	localizedMessageDetails := &errdetails.LocalizedMessage{
		Locale:  locale,
		Message: msg,
	}

	status := status.New(codes.Unimplemented, defaultUnimplementedErrMsg)
	gRPCErrWithoutLocalizedMessage := status.Err()

	zeroLocalizedMsg := LocalizedMessage{}

	statusWithLocalizedMsg, err := status.WithDetails(localizedMessageDetails)
	if err != nil {
		t.Fatal(err)
	}
	gRPCErrWithLocalizedMessage := statusWithLocalizedMsg.Err()

	localizedMsg := LocalizedMessage{
		Locale:  locale,
		Message: msg,
	}

	type args struct {
		gRPCErr error
	}
	tests := []struct {
		name string
		args args
		want LocalizedMessage
	}{
		{
			name: "Should return LocalizedMessage when get gRPCErr with localizedMsgDetails",
			args: args{
				gRPCErr: gRPCErrWithLocalizedMessage,
			},
			want: localizedMsg,
		},
		{
			name: "Should return zeroLocalizedMsg when get gRPCErr without localizedMsgDetails",
			args: args{
				gRPCErr: gRPCErrWithoutLocalizedMessage,
			},
			want: zeroLocalizedMsg,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Given
			assert := assert.New(t)

			// When
			got := LocalizedMessageFrom(tt.args.gRPCErr)

			// Then
			assert(got).Equals(tt.want)
		})
	}
}

func TestNewInvalidArgument(t *testing.T) {
	gRPCErrWithDefaultMsg := status.New(codes.InvalidArgument, defaultInvalidArgumentErrMsg).Err()
	gRPCErrWithoutDetails := status.New(codes.InvalidArgument, "dummy-message").Err()

	violations := []FieldViolation{
		{Field: "dummy-field-1", Description: "dummy-description-1"},
		{Field: "dummy-field-2", Description: "dummy-description-2"},
	}

	badRequestDetails := &errdetails.BadRequest{}
	errDetailsViolations := []*errdetails.BadRequest_FieldViolation{
		{
			Field:       "dummy-field-1",
			Description: "dummy-description-1",
		},
		{
			Field:       "dummy-field-2",
			Description: "dummy-description-2",
		},
	}

	badRequestDetails.FieldViolations = errDetailsViolations
	statusWithDetails := status.New(codes.InvalidArgument, "dummy-message")
	statusWithDetails, err := statusWithDetails.WithDetails(badRequestDetails)
	if err != nil {
		panic(err)
	}
	gRPCErrWithDetails := statusWithDetails.Err()

	type args struct {
		errMsg          string
		fieldViolations []FieldViolation
	}
	tests := []struct {
		name    string
		args    args
		want    error
		wantErr bool
	}{
		{
			name: "should return gRPC error with default errMsg when get empty errMsg",
			args: args{
				errMsg:          "",
				fieldViolations: nil,
			},
			want:    gRPCErrWithDefaultMsg,
			wantErr: false,
		},
		{
			name: "should return gRPC error without details when fieldViolations is nil",
			args: args{
				errMsg:          "dummy-message",
				fieldViolations: nil,
			},
			want:    gRPCErrWithoutDetails,
			wantErr: false,
		},
		{
			name: "should return gRPC error with details when get fieldViolations",
			args: args{
				errMsg:          "dummy-message",
				fieldViolations: violations,
			},
			want:    gRPCErrWithDetails,
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

func TestFieldViolationsFrom(t *testing.T) {
	field1 := "dummy-field-1"
	description1 := "dummy-description-1"
	field2 := "dummy-field-2"
	description2 := "dummy-description-2"
	fieldViolation1 := &errdetails.BadRequest_FieldViolation{
		Field:       field1,
		Description: description1,
	}
	fieldViolation2 := &errdetails.BadRequest_FieldViolation{
		Field:       field2,
		Description: description2,
	}
	badRequestDetails := &errdetails.BadRequest{FieldViolations: []*errdetails.BadRequest_FieldViolation{fieldViolation1, fieldViolation2}}

	status := status.New(codes.Unimplemented, defaultUnimplementedErrMsg)
	gRPCErrWithoutBadRequestDetails := status.Err()

	zeroFieldViolations := []FieldViolation{}

	statusWithBadRequestDetails, err := status.WithDetails(badRequestDetails)
	if err != nil {
		t.Fatal(err)
	}
	gRPCErrWithBadRequestDetails := statusWithBadRequestDetails.Err()

	fieldViolations := []FieldViolation{
		{Field: field1, Description: description1},
		{Field: field2, Description: description2},
	}

	type args struct {
		gRPCErr error
	}
	tests := []struct {
		name string
		args args
		want []FieldViolation
	}{
		{
			name: "Should return []FieldViolations when get gRPCErr with badRequestDetails",
			args: args{
				gRPCErr: gRPCErrWithBadRequestDetails,
			},
			want: fieldViolations,
		},
		{
			name: "Should return zeroFieldViolations when get gRPCErr without badRequestDetails",
			args: args{
				gRPCErr: gRPCErrWithoutBadRequestDetails,
			},
			want: zeroFieldViolations,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Given
			assert := assert.New(t)

			// When
			got := FieldViolationsFrom(tt.args.gRPCErr)

			// Then
			assert(got).Equals(tt.want)
		})
	}
}

func Test_newStatusWithErrorInfo(t *testing.T) {
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
	aborted := status.New(codes.Aborted, defaultAbortedErrMsg)
	abortedWithZeroErrorInfo, err := aborted.WithDetails(&errdetails.ErrorInfo{})
	if err != nil {
		t.Fatal(err)
	}

	type args struct {
		code      codes.Code
		errMsg    string
		errorInfo *ErrorInfo
	}
	tests := []struct {
		name    string
		args    args
		want    *status.Status
		wantErr bool
	}{
		{
			name: "should return Unauthenticated gRPC error when get codes.Unauthenticated",
			args: args{
				code:      codes.Unauthenticated,
				errMsg:    "",
				errorInfo: &ErrorInfo{},
			},
			want:    unauthenticatedWithZeroErrorInfo,
			wantErr: false,
		},
		{
			name: "should return PermissionDenied gRPC error when get codes.PermissionDenied",
			args: args{
				code:      codes.PermissionDenied,
				errMsg:    "",
				errorInfo: &ErrorInfo{},
			},
			want:    permissionDeniedWithZeroErrorInfo,
			wantErr: false,
		},
		{
			name: "should return Aborted gRPC error when get codes.Aborted",
			args: args{
				code:      codes.Aborted,
				errMsg:    "",
				errorInfo: &ErrorInfo{},
			},
			want:    abortedWithZeroErrorInfo,
			wantErr: false,
		},
		{
			name: "should return unmodified Aborted when get codes.Aborted and nil errorInfo",
			args: args{
				code:      codes.Aborted,
				errMsg:    "",
				errorInfo: nil,
			},
			want:    aborted,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Given
			assert := assert.New(t)

			// When
			got, err := newStatusWithErrorInfo(tt.args.code, tt.args.errMsg, tt.args.errorInfo)

			// Then
			assert(err).IsWantedError(tt.wantErr)
			assert(got.Proto()).Equals(tt.want.Proto())
		})
	}
}

func Test_newStatusWithResourceInfo(t *testing.T) {
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
		resourceInfo *ResourceInfo
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
				resourceInfo: &ResourceInfo{},
			},
			want:    notFoundWithZeroResourceInfo,
			wantErr: false,
		},
		{
			name: "should return AlreadyExists gRPC error when get code codes.AlreadyExists",
			args: args{
				code:         codes.AlreadyExists,
				errMsg:       "",
				resourceInfo: &ResourceInfo{},
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
			got, err := newStatusWithResourceInfo(tt.args.code, tt.args.errMsg, tt.args.resourceInfo)

			// Then
			assert(err).IsWantedError(tt.wantErr)
			assert(got.Proto()).Equals(tt.want.Proto())
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
		debugInfo *DebugInfo
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
				debugInfo: &DebugInfo{},
			},
			want:    dataLossWithZeroDebugInfo,
			wantErr: false,
		},
		{
			name: "should return Unknown gRPC error when get code codes.Unknown",
			args: args{
				code:      codes.Unknown,
				errMsg:    "",
				debugInfo: &DebugInfo{},
			},
			want:    unknownWithZeroDebugInfo,
			wantErr: false,
		},
		{
			name: "should return Internal gRPC error when get code codes.Internal",
			args: args{
				code:      codes.Internal,
				errMsg:    "",
				debugInfo: &DebugInfo{},
			},
			want:    internalWithZeroDebugInfo,
			wantErr: false,
		},
		{
			name: "should return Unavailable gRPC error when get code codes.Unavailable",
			args: args{
				code:      codes.Unavailable,
				errMsg:    "",
				debugInfo: &DebugInfo{},
			},
			want:    unavailableWithZeroDebugInfo,
			wantErr: false,
		},
		{
			name: "should return DeadlineExceeded gRPC error when get code codes.DeadlineExceeded",
			args: args{
				code:      codes.DeadlineExceeded,
				errMsg:    "",
				debugInfo: &DebugInfo{},
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
			got, err := newStatusWithDebugInfo(tt.args.code, tt.args.errMsg, tt.args.debugInfo)

			// Then
			assert(err).IsWantedError(tt.wantErr)
			assert(got.Proto()).Equals(tt.want.Proto())
		})
	}
}

func Test_newStatusWithQuotaFailure(t *testing.T) {
	resourceExhaustedWithDefaultErrMsg := status.New(codes.ResourceExhausted, defaultResourceExhaustedErrMsg)
	violations := []QuotaViolation{
		{Subject: "dummy-subject", Description: "dummy-description"},
	}
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
		violations []QuotaViolation
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
				violations: violations,
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
			got, err := newStatusWithQuotaFailure(tt.args.code, tt.args.errMsg, tt.args.violations)

			// Then
			assert(err).IsWantedError(tt.wantErr)
			assert(got.Proto()).Equals(tt.want.Proto())
		})
	}
}

func TestNewOutOfRange(t *testing.T) {
	gRPCErrWithDefaultErrMsg := status.New(codes.OutOfRange, defaultOutOfRangeErrMsg).Err()

	violations := []FieldViolation{
		{
			Field:       "dummy-field-1",
			Description: "dummy-description-1",
		},
		{
			Field:       "dummy-field-2",
			Description: "dummy-description-2",
		},
	}

	badRequestDetails := &errdetails.BadRequest{}
	errDetailsViolations := []*errdetails.BadRequest_FieldViolation{
		{
			Field:       "dummy-field-1",
			Description: "dummy-description-1",
		},
		{
			Field:       "dummy-field-2",
			Description: "dummy-description-2",
		},
	}

	badRequestDetails.FieldViolations = errDetailsViolations
	statusWithDetails := status.New(codes.OutOfRange, "dummy-message")
	statusWithDetails, err := statusWithDetails.WithDetails(badRequestDetails)
	if err != nil {
		panic(err)
	}
	gRPCErrWithDetails := statusWithDetails.Err()

	type args struct {
		errMsg          string
		fieldViolations []FieldViolation
	}
	tests := []struct {
		name    string
		args    args
		want    error
		wantErr bool
	}{
		{
			name: "Should return OutOfRange for valid arguments",
			args: args{
				errMsg:          "dummy-message",
				fieldViolations: violations,
			},
			want:    gRPCErrWithDetails,
			wantErr: false,
		},
		{
			name: "Should return OutOfRange with default errMsg when get empty errMsg",
			args: args{
				errMsg:          "",
				fieldViolations: nil,
			},
			want:    gRPCErrWithDefaultErrMsg,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Given
			assert := assert.New(t)

			// When
			got, err := NewOutOfRange(tt.args.errMsg, tt.args.fieldViolations)

			// Then
			assert(err).IsWantedError(tt.wantErr)
			assert(got).Equals(tt.want)
		})
	}
}

func TestPreconditionFailuresFrom(t *testing.T) {
	type1 := "dummy-type-1"
	subject1 := "dummy-subject-1"
	description1 := "dummy-description-1"
	type2 := "dummy-type-2"
	subject2 := "dummy-subject-2"
	description2 := "dummy-description-2"
	violation1 := &errdetails.PreconditionFailure_Violation{
		Type:        type1,
		Subject:     subject1,
		Description: description1,
	}
	violation2 := &errdetails.PreconditionFailure_Violation{
		Type:        type2,
		Subject:     subject2,
		Description: description2,
	}
	precondFailureDetails := &errdetails.PreconditionFailure{Violations: []*errdetails.PreconditionFailure_Violation{violation1, violation2}}

	status := status.New(codes.Unimplemented, defaultUnimplementedErrMsg)
	gRPCErrWithoutPrecondFailureDetails := status.Err()

	statusWithPrecondFailureDetails, err := status.WithDetails(precondFailureDetails)
	if err != nil {
		t.Fatal(err)
	}
	gRPCErrWithPrecondFailureDetails := statusWithPrecondFailureDetails.Err()

	precondFailures := []PreconditionFailure{
		{Type: type1, Subject: subject1, Description: description1},
		{Type: type2, Subject: subject2, Description: description2},
	}

	zeroPrecondFailures := []PreconditionFailure{}

	type args struct {
		gRPCErr error
	}
	tests := []struct {
		name string
		args args
		want []PreconditionFailure
	}{
		{
			name: "Should return []PreconditionFailure when get gRPCErr with precondFailureDetails",
			args: args{
				gRPCErr: gRPCErrWithPrecondFailureDetails,
			},
			want: precondFailures,
		},
		{
			name: "Should return zeroPrecondFailures when get gRPCErr without precondFailureDetails",
			args: args{
				gRPCErr: gRPCErrWithoutPrecondFailureDetails,
			},
			want: zeroPrecondFailures,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Given
			assert := assert.New(t)

			// When
			got := PreconditionFailuresFrom(tt.args.gRPCErr)

			// Then
			assert(got).Equals(tt.want)
		})
	}
}

func TestErrorInfoFrom(t *testing.T) {
	reason := "dummy-reason"
	domain := "dummy-domain"
	metadata := map[string]string{
		"dummy-key": "dummy-value",
	}
	errorInfoDetails := &errdetails.ErrorInfo{
		Reason:   reason,
		Domain:   domain,
		Metadata: metadata,
	}

	status := status.New(codes.Unimplemented, defaultUnimplementedErrMsg)
	gRPCErrWithoutErrorInfoDetails := status.Err()

	statusWithErrorInfoDetails, err := status.WithDetails(errorInfoDetails)
	if err != nil {
		t.Fatal(err)
	}
	gRPCErrWithErrorInfoDetails := statusWithErrorInfoDetails.Err()

	errorInfo := ErrorInfo{
		Reason:   reason,
		Domain:   domain,
		Metadata: metadata,
	}

	zeroErrorInfo := ErrorInfo{}

	type args struct {
		gRPCErr error
	}
	tests := []struct {
		name string
		args args
		want ErrorInfo
	}{
		{
			name: "Should return ErrorInfo when get gRPCErr with errorInfoDetails",
			args: args{
				gRPCErr: gRPCErrWithErrorInfoDetails,
			},
			want: errorInfo,
		},
		{
			name: "Should return zeroErrorInfo when get gRPCErr without errorInfoDetails",
			args: args{
				gRPCErr: gRPCErrWithoutErrorInfoDetails,
			},
			want: zeroErrorInfo,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Given
			assert := assert.New(t)

			// When
			got := ErrorInfoFrom(tt.args.gRPCErr)

			// Then
			assert(got).Equals(tt.want)
		})
	}
}

func TestResourceInfoFrom(t *testing.T) {
	resourceType := "dummy-resource-type"
	resourceName := "dummy-resource-name"
	owner := "dummy-owner"
	description := "dummy-description"
	resourceInfoDetails := &errdetails.ResourceInfo{
		ResourceType: resourceType,
		ResourceName: resourceName,
		Owner:        owner,
		Description:  description,
	}

	status := status.New(codes.Unimplemented, defaultUnimplementedErrMsg)
	gRPCErrWithoutResourceInfo := status.Err()

	statusWithResourceInfo, err := status.WithDetails(resourceInfoDetails)
	if err != nil {
		t.Fatal(err)
	}
	gRPCErrWithResourceInfo := statusWithResourceInfo.Err()

	resourceInfo := ResourceInfo{
		ResourceType: resourceType,
		ResourceName: resourceName,
		Owner:        owner,
		Description:  description,
	}

	zeroResourceInfo := ResourceInfo{}

	type args struct {
		gRPCErr error
	}
	tests := []struct {
		name string
		args args
		want ResourceInfo
	}{
		{
			name: "Should return ResourceInfo when get gRPCErr with resourceInfoDetails",
			args: args{
				gRPCErr: gRPCErrWithResourceInfo,
			},
			want: resourceInfo,
		},
		{
			name: "Should return zeroResourceInfo when get gRPCErr without resourceInfoDetails",
			args: args{
				gRPCErr: gRPCErrWithoutResourceInfo,
			},
			want: zeroResourceInfo,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Given
			assert := assert.New(t)

			// When
			got := ResourceInfoFrom(tt.args.gRPCErr)

			// Then
			assert(got).Equals(tt.want)
		})
	}
}

func TestQuotaViolationsFrom(t *testing.T) {
	subject1 := "dummy-subject-1"
	description1 := "dummy-description-1"
	subject2 := "dummy-subject-2"
	description2 := "dummy-description-2"

	quotaViolation1 := errdetails.QuotaFailure_Violation{
		Subject:     subject1,
		Description: description1,
	}
	quotaViolation2 := errdetails.QuotaFailure_Violation{
		Subject:     subject2,
		Description: description2,
	}
	quotaViolationDetails := &errdetails.QuotaFailure{Violations: []*errdetails.QuotaFailure_Violation{&quotaViolation1, &quotaViolation2}}

	status := status.New(codes.Unimplemented, defaultUnimplementedErrMsg)
	gRPCErrWithoutQuotaViolations := status.Err()

	statusWithQuotaViolations, err := status.WithDetails(quotaViolationDetails)
	if err != nil {
		t.Fatal(err)
	}
	gRPCErrWithQuotaViolations := statusWithQuotaViolations.Err()

	quotaViolations := []QuotaViolation{
		{Subject: subject1, Description: description1},
		{Subject: subject2, Description: description2},
	}

	zeroQuotaViolations := []QuotaViolation{}

	type args struct {
		gRPCErr error
	}
	tests := []struct {
		name string
		args args
		want []QuotaViolation
	}{
		{
			name: "Should return QuotaViolations when get gRPCErr with quotaViolationDetails",
			args: args{
				gRPCErr: gRPCErrWithQuotaViolations,
			},
			want: quotaViolations,
		},
		{
			name: "Should return zeroQuotaViolations when get gRPCErr without quotaViolationDetails",
			args: args{
				gRPCErr: gRPCErrWithoutQuotaViolations,
			},
			want: zeroQuotaViolations,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Given
			assert := assert.New(t)

			// When
			got := QuotaViolationsFrom(tt.args.gRPCErr)

			// Then
			assert(got).Equals(tt.want)
		})
	}
}

func TestCode(t *testing.T) {
	abortedStatus := status.New(codes.Aborted, defaultAbortedErrMsg)
	abortedGRPCErr := abortedStatus.Err()

	type args struct {
		gRPCErr error
	}
	tests := []struct {
		name string
		args args
		want codes.Code
	}{
		{
			name: "",
			args: args{
				gRPCErr: abortedGRPCErr,
			},
			want: codes.Aborted,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Given
			assert := assert.New(t)

			// When
			got := Code(tt.args.gRPCErr)

			// Then
			assert(got).Equals(tt.want)
		})
	}
}

func TestMessage(t *testing.T) {
	abortedStatusWithDefaultErrMsg := status.New(codes.Aborted, defaultAbortedErrMsg)
	abortedGRPCErrWithDefaultErrMsg := abortedStatusWithDefaultErrMsg.Err()
	abortedStatusWithCustomErrMsg := status.New(codes.Aborted, "dummy-custom-message")
	abortedGRPCErrWithCustomErrMsg := abortedStatusWithCustomErrMsg.Err()

	type args struct {
		gRPCErr error
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "Should return default message when get gRPCErr with default message",
			args: args{
				gRPCErr: abortedGRPCErrWithDefaultErrMsg,
			},
			want: defaultAbortedErrMsg,
		},
		{
			name: "Should return custom message when get gRPCErr with custom message",
			args: args{
				gRPCErr: abortedGRPCErrWithCustomErrMsg,
			},
			want: "dummy-custom-message",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Given
			assert := assert.New(t)

			// When
			got := Message(tt.args.gRPCErr)

			// Then
			assert(got).Equals(tt.want)
		})
	}
}
