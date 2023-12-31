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
	router.HandlerFunc(http.MethodGet, "/restaurants/:restaurant_id/dishes", app.requirePermission("dishes:read", app.listDishesHandler))
	router.HandlerFunc(http.MethodPost, "/restaurants/:restaurant_id/dishes", app.requirePermission("dishes:write", app.createDishHandler))
	router.HandlerFunc(http.MethodGet, "/restaurants/:restaurant_id/dishes/:id", app.requirePermission("dishes:read", app.showDishHandler))
	router.HandlerFunc(http.MethodPatch, "/restaurants/:restaurant_id/dishes/:id", app.requirePermission("dishes:write", app.updateDishHandler))
	router.HandlerFunc(http.MethodDelete, "/restaurants/:restaurant_id/dishes/:id", app.requirePermission("dishes:write", app.deleteDishHandler))
	router.HandlerFunc(http.MethodPost, "/restaurants/:restaurant_id/dishes/:id/photo/", app.requirePermission("dishes:write", app.uploadPhotoHandler))
	router.HandlerFunc(http.MethodGet, "/restaurants/:restaurant_id/dishes/:id/photo/", app.requirePermission("dishes:read", app.servePhotoHandler))

	router.HandlerFunc(http.MethodGet, "/restaurants", app.requirePermission("restaurants:read", app.listRestaurantsHandler))

	// users endpoints
	router.HandlerFunc(http.MethodPost, "/users", app.registerUserHandler)
	router.HandlerFunc(http.MethodPut, "/users/activate", app.activateUserHandler)
	router.HandlerFunc(http.MethodPut, "/users/password", app.updateUserPasswordHandler)
	// router.HandlerFunc(http.MethodPut, "/users/role", app.requireRole("admin", app.updateUserRoleHandler))
	router.HandlerFunc(http.MethodGet, "/users/me", app.requireActivatedUser(app.getUserDataHandler))
	router.HandlerFunc(http.MethodPost, "/users/me/photo", app.requireActivatedUser(app.uploadUserPhotoHandler))
	router.HandlerFunc(http.MethodGet, "/users/me/photo", app.requireActivatedUser(app.serveUserPhotoHandler))
	router.HandlerFunc(http.MethodPatch, "/users/me", app.requireActivatedUser(app.updateUserHandler))

	// orders endpoints
	router.HandlerFunc(http.MethodPost, "/restaurants/:restaurant_id/orders", app.requireRole("customer", app.createOrderHandler))
	router.HandlerFunc(http.MethodGet, "/restaurants/:restaurant_id/orders", app.requireRole("restaurant", app.getOrdersForRestaurantHandler))
	router.HandlerFunc(http.MethodGet, "/users/me/orders", app.requireRole("customer", app.getOrdersForUserHandler))
	router.HandlerFunc(http.MethodGet, "/restaurants/:restaurant_id/orders/:order_id", app.requireRole("restaurant", app.getSingleOrderForRestaurantHandler))
	router.HandlerFunc(http.MethodGet, "/users/me/orders/:order_id", app.requireRole("customer", app.getSingleOrderForUserHandler))
	router.HandlerFunc(http.MethodPatch, "/restaurants/:restaurant_id/orders/:order_id", app.requireRole("restaurant", app.updateOrderHandler))
	router.HandlerFunc(http.MethodPost, "/restaurants/:restaurant_id/orders/:order_id/items", app.requireRole("customer", app.createOrderItemHandler))
	router.HandlerFunc(http.MethodGet, "/restaurants/:restaurant_id/orders/:order_id/items", app.requireRole("restaurant", app.getOrderItemsHandler))
	router.HandlerFunc(http.MethodGet, "/users/me/orders/:order_id/items", app.requireRole("customer", app.getUserOrderItemsHandler))

	// tokens edpoints
	router.HandlerFunc(http.MethodPost, "/tokens/authentication", app.createAuthenticationTokenHandler)
	router.HandlerFunc(http.MethodPost, "/tokens/password-reset", app.createPasswordResetTokenHandler)
	router.HandlerFunc(http.MethodPost, "/tokens/activation", app.createActivationTokenHandler)

	// debug
	router.Handler(http.MethodGet, "/debug/vars", expvar.Handler())

	return app.metrics(app.recoverPanic(app.enableCORS(app.rateLimit(app.authenticate(router)))))
}
