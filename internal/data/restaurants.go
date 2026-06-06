package data

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/xtommas/food-backend/internal/validator"
)

type Restaurant struct {
	ID        int64     `json:"id"`
	Name      string    `json:"name"`
	Photo     string    `json:"photo,omitempty"`
	Address   string    `json:"address"`
	City      string    `json:"city"`
	State     string    `json:"state,omitempty"`
	Province  string    `json:"province,omitempty"`
	Country   string    `json:"country"`
	Latitude  float64   `json:"latitude,omitempty"`
	Longitude float64   `json:"longitude,omitempty"`
	CreatedAt time.Time `json:"created_at"`
	Version   int       `json:"-"`
}

func ValidateRestaurant(v *validator.Validator, r *Restaurant) {
	v.Check(r.Name != "", "name", "must be provided")
	v.Check(len(r.Name) <= 500, "name", "must be no more than 500 characters long")

	v.Check(r.Address != "", "address", "must be provided")
	v.Check(r.City != "", "city", "must be provided")
	v.Check(r.Country != "", "country", "must be provided")

	if r.Latitude != 0 || r.Longitude != 0 {
		v.Check(r.Latitude >= -90 && r.Latitude <= 90, "latitude", "must be between -90 and 90")
		v.Check(r.Longitude >= -180 && r.Longitude <= 180, "longitude", "must be between -180 and 180")
	}
}

type RestaurantModel struct {
	DB *sql.DB
}

func (m RestaurantModel) Insert(restaurant *Restaurant) error {
	query := `
		INSERT INTO restaurants (name, photo, address, city, state, province, country, latitude, longitude)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
		RETURNING id, created_at, version`

	args := []any{
		restaurant.Name,
		restaurant.Photo,
		restaurant.Address,
		restaurant.City,
		restaurant.State,
		restaurant.Province,
		restaurant.Country,
		restaurant.Latitude,
		restaurant.Longitude,
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	return m.DB.QueryRowContext(ctx, query, args...).Scan(
		&restaurant.ID,
		&restaurant.CreatedAt,
		&restaurant.Version,
	)
}

func (m RestaurantModel) Get(id int64) (*Restaurant, error) {
	if id < 1 {
		return nil, ErrRecordNotFound
	}

	query := `
		SELECT id, name, photo, address, city, state, province, country, latitude, longitude, created_at, version
		FROM restaurants
		WHERE id = $1`

	var r Restaurant

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := m.DB.QueryRowContext(ctx, query, id).Scan(
		&r.ID,
		&r.Name,
		&r.Photo,
		&r.Address,
		&r.City,
		&r.State,
		&r.Province,
		&r.Country,
		&r.Latitude,
		&r.Longitude,
		&r.CreatedAt,
		&r.Version,
	)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrRecordNotFound
		default:
			return nil, err
		}
	}

	return &r, nil
}

func (m RestaurantModel) Update(restaurant *Restaurant) error {
	query := `
		UPDATE restaurants
		SET name = $1, photo = $2, address = $3, city = $4, state = $5, province = $6,
		    country = $7, latitude = $8, longitude = $9, version = version + 1
		WHERE id = $10 AND version = $11
		RETURNING version`

	args := []any{
		restaurant.Name,
		restaurant.Photo,
		restaurant.Address,
		restaurant.City,
		restaurant.State,
		restaurant.Province,
		restaurant.Country,
		restaurant.Latitude,
		restaurant.Longitude,
		restaurant.ID,
		restaurant.Version,
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := m.DB.QueryRowContext(ctx, query, args...).Scan(&restaurant.Version)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return ErrEditConflict
		default:
			return err
		}
	}

	return nil
}

func (m RestaurantModel) Delete(id int64) error {
	if id < 1 {
		return ErrRecordNotFound
	}

	query := `DELETE FROM restaurants WHERE id = $1`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	result, err := m.DB.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return ErrRecordNotFound
	}

	return nil
}

func (m RestaurantModel) GetAll() ([]*Restaurant, error) {
	query := `
		SELECT id, name, photo, address, city, state, province, country, latitude, longitude, created_at, version
		FROM restaurants
		ORDER BY name ASC`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	rows, err := m.DB.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var restaurants []*Restaurant

	for rows.Next() {
		var r Restaurant

		err := rows.Scan(
			&r.ID,
			&r.Name,
			&r.Photo,
			&r.Address,
			&r.City,
			&r.State,
			&r.Province,
			&r.Country,
			&r.Latitude,
			&r.Longitude,
			&r.CreatedAt,
			&r.Version,
		)
		if err != nil {
			return nil, err
		}

		restaurants = append(restaurants, &r)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return restaurants, nil
}

func (m RestaurantModel) GetStaff(restaurantID int64) ([]*User, error) {
	query := `
		SELECT u.id, u.photo, u.created_at, u.name, u.email, u.activated, u.version, u.role
		FROM users u
		INNER JOIN restaurant_staff rs ON rs.user_id = u.id
		WHERE rs.restaurant_id = $1`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	rows, err := m.DB.QueryContext(ctx, query, restaurantID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []*User

	for rows.Next() {
		var u User

		err := rows.Scan(
			&u.Id,
			&u.Photo,
			&u.CreatedAt,
			&u.Name,
			&u.Email,
			&u.Activated,
			&u.Version,
			&u.Role,
		)
		if err != nil {
			return nil, err
		}

		users = append(users, &u)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return users, nil
}

func (m RestaurantModel) AddStaff(restaurantID, userID int64, role string) error {
	query := `
		INSERT INTO restaurant_staff (restaurant_id, user_id, role)
		VALUES ($1, $2, $3)
		ON CONFLICT (restaurant_id, user_id) DO UPDATE SET role = EXCLUDED.role`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	_, err := m.DB.ExecContext(ctx, query, restaurantID, userID, role)
	return err
}

func (m RestaurantModel) RemoveStaff(restaurantID, userID int64) error {
	query := `
		DELETE FROM restaurant_staff
		WHERE restaurant_id = $1 AND user_id = $2`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	result, err := m.DB.ExecContext(ctx, query, restaurantID, userID)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return ErrRecordNotFound
	}

	return nil
}

func (m RestaurantModel) IsStaff(restaurantID, userID int64) (bool, error) {
	query := `
		SELECT EXISTS (
			SELECT 1 FROM restaurant_staff
			WHERE restaurant_id = $1 AND user_id = $2
		)`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var exists bool
	err := m.DB.QueryRowContext(ctx, query, restaurantID, userID).Scan(&exists)
	return exists, err
}

func (m RestaurantModel) GetStaffRole(restaurantID, userID int64) (string, error) {
	query := `
		SELECT role
		FROM restaurant_staff
		WHERE restaurant_id = $1 AND user_id = $2`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var role string
	err := m.DB.QueryRowContext(ctx, query, restaurantID, userID).Scan(&role)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return "", ErrRecordNotFound
		default:
			return "", err
		}
	}

	return role, nil
}
