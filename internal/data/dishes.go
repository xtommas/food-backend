package data

import (
	"unicode/utf8"

	"github.com/xtommas/food-backend/internal/validator"
)

type Dish struct {
	Id          int64    `json:"id"`
	Name        string   `json:"name"`
	Price       Price    `json:"price"`
	Description string   `json:"description"`
	Category    []string `json:"category"`
	Photo       string   `json:"photo,omitempty"`
}

func ValidateDish(v *validator.Validator, dish *Dish) {
	v.Check(dish.Name != "", "name", "must be provided")
	v.Check(utf8.RuneCountInString(dish.Name) <= 100, "name", "must be no more than 100 characters long")

	v.Check(dish.Price != 0, "price", "must be provided")
	v.Check(dish.Price > 0, "price", "must be a positive number")

	v.Check(dish.Description != "", "description", "must be provided")
	v.Check(utf8.RuneCountInString(dish.Description) <= 280, "description", "must be no more than 280 characters long")

	v.Check(dish.Category != nil, "category", "must be provided")
	v.Check(len(dish.Category) >= 1, "category", "must contain at least one category")
	v.Check(len(dish.Category) <= 5, "category", "must not contain more than 5 categories")
	v.Check(validator.Unique(dish.Category), "category", "must not contain duplicate values")
}
