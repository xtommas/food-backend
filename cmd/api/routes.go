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

	// GET /restaurants
	// GET /restaurants/:id
	// GET /restaurants/:id/dishes gets all the restaurant's dishes
	// POST /restaurants/:id/dishes creates a new dish associated with the restaurant that makes the request
	// GET /restaurants/:id/dishes/:id get a specific dish
	// PATCH /restaurants/:id/dishes/:id
	// POST /restaurants/:id/dishes/:id/photo upload the photo for the dish
	// GET /restaurants/:id/dishes/:id/photo get the photo for the dish

	// POST /restaurants/ create a restaurant
	// POST /restaurants/activated activate a restaurant account
	// PUT /restaurants/password change the password
	// POST /restaurants/:id/photo
	// GET /restaurants/:id/photo

	// POST /users create a user
	// POST /users/activated activate a user account
	// PUT /users/password change the password
	// GET /users/me get currently logged in user info

	// POST /users/tokens/authentication get authentication token for user
	// POST /users/tokens/password-reset get password reset token for user
	// POST /users/tokens/activation get activation token for user

	// POST /restaurants/tokens/authentication get authentication token for restaurant
	// POST /restaurants/tokens/password-reset get password reset token for restaurant
	// POST /restaurants/tokens/activation get activation token for restaurant

	// debug
	router.Handler(http.MethodGet, "/debug/vars", expvar.Handler())

	return app.metrics(app.recoverPanic(app.enableCORS(app.rateLimit(app.authenticate(router)))))
}
