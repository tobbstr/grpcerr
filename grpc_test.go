package grpcerr

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/tobbstr/testa/assert"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type responseWriterMock struct {
	bytes.Buffer
}

func (m *responseWriterMock) WriteHeader(statusCode int) {

}

func (m *responseWriterMock) Header() http.Header {
	return http.Header{}
}

func TestHttpResponseFormatterAsJSON(t *testing.T) {
	errorInfo := &errdetails.ErrorInfo{
		Reason: "dummy-reason",
		Domain: "dummy-domain",
		Metadata: map[string]string{
			"dummy-key": "dummy-value",
		},
	}
	resourceInfo := &errdetails.ResourceInfo{
		ResourceType: "dummy-resource-type",
		ResourceName: "dummy-resource-name",
		Owner:        "dummy-owner",
		Description:  "dummy-description",
	}
	debugInfo := &errdetails.DebugInfo{
		StackEntries: []string{"dummy-stack-entry"},
		Detail:       "dummy-detail",
	}
	invalidArgument, err := NewInvalidArgument("dummy-msg", []*errdetails.BadRequest_FieldViolation{{Field: "dummy-field-violation-field", Description: "dummy-field-violation-desc"}})
	if err != nil {
		t.Fatal(err)
	}
	failedPrecondition, err := NewFailedPrecondition("dummy-msg", []*errdetails.PreconditionFailure_Violation{{Type: "dummy-failed-precondition-violation-type", Subject: "dummy-failed-precondition-violation-subject", Description: "dummy-failed-precondition-violation-desc"}})
	if err != nil {
		t.Fatal(err)
	}
	outOfRange, err := NewOutOfRange("dummy-msg", []*errdetails.BadRequest_FieldViolation{{Field: "dummy-field-violation-field", Description: "dummy-field-violation-desc"}})
	if err != nil {
		t.Fatal(err)
	}
	unathenticated, err := NewUnauthenticated("dummy-msg", errorInfo)
	if err != nil {
		t.Fatal(err)
	}
	permissionDenied, err := NewPermissionDenied("dummy-msg", errorInfo)
	if err != nil {
		t.Fatal(err)
	}
	notFound, err := NewNotFound("dummy-msg", resourceInfo)
	if err != nil {
		t.Fatal(err)
	}
	aborted, err := NewAborted("dummy-msg", errorInfo)
	if err != nil {
		t.Fatal(err)
	}
	alreadyExists, err := NewAlreadyExists("dummy-msg", resourceInfo)
	if err != nil {
		t.Fatal(err)
	}
	resourceExhausted, err := NewResourceExhausted("dummy-msg", []*errdetails.QuotaFailure_Violation{
		{Subject: "dummy-subject", Description: "dummy-description"},
	})
	if err != nil {
		t.Fatal(err)
	}
	cancelled, err := NewCancelled("dummy-msg")
	if err != nil {
		t.Fatal(err)
	}
	dataLoss, err := NewDataLoss("dummy-msg", debugInfo)
	if err != nil {
		t.Fatal(err)
	}
	unknown, err := NewUnknown("dummy-msg", debugInfo)
	if err != nil {
		t.Fatal(err)
	}
	internal, err := NewInternal("dummy-msg", debugInfo)
	if err != nil {
		t.Fatal(err)
	}
	unimplemented, err := NewUnimplemented("dummy-msg")
	if err != nil {
		t.Fatal(err)
	}
	unavailable, err := NewUnavailable("dummy-msg", debugInfo)
	if err != nil {
		t.Fatal(err)
	}
	deadlineExceeded, err := NewDeadlineExceeded("dummy-msg", debugInfo)
	if err != nil {
		t.Fatal(err)
	}

	type args struct {
		w    *httptest.ResponseRecorder
		opts []ResponseWriterOption
		st   *status.Status
	}
	testCases := []struct {
		name               string
		args               args
		gotResponseWriter  *httptest.ResponseRecorder
		gotGRPCErr         *status.Status
		wantErr            error
		wantHttpStatusCode int
		wantHttpBody       string
	}{
		{
			name: "Should err when get nil gRPC error",
			args: args{
				w:    httptest.NewRecorder(),
				opts: nil,
				st:   nil,
			},
			wantErr:            fmt.Errorf("invalid argument: status was nil"),
			wantHttpStatusCode: http.StatusBadRequest,
			wantHttpBody:       ``,
		},
		{
			name: "Should write correct HTTP response when get InvalidArgument gRPC error",
			args: args{
				w:    httptest.NewRecorder(),
				opts: nil,
				st:   invalidArgument,
			},
			wantErr:            nil,
			wantHttpStatusCode: http.StatusBadRequest,
			wantHttpBody:       `{"code":3, "message":"dummy-msg", "details":[{"@type":"type.googleapis.com/google.rpc.BadRequest", "fieldViolations":[{"field":"dummy-field-violation-field", "description":"dummy-field-violation-desc"}]}]}`,
		},
		{
			name: "Should write correct HTTP response when get FailedPrecondition gRPC error",
			args: args{
				w:    httptest.NewRecorder(),
				opts: nil,
				st:   failedPrecondition,
			},
			wantErr:            nil,
			wantHttpStatusCode: http.StatusBadRequest,
			wantHttpBody:       "{\"code\":9, \"message\":\"dummy-msg\", \"details\":[{\"@type\":\"type.googleapis.com/google.rpc.PreconditionFailure\", \"violations\":[{\"type\":\"dummy-failed-precondition-violation-type\", \"subject\":\"dummy-failed-precondition-violation-subject\", \"description\":\"dummy-failed-precondition-violation-desc\"}]}]}",
		},
		{
			name: "Should write correct HTTP response when get OutOfRange gRPC error",
			args: args{
				w:    httptest.NewRecorder(),
				opts: nil,
				st:   outOfRange,
			},
			wantErr:            nil,
			wantHttpStatusCode: http.StatusBadRequest,
			wantHttpBody:       "{\"code\":11, \"message\":\"dummy-msg\", \"details\":[{\"@type\":\"type.googleapis.com/google.rpc.BadRequest\", \"fieldViolations\":[{\"field\":\"dummy-field-violation-field\", \"description\":\"dummy-field-violation-desc\"}]}]}",
		},
		{
			name: "Should write correct HTTP response when get Unauthenticated gRPC error",
			args: args{
				w:    httptest.NewRecorder(),
				opts: nil,
				st:   unathenticated,
			},
			wantErr:            nil,
			wantHttpStatusCode: http.StatusUnauthorized,
			wantHttpBody:       "{\"code\":16, \"message\":\"dummy-msg\", \"details\":[{\"@type\":\"type.googleapis.com/google.rpc.ErrorInfo\", \"reason\":\"dummy-reason\", \"domain\":\"dummy-domain\", \"metadata\":{\"dummy-key\":\"dummy-value\"}}]}",
		},
		{
			name: "Should write correct HTTP response when get PermissionDenied gRPC error",
			args: args{
				w:    httptest.NewRecorder(),
				opts: nil,
				st:   permissionDenied,
			},
			wantErr:            nil,
			wantHttpStatusCode: http.StatusForbidden,
			wantHttpBody:       "{\"code\":7, \"message\":\"dummy-msg\", \"details\":[{\"@type\":\"type.googleapis.com/google.rpc.ErrorInfo\", \"reason\":\"dummy-reason\", \"domain\":\"dummy-domain\", \"metadata\":{\"dummy-key\":\"dummy-value\"}}]}",
		},
		{
			name: "Should write correct HTTP response when get NotFound gRPC error",
			args: args{
				w:    httptest.NewRecorder(),
				opts: nil,
				st:   notFound,
			},
			wantErr:            nil,
			wantHttpStatusCode: http.StatusNotFound,
			wantHttpBody:       `{"code":5, "message":"dummy-msg", "details":[{"@type":"type.googleapis.com/google.rpc.ResourceInfo", "resourceType":"dummy-resource-type", "resourceName":"dummy-resource-name", "owner":"dummy-owner", "description":"dummy-description"}]}`,
		},
		{
			name: "Should write correct HTTP response when get Aborted gRPC error",
			args: args{
				w:    httptest.NewRecorder(),
				opts: nil,
				st:   aborted,
			},
			wantErr:            nil,
			wantHttpStatusCode: http.StatusConflict,
			wantHttpBody:       "{\"code\":10, \"message\":\"dummy-msg\", \"details\":[{\"@type\":\"type.googleapis.com/google.rpc.ErrorInfo\", \"reason\":\"dummy-reason\", \"domain\":\"dummy-domain\", \"metadata\":{\"dummy-key\":\"dummy-value\"}}]}",
		},
		{
			name: "Should write correct HTTP response when get AlreadyExists gRPC error",
			args: args{
				w:    httptest.NewRecorder(),
				opts: nil,
				st:   alreadyExists,
			},
			wantErr:            nil,
			wantHttpStatusCode: http.StatusConflict,
			wantHttpBody:       `{"code":6,"message":"dummy-msg","details":[{"@type":"type.googleapis.com/google.rpc.ResourceInfo","resourceType":"dummy-resource-type","resourceName":"dummy-resource-name","owner":"dummy-owner","description":"dummy-description"}]}`,
		},
		{
			name: "Should write correct HTTP response when get ResourceExhausted gRPC error",
			args: args{
				w:    httptest.NewRecorder(),
				opts: nil,
				st:   resourceExhausted,
			},
			wantErr:            nil,
			wantHttpStatusCode: http.StatusTooManyRequests,
			wantHttpBody:       `{"code":8,"message":"dummy-msg","details":[{"@type":"type.googleapis.com/google.rpc.QuotaFailure","violations":[{"subject":"dummy-subject","description":"dummy-description"}]}]}`,
		},
		{
			name: "Should write correct HTTP response when get Cancelled gRPC error",
			args: args{
				w:    httptest.NewRecorder(),
				opts: nil,
				st:   cancelled,
			},
			wantErr:            nil,
			wantHttpStatusCode: 499,
			wantHttpBody:       `{"code": 1, "message": "dummy-msg"}`,
		},
		{
			name: "Should write correct HTTP response when get DataLoss gRPC error",
			args: args{
				w:    httptest.NewRecorder(),
				opts: nil,
				st:   dataLoss,
			},
			wantErr:            nil,
			wantHttpStatusCode: http.StatusInternalServerError,
			wantHttpBody:       `{"code":15, "message":"dummy-msg", "details":[{"@type":"type.googleapis.com/google.rpc.DebugInfo", "stackEntries":["dummy-stack-entry"], "detail":"dummy-detail"}]}`,
		},
		{
			name: "Should write correct HTTP response when get Unknown gRPC error",
			args: args{
				w:    httptest.NewRecorder(),
				opts: nil,
				st:   unknown,
			},
			wantErr:            nil,
			wantHttpStatusCode: http.StatusInternalServerError,
			wantHttpBody:       `{"code":2, "message":"dummy-msg", "details":[{"@type":"type.googleapis.com/google.rpc.DebugInfo", "stackEntries":["dummy-stack-entry"], "detail":"dummy-detail"}]}`,
		},
		{
			name: "Should write correct HTTP response when get Internal gRPC error",
			args: args{
				w:    httptest.NewRecorder(),
				opts: nil,
				st:   internal,
			},
			wantErr:            nil,
			wantHttpStatusCode: http.StatusInternalServerError,
			wantHttpBody:       `{"code":13, "message":"dummy-msg", "details":[{"@type":"type.googleapis.com/google.rpc.DebugInfo", "stackEntries":["dummy-stack-entry"], "detail":"dummy-detail"}]}`,
		},
		{
			name: "Should write correct HTTP response when get Unimplemented gRPC error",
			args: args{
				w:    httptest.NewRecorder(),
				opts: nil,
				st:   unimplemented,
			},
			wantErr:            nil,
			wantHttpStatusCode: http.StatusNotImplemented,
			wantHttpBody:       `{"code":12, "message":"dummy-msg"}`,
		},
		{
			name: "Should write correct HTTP response when get Unavailable gRPC error",
			args: args{
				w:    httptest.NewRecorder(),
				opts: nil,
				st:   unavailable,
			},
			wantErr:            nil,
			wantHttpStatusCode: http.StatusServiceUnavailable,
			wantHttpBody:       `{"code":14, "message":"dummy-msg", "details":[{"@type":"type.googleapis.com/google.rpc.DebugInfo", "stackEntries":["dummy-stack-entry"], "detail":"dummy-detail"}]}`,
		},
		{
			name: "Should write correct HTTP response when get DeadlineExceeded gRPC error",
			args: args{
				w:    httptest.NewRecorder(),
				opts: nil,
				st:   deadlineExceeded,
			},
			wantErr:            nil,
			wantHttpStatusCode: http.StatusGatewayTimeout,
			wantHttpBody:       `{"code":4, "message":"dummy-msg", "details":[{"@type":"type.googleapis.com/google.rpc.DebugInfo", "stackEntries":["dummy-stack-entry"], "detail":"dummy-detail"}]}`,
		},
	}

	// t.Parallel()
	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			// Given
			require := assert.NewFatal(t)
			write := HttpResponseWriterFrom(tc.args.w)

			// When
			gotErr := write(tc.args.st).AsJSON()

			// Then
			require(gotErr).Equals(tc.wantErr)
			if gotErr != nil {
				return
			}

			got := tc.args.w.Result()

			require(got.StatusCode).Equals(tc.wantHttpStatusCode)

			require(got.Header.Get("Content-Type")).Equals("application/json")

			defer got.Body.Close()
			gotHttpBody, err := ioutil.ReadAll(got.Body)
			if err != nil {
				t.Fatalf("could not read result body: %v\n", err)
			}

			require(string(gotHttpBody)).IsJSONEqualTo(tc.wantHttpBody)
		})
	}
}

