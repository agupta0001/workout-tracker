package models

import "github.com/jmoiron/sqlx"

type Queries struct {
	GetExercises    *sqlx.Stmt `query:"get-exercises"`
	GetExerciseById *sqlx.Stmt `query:"get-exercise-by-id"`
}
