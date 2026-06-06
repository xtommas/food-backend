package data

import "testing"

func newTestOrder(userID, restaurantID int64) *Order {
	return &Order{
		UserID:       userID,
		RestaurantID: restaurantID,
		Total:        1500,
		Address:      "123 Test Street",
		Status:       "pending",
	}
}

func insertTestOrder(t *testing.T, model OrderModel, userID, restaurantID int64) *Order {
	t.Helper()

	order := newTestOrder(userID, restaurantID)
	if err := model.Insert(order); err != nil {
		t.Fatalf("failed to insert test order: %v", err)
	}

	t.Cleanup(func() {
		testDB.Exec(`DELETE FROM orders WHERE id = $1`, order.ID)
	})

	return order
}

func newTestFilters() Filters {
	return Filters{
		Page:         1,
		PageSize:     10,
		Sort:         "id",
		SortSafelist: []string{"id", "-id", "created_at", "-created_at", "updated_at", "-updated_at", "total", "-total", "status", "-status"},
	}
}

func TestOrderModel_Insert(t *testing.T) {
	userModel := UserModel{DB: testDB}
	orderModel := OrderModel{DB: testDB}
	restaurantID := seedRestaurant(t)
	user := insertTestUser(t, userModel)
	order := newTestOrder(user.Id, restaurantID)

	err := orderModel.Insert(order)
	t.Cleanup(func() {
		testDB.Exec(`DELETE FROM orders WHERE id = $1`, order.ID)
	})

	if err != nil {
		t.Fatalf("Insert() error = %v", err)
	}
	if order.ID == 0 {
		t.Error("Insert() did not set order.ID")
	}
	if order.CreatedAt.IsZero() {
		t.Error("Insert() did not set CreatedAt")
	}
	if order.UpdatedAt.IsZero() {
		t.Error("Insert() did not set UpdatedAt")
	}
}

func TestOrderModel_GetForRestaurant(t *testing.T) {
	userModel := UserModel{DB: testDB}
	orderModel := OrderModel{DB: testDB}
	restaurantID := seedRestaurant(t)
	user := insertTestUser(t, userModel)
	order := insertTestOrder(t, orderModel, user.Id, restaurantID)

	fetched, err := orderModel.GetForRestaurant(order.ID, restaurantID)
	if err != nil {
		t.Fatalf("GetForRestaurant() error = %v", err)
	}

	if fetched.ID != order.ID {
		t.Errorf("GetForRestaurant() ID = %d, want %d", fetched.ID, order.ID)
	}
	if fetched.RestaurantID != restaurantID {
		t.Errorf("GetForRestaurant() RestaurantID = %d, want %d", fetched.RestaurantID, restaurantID)
	}
}

func TestOrderModel_GetForRestaurant_NotFound(t *testing.T) {
	model := OrderModel{DB: testDB}

	_, err := model.GetForRestaurant(999999, 999999)
	if err != ErrRecordNotFound {
		t.Errorf("GetForRestaurant() error = %v, want ErrRecordNotFound", err)
	}

	_, err = model.GetForRestaurant(0, 999999)
	if err != ErrRecordNotFound {
		t.Errorf("GetForRestaurant() with id=0 error = %v, want ErrRecordNotFound", err)
	}
}

func TestOrderModel_GetForUser(t *testing.T) {
	userModel := UserModel{DB: testDB}
	orderModel := OrderModel{DB: testDB}
	restaurantID := seedRestaurant(t)
	user := insertTestUser(t, userModel)
	order := insertTestOrder(t, orderModel, user.Id, restaurantID)

	fetched, err := orderModel.GetForUser(order.ID, user.Id)
	if err != nil {
		t.Fatalf("GetForUser() error = %v", err)
	}

	if fetched.ID != order.ID {
		t.Errorf("GetForUser() ID = %d, want %d", fetched.ID, order.ID)
	}
	if fetched.UserID != user.Id {
		t.Errorf("GetForUser() UserID = %d, want %d", fetched.UserID, user.Id)
	}
}

