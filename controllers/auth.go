// controllers/auth.go
package controllers

import "github.com/beego/beego/v2/server/web/context"

type AuthController struct {
	BaseController
}

// ShowLogin renders the login page. Redirects to dashboard if already authenticated.
func (c *AuthController) ShowLogin() {
	if c.Data["IsLoggedIn"].(bool) {
		c.Redirect("/dashboard", 302)
		return
	}
	c.Data["RedirectTo"] = c.Ctx.GetCookie("redirect_after_login")
	c.Data["ActiveNav"] = ""
	c.TplName = "login.tpl"
	c.Layout = "layout.tpl"
}

// DoLogin creates a session for any non-empty username.
func (c *AuthController) DoLogin() {
	username := c.GetString("username")
	redirectTo := c.GetString("redirect_to")

	if username == "" {
		c.Data["Error"] = "Please enter a username."
		c.Data["ActiveNav"] = ""
		c.TplName = "login.tpl"
		c.Layout = "layout.tpl"
		return
	}

	c.SetSession("username", username)
	c.Ctx.SetCookie("redirect_after_login", "", -1)

	if redirectTo != "" {
		c.Redirect(redirectTo, 302)
		return
	}
	c.Redirect("/dashboard", 302)
}

// DoLogout destroys the session and returns the user to the home page.
func (c *AuthController) DoLogout() {
	c.DestroySession()
	c.Ctx.SetCookie("redirect_after_login", "", -1)
	c.Redirect("/", 302)
}

// SessionUsername reads the authenticated username directly from a filter context.
func SessionUsername(ctx *context.Context) string {
	val := ctx.Input.Session("username")
	if val == nil {
		return ""
	}
	username, _ := val.(string)
	return username
}
