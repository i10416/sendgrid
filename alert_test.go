package sendgrid

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"reflect"
	"strings"
	"testing"
)

func TestGetAlert(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/alerts/1234567", func(w http.ResponseWriter, r *http.Request) {
		if _, err := fmt.Fprint(w, `{
			"id": 1234567,
			"type": "usage_limit",
			"email_to": "dummy@example.com",
			"frequency": "daily",
			"percentage": 90,
			"created_at": 1599999999,
			"updated_at": 1599999999
		}`); err != nil {
			t.Fatal(err)
		}
	})

	expected, err := client.GetAlert(context.Background(), 1234567)
	if err != nil {
		t.Errorf("Unexpected error: %s", err)
		return
	}

	want := &OutputGetAlert{
		ID:         1234567,
		EmailTo:    "dummy@example.com",
		Frequency:  "daily",
		Type:       "usage_limit",
		Percentage: 90,
		CreatedAt:  1599999999,
		UpdatedAt:  1599999999,
	}

	if !reflect.DeepEqual(want, expected) {
		t.Fatal(ErrIncorrectResponse)
	}
}

func TestGetAlert_Failed(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/alerts/1234567", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	})

	_, err := client.GetAlert(context.Background(), 1234567)
	if err == nil {
		t.Fatal("expected an error but got none")
	}
}

func TestGetAlerts(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/alerts", func(w http.ResponseWriter, r *http.Request) {
		if _, err := fmt.Fprint(w, `[
			{
				"id": 1234567,
				"type": "usage_limit",
				"email_to": "dummy@example.com",
				"frequency": "daily",
				"percentage": 90,
				"created_at": 1599999999,
				"updated_at": 1599999999
			}
		]`); err != nil {
			t.Fatal(err)
		}
	})

	expected, err := client.GetAlerts(context.Background())
	if err != nil {
		t.Errorf("Unexpected error: %s", err)
		return
	}

	want := []*Alert{
		{
			ID:         1234567,
			EmailTo:    "dummy@example.com",
			Frequency:  "daily",
			Type:       "usage_limit",
			Percentage: 90,
			CreatedAt:  1599999999,
			UpdatedAt:  1599999999,
		},
	}

	if !reflect.DeepEqual(want, expected) {
		t.Fatal(ErrIncorrectResponse)
	}
}

func TestGetAlerts_Failed(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/alerts", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	})

	_, err := client.GetAlerts(context.Background())
	if err == nil {
		t.Fatal("expected an error but got none")
	}
}

func TestCreateAlert(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/alerts", func(w http.ResponseWriter, r *http.Request) {
		if _, err := fmt.Fprint(w, `{
			"id": 1234567,
			"type": "usage_limit",
			"email_to": "dummy@example.com",
			"frequency": "daily",
			"percentage": 90,
			"created_at": 1599999999,
			"updated_at": 1599999999
		}`); err != nil {
			t.Fatal(err)
		}
	})

	input := &InputCreateAlert{
		Type:       "usage_limit",
		EmailTo:    "dummy@example.com",
		Frequency:  "daily",
		Percentage: 90,
	}
	expected, err := client.CreateAlert(context.Background(), input)
	if err != nil {
		t.Errorf("Unexpected error: %s", err)
		return
	}

	want := &OutputCreateAlert{
		ID:         1234567,
		EmailTo:    "dummy@example.com",
		Frequency:  "daily",
		Type:       "usage_limit",
		Percentage: 90,
		CreatedAt:  1599999999,
		UpdatedAt:  1599999999,
	}

	if !reflect.DeepEqual(want, expected) {
		t.Fatal(ErrIncorrectResponse)
	}
}

func TestCreateAlert_Failed(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/alerts", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	})

	input := &InputCreateAlert{
		Type:       "usage_limit",
		EmailTo:    "dummy@example.com",
		Frequency:  "daily",
		Percentage: 90,
	}
	_, err := client.CreateAlert(context.Background(), input)
	if err == nil {
		t.Fatal("expected an error but got none")
	}
}

func TestUpdateAlert(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/alerts/1234567", func(w http.ResponseWriter, r *http.Request) {
		if _, err := fmt.Fprint(w, `{
			"id": 1234567,
			"type": "usage_limit",
			"email_to": "dummy@example.com",
			"frequency": "daily",
			"percentage": 90,
			"created_at": 1599999999,
			"updated_at": 1599999999
		}`); err != nil {
			t.Fatal(err)
		}
	})

	input := &InputUpdateAlert{
		EmailTo:    "dummy@example.com",
		Frequency:  "daily",
		Percentage: 90,
	}
	expected, err := client.UpdateAlert(context.Background(), 1234567, input)
	if err != nil {
		t.Errorf("Unexpected error: %s", err)
		return
	}

	want := &OutputUpdateAlert{
		ID:         1234567,
		EmailTo:    "dummy@example.com",
		Frequency:  "daily",
		Type:       "usage_limit",
		Percentage: 90,
		CreatedAt:  1599999999,
		UpdatedAt:  1599999999,
	}

	if !reflect.DeepEqual(want, expected) {
		t.Fatal(ErrIncorrectResponse)
	}
}

