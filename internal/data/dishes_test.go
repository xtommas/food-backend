package data

import (
	"database/sql"
	"testing"
	"time"

	_ "github.com/lib/pq"
)

// ---- seeds ----

func seedRestaurant(t *testing.T) int64 {
	t.Helper()
	var id int64
	err := testDB.QueryRow(`
		INSERT INTO restaurants (name, address, city, country)
		VALUES ('Test Restaurant', 'Test Address', 'Test City', 'Test Country')
		RETURNING id`).Scan(&id)
	if err != nil {
		t.Fatalf("seedRestaurant() error = %v", err)
	}
	t.Cleanup(func() {
		testDB.Exec(`DELETE FROM restaurants WHERE id = $1`, id)
	})
	return id
}

// ---- helpers ----

func newTestDish(restaurantID int64) *Dish {
	return &Dish{
		RestaurantID: restaurantID,
		Name:         "Test Dish",
		Price:        500,
		Description:  "A test dish description",
		Categories:   []string{"test", "food"},
	}
}

func insertTestDish(t *testing.T, model DishModel, restaurantID int64) *Dish {
	t.Helper()
	dish := newTestDish(restaurantID)
	if err := model.Insert(dish); err != nil {
		t.Fatalf("failed to insert test dish: %v", err)
	}
	t.Cleanup(func() { model.Delete(dish.ID) })
	return dish
}

// ---- Insert ----

func TestDishModel_Insert(t *testing.T) {
	model := DishModel{DB: testDB}
	restaurantID := seedRestaurant(t)

	dish := newTestDish(restaurantID)
	err := model.Insert(dish)
	t.Cleanup(func() { model.Delete(dish.ID) })

	if err != nil {
		t.Fatalf("Insert() error = %v", err)
	}
	if dish.ID == 0 {
		t.Error("Insert() did not set dish.ID")
	}
	if !dish.Available {
		t.Error("Insert() dish should default to available=true")
	}
	if dish.UpdatedAt.IsZero() {
		t.Error("Insert() did not set UpdatedAt")
	}
}

func TestDishModel_Insert_ValidationEdgeCases(t *testing.T) {
	model := DishModel{DB: testDB}
	restaurantID := seedRestaurant(t)

	tests := []struct {
		name    string
		mutate  func(*Dish)
		wantErr bool
	}{
		{"zero price is not allowed by DB", func(d *Dish) { d.Price = 0 }, true},
		{"empty name is not allowed by DB", func(d *Dish) { d.Name = "" }, true},
		{"duplicate categories are allowed by DB", func(d *Dish) { d.Categories = []string{"a", "a"} }, false},
		{"negative price rejected by DB", func(d *Dish) { d.Price = -1 }, true},
		{"too many categories rejected by DB", func(d *Dish) { d.Categories = []string{"a", "b", "c", "d", "e", "f"} }, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dish := newTestDish(restaurantID)
			tt.mutate(dish)
			err := model.Insert(dish)
			if err == nil {
				t.Cleanup(func() { model.Delete(dish.ID) })
			}
			if (err != nil) != tt.wantErr {
				t.Errorf("Insert() error = %v, wantErr = %v", err, tt.wantErr)
			}
		})
	}
}

// ---- Get ----

func TestDishModel_Get(t *testing.T) {
	model := DishModel{DB: testDB}
	restaurantID := seedRestaurant(t)
	dish := insertTestDish(t, model, restaurantID)

	fetched, err := model.Get(dish.ID)
	if err != nil {
		t.Fatalf("Get() error = %v", err)
	}

	if fetched.ID != dish.ID {
		t.Errorf("Get() ID = %d, want %d", fetched.ID, dish.ID)
	}
	if fetched.Name != dish.Name {
		t.Errorf("Get() Name = %q, want %q", fetched.Name, dish.Name)
	}
	if fetched.Price != dish.Price {
		t.Errorf("Get() Price = %d, want %d", fetched.Price, dish.Price)
	}
	if fetched.RestaurantID != dish.RestaurantID {
		t.Errorf("Get() RestaurantID = %d, want %d", fetched.RestaurantID, dish.RestaurantID)
	}
}

func TestDishModel_Get_NotFound(t *testing.T) {
	model := DishModel{DB: testDB}

	_, err := model.Get(999999)
	if err != ErrRecordNotFound {
		t.Errorf("Get() error = %v, want ErrRecordNotFound", err)
	}
}

