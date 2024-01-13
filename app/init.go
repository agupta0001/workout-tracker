package main

import (
	"fmt"
	"log"
	"os"
	"path"
	"strconv"
	"strings"
	"syscall"
	"time"
	"workout-tracker/models"

	"github.com/jmoiron/sqlx"
	"github.com/knadh/goyesql/v2"
	goyesqlx "github.com/knadh/goyesql/v2/sqlx"
	"github.com/knadh/stuffbin"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	flag "github.com/spf13/pflag"
)

const queryFilePath = "queries.sql"

// initDB initializes the main DB connection pool and parse and loads the app's
// SQL queries into a prepared query map.
func initDB() *sqlx.DB {
	var (
		host    = os.Getenv("DB_HOST")
		user    = os.Getenv("DB_USER")
		pass    = os.Getenv("DB_PASS")
		dbname  = os.Getenv("DB_NAME")
		sslmode = os.Getenv("DB_SSLMODE")
	)

	port, err := strconv.Atoi(os.Getenv("DB_PORT"))

	if err != nil {
		log.Fatalf("Invalid DB_PORT: %s", err)
	}

	fmt.Printf("connecting to db: %s:%d/%s", host, port, dbname)
	db, err := sqlx.Connect("postgres",
		fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s", host, port, user, pass, dbname, sslmode))

	if err != nil {
		log.Fatalf("Error connecting to DB: %s", err)
	}

	return db
}

func initFS(appDir string) stuffbin.FileSystem {
	var appFiles = []string{
		"./queries.sql:queries.sql",
		"./schema.sql:schema.sql",
	}

	execPath, err := os.Executable()
	if err != nil {
		log.Fatalf("Error getting executable path: %v", err)
	}

	hasEmbed := true

	fs, err := stuffbin.UnStuff(execPath)
	if err != nil {
		hasEmbed = false

		// Running in local mode. Load local assets into
		// the in-memory stuffbin.FileSystem.
		fmt.Printf("\nunable to initialize embedded filesystem (%v). Using local filesystem", err)

		fs, err = stuffbin.NewLocalFS("/")
		if err != nil {
			log.Fatalf("\nfailed to initialize local file for assets: %v", err)
		}
	}

	files := []string{}

	if !hasEmbed {
		files = append(files, joinFSPaths(appDir, appFiles)...)
	}

	if len(files) == 0 {
		return fs
	}

	fStatic, err := stuffbin.NewLocalFS("/", files...)
	if err != nil {
		log.Fatalf("error initializing static files: %v", err)
	}

	if err := fs.Merge(fStatic); err != nil {
		log.Fatalf("error merging static files: %v", err)
	}

	return fs
}

func joinFSPaths(root string, paths []string) []string {
	out := make([]string, 0, len(paths))
	for _, p := range paths {
		f := strings.Split(p, ":")

		out = append(out, path.Join(root, f[0])+":"+f[1])
	}

	return out
}

func initFlags() {
	f := flag.NewFlagSet("config", flag.ContinueOnError)
	f.Usage = func() {
		fmt.Println(f.FlagUsages())
		os.Exit(0)
	}

	f.Bool("install", false, "Install the DB schema")
	if err := f.Parse(os.Args[1:]); err != nil {
		log.Fatalf("Error parsing flags: %s", err)
	}

	vp.BindPFlags(f)
}

func readQueries(sqlFile string, db *sqlx.DB, fs stuffbin.FileSystem) goyesql.Queries {
	qB, err := fs.Read(sqlFile)

	if err != nil {
		log.Fatalf("error reading SQL file %s: %v", sqlFile, err)
	}

	qMap, err := goyesql.ParseBytes(qB)

	if err != nil {
		log.Fatalf("error parsing SQL queries: %v", err)
	}

	return qMap
}

func prepareQueries(qMap goyesql.Queries, db *sqlx.DB) *models.Queries {
	var q models.Queries

	if err := goyesqlx.ScanToStruct(&q, qMap, db.Unsafe()); err != nil {
		log.Fatalf("error preparing SQL queries: %v", err)
	}

	return &q
}

func initHTTPServer(app *App) *echo.Echo {
	var srv = echo.New()
	srv.HideBanner = true

	srv.Use(middleware.Logger())
	srv.Use(func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			c.Set("app", app)
			return next(c)
		}
	})

	initHTTPHandlers(srv, app)

	go func() {
		if err := srv.Start(fmt.Sprintf("localhost:%s", os.Getenv("PORT"))); err != nil {
			if strings.Contains(err.Error(), "Server closed") {
				log.Println("HTTP server shut down")
			} else {
				log.Fatalf("error starting HTTP server: %v", err)
			}
		}
	}()

	return srv
}

func awaitReload(sigChan chan os.Signal, closerWait chan bool, closer func()) chan bool {

	out := make(chan bool)

	// Respawn a new process and exit the running one.
	respawn := func() {
		if err := syscall.Exec(os.Args[0], os.Args, os.Environ()); err != nil {
			log.Fatalf("error spawning process: %v", err)
		}
		os.Exit(0)
	}

	// Listen for reload signal.
	go func() {
		for range sigChan {
			log.Println("reloading on signal ...")

			go closer()
			select {
			case <-closerWait:
				// Wait for the closer to finish.
				respawn()
			case <-time.After(time.Second * 3):
				// Or timeout and force close.
				respawn()
			}
		}
	}()

	return out
}