func TestOrderModel_GetForUser_NotFound(t *testing.T) {
	model := OrderModel{DB: testDB}

	_, err := model.GetForUser(999999, 999999)
	if err != ErrRecordNotFound {
		t.Errorf("GetForUser() error = %v, want ErrRecordNotFound", err)
	}

	_, err = model.GetForUser(0, 999999)
	if err != ErrRecordNotFound {
		t.Errorf("GetForUser() with id=0 error = %v, want ErrRecordNotFound", err)
	}
}

func TestOrderModel_Update(t *testing.T) {
	userModel := UserModel{DB: testDB}
	orderModel := OrderModel{DB: testDB}
	restaurantID := seedRestaurant(t)
	user := insertTestUser(t, userModel)
	order := insertTestOrder(t, orderModel, user.Id, restaurantID)

	order.Total = 2200
	order.Status = "confirmed"
	if err := orderModel.Update(order); err != nil {
		t.Fatalf("Update() error = %v", err)
	}

	fetched, err := orderModel.GetForUser(order.ID, user.Id)
	if err != nil {
		t.Fatalf("GetForUser() after Update() error = %v", err)
	}
	if fetched.Total != 2200 {
		t.Errorf("Update() Total = %d, want 2200", fetched.Total)
	}
	if fetched.Status != "confirmed" {
		t.Errorf("Update() Status = %q, want confirmed", fetched.Status)
	}
}

func TestOrderModel_Update_NotFound(t *testing.T) {
	model := OrderModel{DB: testDB}
	order := &Order{
		ID:      999999,
		Total:   100,
		Status:  "confirmed",
		Address: "123 Test Street",
	}

	err := model.Update(order)
	if err != ErrRecordNotFound {
		t.Errorf("Update() error = %v, want ErrRecordNotFound", err)
	}
}

func TestOrderModel_GetAllForRestaurant(t *testing.T) {
	userModel := UserModel{DB: testDB}
	orderModel := OrderModel{DB: testDB}
	restaurantID := seedRestaurant(t)
	user := insertTestUser(t, userModel)
	pending := insertTestOrder(t, orderModel, user.Id, restaurantID)
	confirmed := insertTestOrder(t, orderModel, user.Id, restaurantID)
	confirmed.Status = "confirmed"
	if err := orderModel.Update(confirmed); err != nil {
		t.Fatalf("Update() confirmed order error = %v", err)
	}

	orders, metadata, err := orderModel.GetAllForRestaurant(restaurantID, "pending", newTestFilters())
	if err != nil {
		t.Fatalf("GetAllForRestaurant() error = %v", err)
	}

	if len(orders) != 1 {
		t.Fatalf("GetAllForRestaurant() returned %d orders, want 1", len(orders))
	}
	if orders[0].ID != pending.ID {
		t.Errorf("GetAllForRestaurant() order ID = %d, want %d", orders[0].ID, pending.ID)
	}
	if metadata.TotalRecords != 1 {
		t.Errorf("GetAllForRestaurant() TotalRecords = %d, want 1", metadata.TotalRecords)
	}
}

func TestOrderModel_GetAllForUser(t *testing.T) {
	userModel := UserModel{DB: testDB}
	orderModel := OrderModel{DB: testDB}
	restaurantID := seedRestaurant(t)
	user := insertTestUser(t, userModel)
	first := insertTestOrder(t, orderModel, user.Id, restaurantID)
	second := insertTestOrder(t, orderModel, user.Id, restaurantID)

	orders, metadata, err := orderModel.GetAllForUser(user.Id, "", newTestFilters())
	if err != nil {
		t.Fatalf("GetAllForUser() error = %v", err)
	}

	if len(orders) != 2 {
		t.Fatalf("GetAllForUser() returned %d orders, want 2", len(orders))
	}
	if orders[0].ID != first.ID || orders[1].ID != second.ID {
		t.Errorf("GetAllForUser() IDs = [%d %d], want [%d %d]", orders[0].ID, orders[1].ID, first.ID, second.ID)
	}
	if metadata.TotalRecords != 2 {
		t.Errorf("GetAllForUser() TotalRecords = %d, want 2", metadata.TotalRecords)
	}
}
