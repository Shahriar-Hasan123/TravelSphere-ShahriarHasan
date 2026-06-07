package controllers

import (
	beego "github.com/beego/beego/v2/server/web"
)

type BaseController struct {
	beego.Controller
}

// Prepare runs automatically before Get() or Post() on every request.
// It reads the session and sets template variables available to all pages.
func (c *BaseController) Prepare() {
	username, _ := c.GetSession("username").(string)
	isLoggedIn := username != ""
	c.Data["isLoggedIn"] = isLoggedIn
	c.Data["username"] = username
	if _, exists := c.Data["ActiveNav"]; !exists {
		c.Data["ActiveNav"] = ""
	}
}

// RequireLogin redirects unauthenticated users to the login page.
func (c *BaseController) RequireLogin() bool {
	if !c.Data["isLoggedIn"].(bool) {
		c.Redirect("/login", 302)
		return false
	}
	return true
}
