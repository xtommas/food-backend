package data

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"
	"unicode/utf8"

	"github.com/lib/pq"
	"github.com/xtommas/food-backend/internal/validator"
)

type Dish struct {
	ID           int64     `json:"id"`
	RestaurantID int64     `json:"restaurant_id"`
	Name         string    `json:"name"`
	Price        Price     `json:"price"`
	Description  string    `json:"description"`
	Categories   []string  `json:"categories"`
	Photo        string    `json:"photo,omitempty"`
	Available    bool      `json:"available"`
	UpdatedAt    time.Time `json:"updated_at"`
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

	v.Check(dish.Categories != nil, "categories", "must be provided")
	v.Check(len(dish.Categories) >= 1, "categories", "must contain at least one category")
	v.Check(len(dish.Categories) <= 5, "categories", "must not contain more than 5 categories")
	v.Check(validator.Unique(dish.Categories), "categories", "must not contain duplicate values")
}

func (d DishModel) Insert(dish *Dish) error {
	query := `
		INSERT INTO dishes (restaurant_id, name, price, description, categories, photo)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id, available, updated_at`

	args := []interface{}{dish.RestaurantID, dish.Name, dish.Price, dish.Description, pq.Array(dish.Categories), dish.Photo}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	return d.DB.QueryRowContext(ctx, query, args...).Scan(&dish.ID, &dish.Available, &dish.UpdatedAt)
}

func (d DishModel) Get(id int64) (*Dish, error) {
	if id < 1 {
		return nil, ErrRecordNotFound
	}

	query := `
		SELECT id, restaurant_id, name, price, description, categories, photo, available, updated_at
		FROM dishes
		WHERE id = $1`

	var dish Dish

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := d.DB.QueryRowContext(ctx, query, id).Scan(
		&dish.ID,
		&dish.RestaurantID,
		&dish.Name,
		&dish.Price,
		&dish.Description,
		pq.Array(&dish.Categories),
		&dish.Photo,
		&dish.Available,
		&dish.UpdatedAt,
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

func (d DishModel) GetForRestaurant(id int64, restaurantID int64) ([]*Dish, error) {
	if id < 1 {
		return nil, ErrRecordNotFound
	}

	query := `
		SELECT id, restaurant_id, name, price, description, categories, photo, available, updated_at
		FROM dishes
		WHERE id = $1 AND restaurant_id = $2`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	rows, err := d.DB.QueryContext(ctx, query, id, restaurantID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var dishes []*Dish

	for rows.Next() {
		var dish Dish

		err := rows.Scan(
			&dish.ID,
			&dish.RestaurantID,
			&dish.Name,
			&dish.Price,
			&dish.Description,
			pq.Array(&dish.Categories),
			&dish.Photo,
			&dish.Available,
			&dish.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}

		dishes = append(dishes, &dish)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	if len(dishes) == 0 {
		return nil, ErrRecordNotFound
	}

	return dishes, nil
}

func (d DishModel) Update(dish *Dish) error {
	query := `
		UPDATE dishes
		SET name = $1, price = $2, description = $3, categories = $4, photo = $5, available = $6
		WHERE id = $7
		RETURNING updated_at`

	args := []interface{}{
		dish.Name,
		dish.Price,
		dish.Description,
		pq.Array(dish.Categories),
		dish.Photo,
		dish.Available,
		dish.ID,
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := d.DB.QueryRowContext(ctx, query, args...).Scan(&dish.UpdatedAt)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return ErrRecordNotFound
		default:
			return err
		}
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

func (d DishModel) GetAll(name string, categories []string, available sql.NullBool, filters Filters) ([]*Dish, Metadata, error) {
	query := fmt.Sprintf(`
		SELECT COUNT(*) OVER(), id, restaurant_id, name, price, description, categories, photo, available, updated_at
		FROM dishes
		WHERE (to_tsvector('simple', name) @@ plainto_tsquery('simple', $1) OR $1 = '')
		AND (categories @> $2 OR $2 = '{}')
		AND (available = $3 OR $3 IS NULL)
		ORDER BY %s %s, id ASC
		LIMIT $4 OFFSET $5`, filters.sortColumn(), filters.sortDirection())

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	rows, err := d.DB.QueryContext(ctx, query, name, pq.Array(categories), available, filters.limit(), filters.offset())
	if err != nil {
		return nil, Metadata{}, err
	}
	defer rows.Close()

	totalRecords := 0
	var dishes []*Dish

	for rows.Next() {
		var dish Dish

		err := rows.Scan(
			&totalRecords,
			&dish.ID,
			&dish.RestaurantID,
			&dish.Name,
			&dish.Price,
			&dish.Description,
			pq.Array(&dish.Categories),
			&dish.Photo,
			&dish.Available,
			&dish.UpdatedAt,
		)
		if err != nil {
			return nil, Metadata{}, err
		}

		dishes = append(dishes, &dish)
	}

	if err = rows.Err(); err != nil {
		return nil, Metadata{}, err
	}

	metadata := calculateMetadata(totalRecords, filters.Page, filters.PageSize)

	return dishes, metadata, nil
}

func (d DishModel) GetAllForRestaurant(restaurantID int64, name string, categories []string, available sql.NullBool, filters Filters) ([]*Dish, Metadata, error) {
	query := fmt.Sprintf(`
		SELECT COUNT(*) OVER(), id, restaurant_id, name, price, description, categories, photo, available, updated_at
		FROM dishes
		WHERE (to_tsvector('simple', name) @@ plainto_tsquery('simple', $1) OR $1 = '')
		AND (categories @> $2 OR $2 = '{}')
		AND (available = $3 OR $3 IS NULL)
		AND restaurant_id = $4
		ORDER BY %s %s, id ASC
		LIMIT $5 OFFSET $6`, filters.sortColumn(), filters.sortDirection())

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	rows, err := d.DB.QueryContext(ctx, query, name, pq.Array(categories), available, restaurantID, filters.limit(), filters.offset())
	if err != nil {
		return nil, Metadata{}, err
	}
	defer rows.Close()

	totalRecords := 0
	var dishes []*Dish

	for rows.Next() {
		var dish Dish

		err := rows.Scan(
			&totalRecords,
			&dish.ID,
			&dish.RestaurantID,
			&dish.Name,
			&dish.Price,
			&dish.Description,
			pq.Array(&dish.Categories),
			&dish.Photo,
			&dish.Available,
			&dish.UpdatedAt,
		)
		if err != nil {
			return nil, Metadata{}, err
		}

		dishes = append(dishes, &dish)
	}

	if err = rows.Err(); err != nil {
		return nil, Metadata{}, err
	}

	metadata := calculateMetadata(totalRecords, filters.Page, filters.PageSize)

	return dishes, metadata, nil
}
