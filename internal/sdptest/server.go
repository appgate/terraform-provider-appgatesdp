// Package sdptest provides an in-memory SDP Admin API for tests. It serves
// the collections described by the generated metadata of one API version and
// mimics the behavior the provider depends on: versioned Accept headers,
// login and tokens, server assigned ids and timestamps, write-only fields
// that are never returned, validation errors and 404s.
package sdptest

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"slices"
	"strings"
	"sync"
	"testing"

	"github.com/google/uuid"

	"github.com/appgate/terraform-provider-appgatesdp/internal/apimeta"
)

const (
	Username = "admin"
	Password = "secret"
)

// Request is a request received by the server.
type Request struct {
	Method, Path, Accept string
	Query                string
	Body                 map[string]any
}

// Server is a fake controller for one API version.
type Server struct {
	*httptest.Server
	t       testing.TB
	version int
	// MinVersion and MaxVersion are reported in 406 login responses.
	MinVersion, MaxVersion int
	api                    *apimeta.API

	mu      sync.Mutex
	token   string
	objects map[string]map[string]map[string]any // collection path -> id -> object
	singles map[string]map[string]any
	// endpoints holds responses of read-only data source paths, by pattern
	// (e.g. /sites/{id}/resources).
	endpoints map[string]map[string]any
	requests  []Request
	failures  []failure
	// LowercaseTags mimics the controller lowercasing tag values.
	LowercaseTags bool
	// PendingStatusPolls is the number of /appliances/status responses that
	// report appliances as still registering before they become ready.
	PendingStatusPolls int
	statusPolls        int
}

type failure struct {
	method, path string
	status       int
	body         string
	times        int
}

// New starts a server for the given API version metadata.
func New(t testing.TB, api *apimeta.API) *Server {
	t.Helper()
	var n int
	if _, err := fmt.Sscanf(api.Version, "v%d", &n); err != nil {
		t.Fatalf("bad version %q", api.Version)
	}
	s := &Server{
		t: t, version: n, MinVersion: n, MaxVersion: n, api: api,
		token:     uuid.NewString(),
		objects:   map[string]map[string]map[string]any{},
		singles:   map[string]map[string]any{},
		endpoints: map[string]map[string]any{},
	}
	for _, r := range allResources(api) {
		if r.ReadOnly {
			s.endpoints[r.Path] = map[string]any{}
			continue
		}
		if r.Singleton {
			if s.singles[r.Path] == nil {
				s.singles[r.Path] = map[string]any{}
			}
		} else if s.objects[r.Path] == nil {
			s.objects[r.Path] = map[string]map[string]any{}
		}
	}
	s.Server = httptest.NewServer(http.HandlerFunc(s.serve))
	t.Cleanup(s.Close)
	return s
}

func allResources(api *apimeta.API) []*apimeta.Resource {
	var out []*apimeta.Resource
	for _, m := range []map[string]*apimeta.Resource{api.Resources, api.DataSources} {
		for _, k := range sortedKeys(m) {
			out = append(out, m[k])
		}
	}
	return out
}

// Fail makes the next `times` requests matching method and path fail.
func (s *Server) Fail(method, path string, status int, body string, times int) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.failures = append(s.failures, failure{method, path, status, body, times})
}

// SeedEndpoint sets the response of a read-only endpoint pattern.
func (s *Server) SeedEndpoint(pattern string, obj map[string]any) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.endpoints[pattern] = obj
}

// AcceptOverride reports whether an endpoint documents its own media type.
func (s *Server) AcceptOverride(path, accept string) bool {
	return s.acceptOverride(path, accept)
}

func (s *Server) acceptOverride(path, accept string) bool {
	for _, r := range allResources(s.api) {
		if r.Accept != "" && r.Accept == accept && matchPattern(r.Path, path) {
			return true
		}
	}
	return false
}

// matchPattern reports whether path matches a pattern with {param} segments.
func matchPattern(pattern, path string) bool {
	ps, xs := strings.Split(pattern, "/"), strings.Split(path, "/")
	if len(ps) != len(xs) {
		return false
	}
	for i := range ps {
		if ps[i] != xs[i] && !strings.HasPrefix(ps[i], "{") {
			return false
		}
	}
	return true
}

// Token returns the token issued by successful logins.
func (s *Server) Token() string {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.token
}

// Requests returns all requests received so far.
func (s *Server) Requests() []Request {
	s.mu.Lock()
	defer s.mu.Unlock()
	return slices.Clone(s.requests)
}

