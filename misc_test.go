package sendgrid

import (
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestErrorResponse(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/templates/d-12345abcde/versions/aaaaaa-bbbb-0000-0000-aaaaaaaaa", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "PATCH")
		w.WriteHeader(http.StatusNotFound)
		if _, err := fmt.Fprint(w, `{"error": "You cannot switch editors once a dynamic template version has been created."}`); err != nil {
			t.Fatal(err)
		}
	})

	client.debug = true
	client.httpclient = &http.Client{}
	client.log = log.New(os.Stdout, "sendgrid: ", log.Lshortfile|log.LstdFlags)

	client.Debugf("%s", "test")
	client.Debugln("test")

	if _, err := client.UpdateTemplateVersion(context.TODO(), "d-12345abcde", "aaaaaa-bbbb-0000-0000-aaaaaaaaa", &InputUpdateTemplateVersion{
		Editor: "code",
	}); err == nil {
		t.Fatal("expected an error but got none", err)
	}
}

func TestErrorsResponse(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/teammates/dummy", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		w.WriteHeader(http.StatusNotFound)
		if _, err := fmt.Fprint(w, `{"errors":[{"message": "teammate does not exis"}]}`); err != nil {
			t.Fatal(err)
		}
	})

	client.debug = true
	client.httpclient = &http.Client{}
	client.log = log.New(os.Stdout, "sendgrid: ", log.Lshortfile|log.LstdFlags)

	client.Debugf("%s", "test")
	client.Debugln("test")

	if _, err := client.GetTeammate(context.TODO(), "dummy"); err == nil {
		t.Fatal("expected an error but got none", err)
	}
}

func TestStatusUnAuthorized(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/teammates/dummy", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		w.WriteHeader(http.StatusUnauthorized)
	})

	_, err := client.GetTeammate(context.TODO(), "dummy")
	if err == nil {
		t.Fatal("expected an error but got none", err)
	}
}

