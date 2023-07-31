package data

import (
	"database/sql"
	"errors"
)

var (
	ErrRecordNotFound = errors.New("record not found")
)

type Models struct {
	Dishes DishModel
}

func NewModels(db *sql.DB) Models {
	return Models{
		Dishes: DishModel{DB: db},
	}
}
