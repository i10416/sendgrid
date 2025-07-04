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

func TestGetSSOCertificate(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/sso/certificates/123456", func(w http.ResponseWriter, r *http.Request) {
		if _, err := fmt.Fprint(w, `{
			"id": 123456,
			"public_certificate": "-----BEGIN CERTIFICATE-----\nMIIC0DCCAbigAwIBAgIJAOT==\n-----END CERTIFICATE-----",
			"not_before": 1586137600,
			"not_after": 1586137600,
			"integration_id": "abcdef"
		}`); err != nil {
			t.Fatal(err)
		}
	})

	expected, err := client.GetSSOCertificate(context.TODO(), 123456)
	if err != nil {
		t.Errorf("Unexpected error: %s", err)
		return
	}

	want := &OutputGetSSOCertificate{
		ID:                123456,
		PublicCertificate: "-----BEGIN CERTIFICATE-----\nMIIC0DCCAbigAwIBAgIJAOT==\n-----END CERTIFICATE-----",
		NotBefore:         1586137600,
		NotAfter:          1586137600,
		IntegrationID:     "abcdef",
	}
	if !reflect.DeepEqual(want, expected) {
		t.Fatal(errors.New(pretty.Compare(want, expected)))
	}
}

func TestGetSSOCertificate_Failed(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/sso/certificates/123456", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	})

	_, err := client.GetSSOCertificate(context.TODO(), 123456)
	if err == nil {
		t.Error("expected an error")
	}
}

func TestGetSSOCertificates(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/sso/integrations/abcdef/certificates", func(w http.ResponseWriter, r *http.Request) {
		if _, err := fmt.Fprint(w, `[
			{
				"id": 123456,
				"public_certificate": "-----BEGIN CERTIFICATE-----\nMIIC0DCCAbigAwIBAgIJAOT==\n-----END CERTIFICATE-----",
				"not_before": 1586137600,
				"not_after": 1586137600,
				"integration_id": "abcdef"
			}
		]`); err != nil {
			t.Fatal(err)
		}
	})

	expected, err := client.GetSSOCertificates(context.TODO(), "abcdef")
	if err != nil {
		t.Errorf("Unexpected error: %s", err)
		return
	}

	want := []*SSOCertificate{
		{
			ID:                123456,
			PublicCertificate: "-----BEGIN CERTIFICATE-----\nMIIC0DCCAbigAwIBAgIJAOT==\n-----END CERTIFICATE-----",
			NotBefore:         1586137600,
			NotAfter:          1586137600,
			IntegrationID:     "abcdef",
		},
	}
	if !reflect.DeepEqual(want, expected) {
		t.Fatal(errors.New(pretty.Compare(want, expected)))
	}
}

func TestGetSSOCertificates_Failed(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/sso/integrations/abcdef/certificates", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	})

	_, err := client.GetSSOCertificates(context.TODO(), "abcdef")
	if err == nil {
		t.Error("expected an error")
	}
}

func TestCreateSSOCertificate(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/sso/certificates", func(w http.ResponseWriter, r *http.Request) {
		if _, err := fmt.Fprint(w, `{
			"id": 123456,
			"public_certificate": "-----BEGIN CERTIFICATE-----\nMIIC0DCCAbigAwIBAgIJAOT==\n-----END CERTIFICATE-----",
			"not_before": 1586137600,
			"not_after": 1586137600,
			"integration_id": "abcdef"
		}`); err != nil {
			t.Fatal(err)
		}
	})

	expected, err := client.CreateSSOCertificate(context.TODO(), &InputCreateSSOCertificate{
		PublicCertificate: "-----BEGIN CERTIFICATE-----\nMIIC0DCCAbigAwIBAgIJAOT==\n-----END CERTIFICATE-----",
		Enabled:           true,
		IntegrationID:     "abcdef",
	})
	if err != nil {
		t.Errorf("Unexpected error: %s", err)
		return
	}

	want := &OutputCreateSSOCertificate{
		ID:                123456,
		PublicCertificate: "-----BEGIN CERTIFICATE-----\nMIIC0DCCAbigAwIBAgIJAOT==\n-----END CERTIFICATE-----",
		NotBefore:         1586137600,
		NotAfter:          1586137600,
		IntegrationID:     "abcdef",
	}
	if !reflect.DeepEqual(want, expected) {
		t.Fatal(errors.New(pretty.Compare(want, expected)))
	}
}

func TestCreateSSOCertificate_Failed(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/sso/certificates", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	})

	_, err := client.CreateSSOCertificate(context.TODO(), &InputCreateSSOCertificate{
		PublicCertificate: "-----BEGIN CERTIFICATE-----\nMIIC0DCCAbigAwIBAgIJAOT==\n-----END CERTIFICATE-----",
		Enabled:           true,
		IntegrationID:     "abcdef",
	})
	if err == nil {
		t.Error("expected an error")
	}
}

