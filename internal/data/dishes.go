package data

import (
	"context"
	"database/sql"
	"errors"
	"time"
	"unicode/utf8"

	"github.com/lib/pq"
	"github.com/xtommas/food-backend/internal/validator"
)

type Dish struct {
	Id          int64    `json:"id"`
	Name        string   `json:"name"`
	Price       Price    `json:"price"`
	Description string   `json:"description"`
	Category    []string `json:"category"`
	Photo       string   `json:"photo,omitempty"`
	Available   bool     `json:"available"`
}

type DishModel struct {
	DB *sql.DB
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

	// Add photo validation and require it to be provided in the request, when I figure that out
}

func (d DishModel) Insert(dish *Dish) error {
	query := `
			INSERT INTO dishes (name, price, description, category, photo)
			VALUES ($1, $2, $3, $4, $5)
			RETURNING id, available`

	args := []interface{}{dish.Name, dish.Price, dish.Description, pq.Array(dish.Category), dish.Photo}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)

	defer cancel()

	return d.DB.QueryRowContext(ctx, query, args...).Scan(&dish.Id, &dish.Available)
}

func (d DishModel) Get(id int64) (*Dish, error) {
	if id < 1 {
		return nil, ErrRecordNotFound
	}

	query := `
		SELECT id, name, price, description, category, photo, available
		FROM dishes
		WHERE id = $1`

	var dish Dish

	// create context with a 3 second timeout deadline
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)

	defer cancel()

	// use QueryRowContext to pass the context with the deadline
	err := d.DB.QueryRowContext(ctx, query, id).Scan(
		&dish.Id,
		&dish.Name,
		&dish.Price,
		&dish.Description,
		pq.Array(&dish.Category),
		&dish.Photo,
		&dish.Available,
	)

	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrRecordNotFound
		default:
			return nil, err
		}
	}

	return &dish, nil
}

func (d DishModel) Update(dish *Dish) error {
	return nil
}

func (d DishModel) Delete(id int64) error {
	return nil
}
