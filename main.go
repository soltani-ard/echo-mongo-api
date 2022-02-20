package main

import (
	"github.com/labstack/echo/v4"
	"github.com/soltani-ard/echo-mongo-api/configs"
	"github.com/soltani-ard/echo-mongo-api/routes"
)

func main() {

	e := echo.New()

	// run database
	configs.ConnectDB()

	// routes
	routes.UserRoute(e)

	// start server
	e.Logger.Fatal(e.Start(":8080"))
}
