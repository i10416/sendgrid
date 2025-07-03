package sendgrid

import (
	"bytes"
	"context"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/pkg/errors"
)

const baseURLPath string = "/v3"

var (
	ErrIncorrectResponse = errors.New("response is incorrect")
)

// setup sets up a test HTTP server along with a sendgrid.Client that is
// configured to talk to that test server. Tests should register handlers on
// mux which provide mock responses for the API method being tested.
func setup() (client *Client, mux *http.ServeMux, serverURL string, teardown func()) {
	return setupWithPath()
}

// setupWithPath sets up a test HTTP server along with a sendgrid.Client with the path.
func setupWithPath() (client *Client, mux *http.ServeMux, serverURL string, teardown func()) {
	// mux is the HTTP request multiplexer used with the test server.
	mux = http.NewServeMux()

	apiHandler := http.NewServeMux()
	apiHandler.Handle(baseURLPath+"/", http.StripPrefix(baseURLPath, mux))
	apiHandler.HandleFunc("/", func(w http.ResponseWriter, req *http.Request) {
		fmt.Fprintln(os.Stderr, "\t"+req.URL.String())
		fmt.Fprintln(os.Stderr)
	})

	// server is a test HTTP server used to provide mock API responses.
	server := httptest.NewServer(apiHandler)

	// client is the sendgrid client being tested and is
	// configured to use test server.
	client = New(
		"test-token",
		OptionSubuser("dummy"),
		OptionBaseURL(server.URL+baseURLPath),
		OptionHTTPClient(&http.Client{}),
		OptionDebug(false),
		OptionLog(log.New(os.Stderr, "kenzo0107/sendgrid", log.LstdFlags|log.Lshortfile)),
	)

	return client, mux, server.URL, server.Close
}

func testMethod(t *testing.T, r *http.Request, want string) {
	t.Helper()
	if got := r.Method; got != want {
		t.Errorf("Request method: %v, want %v", got, want)
	}
}

func TestOptionBaseURL(t *testing.T) {
	client := New("test-api-key")
	option := OptionBaseURL("https://custom.example.com")
	option(client)

	if client.baseURL.String() != "https://custom.example.com" {
		t.Errorf("Expected baseURL to be 'https://custom.example.com', got %s", client.baseURL.String())
	}
}

func TestOptionSubuser(t *testing.T) {
	client := New("test-api-key")
	option := OptionSubuser("test-subuser")
	option(client)

	if client.subuser != "test-subuser" {
		t.Errorf("Expected subuser to be 'test-subuser', got %s", client.subuser)
	}
}

func TestOptionHTTPClient(t *testing.T) {
	client := New("test-api-key")
	customClient := &http.Client{}
	option := OptionHTTPClient(customClient)
	option(client)

	if client.httpclient != customClient {
		t.Error("Expected httpclient to be the custom client")
	}
}

func TestOptionDebug(t *testing.T) {
	client := New("test-api-key")
	option := OptionDebug(true)
	option(client)

	if !client.debug {
		t.Errorf("Expected debug to be true, got %v", client.debug)
	}
}

func TestOptionLog(t *testing.T) {
	client := New("test-api-key")
	customLogger := log.New(os.Stderr, "custom-", log.LstdFlags)
	option := OptionLog(customLogger)
	option(client)

	if client.log == nil {
		t.Error("Expected log to not be nil")
	}
}

func TestDebugMethods(t *testing.T) {
	client := New("test-api-key", OptionDebug(true))

	if !client.Debug() {
		t.Error("Expected Debug() to return true")
	}

	// Test Debugf and Debugln (these should not panic)
	client.Debugf("Test debug format: %s", "message")
	client.Debugln("Test debug line")
}

func TestDebugMethodsDisabled(t *testing.T) {
	client := New("test-api-key", OptionDebug(false))

	if client.Debug() {
		t.Error("Expected Debug() to return false")
	}

	// Test Debugf and Debugln when debug is disabled (these should not panic)
	client.Debugf("Test debug format: %s", "message")
	client.Debugln("Test debug line")
}

func TestNewRequestWithTrailingSlashError(t *testing.T) {
	client := New("test-api-key", OptionBaseURL("https://api.sendgrid.com/v3/"))

	_, err := client.NewRequest("GET", "/test", nil)
	if err == nil {
		t.Error("Expected error for baseURL with trailing slash")
	}
}

func TestAddOptions(t *testing.T) {
	client := New("test-api-key")

	type testOpts struct {
		Limit  int    `url:"limit,omitempty"`
		Offset int    `url:"offset,omitempty"`
		Filter string `url:"filter,omitempty"`
	}

	opts := testOpts{
		Limit:  10,
		Offset: 0,
		Filter: "test",
	}

	result, err := client.AddOptions("/test", opts)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if result != "/test?filter=test&limit=10" {
		t.Errorf("Expected '/test?filter=test&limit=10', got %s", result)
	}
}

func TestAddOptionsWithNilPointer(t *testing.T) {
	client := New("test-api-key")

	var opts *struct{}
	result, err := client.AddOptions("/test", opts)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if result != "/test" {
		t.Errorf("Expected '/test', got %s", result)
	}
}

func TestAddOptionsWithInvalidURL(t *testing.T) {
	client := New("test-api-key")

	opts := struct{}{}
	_, err := client.AddOptions("://invalid-url", opts)
	if err == nil {
		t.Error("Expected error for invalid URL")
	}
}

func TestDoWithNilContext(t *testing.T) {
	client := New("test-api-key")

	req, err := client.NewRequest("GET", "/test", nil)
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}

	err = client.Do(nil, req, nil) //nolint:staticcheck
	if err == nil {
		t.Error("Expected error for nil context")
	}

	if err.Error() != "context must be non-nil" {
		t.Errorf("Expected error message 'context must be non-nil', got %v", err.Error())
	}
}

func TestDoWithIOWriter(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/test", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, `{"result": "success"}`)
	})

	req, err := client.NewRequest("GET", "/test", nil)
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}

	var buf bytes.Buffer
	err = client.Do(req.Context(), req, &buf)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if buf.String() != `{"result": "success"}` {
		t.Errorf("Expected buffer to contain JSON, got %v", buf.String())
	}
}

func TestDoWithEmptyResponse(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/test", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		w.WriteHeader(http.StatusOK)
	})

	req, err := client.NewRequest("GET", "/test", nil)
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}

	var result map[string]string
	err = client.Do(req.Context(), req, &result)
	if err != nil {
		t.Errorf("Expected no error for empty response, got %v", err)
	}
}

func TestDoWithNilResponse(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/test", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, `{"result": "success"}`)
	})

	req, err := client.NewRequest("GET", "/test", nil)
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}

	err = client.Do(req.Context(), req, nil)
	if err != nil {
		t.Errorf("Expected no error for nil response, got %v", err)
	}
}

func TestDoWithCancelledContext(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/test", func(w http.ResponseWriter, r *http.Request) {
		// Simulate a slow response
		select {
		case <-r.Context().Done():
			return
		case <-time.After(100 * time.Millisecond):
			w.WriteHeader(http.StatusOK)
		}
	})

	req, err := client.NewRequest("GET", "/test", nil)
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	cancel() // Cancel immediately

	err = client.Do(ctx, req, nil)
	if err == nil {
		t.Error("Expected error for cancelled context")
	}

	if err != context.Canceled {
		t.Errorf("Expected context.Canceled error, got %v", err)
	}
}
