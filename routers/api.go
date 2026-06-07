package routers

import (
	apicontrollers "TravelSphere/controllers/api"

	beego "github.com/beego/beego/v2/server/web"
)

func registerAPIRoutes() {
	// Country API
	beego.Router("/api/countries", &apicontrollers.CountriesAPIController{})
	beego.Router("/api/countries/suggestions", &apicontrollers.CountriesAPIController{}, "get:Suggestions")
	beego.Router("/api/countries/:slug", &apicontrollers.CountriesAPIController{}, "get:Detail")

	// Wishlist CRUD API
	beego.Router("/api/wishlist", &apicontrollers.WishlistAPIController{})
	beego.Router("/api/wishlist/:id", &apicontrollers.WishlistAPIController{}, "put:Update;delete:Delete")

	// Dashboard summary API
	beego.Router("/api/dashboard/summary", &apicontrollers.DashboardAPIController{})

	// Attractions API
	beego.Router("/api/attractions", &apicontrollers.AttractionsAPIController{})
}
