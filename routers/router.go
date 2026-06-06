// SSR page routes - every route here must return text/html.
package routers

import (
	beego "github.com/beego/beego/v2/server/web"
)

// registerSSRRoutes maps browser-navigable URLs to SSR controllers.
func registerSSRRoutes() {
	_ = beego.Router
}