func TestHttpResponseFormatterWithOptionAsJSON(t *testing.T) {
	unimplemented, err := NewUnimplemented("dummy-msg")
	if err != nil {
		t.Fatal(err)
	}

	type args struct {
		w    *httptest.ResponseRecorder
		opts []ResponseWriterOption
		st   *status.Status
	}
	testCases := []struct {
		name               string
		args               args
		gotResponseWriter  *httptest.ResponseRecorder
		gotGRPCErr         *status.Status
		gotCustomOption    func(http.ResponseWriter)
		wantErr            error
		wantHttpStatusCode int
		wantHttpBody       string
		wantCustomHeader   string
	}{
		{
			name: "Should err when get nil gRPC error and custom option",
			args: args{
				w: httptest.NewRecorder(),
				opts: []ResponseWriterOption{
					func(w http.ResponseWriter) { w.WriteHeader(http.StatusOK) },
				},
				st: nil,
			},
			wantErr:            fmt.Errorf("invalid argument: status was nil"),
			wantHttpStatusCode: http.StatusOK,
			wantHttpBody:       ``,
		},
		{
			name: "Should write correct HTTP response for custom content-type",
			args: args{
				w: httptest.NewRecorder(),
				opts: []ResponseWriterOption{
					func(w http.ResponseWriter) { w.Header().Set("Content-Type", "dummy-content-type-value") },
				},
				st: unimplemented,
			},
			wantErr:            nil,
			wantHttpStatusCode: http.StatusNotImplemented,
			wantHttpBody:       `{"code":12, "message":"dummy-msg"}`,
			wantCustomHeader:   "dummy-content-type-value",
		},
	}

	// t.Parallel()
	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			// Given
			require := assert.NewFatal(t)
			write := HttpResponseWriterFrom(tc.args.w, tc.args.opts...)

			// When
			gotErr := write(tc.args.st).AsJSON()

			// Then
			require(gotErr).Equals(tc.wantErr)
			if gotErr != nil {
				return
			}

			got := tc.args.w.Result()

			require(got.StatusCode).Equals(tc.wantHttpStatusCode)

			require(got.Header.Get("Content-Type")).Equals(tc.wantCustomHeader)

			defer got.Body.Close()
			gotHttpBody, err := ioutil.ReadAll(got.Body)
			if err != nil {
				t.Fatalf("could not read result body: %v\n", err)
			}

			require(string(gotHttpBody)).IsJSONEqualTo(tc.wantHttpBody)
		})
	}
}

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
			want:    gRPCErrWithDebugInfo,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Given
			require := assert.NewFatal(t)
			assert := assert.New(t)

			// When
			got, err := AddDebugInfo(tt.args.gRPCErr, tt.args.debugInfo)

			// Then
			require((err != nil)).Equals(tt.wantErr)
			if err != nil {
				return
			}
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
			require := assert.NewFatal(t)
			assert := assert.New(t)

			// When
			got, err := AddRequestInfo(tt.args.gRPCErr, tt.args.requestID, tt.args.servingData)

			// Then
			require((err != nil)).Equals(tt.wantErr)
			if err != nil {
				return
			}
			assert(got).Equals(tt.want)
		})
	}
}

