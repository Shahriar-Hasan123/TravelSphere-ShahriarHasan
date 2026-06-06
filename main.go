package main

import (
	"TravelSphere/routers"
	beego "github.com/beego/beego/v2/server/web"
)

func main() {
	// Register SSR page routes and JSON API routes.
	routers.Init()
	beego.Run()
}
