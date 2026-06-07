// SSR page routes - every route here must return text/html.
package routers

import (
	"TravelSphere/controllers"
	beego "github.com/beego/beego/v2/server/web"
)

// registerSSRRoutes maps browser-navigable URLs to SSR controllers.
func registerSSRRoutes() {
	// Home page
	beego.Router("/", &controllers.HomeController{})

	// Country Explorer and Destination Detail
	beego.Router("countries/", &controllers.CountryController{})
	beego.Router("countries/:slug", &controllers.CountryController{}, "get:Detail")

	// Auth routes
	beego.Router("/login", &controllers.AuthController{}, "get:ShowLogin;post:DoLogin")
	beego.Router("/logout", &controllers.AuthController{}, "get:DoLogout")

	
	beego.Router("/wishlist", &controllers.WishlistController{})
	beego.Router("/dashboard", &controllers.DashboardController{})
}
