// CountryController handles SSR routes for /countries and /countries/:slug.
package controllers

import (
	"TravelSphere/services"
	"log"
)

// CountryController serves the Country Explorer and Destination Detail pages.
type CountryController struct {
	BaseController
}

// Get renders the Country Explorer page (/countries).
func (c *CountryController) Get() {
	countryService := services.NewCountryService()
	search := c.GetString("search")
	region := c.GetString("region")

	countries, err := countryService.GetAllCountries(search, region)
	if err != nil {
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
	countryService := services.NewCountryService()
	slug := c.Ctx.Input.Param(":slug")

	country, err := countryService.GetCountryBySlug(slug)
	if err != nil {
		log.Printf("country detail error for slug %q: %v", slug, err)
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

	// Fetch attractions when coordinates are available.
	var attractions interface{}
	attractionService := services.NewAttractionService()
	if len(country.Latlng) == 2 {
		list, err := attractionService.GetAttractionsByCoords(
			country.Latlng[0], country.Latlng[1],
		)
		if err != nil {
			// Degrade gracefully — show page without attractions.
			log.Printf("attractions fetch error for %s: %v", country.Name, err)
		} else {
			attractions = list
		}
	}

	// Fetch weather — nil when key not set or city not found.
	weatherService := services.NewWeatherService()
	var weather interface{}
	if weatherService.IsConfigured() && country.Capital != "" {
		w, err := weatherService.GetWeather(country.Capital)
		if err != nil {
			log.Printf("weather fetch error for %s: %v", country.Capital, err)
		} else {
			weather = w
		}
	}

	c.Data["Country"] = country
	c.Data["Attractions"] = attractions
	c.Data["Weather"] = weather
	c.Data["ActiveNav"] = "countries"
	c.TplName = "destination.tpl"
	c.Layout = "layout.tpl"
}
