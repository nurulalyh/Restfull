package routes

import (
	"test/configs"
	"test/features/users"

	"github.com/labstack/echo/v4"
)

func RouteUser(e *echo.Echo, uc users.UserHandlerInterface, cfg configs.ProgramConfig) {
	e.POST("/users", uc.Register())
	e.POST("/login", uc.Login())
	// e.GET("/users", uc.MyProfile(), echojwt.JWT([]byte(cfg.Secret)))
	// // e.GET("/users/:id",)
	// e.POST("/refresh", uc.RefreshToken(), echojwt.JWT([]byte(cfg.RefreshSecret)))
}
