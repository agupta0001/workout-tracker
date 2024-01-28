package main

import (
	"net/http"
	"workout-tracker/models"

	"github.com/labstack/echo/v4"
)

func handleGetWorkouts(c echo.Context) error {
	var app = c.Get("app").(*App)

	out := models.Exercise{}

	err := app.queries.GetExerciseById.Select(&out, 1)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, okResp{out})
}
