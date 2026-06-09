package controllers

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"TravelSphere/models"
	"TravelSphere/utils/clients"

	"github.com/beego/beego/v2/server/web/context"
)

func newCountryController(t *testing.T, method, path string) (*CountryController, *context.Context, *httptest.ResponseRecorder) {
	t.Helper()
	w := httptest.NewRecorder()
	req := httptest.NewRequest(method, path, nil)
	ctx := context.NewContext()
	ctx.Reset(w, req)
	ctx.Input.CruSession = &mockSession{data: make(map[interface{}]interface{})}
	ctrl := &CountryController{}
	ctrl.Ctx = ctx
	ctrl.Data = make(map[interface{}]interface{})
	return ctrl, ctx, w
}

func mockServer(t *testing.T, status int, payload interface{}) *httptest.Server {
	t.Helper()
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(status)
		if payload != nil {
			if err := json.NewEncoder(w).Encode(payload); err != nil {
				t.Fatalf("encode payload: %v", err)
			}
		}
	}))
}

func TestCountryController_Get_ShowsErrorWhenServiceFails(t *testing.T) {
	server := mockServer(t, http.StatusInternalServerError, nil)
	defer server.Close()
	t.Setenv("RESTCOUNTRIES_BASE_URL", server.URL)

	ctrl, _, _ := newCountryController(t, http.MethodGet, "/countries")
	ctrl.Get()

	if got := ctrl.Data["Error"].(string); got != "Unable to load countries. Please try again later." {
		t.Errorf("expected error message, got %q", got)
	}
	if ctrl.Data["Countries"] != nil {
		t.Error("expected Countries nil when service fails")
	}
}

func TestCountryController_Get_RenderedWhenCountriesAvailable(t *testing.T) {
	server := mockServer(t, http.StatusOK, []clients.RawCountry{{
		Name: struct {
			Common   string "json:\"common\""
			Official string "json:\"official\""
		}{Common: "United States", Official: "United States of America"},
		Flags: struct {
			PNG string "json:\"png\""
			SVG string "json:\"svg\""
		}{PNG: "https://example.com/us.png"},
		Capital:    []string{"Washington"},
		Population: 331002651,
		Region:     "Americas",
		Subregion:  "Northern America",
		Currencies: map[string]struct {
			Name   string "json:\"name\""
			Symbol string "json:\"symbol\""
		}{"USD": {Name: "United States dollar", Symbol: "$"}},
		Languages: map[string]string{"eng": "English"},
		Latlng:    []float64{38.8977, -77.0365},
	}})
	defer server.Close()
	t.Setenv("RESTCOUNTRIES_BASE_URL", server.URL)

	ctrl, _, _ := newCountryController(t, http.MethodGet, "/countries?search=Uni&region=Americas")
	ctrl.Ctx.Request.URL.RawQuery = "search=Uni&region=Americas"
	ctrl.Get()

	if got := ctrl.Data["SearchQuery"].(string); got != "Uni" {
		t.Errorf("expected SearchQuery Uni, got %q", got)
	}
	if got := ctrl.Data["RegionFilter"].(string); got != "Americas" {
		t.Errorf("expected RegionFilter Americas, got %q", got)
	}
	countries, ok := ctrl.Data["Countries"].([]models.Country)
	if !ok {
		t.Fatal("expected Countries to be a []models.Country")
	}
	if len(countries) != 1 {
		t.Fatalf("expected 1 country, got %d", len(countries))
	}
	if got := ctrl.TplName; got != "countries.tpl" {
		t.Errorf("expected countries.tpl, got %q", got)
	}
}

