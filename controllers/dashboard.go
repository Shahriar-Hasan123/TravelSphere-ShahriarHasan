// DashboardController handles GET /dashboard - protected SSR page

package controllers

type DashboardController struct {
	BaseController
}

// Get renders the dashboard page for authenticated users
func (c *DashboardController) Get() {
	if !c.RequireLogin() {
		return
	}
	c.Data["ActiveNav"] = "dashboard"
	c.TplName = "dashboard.tpl"
	c.Layout = "layout.tpl"
}
