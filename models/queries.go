package models

import "github.com/jmoiron/sqlx"

type Queries struct {
	// User queries
	GetUserById   *sqlx.Stmt `query:"get-user-by-id"`
	GetUserByUUID *sqlx.Stmt `query:"get-user-by-uuid"`
	CreateUser    *sqlx.Stmt `query:"create-user"`

	// Exercise queries
	GetExercises    *sqlx.Stmt `query:"get-exercises"`
	GetExerciseById *sqlx.Stmt `query:"get-exercise-by-id"`

	// Workout queries
	GetWorkouts *sqlx.Stmt `query:"get-workouts"`
}
