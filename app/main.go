package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"workout-tracker/models"

	"github.com/jmoiron/sqlx"
	"github.com/joho/godotenv"
	"github.com/knadh/stuffbin"
	"github.com/spf13/viper"
)

type App struct {
	db      *sqlx.DB
	queries *models.Queries
	fs      stuffbin.FileSystem

	chReload chan os.Signal
}

var (
	db      *sqlx.DB
	queries *models.Queries
	fs      stuffbin.FileSystem
	vp      *viper.Viper = viper.New()

	appDir string = "."
)

func init() {
	initFlags()

	err := godotenv.Load()

	if err != nil {
		log.Fatalf("Error loading .env file: %s", err)
	}

	db = initDB()
	fs = initFS(appDir)

	if vp.GetBool("install") {
		install(db, fs)
		os.Exit(0)
	}

	if ok, err := checkSchema(db); err != nil {
		log.Fatalf("\nError checking DB schema: %s", err)
	} else if !ok {
		log.Fatalf("\nDB schema is not up to date")
	}

	qMap := readQueries(queryFilePath, db, fs)

	queries = prepareQueries(qMap, db)
}

func main() {
	app := &App{
		db:      db,
		queries: queries,
		fs:      fs,
	}

	srv := initHTTPServer(app)

	app.chReload = make(chan os.Signal)
	signal.Notify(app.chReload, syscall.SIGHUP)

	closerWait := make(chan bool)
	<-awaitReload(app.chReload, closerWait, func() {
		ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
		defer cancel()
		srv.Shutdown(ctx)

		// Close the DB pool.
		app.db.DB.Close()

		// Signal the close.
		closerWait <- true
	})
}
