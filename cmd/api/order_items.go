package main

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/xtommas/food-backend/internal/data"
	"github.com/xtommas/food-backend/internal/validator"
)

func (app *application) createOrderItemHandler(w http.ResponseWriter, r *http.Request) {
	restaurant_id, err := app.readIdParam(r, "restaurant_id")
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}

	order_id, err := app.readIdParam(r, "order_id")
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}

	var input struct {
		Dish_id  int64 `json:"dish_id"`
		Quantity int   `json:"quantity"`
	}

	err = app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	dish, err := app.models.Dishes.Get(input.Dish_id)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	if dish.Restaurant_id != restaurant_id {
		app.notFoundResponse(w, r)
		return
	}

	subtotal := dish.Price * data.Price(input.Quantity)

	order_item := &data.OrderItem{
		Order_id: order_id,
		Dish_id:  input.Dish_id,
		Quantity: input.Quantity,
		Subtotal: subtotal,
	}

	v := validator.New()

	if data.ValidateOrderItem(v, order_item); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	user := app.contextGetUser(r)

	order, err := app.models.Orders.GetForUser(order_id, user.Id)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	if order.Status == "delivered" {
		app.editConflictResponse(w, r)
		return
	}

	if order.Restaurant_id != restaurant_id {
		app.badRequestResponse(w, r, errors.New("invalid restaurant"))
		return
	}

	err = app.models.OrderItems.Insert(order_item)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	order.Total += order_item.Subtotal

	err = app.models.Orders.Update(order)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrEditConflict):
			app.editConflictResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	headers := make(http.Header)
	headers.Set("Location", fmt.Sprintf("/restaurant/%d/orders/orders/%d/items/%d", restaurant_id, order_id, order_item.Id))

	err = app.writeJSON(w, http.StatusCreated, envelope{"order_item": order_item}, headers)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) getOrderItemsHandler(w http.ResponseWriter, r *http.Request) {
	restaurant_id, err := app.readIdParam(r, "restaurant_id")
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}

	order_id, err := app.readIdParam(r, "order_id")
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}

	restaurant := app.contextGetUser(r)

	if restaurant.Id != restaurant_id {
		app.notPermittedResponse(w, r)
		return
	}

	order, err := app.models.Orders.GetForRestaurant(order_id, restaurant_id)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	items, err := app.models.OrderItems.GetForOrder(order.Id)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	type order_item struct {
		Name     string     `json:"dish"`
		Quantity int        `json:"quantity"`
		Subtotal data.Price `json:"subtotal"`
	}

	order_items := []order_item{}

	for _, item := range items {
		dish, err := app.models.Dishes.Get(item.Dish_id)
		if err != nil {
			switch {
			case errors.Is(err, data.ErrRecordNotFound):
				app.notFoundResponse(w, r)
			default:
				app.serverErrorResponse(w, r, err)
			}
			return
		}

		order_items = append(order_items, order_item{Name: dish.Name, Quantity: item.Quantity, Subtotal: item.Subtotal})
	}

	err = app.writeJSON(w, http.StatusCreated, envelope{"items": order_items}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) getUserOrderItemsHandler(w http.ResponseWriter, r *http.Request) {
	order_id, err := app.readIdParam(r, "order_id")
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}

	user := app.contextGetUser(r)

	order, err := app.models.Orders.GetForUser(order_id, user.Id)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	items, err := app.models.OrderItems.GetForOrder(order.Id)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	type order_item struct {
		Name     string     `json:"dish"`
		Quantity int        `json:"quantity"`
		Subtotal data.Price `json:"subtotal"`
	}

	order_items := []order_item{}

	for _, item := range items {
		dish, err := app.models.Dishes.Get(item.Dish_id)
		if err != nil {
			switch {
			case errors.Is(err, data.ErrRecordNotFound):
				app.notFoundResponse(w, r)
			default:
				app.serverErrorResponse(w, r, err)
			}
			return
		}

		order_items = append(order_items, order_item{Name: dish.Name, Quantity: item.Quantity, Subtotal: item.Subtotal})
	}

	err = app.writeJSON(w, http.StatusCreated, envelope{"items": order_items}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}
