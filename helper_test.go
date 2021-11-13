package errorcodeutil_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"strconv"
	"testing"

	"github.com/savannahghi/errorcodeutil"
	"github.com/stretchr/testify/assert"
)

// SETUP
type FakeResponse struct {
	t       *testing.T
	headers http.Header
	body    []byte
	status  int
}

func NewFakeRecorder(t *testing.T) *FakeResponse {
	return &FakeResponse{
		t:       t,
		headers: make(http.Header),
	}
}

func (r *FakeResponse) Header() http.Header {
	return r.headers
}

func (r *FakeResponse) Write(body []byte) (int, error) {
	r.body = body
	/*
	  For testing purposes only
	  this will activate the part of the code that checks for errors
	  on http.ResponseWriter
	*/
	if body == nil {
		return 0, fmt.Errorf("body is nil")
	}
	return len(body), nil
}

func (r *FakeResponse) WriteHeader(status int) {
	r.status = status
}

func (r *FakeResponse) Assert(status int, body string) {
	if r.status != status {
		r.t.Errorf("expected status %+v to equal %+v", r.status, status)
	}
	if string(r.body) != body {
		r.t.Errorf("expected body %+v to equal %+v", string(r.body), body)
	}
}

func TestReportErr(t *testing.T) {
	originalDebug := os.Getenv("DEBUG")

	type args struct {
		w      http.ResponseWriter
		err    error
		status int
		debug  bool
	}

	tests := []struct {
		name string
		args args
	}{
		{
			name: "happy case :) log the error with debug on",
			args: args{
				w:      httptest.NewRecorder(),
				err:    fmt.Errorf("a test error"),
				status: http.StatusBadRequest,
				debug:  true,
			},
		},
		{
			name: "happy case :) setting the status from the error reporter",
			args: args{
				w:      httptest.NewRecorder(),
				err:    fmt.Errorf("a test error"),
				status: http.StatusBadRequest,
				debug:  false,
			},
		},
	}
	for _, tt := range tests {
		// set debug environment
		os.Setenv("DEBUG", strconv.FormatBool(tt.args.debug))

		errorcodeutil.ReportErr(tt.args.w, tt.args.err, tt.args.status)
		rw, ok := tt.args.w.(*httptest.ResponseRecorder)
		assert.True(t, ok)
		assert.NotNil(t, rw)
		assert.Equal(t, tt.args.status, rw.Code)
	}
	// Reset debug environment
	os.Setenv("DEBUG", originalDebug)
}

func TestErrorMap(t *testing.T) {
	err := fmt.Errorf("test error")
	errMap := errorcodeutil.ErrorMap(err)
	if errMap["error"] == "" {
		t.Errorf("empty error key in error map")
	}
	if errMap["error"] != "test error" {
		t.Errorf("expected http error value to be %s, got %s", "test error", errMap["error"])
	}
}

func TestRespondWithError(t *testing.T) {
	validRecorder := httptest.NewRecorder()
	type args struct {
		w    http.ResponseWriter
		code int
		err  error
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "responds with correct status code",
			args: args{
				w:    validRecorder,
				code: http.StatusNotFound,
				err:  fmt.Errorf("not found error"),
			},
		},
	}
	for _, tt := range tests {
		errorcodeutil.RespondWithError(tt.args.w, tt.args.code, tt.args.err)
		rw, ok := tt.args.w.(*httptest.ResponseRecorder)
		assert.True(t, ok)
		assert.NotNil(t, rw)
		assert.Equal(t, tt.args.code, rw.Code)
	}
}

func TestRespondWithJSON(t *testing.T) {
	validRecorder := NewFakeRecorder(t)
	w := validRecorder
	code := 400
	errorcodeutil.RespondWithJSON(w, code, nil)
}
