package data

import (
	"database/sql"
	"errors"
)

var (
	ErrRecordNotFound = errors.New("record not found")
	ErrEditConflict   = errors.New("edit conflict")
)

type Models struct {
	Dishes      DishModelInterface
	Users       UserModelInterface
	Permissions PermissionModelInterface
	Tokens      TokenModelInterface
	Orders      OrderModelInterface
	OrderItems  OrderItemModelInterface
	Restaurants RestaurantModelInterface
}

func NewModels(db *sql.DB) Models {
	return Models{
		Dishes:      DishModel{DB: db},
		Users:       UserModel{DB: db},
		Permissions: PermissionModel{DB: db},
		Tokens:      TokenModel{DB: db},
		Orders:      OrderModel{DB: db},
		OrderItems:  OrderItemModel{DB: db},
		Restaurants: RestaurantModel{DB: db},
	}
}
