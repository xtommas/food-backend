package main

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func (app *application) routes() *httprouter.Router {
	router := httprouter.New()

	router.HandlerFunc(http.MethodGet, "/healthcheck", app.healthcheckHandler)
	router.HandlerFunc(http.MethodPost, "/dishes", app.createDishHandler)
	router.HandlerFunc(http.MethodGet, "/dishes/:id", app.showDishHandler)

	return router
}
