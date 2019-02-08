package routers

import (
	"fmt"
	"net/http"
	"workshop01/controllers"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"github.com/spf13/viper"
)

//Init func
func Init(e *echo.Echo) {

	port := fmt.Sprintf(":%v", viper.GetString("port"))

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	fmt.Printf("Starting...")

	// CORS
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		//AllowOrigins: []string{"*"},
		//AllowMethods: []string{echo.GET, echo.HEAD, echo.PUT, echo.PATCH, echo.POST, echo.DELETE},
	}))

	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "OK")
	})

	e.GET("/todos", controllers.List)
	e.POST("/todos", controllers.Create)
	e.GET("/todos/:id", controllers.View)
	e.PUT("/todos/:id", controllers.Done)
	e.DELETE("/todos/:id", controllers.Delete)

	//GoRoutine
	e.GET("/hello", controllers.CallHelloRoutine)

	e.Logger.Fatal(e.Start(port))
}
