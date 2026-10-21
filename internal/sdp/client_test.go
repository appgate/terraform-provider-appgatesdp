package sdp

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"sync/atomic"
	"testing"
	"time"
)

func TestNormalizeURL(t *testing.T) {
	tests := map[string]string{
		"https://ctrl.example.com":                 "https://ctrl.example.com:8443/admin",
		"https://ctrl.example.com:444/admin":       "https://ctrl.example.com:444/admin",
		"https://ctrl.example.com:8443/whatever?x": "https://ctrl.example.com:8443/admin",
		"http://10.0.0.1:8080":                     "http://10.0.0.1:8080/admin",
	}
	for in, want := range tests {
		got, err := NormalizeURL(in)
		if err != nil || got.String() != want {
			t.Errorf("NormalizeURL(%q) = %v, %v; want %s", in, got, err, want)
		}
	}
	for _, bad := range []string{"", "ctrl.example.com", "ftp://x", "https://"} {
		if _, err := NormalizeURL(bad); err == nil {
			t.Errorf("NormalizeURL(%q) should fail", bad)
		}
	}
}

func TestNewRequiresCredentials(t *testing.T) {
	if _, err := New(Config{URL: "https://c"}, 25); err == nil || !strings.Contains(err.Error(), "username and password") {
		t.Fatalf("got %v", err)
	}
	if _, err := New(Config{URL: "https://c", BearerToken: "t"}, 25); err != nil {
		t.Fatal(err)
	}
	if _, err := New(Config{URL: "https://c", Username: "u", Password: "p", PEMFile: "/does/not/exist"}, 25); err == nil {
		t.Fatal("missing pem file must fail")
	}
}

// fakeController counts logins and fails requests as instructed.
type fakeController struct {
	logins, calls atomic.Int32
	token         atomic.Value
	loginStatus   []int // statuses returned by successive logins before success
}

func (f *fakeController) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Header.Get("Accept") != "application/vnd.appgate.peer-v24+json" && r.URL.Path == "/admin/login" {
		w.WriteHeader(http.StatusNotAcceptable)
		_, _ = w.Write([]byte(`{"id":"not acceptable","minSupportedVersion":20,"maxSupportedVersion":24}`))
		return
	}
	if r.URL.Path == "/admin/login" {
		n := int(f.logins.Add(1))
		if n <= len(f.loginStatus) {
			w.WriteHeader(f.loginStatus[n-1])
			return
		}
		tok := "token-" + string(rune('0'+n))
		f.token.Store(tok)
		_ = json.NewEncoder(w).Encode(map[string]string{"token": tok})
		return
	}
	f.calls.Add(1)
	if tok, _ := f.token.Load().(string); r.Header.Get("Authorization") != "Bearer "+tok {
		w.WriteHeader(http.StatusUnauthorized)
		_, _ = w.Write([]byte(`{"id":"unauthorized","message":"expired"}`))
		return
	}
	switch r.URL.Path {
	case "/admin/conditions/missing":
		w.WriteHeader(http.StatusNotFound)
		_, _ = w.Write([]byte(`{"id":"not found","message":"gone"}`))
	case "/admin/conditions/invalid":
		w.WriteHeader(http.StatusUnprocessableEntity)
		_, _ = w.Write([]byte(`{"id":"validation error","message":"Validation error.","errors":[{"field":"name","message":"may not be null"}]}`))
	default:
		_, _ = w.Write([]byte(`{"id":"1","count":12345678901234567}`))
	}
}

func newTestClient(t *testing.T, f *fakeController, cfg Config) *Client {
	t.Helper()
	srv := httptest.NewServer(f)
	t.Cleanup(srv.Close)
	cfg.URL = srv.URL
	if cfg.BearerToken == "" {
		cfg.Username, cfg.Password = "admin", "admin"
	}
	c, err := New(cfg, 24)
	if err != nil {
		t.Fatal(err)
	}
	return c
}

func TestTokenIsCachedAndRefreshed(t *testing.T) {
	f := &fakeController{}
	c := newTestClient(t, f, Config{})
	ctx := context.Background()
	var out map[string]any
	for range 3 {
		if err := c.Do(ctx, http.MethodGet, "/conditions/1", nil, nil, &out); err != nil {
			t.Fatal(err)
		}
	}
	if n := f.logins.Load(); n != 1 {
		t.Fatalf("logged in %d times, want 1", n)
	}
	if out["count"].(json.Number).String() != "12345678901234567" {
		t.Fatalf("numbers must keep full precision, got %v", out["count"])
	}
	// The controller invalidates the token: the client logs in once more
	// and retries the request.
	f.token.Store("rotated")
	if err := c.Do(ctx, http.MethodGet, "/conditions/1", nil, nil, &out); err != nil {
		t.Fatal(err)
	}
	if n := f.logins.Load(); n != 2 {
		t.Fatalf("logged in %d times, want 2", n)
	}
}

