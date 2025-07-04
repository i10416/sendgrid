package sendgrid

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"reflect"
	"strings"
	"testing"

	"github.com/kylelemons/godebug/pretty"
	"github.com/pkg/errors"
)

func TestGetEnforceTLS(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/user/settings/enforced_tls", func(w http.ResponseWriter, r *http.Request) {
		if _, err := fmt.Fprint(w, `{
			"require_tls": false,
			"require_valid_cert": false,
			"version": 1.1
		}`); err != nil {
			t.Fatal(err)
		}
	})

	expected, err := client.GetEnforceTLS(context.TODO())
	if err != nil {
		t.Errorf("Unexpected error: %s", err)
		return
	}

	want := &OutputGetEnforceTLS{
		RequireTLS:       false,
		RequireValidCert: false,
		Version:          1.1,
	}

	if !reflect.DeepEqual(want, expected) {
		t.Fatal(ErrIncorrectResponse, errors.New(pretty.Compare(want, expected)))
	}
}

func TestGetEnforceTLS_Failed(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/user/settings/enforced_tls", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	})

	_, err := client.GetEnforceTLS(context.TODO())
	if err == nil {
		t.Fatal("expected an error but got none")
	}
}

func TestUpdateEnforceTLS(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/user/settings/enforced_tls", func(w http.ResponseWriter, r *http.Request) {
		if _, err := fmt.Fprint(w, `{
			"require_tls": false,
			"require_valid_cert": false,
			"version": 1.1
		}`); err != nil {
			t.Fatal(err)
		}
	})

	expected, err := client.UpdateEnforceTLS(context.TODO(), &InputUpdateEnforceTLS{
		RequireTLS:       false,
		RequireValidCert: false,
		Version:          1.1,
	})
	if err != nil {
		t.Errorf("Unexpected error: %s", err)
		return
	}

	want := &OutputUpdateEnforceTLS{
		RequireTLS:       false,
		RequireValidCert: false,
		Version:          1.1,
	}

	if !reflect.DeepEqual(want, expected) {
		t.Fatal(ErrIncorrectResponse, errors.New(pretty.Compare(want, expected)))
	}
}

func TestUpdateEnforceTLS_Failed(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/user/settings/enforced_tls", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	})

	_, err := client.UpdateEnforceTLS(context.TODO(), &InputUpdateEnforceTLS{
		RequireTLS:       false,
		RequireValidCert: false,
		Version:          1.1,
	})
	if err == nil {
		t.Fatal("expected an error but got none")
	}
}

// NewRequest Error Tests for Enforce TLS methods
func TestGetEnforceTLS_NewRequestError(t *testing.T) {
	client, _, _, teardown := setup()
	defer teardown()

	originalBaseURL := client.baseURL
	invalidURL, _ := url.Parse("https://api.example.com/v3/")
	client.baseURL = invalidURL

	_, err := client.GetEnforceTLS(context.TODO())
	if err == nil {
		t.Error("Expected error for invalid baseURL")
	}
	if err != nil && !strings.Contains(err.Error(), "trailing slash") {
		t.Errorf("Expected error message to contain 'trailing slash', got %v", err.Error())
	}

	client.baseURL = originalBaseURL
}

func TestUpdateEnforceTLS_NewRequestError(t *testing.T) {
	client, _, _, teardown := setup()
	defer teardown()

	originalBaseURL := client.baseURL
	invalidURL, _ := url.Parse("https://api.example.com/v3/")
	client.baseURL = invalidURL

	input := &InputUpdateEnforceTLS{
		RequireTLS:       true,
		RequireValidCert: true,
		Version:          1.2,
	}
	_, err := client.UpdateEnforceTLS(context.TODO(), input)
	if err == nil {
		t.Error("Expected error for invalid baseURL")
	}
	if err != nil && !strings.Contains(err.Error(), "trailing slash") {
		t.Errorf("Expected error message to contain 'trailing slash', got %v", err.Error())
	}

	client.baseURL = originalBaseURL
}
