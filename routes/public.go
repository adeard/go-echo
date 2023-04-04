package routes

import (
	"go-echo/controller"

	"github.com/labstack/echo/v4"
)

func PublicRoute(route *echo.Echo) {

	route.POST("/login", controller.SignIn())
	route.GET("/protected", controller.ProtectedVerify())
	route.POST("/register", controller.UserCreate())
}