// Seed stores an object in a collection, returning its id.
func (s *Server) Seed(collection string, obj map[string]any) string {
	s.mu.Lock()
	defer s.mu.Unlock()
	id, _ := obj["id"].(string)
	if id == "" {
		id = uuid.NewString()
		obj["id"] = id
	}
	if s.objects[collection] == nil {
		s.objects[collection] = map[string]map[string]any{}
	}
	s.objects[collection][id] = obj
	return id
}

// Object returns a copy of a stored object, or nil.
func (s *Server) Object(collection, id string) map[string]any {
	s.mu.Lock()
	defer s.mu.Unlock()
	if collection != "" && id == "" {
		return cloneJSON(s.singles[collection])
	}
	return cloneJSON(s.objects[collection][id])
}

// Update modifies a stored object out of band, simulating drift.
func (s *Server) Update(collection, id string, fn func(map[string]any)) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if id == "" {
		fn(s.singles[collection])
		return
	}
	fn(s.objects[collection][id])
}

// Delete removes an object out of band.
func (s *Server) Delete(collection, id string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.objects[collection], id)
}

// IDs returns the ids of a collection's objects.
func (s *Server) IDs(collection string) []string {
	s.mu.Lock()
	defer s.mu.Unlock()
	return sortedKeys(s.objects[collection])
}

// Count returns the number of objects in a collection.
func (s *Server) Count(collection string) int {
	s.mu.Lock()
	defer s.mu.Unlock()
	return len(s.objects[collection])
}

func (s *Server) serve(w http.ResponseWriter, r *http.Request) {
	s.mu.Lock()
	defer s.mu.Unlock()
	path := strings.TrimPrefix(r.URL.Path, "/admin")
	var body map[string]any
	if b, _ := io.ReadAll(r.Body); len(b) > 0 {
		dec := json.NewDecoder(bytes.NewReader(b))
		dec.UseNumber()
		if err := dec.Decode(&body); err != nil {
			writeJSON(w, http.StatusBadRequest, map[string]any{"id": "json error", "message": err.Error()})
			return
		}
	}
	s.requests = append(s.requests, Request{Method: r.Method, Path: path, Accept: r.Header.Get("Accept"), Query: r.URL.RawQuery, Body: body})

	for i, f := range s.failures {
		if f.times > 0 && f.method == r.Method && f.path == path {
			s.failures[i].times--
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(f.status)
			_, _ = io.WriteString(w, f.body)
			return
		}
	}

	var v int
	if _, err := fmt.Sscanf(r.Header.Get("Accept"), "application/vnd.appgate.peer-v%d+json", &v); (err != nil || v < s.MinVersion || v > s.MaxVersion) && !s.acceptOverride(path, r.Header.Get("Accept")) {
		writeJSON(w, http.StatusNotAcceptable, map[string]any{
			"id": "not acceptable", "message": "Invalid Accept header",
			"minSupportedVersion": s.MinVersion, "maxSupportedVersion": s.MaxVersion,
		})
		return
	}
	if path == "/login" && r.Method == http.MethodPost {
		if body["username"] != Username || body["password"] != Password {
			writeJSON(w, http.StatusUnauthorized, map[string]any{"id": "unauthorized", "message": "Invalid username or password."})
			return
		}
		writeJSON(w, http.StatusOK, map[string]any{"token": s.token, "expires": "2099-01-01T00:00:00Z"})
		return
	}
	if r.Header.Get("Authorization") != "Bearer "+s.token {
		writeJSON(w, http.StatusUnauthorized, map[string]any{"id": "unauthorized", "message": "Token is invalid."})
		return
	}
	if path == "/appliances/status" && r.Method == http.MethodGet {
		s.statusPolls++
		var data []any
		for _, id := range sortedKeys(s.objects["/appliances"]) {
			obj := s.objects["/appliances"][id]
			state := "appliance_ready"
			if ctrl, _ := obj["controller"].(map[string]any); ctrl["enabled"] == true {
				state = "controller_ready"
			}
			if s.statusPolls <= s.PendingStatusPolls {
				state = "appliance_registering"
			}
			data = append(data, map[string]any{"id": id, "name": obj["name"], "state": state})
		}
		writeJSON(w, http.StatusOK, map[string]any{"data": data})
		return
	}
	if strings.HasSuffix(path, "/deactivate") && r.Method == http.MethodPost {
		coll, id := splitObjectPath(strings.TrimSuffix(path, "/deactivate"))
		if obj := s.objects[coll][id]; obj != nil {
			obj["activated"] = false
			w.WriteHeader(http.StatusAccepted)
			return
		}
	}
	if strings.HasSuffix(path, "/export") && r.Method == http.MethodPost {
		coll, id := splitObjectPath(strings.TrimSuffix(path, "/export"))
		if obj := s.objects[coll][id]; obj != nil {
			writeJSON(w, http.StatusOK, map[string]any{"applianceId": id, "sshConfig": body})
			return
		}
	}
	if r.Method == http.MethodGet {
		if obj, ok := s.endpoints[path]; ok {
			writeJSON(w, http.StatusOK, obj)
			return
		}
	}
	// Exact paths win over patterns (/discovered-apps/configuration is a
	// singleton, /discovered-apps/{id} an endpoint).
	if single, ok := s.singles[path]; ok {
		s.serveSingleton(w, r, path, single, body)
		return
	}
	if coll, ok := s.objects[path]; ok {
		s.serveCollection(w, r, path, coll, body)
		return
	}
	if r.Method == http.MethodGet {
		for _, pattern := range sortedKeys(s.endpoints) {
			if strings.Contains(pattern, "{") && matchPattern(pattern, path) {
				writeJSON(w, http.StatusOK, s.endpoints[pattern])
				return
			}
		}
	}
	if collPath, id := splitObjectPath(path); s.objects[collPath] != nil {
		s.serveObject(w, r, collPath, id, body)
		return
	}
	writeJSON(w, http.StatusNotFound, map[string]any{"id": "not found", "message": "No such endpoint " + path})
}

