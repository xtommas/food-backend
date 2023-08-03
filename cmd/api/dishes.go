package main

import (
	"database/sql"
	"errors"
	"fmt"
	"net/http"

	"github.com/xtommas/food-backend/internal/data"
	"github.com/xtommas/food-backend/internal/validator"
)

func (app *application) createDishHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Name        string     `json:"name"`
		Price       data.Price `json:"price"`
		Description string     `json:"description"`
		Category    []string   `json:"category"`
		Photo       string     `json:"photo"`
	}

	err := app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	dish := &data.Dish{
		Name:        input.Name,
		Price:       input.Price,
		Description: input.Description,
		Category:    input.Category,
		Photo:       input.Photo,
	}

	v := validator.New()

	if data.ValidateDish(v, dish); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	err = app.models.Dishes.Insert(dish)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	headers := make(http.Header)
	headers.Set("Location", fmt.Sprintf("/dishes/%d", dish.Id))

	err = app.writeJSON(w, http.StatusCreated, envelope{"dish": dish}, headers)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}

}

func (app *application) showDishHandler(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIdParam(r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}

	dish, err := app.models.Dishes.Get(id)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"dish": dish}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) updateDishHandler(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIdParam(r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}

	dish, err := app.models.Dishes.Get(id)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	var input struct {
		Name        *string     `json:"name"`
		Price       *data.Price `json:"price"`
		Description *string     `json:"description"`
		Category    []string    `json:"category"`
		Photo       *string     `json:"photo"`
		Available   *bool       `json:"available"`
	}

	// read the body data into the input struct
	err = app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	// if the input values are not nil, update the movie record with the new value
	if input.Name != nil {
		dish.Name = *input.Name
	}

	if input.Price != nil {
		dish.Price = *input.Price
	}

	if input.Description != nil {
		dish.Description = *input.Description
	}

	if input.Category != nil {
		dish.Category = input.Category
	}

	if input.Photo != nil {
		dish.Photo = *input.Photo
	}

	if input.Available != nil {
		dish.Available = *input.Available
	}

	// validate the updated movie
	v := validator.New()

	if data.ValidateDish(v, dish); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	// pass the updated record
	err = app.models.Dishes.Update(dish)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrEditConflict):
			app.editConflictResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"dish": dish}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) deleteDishHandler(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIdParam(r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}

	err = app.models.Dishes.Delete(id)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"message": "dish successfully deleted"}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) listDishesHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Name      string
		Category  []string
		Available sql.NullBool
		data.Filters
	}

	v := validator.New()

	qs := r.URL.Query()

	input.Name = app.readString(qs, "name", "")
	input.Available = app.readBool(qs, "available", v)
	input.Category = app.readCSV(qs, "category", []string{})

	input.Filters.Page = app.readInt(qs, "page", 1, v)
	input.Filters.PageSize = app.readInt(qs, "page_size", 20, v)

	input.Filters.Sort = app.readString(qs, "sort", "id")

	input.Filters.SortSafelist = []string{"id", "name", "price", "available", "-id", "-name", "-price", "-available"}

	if data.ValidateFilters(v, input.Filters); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	dishes, err := app.models.Dishes.GetAll(input.Name, input.Category, input.Available, input.Filters)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"dishes": dishes}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}
