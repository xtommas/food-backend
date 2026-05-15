package main

import (
	"expvar"
	"net/http"
)

func (app *application) routes() http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /healthcheck", app.healthcheckHandler)

	// dishes endpoints
	mux.HandleFunc("GET /restaurants/:restaurant_id/dishes", app.requirePermission("dishes:read", app.listDishesHandler))
	mux.HandleFunc("POST /restaurants/:restaurant_id/dishes", app.requirePermission("dishes:write", app.createDishHandler))
	mux.HandleFunc("GET /restaurants/:restaurant_id/dishes/:id", app.requirePermission("dishes:read", app.showDishHandler))
	mux.HandleFunc("PATCH /restaurants/:restaurant_id/dishes/:id", app.requirePermission("dishes:write", app.updateDishHandler))
	mux.HandleFunc("DELETE /restaurants/:restaurant_id/dishes/:id", app.requirePermission("dishes:write", app.deleteDishHandler))
	mux.HandleFunc("POST /restaurants/:restaurant_id/dishes/:id/photo/", app.requirePermission("dishes:write", app.uploadPhotoHandler))
	mux.HandleFunc("GET /restaurants/:restaurant_id/dishes/:id/photo/", app.requirePermission("dishes:read", app.servePhotoHandler))

	mux.HandleFunc("GET /restaurants", app.requirePermission("restaurants:read", app.listRestaurantsHandler))

	// users endpoints
	mux.HandleFunc("POST /users", app.registerUserHandler)
	mux.HandleFunc("PUT /users/activate", app.activateUserHandler)
	mux.HandleFunc("PUT /users/password", app.updateUserPasswordHandler)
	mux.HandleFunc("GET /users/me", app.requireActivatedUser(app.getUserDataHandler))
	mux.HandleFunc("POST /users/me/photo", app.requireActivatedUser(app.uploadUserPhotoHandler))
	mux.HandleFunc("GET /users/me/photo", app.requireActivatedUser(app.serveUserPhotoHandler))
	mux.HandleFunc("PATCH /users/me", app.requireActivatedUser(app.updateUserHandler))

	// orders endpoints
	mux.HandleFunc("POST /restaurants/:restaurant_id/orders", app.requireRole("customer", app.createOrderHandler))
	mux.HandleFunc("GET /restaurants/:restaurant_id/orders", app.requireRole("restaurant", app.getOrdersForRestaurantHandler))
	mux.HandleFunc("GET /users/me/orders", app.requireRole("customer", app.getOrdersForUserHandler))
	mux.HandleFunc("GET /restaurants/:restaurant_id/orders/:order_id", app.requireRole("restaurant", app.getSingleOrderForRestaurantHandler))
	mux.HandleFunc("GET /users/me/orders/:order_id", app.requireRole("customer", app.getSingleOrderForUserHandler))
	mux.HandleFunc("PATCH /restaurants/:restaurant_id/orders/:order_id", app.requireRole("restaurant", app.updateOrderHandler))
	mux.HandleFunc("POST /restaurants/:restaurant_id/orders/:order_id/items", app.requireRole("customer", app.createOrderItemHandler))
	mux.HandleFunc("GET /restaurants/:restaurant_id/orders/:order_id/items", app.requireRole("restaurant", app.getOrderItemsHandler))
	mux.HandleFunc("GET /users/me/orders/:order_id/items", app.requireRole("customer", app.getUserOrderItemsHandler))

	// tokens edpoints
	mux.HandleFunc("POST /tokens/authentication", app.createAuthenticationTokenHandler)
	mux.HandleFunc("POST /tokens/password-reset", app.createPasswordResetTokenHandler)
	mux.HandleFunc("POST /tokens/activation", app.createActivationTokenHandler)

	// debug
	mux.Handle("GET /debug/vars", expvar.Handler())

	// Custom not-found and method-not-allowed responses
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		app.notFoundResponse(w, r)
	})

	return app.metrics(app.recoverPanic(app.enableCORS(app.rateLimit(app.authenticate(mux)))))
}
