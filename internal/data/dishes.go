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
	query := `
		UPDATE dishes
		SET name = $1, price = $2, description = $3, category = $4, photo = $5, available = $6
		WHERE id = $7`

	args := []interface{}{
		dish.Name,
		dish.Price,
		dish.Description,
		pq.Array(dish.Category),
		dish.Photo,
		dish.Available,
		dish.Id,
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)

	defer cancel()

	result, err := d.DB.ExecContext(ctx, query, args...)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return ErrEditConflict
		default:
			return err
		}
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return ErrRecordNotFound
	}

	return nil
}

func (d DishModel) Delete(id int64) error {
	if id < 1 {
		return ErrRecordNotFound
	}

	query := `
		DELETE FROM dishes
		WHERE id = $1`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)

	defer cancel()

	result, err := d.DB.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return ErrRecordNotFound
	}

	return nil
}

func (d DishModel) GetAll(name string, category []string, available string, filters Filters) ([]*Dish, error) {
	query := `
		SELECT id, name, price, description, category, photo, available
		FROM dishes
		WHERE (to_tsvector('simple', name) @@ plainto_tsquery('simple', $1) OR $1 = '')
		AND (category @> $2 OR $2 = '{}')
		ORDER BY id`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	rows, err := d.DB.QueryContext(ctx, query, name, pq.Array(category))
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	dishes := []*Dish{}

	for rows.Next() {
		var dish Dish

		err := rows.Scan(
			&dish.Id,
			&dish.Name,
			&dish.Price,
			&dish.Description,
			pq.Array(&dish.Category),
			&dish.Photo,
			&dish.Available,
		)

		if err != nil {
			return nil, err
		}

		dishes = append(dishes, &dish)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return dishes, nil
}
