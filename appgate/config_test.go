package appgate

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/appgate/sdp-api-client-go/api/v14/openapi"

	"github.com/google/uuid"
)

func TestNewClient(t *testing.T) {
	clientCfg := openapi.NewConfiguration()
	c := openapi.NewAPIClient(clientCfg)

	cfg := c.GetConfig()

	if cfg.UserAgent != clientCfg.UserAgent {
		t.Fatal("Expected same base path.")
	}
}

func setup() (*openapi.APIClient, *openapi.Configuration, *http.ServeMux, *httptest.Server, func()) {
	mux := http.NewServeMux()
	server := httptest.NewServer(mux)

	clientCfg := openapi.NewConfiguration()

	clientCfg.Debug = false
	u, _ := url.Parse(server.URL)
	clientCfg.Servers = []openapi.ServerConfiguration{
		{
			URL: u.String(),
		},
	}

	c := openapi.NewAPIClient(clientCfg)

	return c, clientCfg, mux, server, server.Close
}

func testMethod(t *testing.T, r *http.Request, want string) {
	if got := r.Method; got != want {
		t.Errorf("Request method: %v, want %v", got, want)
	}
}

func TestLoginInternalServerError(t *testing.T) {
	client, _, mux, _, teardown := setup()
	defer teardown()
	mux.HandleFunc("/login", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(w, `{
  "id": "string",
  "message": "An unexpected error occurred."
}`)
		testMethod(t, r, http.MethodPost)

	})

	loginOpts := openapi.LoginRequest{
		ProviderName: "local",
		Username:     openapi.PtrString("admin"),
		Password:     openapi.PtrString("admin"),
		DeviceId:     uuid.New().String(),
	}
	_, resp, err := client.LoginApi.LoginPost(context.Background()).LoginRequest(loginOpts).Execute()
	if err == nil {
		t.Errorf("Expected error, got none")
	}

	if resp != nil && resp.StatusCode != http.StatusInternalServerError {
		t.Fatalf("Expected HTTP 500, got %v", resp.StatusCode)
	}
	oerr, ok := err.(openapi.GenericOpenAPIError)
	if !ok {
		t.Fatalf("Expected GenericOpenAPIError, got %+v", err)
	}

	m, ok := oerr.Model().(openapi.Error)
	if !ok {
		t.Fatalf("Expected openapi Error, got %+v", m)
	}
	if *m.Message != "An unexpected error occurred." {
		t.Fatalf("Expected error message 'An unexpected error occurred.', got %s", *m.Message)
	}
}

func TestGetToken200(t *testing.T) {
	client, _, mux, _, teardown := setup()
	defer teardown()
	mux.HandleFunc("/login", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, `{
  "version": "4.3.0-20000",
  "user": {
    "name": "admin",
    "needTwoFactorAuth": false,
    "canAccessAuditLogs": true,
    "privileges": [
      {
        "type": "All",
        "target": "All",
        "scope": {
          "all": true,
          "ids": [
            "4c07bc67-57ea-42dd-b702-c2d6c45419fc"
          ],
          "tags": [
            "tag"
          ]
        },
        "defaultTags": [
          "api-created"
        ]
      }
    ]
  },
  "token": "very-long-string",
  "expires": "2020-01-27T08:50:34Z",
  "messageOfTheDay": "Welcome to Appgate SDP."
}`)
		testMethod(t, r, http.MethodPost)

	})
	c := &Config{
		URL:      "http://appgate.com/admin",
		Username: "admin",
		Password: "admin",
	}
	token, err := getToken(client, c)
	if err != nil {
		t.Fatalf("Unexpected error, got %+v", err)
	}
	if token != "Bearer very-long-string" {
		t.Fatalf("Expected token very-long-string, got %s", token)
	}
}
