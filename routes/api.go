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

	g := route.Group("/user")
	g.Use(controller.TokenRefresherMiddleware)
	g.Use(middleware.JWTWithConfig(middleware.JWTConfig{
		Claims:                  &controller.Claims{},
		SigningKey:              []byte(controller.GetJWTSecret()),
		TokenLookup:             "cookie:access-token", // "<source>:<name>"
		ErrorHandlerWithContext: controller.JWTErrorChecker,
	}))

	g.GET("", controller.UserGet())
	g.POST("", controller.UserCreate())
	g.PUT("/:email", controller.UserUpdate())
	g.DELETE("/:email", controller.UserDelete())

}