func TestLoginRetriesWhileControllerStarts(t *testing.T) {
	f := &fakeController{loginStatus: []int{502, 503}}
	c := newTestClient(t, f, Config{LoginTimeout: time.Minute})
	if err := c.Do(context.Background(), http.MethodGet, "/conditions/1", nil, nil, nil); err != nil {
		t.Fatal(err)
	}
	if n := f.logins.Load(); n != 3 {
		t.Fatalf("logins = %d, want 3", n)
	}
}

func TestLoginFailsFastOnClientErrors(t *testing.T) {
	f := &fakeController{loginStatus: []int{401}}
	c := newTestClient(t, f, Config{LoginTimeout: time.Minute})
	start := time.Now()
	err := c.Do(context.Background(), http.MethodGet, "/conditions/1", nil, nil, nil)
	if err == nil || !strings.Contains(err.Error(), "HTTP 401") || time.Since(start) > 5*time.Second {
		t.Fatalf("got %v after %s", err, time.Since(start))
	}
}

func TestLoginTimeout(t *testing.T) {
	f := &fakeController{loginStatus: []int{500, 500, 500, 500, 500, 500}}
	c := newTestClient(t, f, Config{LoginTimeout: time.Second})
	if err := c.Do(context.Background(), http.MethodGet, "/x", nil, nil, nil); err == nil || !strings.Contains(err.Error(), "login failed") {
		t.Fatalf("got %v", err)
	}
}

func TestBearerTokenSkipsLogin(t *testing.T) {
	f := &fakeController{}
	f.token.Store("external")
	c := newTestClient(t, f, Config{BearerToken: "external"})
	if err := c.Do(context.Background(), http.MethodGet, "/conditions/1", nil, nil, nil); err != nil {
		t.Fatal(err)
	}
	if f.logins.Load() != 0 {
		t.Fatal("bearer token must not log in")
	}
}

func TestErrors(t *testing.T) {
	f := &fakeController{}
	c := newTestClient(t, f, Config{})
	ctx := context.Background()
	err := c.Do(ctx, http.MethodGet, "/conditions/missing", nil, nil, nil)
	if !IsNotFound(err) {
		t.Fatalf("want not found, got %v", err)
	}
	err = c.Do(ctx, http.MethodPut, "/conditions/invalid", nil, map[string]any{}, nil)
	var apiErr *Error
	if !errors.As(err, &apiErr) || apiErr.Status != 422 || len(apiErr.Errors) != 1 {
		t.Fatalf("got %#v", err)
	}
	if want := "PUT /conditions/invalid: HTTP 422 validation error: Validation error.\n  name: may not be null"; err.Error() != want {
		t.Fatalf("error message\n%s\nwant\n%s", err, want)
	}
	if IsNotFound(err) || IsNotFound(errors.New("x")) {
		t.Fatal("only 404 is not found")
	}
}

func TestSupportedVersions(t *testing.T) {
	c := newTestClient(t, &fakeController{}, Config{})
	lo, hi, err := c.SupportedVersions(context.Background())
	if err != nil || lo != 20 || hi != 24 {
		t.Fatalf("got %d %d %v", lo, hi, err)
	}
}

func TestUnknownAuthority(t *testing.T) {
	srv := httptest.NewUnstartedServer(http.NotFoundHandler())
	srv.Config.ErrorLog = log.New(io.Discard, "", 0)
	srv.StartTLS()
	defer srv.Close()
	c, err := New(Config{URL: srv.URL, Username: "u", Password: "p"}, 25)
	if err != nil {
		t.Fatal(err)
	}
	err = c.Do(context.Background(), http.MethodGet, "/x", nil, nil, nil)
	if err == nil || !strings.Contains(err.Error(), "pem_filepath") {
		t.Fatalf("got %v", err)
	}
	insecure, _ := New(Config{URL: srv.URL, Username: "u", Password: "p", Insecure: true, LoginTimeout: time.Millisecond}, 25)
	if err := insecure.Do(context.Background(), http.MethodGet, "/x", nil, nil, nil); err == nil || strings.Contains(err.Error(), "pem_filepath") {
		t.Fatalf("insecure must skip verification, got %v", err)
	}
}

func TestEscapedPathsAreSentOnce(t *testing.T) {
	var got string
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		got = r.URL.EscapedPath()
		w.WriteHeader(http.StatusNoContent)
	}))
	defer srv.Close()
	c, err := New(Config{URL: srv.URL, BearerToken: "t"}, 25)
	if err != nil {
		t.Fatal(err)
	}
	if err := c.Do(context.Background(), http.MethodDelete, "/blacklist/"+url.PathEscape("CN=a b,OU=ldap/x"), nil, nil, nil); err != nil {
		t.Fatal(err)
	}
	if want := "/admin/blacklist/CN=a%20b%2COU=ldap%2Fx"; got != want {
		t.Fatalf("path %s, want %s", got, want)
	}
}
