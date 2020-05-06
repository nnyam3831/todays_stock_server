package main

import (
	"os"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

func main() {
	e := echo.New()
	port := os.Getenv("PORT")
	e.GET("/", Home)
	e.GET("/golden", GetGQ)
	e.GET("/kos", GetKOS)
	e.GET("/rise", GetRise)
	e.GET("/search", GetSearch)

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{echo.GET, echo.HEAD, echo.PUT, echo.PATCH, echo.POST, echo.DELETE},
	}))

	e.Logger.Fatal(e.Start(":" + port))
}
