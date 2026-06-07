// DashboardAPIController handles GET /api/dashboard/summary

package apicontrollers

import beego "github.com/beego/beego/v2/server/web"

type DashboardAPIController struct{
	beego.Controller
}

// Get returns wishlist counts for the dashboard stats panel (stub)
func (c *DashboardAPIController) Get() {
	c.Data["json"] = map[string]string{"status": "ok", "message": "dashboard stub"}
	c.ServeJSON()
}