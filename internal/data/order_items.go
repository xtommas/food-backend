package data

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/xtommas/food-backend/internal/validator"
)

type OrderItem struct {
	Id       int64 `json:"id"`
	Order_id int64 `json:"order_id"`
	Dish_id  int64 `json:"dish_id"`
	Quantity int   `json:"quantity"`
	Subtotal Price `json:"subtotal"`
}

func ValidateQuantity(v *validator.Validator, quantity int) {
	v.Check(quantity > 0, "address", "must be a positive number")
}

func ValidateOrderItem(v *validator.Validator, item *OrderItem) {
	ValidateQuantity(v, item.Quantity)
}

type OrderItemModel struct {
	DB *sql.DB
}

func (i OrderItemModel) Insert(orderItem *OrderItem) error {
	query := `
		INSERT INTO order_items (order_id, dish_id, quantity, subtotal)
		VALUES ($1, $2, $3, $4)
		RETURNING id`

	args := []interface{}{orderItem.Order_id, orderItem.Dish_id, orderItem.Quantity, orderItem.Subtotal}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	return i.DB.QueryRowContext(ctx, query, args...).Scan(&orderItem.Id)
}

func (i OrderItemModel) Update(orderItem *OrderItem) error {
	query := `
		UPDATE order_items
		SET quantity = $1, subtotal = $2
		WHERE id = $3`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	result, err := i.DB.ExecContext(ctx, query, orderItem.Quantity, orderItem.Subtotal, orderItem.Id)
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

func (i OrderItemModel) GetForOrder(order_id int64) ([]*OrderItem, error) {
	query := `
		SELECT id, order_id, dish_id, quantity, subtotal
		FROM order_items
		WHERE order_id = $1`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	rows, err := i.DB.QueryContext(ctx, query, order_id)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	orderItems := []*OrderItem{}

	for rows.Next() {
		var item OrderItem

		err := rows.Scan(
			&item.Id,
			&item.Order_id,
			&item.Dish_id,
			&item.Quantity,
			&item.Subtotal,
		)

		if err != nil {
			return nil, err
		}

		orderItems = append(orderItems, &item)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return orderItems, nil
}
