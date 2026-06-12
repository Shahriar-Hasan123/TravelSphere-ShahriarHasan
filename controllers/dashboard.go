// DashboardController serves GET /dashboard — protected by AuthFilter.
package controllers

import "TravelSphere/services"

type DashboardController struct {
	BaseController
}

// Get renders the dashboard with live stats and the saved destinations list.
func (c *DashboardController) Get() {
	username := c.GetSession("username").(string)
	svc := services.NewDashboardService()

	total, planned, visited := svc.Summary(username)
	items := svc.GetItems(username)

	c.Data["ActiveNav"] = "dashboard"
	c.Data["TotalSaved"] = total
	c.Data["Planned"] = planned
	c.Data["Visited"] = visited
	c.Data["WishlistItems"] = items
	c.TplName = "dashboard.tpl"
	c.Layout = "layout.tpl"
}
