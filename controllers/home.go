// HomeController handles GET / - the TravelSphere home page.

package controllers

import (
	"TravelSphere/models"
	"TravelSphere/services"
)

type HomeController struct {
	BaseController
}

var popularAttractions = []models.Attraction{
	{Name: "Eiffel Tower",      Kinds: []string{"architecture", "historic"}},
	{Name: "Grand Canyon",      Kinds: []string{"natural"}},
	{Name: "Sydney Opera House",Kinds: []string{"architecture", "theatre"}},
	{Name: "Colosseum",         Kinds: []string{"historic", "architecture"}},
}

// Get renders the home page with featured countries and popular attractions.
func (c *HomeController) Get() {
	featured, err := services.NewCountryService().GetFeaturedCountries()
	if err != nil {
		// show the page without featured countries.
		featured = nil
	}

	c.Data["FeaturedCountries"]  = featured
	c.Data["PopularAttractions"] = popularAttractions
	c.Data["ActiveNav"] = "home"
	c.TplName = "home.tpl"
	c.Layout = "layout.tpl"
}