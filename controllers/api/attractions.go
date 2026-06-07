// AttractionsAPIController handles GET /api/attractions.

package apicontrollers

import beego "github.com/beego/beego/v2/server/web"

type AttractionsAPIController struct {
	beego.Controller
}

// Get returns attractions by country/coordinates (stub)
func (c *AttractionsAPIController) Get() {
	c.Data["json"] = map[string]string{"status": "ok", "message": "attractions stub"}
	c.ServeJSON()
}
