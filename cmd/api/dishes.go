package main

import (
	"fmt"
	"net/http"

	"github.com/xtommas/food-backend/internal/data"
)

func (app *application) createDishHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "create a new dish")
}

func (app *application) showDishHandler(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIdParam(r)
	if err != nil {
		http.NotFound(w, r)
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

	err = app.writeJSON(w, http.StatusOK, dish, nil)
	if err != nil {
		app.logger.Println(err)
		http.Error(w, "The server encountered a problem and could not process your request", http.StatusInternalServerError)
	}
}
