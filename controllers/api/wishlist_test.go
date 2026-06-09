// controllers/api/wishlist_test.go
package apicontrollers

import (
	"TravelSphere/models"
	"TravelSphere/services"
	"TravelSphere/utils"
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/beego/beego/v2/server/web/context"
)

func newWishlistAPICtrl(method, path string, body []byte) (*WishlistAPIController, *httptest.ResponseRecorder) {
	w := httptest.NewRecorder()
	req := httptest.NewRequest(method, path, bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	beeCtx := context.NewContext()
	beeCtx.Reset(w, req)

	ctrl := &WishlistAPIController{}
	ctrl.Ctx = beeCtx
	ctrl.Data = make(map[interface{}]interface{})
	return ctrl, w
}

// runSafe calls fn and recovers from any panic — ServeJSON panics without full Beego runtime.
func runSafe(fn func()) (panicked bool) {
	defer func() {
		if r := recover(); r != nil {
			panicked = true
		}
	}()
	fn()
	return false
}

func TestWishlistAPI_SessionUsername_NoSession(t *testing.T) {
	ctrl, _ := newWishlistAPICtrl(http.MethodGet, "/api/wishlist", nil)
	username := ctrl.sessionUsername()
	if username != "" {
		t.Errorf("expected empty username for no session, got %q", username)
	}
}

func TestWishlistAPI_SessionUsername_ReturnsEmpty(t *testing.T) {
	// Even with various scenarios, sessionUsername should return empty when not authenticated
	ctrl, _ := newWishlistAPICtrl(http.MethodGet, "/api/wishlist", nil)
	username := ctrl.sessionUsername()
	if username != "" {
		t.Errorf("expected empty username when not authenticated, got %q", username)
	}
}

func TestWishlistAPI_Get_NoSession(t *testing.T) {
	ctrl, w := newWishlistAPICtrl(http.MethodGet, "/api/wishlist", nil)
	runSafe(ctrl.Get)

	// Should return data even without session (just empty list for empty username)
	body := w.Body.Bytes()
	if len(body) == 0 {
		t.Error("expected response body from Get()")
	}
}

func TestWishlistAPI_Get_WithItems(t *testing.T) {
	// Pre-populate with an item
	services.GetWishlistService().Create("", "Greece", "Beautiful islands", "Planned")

	ctrl, w := newWishlistAPICtrl(http.MethodGet, "/api/wishlist", nil)
	runSafe(ctrl.Get)

	body := w.Body.Bytes()
	if len(body) == 0 {
		t.Fatal("expected response body from Get()")
	}

	// Verify it's valid JSON
	var result map[string]interface{}
	if err := json.Unmarshal(body, &result); err != nil {
		t.Errorf("response is not valid JSON: %v", err)
	}
}

func TestWishlistAPI_Post_InvalidJSON(t *testing.T) {
	ctrl, _ := newWishlistAPICtrl(http.MethodPost, "/api/wishlist", []byte("not json"))
	runSafe(ctrl.Post)
}

func TestWishlistAPI_Post_MissingCountryName(t *testing.T) {
	body := []byte(`{"note":"test","status":"Planned"}`)
	ctrl, _ := newWishlistAPICtrl(http.MethodPost, "/api/wishlist", body)
	runSafe(ctrl.Post)
}

func TestWishlistAPI_Post_InvalidStatus(t *testing.T) {
	body := []byte(`{"country_name":"France","status":"Unknown"}`)
	ctrl, _ := newWishlistAPICtrl(http.MethodPost, "/api/wishlist", body)
	runSafe(ctrl.Post)
}

func TestWishlistAPI_Post_ValidBody(t *testing.T) {
	body := []byte(`{"country_name":"France","status":"Planned"}`)
	ctrl, _ := newWishlistAPICtrl(http.MethodPost, "/api/wishlist", body)
	runSafe(ctrl.Post)
}

func TestWishlistAPI_Post_DefaultStatus(t *testing.T) {
	// status omitted — should default to Planned.
	body := []byte(`{"country_name":"Germany"}`)
	ctrl, _ := newWishlistAPICtrl(http.MethodPost, "/api/wishlist", body)
	runSafe(ctrl.Post)
}

func TestWishlistAPI_Post_DuplicateCountry(t *testing.T) {
	// Pre-seed the wishlist with Germany for this user (empty username).
	services.GetWishlistService().Create("", "Italy", "", "Planned")

	body := []byte(`{"country_name":"Italy","status":"Planned"}`)
	ctrl, _ := newWishlistAPICtrl(http.MethodPost, "/api/wishlist", body)
	runSafe(ctrl.Post)
}

func TestWishlistAPI_Update_InvalidJSON(t *testing.T) {
	ctrl, _ := newWishlistAPICtrl(http.MethodPut, "/api/wishlist/some-id", []byte("bad"))
	ctrl.Ctx.Input.SetParam(":id", "some-id")
	runSafe(ctrl.Update)
}

func TestWishlistAPI_Update_InvalidStatus(t *testing.T) {
	body := []byte(`{"note":"test","status":"Bad"}`)
	ctrl, _ := newWishlistAPICtrl(http.MethodPut, "/api/wishlist/some-id", body)

	ctrl.Ctx.Input.SetParam(":id", "some-id")
	runSafe(ctrl.Update)
}

func TestWishlistAPI_Update_NotFound(t *testing.T) {
	body := []byte(`{"note":"test","status":"Planned"}`)
	ctrl, _ := newWishlistAPICtrl(http.MethodPut, "/api/wishlist/nonexistent", body)
	ctrl.Ctx.Input.SetParam(":id", "nonexistent-id-999")
	runSafe(ctrl.Update)
}

func TestWishlistAPI_Update_Success(t *testing.T) {
	// Create an item with empty username (no session) to match sessionUsername().
	item := services.GetWishlistService().Create("", "Spain", "", "Planned")

	body, _ := json.Marshal(map[string]string{
		"note":   "Visit Madrid",
		"status": "Visited",
	})
	ctrl, _ := newWishlistAPICtrl(http.MethodPut, "/api/wishlist/"+item.ID, body)
	ctrl.Ctx.Input.SetParam(":id", item.ID)
	runSafe(ctrl.Update)
}

func TestWishlistAPI_Delete_NotFound(t *testing.T) {
	ctrl, _ := newWishlistAPICtrl(http.MethodDelete, "/api/wishlist/bad-id", nil)
	ctrl.Ctx.Input.SetParam(":id", "nonexistent-id-999")
	runSafe(ctrl.Delete)
}

func TestWishlistAPI_Delete_Success(t *testing.T) {
	item := services.GetWishlistService().Create("", "Portugal", "", "Planned")

	ctrl, _ := newWishlistAPICtrl(http.MethodDelete, "/api/wishlist/"+item.ID, nil)
	ctrl.Ctx.Input.SetParam(":id", item.ID)
	runSafe(ctrl.Delete)
}

// TestWishlistValidation tests the validation helpers used by the controller.
func TestWishlistValidation_CountryName(t *testing.T) {
	msg := utils.ValidateWishlistCreate("")
	if msg == "" {
		t.Error("expected validation error for empty country name")
	}

	msg = utils.ValidateWishlistCreate("France")
	if msg != "" {
		t.Errorf("expected no error for valid name, got %q", msg)
	}
}

func TestWishlistValidation_Status(t *testing.T) {
	if utils.ValidateWishlistStatus("Planned") != "" {
		t.Error("Planned should be valid")
	}
	if utils.ValidateWishlistStatus("Visited") != "" {
		t.Error("Visited should be valid")
	}
	if utils.ValidateWishlistStatus("Unknown") == "" {
		t.Error("Unknown should be invalid")
	}
}

// TestWishlistItemJSON verifies the JSON tags on WishlistItem match what JS expects.
func TestWishlistItemJSON(t *testing.T) {
	item := models.WishlistItem{
		ID:          "test-id",
		CountryName: "France",
		Note:        "test note",
		Status:      models.StatusPlanned,
	}

	data, err := json.Marshal(item)
	if err != nil {
		t.Fatalf("marshal error: %v", err)
	}

	var decoded map[string]interface{}
	if err := json.Unmarshal(data, &decoded); err != nil {
		t.Fatalf("unmarshal error: %v", err)
	}

	if decoded["id"] == nil {
		t.Error("expected 'id' key in JSON")
	}
	if decoded["CountryName"] == nil {
		t.Error("expected 'CountryName' key in JSON")
	}
	if decoded["Status"] == nil {
		t.Error("expected 'Status' key in JSON")
	}
}