func splitObjectPath(path string) (string, string) {
	i := strings.LastIndex(path, "/")
	if i <= 0 {
		return path, ""
	}
	return path[:i], path[i+1:]
}

func (s *Server) resource(collection string, obj map[string]any) *apimeta.Resource {
	var fallback *apimeta.Resource
	for _, r := range allResources(s.api) {
		if r.Path != collection {
			continue
		}
		if r.Discriminator == nil {
			fallback = r
			continue
		}
		if obj != nil && obj[r.Discriminator.Property] == r.Discriminator.Value {
			return r
		}
	}
	if fallback == nil {
		for _, r := range allResources(s.api) {
			if r.Path == collection {
				return r
			}
		}
	}
	return fallback
}

func (s *Server) serveSingleton(w http.ResponseWriter, r *http.Request, path string, obj, body map[string]any) {
	res := s.resource(path, nil)
	switch r.Method {
	case http.MethodGet:
		writeJSON(w, http.StatusOK, s.view(res, obj))
	case http.MethodPut, http.MethodPost:
		// Wrapped bodies are free-form and accepted as they are.
		if errs := validate(res.Fields, body, ""); len(errs) > 0 && !res.Wrapped {
			writeValidation(w, errs)
			return
		}
		s.singles[path] = s.normalize(body)
		writeJSON(w, http.StatusOK, s.view(res, s.singles[path]))
	case http.MethodDelete:
		s.singles[path] = map[string]any{}
		w.WriteHeader(http.StatusNoContent)
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func (s *Server) serveCollection(w http.ResponseWriter, r *http.Request, path string, coll map[string]map[string]any, body map[string]any) {
	switch r.Method {
	case http.MethodGet:
		query := r.URL.Query().Get("query")
		var data []any
		for _, id := range sortedKeys(coll) {
			obj := coll[id]
			if name, _ := obj["name"].(string); query == "" || strings.Contains(strings.ToLower(name), strings.ToLower(query)) {
				data = append(data, s.view(s.resource(path, obj), obj))
			}
		}
		listKey := "data"
		if res := s.resource(path, nil); res != nil {
			listKey = res.ListProperty()
		}
		writeJSON(w, http.StatusOK, map[string]any{listKey: data, "range": fmt.Sprintf("0-%d/%d", len(data), len(data))})
	case http.MethodPost:
		res := s.resource(path, body)
		if errs := validate(res.Fields, body, ""); len(errs) > 0 {
			writeValidation(w, errs)
			return
		}
		obj := s.normalize(body)
		idKey := res.ID()
		id, _ := obj[idKey].(string)
		if id == "" {
			id = uuid.NewString()
		} else if _, exists := coll[id]; exists {
			writeJSON(w, http.StatusConflict, map[string]any{"id": "conflict", "message": "Object with id already exists."})
			return
		}
		obj[idKey] = id
		if _, ok := res.Fields["created"]; ok {
			obj["created"] = "2026-01-01T00:00:00Z"
		}
		if _, ok := res.Fields["updated"]; ok {
			obj["updated"] = "2026-01-01T00:00:00Z"
		}
		coll[id] = obj
		writeJSON(w, http.StatusOK, s.view(res, obj))
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func (s *Server) serveObject(w http.ResponseWriter, r *http.Request, path, id string, body map[string]any) {
	obj, ok := s.objects[path][id]
	if !ok {
		writeJSON(w, http.StatusNotFound, map[string]any{"id": "not found", "message": "Object not found."})
		return
	}
	res := s.resource(path, obj)
	switch r.Method {
	case http.MethodGet:
		writeJSON(w, http.StatusOK, s.view(res, obj))
	case http.MethodPut:
		if errs := validate(res.Fields, body, ""); len(errs) > 0 {
			writeValidation(w, errs)
			return
		}
		updated := s.normalize(body)
		for _, k := range []string{"id", "created"} {
			if v, ok := obj[k]; ok {
				updated[k] = v
			}
		}
		if _, ok := res.Fields["updated"]; ok {
			updated["updated"] = "2026-01-02T00:00:00Z"
		}
		s.objects[path][id] = updated
		writeJSON(w, http.StatusOK, s.view(res, updated))
	case http.MethodDelete:
		delete(s.objects[path], id)
		w.WriteHeader(http.StatusNoContent)
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

// normalize applies server side normalization to a request body.
func (s *Server) normalize(body map[string]any) map[string]any {
	obj := cloneJSON(body)
	if obj == nil {
		obj = map[string]any{}
	}
	if tags, ok := obj["tags"].([]any); ok && s.LowercaseTags {
		for i, t := range tags {
			if str, ok := t.(string); ok {
				tags[i] = strings.ToLower(str)
			}
		}
		slices.SortFunc(tags, func(a, b any) int { return strings.Compare(fmt.Sprint(a), fmt.Sprint(b)) })
	}
	return obj
}

// view is the response representation: write-only fields are never returned.
func (s *Server) view(res *apimeta.Resource, obj map[string]any) map[string]any {
	out := cloneJSON(obj)
	if res != nil {
		stripWriteOnly(out, res.Fields)
	}
	return out
}

func stripWriteOnly(obj map[string]any, fields map[string]*apimeta.Field) {
	for _, f := range fields {
		if f.WriteOnly {
			delete(obj, f.JSON)
			continue
		}
		switch v := obj[f.JSON].(type) {
		case map[string]any:
			if f.Fields != nil {
				stripWriteOnly(v, f.Fields)
			}
		case []any:
			for _, e := range v {
				if m, ok := e.(map[string]any); ok && f.Fields != nil {
					stripWriteOnly(m, f.Fields)
				}
			}
		}
	}
}

// validate rejects missing required fields, read-only fields and fields the
// API version does not know, like a strict controller would.
func validate(fields map[string]*apimeta.Field, body map[string]any, prefix string) []map[string]any {
	var errs []map[string]any
	byJSON := map[string]*apimeta.Field{}
	for _, f := range fields {
		byJSON[f.JSON] = f
		if f.Required && !f.ReadOnly {
			if v, ok := body[f.JSON]; !ok || v == nil {
				errs = append(errs, map[string]any{"field": prefix + f.JSON, "message": "may not be null"})
			}
		}
	}
	for _, k := range sortedKeys(body) {
		f, ok := byJSON[k]
		if !ok {
			if k != "id" && k != "type" {
				errs = append(errs, map[string]any{"field": prefix + k, "message": "unknown field"})
			}
			continue
		}
		if f.Fields == nil {
			continue
		}
		switch v := body[k].(type) {
		case map[string]any:
			errs = append(errs, validate(f.Fields, v, prefix+k+".")...)
		case []any:
			for _, e := range v {
				if m, ok := e.(map[string]any); ok {
					errs = append(errs, validate(f.Fields, m, prefix+k+".")...)
				}
			}
		}
	}
	return errs
}

func writeValidation(w http.ResponseWriter, errs []map[string]any) {
	writeJSON(w, http.StatusUnprocessableEntity, map[string]any{"id": "validation error", "message": "Validation error.", "errors": errs})
}

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}

func cloneJSON(m map[string]any) map[string]any {
	if m == nil {
		return nil
	}
	b, _ := json.Marshal(m)
	dec := json.NewDecoder(bytes.NewReader(b))
	dec.UseNumber()
	var out map[string]any
	_ = dec.Decode(&out)
	return out
}

func sortedKeys[V any](m map[string]V) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	slices.Sort(keys)
	return keys
}
