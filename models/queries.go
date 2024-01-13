package models

import "github.com/jmoiron/sqlx"

type Queries struct {
	GetWeights *sqlx.Stmt `query:"get-workouts"`
}
