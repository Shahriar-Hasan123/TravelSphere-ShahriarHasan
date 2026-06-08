// AttractionsAPIController handles GET /api/attractions.
// Returns attraction data by lat/lon coordinates — used by AJAX on the destination detail page if needed, and available for unit tests.

package apicontrollers

import (
	"TravelSphere/services"
	"TravelSphere/utils"
	"strconv"

	beego "github.com/beego/beego/v2/server/web"
)

// AttractionsAPIController returns JSON attraction data by coordinates.
type AttractionsAPIController struct {
	beego.Controller
}

var attractionSvc = services.NewAttractionService()

// Get returns attractions near the given lat/lon coordinates.
// Query params: lat (float), lon (float).
func (c *AttractionsAPIController) Get() {
	latStr := c.GetString("lat")
	lonStr := c.GetString("lon")

	if latStr == "" || lonStr == "" {
		c.Ctx.Output.SetStatus(400)
		c.Data["json"] = utils.ErrorResponse("lat and lon query params are required", 400)
		c.ServeJSON()
		return
	}

	lat, err := strconv.ParseFloat(latStr, 64)
	if err != nil {
		c.Ctx.Output.SetStatus(400)
		c.Data["json"] = utils.ErrorResponse("invalid lat value", 400)
		c.ServeJSON()
		return
	}

	lon, err := strconv.ParseFloat(lonStr, 64)
	if err != nil {
		c.Ctx.Output.SetStatus(400)
		c.Data["json"] = utils.ErrorResponse("invalid lon value", 400)
		c.ServeJSON()
		return
	}

	attractions, err := attractionSvc.GetAttractionsByCoords(lat, lon)
	if err != nil {
		c.Ctx.Output.SetStatus(500)
		c.Data["json"] = utils.ErrorResponse("Failed to fetch attractions", 500)
		c.ServeJSON()
		return
	}

	c.Data["json"] = utils.OKResponse(attractions)
	c.ServeJSON()
}
