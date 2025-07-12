package sendgrid

import (
	"context"
	"fmt"
	"net/http"
	"reflect"
	"testing"

	"github.com/kylelemons/godebug/pretty"
	"github.com/pkg/errors"
	"net/url"
	"strings"
)

func TestGetInboundParseWebhooks(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/user/webhooks/parse/settings", func(w http.ResponseWriter, r *http.Request) {
		if _, err := fmt.Fprint(w, `{
			"result": [
				{
					"hostname": "bar.foo",
					"url": "https://example.com",
					"spam_check": false,
					"send_raw": false
				}
			]
		}`); err != nil {
			t.Fatal(err)
		}
	})

	expected, err := client.GetInboundParseWebhooks(context.TODO())
	if err != nil {
		t.Errorf("Unexpected error: %s", err)
		return
	}

	want := []*InboundParseWebhook{
		{
			URL:       "https://example.com",
			Hostname:  "bar.foo",
			SpamCheck: false,
			SendRaw:   false,
		},
	}

	if !reflect.DeepEqual(want, expected) {
		t.Fatal(ErrIncorrectResponse, errors.New(pretty.Compare(want, expected)))
	}
}

func TestGetInboundParseWebhooks_Failed(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/user/webhooks/parse/settings", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	})

	_, err := client.GetInboundParseWebhooks(context.TODO())
	if err == nil {
		t.Fatal("expected an error but got nil")
	}
}

func TestGetInboundParseWebhook(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/user/webhooks/parse/settings/bar.foo", func(w http.ResponseWriter, r *http.Request) {
		if _, err := fmt.Fprint(w, `{
					"hostname": "bar.foo",
					"url": "https://example.com",
					"spam_check": false,
					"send_raw": false
		}`); err != nil {
			t.Fatal(err)
		}
	})

	expected, err := client.GetInboundParseWebhook(context.TODO(), "bar.foo")
	if err != nil {
		t.Errorf("Unexpected error: %s", err)
		return
	}

	want := &OutputGetInboundParseWebhook{
		URL:       "https://example.com",
		Hostname:  "bar.foo",
		SpamCheck: false,
		SendRaw:   false,
	}

	if !reflect.DeepEqual(want, expected) {
		t.Fatal(ErrIncorrectResponse, errors.New(pretty.Compare(want, expected)))
	}
}

func TestGetInboundParseWebhook_Failed(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/user/webhooks/parse/settings/bar.foo", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	})

	_, err := client.GetInboundParseWebhook(context.TODO(), "bar.foo")
	if err == nil {
		t.Fatal("expected an error but got nil")
	}
}

func TestCreateInboundParseWebhook(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/user/webhooks/parse/settings", func(w http.ResponseWriter, r *http.Request) {
		if _, err := fmt.Fprint(w, `{
			"url": "https://example.com",
			"hostname": "foo.bar",
			"spam_check": false,
			"send_raw": false
		}`); err != nil {
			t.Fatal(err)
		}
	})

	expected, err := client.CreateInboundParseWebhook(context.TODO(), &InputCreateInboundParseWebhook{
		URL:       "https://example.com",
		Hostname:  "foo.bar",
		SpamCheck: false,
		SendRaw:   false,
	})
	if err != nil {
		t.Errorf("Unexpected error: %s", err)
		return
	}

	want := &OutputCreateInboundParseWebhook{
		URL:       "https://example.com",
		Hostname:  "foo.bar",
		SpamCheck: false,
		SendRaw:   false,
	}

	if !reflect.DeepEqual(want, expected) {
		t.Fatal(ErrIncorrectResponse, errors.New(pretty.Compare(want, expected)))
	}
}

func TestCreateInboundParseWebhook_Failed(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/user/webhooks/parse/settings", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	})

	_, err := client.CreateInboundParseWebhook(context.TODO(), &InputCreateInboundParseWebhook{
		URL:       "https://example.com",
		Hostname:  "foo.bar",
		SpamCheck: false,
		SendRaw:   false,
	})
	if err == nil {
		t.Fatal("expected an error but got nil")
	}
}

func TestUpdateInboundParseWebhook(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/user/webhooks/parse/settings/foo.bar", func(w http.ResponseWriter, r *http.Request) {
		if _, err := fmt.Fprint(w, `{
			"url": "https://example.com",
			"hostname": "foo.bar",
			"spam_check": false,
			"send_raw": false
		}`); err != nil {
			t.Fatal(err)
		}
	})

	expected, err := client.UpdateInboundParseWebhook(context.TODO(), "foo.bar", &InputUpdateInboundParseWebhook{
		URL:       "https://example.com",
		SpamCheck: false,
		SendRaw:   false,
	})
	if err != nil {
		t.Errorf("Unexpected error: %s", err)
		return
	}

	want := &OutputUpdateInboundParseWebhook{
		URL:       "https://example.com",
		Hostname:  "foo.bar",
		SpamCheck: false,
		SendRaw:   false,
	}

	if !reflect.DeepEqual(want, expected) {
		t.Fatal(ErrIncorrectResponse, errors.New(pretty.Compare(want, expected)))
	}
}