func TestUpdateAlert_Failed(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/alerts/1234567", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	})

	input := &InputUpdateAlert{
		EmailTo:    "dummy@example.com",
		Frequency:  "daily",
		Percentage: 90,
	}
	_, err := client.UpdateAlert(context.Background(), 1234567, input)
	if err == nil {
		t.Fatal("expected an error but got none")
	}
}

func TestDeleteAlert(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/alerts/1234567", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNoContent)
	})

	err := client.DeleteAlert(context.Background(), 1234567)
	if err != nil {
		t.Errorf("Unexpected error: %s", err)
		return
	}
}

func TestDeleteAlert_Failed(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/alerts/1234567", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	})

	err := client.DeleteAlert(context.Background(), 1234567)
	if err == nil {
		t.Fatal("expected an error but got none")
	}
}

// NewRequest Error Tests for Alert methods
func TestGetAlert_NewRequestError(t *testing.T) {
	client, _, _, teardown := setup()
	defer teardown()

	originalBaseURL := client.baseURL
	invalidURL, _ := url.Parse("https://api.example.com/v3/")
	client.baseURL = invalidURL

	_, err := client.GetAlert(context.TODO(), 1234567)
	if err == nil {
		t.Error("Expected error for invalid baseURL")
	}
	if err != nil && !strings.Contains(err.Error(), "trailing slash") {
		t.Errorf("Expected error message to contain 'trailing slash', got %v", err.Error())
	}

	client.baseURL = originalBaseURL
}

func TestGetAlerts_NewRequestError(t *testing.T) {
	client, _, _, teardown := setup()
	defer teardown()

	originalBaseURL := client.baseURL
	invalidURL, _ := url.Parse("https://api.example.com/v3/")
	client.baseURL = invalidURL

	_, err := client.GetAlerts(context.TODO())
	if err == nil {
		t.Error("Expected error for invalid baseURL")
	}
	if err != nil && !strings.Contains(err.Error(), "trailing slash") {
		t.Errorf("Expected error message to contain 'trailing slash', got %v", err.Error())
	}

	client.baseURL = originalBaseURL
}

func TestCreateAlert_NewRequestError(t *testing.T) {
	client, _, _, teardown := setup()
	defer teardown()

	originalBaseURL := client.baseURL
	invalidURL, _ := url.Parse("https://api.example.com/v3/")
	client.baseURL = invalidURL

	input := &InputCreateAlert{
		Type:      "usage_limit",
		EmailTo:   "test@example.com",
		Frequency: "daily",
		Percentage: 90,
	}
	_, err := client.CreateAlert(context.TODO(), input)
	if err == nil {
		t.Error("Expected error for invalid baseURL")
	}
	if err != nil && !strings.Contains(err.Error(), "trailing slash") {
		t.Errorf("Expected error message to contain 'trailing slash', got %v", err.Error())
	}

	client.baseURL = originalBaseURL
}

func TestUpdateAlert_NewRequestError(t *testing.T) {
	client, _, _, teardown := setup()
	defer teardown()

	originalBaseURL := client.baseURL
	invalidURL, _ := url.Parse("https://api.example.com/v3/")
	client.baseURL = invalidURL

	input := &InputUpdateAlert{
		EmailTo:   "updated@example.com",
		Frequency: "weekly",
	}
	_, err := client.UpdateAlert(context.TODO(), 1234567, input)
	if err == nil {
		t.Error("Expected error for invalid baseURL")
	}
	if err != nil && !strings.Contains(err.Error(), "trailing slash") {
		t.Errorf("Expected error message to contain 'trailing slash', got %v", err.Error())
	}

	client.baseURL = originalBaseURL
}

func TestDeleteAlert_NewRequestError(t *testing.T) {
	client, _, _, teardown := setup()
	defer teardown()

	originalBaseURL := client.baseURL
	invalidURL, _ := url.Parse("https://api.example.com/v3/")
	client.baseURL = invalidURL

	err := client.DeleteAlert(context.TODO(), 1234567)
	if err == nil {
		t.Error("Expected error for invalid baseURL")
	}
	if err != nil && !strings.Contains(err.Error(), "trailing slash") {
		t.Errorf("Expected error message to contain 'trailing slash', got %v", err.Error())
	}

	client.baseURL = originalBaseURL
}

