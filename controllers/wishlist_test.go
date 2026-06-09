package controllers

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"TravelSphere/models"
	"TravelSphere/services"

	"github.com/beego/beego/v2/server/web/context"
)

func newWishlistController(t *testing.T) (*WishlistController, *context.Context, *httptest.ResponseRecorder) {
	t.Helper()
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/wishlist", nil)
	ctx := context.NewContext()
	ctx.Reset(w, req)
	ctx.Input.CruSession = &mockSession{data: map[interface{}]interface{}{"username": "wishlistuser"}}
	ctrl := &WishlistController{}
	ctrl.Ctx = ctx
	ctrl.Data = make(map[interface{}]interface{})
	return ctrl, ctx, w
}

func TestWishlistController_Get_PopulatesWishlistItems(t *testing.T) {
	svc := services.GetWishlistService()
	item := svc.Create("wishlistuser", "Italy", "notes", "Planned")
	defer svc.Delete("wishlistuser", item.ID)

	ctrl, _, _ := newWishlistController(t)
	ctrl.Get()

	if got := ctrl.Data["ActiveNav"].(string); got != "wishlist" {
		t.Errorf("expected ActiveNav wishlist, got %q", got)
	}
	items, ok := ctrl.Data["WishlistItems"].([]models.WishlistItem)
	if !ok {
		t.Fatal("expected WishlistItems to be a []WishlistItem")
	}
	if len(items) != 1 {
		t.Fatalf("expected 1 wishlist item, got %d", len(items))
	}
	if got := ctrl.TplName; got != "wishlist.tpl" {
		t.Errorf("expected wishlist.tpl, got %q", got)
	}
}

func TestWishlistController_Get_EmptyWishlist(t *testing.T) {
	ctrl, _, _ := newWishlistController(t)
	ctrl.Get()

	items, ok := ctrl.Data["WishlistItems"].([]models.WishlistItem)
	if !ok {
		t.Fatal("expected WishlistItems to be a []WishlistItem")
	}
	if len(items) != 0 {
		t.Errorf("expected empty wishlist, got %d", len(items))
	}
}