func TestDishModel_Get_InvalidID(t *testing.T) {
	model := DishModel{DB: testDB}

	_, err := model.Get(0)
	if err != ErrRecordNotFound {
		t.Errorf("Get() with id=0 error = %v, want ErrRecordNotFound", err)
	}

	_, err = model.Get(-1)
	if err != ErrRecordNotFound {
		t.Errorf("Get() with id=-1 error = %v, want ErrRecordNotFound", err)
	}
}

// ---- Update ----

func TestDishModel_Update(t *testing.T) {
	model := DishModel{DB: testDB}
	restaurantID := seedRestaurant(t)
	dish := insertTestDish(t, model, restaurantID)

	dish.Name = "Updated Name"
	dish.Price = 999
	dish.Available = false

	if err := model.Update(dish); err != nil {
		t.Fatalf("Update() error = %v", err)
	}

	fetched, err := model.Get(dish.ID)
	if err != nil {
		t.Fatalf("Get() after Update() error = %v", err)
	}

	if fetched.Name != "Updated Name" {
		t.Errorf("Update() Name = %q, want %q", fetched.Name, "Updated Name")
	}
	if fetched.Price != 999 {
		t.Errorf("Update() Price = %d, want 999", fetched.Price)
	}
	if fetched.Available != false {
		t.Error("Update() Available should be false")
	}
}

func TestDishModel_Update_SetsUpdatedAt(t *testing.T) {
	model := DishModel{DB: testDB}
	restaurantID := seedRestaurant(t)
	dish := insertTestDish(t, model, restaurantID)
	originalUpdatedAt := dish.UpdatedAt

	time.Sleep(1100 * time.Millisecond)

	dish.Name = "New Name"
	if err := model.Update(dish); err != nil {
		t.Fatalf("Update() error = %v", err)
	}

	if !dish.UpdatedAt.After(originalUpdatedAt) {
		t.Error("Update() should advance UpdatedAt")
	}
}

func TestDishModel_Update_NotFound(t *testing.T) {
	model := DishModel{DB: testDB}
	restaurantID := seedRestaurant(t)

	ghost := &Dish{
		ID:           999999,
		RestaurantID: restaurantID,
		Name:         "Ghost",
		Price:        100,
		Description:  "Does not exist",
		Categories:   []string{"none"},
	}

	err := model.Update(ghost)
	if err != ErrRecordNotFound {
		t.Errorf("Update() error = %v, want ErrRecordNotFound", err)
	}
}

// ---- Delete ----

func TestDishModel_Delete(t *testing.T) {
	model := DishModel{DB: testDB}
	restaurantID := seedRestaurant(t)

	dish := newTestDish(restaurantID)
	if err := model.Insert(dish); err != nil {
		t.Fatalf("Insert() error = %v", err)
	}

	if err := model.Delete(dish.ID); err != nil {
		t.Fatalf("Delete() error = %v", err)
	}

	_, err := model.Get(dish.ID)
	if err != ErrRecordNotFound {
		t.Errorf("Get() after Delete() error = %v, want ErrRecordNotFound", err)
	}
}

func TestDishModel_Delete_NotFound(t *testing.T) {
	model := DishModel{DB: testDB}

	err := model.Delete(999999)
	if err != ErrRecordNotFound {
		t.Errorf("Delete() error = %v, want ErrRecordNotFound", err)
	}
}

func TestDishModel_Delete_InvalidID(t *testing.T) {
	model := DishModel{DB: testDB}

	err := model.Delete(0)
	if err != ErrRecordNotFound {
		t.Errorf("Delete() with id=0 error = %v, want ErrRecordNotFound", err)
	}
}

// ---- GetAllForRestaurant ----

func TestDishModel_GetAllForRestaurant(t *testing.T) {
	model := DishModel{DB: testDB}
	restaurantID := seedRestaurant(t)

	insertTestDish(t, model, restaurantID)
	insertTestDish(t, model, restaurantID)

	// decoy restaurant — dishes should NOT appear in results
	otherRestaurantID := seedRestaurant(t)
	insertTestDish(t, model, otherRestaurantID)

	filters := Filters{Page: 1, PageSize: 20, Sort: "id", SortSafelist: []string{"id"}}
	dishes, metadata, err := model.GetAllForRestaurant(restaurantID, "", []string{}, sql.NullBool{}, filters)
	if err != nil {
		t.Fatalf("GetAllForRestaurant() error = %v", err)
	}

	for _, d := range dishes {
		if d.RestaurantID != restaurantID {
			t.Errorf("GetAllForRestaurant() returned dish with RestaurantID = %d, want %d", d.RestaurantID, restaurantID)
		}
	}
	if metadata.TotalRecords < 2 {
		t.Errorf("GetAllForRestaurant() TotalRecords = %d, want >= 2", metadata.TotalRecords)
	}
}

