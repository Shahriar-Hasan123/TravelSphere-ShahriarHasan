// Init sets up middleware, routes, and static assets.

package routers

import (
	"TravelSphere/filters"

	beego "github.com/beego/beego/v2/server/web"
	"github.com/beego/beego/v2/server/web/context"
)

func Init() {
	registerFilters()
	registerSSRRoutes()
	registerAPIRoutes()

	beego.SetStaticPath("/static", "static")
}

func registerFilters() {
	// Request logging.
	beego.InsertFilter("/*", beego.BeforeRouter, loggingFilterHandler)

	// Auth for protected SSR routes.
	beego.InsertFilter("/wishlist", beego.BeforeExec, authFilterHandler)
	beego.InsertFilter("/dashboard", beego.BeforeExec, authFilterHandler)

	// Auth for wishlist API.
	beego.InsertFilter("/api/wishlist", beego.BeforeExec, authFilterHandler)
	beego.InsertFilter("/api/wishlist/*", beego.BeforeExec, authFilterHandler)
}

func loggingFilterHandler(ctx *context.Context) {
	filters.LoggingFilter(ctx)
}

func authFilterHandler(ctx *context.Context) {
	filters.AuthFilter(ctx)
}
