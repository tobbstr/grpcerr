package grpcerr

import (
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/tobbstr/testa/assert"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func TestHttpResponseEncodeWriteAsJSON(t *testing.T) {
	errorInfo := &ErrorInfo{
		Reason: "dummy-reason",
		Domain: "dummy-domain",
		Metadata: map[string]string{
			"dummy-key": "dummy-value",
		},
	}
	resourceInfo := &ResourceInfo{
		ResourceType: "dummy-resource-type",
		ResourceName: "dummy-resource-name",
		Owner:        "dummy-owner",
		Description:  "dummy-description",
	}
	debugInfo := &DebugInfo{
		StackEntries: []string{"dummy-stack-entry"},
		Detail:       "dummy-detail",
	}
	invalidArgument, err := NewInvalidArgument("dummy-msg", []FieldViolation{{Field: "dummy-field-violation-field", Description: "dummy-field-violation-desc"}})
	if err != nil {
		t.Fatal(err)
	}
	failedPrecondition, err := NewFailedPrecondition("dummy-msg", []PreconditionFailure{{Type: "dummy-failed-precondition-violation-type", Subject: "dummy-failed-precondition-violation-subject", Description: "dummy-failed-precondition-violation-desc"}})
	if err != nil {
		t.Fatal(err)
	}
	outOfRange, err := NewOutOfRange("dummy-msg", []FieldViolation{{Field: "dummy-field-violation-field", Description: "dummy-field-violation-desc"}})
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
	resourceExhausted, err := NewResourceExhausted("dummy-msg", []QuotaViolation{
		{Subject: "dummy-subject", Description: "dummy-description"},
	})
	if err != nil {
		t.Fatal(err)
	}
	cancelled := NewCancelled("dummy-msg")

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
	unimplemented := NewUnimplemented("dummy-msg")

	unavailable, err := NewUnavailable("dummy-msg", debugInfo)
	if err != nil {
		t.Fatal(err)
	}
	deadlineExceeded, err := NewDeadlineExceeded("dummy-msg", debugInfo)
	if err != nil {
		t.Fatal(err)
	}

	type args struct {
		w       *httptest.ResponseRecorder
		opts    []ResponseWriterOption
		gRPCErr error
	}
	tests := []struct {
		name    string
		args    args
		want    *http.Response
		wantErr error
	}{
		{
			name: "Should err when get nil gRPC error",
			args: args{
				w:       httptest.NewRecorder(),
				opts:    nil,
				gRPCErr: nil,
			},
			want: &http.Response{
				StatusCode: http.StatusInternalServerError,
				Body:       io.NopCloser(strings.NewReader("")),
				Header: map[string][]string{
					"Content-Type": {""},
				},
			},
			wantErr: fmt.Errorf("invalid argument: gRPCErr was nil"),
		},
		{
			name: "Should write correct HTTP response when get InvalidArgument gRPC error",
			args: args{
				w:       httptest.NewRecorder(),
				opts:    nil,
				gRPCErr: invalidArgument,
			},
			want: &http.Response{
				StatusCode: http.StatusBadRequest,
				Body:       io.NopCloser(strings.NewReader(`{"code":3, "message":"dummy-msg", "details":[{"@type":"type.googleapis.com/google.rpc.BadRequest", "fieldViolations":[{"field":"dummy-field-violation-field", "description":"dummy-field-violation-desc"}]}]}`)),
				Header: map[string][]string{
					"Content-Type": {"application/json"},
				},
			},
			wantErr: nil,
		},
		{
			name: "Should write correct HTTP response when get FailedPrecondition gRPC error",
			args: args{
				w:       httptest.NewRecorder(),
				opts:    nil,
				gRPCErr: failedPrecondition,
			},
			want: &http.Response{
				StatusCode: http.StatusBadRequest,
				Body:       io.NopCloser(strings.NewReader("{\"code\":9, \"message\":\"dummy-msg\", \"details\":[{\"@type\":\"type.googleapis.com/google.rpc.PreconditionFailure\", \"violations\":[{\"type\":\"dummy-failed-precondition-violation-type\", \"subject\":\"dummy-failed-precondition-violation-subject\", \"description\":\"dummy-failed-precondition-violation-desc\"}]}]}")),
				Header: map[string][]string{
					"Content-Type": {"application/json"},
				},
			},
			wantErr: nil,
		},
		{
			name: "Should write correct HTTP response when get OutOfRange gRPC error",
			args: args{
				w:       httptest.NewRecorder(),
				opts:    nil,
				gRPCErr: outOfRange,
			},
			want: &http.Response{
				StatusCode: http.StatusBadRequest,
				Body:       io.NopCloser(strings.NewReader("{\"code\":11, \"message\":\"dummy-msg\", \"details\":[{\"@type\":\"type.googleapis.com/google.rpc.BadRequest\", \"fieldViolations\":[{\"field\":\"dummy-field-violation-field\", \"description\":\"dummy-field-violation-desc\"}]}]}")),
				Header: map[string][]string{
					"Content-Type": {"application/json"},
				},
			},
			wantErr: nil,
		},
		{
			name: "Should write correct HTTP response when get Unauthenticated gRPC error",
			args: args{
				w:       httptest.NewRecorder(),
				opts:    nil,
				gRPCErr: unathenticated,
			},
			want: &http.Response{
				StatusCode: http.StatusUnauthorized,
				Body:       io.NopCloser(strings.NewReader("{\"code\":16, \"message\":\"dummy-msg\", \"details\":[{\"@type\":\"type.googleapis.com/google.rpc.ErrorInfo\", \"reason\":\"dummy-reason\", \"domain\":\"dummy-domain\", \"metadata\":{\"dummy-key\":\"dummy-value\"}}]}")),
				Header: map[string][]string{
					"Content-Type": {"application/json"},
				},
			},
			wantErr: nil,
		},
		{
			name: "Should write correct HTTP response when get PermissionDenied gRPC error",
			args: args{
				w:       httptest.NewRecorder(),
				opts:    nil,
				gRPCErr: permissionDenied,
			},
			want: &http.Response{
				StatusCode: http.StatusForbidden,
				Body:       io.NopCloser(strings.NewReader("{\"code\":7, \"message\":\"dummy-msg\", \"details\":[{\"@type\":\"type.googleapis.com/google.rpc.ErrorInfo\", \"reason\":\"dummy-reason\", \"domain\":\"dummy-domain\", \"metadata\":{\"dummy-key\":\"dummy-value\"}}]}")),
				Header: map[string][]string{
					"Content-Type": {"application/json"},
				},
			},
			wantErr: nil,
		},
		{
			name: "Should write correct HTTP response when get NotFound gRPC error",
			args: args{
				w:       httptest.NewRecorder(),
				opts:    nil,
				gRPCErr: notFound,
			},
			want: &http.Response{
				StatusCode: http.StatusNotFound,
				Body:       io.NopCloser(strings.NewReader(`{"code":5, "message":"dummy-msg", "details":[{"@type":"type.googleapis.com/google.rpc.ResourceInfo", "resourceType":"dummy-resource-type", "resourceName":"dummy-resource-name", "owner":"dummy-owner", "description":"dummy-description"}]}`)),
				Header: map[string][]string{
					"Content-Type": {"application/json"},
				},
			},
			wantErr: nil,
		},
		{
			name: "Should write correct HTTP response when get Aborted gRPC error",
			args: args{
				w:       httptest.NewRecorder(),
				opts:    nil,
				gRPCErr: aborted,
			},
			want: &http.Response{
				StatusCode: http.StatusConflict,
				Body:       io.NopCloser(strings.NewReader("{\"code\":10, \"message\":\"dummy-msg\", \"details\":[{\"@type\":\"type.googleapis.com/google.rpc.ErrorInfo\", \"reason\":\"dummy-reason\", \"domain\":\"dummy-domain\", \"metadata\":{\"dummy-key\":\"dummy-value\"}}]}")),
				Header: map[string][]string{
					"Content-Type": {"application/json"},
				},
			},
			wantErr: nil,
		},
		{
			name: "Should write correct HTTP response when get AlreadyExists gRPC error",
			args: args{
				w:       httptest.NewRecorder(),
				opts:    nil,
				gRPCErr: alreadyExists,
			},
			want: &http.Response{
				StatusCode: http.StatusConflict,
				Body:       io.NopCloser(strings.NewReader(`{"code":6, "message":"dummy-msg", "details":[{"@type":"type.googleapis.com/google.rpc.ResourceInfo", "resourceType":"dummy-resource-type", "resourceName":"dummy-resource-name", "owner":"dummy-owner", "description":"dummy-description"}]}`)),
				Header: map[string][]string{
					"Content-Type": {"application/json"},
				},
			},
			wantErr: nil,
		},
		{
			name: "Should write correct HTTP response when get ResourceExhausted gRPC error",
			args: args{
				w:       httptest.NewRecorder(),
				opts:    nil,
				gRPCErr: resourceExhausted,
			},
			want: &http.Response{
				StatusCode: http.StatusTooManyRequests,
				Body:       io.NopCloser(strings.NewReader(`{"code":8, "message":"dummy-msg", "details":[{"@type":"type.googleapis.com/google.rpc.QuotaFailure", "violations":[{"subject":"dummy-subject", "description":"dummy-description"}]}]}`)),
				Header: map[string][]string{
					"Content-Type": {"application/json"},
				},
			},
			wantErr: nil,
		},
		{
			name: "Should write correct HTTP response when get Cancelled gRPC error",
			args: args{
				w:       httptest.NewRecorder(),
				opts:    nil,
				gRPCErr: cancelled,
			},
			want: &http.Response{
				StatusCode: 499,
				Body:       io.NopCloser(strings.NewReader(`{"code":1, "message":"dummy-msg"}`)),
				Header: map[string][]string{
					"Content-Type": {"application/json"},
				},
			},
			wantErr: nil,
		},
		{
			name: "Should write correct HTTP response when get DataLoss gRPC error",
			args: args{
				w:       httptest.NewRecorder(),
				opts:    nil,
				gRPCErr: dataLoss,
			},
			want: &http.Response{
				StatusCode: http.StatusInternalServerError,
				Body:       io.NopCloser(strings.NewReader(`{"code":15, "message":"dummy-msg", "details":[{"@type":"type.googleapis.com/google.rpc.DebugInfo", "stackEntries":["dummy-stack-entry"], "detail":"dummy-detail"}]}`)),
				Header: map[string][]string{
					"Content-Type": {"application/json"},
				},
			},
			wantErr: nil,
		},
		{
			name: "Should write correct HTTP response when get Unknown gRPC error",
			args: args{
				w:       httptest.NewRecorder(),
				opts:    nil,
				gRPCErr: unknown,
			},
			want: &http.Response{
				StatusCode: http.StatusInternalServerError,
				Body:       io.NopCloser(strings.NewReader(`{"code":2, "message":"dummy-msg", "details":[{"@type":"type.googleapis.com/google.rpc.DebugInfo", "stackEntries":["dummy-stack-entry"], "detail":"dummy-detail"}]}`)),
				Header: map[string][]string{
					"Content-Type": {"application/json"},
				},
			},
			wantErr: nil,
		},
		{
			name: "Should write correct HTTP response when get Internal gRPC error",
			args: args{
				w:       httptest.NewRecorder(),
				opts:    nil,
				gRPCErr: internal,
			},
			want: &http.Response{
				StatusCode: http.StatusInternalServerError,
				Body:       io.NopCloser(strings.NewReader(`{"code":13, "message":"dummy-msg", "details":[{"@type":"type.googleapis.com/google.rpc.DebugInfo", "stackEntries":["dummy-stack-entry"], "detail":"dummy-detail"}]}`)),
				Header: map[string][]string{
					"Content-Type": {"application/json"},
				},
			},
			wantErr: nil,
		},
		{
			name: "Should write correct HTTP response when get Unimplemented gRPC error",
			args: args{
				w:       httptest.NewRecorder(),
				opts:    nil,
				gRPCErr: unimplemented,
			},
			want: &http.Response{
				StatusCode: http.StatusNotImplemented,
				Body:       io.NopCloser(strings.NewReader(`{"code":12, "message":"dummy-msg"}`)),
				Header: map[string][]string{
					"Content-Type": {"application/json"},
				},
			},
			wantErr: nil,
		},
		{
			name: "Should write correct HTTP response when get Unavailable gRPC error",
			args: args{
				w:       httptest.NewRecorder(),
				opts:    nil,
				gRPCErr: unavailable,
			},
			want: &http.Response{
				StatusCode: http.StatusServiceUnavailable,
				Body:       io.NopCloser(strings.NewReader(`{"code":14, "message":"dummy-msg", "details":[{"@type":"type.googleapis.com/google.rpc.DebugInfo", "stackEntries":["dummy-stack-entry"], "detail":"dummy-detail"}]}`)),
				Header: map[string][]string{
					"Content-Type": {"application/json"},
				},
			},
			wantErr: nil,
		},
		{
			name: "Should write correct HTTP response when get DeadlineExceeded gRPC error",
			args: args{
				w:       httptest.NewRecorder(),
				opts:    nil,
				gRPCErr: deadlineExceeded,
			},
			want: &http.Response{
				StatusCode: http.StatusGatewayTimeout,
				Body:       io.NopCloser(strings.NewReader(`{"code":4, "message":"dummy-msg", "details":[{"@type":"type.googleapis.com/google.rpc.DebugInfo", "stackEntries":["dummy-stack-entry"], "detail":"dummy-detail"}]}`)),
				Header: map[string][]string{
					"Content-Type": {"application/json"},
				},
			},
			wantErr: nil,
		},
		{
			name: "Should write correct HTTP response when setting custom option",
			args: args{
				w: httptest.NewRecorder(),
				opts: []ResponseWriterOption{
					func(w http.ResponseWriter) { w.Header().Set("Content-Type", "dummy-content-type-value") },
					func(w http.ResponseWriter) { w.WriteHeader(http.StatusOK) },
				},
				gRPCErr: unathenticated,
			},
			want: &http.Response{
				StatusCode: http.StatusOK,
				Body:       io.NopCloser(strings.NewReader(`{"code":16,"message":"dummy-msg","details":[{"@type":"type.googleapis.com/google.rpc.ErrorInfo","reason":"dummy-reason","domain":"dummy-domain","metadata":{"dummy-key":"dummy-value"}}]}`)),
				Header: map[string][]string{
					"Content-Type": {"dummy-content-type-value"},
				},
			},
			wantErr: nil,
		},
		{
			name: "Should return error when gRPCErr does not have GRPCStatus() method",
			args: args{
				w:       httptest.NewRecorder(),
				opts:    nil,
				gRPCErr: fmt.Errorf("dummy-error-that-does-not-have-grpcstatus-method"),
			},
			want: &http.Response{
				StatusCode: http.StatusInternalServerError,
				Body:       io.NopCloser(strings.NewReader("")),
				Header:     map[string][]string{},
			},
			wantErr: fmt.Errorf("invalid argument: gRPCErr's root error must have the GRPCStatus() method"),
		},
		{
			name: "Should write correct HTTP response when gRPCErr is wrapped error",
			args: args{
				w:       httptest.NewRecorder(),
				opts:    nil,
				gRPCErr: fmt.Errorf("dummy-wrapping-error: %w", unimplemented),
			},
			want: &http.Response{
				StatusCode: http.StatusNotImplemented,
				Body:       io.NopCloser(strings.NewReader(`{"code":12, "message":"dummy-msg"}`)),
				Header: map[string][]string{
					"Content-Type": {"application/json"},
				},
			},
			wantErr: nil,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			// Given
			assert := assert.New(t)
			encodeAndWrite := NewHttpResponseEncodeWriter(tt.args.w, tt.args.opts...)

			// When
			gotErr := encodeAndWrite(tt.args.gRPCErr).AsJSON()
			got := tt.args.w.Result()

			// Then
			assert(gotErr).Equals(tt.wantErr)
			assert(got.StatusCode).Equals(tt.want.StatusCode)
			assert(got.Header.Get("Content-Type")).Equals(tt.want.Header.Get("Content-Type"))

			defer got.Body.Close()
			gotHttpBody, err := ioutil.ReadAll(got.Body)
			assert(err).IsNil()

			defer tt.want.Body.Close()
			wantHttpBody, err := ioutil.ReadAll(tt.want.Body)
			assert(err).IsNil()

			assert(string(gotHttpBody)).IsJSONEqualTo(string(wantHttpBody))
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
			name: "should return InternalServerError when get gRPC error with illegal code",
			args: args{
				st: status.New(codes.Code(9999), "dummy-msg"),
			},
			want: http.StatusInternalServerError,
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
