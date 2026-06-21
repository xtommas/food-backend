package main

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"os"
	"time"

	_ "github.com/lib/pq"

	"github.com/xtommas/food-backend/internal/data"
)

func main() {
	dsn := os.Getenv("DB_DSN")
	if dsn == "" {
		log.Fatal("DB_DSN is not set")
	}

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		log.Fatalf("opening db: %v", err)
	}
	defer db.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := db.PingContext(ctx); err != nil {
		log.Fatalf("pinging db: %v", err)
	}

	// If any restaurant already exists, assume
	// this environment has already been seeded and bail out
	var restaurantCount int
	err = db.QueryRow(`SELECT COUNT(*) FROM restaurants`).Scan(&restaurantCount)
	if err != nil {
		log.Fatalf("checking existing data: %v", err)
	}
	if restaurantCount > 0 {
		log.Println("database already has data, skipping seed")
		return
	}

	models := data.NewModels(db)

	log.Println("seeding users...")

	admin, err := upsertUser(models, "Admin", "admin@example.com", "admin", "password123")
	if err != nil {
		log.Fatalf("seeding admin: %v", err)
	}

	owner, err := upsertUser(models, "Mario Rossi", "owner@example.com", "customer", "password123")
	if err != nil {
		log.Fatalf("seeding owner: %v", err)
	}

	staff, err := upsertUser(models, "Luigi Bianchi", "staff@example.com", "customer", "password123")
	if err != nil {
		log.Fatalf("seeding staff: %v", err)
	}

	customer, err := upsertUser(models, "Jane Customer", "customer@example.com", "customer", "password123")
	if err != nil {
		log.Fatalf("seeding customer: %v", err)
	}

	log.Println("seeding restaurant...")

	restaurant := &data.Restaurant{
		Name:      "Trattoria Roma",
		Address:   "Via del Corso 1",
		City:      "Rome",
		Country:   "Italy",
		Latitude:  41.902782,
		Longitude: 12.496366,
	}
	if err := models.Restaurants.Insert(restaurant); err != nil {
		log.Fatalf("inserting restaurant: %v", err)
	}

	log.Println("assigning restaurant staff...")

	if err := models.Restaurants.AddStaff(restaurant.ID, owner.Id, "owner"); err != nil {
		log.Fatalf("assigning owner: %v", err)
	}
	if err := models.Restaurants.AddStaff(restaurant.ID, staff.Id, "staff"); err != nil {
		log.Fatalf("assigning staff: %v", err)
	}

	log.Println("seeding dishes...")

	dishes := []*data.Dish{
		{
			RestaurantID: restaurant.ID,
			Name:         "Margherita Pizza",
			Price:        1200,
			Description:  "Classic pizza with tomato, mozzarella, and basil.",
			Categories:   []string{"pizza", "vegetarian"},
		},
		{
			RestaurantID: restaurant.ID,
			Name:         "Spaghetti Carbonara",
			Price:        1450,
			Description:  "Spaghetti with egg, pecorino, guanciale, and black pepper.",
			Categories:   []string{"pasta"},
		},
		{
			RestaurantID: restaurant.ID,
			Name:         "Tiramisu",
			Price:        700,
			Description:  "Coffee-soaked ladyfingers layered with mascarpone cream.",
			Categories:   []string{"dessert"},
		},
	}
	for _, dish := range dishes {
		if err := models.Dishes.Insert(dish); err != nil {
			log.Fatalf("inserting dish %q: %v", dish.Name, err)
		}
	}

	log.Println("granting default permissions...")

	codes := []string{"dishes:read", "restaurants:read"}
	userIDs := []int64{admin.Id, owner.Id, staff.Id, customer.Id}

	for _, code := range codes {
		permissionID, err := ensurePermission(db, code)
		if err != nil {
			log.Fatalf("ensuring permission %q: %v", code, err)
		}

		for _, userID := range userIDs {
			_, err := db.Exec(
				`INSERT INTO users_permissions (user_id, permission_id) VALUES ($1, $2)
				 ON CONFLICT DO NOTHING`,
				userID, permissionID,
			)
			if err != nil {
				log.Fatalf("granting %q to user %d: %v", code, userID, err)
			}
		}
	}

	fmt.Println("seed complete:")
	fmt.Println("  admin:    admin@example.com    / password123")
	fmt.Printf("  owner:    owner@example.com    / password123  (owns %q)\n", restaurant.Name)
	fmt.Println("  staff:    staff@example.com    / password123")
	fmt.Println("  customer: customer@example.com / password123")
}

// upsertUser inserts a user, or fetches the existing one if the email
// is already taken (so a partially-failed seed run can be re-run safely).
func upsertUser(models data.Models, name, email, role, plaintextPassword string) (*data.User, error) {
	user := &data.User{
		Name:      name,
		Email:     email,
		Activated: true,
		Role:      role,
	}
	if err := user.Password.Set(plaintextPassword); err != nil {
		return nil, fmt.Errorf("hashing password for %s: %w", email, err)
	}

	err := models.Users.Insert(user)
	if err == nil {
		return user, nil
	}
	if errors.Is(err, data.ErrDuplicateEmail) {
		return models.Users.GetByEmail(email)
	}
	return nil, err
}

func ensurePermission(db *sql.DB, code string) (int64, error) {
	var id int64
	err := db.QueryRow(`SELECT id FROM permissions WHERE code = $1`, code).Scan(&id)
	if err == nil {
		return id, nil
	}
	if !errors.Is(err, sql.ErrNoRows) {
		return 0, err
	}
	err = db.QueryRow(`INSERT INTO permissions (code) VALUES ($1) RETURNING id`, code).Scan(&id)
	return id, err
}
