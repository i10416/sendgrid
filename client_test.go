package sendgrid

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestNew(t *testing.T) {
	client := New("test-key")
	if client == nil {
		t.Fatal("expected client to be non-nil")
	}
	if client.apiKey != "test-key" {
		t.Errorf("expected apiKey to be 'test-key', got %s", client.apiKey)
	}
	if client.baseURL.String() != "https://api.sendgrid.com/v3" {
		t.Errorf("expected baseURL to be 'https://api.sendgrid.com/v3', got %s", client.baseURL.String())
	}
}

func TestNewWithOptions(t *testing.T) {
	client := New("test-key",
		OptionBaseURL("https://custom.api.com"),
		OptionSubuser("test-subuser"),
		OptionDebug(true),
	)

	if client.baseURL.String() != "https://custom.api.com" {
		t.Errorf("expected baseURL to be 'https://custom.api.com', got %s", client.baseURL.String())
	}
	if client.subuser != "test-subuser" {
		t.Errorf("expected subuser to be 'test-subuser', got %s", client.subuser)
	}
	if !client.debug {
		t.Error("expected debug to be true")
	}
}

func TestClient_NewRequest(t *testing.T) {
	client := New("test-key")

	req, err := client.NewRequest("GET", "/test", nil)
	if err != nil {
		t.Fatal(err)
	}

	if req.Method != "GET" {
		t.Errorf("expected method to be 'GET', got %s", req.Method)
	}
	if req.URL.Path != "/v3/test" {
		t.Errorf("expected path to be '/v3/test', got %s", req.URL.Path)
	}
	if req.Header.Get("Authorization") != "Bearer test-key" {
		t.Errorf("expected authorization header to be 'Bearer test-key', got %s", req.Header.Get("Authorization"))
	}
}

func TestClient_NewRequestWithBody(t *testing.T) {
	client := New("test-key")

	body := map[string]string{"test": "value"}
	req, err := client.NewRequest("POST", "/test", body)
	if err != nil {
		t.Fatal(err)
	}

	if req.Method != "POST" {
		t.Errorf("expected method to be 'POST', got %s", req.Method)
	}
	if req.Header.Get("Content-Type") != "application/json" {
		t.Errorf("expected content-type to be 'application/json', got %s", req.Header.Get("Content-Type"))
	}
}

func TestClient_NewRequestWithSubuser(t *testing.T) {
	client := New("test-key", OptionSubuser("test-subuser"))

	req, err := client.NewRequest("GET", "/test", nil)
	if err != nil {
		t.Fatal(err)
	}

	if req.Header.Get("On-Behalf-Of") != "test-subuser" {
		t.Errorf("expected On-Behalf-Of header to be 'test-subuser', got %s", req.Header.Get("On-Behalf-Of"))
	}
}

func TestClient_Do(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"test": "value"}`))
	}))
	defer server.Close()

	client := New("test-key", OptionBaseURL(server.URL))

	req, err := client.NewRequest("GET", "/test", nil)
	if err != nil {
		t.Fatal(err)
	}

	var result map[string]string
	err = client.Do(context.Background(), req, &result)
	if err != nil {
		t.Fatal(err)
	}

	if result["test"] != "value" {
		t.Errorf("expected result to be 'value', got %s", result["test"])
	}
}

func TestClient_DoWithError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte(`{"error": "test error"}`))
	}))
	defer server.Close()

	client := New("test-key", OptionBaseURL(server.URL))

	req, err := client.NewRequest("GET", "/test", nil)
	if err != nil {
		t.Fatal(err)
	}

	var result map[string]string
	err = client.Do(context.Background(), req, &result)
	if err == nil {
		t.Fatal("expected error but got none")
	}
}

func TestClient_Debug(t *testing.T) {
	client := New("test-key", OptionDebug(true))
	if !client.Debug() {
		t.Error("expected debug to be true")
	}

	client = New("test-key", OptionDebug(false))
	if client.Debug() {
		t.Error("expected debug to be false")
	}
}

func TestClient_Debugf(t *testing.T) {
	var buf bytes.Buffer
	client := New("test-key", OptionDebug(true))
	client.log = &mockLogger{writer: &buf}

	client.Debugf("test %s", "message")

	if buf.Len() == 0 {
		t.Error("expected debug message to be written")
	}
}

func TestClient_Debugln(t *testing.T) {
	var buf bytes.Buffer
	client := New("test-key", OptionDebug(true))
	client.log = &mockLogger{writer: &buf}

	client.Debugln("test", "message")

	if buf.Len() == 0 {
		t.Error("expected debug message to be written")
	}
}

type mockLogger struct {
	writer io.Writer
}

func (m *mockLogger) Output(calldepth int, s string) error {
	_, err := m.writer.Write([]byte(s))
	return err
}

func (m *mockLogger) Print(v ...interface{}) {
	_, _ = m.writer.Write([]byte(fmt.Sprint(v...)))
}

func (m *mockLogger) Printf(format string, v ...interface{}) {
	_, _ = m.writer.Write([]byte(fmt.Sprintf(format, v...)))
}

func (m *mockLogger) Println(v ...interface{}) {
	_, _ = m.writer.Write([]byte(fmt.Sprintln(v...)))
}

func TestBool(t *testing.T) {
	b := Bool(true)
	if b == nil {
		t.Fatal("expected non-nil bool pointer")
	}
	if *b != true {
		t.Error("expected bool value to be true")
	}
}

func TestString(t *testing.T) {
	s := String("test")
	if s == nil {
		t.Fatal("expected non-nil string pointer")
	}
	if *s != "test" {
		t.Error("expected string value to be 'test'")
	}
}
