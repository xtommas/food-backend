package data

import "testing"

func TestOrderItemModel_Insert(t *testing.T) {
	userModel := UserModel{DB: testDB}
	orderModel := OrderModel{DB: testDB}
	itemModel := OrderItemModel{DB: testDB}
	restaurantID := seedRestaurant(t)
	user := insertTestUser(t, userModel)
	order := insertTestOrder(t, orderModel, user.Id, restaurantID)
	dish := insertTestDish(t, DishModel{DB: testDB}, restaurantID)

	item := &OrderItem{
		OrderID:   order.ID,
		DishID:    dish.ID,
		DishName:  dish.Name,
		UnitPrice: dish.Price,
		Quantity:  2,
		Subtotal:  dish.Price * 2,
	}

	if err := itemModel.Insert(item); err != nil {
		t.Fatalf("Insert() error = %v", err)
	}

	if item.ID == 0 {
		t.Error("Insert() did not set item.ID")
	}
}

func TestOrderItemModel_InsertFromDish(t *testing.T) {
	userModel := UserModel{DB: testDB}
	orderModel := OrderModel{DB: testDB}
	itemModel := OrderItemModel{DB: testDB}
	restaurantID := seedRestaurant(t)
	user := insertTestUser(t, userModel)
	order := insertTestOrder(t, orderModel, user.Id, restaurantID)
	dish := insertTestDish(t, DishModel{DB: testDB}, restaurantID)

	item, err := itemModel.InsertFromDish(order.ID, dish, 3)
	if err != nil {
		t.Fatalf("InsertFromDish() error = %v", err)
	}

	if item.ID == 0 {
		t.Error("InsertFromDish() did not set item.ID")
	}
	if item.DishName != dish.Name {
		t.Errorf("InsertFromDish() DishName = %q, want %q", item.DishName, dish.Name)
	}
	if item.UnitPrice != dish.Price {
		t.Errorf("InsertFromDish() UnitPrice = %d, want %d", item.UnitPrice, dish.Price)
	}
	if item.Subtotal != dish.Price*3 {
		t.Errorf("InsertFromDish() Subtotal = %d, want %d", item.Subtotal, dish.Price*3)
	}
}

func TestOrderItemModel_GetForOrder(t *testing.T) {
	userModel := UserModel{DB: testDB}
	orderModel := OrderModel{DB: testDB}
	itemModel := OrderItemModel{DB: testDB}
	restaurantID := seedRestaurant(t)
	user := insertTestUser(t, userModel)
	order := insertTestOrder(t, orderModel, user.Id, restaurantID)
	dish := insertTestDish(t, DishModel{DB: testDB}, restaurantID)
	item, err := itemModel.InsertFromDish(order.ID, dish, 2)
	if err != nil {
		t.Fatalf("InsertFromDish() error = %v", err)
	}

	items, err := itemModel.GetForOrder(order.ID)
	if err != nil {
		t.Fatalf("GetForOrder() error = %v", err)
	}

	if len(items) != 1 {
		t.Fatalf("GetForOrder() returned %d items, want 1", len(items))
	}
	if items[0].ID != item.ID {
		t.Errorf("GetForOrder() item ID = %d, want %d", items[0].ID, item.ID)
	}
}

func TestOrderItemModel_Update(t *testing.T) {
	userModel := UserModel{DB: testDB}
	orderModel := OrderModel{DB: testDB}
	itemModel := OrderItemModel{DB: testDB}
	restaurantID := seedRestaurant(t)
	user := insertTestUser(t, userModel)
	order := insertTestOrder(t, orderModel, user.Id, restaurantID)
	dish := insertTestDish(t, DishModel{DB: testDB}, restaurantID)
	item, err := itemModel.InsertFromDish(order.ID, dish, 2)
	if err != nil {
		t.Fatalf("InsertFromDish() error = %v", err)
	}

	item.Quantity = 5
	item.Subtotal = item.UnitPrice * 5
	if err := itemModel.Update(item); err != nil {
		t.Fatalf("Update() error = %v", err)
	}

	items, err := itemModel.GetForOrder(order.ID)
	if err != nil {
		t.Fatalf("GetForOrder() after Update() error = %v", err)
	}
	if items[0].Quantity != 5 {
		t.Errorf("Update() Quantity = %d, want 5", items[0].Quantity)
	}
	if items[0].Subtotal != item.UnitPrice*5 {
		t.Errorf("Update() Subtotal = %d, want %d", items[0].Subtotal, item.UnitPrice*5)
	}
}

func TestOrderItemModel_Update_NotFound(t *testing.T) {
	model := OrderItemModel{DB: testDB}
	item := &OrderItem{ID: 999999, Quantity: 1, Subtotal: 100}

	err := model.Update(item)
	if err != ErrRecordNotFound {
		t.Errorf("Update() error = %v, want ErrRecordNotFound", err)
	}
}

func TestOrderItemModel_DeleteForOrder(t *testing.T) {
	userModel := UserModel{DB: testDB}
	orderModel := OrderModel{DB: testDB}
	itemModel := OrderItemModel{DB: testDB}
	restaurantID := seedRestaurant(t)
	user := insertTestUser(t, userModel)
	order := insertTestOrder(t, orderModel, user.Id, restaurantID)
	dish := insertTestDish(t, DishModel{DB: testDB}, restaurantID)
	if _, err := itemModel.InsertFromDish(order.ID, dish, 2); err != nil {
		t.Fatalf("InsertFromDish() error = %v", err)
	}

	if err := itemModel.DeleteForOrder(order.ID); err != nil {
		t.Fatalf("DeleteForOrder() error = %v", err)
	}

	items, err := itemModel.GetForOrder(order.ID)
	if err != nil {
		t.Fatalf("GetForOrder() after DeleteForOrder() error = %v", err)
	}
	if len(items) != 0 {
		t.Errorf("GetForOrder() after DeleteForOrder() returned %d items, want 0", len(items))
	}
}

func TestCalculateTotal(t *testing.T) {
	items := []*OrderItem{
		{Subtotal: 100},
		{Subtotal: 250},
		{Subtotal: 400},
	}

	total := CalculateTotal(items)
	if total != 750 {
		t.Errorf("CalculateTotal() = %d, want 750", total)
	}
}
