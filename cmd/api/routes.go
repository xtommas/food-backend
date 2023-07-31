package main

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func (app *application) routes() *httprouter.Router {
	router := httprouter.New()

	router.NotFound = http.HandlerFunc(app.notFoundResponse)
	router.MethodNotAllowed = http.HandlerFunc(app.methodNotAllowedResponse)

	router.HandlerFunc(http.MethodGet, "/healthcheck", app.healthcheckHandler)
	router.HandlerFunc(http.MethodPost, "/dishes", app.createDishHandler)
	router.HandlerFunc(http.MethodGet, "/dishes/:id", app.showDishHandler)
	router.HandlerFunc(http.MethodPut, "/dishes/:id", app.updateDishHandler)
	router.HandlerFunc(http.MethodDelete, "/dishes/:id", app.deleteDishHandler)

	return router
}
