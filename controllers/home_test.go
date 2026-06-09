package controllers

import (
	"TravelSphere/models"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/beego/beego/v2/server/web/context"
)

func newHomeController(t *testing.T) (*HomeController, *http.Request, *httptest.ResponseRecorder) {
	t.Helper()
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	ctx := context.NewContext()
	ctx.Reset(w, req)
	ctx.Input.CruSession = &mockSession{data: make(map[interface{}]interface{})}
	ctrl := &HomeController{}
	ctrl.Ctx = ctx
	ctrl.Data = make(map[interface{}]interface{})
	return ctrl, req, w
}

func TestHomeController_Get_FallsBackWhenCountryServiceFails(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer server.Close()
	t.Setenv("RESTCOUNTRIES_BASE_URL", server.URL)
	t.Setenv("OPENTRIPMAP_BASE_URL", server.URL)
	t.Setenv("OPENTRIPMAP_API_KEY", "test-key")

	ctrl, _, _ := newHomeController(t)
	ctrl.Get()

	if featured, ok := ctrl.Data["FeaturedCountries"].([]models.FeaturedCountry); ok {
		if len(featured) != 0 {
			t.Errorf("expected FeaturedCountries empty when service fails, got %d", len(featured))
		}
	} else if ctrl.Data["FeaturedCountries"] != nil {
		t.Error("expected FeaturedCountries nil or empty when service fails")
	}
	if attractions, ok := ctrl.Data["PopularAttractions"].([]models.Attraction); !ok || len(attractions) == 0 {
		t.Fatal("expected fallback PopularAttractions to be populated")
	}
	if got := ctrl.Data["ActiveNav"].(string); got != "home" {
		t.Errorf("expected ActiveNav home, got %q", got)
	}
	if got := ctrl.TplName; got != "home.tpl" {
		t.Errorf("expected home.tpl, got %q", got)
	}
}

func TestFetchHomeAttractions_ReturnsStaticFallbackWhenExternalFails(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer server.Close()
	t.Setenv("OPENTRIPMAP_BASE_URL", server.URL)
	t.Setenv("OPENTRIPMAP_API_KEY", "test-key")

	attractions := fetchHomeAttractions()
	if len(attractions) == 0 {
		t.Fatal("expected fallback attractions, got empty list")
	}
	if attractions[0].Name != "Eiffel Tower" {
		t.Errorf("expected first fallback attraction to be Eiffel Tower, got %q", attractions[0].Name)
	}
}

func TestFetchHomeAttractions_UsesOpenTripMapWhenAvailable(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"features": []map[string]interface{}{
				{"properties": map[string]interface{}{"name": "Home Attraction", "kinds": "historic,architecture", "dist": 1.0, "xid": "home-xid"}},
			},
		})
	}))
	defer server.Close()
	t.Setenv("OPENTRIPMAP_BASE_URL", server.URL)
	t.Setenv("OPENTRIPMAP_API_KEY", "test-key")

	attractions := fetchHomeAttractions()
	if len(attractions) == 0 {
		t.Fatal("expected attractions from OpenTripMap, got empty list")
	}
	if attractions[0].Name != "Home Attraction" {
		t.Errorf("expected Home Attraction, got %q", attractions[0].Name)
	}
}
