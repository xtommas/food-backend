package main

import (
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

	fmt.Fprintf(w, "%+v\n", input)
}

func (app *application) showDishHandler(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIdParam(r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}

	dish := data.Dish{
		Id:          id,
		Name:        "Pizza",
		Price:       1500,
		Category:    []string{"Pizzas"},
		Description: "Pizza de Muzzarella de 8 porciones",
		Photo:       "",
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"dish": dish}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}
