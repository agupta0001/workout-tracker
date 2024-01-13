package models

import (
	"github.com/lib/pq"
	"gopkg.in/volatiletech/null.v6"
)

type Base struct {
	ID        int       `db:"id" json:"id"`
	CreatedAt null.Time `db:"created_at" json:"created_at"`
	UpdatedAt null.Time `db:"updated_at" json:"updated_at"`
}

type User struct {
	Base

	UUID   string      `db:"uuid" json:"uuid"`
	Email  string      `db:"email" json:"email"`
	Name   string      `db:"name" json:"name"`
	Avatar null.String `db:"avatar" json:"photo_url"`
	Token  null.String `db:"token" json:"token"`
}

type Weight struct {
	Base

	Measure float64 `db:"measure" json:"measure"`
}

type Exercise struct {
	Base

	Name string         `db:"name" json:"name"`
	Tags pq.StringArray `db:"tags" json:"tags"`
}

type Workout struct {
	Base

	ExerciseID int `db:"exercise_id" json:"exercise_id"`
	UserID     int `db:"user_id" json:"user_id"`
	WeightID   int `db:"weight_id" json:"weight_id"`
	SetNo      int `db:"set_no" json:"set_no"`
	Reps       int `db:"reps" json:"reps"`
}
