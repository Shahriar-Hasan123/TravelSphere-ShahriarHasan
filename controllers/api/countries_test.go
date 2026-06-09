// controllers/api/countries_test.go
package apicontrollers

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"TravelSphere/services"

	"github.com/beego/beego/v2/server/web/context"
)

func newCountriesAPIContext(method, path string) *CountriesAPIController {
	w := httptest.NewRecorder()
	req := httptest.NewRequest(method, path, nil)

	beeCtx := context.NewContext()
	beeCtx.Reset(w, req)

	ctrl := &CountriesAPIController{}
	ctrl.Ctx = beeCtx
	ctrl.Data = make(map[interface{}]interface{})

	return ctrl
}

func TestCountriesAPI_Get_NoParams(t *testing.T) {
	ctrl := newCountriesAPIContext(http.MethodGet, "/api/countries")
	// Service will call live API — just verify no panic.
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("Get() panicked: %v", r)
		}
	}()
	ctrl.Get()
}

func TestCountriesAPI_Suggestions_EmptyQuery(t *testing.T) {
	ctrl := newCountriesAPIContext(http.MethodGet, "/api/countries/suggestions?q=")
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("Suggestions() panicked: %v", r)
		}
	}()
	ctrl.Suggestions()
}

func TestCountriesAPI_Detail_InvalidSlug(t *testing.T) {
	ctrl := newCountriesAPIContext(http.MethodGet, "/api/countries/zzz-invalid")
	ctrl.Ctx.Input = context.NewInput()
	ctrl.Ctx.Input.SetParam(":slug", "zzz-invalid-slug-nobody-has")

	defer func() {
		if r := recover(); r != nil {
			t.Errorf("Detail() panicked: %v", r)
		}
	}()
	ctrl.Detail()
}

func TestCountriesAPI_Get_WithSearchParam(t *testing.T) {
	ctrl := newCountriesAPIContext(http.MethodGet, "/api/countries?search=france")
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("Get() with search panicked: %v", r)
		}
	}()
	ctrl.Get()
}

func TestCountriesAPI_Get_WithRegionParam(t *testing.T) {
	ctrl := newCountriesAPIContext(http.MethodGet, "/api/countries?region=Asia")
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("Get() with region panicked: %v", r)
		}
	}()
	ctrl.Get()
}

func TestCountriesAPI_Get_WithBothParams(t *testing.T) {
	ctrl := newCountriesAPIContext(http.MethodGet, "/api/countries?search=india&region=Asia")
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("Get() with both params panicked: %v", r)
		}
	}()
	ctrl.Get()
}

func TestCountriesAPI_Suggestions_WithQuery(t *testing.T) {
	ctrl := newCountriesAPIContext(http.MethodGet, "/api/countries/suggestions?q=fr")
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("Suggestions() with query panicked: %v", r)
		}
	}()
	ctrl.Suggestions()
}

func TestCountriesAPI_Detail_ValidSlug(t *testing.T) {
	ctrl := newCountriesAPIContext(http.MethodGet, "/api/countries/france")
	ctrl.Ctx.Input.SetParam(":slug", "france")

	defer func() {
		if r := recover(); r != nil {
			t.Errorf("Detail() with valid slug panicked: %v", r)
		}
	}()
	ctrl.Detail()
}

func TestCountriesAPI_Get_BadSearch(t *testing.T) {
	// Test with search param that has no matches
	ctrl := newCountriesAPIContext(http.MethodGet, "/api/countries?search=zzzzzzzznoplace")
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("Get() with bad search panicked: %v", r)
		}
	}()
	ctrl.Get()
}

func TestCountriesAPI_Get_BadRegion(t *testing.T) {
	// Test with invalid region
	ctrl := newCountriesAPIContext(http.MethodGet, "/api/countries?region=Atlantis")
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("Get() with bad region panicked: %v", r)
		}
	}()
	ctrl.Get()
}

func TestCountriesAPI_Suggestions_NoMatch(t *testing.T) {
	// Query that shouldn't match anything
	ctrl := newCountriesAPIContext(http.MethodGet, "/api/countries/suggestions?q=zzzzzzzz")
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("Suggestions() with no match panicked: %v", r)
		}
	}()
	ctrl.Suggestions()
}

func TestCountriesAPI_Get_ReturnsServerErrorOnBackendFailure(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{}`))
	}))
	defer server.Close()
	t.Setenv("RESTCOUNTRIES_BASE_URL", server.URL)
	countrySvc = services.NewCountryService()

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/api/countries", nil)
	beeCtx := context.NewContext()
	beeCtx.Reset(w, req)
	ctrl := &CountriesAPIController{}
	ctrl.Ctx = beeCtx
	ctrl.Data = make(map[interface{}]interface{})

	ctrl.Get()
	if w.Code != http.StatusInternalServerError {
		t.Fatalf("expected 500, got %d", w.Code)
	}

	var result map[string]interface{}
	if err := json.Unmarshal(w.Body.Bytes(), &result); err != nil {
		t.Fatalf("expected JSON body, got error: %v", err)
	}
	if msg, ok := result["message"].(string); !ok || msg != "Failed to fetch countries" {
		t.Fatalf("unexpected error response: %v", result)
	}
}

func TestCountriesAPI_Detail_ReturnsServerErrorOnBackendFailure(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{}`))
	}))
	defer server.Close()
	t.Setenv("RESTCOUNTRIES_BASE_URL", server.URL)
	countrySvc = services.NewCountryService()

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/api/countries/france", nil)
	beeCtx := context.NewContext()
	beeCtx.Reset(w, req)
	ctrl := &CountriesAPIController{}
	ctrl.Ctx = beeCtx
	ctrl.Data = make(map[interface{}]interface{})
	ctrl.Ctx.Input.SetParam(":slug", "france")

	ctrl.Detail()
	if w.Code != http.StatusInternalServerError {
		t.Fatalf("expected 500, got %d", w.Code)
	}
}

func TestCountriesAPI_Suggestions_ReturnsServerErrorOnBackendFailure(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{}`))
	}))
	defer server.Close()
	t.Setenv("RESTCOUNTRIES_BASE_URL", server.URL)
	countrySvc = services.NewCountryService()

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/api/countries/suggestions?q=fr", nil)
	beeCtx := context.NewContext()
	beeCtx.Reset(w, req)
	ctrl := &CountriesAPIController{}
	ctrl.Ctx = beeCtx
	ctrl.Data = make(map[interface{}]interface{})

	ctrl.Suggestions()
	if w.Code != http.StatusInternalServerError {
		t.Fatalf("expected 500, got %d", w.Code)
	}
}
