package data

import "testing"

func newTestRestaurant() *Restaurant {
	return &Restaurant{
		Name:      "Test Restaurant",
		Photo:     "test-restaurant.jpg",
		Address:   "123 Test Street",
		City:      "Test City",
		State:     "Test State",
		Province:  "Test Province",
		Country:   "Test Country",
		Latitude:  -34.603722,
		Longitude: -58.381592,
	}
}

func insertTestRestaurant(t *testing.T, model RestaurantModel) *Restaurant {
	t.Helper()

	restaurant := newTestRestaurant()
	if err := model.Insert(restaurant); err != nil {
		t.Fatalf("failed to insert test restaurant: %v", err)
	}

	t.Cleanup(func() {
		model.Delete(restaurant.ID)
	})

	return restaurant
}

func TestRestaurantModel_Insert(t *testing.T) {
	model := RestaurantModel{DB: testDB}
	restaurant := newTestRestaurant()

	err := model.Insert(restaurant)
	t.Cleanup(func() {
		model.Delete(restaurant.ID)
	})

	if err != nil {
		t.Fatalf("Insert() error = %v", err)
	}
	if restaurant.ID == 0 {
		t.Error("Insert() did not set restaurant.ID")
	}
	if restaurant.CreatedAt.IsZero() {
		t.Error("Insert() did not set CreatedAt")
	}
	if restaurant.Version != 1 {
		t.Errorf("Insert() Version = %d, want 1", restaurant.Version)
	}
}

func TestRestaurantModel_Get(t *testing.T) {
	model := RestaurantModel{DB: testDB}
	restaurant := insertTestRestaurant(t, model)

	fetched, err := model.Get(restaurant.ID)
	if err != nil {
		t.Fatalf("Get() error = %v", err)
	}

	if fetched.ID != restaurant.ID {
		t.Errorf("Get() ID = %d, want %d", fetched.ID, restaurant.ID)
	}
	if fetched.Name != restaurant.Name {
		t.Errorf("Get() Name = %q, want %q", fetched.Name, restaurant.Name)
	}
	if fetched.City != restaurant.City {
		t.Errorf("Get() City = %q, want %q", fetched.City, restaurant.City)
	}
}

func TestRestaurantModel_Get_NotFound(t *testing.T) {
	model := RestaurantModel{DB: testDB}

	_, err := model.Get(999999)
	if err != ErrRecordNotFound {
		t.Errorf("Get() error = %v, want ErrRecordNotFound", err)
	}

	_, err = model.Get(0)
	if err != ErrRecordNotFound {
		t.Errorf("Get() with id=0 error = %v, want ErrRecordNotFound", err)
	}
}

func TestRestaurantModel_Update(t *testing.T) {
	model := RestaurantModel{DB: testDB}
	restaurant := insertTestRestaurant(t, model)

	restaurant.Name = "Updated Restaurant"
	restaurant.City = "Updated City"
	restaurant.Latitude = -33.44889
	restaurant.Longitude = -70.66927
	oldVersion := restaurant.Version

	if err := model.Update(restaurant); err != nil {
		t.Fatalf("Update() error = %v", err)
	}

	if restaurant.Version != oldVersion+1 {
		t.Errorf("Update() Version = %d, want %d", restaurant.Version, oldVersion+1)
	}

	fetched, err := model.Get(restaurant.ID)
	if err != nil {
		t.Fatalf("Get() after Update() error = %v", err)
	}
	if fetched.Name != "Updated Restaurant" {
		t.Errorf("Update() Name = %q, want %q", fetched.Name, "Updated Restaurant")
	}
	if fetched.City != "Updated City" {
		t.Errorf("Update() City = %q, want %q", fetched.City, "Updated City")
	}
}

func TestRestaurantModel_Update_EditConflict(t *testing.T) {
	model := RestaurantModel{DB: testDB}
	restaurant := insertTestRestaurant(t, model)

	stale := *restaurant

	restaurant.Name = "Current Restaurant"
	if err := model.Update(restaurant); err != nil {
		t.Fatalf("Update() current restaurant error = %v", err)
	}

	stale.Name = "Stale Restaurant"
	err := model.Update(&stale)
	if err != ErrEditConflict {
		t.Errorf("Update() stale restaurant error = %v, want ErrEditConflict", err)
	}
}

func TestRestaurantModel_Delete(t *testing.T) {
	model := RestaurantModel{DB: testDB}
	restaurant := newTestRestaurant()
	if err := model.Insert(restaurant); err != nil {
		t.Fatalf("Insert() error = %v", err)
	}

	if err := model.Delete(restaurant.ID); err != nil {
		t.Fatalf("Delete() error = %v", err)
	}

	_, err := model.Get(restaurant.ID)
	if err != ErrRecordNotFound {
		t.Errorf("Get() after Delete() error = %v, want ErrRecordNotFound", err)
	}
}

