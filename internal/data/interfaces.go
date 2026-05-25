package data

import (
	"database/sql"
	"time"
)

type DishModelInterface interface {
	Insert(dish *Dish) error
	Get(id int64) (*Dish, error)
	Update(dish *Dish) error
	Delete(id int64) error
	GetAllForRestaurant(restaurantID int64, name string, categories []string, available sql.NullBool, filters Filters) ([]*Dish, Metadata, error)
}

type OrderItemModelInterface interface {
	Insert(orderItem *OrderItem) error
	InsertFromDish(orderId int64, dish *Dish, quantity int) (*OrderItem, error)
	Update(orderItem *OrderItem) error
	GetForOrder(orderID int64) ([]*OrderItem, error)
	DeleteForOrder(orderID int64) error
}

type OrderModelInterface interface {
	Insert(order *Order) error
	GetForRestaurant(id int64, restaurantID int64) (*Order, error)
	GetForUser(id int64, userID int64) (*Order, error)
	Update(order *Order) error
	GetAllForRestaurant(restaurantID int64, status string, filters Filters) ([]*Order, Metadata, error)
	GetAllForUser(userID int64, status string, filters Filters) ([]*Order, Metadata, error)
}

type PermissionModelInterface interface {
	GetAllForUser(userId int64) (Permissions, error)
	AddForUser(userId int64, codes ...string) error
	DeleteForUser(userId int64, code string) error
}

type RestaurantModelInterface interface {
	Insert(restaurant *Restaurant) error
	Get(id int64) (*Restaurant, error)
	Update(restaurant *Restaurant) error
	Delete(id int64) error
	GetAll() ([]*Restaurant, error)
	GetStaff(restaurantID int64) ([]*User, error)
	AddStaff(restaurantID, userID int64, role string) error
	RemoveStaff(restaurantID, userID int64) error
	IsStaff(restaurantID, userID int64) (bool, error)
	GetStaffRole(restaurantID, userID int64) (string, error)
}

type TokenModelInterface interface {
	New(userID int64, ttl time.Duration, scope string) (*Token, error)
	Insert(token *Token) error
	DeleteAllForUser(scope string, userID int64) error
}

type UserModelInterface interface {
	Insert(user *User) error
	GetByEmail(email string) (*User, error)
	Update(user *User) error
	GetForToken(tokenScope, tokenPlaintext string) (*User, error)
	Get(id int64) (*User, error)
}
