package routes

import (
	"github.com/labstack/echo/v4"
	"github.com/soltani-ard/echo-mongo-api/controllers"
)

// UserRoute => All routes related to users comes here
func UserRoute(e *echo.Echo) {
	e.POST("/user", controllers.CreateUser)
	e.GET("/user/:userId", controllers.GetUser)
	e.PUT("/user/:userId", controllers.EditUser)
	e.DELETE("/user/:userId", controllers.DeleteUser)
	e.GET("/users", controllers.GetAllUsers)
}