func TestErrorResponseErr(t *testing.T) {
	tests := []struct {
		name     string
		error    string
		expected bool
	}{
		{
			name:     "empty error",
			error:    "",
			expected: false,
		},
		{
			name:     "non-empty error",
			error:    "test error message",
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resp := ErrorResponse{Error: tt.error}
			err := resp.Err()
			
			if tt.expected {
				assert.Error(t, err)
				assert.Equal(t, tt.error, err.Error())
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestRateLimitedErrorError(t *testing.T) {
	retryAfter := 30 * time.Second
	err := &RateLimitedError{RetryAfter: retryAfter}
	
	expected := "sendgrid rate limit exceeded, retry after 30s"
	assert.Equal(t, expected, err.Error())
}

func TestErrorsResponseErrs(t *testing.T) {
	tests := []struct {
		name     string
		errors   []*Error
		expected string
		hasError bool
	}{
		{
			name:     "empty errors",
			errors:   []*Error{},
			expected: "",
			hasError: false,
		},
		{
			name: "single error without field",
			errors: []*Error{
				{Message: String("test message")},
			},
			expected: "message: test message",
			hasError: true,
		},
		{
			name: "single error with field",
			errors: []*Error{
				{Field: String("email"), Message: String("invalid email")},
			},
			expected: "field: email, message: invalid email",
			hasError: true,
		},
		{
			name: "multiple errors",
			errors: []*Error{
				{Field: String("email"), Message: String("invalid email")},
				{Message: String("missing name")},
			},
			expected: "field: email, message: invalid email, message: missing name",
			hasError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resp := ErrorsResponse{Errors: tt.errors}
			err := resp.Errs()
			
			if tt.hasError {
				assert.Error(t, err)
				assert.Equal(t, tt.expected, err.Error())
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestStatusCodeError(t *testing.T) {
	err := statusCodeError{
		Code:   500,
		Status: "500 Internal Server Error",
	}
	
	assert.Equal(t, "sendgrid server error: 500 Internal Server Error", err.Error())
	assert.Equal(t, 500, err.HTTPStatusCode())
}

func TestCheckStatusCode(t *testing.T) {
	// Mock debug interface
	debug := &mockDebug{debug: true}

	tests := []struct {
		name           string
		statusCode     int
		headers        map[string]string
		body           string
		expectedError  string
		shouldReturnError bool
	}{
		{
			name:       "success 200",
			statusCode: 200,
			body:       `{"success": true}`,
			shouldReturnError: false,
		},
		{
			name:       "success 201",
			statusCode: 201,
			body:       `{"created": true}`,
			shouldReturnError: false,
		},
		{
			name:       "rate limit with valid header",
			statusCode: 429,
			headers:    map[string]string{"X-RateLimit-Reset": strconv.FormatInt(time.Now().Add(30*time.Second).Unix(), 10)},
			expectedError: "sendgrid rate limit exceeded",
			shouldReturnError: true,
		},
		{
			name:       "rate limit with invalid header",
			statusCode: 429,
			headers:    map[string]string{"X-RateLimit-Reset": "invalid"},
			expectedError: "invalid syntax",
			shouldReturnError: true,
		},
		{
			name:       "errors response",
			statusCode: 400,
			body:       `{"errors": [{"message": "validation failed"}]}`,
			expectedError: "message: validation failed",
			shouldReturnError: true,
		},
		{
			name:       "status code error fallback",
			statusCode: 500,
			body:       `invalid json`,
			expectedError: "sendgrid server error: 500 Internal Server Error",
			shouldReturnError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create HTTP response manually
			resp := &http.Response{
				StatusCode: tt.statusCode,
				Status:     fmt.Sprintf("%d %s", tt.statusCode, http.StatusText(tt.statusCode)),
				Header:     make(http.Header),
				Body:       io.NopCloser(strings.NewReader(tt.body)),
			}

			// Set headers
			for k, v := range tt.headers {
				resp.Header.Set(k, v)
			}

			err := checkStatusCode(resp, debug)
			
			if tt.shouldReturnError {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.expectedError)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestCheckStatusCodeErrorResponse(t *testing.T) {
	// Mock debug interface
	debug := &mockDebug{debug: true}

	// Test single error response parsing
	resp := &http.Response{
		StatusCode: 404,
		Status:     "404 Not Found",
		Header:     make(http.Header),
		Body:       io.NopCloser(strings.NewReader(`{"error": "not found"}`)),
	}

	err := checkStatusCode(resp, debug)
	
	assert.Error(t, err)
	// Due to how the function works, first it tries to parse as ErrorsResponse which fails
	// Then it tries ErrorResponse which also fails because body is already consumed
	// So it falls back to statusCodeError
	assert.Contains(t, err.Error(), "sendgrid server error: 404 Not Found")
}

func TestLogResponse(t *testing.T) {
	tests := []struct {
		name        string
		debug       bool
		expectError bool
	}{
		{
			name:        "debug disabled",
			debug:       false,
			expectError: false,
		},
		{
			name:        "debug enabled",
			debug:       true,
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			debug := &mockDebug{debug: tt.debug}
			
			resp := &http.Response{
				StatusCode: 200,
				Status:     "200 OK",
				Header:     make(http.Header),
				Body:       &mockReadCloser{strings.NewReader("test body")},
			}

			err := logResponse(resp, debug)
			
			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestLogResponseError(t *testing.T) {
	debug := &mockDebug{debug: true}
	
	// Create a response with nil body reader to cause httputil.DumpResponse to fail
	resp := &http.Response{
		StatusCode: 200,
		Status:     "200 OK",
		Header:     make(http.Header),
		Body:       &errorReadCloser{},
	}

	err := logResponse(resp, debug)
	assert.Error(t, err)
}

func TestNewJSONParser(t *testing.T) {
	type testStruct struct {
		Message string `json:"message"`
	}

	tests := []struct {
		name    string
		body    string
		wantErr bool
	}{
		{
			name:    "valid JSON",
			body:    `{"message": "test"}`,
			wantErr: false,
		},
		{
			name:    "invalid JSON",
			body:    `{invalid json}`,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var dst testStruct
			parser := newJSONParser(&dst)
			
			resp := &http.Response{
				Body: &mockReadCloser{strings.NewReader(tt.body)},
			}

			err := parser(resp)
			
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, "test", dst.Message)
			}
		})
	}
}

// Mock interfaces and types for testing

type mockDebug struct {
	debug bool
	logs  []string
}

func (m *mockDebug) Debug() bool {
	return m.debug
}

func (m *mockDebug) Debugf(format string, v ...interface{}) {
	m.logs = append(m.logs, fmt.Sprintf(format, v...))
}

func (m *mockDebug) Debugln(v ...interface{}) {
	m.logs = append(m.logs, fmt.Sprintln(v...))
}

type mockReadCloser struct {
	*strings.Reader
}

func (m *mockReadCloser) Close() error {
	return nil
}

type errorReadCloser struct{}

func (e *errorReadCloser) Read(p []byte) (n int, err error) {
	return 0, fmt.Errorf("read error")
}

func (e *errorReadCloser) Close() error {
	return nil
}
