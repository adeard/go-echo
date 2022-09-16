package routes

import (
	"go-echo/controller"

	"github.com/labstack/echo/v4"
)

func PublicRoute(route *echo.Echo) {

	route.POST("/user/signin", controller.SignIn())

}
