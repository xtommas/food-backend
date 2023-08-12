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
	Id            int64     `json:"id"`
	User_id       int64     `json:"user_id"`
	Restaurant_id int64     `json:"restaurant_id"`
	Total         Price     `json:"total"`
	Address       string    `json:"address"`
	Created_at    time.Time `json:"created_at"`
	Status        string    `json:"status"`
}

func ValidateAddress(v *validator.Validator, address string) {
	v.Check(address != "", "address", "must be provided")
}

func ValidateStatus(v *validator.Validator, status string) {
	v.Check(status != "", "status", "must be provided")
	v.Check(status == "created" || status == "in progress" || status == "ready" || status == "delivered" || status == "cancelled", "status", "invalid status")
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
		RETURNING id`

	args := []interface{}{order.User_id, order.Restaurant_id, order.Total, order.Address, order.Status}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	return o.DB.QueryRowContext(ctx, query, args...).Scan(&order.Id)
}

func (o OrderModel) Get(id int64) (*Order, error) {
	if id < 1 {
		return nil, ErrRecordNotFound
	}

	query := `
		SELECT id, user_id, restaurant_id, total, address, created_at, status
		FROM orders
		WHERE id = $1`

	var order Order

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := o.DB.QueryRowContext(ctx, query, id).Scan(
		&order.Id,
		&order.User_id,
		&order.Restaurant_id,
		&order.Total,
		&order.Address,
		&order.Created_at,
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
		WHERE id = $3`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	result, err := o.DB.ExecContext(ctx, query, order.Total, order.Status, order.Id)
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

func (o OrderModel) GetAllForRestaurant(restaurant_id int64, status string, filters Filters) ([]*Order, Metadata, error) {
	query := fmt.Sprintf(`
		SELECT COUNT(*) OVER(), id, user_id, restaurant_id, total, address, created_at, status
		FROM orders
		WHERE restaurant_id = $1
		AND (status = $2 OR $2 = '')
		ORDER BY %s %s, id ASC
		LIMIT $3 OFFSET $4`, filters.sortColumn(), filters.sortDirection())

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	rows, err := o.DB.QueryContext(ctx, query, restaurant_id, status, filters.limit(), filters.offset())
	if err != nil {
		return nil, Metadata{}, err
	}

	defer rows.Close()

	totalRecords := 0
	orders := []*Order{}

	for rows.Next() {
		var order Order

		err := rows.Scan(
			&totalRecords,
			&order.Id,
			&order.User_id,
			&order.Restaurant_id,
			&order.Total,
			&order.Address,
			&order.Created_at,
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

func (o OrderModel) GetAllForUser(user_id int64, status string, filters Filters) ([]*Order, Metadata, error) {
	query := fmt.Sprintf(`
		SELECT COUNT(*) OVER(), id, user_id, restaurant_id, total, address, created_at, status
		FROM orders
		WHERE user_id = $1
		AND (status = $2 OR $2 = '')
		ORDER BY %s %s, id ASC
		LIMIT $3 OFFSET $4`, filters.sortColumn(), filters.sortDirection())

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	rows, err := o.DB.QueryContext(ctx, query, user_id, status, filters.limit(), filters.offset())
	if err != nil {
		return nil, Metadata{}, err
	}

	defer rows.Close()

	totalRecords := 0
	orders := []*Order{}

	for rows.Next() {
		var order Order

		err := rows.Scan(
			&totalRecords,
			&order.Id,
			&order.User_id,
			&order.Restaurant_id,
			&order.Total,
			&order.Address,
			&order.Created_at,
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
