// CountriesAPIController serves JSON responses consumed by AJAX calls.
// GET /api/countries - list with optional search and region filters
// GET /api/countries/:slug - single country detail

package apicontrollers

import (
	"TravelSphere/services"
	"TravelSphere/utils"

	beego "github.com/beego/beego/v2/server/web"
)

// CountriesAPIController handles all /api/countries JSON endpoints.
type CountriesAPIController struct {
	beego.Controller
}

// countrySvc is a package-level service instance (stateless, safe to share).
var countrySvc = services.NewCountryService()

// Get returns a filtered JSON array of countries.
// Query params: search (string), region (string).
// Called by countries.js AJAX on every search/filter change.
func (c *CountriesAPIController) Get() {
	search := c.GetString("search")
	region := c.GetString("region")

	countries, err := countrySvc.GetAllCountries(search, region)
	if err != nil {
		c.Ctx.Output.SetStatus(500)
		c.Data["json"] = utils.ErrorResponse("Failed to fetch countries", 500)
		c.ServeJSON()
		return
	}

	c.Data["json"] = utils.OKResponse(countries)
	c.ServeJSON()
}

// Detail returns JSON for a single country identified by its slug.
// Returns 404 when the slug is unknown.
func (c *CountriesAPIController) Detail() {
	slug := c.Ctx.Input.Param(":slug")

	country, err := countrySvc.GetCountryBySlug(slug)
	if err != nil {
		c.Ctx.Output.SetStatus(500)
		c.Data["json"] = utils.ErrorResponse("Failed to fetch country", 500)
		c.ServeJSON()
		return
	}

	if country == nil {
		c.Ctx.Output.SetStatus(404)
		c.Data["json"] = utils.ErrorResponse("Country not found", 404)
		c.ServeJSON()
		return
	}

	c.Data["json"] = utils.OKResponse(country)
	c.ServeJSON()
}

// Suggestions returns lightweight country data for home page autocomplete.
// Query param: q (string). Called by home.js AJAX.
func (c *CountriesAPIController) Suggestions() {
	query := c.GetString("q")

	suggestions, err := countrySvc.SearchSuggestions(query)
	if err != nil {
		c.Ctx.Output.SetStatus(500)
		c.Data["json"] = utils.ErrorResponse("Failed to fetch suggestions", 500)
		c.ServeJSON()
		return
	}

	c.Data["json"] = utils.OKResponse(suggestions)
	c.ServeJSON()
}
