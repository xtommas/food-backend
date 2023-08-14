package main

import (
	"errors"
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

	restaurant, err := app.models.Users.Get(restaurant_id)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	if restaurant.Role != "restaurant" {
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

func (app *application) getOrdersForRestaurantHandler(w http.ResponseWriter, r *http.Request) {
	restaurant_id, err := app.readIdParam(r, "restaurant_id")
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}

	var input struct {
		Status string
		data.Filters
	}

	v := validator.New()

	qs := r.URL.Query()

	input.Status = app.readString(qs, "status", "")

	input.Filters.Page = app.readInt(qs, "page", 1, v)
	input.Filters.PageSize = app.readInt(qs, "page_size", 50, v)

	input.Filters.Sort = app.readString(qs, "sort", "id")

	input.Filters.SortSafelist = []string{"id", "total", "status", "-id", "-total", "-status"}

	if data.ValidateFilters(v, input.Filters); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	user := app.contextGetUser(r)
	if user.Id != restaurant_id {
		app.notFoundResponse(w, r)
		return
	}

	orders, metadata, err := app.models.Orders.GetAllForRestaurant(restaurant_id, input.Status, input.Filters)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"orders": orders, "metadata": metadata}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) getOrdersForUserHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Status string
		data.Filters
	}

	v := validator.New()

	qs := r.URL.Query()

	input.Status = app.readString(qs, "status", "")

	input.Filters.Page = app.readInt(qs, "page", 1, v)
	input.Filters.PageSize = app.readInt(qs, "page_size", 50, v)

	input.Filters.Sort = app.readString(qs, "sort", "id")

	input.Filters.SortSafelist = []string{"id", "total", "status", "-id", "-total", "-status"}

	if data.ValidateFilters(v, input.Filters); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	user := app.contextGetUser(r)

	orders, metadata, err := app.models.Orders.GetAllForUser(user.Id, input.Status, input.Filters)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"orders": orders, "metadata": metadata}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}
