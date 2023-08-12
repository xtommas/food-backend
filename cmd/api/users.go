package main

import (
	"errors"
	"io"
	"mime"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/pascaldekloe/jwt"
	"github.com/xtommas/food-backend/internal/data"
	"github.com/xtommas/food-backend/internal/validator"
)

func (app *application) registerUserHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Name     string `json:"name"`
		Email    string `json:"email"`
		Password string `json:"password"`
		Role     string `json:"role"`
	}

	err := app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	user := &data.User{
		Name:      input.Name,
		Photo:     "",
		Email:     input.Email,
		Activated: false,
		Role:      input.Role,
	}

	err = user.Password.Set(input.Password)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	v := validator.New()

	if data.ValidateUser(v, user); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	err = app.models.Users.Insert(user)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrDuplicateEmail):
			v.AddError("email", "a user with this email address already exists")
			app.failedValidationResponse(w, r, v.Errors)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	if user.Role == "restaurant" || user.Role == "admin" {
		err = app.models.Permissions.AddForUser(user.Id, "dishes:read", "dishes:write")
		if err != nil {
			app.serverErrorResponse(w, r, err)
			return
		}
	} else {
		err = app.models.Permissions.AddForUser(user.Id, "dishes:read", "restaurant:read")
		if err != nil {
			app.serverErrorResponse(w, r, err)
			return
		}
	}

	var claims jwt.Claims
	claims.Subject = strconv.FormatInt(user.Id, 10)
	claims.Issued = jwt.NewNumericTime(time.Now())
	claims.NotBefore = jwt.NewNumericTime(time.Now())
	claims.Expires = jwt.NewNumericTime(time.Now().Add(24 * time.Hour))
	claims.Issuer = "github.com/xtommas/food-backend"
	claims.Audiences = []string{"github.com/xtommas/food-backend"}

	claims.Set = map[string]interface{}{"scope": "activation"}

	jwtBytes, err := claims.HMACSign(jwt.HS256, []byte(app.config.jwt.secret))
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	err = app.writeJSON(w, http.StatusAccepted, envelope{"user": user, "activation_token": string(jwtBytes)}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) activateUserHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		TokenPlaintext string `json:"token"`
	}

	err := app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	claims, err := jwt.HMACCheck([]byte(input.TokenPlaintext), []byte(app.config.jwt.secret))
	if err != nil {
		app.invalidAuthenticationTokenResponse(w, r)
		return
	}

	if !claims.Valid(time.Now()) {
		app.invalidAuthenticationTokenResponse(w, r)
		return
	}

	if claims.Issuer != "github.com/xtommas/food-backend" {
		app.invalidAuthenticationTokenResponse(w, r)
		return
	}

	if !claims.AcceptAudience("github.com/xtommas/food-backend") {
		app.invalidAuthenticationTokenResponse(w, r)
		return
	}

	scope, _ := claims.Set["scope"].(string)
	if scope != "activation" {
		app.invalidAuthenticationTokenResponse(w, r)
		return
	}

	userId, err := strconv.ParseInt(claims.Subject, 10, 64)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	user, err := app.models.Users.Get(userId)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.invalidAuthenticationTokenResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	user.Activated = true

	err = app.models.Users.Update(user)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrEditConflict):
			app.editConflictResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"user": user}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) updateUserPasswordHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Password       string `json:"password"`
		TokenPlaintext string `json:"token"`
	}

	err := app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	v := validator.New()

	data.ValidatePasswordPlaintext(v, input.Password)
	//data.ValidateTokenPlaintext(v, input.TokenPlaintext)

	if !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	claims, err := jwt.HMACCheck([]byte(input.TokenPlaintext), []byte(app.config.jwt.secret))
	if err != nil {
		app.invalidAuthenticationTokenResponse(w, r)
		return
	}

	if !claims.Valid(time.Now()) {
		app.invalidAuthenticationTokenResponse(w, r)
		return
	}

	if claims.Issuer != "github.com/xtommas/food-backend" {
		app.invalidAuthenticationTokenResponse(w, r)
		return
	}

	if !claims.AcceptAudience("github.com/xtommas/food-backend") {
		app.invalidAuthenticationTokenResponse(w, r)
		return
	}

	scope, _ := claims.Set["scope"].(string)
	if scope != "password-reset" {
		app.invalidAuthenticationTokenResponse(w, r)
		return
	}

	userId, err := strconv.ParseInt(claims.Subject, 10, 64)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	user, err := app.models.Users.Get(userId)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.invalidAuthenticationTokenResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	err = user.Password.Set(input.Password)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	err = app.models.Users.Update(user)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrEditConflict):
			app.editConflictResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	env := envelope{"message": "your password was successfully reset"}

	err = app.writeJSON(w, http.StatusOK, env, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) updateUserRoleHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Email string `json:"email"`
		Role  string `json:"role"`
	}

	err := app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	v := validator.New()

	data.ValidateEmail(v, input.Email)
	data.ValidateRole(v, input.Role)

	if !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	user, err := app.models.Users.GetByEmail(input.Email)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			v.AddError("email", "invalid email")
			app.failedValidationResponse(w, r, v.Errors)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	oldRole := user.Role

	user.Role = input.Role

	err = app.models.Users.Update(user)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrEditConflict):
			app.editConflictResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	if user.Role != "customer" {
		err = app.models.Permissions.AddForUser(user.Id, "dishes:write")
		if err != nil {
			app.serverErrorResponse(w, r, err)
			return
		}
	} else {
		if oldRole != "customer" {
			err = app.models.Permissions.DeleteForUser(user.Id, "dishes:write")
			if err != nil {
				app.serverErrorResponse(w, r, err)
				return
			}
		}
	}
	env := envelope{"message": "role successfully updated"}

	err = app.writeJSON(w, http.StatusOK, env, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) getUserDataHandler(w http.ResponseWriter, r *http.Request) {
	user := app.contextGetUser(r)

	err := app.writeJSON(w, http.StatusOK, envelope{"user": user}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) listRestaurantsHandler(w http.ResponseWriter, r *http.Request) {
	restaurants, err := app.models.Users.GetRestaurants()
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"restaurants": restaurants}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) uploadUserPhotoHandler(w http.ResponseWriter, r *http.Request) {
	user := app.contextGetUser(r)

	err := r.ParseMultipartForm(10 << 20) //limit size to 10 MB
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	fileName := strconv.FormatInt(int64(user.Id), 10) + ".jpg"
	folder := "images/users/"

	user.Photo, err = app.storeImage(w, r, folder, fileName)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	err = app.models.Users.Update(user)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrEditConflict):
			app.editConflictResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"user": user}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) serveUserPhotoHandler(w http.ResponseWriter, r *http.Request) {
	user := app.contextGetUser(r)

	imagePath := user.Photo

	// Open the image file
	imageFile, err := os.Open(imagePath)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}
	defer imageFile.Close()

	// Get the image's content type
	contentType := mime.TypeByExtension(filepath.Ext(imagePath))

	// Set the content type header
	w.Header().Set("Content-Type", contentType)

	// Copy the image contents to the response
	_, err = io.Copy(w, imageFile)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) updateUserHandler(w http.ResponseWriter, r *http.Request) {
	user := app.contextGetUser(r)

	var input struct {
		Name  *string `json:"name"`
		Email *string `json:"email"`
	}

	err := app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	if input.Name != nil {
		user.Name = *input.Name
	}

	if input.Email != nil {
		user.Email = *input.Email
	}

	v := validator.New()

	if data.ValidateUser(v, user); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	err = app.models.Users.Update(user)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrEditConflict):
			app.editConflictResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"user": user}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}
