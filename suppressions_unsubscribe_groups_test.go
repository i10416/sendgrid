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

func TestGetSuppressionGroup(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/asm/groups/12345", func(w http.ResponseWriter, r *http.Request) {
		if _, err := fmt.Fprint(w, `{
			"id":12345,
			"name":"dummy",
			"description":"dummy description",
			"is_default": false,
			"unsubscribes": 0,
			"last_email_sent_at": ""
		}`); err != nil {
			t.Fatal(err)
		}
	})

	expected, err := client.GetSuppressionGroup(context.TODO(), 12345)
	if err != nil {
		t.Errorf("Unexpected error: %s", err)
		return
	}

	want := &SuppressionGroup{
		ID:              12345,
		Name:            "dummy",
		Description:     "dummy description",
		IsDefault:       false,
		Unsubscribes:    0,
		LastEmailSentAt: "",
	}
	if !reflect.DeepEqual(want, expected) {
		t.Fatal(ErrIncorrectResponse)
	}
}

func TestGetSuppressionGroup_Failed(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/asm/groups/12345", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	})

	_, err := client.GetSuppressionGroup(context.TODO(), 12345)
	if err == nil {
		t.Fatal("expected an error but got none")
	}
}

func TestGetSuppressionGroups(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/asm/groups", func(w http.ResponseWriter, r *http.Request) {
		if _, err := fmt.Fprint(w, `[{
			"id":12345,
			"name":"dummy",
			"description":"dummy description",
			"is_default": false,
			"unsubscribes": 0,
			"last_email_sent_at": ""
		}]`); err != nil {
			t.Fatal(err)
		}
	})

	expected, err := client.GetSuppressionGroups(context.TODO())
	if err != nil {
		t.Errorf("Unexpected error: %s", err)
		return
	}

	want := []*SuppressionGroup{
		{
			ID:              12345,
			Name:            "dummy",
			Description:     "dummy description",
			IsDefault:       false,
			Unsubscribes:    0,
			LastEmailSentAt: "",
		},
	}
	if !reflect.DeepEqual(want, expected) {
		t.Fatal(ErrIncorrectResponse)
	}
}

func TestGetSuppressionGroups_Failed(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/asm/groups", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	})

	_, err := client.GetSuppressionGroups(context.TODO())
	if err == nil {
		t.Fatal("expected an error but got none")
	}
}

func TestCreateSuppressionGroup(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/asm/groups", func(w http.ResponseWriter, r *http.Request) {
		if _, err := fmt.Fprint(w, `{
			"id":12345,
			"name":"dummy",
			"description":"dummy description",
			"is_default": false
		}`); err != nil {
			t.Fatal(err)
		}
	})

	expected, err := client.CreateSuppressionGroup(context.TODO(), &InputCreateSuppressionGroup{
		Name:        "dummy",
		Description: "dummy description",
		IsDefault:   false,
	})
	if err != nil {
		t.Errorf("Unexpected error: %s", err)
		return
	}

	want := &OutputCreateSuppressionGroup{
		ID:          12345,
		Name:        "dummy",
		Description: "dummy description",
		IsDefault:   false,
	}
	if !reflect.DeepEqual(want, expected) {
		t.Fatal(ErrIncorrectResponse)
	}
}

func TestCreateSuppressionGroup_Failed(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/asm/groups", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	})

	_, err := client.CreateSuppressionGroup(context.TODO(), &InputCreateSuppressionGroup{
		Name:        "dummy",
		Description: "dummy description",
		IsDefault:   false,
	})
	if err == nil {
		t.Fatal("expected an error but got none")
	}
}

func TestUpdateSuppressionGroup(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/asm/groups/12345", func(w http.ResponseWriter, r *http.Request) {
		if _, err := fmt.Fprint(w, `{
			"id":12345,
			"name":"dummy",
			"description":"dummy description",
			"is_default": false
		}`); err != nil {
			t.Fatal(err)
		}
	})

	expected, err := client.UpdateSuppressionGroup(context.TODO(), 12345, &InputUpdateSuppressionGroup{
		Name:        "dummy",
		Description: "dummy description",
		IsDefault:   false,
	})
	if err != nil {
		t.Errorf("Unexpected error: %s", err)
		return
	}

	want := &OutputUpdateSuppressionGroup{
		ID:              12345,
		Name:            "dummy",
		Description:     "dummy description",
		IsDefault:       false,
		LastEmailSentAt: "",
		Unsubscribes:    0,
	}
	if !reflect.DeepEqual(want, expected) {
		t.Fatal(ErrIncorrectResponse)
	}
}

