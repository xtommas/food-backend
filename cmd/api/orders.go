package main

import (
	"fmt"
	"net/http"

	"github.com/xtommas/food-backend/internal/data"
	"github.com/xtommas/food-backend/internal/validator"
)

func (app *application) createOrderHandler(w http.ResponseWriter, r *http.Request) {
	restaurant_id, err := app.readIdParam(r, "restaurant_id")
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}

	user := app.contextGetUser(r)

	var input struct {
		Total   data.Price `json:"price"`
		Address string     `json:"address"`
	}

	err = app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	order := &data.Order{
		User_id:       user.Id,
		Restaurant_id: restaurant_id,
		Total:         input.Total,
		Address:       input.Address,
		Status:        "created",
	}

	v := validator.New()

	if data.ValidateOrder(v, order); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	err = app.models.Orders.Insert(order)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	headers := make(http.Header)
	headers.Set("Location", fmt.Sprintf("/restaurant/%d/orders/orders/%d", order.Restaurant_id, order.Id))

	err = app.writeJSON(w, http.StatusCreated, envelope{"order": order}, headers)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}
