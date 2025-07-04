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

func TestGetAPIKeys(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/api_keys", func(w http.ResponseWriter, r *http.Request) {
		if _, err := fmt.Fprint(w, `{"result":[{
			"api_key_id": "abcdefghijklmnopqrstuv",
			"name": "full-access"
		  }]}`); err != nil {
			t.Fatal(err)
		}
	})

	expected, err := client.GetAPIKeys(context.TODO())
	if err != nil {
		t.Errorf("Unexpected error: %s", err)
		return
	}

	want := &OutputGetAPIKeys{
		APIKeys: []APIKey{
			{
				ApiKeyId: "abcdefghijklmnopqrstuv",
				Name:     "full-access",
			},
		},
	}
	if !reflect.DeepEqual(want, expected) {
		t.Fatal(ErrIncorrectResponse)
	}
}

func TestGetAPIKeys_Failed(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/api_keys", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	})

	_, err := client.GetAPIKeys(context.TODO())
	if err == nil {
		t.Fatal("expected an error but got none")
	}
}

func TestGetAPIKey(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/api_keys/dummy", func(w http.ResponseWriter, r *http.Request) {
		if _, err := fmt.Fprint(w, `{
			"scopes": [],
			"api_key_id": "dummy",
			"name": "full-accesses"
		  }`); err != nil {
			t.Fatal(err)
		}
	})

	expected, err := client.GetAPIKey(context.TODO(), "dummy")
	if err != nil {
		t.Errorf("Unexpected error: %s", err)
		return
	}

	want := &OutputGetAPIKey{
		ApiKeyId: "dummy",
		Scopes:   []string{},
		Name:     "full-accesses",
	}
	if !reflect.DeepEqual(want, expected) {
		t.Fatal(ErrIncorrectResponse)
	}
}

func TestGetAPIKey_Failed(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/api_keys/dummy", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	})

	_, err := client.GetAPIKey(context.TODO(), "dummy")
	if err == nil {
		t.Fatal("expected an error but got none")
	}
}

func TestCreateAPIKey(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/api_keys", func(w http.ResponseWriter, r *http.Request) {
		if _, err := fmt.Fprint(w, `{
			"api_key": "SG.abcdefghi",
			"api_key_id": "dummy",
			"name": "dummy",
			"scopes":[
				"user.profile.read",
				"user.profile.update"
			]
		  }`); err != nil {
			t.Fatal(err)
		}
	})

	expected, err := client.CreateAPIKey(context.TODO(), &InputCreateAPIKey{
		Name: "dummy",
		Scopes: []string{
			"user.profile.read",
			"user.profile.update",
		},
	})
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
		return
	}

	want := &OutputCreateAPIKey{
		ApiKey:   "SG.abcdefghi",
		ApiKeyId: "dummy",
		Name:     "dummy",
		Scopes: []string{
			"user.profile.read",
			"user.profile.update",
		},
	}

	if !reflect.DeepEqual(want, expected) {
		t.Fatal(ErrIncorrectResponse, errors.New(pretty.Compare(want, expected)))
	}
}

func TestCreateAPIKey_Failed(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/api_keys", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	})

	_, err := client.CreateAPIKey(context.TODO(), &InputCreateAPIKey{
		Name: "dummy",
		Scopes: []string{
			"user.profile.read",
		},
	})
	if err == nil {
		t.Fatal("expected an error but got none")
	}
}

func TestUpdateAPIKeyName(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/api_keys/dummy", func(w http.ResponseWriter, r *http.Request) {
		if _, err := fmt.Fprint(w, `{
			"api_key_id": "dummy",
			"name": "full-accesses-dummy"
		}`); err != nil {
			t.Fatal(err)
		}
	})

	expected, err := client.UpdateAPIKeyName(context.TODO(), "dummy", &InputUpdateAPIKeyName{
		Name: "full-accesses-dummy",
	})
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
		return
	}

	want := &OutputUpdateAPIKeyName{
		ApiKeyId: "dummy",
		Name:     "full-accesses-dummy",
	}
	if !reflect.DeepEqual(want, expected) {
		t.Fatal(ErrIncorrectResponse, errors.New(pretty.Compare(want, expected)))
	}
}

func TestUpdateAPIKeyName_Failed(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/api_keys/dummy", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	})

	_, err := client.UpdateAPIKeyName(context.TODO(), "dummy", &InputUpdateAPIKeyName{
		Name: "full-accesses-dummy",
	})
	if err == nil {
		t.Fatal("expected an error but got none")
	}
}

func TestUpdateAPIKeyNameAndScopes(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/api_keys/dummy", func(w http.ResponseWriter, r *http.Request) {
		if _, err := fmt.Fprint(w, `{
			"scopes": [
				"user.profile.read",
				"user.profile.update"
			],
			"api_key_id": "dummy",
			"name": "full-accesses-dummy"
		}`); err != nil {
			t.Fatal(err)
		}
	})

	expected, err := client.UpdateAPIKeyNameAndScopes(context.TODO(), "dummy", &InputUpdateAPIKeyNameAndScopes{
		Name: "full-accesses-dummy",
		Scopes: []string{
			"user.profile.read",
			"user.profile.update",
		},
	})
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
		return
	}

	want := &OutputUpdateAPIKeyNameAndScopes{
		ApiKeyId: "dummy",
		Name:     "full-accesses-dummy",
		Scopes: []string{
			"user.profile.read",
			"user.profile.update",
		},
	}
	if !reflect.DeepEqual(want, expected) {
		t.Fatal(ErrIncorrectResponse, errors.New(pretty.Compare(want, expected)))
	}
}