func TestUpdateSuppressionGroup_Failed(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/asm/groups", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	})

	_, err := client.UpdateSuppressionGroup(context.TODO(), 12345, &InputUpdateSuppressionGroup{
		Name:        "dummy",
		Description: "dummy description",
		IsDefault:   false,
	})
	if err == nil {
		t.Fatal("expected an error but got none")
	}
}

func TestDeleteSuppressionGroup(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/asm/groups/12345", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNoContent)
	})

	if err := client.DeleteSuppressionGroup(context.TODO(), 12345); err != nil {
		t.Errorf("Unexpected error: %s", err)
		return
	}
}

func TestDeleteSuppressionGroup_Failed(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/asm/groups", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	})

	err := client.DeleteSuppressionGroup(context.TODO(), 12345)
	if err == nil {
		t.Fatal("expected an error but got none")
	}
}

// NewRequest Error Tests
func TestGetSuppressionGroup_NewRequestError(t *testing.T) {
	client, _, _, teardown := setup()
	defer teardown()

	originalBaseURL := client.baseURL
	invalidURL, _ := url.Parse("https://api.example.com/v3/")
	client.baseURL = invalidURL

	_, err := client.GetSuppressionGroup(context.TODO(), 12345)
	if err == nil {
		t.Error("Expected error for invalid baseURL")
	}
	if err != nil && !strings.Contains(err.Error(), "trailing slash") {
		t.Errorf("Expected error message to contain 'trailing slash', got %v", err.Error())
	}

	client.baseURL = originalBaseURL
}

func TestGetSuppressionGroups_NewRequestError(t *testing.T) {
	client, _, _, teardown := setup()
	defer teardown()

	originalBaseURL := client.baseURL
	invalidURL, _ := url.Parse("https://api.example.com/v3/")
	client.baseURL = invalidURL

	_, err := client.GetSuppressionGroups(context.TODO())
	if err == nil {
		t.Error("Expected error for invalid baseURL")
	}
	if err != nil && !strings.Contains(err.Error(), "trailing slash") {
		t.Errorf("Expected error message to contain 'trailing slash', got %v", err.Error())
	}

	client.baseURL = originalBaseURL
}

func TestCreateSuppressionGroup_NewRequestError(t *testing.T) {
	client, _, _, teardown := setup()
	defer teardown()

	originalBaseURL := client.baseURL
	invalidURL, _ := url.Parse("https://api.example.com/v3/")
	client.baseURL = invalidURL

	input := &InputCreateSuppressionGroup{
		Name:        "test",
		Description: "test description",
	}
	_, err := client.CreateSuppressionGroup(context.TODO(), input)
	if err == nil {
		t.Error("Expected error for invalid baseURL")
	}
	if err != nil && !strings.Contains(err.Error(), "trailing slash") {
		t.Errorf("Expected error message to contain 'trailing slash', got %v", err.Error())
	}

	client.baseURL = originalBaseURL
}

func TestUpdateSuppressionGroup_NewRequestError(t *testing.T) {
	client, _, _, teardown := setup()
	defer teardown()

	originalBaseURL := client.baseURL
	invalidURL, _ := url.Parse("https://api.example.com/v3/")
	client.baseURL = invalidURL

	input := &InputUpdateSuppressionGroup{
		Name:        "updated test",
		Description: "updated description",
	}
	_, err := client.UpdateSuppressionGroup(context.TODO(), 12345, input)
	if err == nil {
		t.Error("Expected error for invalid baseURL")
	}
	if err != nil && !strings.Contains(err.Error(), "trailing slash") {
		t.Errorf("Expected error message to contain 'trailing slash', got %v", err.Error())
	}

	client.baseURL = originalBaseURL
}

func TestDeleteSuppressionGroup_NewRequestError(t *testing.T) {
	client, _, _, teardown := setup()
	defer teardown()

	originalBaseURL := client.baseURL
	invalidURL, _ := url.Parse("https://api.example.com/v3/")
	client.baseURL = invalidURL

	err := client.DeleteSuppressionGroup(context.TODO(), 12345)
	if err == nil {
		t.Error("Expected error for invalid baseURL")
	}
	if err != nil && !strings.Contains(err.Error(), "trailing slash") {
		t.Errorf("Expected error message to contain 'trailing slash', got %v", err.Error())
	}

	client.baseURL = originalBaseURL
}
