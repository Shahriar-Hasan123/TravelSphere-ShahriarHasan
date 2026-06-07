// AuthController handles the login and logout SSR routes. Authentication uses session storage — no database is involved.

package controllers

import "github.com/beego/beego/v2/server/web/context"

type AuthController struct {
	BaseController
}

// validCredentials defines the accepted username/password pairs.
var validCredentials = map[string]string{
	"beta": "beta123",
}

// ShowLogin renders the login page (GET /login).
func (c *AuthController) ShowLogin() {
	isLoggedIn, _ := c.Data["isLoggedIn"].(bool)
	if isLoggedIn {
		c.Redirect("/dashboard", 302)
		return
	}

	// Check if there is a pending redirect target from the auth filter cookie.
	redirectTo := c.Ctx.GetCookie("redirect_after_login")

	c.Data["RedirectTo"] = redirectTo
	c.Data["ActiveNav"] = ""
	c.TplName = "login.tpl"
	c.Layout = "layout.tpl"
}

// DoLogin processes POST /login, sets session on success, and redirects or shows an error.
func (c *AuthController) DoLogin() {
	username := c.GetString("username")
	password := c.GetString("password")
	redirectTo := c.GetString("redirect_to")

	// Validate credentials against the hardcoded store.
	expectedPassword, exists := validCredentials[username]
	if !exists || expectedPassword != password {
		c.Data["Error"] = "Invalid username or password."
		c.Data["Username"] = username // Pre-fill username field on error.
		c.Data["ActiveNav"] = ""
		c.TplName = "login.tpl"
		c.Layout = "layout.tpl"
		return
	}

	// Credentials valid — establish session.
	c.SetSession("username", username)

	// Clear the redirect cookie now that login succeeded.
	c.Ctx.SetCookie("redirect_after_login", "", -1)

	// Redirect to the originally requested page or default to dashboard.
	if redirectTo != "" {
		c.Redirect(redirectTo, 302)
		return
	}
	c.Redirect("/dashboard", 302)
}

// DoLogout destroys the session and redirects to the home page (GET /logout).
func (c *AuthController) DoLogout() {
	c.DestroySession()

	// Clear any leftover redirect cookie.
	c.Ctx.SetCookie("redirect_after_login", "", -1)

	c.Redirect("/", 302)
}

// SessionUsername is a helper used by filters to read the session username
func SessionUsername(ctx *context.Context) string {
	val := ctx.Input.Session("username")
	if val == nil {
		return ""
	}
	username, _ := val.(string)
	return username
}
