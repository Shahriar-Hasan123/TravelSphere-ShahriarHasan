package controllers

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/beego/beego/v2/server/web"
	beecontext "github.com/beego/beego/v2/server/web/context"
	"github.com/beego/beego/v2/server/web/session"
)

type mockSession struct {
	data map[interface{}]interface{}
}

func (m *mockSession) Set(ctx context.Context, key, val interface{}) error {
	m.data[key] = val
	return nil
}

func (m *mockSession) Get(ctx context.Context, key interface{}) interface{} {
	return m.data[key]
}

func (m *mockSession) Delete(ctx context.Context, key interface{}) error {
	delete(m.data, key)
	return nil
}

func (m *mockSession) SessionID(ctx context.Context) string { return "mock-session-id" }

func (m *mockSession) SessionRelease(ctx context.Context, w http.ResponseWriter) {}

func (m *mockSession) SessionReleaseIfPresent(ctx context.Context, w http.ResponseWriter) {}

func (m *mockSession) Flush(ctx context.Context) error {
	m.data = make(map[interface{}]interface{})
	return nil
}

func ensureBeegoGlobalSession(t *testing.T) {
	t.Helper()
	if web.GlobalSessions != nil {
		return
	}
	mgr, err := session.NewManager("memory", session.NewManagerConfig())
	if err != nil {
		t.Fatalf("unable to create global session manager: %v", err)
	}
	web.GlobalSessions = mgr
}

func newControllerContext(t *testing.T, method, path string) (*beecontext.Context, *httptest.ResponseRecorder) {
	t.Helper()
	w := httptest.NewRecorder()
	req := httptest.NewRequest(method, path, nil)
	ctx := beecontext.NewContext()
	ctx.Reset(w, req)
	ctx.Input.CruSession = &mockSession{data: make(map[interface{}]interface{})}
	return ctx, w
}

func newBaseController(t *testing.T) (*BaseController, *beecontext.Context, *httptest.ResponseRecorder) {
	ctx, w := newControllerContext(t, http.MethodGet, "/")
	ctrl := &BaseController{}
	ctrl.Ctx = ctx
	ctrl.Data = make(map[interface{}]interface{})
	return ctrl, ctx, w
}

func TestBaseController_Prepare_DefaultsForAnonymousUsers(t *testing.T) {
	ctrl, _, _ := newBaseController(t)
	ctrl.Prepare()

	if got := ctrl.Data["isLoggedIn"].(bool); got {
		t.Errorf("expected isLoggedIn false, got %v", got)
	}
	if got := ctrl.Data["IsLoggedIn"].(bool); got {
		t.Errorf("expected IsLoggedIn false, got %v", got)
	}
	if got := ctrl.Data["username"].(string); got != "" {
		t.Errorf("expected empty username, got %q", got)
	}
	if got := ctrl.Data["ActiveNav"].(string); got != "" {
		t.Errorf("expected empty ActiveNav, got %q", got)
	}
}

func TestBaseController_Prepare_UsesSessionUsername(t *testing.T) {
	ctrl, ctx, _ := newBaseController(t)
	ctx.Input.CruSession.Set(context.Background(), "username", "traveler")

	ctrl.Prepare()

	if got := ctrl.Data["username"].(string); got != "traveler" {
		t.Errorf("expected username traveler, got %q", got)
	}
	if got := ctrl.Data["isLoggedIn"].(bool); !got {
		t.Errorf("expected isLoggedIn true, got %v", got)
	}
}

func TestBaseController_RequireLogin_RedirectsWhenUnauthenticated(t *testing.T) {
	ctrl, _, w := newBaseController(t)
	ok := ctrl.RequireLogin()

	if ok {
		t.Error("expected RequireLogin to return false for unauthenticated users")
	}
	if w.Result().StatusCode != http.StatusFound {
		t.Errorf("expected redirect status 302, got %d", w.Result().StatusCode)
	}
	if got := w.Header().Get("Location"); got != "/login" {
		t.Errorf("expected Location /login, got %q", got)
	}
}

func TestBaseController_RequireLogin_AllowsWhenAuthenticated(t *testing.T) {
	ctrl, _, _ := newBaseController(t)
	ctrl.Data["isLoggedIn"] = true

	if ok := ctrl.RequireLogin(); !ok {
		t.Error("expected RequireLogin to return true for authenticated users")
	}
}
