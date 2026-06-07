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
	beego.InsertFilter("/*", beego.BeforeRouter,
		func(ctx *context.Context) { filters.LoggingFilter(ctx) },
	)

	// Auth for protected SSR routes.
	beego.InsertFilter("/wishlist", beego.BeforeExec,
		func(ctx *context.Context) { filters.AuthFilter(ctx) },
	)
	beego.InsertFilter("/dashboard", beego.BeforeExec,
		func(ctx *context.Context) { filters.AuthFilter(ctx) },
	)

	// Auth for wishlist API.
	beego.InsertFilter("/api/wishlist", beego.BeforeExec,
		func(ctx *context.Context) { filters.AuthFilter(ctx) },
	)
	beego.InsertFilter("/api/wishlist/*", beego.BeforeExec,
		func(ctx *context.Context) { filters.AuthFilter(ctx) },
	)
}