func TestUpdateAPIKeyNameAndScopes_Failed(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/api_keys/dummy", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	})

	_, err := client.UpdateAPIKeyNameAndScopes(context.TODO(), "dummy", &InputUpdateAPIKeyNameAndScopes{
		Name: "full-accesses-dummy",
		Scopes: []string{
			"user.profile.read",
			"user.profile.update",
		},
	})
	if err == nil {
		t.Fatal("expected an error but got none")
	}
}

func TestDeleteAPIKey(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/api_keys/dummy", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNoContent)
	})

	err := client.DeleteAPIKey(context.TODO(), "dummy")
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
		return
	}
}

func TestDeleteAPIKey_Failed(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/api_keys/dummy", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	})

	err := client.DeleteAPIKey(context.TODO(), "dummy")
	if err == nil {
		t.Fatal("expected an error but got none")
	}
}

// NewRequest Error Tests for API Key methods
func TestGetAPIKeys_NewRequestError(t *testing.T) {
	client, _, _, teardown := setup()
	defer teardown()

	originalBaseURL := client.baseURL
	invalidURL, _ := url.Parse("https://api.example.com/v3/")
	client.baseURL = invalidURL

	_, err := client.GetAPIKeys(context.TODO())
	if err == nil {
		t.Error("Expected error for invalid baseURL")
	}
	if err != nil && !strings.Contains(err.Error(), "trailing slash") {
		t.Errorf("Expected error message to contain 'trailing slash', got %v", err.Error())
	}

	client.baseURL = originalBaseURL
}

func TestGetAPIKey_NewRequestError(t *testing.T) {
	client, _, _, teardown := setup()
	defer teardown()

	originalBaseURL := client.baseURL
	invalidURL, _ := url.Parse("https://api.example.com/v3/")
	client.baseURL = invalidURL

	_, err := client.GetAPIKey(context.TODO(), "test-api-key-id")
	if err == nil {
		t.Error("Expected error for invalid baseURL")
	}
	if err != nil && !strings.Contains(err.Error(), "trailing slash") {
		t.Errorf("Expected error message to contain 'trailing slash', got %v", err.Error())
	}

	client.baseURL = originalBaseURL
}

func TestCreateAPIKey_NewRequestError(t *testing.T) {
	client, _, _, teardown := setup()
	defer teardown()

	originalBaseURL := client.baseURL
	invalidURL, _ := url.Parse("https://api.example.com/v3/")
	client.baseURL = invalidURL

	input := &InputCreateAPIKey{
		Name:   "Test API Key",
		Scopes: []string{"mail.send"},
	}
	_, err := client.CreateAPIKey(context.TODO(), input)
	if err == nil {
		t.Error("Expected error for invalid baseURL")
	}
	if err != nil && !strings.Contains(err.Error(), "trailing slash") {
		t.Errorf("Expected error message to contain 'trailing slash', got %v", err.Error())
	}

	client.baseURL = originalBaseURL
}

func TestUpdateAPIKeyName_NewRequestError(t *testing.T) {
	client, _, _, teardown := setup()
	defer teardown()

	originalBaseURL := client.baseURL
	invalidURL, _ := url.Parse("https://api.example.com/v3/")
	client.baseURL = invalidURL

	input := &InputUpdateAPIKeyName{
		Name: "Updated API Key Name",
	}
	_, err := client.UpdateAPIKeyName(context.TODO(), "test-api-key-id", input)
	if err == nil {
		t.Error("Expected error for invalid baseURL")
	}
	if err != nil && !strings.Contains(err.Error(), "trailing slash") {
		t.Errorf("Expected error message to contain 'trailing slash', got %v", err.Error())
	}

	client.baseURL = originalBaseURL
}

func TestUpdateAPIKeyNameAndScopes_NewRequestError(t *testing.T) {
	client, _, _, teardown := setup()
	defer teardown()

	originalBaseURL := client.baseURL
	invalidURL, _ := url.Parse("https://api.example.com/v3/")
	client.baseURL = invalidURL

	input := &InputUpdateAPIKeyNameAndScopes{
		Name:   "Updated API Key",
		Scopes: []string{"mail.send", "mail.schedule.write"},
	}
	_, err := client.UpdateAPIKeyNameAndScopes(context.TODO(), "test-api-key-id", input)
	if err == nil {
		t.Error("Expected error for invalid baseURL")
	}
	if err != nil && !strings.Contains(err.Error(), "trailing slash") {
		t.Errorf("Expected error message to contain 'trailing slash', got %v", err.Error())
	}

	client.baseURL = originalBaseURL
}

func TestDeleteAPIKey_NewRequestError(t *testing.T) {
	client, _, _, teardown := setup()
	defer teardown()

	originalBaseURL := client.baseURL
	invalidURL, _ := url.Parse("https://api.example.com/v3/")
	client.baseURL = invalidURL

	err := client.DeleteAPIKey(context.TODO(), "test-api-key-id")
	if err == nil {
		t.Error("Expected error for invalid baseURL")
	}
	if err != nil && !strings.Contains(err.Error(), "trailing slash") {
		t.Errorf("Expected error message to contain 'trailing slash', got %v", err.Error())
	}

	client.baseURL = originalBaseURL
}
