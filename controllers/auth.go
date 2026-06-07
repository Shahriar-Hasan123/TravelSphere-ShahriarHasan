// AuthController handles login and logout SSR routes

package controllers

type AuthController struct {
	BaseController
}

// ShowLogin renders the login page (GET /login)
func (c *AuthController) ShowLogin() {
	c.Data["ActiveNav"] = ""
	c.TplName = "login.tpl"
	c.Layout = "layout.tpl"
}

// DoLogin processes the login form (POST /login)
func (c *AuthController) DoLogin() {
	username := c.GetString("username")
	password := c.GetString("password")

	if username != "" && password != "" {
		c.SetSession("username", username)
		c.Redirect("/dashboard", 302)
		return
	}

	c.Data["Error"] = "Invalid credentials"
	c.Data["ActiveNav"] = ""
	c.TplName = "login.tpl"
	c.Layout = "layout.tpl"
}

// DoLogout clears the session and redirects to home (GET /logout)
func (c *AuthController) DoLogout() {
	c.DestroySession()
	c.Redirect("/", 302)
}
