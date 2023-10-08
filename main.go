package main

import (
	"fmt"
	"test/configs"
	"test/features/users/data"
	"test/features/users/handler"
	"test/features/users/service"
	"test/helper"
	"test/routes"
	"test/utils/database"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	e := echo.New()
	var config = configs.InitConfig()

	db := database.InitDB(*config)
	database.Migrate(db)

	userModel := data.New(db)
	generator := helper.NewGenerator()
	jwtInterface := helper.New(config.Secret, config.RefreshSecret)
	userServices := service.New(userModel, generator, jwtInterface)

	userControll := handler.NewHandler(userServices)


	e.Pre(middleware.RemoveTrailingSlash())

	e.Use(middleware.CORS())
	e.Use(middleware.LoggerWithConfig(
		middleware.LoggerConfig{
			Format: "method=${method}, uri=${uri}, status=${status}, time=${time_rfc3339}\n",
		}))

	routes.RouteUser(e, userControll, *config)

	e.Logger.Fatal(e.Start(fmt.Sprintf(":%d", config.ServerPort)).Error())
}
