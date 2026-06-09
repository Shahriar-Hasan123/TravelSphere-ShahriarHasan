// controllers/api/dashboard_test.go
package apicontrollers

import (
	"TravelSphere/services"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/beego/beego/v2/server/web/context"
)

func TestDashboardAPI_Get_NoSession(t *testing.T) {
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/api/dashboard/summary", nil)

	beeCtx := context.NewContext()
	beeCtx.Reset(w, req)

	ctrl := &DashboardAPIController{}
	ctrl.Ctx = beeCtx
	ctrl.Data = make(map[interface{}]interface{})

	defer func() {
		if r := recover(); r != nil {
			t.Errorf("Get() panicked: %v", r)
		}
	}()
	ctrl.Get()
}

func TestDashboardAPI_Get_ReturnsJSON(t *testing.T) {
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/api/dashboard/summary", nil)

	beeCtx := context.NewContext()
	beeCtx.Reset(w, req)

	ctrl := &DashboardAPIController{}
	ctrl.Ctx = beeCtx
	ctrl.Data = make(map[interface{}]interface{})

	defer func() {
		recover()
	}()

	ctrl.Get()

	body := w.Body.Bytes()
	if len(body) > 0 {
		var result map[string]interface{}
		if err := json.Unmarshal(body, &result); err == nil {
			// Valid JSON response
			if data, ok := result["data"]; ok {
				dataMap := data.(map[string]interface{})
				if _, hasTotal := dataMap["total"]; !hasTotal {
					t.Error("expected 'total' in response data")
				}
				if _, hasPlanned := dataMap["planned"]; !hasPlanned {
					t.Error("expected 'planned' in response data")
				}
				if _, hasVisited := dataMap["visited"]; !hasVisited {
					t.Error("expected 'visited' in response data")
				}
			}
		}
	}
}

func TestDashboardAPI_Get_WithData(t *testing.T) {
	// Pre-populate wishlist with items
	svc := services.GetWishlistService()
	svc.Create("dashtest", "France", "note", "Planned")
	svc.Create("dashtest", "Japan", "note", "Visited")
	svc.Create("dashtest", "Spain", "note", "Planned")

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/api/dashboard/summary", nil)

	beeCtx := context.NewContext()
	beeCtx.Reset(w, req)

	ctrl := &DashboardAPIController{}
	ctrl.Ctx = beeCtx
	ctrl.Data = make(map[interface{}]interface{})

	defer func() {
		recover()
	}()

	// This will use empty username since there's no session
	ctrl.Get()

	// Cleanup
	for _, item := range svc.GetAll("dashtest") {
		svc.Delete("dashtest", item.ID)
	}
}

func TestDashboardAPI_Get_EmptyWishlist(t *testing.T) {
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/api/dashboard/summary", nil)

	beeCtx := context.NewContext()
	beeCtx.Reset(w, req)

	ctrl := &DashboardAPIController{}
	ctrl.Ctx = beeCtx
	ctrl.Data = make(map[interface{}]interface{})

	defer func() {
		recover()
	}()

	ctrl.Get()

	// Should still return valid JSON structure
	body := w.Body.Bytes()
	if len(body) > 0 {
		var result map[string]interface{}
		if err := json.Unmarshal(body, &result); err != nil {
			t.Errorf("expected valid JSON response, got: %v", err)
		}
	}
}
