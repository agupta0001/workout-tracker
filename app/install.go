package main

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/jmoiron/sqlx"
	"github.com/knadh/stuffbin"
	"github.com/lib/pq"
)

func isTableNotExistError(err error) bool {
	if p, ok := err.(*pq.Error); ok {
		if p.Code == "42P01" {
			return true
		}
	}

	return false
}

func checkSchema(db *sqlx.DB) (bool, error) {
	if _, err := db.Exec("SELECT * FROM users"); err != nil {
		if isTableNotExistError(err) {
			return false, nil
		}

		return false, err
	}

	return true, nil
}

func installSchema(db *sqlx.DB, fs stuffbin.FileSystem) error {
	q, err := fs.Read("./schema.sql")
	if err != nil {
		return err
	}

	if _, err := db.Exec(string(q)); err != nil {
		return err
	}

	return nil
}

func install(db *sqlx.DB, fs stuffbin.FileSystem) {
	fmt.Println("")

	fmt.Println("** first time installation **")
	fmt.Printf("** IMPORTANT: This will wipe existing listmonk tables and types in the DB '%s' **",
		os.Getenv("DB_NAME"))

	fmt.Println("")

	var ok string
	fmt.Print("continue (y/N)?  ")
	if _, err := fmt.Scanf("%s", &ok); err != nil {
		log.Fatalf("error reading value from terminal: %v", err)
	}
	if strings.ToLower(ok) != "y" {
		fmt.Println("install cancelled.")
		return
	}

	if err := installSchema(db, fs); err != nil {
		log.Fatalf("error migrating DB schema: %v", err)
	}

	log.Printf("setup complete")
	log.Printf(`run the program and access the dashboard at %s`, os.Getenv("PORT"))
}
