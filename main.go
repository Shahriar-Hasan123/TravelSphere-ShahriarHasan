package main

import (
	"TravelSphere/routers"
	"log"

	beego "github.com/beego/beego/v2/server/web"
	"github.com/joho/godotenv"
)

func main() {
	// Load .env file into environment variables.
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found — using system environment variables")
	}
	routers.Init()
	beego.Run()
}
