package data

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/xtommas/food-backend/internal/validator"
)

type Order struct {
	ID           int64     `json:"id"`
	UserID       int64     `json:"user_id"`
	RestaurantID int64     `json:"restaurant_id"`
	Total        Price     `json:"total"`
	Address      string    `json:"address"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
	Status       string    `json:"status"`
}

var validStatuses = []string{"pending", "confirmed", "preparing", "ready", "delivered", "cancelled"}

var validTransitions = map[string][]string{
	"pending":   {"confirmed", "cancelled"},
	"confirmed": {"preparing", "cancelled"},
	"preparing": {"ready"},
	"ready":     {"delivered"},
	"delivered": {},
	"cancelled": {},
}

func ValidateAddress(v *validator.Validator, address string) {
	v.Check(address != "", "address", "must be provided")
}

func ValidateStatus(v *validator.Validator, status string) {
	v.Check(status != "", "status", "must be provided")
	v.Check(validator.PermittedValue(status, validStatuses...), "status", "invalid status")
}

func ValidateStatusTransition(v *validator.Validator, from, to string) {
	allowed, ok := validTransitions[from]
	if !ok {
		v.AddError("status", "current status is unrecognised")
		return
	}
	v.Check(validator.PermittedValue(to, allowed...), "status", "invalid transition from "+from+" to "+to)
}

func ValidateOrder(v *validator.Validator, order *Order) {
	ValidateAddress(v, order.Address)
	ValidateStatus(v, order.Status)
}

type OrderModel struct {
	DB *sql.DB
}

func (o OrderModel) Insert(order *Order) error {
	query := `
		INSERT INTO orders (user_id, restaurant_id, total, address, status)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id, created_at, updated_at`

	args := []interface{}{order.UserID, order.RestaurantID, order.Total, order.Address, order.Status}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	return o.DB.QueryRowContext(ctx, query, args...).Scan(&order.ID, &order.CreatedAt, &order.UpdatedAt)
}

func (o OrderModel) GetForRestaurant(id int64, restaurantID int64) (*Order, error) {
	if id < 1 {
		return nil, ErrRecordNotFound
	}

	query := `
		SELECT id, user_id, restaurant_id, total, address, created_at, updated_at, status
		FROM orders
		WHERE id = $1 AND restaurant_id = $2`

	var order Order

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := o.DB.QueryRowContext(ctx, query, id, restaurantID).Scan(
		&order.ID,
		&order.UserID,
		&order.RestaurantID,
		&order.Total,
		&order.Address,
		&order.CreatedAt,
		&order.UpdatedAt,
		&order.Status,
	)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrRecordNotFound
		default:
			return nil, err
		}
	}

	return &order, nil
}

func (o OrderModel) GetForUser(id int64, userID int64) (*Order, error) {
	if id < 1 {
		return nil, ErrRecordNotFound
	}

	query := `
		SELECT id, user_id, restaurant_id, total, address, created_at, updated_at, status
		FROM orders
		WHERE id = $1 AND user_id = $2`

	var order Order

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := o.DB.QueryRowContext(ctx, query, id, userID).Scan(
		&order.ID,
		&order.UserID,
		&order.RestaurantID,
		&order.Total,
		&order.Address,
		&order.CreatedAt,
		&order.UpdatedAt,
		&order.Status,
	)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrRecordNotFound
		default:
			return nil, err
		}
	}

	return &order, nil
}

func (o OrderModel) Update(order *Order) error {
	query := `
		UPDATE orders
		SET total = $1, status = $2
		WHERE id = $3
		RETURNING updated_at`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := o.DB.QueryRowContext(ctx, query, order.Total, order.Status, order.ID).Scan(&order.UpdatedAt)
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

func (o OrderModel) GetAllForRestaurant(restaurantID int64, status string, filters Filters) ([]*Order, Metadata, error) {
	query := fmt.Sprintf(`
		SELECT COUNT(*) OVER(), id, user_id, restaurant_id, total, address, created_at, updated_at, status
		FROM orders
		WHERE restaurant_id = $1
		AND (status = $2 OR $2 = '')
		ORDER BY %s %s, id ASC
		LIMIT $3 OFFSET $4`, filters.sortColumn(), filters.sortDirection())

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	rows, err := o.DB.QueryContext(ctx, query, restaurantID, status, filters.limit(), filters.offset())
	if err != nil {
		return nil, Metadata{}, err
	}
	defer rows.Close()

	totalRecords := 0
	var orders []*Order

	for rows.Next() {
		var order Order

		err := rows.Scan(
			&totalRecords,
			&order.ID,
			&order.UserID,
			&order.RestaurantID,
			&order.Total,
			&order.Address,
			&order.CreatedAt,
			&order.UpdatedAt,
			&order.Status,
		)
		if err != nil {
			return nil, Metadata{}, err
		}

		orders = append(orders, &order)
	}

	if err = rows.Err(); err != nil {
		return nil, Metadata{}, err
	}

	metadata := calculateMetadata(totalRecords, filters.Page, filters.PageSize)

	return orders, metadata, nil
}

func (o OrderModel) GetAllForUser(userID int64, status string, filters Filters) ([]*Order, Metadata, error) {
	query := fmt.Sprintf(`
		SELECT COUNT(*) OVER(), id, user_id, restaurant_id, total, address, created_at, updated_at, status
		FROM orders
		WHERE user_id = $1
		AND (status = $2 OR $2 = '')
		ORDER BY %s %s, id ASC
		LIMIT $3 OFFSET $4`, filters.sortColumn(), filters.sortDirection())

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	rows, err := o.DB.QueryContext(ctx, query, userID, status, filters.limit(), filters.offset())
	if err != nil {
		return nil, Metadata{}, err
	}
	defer rows.Close()

	totalRecords := 0
	var orders []*Order

	for rows.Next() {
		var order Order

		err := rows.Scan(
			&totalRecords,
			&order.ID,
			&order.UserID,
			&order.RestaurantID,
			&order.Total,
			&order.Address,
			&order.CreatedAt,
			&order.UpdatedAt,
			&order.Status,
		)
		if err != nil {
			return nil, Metadata{}, err
		}

		orders = append(orders, &order)
	}

	if err = rows.Err(); err != nil {
		return nil, Metadata{}, err
	}

	metadata := calculateMetadata(totalRecords, filters.Page, filters.PageSize)

	return orders, metadata, nil
}
