package routes

import (
	"api_testing/controllers"

	"github.com/labstack/echo/v4"
)

func SetBookRoutes(e *echo.Echo) {
	e.GET("/books", controllers.GetBooksController)
	e.GET("/books/:id", controllers.GetBookController)
	e.POST("/books", controllers.CreateBookController)
	e.DELETE("/books/:id", controllers.DeleteBookController)
	e.PUT("/books/:id", controllers.UpdateBookController)
}
