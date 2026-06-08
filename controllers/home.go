package controllers

import (
	"TravelSphere/models"
	"TravelSphere/services"
	"log"
)

type HomeController struct {
	BaseController
}

// staticAttractions is the fallback list when OpenTripMap is unavailable.
var staticAttractions = []models.Attraction{
	{Name: "Eiffel Tower",       Kinds: []string{"architecture", "historic"}},
	{Name: "Grand Canyon",       Kinds: []string{"natural"}},
	{Name: "Sydney Opera House", Kinds: []string{"architecture", "theatre"}},
	{Name: "Colosseum",          Kinds: []string{"historic", "architecture"}},
}

// famousCoords holds coordinates for globally recognised attraction hotspots.
// Used to seed the home page Popular attractions section from OpenTripMap.
var famousCoords = []struct{ lat, lon float64 }{
	{48.8584, 2.2945},   // Paris — Eiffel Tower area
	{36.1069, -112.1129}, // Grand Canyon
	{41.8902, 12.4922},  // Rome — Colosseum area
	{-33.8568, 151.2153}, // Sydney Opera House area
}

// Get renders the home page with featured countries and popular attractions.
func (c *HomeController) Get() {
	countrySvc := services.NewCountryService()

	featured, err := countrySvc.GetFeaturedCountries()
	if err != nil {
		log.Printf("home: featured countries error: %v", err)
		featured = nil
	}

	attractions := fetchHomeAttractions()

	c.Data["ActiveNav"]          = "home"
	c.Data["FeaturedCountries"]  = featured
	c.Data["PopularAttractions"] = attractions
	c.TplName = "home.tpl"
	c.Layout  = "layout.tpl"
}

// fetchHomeAttractions attempts to load real attractions from OpenTripMap. Falls back to the static list on any error so the home page never breaks.
func fetchHomeAttractions() []models.Attraction {
	attractionSvc := services.NewAttractionService()
	seen          := map[string]bool{}
	var results   []models.Attraction

	for _, coord := range famousCoords {
		list, err := attractionSvc.GetAttractionsByCoords(coord.lat, coord.lon)
		if err != nil {
			log.Printf("home: attractions fetch error: %v", err)
			continue
		}
		for _, a := range list {
			if !seen[a.Name] && len(results) < 8 {
				seen[a.Name] = true
				results = append(results, a)
			}
		}
		if len(results) >= 8 {
			break
		}
	}

	if len(results) == 0 {
		return staticAttractions
	}
	return results
}