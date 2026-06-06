package main

import (
	"expvar"
	"net/http"
)

func (app *application) routes() http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /healthcheck", app.healthcheckHandler)

	// dishes endpoints
	mux.HandleFunc("GET /restaurants/{restaurant_id}/dishes", app.requirePermission("dishes:read", app.listDishesHandler))
	mux.HandleFunc("POST /restaurants/{restaurant_id}/dishes", app.requireRestaurantStaff(app.createDishHandler))
	mux.HandleFunc("GET /restaurants/{restaurant_id}/dishes/{id}", app.requirePermission("dishes:read", app.showDishHandler))
	mux.HandleFunc("PATCH /restaurants/{restaurant_id}/dishes/{id}", app.requireRestaurantStaff(app.updateDishHandler))
	mux.HandleFunc("DELETE /restaurants/{restaurant_id}/dishes/{id}", app.requireRestaurantStaff(app.deleteDishHandler))
	mux.HandleFunc("POST /restaurants/{restaurant_id}/dishes/{id}/photo/", app.requireRestaurantStaff(app.uploadPhotoHandler))
	mux.HandleFunc("GET /restaurants/{restaurant_id}/dishes/{id}/photo/", app.requirePermission("dishes:read", app.servePhotoHandler))

	// restaurants endpoints
	mux.HandleFunc("GET /restaurants", app.requirePermission("restaurants:read", app.listRestaurantsHandler))
	mux.HandleFunc("POST /restaurants", app.requireAdmin(app.createRestaurantHandler))
	mux.HandleFunc("GET /restaurants/{restaurant_id}", app.requirePermission("restaurants:read", app.showRestaurantHandler))
	mux.HandleFunc("PATCH /restaurants/{restaurant_id}", app.requireAdmin(app.updateRestaurantHandler))
	mux.HandleFunc("DELETE /restaurants/{restaurant_id}", app.requireAdmin(app.deleteRestaurantHandler))

	// restaurant staff endpoints
	mux.HandleFunc("GET /restaurants/{restaurant_id}/staff", app.requireRestaurantOwner(app.listRestaurantStaffHandler))
	mux.HandleFunc("POST /restaurants/{restaurant_id}/staff", app.requireRestaurantOwner(app.addRestaurantStaffHandler))
	mux.HandleFunc("DELETE /restaurants/{restaurant_id}/staff/{user_id}", app.requireRestaurantOwner(app.removeRestaurantStaffHandler))

	// users endpoints
	mux.HandleFunc("POST /users", app.registerUserHandler)
	mux.HandleFunc("PUT /users/activate", app.activateUserHandler)
	mux.HandleFunc("PUT /users/password", app.updateUserPasswordHandler)
	mux.HandleFunc("GET /users/me", app.requireActivatedUser(app.getUserDataHandler))
	mux.HandleFunc("POST /users/me/photo", app.requireActivatedUser(app.uploadUserPhotoHandler))
	mux.HandleFunc("GET /users/me/photo", app.requireActivatedUser(app.serveUserPhotoHandler))
	mux.HandleFunc("PATCH /users/me", app.requireActivatedUser(app.updateUserHandler))

	// admin endpoints
	mux.HandleFunc("POST /admin/promote", app.requireAdmin(app.promoteUserHandler))

	// orders endpoints
	mux.HandleFunc("POST /restaurants/{restaurant_id}/orders", app.requireActivatedUser(app.createOrderHandler))
	mux.HandleFunc("GET /restaurants/{restaurant_id}/orders", app.requireRestaurantStaff(app.getOrdersForRestaurantHandler))
	mux.HandleFunc("GET /users/me/orders", app.requireActivatedUser(app.getOrdersForUserHandler))
	mux.HandleFunc("GET /restaurants/{restaurant_id}/orders/{order_id}", app.requireRestaurantStaff(app.getSingleOrderForRestaurantHandler))
	mux.HandleFunc("GET /users/me/orders/{order_id}", app.requireActivatedUser(app.getSingleOrderForUserHandler))
	mux.HandleFunc("PATCH /restaurants/{restaurant_id}/orders/{order_id}", app.requireRestaurantStaff(app.updateOrderHandler))
	mux.HandleFunc("POST /restaurants/{restaurant_id}/orders/{order_id}/items", app.requireActivatedUser(app.createOrderItemHandler))
	mux.HandleFunc("GET /restaurants/{restaurant_id}/orders/{order_id}/items", app.requireRestaurantStaff(app.getOrderItemsHandler))
	mux.HandleFunc("GET /users/me/orders/{order_id}/items", app.requireActivatedUser(app.getUserOrderItemsHandler))

	// tokens endpoints
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
