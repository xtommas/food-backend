package main

import (
	"expvar"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func (app *application) routes() http.Handler {
	router := httprouter.New()

	router.NotFound = http.HandlerFunc(app.notFoundResponse)
	router.MethodNotAllowed = http.HandlerFunc(app.methodNotAllowedResponse)

	router.HandlerFunc(http.MethodGet, "/healthcheck", app.healthcheckHandler)

	// dishes endpoints
	router.HandlerFunc(http.MethodGet, "/dishes", app.requirePermission("dishes:read", app.listDishesHandler))
	router.HandlerFunc(http.MethodPost, "/dishes", app.requirePermission("dishes:write", app.createDishHandler))
	router.HandlerFunc(http.MethodGet, "/dishes/:id", app.requirePermission("dishes:read", app.showDishHandler))
	router.HandlerFunc(http.MethodPatch, "/dishes/:id", app.requirePermission("dishes:write", app.updateDishHandler))
	router.HandlerFunc(http.MethodDelete, "/dishes/:id", app.requirePermission("dishes:write", app.deleteDishHandler))
	router.HandlerFunc(http.MethodPost, "/dishes/photo/:id", app.requirePermission("dishes:write", app.uploadPhotoHandler))
	router.HandlerFunc(http.MethodGet, "/dishes/:id/photo/", app.requirePermission("dishes:read", app.servePhotoHandler))

	// users endpoints
	router.HandlerFunc(http.MethodPost, "/users", app.registerUserHandler)
	router.HandlerFunc(http.MethodPut, "/users/activated", app.activateUserHandler)
	router.HandlerFunc(http.MethodPut, "/users/password", app.updateUserPasswordHandler)
	router.HandlerFunc(http.MethodPut, "/users/role", app.requireAdmin(app.updateUserRoleHandler))
	router.HandlerFunc(http.MethodGet, "/users/me", app.requireActivatedUser(app.getUserDataHandler))

	// tokens edpoints
	router.HandlerFunc(http.MethodPost, "/tokens/authentication", app.createAuthenticationTokenHandler)
	router.HandlerFunc(http.MethodPost, "/tokens/password-reset", app.createPasswordResetTokenHandler)
	router.HandlerFunc(http.MethodPost, "/tokens/activation", app.createActivationTokenHandler)

	// debug
	router.Handler(http.MethodGet, "/debug/vars", expvar.Handler())

	return app.metrics(app.recoverPanic(app.enableCORS(app.rateLimit(app.authenticate(router)))))
}
