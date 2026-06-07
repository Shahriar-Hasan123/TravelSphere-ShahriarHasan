// AuthFilter protects routes that require an authenticated session.

package filters

import (
	"encoding/json"

	"github.com/beego/beego/v2/server/web/context"
)

// AuthFilter is a Beego filter function that enforces session authentication.
func AuthFilter(ctx *context.Context) {
	username := ctx.Input.Session("username")

	if username != nil && username.(string) != "" {
		return
	}

	// API routes always get JSON; SSR routes get a redirect.
	if isAPIRequest(ctx) {
		respondUnauthorized(ctx)
		return
	}

	// Store the originally requested URL so we can redirect back after login.
	ctx.SetCookie("redirect_after_login", ctx.Request.URL.Path, 300)
	ctx.Redirect(302, "/login")
}

// isAPIRequest returns true when the request path starts with /api/. These requests come from AJAX and expect JSON, not an HTML redirect.
func isAPIRequest(ctx *context.Context) bool {
	path := ctx.Request.URL.Path
	return len(path) >= 5 && path[:5] == "/api/"
}

// respondUnauthorized writes a 401 JSON error response for API requests.
func respondUnauthorized(ctx *context.Context) {
	ctx.Output.SetStatus(401)
	ctx.Output.Header("Content-Type", "application/json")
	body, _ := json.Marshal(map[string]interface{}{
		"status":  "error",
		"message": "Authentication required",
		"code":    401,
	})
	ctx.Output.Body(body)
}