func TestUpdateSSOCertificate(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/sso/certificates/123456", func(w http.ResponseWriter, r *http.Request) {
		if _, err := fmt.Fprint(w, `{
			"id": 123456,
			"public_certificate": "-----BEGIN CERTIFICATE-----\nMIIC0DCCAbigAwIBAgIJAOT==\n-----END CERTIFICATE-----",
			"not_before": 1586137600,
			"not_after": 1586137600,
			"integration_id": "abcdef"
		}`); err != nil {
			t.Fatal(err)
		}
	})

	expected, err := client.UpdateSSOCertificate(context.TODO(), 123456, &InputUpdateSSOCertificate{
		PublicCertificate: "-----BEGIN CERTIFICATE-----\nMIIC0DCCAbigAwIBAgIJAOT==\n-----END CERTIFICATE-----",
		Enabled:           true,
		IntegrationID:     "abcdef",
	})
	if err != nil {
		t.Errorf("Unexpected error: %s", err)
		return
	}

	want := &OutputUpdateSSOCertificate{
		ID:                123456,
		PublicCertificate: "-----BEGIN CERTIFICATE-----\nMIIC0DCCAbigAwIBAgIJAOT==\n-----END CERTIFICATE-----",
		NotBefore:         1586137600,
		NotAfter:          1586137600,
		IntegrationID:     "abcdef",
	}
	if !reflect.DeepEqual(want, expected) {
		t.Fatal(errors.New(pretty.Compare(want, expected)))
	}
}

func TestUpdateSSOCertificate_Failed(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/sso/certificates/123456", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	})

	_, err := client.UpdateSSOCertificate(context.TODO(), 123456, &InputUpdateSSOCertificate{
		PublicCertificate: "-----BEGIN CERTIFICATE-----\nMIIC0DCCAbigAwIBAgIJAOT==\n-----END CERTIFICATE-----",
		Enabled:           true,
		IntegrationID:     "abcdef",
	})
	if err == nil {
		t.Error("expected an error")
	}
}

func TestDeleteSSOCertificate(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/sso/certificates/123456", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNoContent)
	})

	if err := client.DeleteSSOCertificate(context.TODO(), 123456); err != nil {
		t.Errorf("Unexpected error: %s", err)
		return
	}
}

func TestDeleteSSOCertificate_Failed(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/sso/certificates/123456", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	})

	err := client.DeleteSSOCertificate(context.TODO(), 123456)
	if err == nil {
		t.Error("expected an error")
	}
}

// NewRequest Error Tests
func TestGetSSOCertificate_NewRequestError(t *testing.T) {
	client, _, _, teardown := setup()
	defer teardown()

	originalBaseURL := client.baseURL
	invalidURL, _ := url.Parse("https://api.example.com/v3/")
	client.baseURL = invalidURL

	_, err := client.GetSSOCertificate(context.TODO(), 12345)
	if err == nil {
		t.Error("Expected error for invalid baseURL")
	}
	if err != nil && !strings.Contains(err.Error(), "trailing slash") {
		t.Errorf("Expected error message to contain 'trailing slash', got %v", err.Error())
	}

	client.baseURL = originalBaseURL
}

func TestGetSSOCertificates_NewRequestError(t *testing.T) {
	client, _, _, teardown := setup()
	defer teardown()

	originalBaseURL := client.baseURL
	invalidURL, _ := url.Parse("https://api.example.com/v3/")
	client.baseURL = invalidURL

	_, err := client.GetSSOCertificates(context.TODO(), "test-integration-id")
	if err == nil {
		t.Error("Expected error for invalid baseURL")
	}
	if err != nil && !strings.Contains(err.Error(), "trailing slash") {
		t.Errorf("Expected error message to contain 'trailing slash', got %v", err.Error())
	}

	client.baseURL = originalBaseURL
}

func TestCreateSSOCertificate_NewRequestError(t *testing.T) {
	client, _, _, teardown := setup()
	defer teardown()

	originalBaseURL := client.baseURL
	invalidURL, _ := url.Parse("https://api.example.com/v3/")
	client.baseURL = invalidURL

	input := &InputCreateSSOCertificate{
		PublicCertificate: "-----BEGIN CERTIFICATE-----\ntest certificate\n-----END CERTIFICATE-----",
		IntegrationID:     "test-integration-id",
	}
	_, err := client.CreateSSOCertificate(context.TODO(), input)
	if err == nil {
		t.Error("Expected error for invalid baseURL")
	}
	if err != nil && !strings.Contains(err.Error(), "trailing slash") {
		t.Errorf("Expected error message to contain 'trailing slash', got %v", err.Error())
	}

	client.baseURL = originalBaseURL
}

func TestUpdateSSOCertificate_NewRequestError(t *testing.T) {
	client, _, _, teardown := setup()
	defer teardown()

	originalBaseURL := client.baseURL
	invalidURL, _ := url.Parse("https://api.example.com/v3/")
	client.baseURL = invalidURL

	input := &InputUpdateSSOCertificate{
		PublicCertificate: "-----BEGIN CERTIFICATE-----\nupdated certificate\n-----END CERTIFICATE-----",
		IntegrationID:     "updated-integration-id",
	}
	_, err := client.UpdateSSOCertificate(context.TODO(), 12345, input)
	if err == nil {
		t.Error("Expected error for invalid baseURL")
	}
	if err != nil && !strings.Contains(err.Error(), "trailing slash") {
		t.Errorf("Expected error message to contain 'trailing slash', got %v", err.Error())
	}

	client.baseURL = originalBaseURL
}

func TestDeleteSSOCertificate_NewRequestError(t *testing.T) {
	client, _, _, teardown := setup()
	defer teardown()

	originalBaseURL := client.baseURL
	invalidURL, _ := url.Parse("https://api.example.com/v3/")
	client.baseURL = invalidURL

	err := client.DeleteSSOCertificate(context.TODO(), 12345)
	if err == nil {
		t.Error("Expected error for invalid baseURL")
	}
	if err != nil && !strings.Contains(err.Error(), "trailing slash") {
		t.Errorf("Expected error message to contain 'trailing slash', got %v", err.Error())
	}

	client.baseURL = originalBaseURL
}
