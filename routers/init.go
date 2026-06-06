// Init is called once from main.go to register all application routes.
package routers

import beego "github.com/beego/beego/v2/server/web"

// Init registers SSR page routes, JSON API routes, and static file serving.
func Init() {
	registerSSRRoutes()
	registerAPIRoutes()

	// Serve everything under /static/* from the static/ directory.
	beego.SetStaticPath("/static", "static")
}