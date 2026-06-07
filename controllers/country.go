// CountryController handles SSR routes for /countries and /countries/:slug.
package controllers

import (
	"TravelSphere/services"
)

// CountryController serves the Country Explorer and Destination Detail pages.
type CountryController struct {
	BaseController
}

// countryService is a package-level instance shared across requests.
var countryService = services.NewCountryService()

// Fetches the full country list server-side for the initial SSR load.
func (c *CountryController) Get() {
	search := c.GetString("search")
	region := c.GetString("region")

	countries, err := countryService.GetAllCountries(search, region)
	if err != nil {
		// Render the page with an empty list and a user-facing error notice.
		c.Data["Error"] = "Unable to load countries. Please try again later."
		c.Data["Countries"] = nil
	} else {
		c.Data["Countries"] = countries
	}

	c.Data["SearchQuery"] = search
	c.Data["RegionFilter"] = region
	c.Data["ActiveNav"] = "countries"
	c.TplName = "countries.tpl"
	c.Layout = "layout.tpl"
}

// Detail renders the Destination Detail page (/countries/:slug).
func (c *CountryController) Detail() {
	slug := c.Ctx.Input.Param(":slug")

	country, err := countryService.GetCountryBySlug(slug)
	if err != nil {
		c.Data["Error"] = "Unable to load country details. Please try again later."
		c.Data["Country"] = nil
		c.Data["ActiveNav"] = "countries"
		c.TplName = "destination.tpl"
		c.Layout = "layout.tpl"
		return
	}

	if country == nil {
		// Slug did not match any country - show user-friendly 404 page.
		c.Data["ActiveNav"] = "countries"
		c.TplName = "404.tpl"
		c.Layout = "layout.tpl"
		c.Ctx.Output.SetStatus(404)
		if err := c.Render(); err != nil {
			c.Ctx.WriteString("404 Not Found")
		}
		return
	}

	c.Data["Country"] = country
	c.Data["Attractions"] = nil
	c.Data["Weather"] = nil
	c.Data["ActiveNav"] = "countries"
	c.TplName = "destination.tpl"
	c.Layout = "layout.tpl"
}
