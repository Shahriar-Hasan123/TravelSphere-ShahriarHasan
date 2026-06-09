// filters/auth_test.go
package filters

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/beego/beego/v2/server/web/context"
)

func newAuthTestContext(method, path string) (*context.Context, *httptest.ResponseRecorder) {
	w := httptest.NewRecorder()
	req := httptest.NewRequest(method, path, nil)
	ctx := context.NewContext()
	ctx.Reset(w, req)
	return ctx, w
}

func TestIsAPIRequest_True(t *testing.T) {
	paths := []string{
		"/api/wishlist",
		"/api/countries",
		"/api/dashboard/summary",
		"/api/attractions",
	}
	for _, path := range paths {
		t.Run(path, func(t *testing.T) {
			ctx, _ := newAuthTestContext(http.MethodGet, path)
			if !isAPIRequest(ctx) {
				t.Errorf("expected %q to be an API request", path)
			}
		})
	}
}

func TestIsAPIRequest_False(t *testing.T) {
	paths := []string{"/wishlist", "/dashboard", "/countries", "/", "/login"}
	for _, path := range paths {
		t.Run(path, func(t *testing.T) {
			ctx, _ := newAuthTestContext(http.MethodGet, path)
			if isAPIRequest(ctx) {
				t.Errorf("expected %q to not be an API request", path)
			}
		})
	}
}

func TestIsAPIRequest_ShortPath(t *testing.T) {
	ctx, _ := newAuthTestContext(http.MethodGet, "/api")
	if isAPIRequest(ctx) {
		t.Error("path '/api' (len 4) should not match")
	}
}

func TestRespondUnauthorized_WritesJSON(t *testing.T) {
	ctx, w := newAuthTestContext(http.MethodGet, "/api/wishlist")
	respondUnauthorized(ctx)

	// Read from the raw recorder — Beego writes through to it.
	body := w.Body.Bytes()
	if len(body) == 0 {
		t.Fatal("expected non-empty response body")
	}

	var result map[string]interface{}
	if err := json.Unmarshal(body, &result); err != nil {
		t.Fatalf("body is not valid JSON: %v", err)
	}
	if result["status"] != "error" {
		t.Errorf("expected status 'error', got %v", result["status"])
	}
	if result["message"] != "Authentication required" {
		t.Errorf("unexpected message: %v", result["message"])
	}
}

func TestRespondUnauthorized_ContentType(t *testing.T) {
	ctx, w := newAuthTestContext(http.MethodGet, "/api/wishlist")
	respondUnauthorized(ctx)

	ct := w.Header().Get("Content-Type")
	if ct != "application/json" {
		t.Errorf("expected application/json content-type, got %q", ct)
	}
}

func TestAuthFilter_NoSession_APIRoute(t *testing.T) {
	ctx, w := newAuthTestContext(http.MethodGet, "/api/wishlist")

	// Run the filter — no session set, API route — should write 401 JSON.
	AuthFilter(ctx)

	body := w.Body.Bytes()
	if len(body) == 0 {
		t.Fatal("expected JSON response body for unauthenticated API request")
	}

	var result map[string]interface{}
	if err := json.Unmarshal(body, &result); err != nil {
		t.Fatalf("response is not valid JSON: %v", err)
	}
	if result["status"] != "error" {
		t.Errorf("expected error status, got %v", result["status"])
	}
}

func TestAuthFilter_NoSession_SSRRoute(t *testing.T) {
	ctx, w := newAuthTestContext(http.MethodGet, "/dashboard")
	AuthFilter(ctx)

	// SSR route should redirect to /login.
	location := w.Header().Get("Location")
	if location != "/login" {
		t.Errorf("expected redirect to /login, got %q", location)
	}
}

func TestAuthFilter_WithSession_Allowed(t *testing.T) {
	ctx, w := newAuthTestContext(http.MethodGet, "/api/wishlist")

	// Simulate an authenticated session by setting the session value.
	ctx.Input.CruSession = &mockSession{data: map[interface{}]interface{}{
		"username": "beta",
	}}

	AuthFilter(ctx)

	// With a valid session the filter should pass through — no 401 body.
	body := w.Body.Bytes()
	if len(body) > 0 {
		t.Errorf("expected no response body for authenticated request, got: %s", body)
	}
}

func TestAuthFilter_SessionWithNilUsername_APIRoute(t *testing.T) {
	ctx, w := newAuthTestContext(http.MethodGet, "/api/wishlist")

	// Session exists but username is nil
	ctx.Input.CruSession = &mockSession{data: make(map[interface{}]interface{})}

	AuthFilter(ctx)

	// Should respond with 401 JSON for API route
	body := w.Body.Bytes()
	if len(body) == 0 {
		t.Fatal("expected JSON response body for nil username API request")
	}

	var result map[string]interface{}
	if err := json.Unmarshal(body, &result); err != nil {
		t.Fatalf("response is not valid JSON: %v", err)
	}
	if result["status"] != "error" {
		t.Errorf("expected error status, got %v", result["status"])
	}
}

func TestAuthFilter_SessionWithNilUsername_SSRRoute(t *testing.T) {
	ctx, w := newAuthTestContext(http.MethodGet, "/dashboard")

	// Session exists but username is nil
	ctx.Input.CruSession = &mockSession{data: make(map[interface{}]interface{})}

	AuthFilter(ctx)

	// SSR route should redirect to /login
	location := w.Header().Get("Location")
	if location != "/login" {
		t.Errorf("expected redirect to /login, got %q", location)
	}
}

func TestAuthFilter_SessionWithEmptyUsername_APIRoute(t *testing.T) {
	ctx, w := newAuthTestContext(http.MethodGet, "/api/countries")

	// Session exists but username is empty string
	ctx.Input.CruSession = &mockSession{data: map[interface{}]interface{}{
		"username": "",
	}}

	AuthFilter(ctx)

	// Should respond with 401 JSON for API route
	body := w.Body.Bytes()
	if len(body) == 0 {
		t.Fatal("expected JSON response body for empty username API request")
	}

	var result map[string]interface{}
	if err := json.Unmarshal(body, &result); err != nil {
		t.Fatalf("response is not valid JSON: %v", err)
	}
	if result["status"] != "error" {
		t.Errorf("expected error status, got %v", result["status"])
	}
}

func TestAuthFilter_SessionWithEmptyUsername_SSRRoute(t *testing.T) {
	ctx, w := newAuthTestContext(http.MethodGet, "/wishlist")

	// Session exists but username is empty string
	ctx.Input.CruSession = &mockSession{data: map[interface{}]interface{}{
		"username": "",
	}}

	AuthFilter(ctx)

	// SSR route should redirect to /login
	location := w.Header().Get("Location")
	if location != "/login" {
		t.Errorf("expected redirect to /login, got %q", location)
	}
}
