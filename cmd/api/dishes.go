package main

import (
	"database/sql"
	"errors"
	"fmt"
	"io"
	"mime"
	"net/http"
	"os"
	"path/filepath"
	"strconv"

	"github.com/xtommas/food-backend/internal/data"
	"github.com/xtommas/food-backend/internal/validator"
)

func (app *application) createDishHandler(w http.ResponseWriter, r *http.Request) {
	restaurantID, err := app.readIdParam(r, "restaurant_id")
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}

	_, err = app.models.Restaurants.Get(restaurantID)
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
		Name        string   `json:"name"`
		Price       int64    `json:"price"`
		Description string   `json:"description"`
		Categories  []string `json:"categories"`
	}

	err = app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	dish := &data.Dish{
		Name:         input.Name,
		RestaurantID: restaurantID,
		Price:        input.Price,
		Description:  input.Description,
		Categories:   input.Categories,
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
	headers.Set("Location", fmt.Sprintf("/restaurants/%d/dishes/%d", dish.RestaurantID, dish.ID))

	err = app.writeJSON(w, http.StatusCreated, envelope{"dish": dish}, headers)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) showDishHandler(w http.ResponseWriter, r *http.Request) {
	restaurantID, err := app.readIdParam(r, "restaurant_id")
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}

	id, err := app.readIdParam(r, "id")
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}

	_, err = app.models.Restaurants.Get(restaurantID)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
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

	if dish.RestaurantID != restaurantID {
		app.notFoundResponse(w, r)
		return
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"dish": dish}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) updateDishHandler(w http.ResponseWriter, r *http.Request) {
	restaurantID, err := app.readIdParam(r, "restaurant_id")
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}

	id, err := app.readIdParam(r, "id")
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

	if dish.RestaurantID != restaurantID {
		app.notFoundResponse(w, r)
		return
	}

	var input struct {
		Name        *string  `json:"name"`
		Price       *int64   `json:"price"`
		Description *string  `json:"description"`
		Categories  []string `json:"categories"`
		Photo       *string  `json:"photo"`
		Available   *bool    `json:"available"`
	}

	err = app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	if input.Name != nil {
		dish.Name = *input.Name
	}
	if input.Price != nil {
		dish.Price = *input.Price
	}
	if input.Description != nil {
		dish.Description = *input.Description
	}
	if input.Categories != nil {
		dish.Categories = input.Categories
	}
	if input.Photo != nil {
		dish.Photo = *input.Photo
	}
	if input.Available != nil {
		dish.Available = *input.Available
	}

	v := validator.New()

	if data.ValidateDish(v, dish); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

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
	restaurantID, err := app.readIdParam(r, "restaurant_id")
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}

	id, err := app.readIdParam(r, "id")
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

	if dish.RestaurantID != restaurantID {
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

	if dish.Photo != "" {
		if _, err := os.Stat(dish.Photo); err == nil {
			if err := os.Remove(dish.Photo); err != nil {
				app.serverErrorResponse(w, r, err)
				return
			}
		}
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"message": "dish successfully deleted"}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) listDishesHandler(w http.ResponseWriter, r *http.Request) {
	restaurantID, err := app.readIdParam(r, "restaurant_id")
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}

	_, err = app.models.Restaurants.Get(restaurantID)
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
		Name       string
		Categories []string
		Available  sql.NullBool
		data.Filters
	}

	v := validator.New()

	qs := r.URL.Query()

	input.Name = app.readString(qs, "name", "")
	input.Available = app.readBool(qs, "available", v)
	input.Categories = app.readCSV(qs, "categories", []string{})

	input.Filters.Page = app.readInt(qs, "page", 1, v)
	input.Filters.PageSize = app.readInt(qs, "page_size", 20, v)
	input.Filters.Sort = app.readString(qs, "sort", "id")
	input.Filters.SortSafelist = []string{"id", "name", "price", "available", "-id", "-name", "-price", "-available"}

	if data.ValidateFilters(v, input.Filters); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	dishes, metadata, err := app.models.Dishes.GetAllForRestaurant(restaurantID, input.Name, input.Categories, input.Available, input.Filters)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"dishes": dishes, "metadata": metadata}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) uploadPhotoHandler(w http.ResponseWriter, r *http.Request) {
	restaurantID, err := app.readIdParam(r, "restaurant_id")
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}

	id, err := app.readIdParam(r, "id")
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

	if dish.RestaurantID != restaurantID {
		app.notFoundResponse(w, r)
		return
	}

	err = r.ParseMultipartForm(10 << 20)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	fileName := strconv.FormatInt(dish.ID, 10) + ".jpg"
	folder := "images/dishes/"

	dish.Photo, err = app.storeImage(w, r, folder, fileName)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

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

func (app *application) servePhotoHandler(w http.ResponseWriter, r *http.Request) {
	restaurantID, err := app.readIdParam(r, "restaurant_id")
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}

	id, err := app.readIdParam(r, "id")
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

	if dish.RestaurantID != restaurantID {
		app.notFoundResponse(w, r)
		return
	}

	imageFile, err := os.Open(dish.Photo)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}
	defer imageFile.Close()

	contentType := mime.TypeByExtension(filepath.Ext(dish.Photo))
	w.Header().Set("Content-Type", contentType)

	_, err = io.Copy(w, imageFile)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}
