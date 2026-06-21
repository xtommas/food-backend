package main

import "time"

type orderCreatedEvent struct {
	OrderID      int64     `json:"order_id"`
	UserID       int64     `json:"user_id"`
	RestaurantID int64     `json:"restaurant_id"`
	Address      string    `json:"address"`
	Status       string    `json:"status"`
	CreatedAt    time.Time `json:"created_at"`
}

type orderStatusChangedEvent struct {
	OrderID      int64     `json:"order_id"`
	RestaurantID int64     `json:"restaurant_id"`
	UserID       int64     `json:"user_id"`
	FromStatus   string    `json:"from_status"`
	ToStatus     string    `json:"to_status"`
	ChangedAt    time.Time `json:"changed_at"`
}