func TestRestaurantModel_Delete_NotFound(t *testing.T) {
	model := RestaurantModel{DB: testDB}

	err := model.Delete(999999)
	if err != ErrRecordNotFound {
		t.Errorf("Delete() error = %v, want ErrRecordNotFound", err)
	}

	err = model.Delete(0)
	if err != ErrRecordNotFound {
		t.Errorf("Delete() with id=0 error = %v, want ErrRecordNotFound", err)
	}
}

func TestRestaurantModel_GetAll(t *testing.T) {
	model := RestaurantModel{DB: testDB}
	first := insertTestRestaurant(t, model)
	second := insertTestRestaurant(t, model)

	first.Name = "AAA Test Restaurant"
	second.Name = "ZZZ Test Restaurant"
	if err := model.Update(first); err != nil {
		t.Fatalf("Update() first restaurant error = %v", err)
	}
	if err := model.Update(second); err != nil {
		t.Fatalf("Update() second restaurant error = %v", err)
	}

	restaurants, err := model.GetAll()
	if err != nil {
		t.Fatalf("GetAll() error = %v", err)
	}

	firstIndex := -1
	secondIndex := -1
	for i, restaurant := range restaurants {
		if restaurant.ID == first.ID {
			firstIndex = i
		}
		if restaurant.ID == second.ID {
			secondIndex = i
		}
	}

	if firstIndex == -1 {
		t.Fatal("GetAll() did not include first test restaurant")
	}
	if secondIndex == -1 {
		t.Fatal("GetAll() did not include second test restaurant")
	}
	if firstIndex > secondIndex {
		t.Errorf("GetAll() order placed %q after %q", first.Name, second.Name)
	}
}

func TestRestaurantModel_Staff(t *testing.T) {
	restaurantModel := RestaurantModel{DB: testDB}
	userModel := UserModel{DB: testDB}
	restaurant := insertTestRestaurant(t, restaurantModel)
	user := insertTestUser(t, userModel)

	isStaff, err := restaurantModel.IsStaff(restaurant.ID, user.Id)
	if err != nil {
		t.Fatalf("IsStaff() before AddStaff() error = %v", err)
	}
	if isStaff {
		t.Error("IsStaff() should be false before AddStaff()")
	}

	if err := restaurantModel.AddStaff(restaurant.ID, user.Id, "manager"); err != nil {
		t.Fatalf("AddStaff() error = %v", err)
	}
	t.Cleanup(func() {
		restaurantModel.RemoveStaff(restaurant.ID, user.Id)
	})

	isStaff, err = restaurantModel.IsStaff(restaurant.ID, user.Id)
	if err != nil {
		t.Fatalf("IsStaff() after AddStaff() error = %v", err)
	}
	if !isStaff {
		t.Error("IsStaff() should be true after AddStaff()")
	}

	role, err := restaurantModel.GetStaffRole(restaurant.ID, user.Id)
	if err != nil {
		t.Fatalf("GetStaffRole() error = %v", err)
	}
	if role != "manager" {
		t.Errorf("GetStaffRole() = %q, want manager", role)
	}

	staff, err := restaurantModel.GetStaff(restaurant.ID)
	if err != nil {
		t.Fatalf("GetStaff() error = %v", err)
	}
	if len(staff) != 1 {
		t.Fatalf("GetStaff() returned %d users, want 1", len(staff))
	}
	if staff[0].Id != user.Id {
		t.Errorf("GetStaff() user id = %d, want %d", staff[0].Id, user.Id)
	}

	if err := restaurantModel.RemoveStaff(restaurant.ID, user.Id); err != nil {
		t.Fatalf("RemoveStaff() error = %v", err)
	}

	role, err = restaurantModel.GetStaffRole(restaurant.ID, user.Id)
	if err != ErrRecordNotFound {
		t.Errorf("GetStaffRole() after RemoveStaff() error = %v, want ErrRecordNotFound", err)
	}
	if role != "" {
		t.Errorf("GetStaffRole() after RemoveStaff() role = %q, want empty string", role)
	}
}

func TestRestaurantModel_RemoveStaff_NotFound(t *testing.T) {
	model := RestaurantModel{DB: testDB}

	err := model.RemoveStaff(999999, 999999)
	if err != ErrRecordNotFound {
		t.Errorf("RemoveStaff() error = %v, want ErrRecordNotFound", err)
	}
}