func TestUpdateInboundParseWebhook_Failed(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/user/webhooks/parse/settings/foo.bar", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	})

	_, err := client.UpdateInboundParseWebhook(context.TODO(), "foo.bar", &InputUpdateInboundParseWebhook{
		URL:       "https://example.com",
		SpamCheck: false,
		SendRaw:   false,
	})
	if err == nil {
		t.Fatal("expected an error but got nil")
	}
}

func TestDeleteInboundParseWebhook(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/user/webhooks/parse/settings/foo.bar", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNoContent)
	})

	err := client.DeleteInboundParseWebhook(context.TODO(), "foo.bar")
	if err != nil {
		t.Errorf("Unexpected error: %s", err)
		return
	}
}

func TestDeleteInboundParseWebhook_Failed(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/user/webhooks/parse/settings/foo.bar", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	})

	err := client.DeleteInboundParseWebhook(context.TODO(), "foo.bar")
	if err == nil {
		t.Fatal("expected an error but got nil")
	}
}

// NewRequest Error Tests
func TestGetInboundParseWebhooks_NewRequestError(t *testing.T) {
	client, _, _, teardown := setup()
	defer teardown()

	originalBaseURL := client.baseURL
	invalidURL, _ := url.Parse("https://api.example.com/v3/")
	client.baseURL = invalidURL

	// Call method with appropriate parameters
	_, err := client.GetInboundParseWebhooks(context.TODO())
	if err == nil {
		t.Error("Expected error for invalid baseURL")
	}
	if err != nil && !strings.Contains(err.Error(), "trailing slash") {
		t.Errorf("Expected error message to contain 'trailing slash', got %v", err.Error())
	}

	client.baseURL = originalBaseURL
}

func TestGetInboundParseWebhook_NewRequestError(t *testing.T) {
	client, _, _, teardown := setup()
	defer teardown()

	originalBaseURL := client.baseURL
	invalidURL, _ := url.Parse("https://api.example.com/v3/")
	client.baseURL = invalidURL

	// Call method with appropriate parameters
	_, err := client.GetInboundParseWebhook(context.TODO(), "mail.example.com")
	if err == nil {
		t.Error("Expected error for invalid baseURL")
	}
	if err != nil && !strings.Contains(err.Error(), "trailing slash") {
		t.Errorf("Expected error message to contain 'trailing slash', got %v", err.Error())
	}

	client.baseURL = originalBaseURL
}

func TestCreateInboundParseWebhook_NewRequestError(t *testing.T) {
	client, _, _, teardown := setup()
	defer teardown()

	originalBaseURL := client.baseURL
	invalidURL, _ := url.Parse("https://api.example.com/v3/")
	client.baseURL = invalidURL

	// Call method with appropriate parameters
	input := &InputCreateInboundParseWebhook{
		Hostname: "mail.example.com",
		URL:      "https://example.com/parse",
	}
	_, err := client.CreateInboundParseWebhook(context.TODO(), input)
	if err == nil {
		t.Error("Expected error for invalid baseURL")
	}
	if err != nil && !strings.Contains(err.Error(), "trailing slash") {
		t.Errorf("Expected error message to contain 'trailing slash', got %v", err.Error())
	}

	client.baseURL = originalBaseURL
}

func TestUpdateInboundParseWebhook_NewRequestError(t *testing.T) {
	client, _, _, teardown := setup()
	defer teardown()

	originalBaseURL := client.baseURL
	invalidURL, _ := url.Parse("https://api.example.com/v3/")
	client.baseURL = invalidURL

	// Call method with appropriate parameters
	input := &InputUpdateInboundParseWebhook{
		URL:       "https://example.com/updated-parse",
		SpamCheck: true,
	}
	_, err := client.UpdateInboundParseWebhook(context.TODO(), "mail.example.com", input)
	if err == nil {
		t.Error("Expected error for invalid baseURL")
	}
	if err != nil && !strings.Contains(err.Error(), "trailing slash") {
		t.Errorf("Expected error message to contain 'trailing slash', got %v", err.Error())
	}

	client.baseURL = originalBaseURL
}

func TestDeleteInboundParseWebhook_NewRequestError(t *testing.T) {
	client, _, _, teardown := setup()
	defer teardown()

	originalBaseURL := client.baseURL
	invalidURL, _ := url.Parse("https://api.example.com/v3/")
	client.baseURL = invalidURL

	// Call method with appropriate parameters
	err := client.DeleteInboundParseWebhook(context.TODO(), "mail.example.com")
	if err == nil {
		t.Error("Expected error for invalid baseURL")
	}
	if err != nil && !strings.Contains(err.Error(), "trailing slash") {
		t.Errorf("Expected error message to contain 'trailing slash', got %v", err.Error())
	}

	client.baseURL = originalBaseURL
}
