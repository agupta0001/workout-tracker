package main

import (
	"net/http"
	"workout-tracker/models"

	"github.com/labstack/echo/v4"
)

func handleGetExercises(c echo.Context) error {
	var app = c.Get("app").(*App)

	var out models.Exercises

	err := app.queries.GetExercises.Select(&out)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, okResp{out})
}
