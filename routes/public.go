package routes

import (
	"go-echo/controller"

	"github.com/labstack/echo/v4"
)

func PublicRoute(route *echo.Echo) {

	route.POST("/user/signin", controller.SignIn())

	route.GET("/inventory", controller.InventoryGet())
	route.PUT("/inventory/:id", controller.InventoryUpdate())
	route.POST("/inventory", controller.InventoryCreate())

}
