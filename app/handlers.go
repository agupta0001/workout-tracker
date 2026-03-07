package main

import (
	"log"
	"os"

	"github.com/golang-jwt/jwt/v5"
	echojwt "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
)

const (
	sortAsc  = "asc"
	sortDesc = "desc"
)

type okResp struct {
	Data interface{} `json:"data"`
}

func initHTTPHandlers(e *echo.Echo, app *App) {
	e.HTTPErrorHandler = func(err error, c echo.Context) {

		if _, ok := err.(*echo.HTTPError); !ok {
			log.Println(err.Error())
		}
		e.DefaultHTTPErrorHandler(err, c)
	}

	api := e.Group("/api")

	api.GET("/health", handleHealthCheck)

	firebaseRoute := api.Group("")

	firebaseRoute.Use(firebaseAuthMiddleware)

	firebaseRoute.POST("/login", handleLogin)

	restricted := api.Group("")

	config := echojwt.Config{
		NewClaimsFunc: func(c echo.Context) jwt.Claims {
			return new(jwtClaims)
		},
		SigningKey: []byte(os.Getenv("SECRET")),
	}

	restricted.Use(echojwt.WithConfig(config))
	restricted.Use(authenticateRequest)

	restricted.GET("/api/exercises", handleGetExercises)
	restricted.GET("/api/workouts", handleGetWorkouts)
}

func handleHealthCheck(c echo.Context) error {
	return c.JSON(200, okResp{Data: "ok"})
}
