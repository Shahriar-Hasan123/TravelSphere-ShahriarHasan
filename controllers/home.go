// HomeController handles GET / — the home page.
package controllers

type HomeController struct {
	BaseController
}

// Get renders the home page with featured countries and popular attractions.
func (c *HomeController) Get() {
	c.Data["ActiveNav"] = "home"
	c.TplName = "home.tpl"
	c.Layout = "layout.tpl"
}
