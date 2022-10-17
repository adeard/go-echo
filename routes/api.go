package routes

import (
	"go-echo/controller"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type CustomValidator struct {
	validator *validator.Validate
}

func (cv *CustomValidator) Validate(i interface{}) error {
	return cv.validator.Struct(i)
}

func ApiRoute(route *echo.Echo) {

	route.Validator = &CustomValidator{validator: validator.New()}

	g := route.Group("/v1")
	g.Use(controller.TokenRefresherMiddleware)
	g.Use(middleware.JWTWithConfig(middleware.JWTConfig{
		Claims:                  &controller.Claims{},
		SigningKey:              []byte(controller.GetJWTSecret()),
		TokenLookup:             "cookie:access-token", // "<source>:<name>"
		ErrorHandlerWithContext: controller.JWTErrorChecker,
	}))

	g.GET("/user", controller.UserGet())
	g.POST("/user", controller.UserCreate())
	g.PUT("/user/:email", controller.UserUpdate())
	g.DELETE("/user/:email", controller.UserDelete())

	g.GET("/inventory", controller.InventoryGet())
	g.POST("/inventory", controller.InventoryCreate())
	g.PUT("/inventory/:id", controller.InventoryUpdate())
	g.DELETE("/inventory/:id", controller.InventoryDelete())

}
