package main

import (
	"log"

	"github.com/labstack/echo/v4"
)

const (
	sortAsc  = "asc"
	sortDesc = "desc"
)

type okResp struct {
	Data interface{} `json:"data"`
}

// registerHandlers registers HTTP handlers.
func initHTTPHandlers(e *echo.Echo, app *App) {
	e.HTTPErrorHandler = func(err error, c echo.Context) {

		if _, ok := err.(*echo.HTTPError); !ok {
			log.Println(err.Error())
		}
		e.DefaultHTTPErrorHandler(err, c)
	}

	e.GET("/api/health", handleHealthCheck)
}

func handleHealthCheck(c echo.Context) error {
	return c.JSON(200, okResp{Data: "ok"})
}
