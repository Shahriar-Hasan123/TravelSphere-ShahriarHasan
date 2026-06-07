// CountryController handles SSR routes for /countries and /countries/:slug

package controllers

type CountryController struct {
	BaseController
}

// Get renders the Country Explorer page (/countries)
func (c *CountryController) Get() {
	c.Data["ActiveNav"] = "countries"
	c.TplName = "countries.tpl"
	c.Layout = "layout.tpl"
}

// Detail renders the Destination Detail page (/countries/:slug)
func (c *CountryController) Detail() {
	c.Data["ActiveNav"] = "countries"
	c.TplName = "destination.tpl"
	c.Layout = "layout.tpl"
}
