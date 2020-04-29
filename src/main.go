package main

import (
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"github.com/nnyam3831/todays_stock_server/src/api"
)

func main() {
	e := echo.New()
	e.GET("/", api.Home)
	e.GET("/golden", api.GetGQ)
	e.GET("/kos", api.GetKOS)
	e.GET("/rise", api.GetRise)
	e.GET("/search", api.GetSearch)

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{echo.GET, echo.HEAD, echo.PUT, echo.PATCH, echo.POST, echo.DELETE},
	}))

	e.Logger.Fatal(e.Start(":1323"))
}
