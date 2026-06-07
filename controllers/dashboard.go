// DashboardController handles GET /dashboard — protected SSR page.
package controllers

type DashboardController struct {
	BaseController
}

// Get renders the dashboard page. Auth is guaranteed by the filter — no need to check here.
func (c *DashboardController) Get() {
	c.Data["ActiveNav"] = "dashboard"
	c.Data["TotalSaved"] = 0  
	c.Data["Planned"]    = 0
	c.Data["Visited"]    = 0
	c.Data["WishlistItems"] = nil
	c.TplName = "dashboard.tpl"
	c.Layout = "layout.tpl"
}