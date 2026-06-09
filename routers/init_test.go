package routers

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/beego/beego/v2/server/web/context"
)

func TestInit_RegistersRoutesWithoutPanic(t *testing.T) {
	defer func() {
		if r := recover(); r != nil {
			t.Fatalf("Init() panicked: %v", r)
		}
	}()
	Init()
}

func TestRegisterRoutesHelpers_DoNotPanic(t *testing.T) {
	defer func() {
		if r := recover(); r != nil {
			t.Fatalf("route registration panicked: %v", r)
		}
	}()

	registerSSRRoutes()
	registerAPIRoutes()
}

func TestRegisterFilters_DoesNotPanic(t *testing.T) {
	defer func() {
		if r := recover(); r != nil {
			t.Fatalf("registerFilters panicked: %v", r)
		}
	}()

	registerFilters()
}

func TestLoggingFilterHandler_DoesNotPanic(t *testing.T) {
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	beeCtx := context.NewContext()
	beeCtx.Reset(w, req)

	defer func() {
		if r := recover(); r != nil {
			t.Fatalf("loggingFilterHandler panicked: %v", r)
		}
	}()

	loggingFilterHandler(beeCtx)
}

func TestAuthFilterHandler_ForAPIRequestWritesUnauthorized(t *testing.T) {
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/api/wishlist", nil)
	beeCtx := context.NewContext()
	beeCtx.Reset(w, req)

	defer func() {
		if r := recover(); r != nil {
			t.Fatalf("authFilterHandler panicked: %v", r)
		}
	}()

	authFilterHandler(beeCtx)
	if w.Code != http.StatusUnauthorized {
		t.Fatalf("expected 401 unauthorized, got %d", w.Code)
	}
}
