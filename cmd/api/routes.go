package main

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func (app *application) routes() http.Handler {
	router := httprouter.New()

	router.NotFound = http.HandlerFunc(app.notFoundResponse)
	router.MethodNotAllowed = http.HandlerFunc(app.methodNotAllowedResponse)

	router.HandlerFunc(http.MethodGet, "/healthcheck", app.healthcheckHandler)

	router.HandlerFunc(http.MethodGet, "/dishes", app.listDishesHandler)
	router.HandlerFunc(http.MethodPost, "/dishes", app.createDishHandler)
	router.HandlerFunc(http.MethodGet, "/dishes/:id", app.showDishHandler)
	router.HandlerFunc(http.MethodPatch, "/dishes/:id", app.updateDishHandler)
	router.HandlerFunc(http.MethodDelete, "/dishes/:id", app.deleteDishHandler)

	router.HandlerFunc(http.MethodPost, "/users", app.registerUserHandler)
	router.HandlerFunc(http.MethodPut, "/users/activated", app.activateUserHandler)

	return app.recoverPanic(app.rateLimit(router))
}
