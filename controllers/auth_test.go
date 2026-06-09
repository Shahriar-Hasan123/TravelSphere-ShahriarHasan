package controllers

import (
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	beecontext "github.com/beego/beego/v2/server/web/context"
)

func newAuthController(t *testing.T, method, path string, body string, cookies []*http.Cookie) (*AuthController, *beecontext.Context, *httptest.ResponseRecorder) {
	t.Helper()
	var reader io.Reader
	if body != "" {
		reader = strings.NewReader(body)
	}
	w := httptest.NewRecorder()
	req := httptest.NewRequest(method, path, reader)
	if body != "" {
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	}
	for _, cookie := range cookies {
		req.AddCookie(cookie)
	}
	ctx := beecontext.NewContext()
	ctx.Reset(w, req)
	ensureBeegoGlobalSession(t)
	ctx.Input.CruSession = &mockSession{data: make(map[interface{}]interface{})}

	ctrl := &AuthController{}
	ctrl.Ctx = ctx
	ctrl.Data = make(map[interface{}]interface{})
	return ctrl, ctx, w
}

func TestAuthController_ShowLogin_RendersLoginPage(t *testing.T) {
	cookie := &http.Cookie{Name: "redirect_after_login", Value: "/wishlist"}
	ctrl, _, _ := newAuthController(t, http.MethodGet, "/login", "", []*http.Cookie{cookie})
	ctrl.Data["IsLoggedIn"] = false

	ctrl.ShowLogin()

	if got := ctrl.TplName; got != "login.tpl" {
		t.Errorf("expected login.tpl, got %q", got)
	}
	if got := ctrl.Data["RedirectTo"].(string); got != "/wishlist" {
		t.Errorf("expected RedirectTo /wishlist, got %q", got)
	}
}

func TestAuthController_ShowLogin_RedirectsWhenAuthenticated(t *testing.T) {
	ctrl, _, w := newAuthController(t, http.MethodGet, "/login", "", nil)
	ctrl.Data["IsLoggedIn"] = true

	ctrl.ShowLogin()

	if w.Result().StatusCode != http.StatusFound {
		t.Errorf("expected redirect status 302, got %d", w.Result().StatusCode)
	}
	if got := w.Header().Get("Location"); got != "/dashboard" {
		t.Errorf("expected Location /dashboard, got %q", got)
	}
}

func TestAuthController_DoLogin_ShowsErrorForMissingUsername(t *testing.T) {
	body := url.Values{"username": {""}, "redirect_to": {""}}.Encode()
	ctrl, _, _ := newAuthController(t, http.MethodPost, "/login", body, nil)

	ctrl.DoLogin()

	if got := ctrl.TplName; got != "login.tpl" {
		t.Errorf("expected login.tpl, got %q", got)
	}
	if got := ctrl.Data["Error"].(string); got != "Please enter a username." {
		t.Errorf("expected validation error, got %q", got)
	}
}

func TestAuthController_DoLogin_RedirectsToRequestedPage(t *testing.T) {
	body := url.Values{"username": {"traveler"}, "redirect_to": {"/dashboard"}}.Encode()
	ctrl, ctx, w := newAuthController(t, http.MethodPost, "/login", body, nil)

	ctrl.DoLogin()

	if got := w.Result().StatusCode; got != http.StatusFound {
		t.Errorf("expected redirect status 302, got %d", got)
	}
	if got := w.Header().Get("Location"); got != "/dashboard" {
		t.Errorf("expected Location /dashboard, got %q", got)
	}
	if got := ctx.Input.Session("username"); got != "traveler" {
		t.Errorf("expected session username traveler, got %v", got)
	}
}

func TestAuthController_DoLogout_ClearsSession(t *testing.T) {
	ctrl, ctx, w := newAuthController(t, http.MethodPost, "/logout", "", nil)
	ctx.Input.CruSession.Set(context.Background(), "username", "traveler")

	ctrl.DoLogout()

	if got := w.Result().StatusCode; got != http.StatusFound {
		t.Errorf("expected redirect status 302, got %d", got)
	}
	if got := w.Header().Get("Location"); got != "/" {
		t.Errorf("expected Location /, got %q", got)
	}
	if ctx.Input.CruSession != nil {
		t.Errorf("expected session to be destroyed, but CruSession was still present")
	}
}

func TestSessionUsername_ReadsSessionValue(t *testing.T) {
	_, ctx, _ := newAuthController(t, http.MethodGet, "/login", "", nil)
	ctx.Input.CruSession.Set(context.Background(), "username", "traveler")

	if got := SessionUsername(ctx); got != "traveler" {
		t.Errorf("expected SessionUsername traveler, got %q", got)
	}
}

func TestSessionUsername_ReturnsEmptyWhenMissing(t *testing.T) {
	_, ctx, _ := newAuthController(t, http.MethodGet, "/login", "", nil)

	if got := SessionUsername(ctx); got != "" {
		t.Errorf("expected empty username, got %q", got)
	}
}
