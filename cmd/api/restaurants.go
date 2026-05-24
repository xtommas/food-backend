package main

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/xtommas/food-backend/internal/data"
	"github.com/xtommas/food-backend/internal/validator"
)

func (app *application) listRestaurantsHandler(w http.ResponseWriter, r *http.Request) {
	restaurants, err := app.models.Restaurants.GetAll()
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"restaurants": restaurants}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) showRestaurantHandler(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(r.PathValue("restaurant_id"), 10, 64)
	if err != nil || id < 1 {
		app.notFoundResponse(w, r)
		return
	}

	restaurant, err := app.models.Restaurants.Get(id)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"restaurant": restaurant}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) createRestaurantHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Name      string  `json:"name"`
		Photo     string  `json:"photo"`
		Address   string  `json:"address"`
		City      string  `json:"city"`
		State     string  `json:"state"`
		Province  string  `json:"province"`
		Country   string  `json:"country"`
		Latitude  float64 `json:"latitude"`
		Longitude float64 `json:"longitude"`
	}

	err := app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	restaurant := &data.Restaurant{
		Name:      input.Name,
		Photo:     input.Photo,
		Address:   input.Address,
		City:      input.City,
		State:     input.State,
		Province:  input.Province,
		Country:   input.Country,
		Latitude:  input.Latitude,
		Longitude: input.Longitude,
	}

	v := validator.New()
	if data.ValidateRestaurant(v, restaurant); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	err = app.models.Restaurants.Insert(restaurant)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	err = app.writeJSON(w, http.StatusCreated, envelope{"restaurant": restaurant}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) updateRestaurantHandler(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(r.PathValue("restaurant_id"), 10, 64)
	if err != nil || id < 1 {
		app.notFoundResponse(w, r)
		return
	}

	restaurant, err := app.models.Restaurants.Get(id)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	var input struct {
		Name      *string  `json:"name"`
		Photo     *string  `json:"photo"`
		Address   *string  `json:"address"`
		City      *string  `json:"city"`
		State     *string  `json:"state"`
		Province  *string  `json:"province"`
		Country   *string  `json:"country"`
		Latitude  *float64 `json:"latitude"`
		Longitude *float64 `json:"longitude"`
	}

	err = app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	if input.Name != nil {
		restaurant.Name = *input.Name
	}
	if input.Photo != nil {
		restaurant.Photo = *input.Photo
	}
	if input.Address != nil {
		restaurant.Address = *input.Address
	}
	if input.City != nil {
		restaurant.City = *input.City
	}
	if input.State != nil {
		restaurant.State = *input.State
	}
	if input.Province != nil {
		restaurant.Province = *input.Province
	}
	if input.Country != nil {
		restaurant.Country = *input.Country
	}
	if input.Latitude != nil {
		restaurant.Latitude = *input.Latitude
	}
	if input.Longitude != nil {
		restaurant.Longitude = *input.Longitude
	}

	v := validator.New()
	if data.ValidateRestaurant(v, restaurant); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	err = app.models.Restaurants.Update(restaurant)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrEditConflict):
			app.editConflictResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"restaurant": restaurant}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) deleteRestaurantHandler(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(r.PathValue("restaurant_id"), 10, 64)
	if err != nil || id < 1 {
		app.notFoundResponse(w, r)
		return
	}

	err = app.models.Restaurants.Delete(id)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"message": "restaurant successfully deleted"}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

// Restaurant staff handlers

func (app *application) listRestaurantStaffHandler(w http.ResponseWriter, r *http.Request) {
	restaurantID, err := strconv.ParseInt(r.PathValue("restaurant_id"), 10, 64)
	if err != nil || restaurantID < 1 {
		app.notFoundResponse(w, r)
		return
	}

	_, err = app.models.Restaurants.Get(restaurantID)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	staff, err := app.models.Restaurants.GetStaff(restaurantID)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"staff": staff}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) addRestaurantStaffHandler(w http.ResponseWriter, r *http.Request) {
	restaurantID, err := strconv.ParseInt(r.PathValue("restaurant_id"), 10, 64)
	if err != nil || restaurantID < 1 {
		app.notFoundResponse(w, r)
		return
	}

	var input struct {
		Email string `json:"email"`
		Role  string `json:"role"`
	}

	err = app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	v := validator.New()
	data.ValidateEmail(v, input.Email)
	v.Check(
		validator.PermittedValue(input.Role, "owner", "staff"),
		"role", "must be 'owner' or 'staff'",
	)
	if !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	_, err = app.models.Restaurants.Get(restaurantID)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	user, err := app.models.Users.GetByEmail(input.Email)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			v.AddError("email", "no user found with this email address")
			app.failedValidationResponse(w, r, v.Errors)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	err = app.models.Restaurants.AddStaff(restaurantID, user.Id, input.Role)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	err = app.models.Permissions.AddForUser(user.Id, "dishes:write", "orders:read", "orders:write")
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"message": "staff member successfully added"}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) removeRestaurantStaffHandler(w http.ResponseWriter, r *http.Request) {
	restaurantID, err := strconv.ParseInt(r.PathValue("restaurant_id"), 10, 64)
	if err != nil || restaurantID < 1 {
		app.notFoundResponse(w, r)
		return
	}

	userID, err := strconv.ParseInt(r.PathValue("user_id"), 10, 64)
	if err != nil || userID < 1 {
		app.notFoundResponse(w, r)
		return
	}

	err = app.models.Restaurants.RemoveStaff(restaurantID, userID)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"message": "staff member successfully removed"}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}
