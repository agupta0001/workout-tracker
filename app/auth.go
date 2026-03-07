package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"time"
	"workout-tracker/models"

	"firebase.google.com/go/v4/auth"
	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
)

type jwtClaims struct {
	UID string `json:"uid"`

	jwt.RegisteredClaims
}

func authenticateUser(app *App, idToken string) (*auth.UserRecord, error) {
	ctx := context.Background()
	auth, err := app.fbAuth.Auth(ctx)

	if err != nil {
		return nil, err
	}

	verifiedToken, err := auth.VerifyIDToken(ctx, idToken)

	if err != nil {
		return nil, err
	}

	user, err := auth.GetUser(ctx, verifiedToken.UID)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func firebaseAuthMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		app := c.Get("app").(*App)
		authHeader := c.Request().Header.Get("Authorization")

		if authHeader == "" {
			return echo.NewHTTPError(401, "Missing Authorization header")
		}

		idToken := authHeader[len("Bearer "):]

		user, err := authenticateUser(app, idToken)

		if err != nil {
			return echo.NewHTTPError(401, "Invalid or expired token")
		}

		c.Set("currentUser", user)
		return next(c)
	}
}

func authenticateRequest(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		app := c.Get("app").(*App)
		token := c.Get("user").(*jwt.Token)
		claims := token.Claims.(*jwtClaims)

		var user models.User

		err := app.queries.GetUserByUUID.Get(&user, claims.UID)

		if err != nil {
			log.Printf("Failed to fetch user with UID %s: %v", claims.UID, err)
			return echo.NewHTTPError(401, "User not found")
		}

		c.Set("currentUser", &user)

		return next(c)
	}
}

func handleLogin(c echo.Context) error {
	app := c.Get("app").(*App)
	requestUser := c.Get("currentUser").(*auth.UserRecord)

	user := models.User{
		UID:    requestUser.UID,
		Email:  requestUser.Email,
		Name:   requestUser.DisplayName,
		Avatar: requestUser.PhotoURL,
	}

	var out models.User

	err := app.queries.GetUserByUUID.Get(&out, user.UID)
	if err != nil {
		err = app.queries.CreateUser.Get(&out, user.UID, user.Email, user.Name, user.Avatar)
		if err != nil {
			return err
		}
	}

	tokenExpiresIn := time.Now().Add(24 * 1 * time.Hour)

	claim := &jwtClaims{
		out.UID,
		jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(tokenExpiresIn),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claim)

	signedToken, err := token.SignedString([]byte(os.Getenv("SECRET")))

	if err != nil {
		return echo.NewHTTPError(500, "Failed to generate JWT")
	}

	return c.JSON(http.StatusOK, okResp{Data: map[string]interface{}{
		"token":          signedToken,
		"user":           out,
		"tokenExpiresIn": tokenExpiresIn,
	}})

}
