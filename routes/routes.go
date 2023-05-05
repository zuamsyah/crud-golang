package routes

import (
	"crud-golang/handler"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func Init() *echo.Echo {
	e := echo.New()

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	e.POST("/login", handler.Login)
	e.POST("/register", handler.Register)

	api := e.Group("/api")
	api.Use(handler.AuthMiddleware)

	// students endpoints
	api.GET("/students", handler.GetStudents)
	api.POST("/students", handler.CreateStudent)
	api.GET("/students/:id", handler.GetStudent)
	api.PUT("/students/:id", handler.UpdateStudent)
	api.DELETE("/students/:id", handler.DeleteStudent)

	// courses endpoints
	api.GET("/courses", handler.GetCourses)
	api.POST("/courses", handler.CreateCourse)
	api.GET("/courses/:id", handler.GetCourse)
	api.PUT("/courses/:id", handler.UpdateCourse)
	api.DELETE("/courses/:id", handler.DeleteCourse)

	api.GET("/profile", handler.GetUserProfile)

	api.GET("/users", handler.GetUsers)

	return e
}
