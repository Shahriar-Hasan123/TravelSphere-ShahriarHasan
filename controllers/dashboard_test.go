package controllers

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"TravelSphere/models"
	"TravelSphere/services"

	"github.com/beego/beego/v2/server/web/context"
)

func newDashboardController(t *testing.T) (*DashboardController, *context.Context, *httptest.ResponseRecorder) {
	t.Helper()
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/dashboard", nil)
	ctx := context.NewContext()
	ctx.Reset(w, req)
	ctx.Input.CruSession = &mockSession{data: map[interface{}]interface{}{"username": "dashuser"}}
	ctrl := &DashboardController{}
	ctrl.Ctx = ctx
	ctrl.Data = make(map[interface{}]interface{})
	return ctrl, ctx, w
}

func TestDashboardController_Get_PopulatesDashboardData(t *testing.T) {
	svc := services.GetWishlistService()
	item1 := svc.Create("dashuser", "Canada", "note", "Planned")
	item2 := svc.Create("dashuser", "Spain", "note", "Visited")
	defer func() {
		svc.Delete("dashuser", item1.ID)
		svc.Delete("dashuser", item2.ID)
	}()

	ctrl, _, _ := newDashboardController(t)
	ctrl.Get()

	if got := ctrl.Data["ActiveNav"].(string); got != "dashboard" {
		t.Errorf("expected ActiveNav dashboard, got %q", got)
	}
	if got := ctrl.Data["TotalSaved"].(int); got != 2 {
		t.Errorf("expected TotalSaved 2, got %d", got)
	}
	if got := ctrl.Data["Planned"].(int); got != 1 {
		t.Errorf("expected Planned 1, got %d", got)
	}
	if got := ctrl.Data["Visited"].(int); got != 1 {
		t.Errorf("expected Visited 1, got %d", got)
	}
	if items, ok := ctrl.Data["WishlistItems"].([]models.WishlistItem); ok {
		if len(items) != 2 {
			t.Errorf("expected 2 wishlist items, got %d", len(items))
		}
	} else {
		t.Fatal("expected WishlistItems to be a slice of WishlistItem")
	}
	if got := ctrl.TplName; got != "dashboard.tpl" {
		t.Errorf("expected dashboard.tpl, got %q", got)
	}
}