func Test_httpStatusCodeFrom(t *testing.T) {
	type args struct {
		st *status.Status
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		{
			name: "should return StatusUnprocessableEntity when get gRPC error with illegal code",
			args: args{
				st: status.New(codes.Code(9999), "dummy-msg"),
			},
			want: http.StatusUnprocessableEntity,
		},
		// The rest of the status codes are tested by TestHttpResponseFormatterAsJSON
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Given
			assert := assert.New(t)

			// When
			got := httpStatusCodeFrom(tt.args.st)

			// Then
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
			require := assert.NewFatal(t)

			// When
			got, err := AddHelp(tt.args.gRPCErr, tt.args.links)

			// Then
			require((err != nil)).Equals(tt.wantErr)
			if err != nil {
				return
			}
			require(got).Equals(tt.want)
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
			require := assert.NewFatal(t)

			// When
			got, err := AddLocalizedMessage(tt.args.gRPCErr, tt.args.locale, tt.args.msg)

			// Then
			require((err != nil)).Equals(tt.wantErr)
			if err != nil {
				return
			}
			require(got).Equals(tt.want)
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
			require := assert.NewFatal(t)

			// When
			got, err := NewInvalidArgument(tt.args.errMsg, tt.args.fieldViolations)

			// Then
			require((err != nil)).Equals(tt.wantErr)
			if err != nil {
				return
			}
			require(got).Equals(tt.want)
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
			require := assert.NewFatal(t)

			// When
			got, err := newGRPCErrorWithErrorInfo(tt.args.code, tt.args.errMsg, tt.args.errorInfo)

			// Then
			require((err != nil)).Equals(tt.wantErr)
			if err != nil {
				return
			}
			require(got).Equals(tt.want)
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
			require := assert.NewFatal(t)

			// When
			got, err := newGRPCErrorWithResourceInfo(tt.args.code, tt.args.errMsg, tt.args.resourceInfo)

			// Then
			require((err != nil)).Equals(tt.wantErr)
			if err != nil {
				return
			}
			require(got).Equals(tt.want)
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
			require := assert.NewFatal(t)

			// When
			got, err := newGRPCErrorWithDebugInfo(tt.args.code, tt.args.errMsg, tt.args.debugInfo)

			// Then
			require((err != nil)).Equals(tt.wantErr)
			if err != nil {
				return
			}
			require(got).Equals(tt.want)
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
			require := assert.NewFatal(t)

			// When
			got, err := newGRPCErrorWithQuotaFailure(tt.args.code, tt.args.errMsg, tt.args.violations)

			// Then
			require((err != nil)).Equals(tt.wantErr)
			if err != nil {
				return
			}
			require(got).Equals(tt.want)
		})
	}
}