func TestDishModel_GetAllForRestaurant_FilterByName(t *testing.T) {
	model := DishModel{DB: testDB}
	restaurantID := seedRestaurant(t)
	dish := insertTestDish(t, model, restaurantID)

	dish.Name = "Unique Searchable Name"
	model.Update(dish)

	filters := Filters{Page: 1, PageSize: 20, Sort: "id", SortSafelist: []string{"id"}}
	results, _, err := model.GetAllForRestaurant(restaurantID, "Unique Searchable", []string{}, sql.NullBool{}, filters)
	if err != nil {
		t.Fatalf("GetAllForRestaurant() error = %v", err)
	}

	found := false
	for _, d := range results {
		if d.ID == dish.ID {
			found = true
		}
	}
	if !found {
		t.Error("GetAllForRestaurant() name filter did not return the expected dish")
	}
}

func TestDishModel_GetAllForRestaurant_FilterByAvailability(t *testing.T) {
	model := DishModel{DB: testDB}
	restaurantID := seedRestaurant(t)
	dish := insertTestDish(t, model, restaurantID)

	dish.Available = false
	model.Update(dish)

	filters := Filters{Page: 1, PageSize: 20, Sort: "id", SortSafelist: []string{"id"}}

	unavailable, _, err := model.GetAllForRestaurant(restaurantID, "", []string{}, sql.NullBool{Valid: true, Bool: false}, filters)
	if err != nil {
		t.Fatalf("GetAllForRestaurant() error = %v", err)
	}
	for _, d := range unavailable {
		if d.Available {
			t.Errorf("GetAllForRestaurant() available=false filter returned an available dish (id=%d)", d.ID)
		}
	}

	available, _, err := model.GetAllForRestaurant(restaurantID, "", []string{}, sql.NullBool{Valid: true, Bool: true}, filters)
	if err != nil {
		t.Fatalf("GetAllForRestaurant() error = %v", err)
	}
	for _, d := range available {
		if d.ID == dish.ID {
			t.Error("GetAllForRestaurant() available=true filter returned an unavailable dish")
		}
	}
}

func TestDishModel_GetAllForRestaurant_FilterByCategory(t *testing.T) {
	model := DishModel{DB: testDB}
	restaurantID := seedRestaurant(t)
	dish := insertTestDish(t, model, restaurantID)

	dish.Categories = []string{"vegan", "salad"}
	model.Update(dish)

	filters := Filters{Page: 1, PageSize: 20, Sort: "id", SortSafelist: []string{"id"}}
	results, _, err := model.GetAllForRestaurant(restaurantID, "", []string{"vegan"}, sql.NullBool{}, filters)
	if err != nil {
		t.Fatalf("GetAllForRestaurant() error = %v", err)
	}

	found := false
	for _, d := range results {
		if d.ID == dish.ID {
			found = true
		}
	}
	if !found {
		t.Error("GetAllForRestaurant() category filter did not return the expected dish")
	}
}

func TestDishModel_GetAllForRestaurant_Pagination(t *testing.T) {
	model := DishModel{DB: testDB}
	restaurantID := seedRestaurant(t)

	for i := 0; i < 3; i++ {
		insertTestDish(t, model, restaurantID)
	}

	page1Filters := Filters{Page: 1, PageSize: 2, Sort: "id", SortSafelist: []string{"id"}}
	page2Filters := Filters{Page: 2, PageSize: 2, Sort: "id", SortSafelist: []string{"id"}}

	page1, meta1, err := model.GetAllForRestaurant(restaurantID, "", []string{}, sql.NullBool{}, page1Filters)
	if err != nil {
		t.Fatalf("GetAllForRestaurant() page 1 error = %v", err)
	}
	page2, _, err := model.GetAllForRestaurant(restaurantID, "", []string{}, sql.NullBool{}, page2Filters)
	if err != nil {
		t.Fatalf("GetAllForRestaurant() page 2 error = %v", err)
	}

	if len(page1) > 2 {
		t.Errorf("page 1 returned %d results, want <= 2", len(page1))
	}
	if meta1.CurrentPage != 1 {
		t.Errorf("meta CurrentPage = %d, want 1", meta1.CurrentPage)
	}

	p1IDs := map[int64]bool{}
	for _, d := range page1 {
		p1IDs[d.ID] = true
	}
	for _, d := range page2 {
		if p1IDs[d.ID] {
			t.Errorf("dish id %d appeared on both pages", d.ID)
		}
	}
}
