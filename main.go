package main

import (
	"github.com/labstack/echo/v4"
	"github.com/soltani-ard/echo-mongo-api/configs"
)

func main() {

	e := echo.New()

	// run database
	configs.ConnectDB()

	// start server
	e.Logger.Fatal(e.Start(":8080"))
}