func TestCountryController_Detail_Returns404WhenSlugMissing(t *testing.T) {
	server := mockServer(t, http.StatusOK, []clients.RawCountry{})
	defer server.Close()
	t.Setenv("RESTCOUNTRIES_BASE_URL", server.URL)

	ctrl, _, _ := newCountryController(t, http.MethodGet, "/countries/unknown")
	ctrl.Ctx.Input.SetParam(":slug", "unknown")
	ctrl.Detail()

	if got := ctrl.TplName; got != "404.tpl" {
		t.Errorf("expected 404.tpl, got %q", got)
	}
	if status := ctrl.Ctx.Output.Status; status != http.StatusNotFound {
		t.Errorf("expected status 404, got %d", status)
	}
}

func TestCountryController_Detail_RendersDestinationWithAttractionsAndWeather(t *testing.T) {
	countryServer := mockServer(t, http.StatusOK, []clients.RawCountry{{
		Name: struct {
			Common   string "json:\"common\""
			Official string "json:\"official\""
		}{Common: "Testland", Official: "Testland Republic"},
		Flags: struct {
			PNG string "json:\"png\""
			SVG string "json:\"svg\""
		}{PNG: "https://example.com/testland.png"},
		Capital:    []string{"Testville"},
		Population: 123456,
		Region:     "Test Region",
		Subregion:  "Unit Test Area",
		Currencies: map[string]struct {
			Name   string "json:\"name\""
			Symbol string "json:\"symbol\""
		}{"TST": {Name: "Test Dollar", Symbol: "T$"}},
		Languages: map[string]string{"eng": "English"},
		Latlng:    []float64{12.34, 56.78},
	}})
	defer countryServer.Close()

	attractionServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"features": []map[string]interface{}{
				{"properties": map[string]interface{}{"name": "Test Attraction", "kinds": "historic,architecture", "dist": 100.0, "xid": "test-xid"}},
			},
		})
	}))
	defer attractionServer.Close()

	weatherServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"location": map[string]interface{}{"name": "Testville"},
			"current": map[string]interface{}{
				"temp_c":      25.0,
				"temp_f":      77.0,
				"feelslike_c": 26.0,
				"humidity":    50,
				"wind_kph":    10.0,
				"condition": map[string]interface{}{
					"text": "Sunny",
					"icon": "//cdn.example.com/icon.png",
				},
			},
			"forecast": map[string]interface{}{
				"forecastday": []map[string]interface{}{
					{
						"date": "2026-01-01",
						"day": map[string]interface{}{
							"maxtemp_c":   28.0,
							"mintemp_c":   18.0,
							"maxtemp_f":   82.0,
							"mintemp_f":   64.0,
							"avghumidity": 40.0,
							"condition": map[string]interface{}{
								"text": "Clear",
								"icon": "//cdn.example.com/tomorrow.png",
							},
						},
					},
				},
			},
		})
	}))
	defer weatherServer.Close()

	t.Setenv("RESTCOUNTRIES_BASE_URL", countryServer.URL)
	t.Setenv("OPENTRIPMAP_BASE_URL", attractionServer.URL)
	t.Setenv("OPENTRIPMAP_API_KEY", "test-key")
	t.Setenv("WEATHERAPI_BASE_URL", weatherServer.URL)
	t.Setenv("WEATHERAPI_KEY", "test-key")

	ctrl, _, _ := newCountryController(t, http.MethodGet, "/countries/testland")
	ctrl.Ctx.Input.SetParam(":slug", "testland")
	ctrl.Ctx.Request.URL.Path = "/countries/testland"
	ctrl.Detail()

	if got := ctrl.TplName; got != "destination.tpl" {
		t.Errorf("expected destination.tpl, got %q", got)
	}
	if got := ctrl.Data["ActiveNav"].(string); got != "countries" {
		t.Errorf("expected ActiveNav countries, got %q", got)
	}
	if ctrl.Data["Country"] == nil {
		t.Fatal("expected Country to be set")
	}
	if got := ctrl.Data["Attractions"]; got == nil {
		t.Fatal("expected Attractions to be set")
	}
	if got := ctrl.Data["Weather"]; got == nil {
		t.Fatal("expected Weather to be set")
	}
}
