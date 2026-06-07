// CountriesAPIController handles JSON API routes for country data.

package apicontrollers

import (
	beego "github.com/beego/beego/v2/server/web"
)

type CountriesAPIController struct {
	beego.Controller
}

// Get returns a JSON list of countries
func (c *CountriesAPIController) Get() {
	c.Data["json"] = map[string]string{"status": "ok", "message": "countries stub"}
	c.ServeJSON()
}

// Detail returns JSON detail for a single country by slug (stub)
func (c *CountriesAPIController) Detail() {
	c.Data["json"] = map[string]string{"status": "ok", "message": "country detail stub"}
	c.ServeJSON()
}
