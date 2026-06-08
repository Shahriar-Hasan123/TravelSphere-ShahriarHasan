// DashboardAPIController serves GET /api/dashboard/summary for AJAX refresh.
package apicontrollers

import (
	"TravelSphere/services"
	"TravelSphere/utils"

	beego "github.com/beego/beego/v2/server/web"
)

type DashboardAPIController struct {
	beego.Controller
}

// Get returns wishlist summary counts for the authenticated user.
func (c *DashboardAPIController) Get() {
	username, _ := c.GetSession("username").(string)

	total, planned, visited := services.NewDashboardService().Summary(username)

	c.Data["json"] = utils.OKResponse(map[string]int{
		"total":   total,
		"planned": planned,
		"visited": visited,
	})
	c.ServeJSON()
}
