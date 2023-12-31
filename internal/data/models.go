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
	Dishes      DishModel
	Users       UserModel
	Permissions PermissionModel
	Tokens      TokenModel
	Orders      OrderModel
	OrderItems  OrderItemModel
}

func NewModels(db *sql.DB) Models {
	return Models{
		Dishes:      DishModel{DB: db},
		Users:       UserModel{DB: db},
		Permissions: PermissionModel{DB: db},
		Tokens:      TokenModel{DB: db},
		Orders:      OrderModel{DB: db},
		OrderItems:  OrderItemModel{DB: db},
	}
}
