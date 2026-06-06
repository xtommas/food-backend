<div align="center">
  <h1>Food delivery backend</h1>

![Go](https://img.shields.io/badge/go-%2300ADD8.svg?style=for-the-badge&logo=go&logoColor=white)
![Postgres](https://img.shields.io/badge/postgres-%23316192.svg?style=for-the-badge&logo=postgresql&logoColor=white)
![JWT](https://img.shields.io/badge/JWT-black?style=for-the-badge&logo=JSON%20web%20tokens)
![Postman](https://img.shields.io/badge/Postman-FF6C37?style=for-the-badge&logo=postman&logoColor=white)

<p align="left">Backend for a food delivery application, where admins manage restaurants, restaurant staff manage dishes and orders, and customers place orders.</p>

</div>

The API supports the following endpoints and actions:

| Method | URL Pattern                                        | Action                                          | Auth |
| :----- | :------------------------------------------------- | :---------------------------------------------- | :--- |
| GET    | /healthcheck                                       | Show application health and version information | Public |
| POST   | /users                                             | Register a customer account                     | Public |
| PUT    | /users/activate                                    | Activate an account                             | Public |
| PUT    | /users/password                                    | Reset a password with a password-reset token    | Public |
| GET    | /users/me                                          | Get the authenticated user                      | Activated user |
| PATCH  | /users/me                                          | Update the authenticated user's name or email   | Activated user |
| POST   | /users/me/photo                                    | Upload a user profile photo                     | Activated user |
| GET    | /users/me/photo                                    | Download the authenticated user's photo         | Activated user |
| POST   | /admin/promote                                     | Promote a user to admin                         | Admin |
| GET    | /restaurants                                       | List restaurants                                | `restaurants:read` |
| POST   | /restaurants                                       | Create a restaurant                             | Admin |
| GET    | /restaurants/:restaurant_id                        | Get one restaurant                              | `restaurants:read` |
| PATCH  | /restaurants/:restaurant_id                        | Update a restaurant                             | Admin |
| DELETE | /restaurants/:restaurant_id                        | Delete a restaurant                             | Admin |
| GET    | /restaurants/:restaurant_id/staff                  | List restaurant staff                           | Restaurant owner or admin |
| POST   | /restaurants/:restaurant_id/staff                  | Add or update a staff member                    | Restaurant owner or admin |
| DELETE | /restaurants/:restaurant_id/staff/:user_id         | Remove a staff member                           | Restaurant owner or admin |
| GET    | /restaurants/:restaurant_id/dishes                 | List dishes for a restaurant                    | `dishes:read` |
| POST   | /restaurants/:restaurant_id/dishes                 | Add a dish                                      | Restaurant staff or admin |
| GET    | /restaurants/:restaurant_id/dishes/:id             | Get one dish                                    | `dishes:read` |
| PATCH  | /restaurants/:restaurant_id/dishes/:id             | Update a dish                                   | Restaurant staff or admin |
| DELETE | /restaurants/:restaurant_id/dishes/:id             | Delete a dish                                   | Restaurant staff or admin |
| POST   | /restaurants/:restaurant_id/dishes/:id/photo/      | Upload a dish photo                             | Restaurant staff or admin |
| GET    | /restaurants/:restaurant_id/dishes/:id/photo/      | Download a dish photo                           | `dishes:read` |
| POST   | /restaurants/:restaurant_id/orders                 | Create an order for a restaurant                | Activated user |
| GET    | /restaurants/:restaurant_id/orders                 | List restaurant orders                          | Restaurant staff or admin |
| GET    | /restaurants/:restaurant_id/orders/:order_id       | Get one restaurant order with items             | Restaurant staff or admin |
| PATCH  | /restaurants/:restaurant_id/orders/:order_id       | Update an order status                          | Restaurant staff or admin |
| POST   | /restaurants/:restaurant_id/orders/:order_id/items | Add an item to an order                         | Activated user |
| GET    | /restaurants/:restaurant_id/orders/:order_id/items | List items for one restaurant order             | Restaurant staff or admin |
| GET    | /users/me/orders                                   | List the authenticated user's orders            | Activated user |
| GET    | /users/me/orders/:order_id                         | Get one authenticated-user order with items     | Activated user |
| GET    | /users/me/orders/:order_id/items                   | List items for one authenticated-user order     | Activated user |
| POST   | /tokens/authentication                             | Create an authentication token                  | Public |
| POST   | /tokens/password-reset                             | Create a password-reset token                   | Public |
| POST   | /tokens/activation                                 | Create a new activation token                   | Public |
| GET    | /debug/vars                                        | Display application metrics                     | Public |

Filtering:

- Dishes: `?available=true/false`, `?name=pizza`, `?categories=pizza,vegetarian`, `?sort=id/-id/name/-name/price/-price/available/-available`.
- Orders: `?status=pending/confirmed/preparing/ready/delivered/cancelled`, `?sort=id/-id/total/-total/status/-status`.
- Pagination uses `?page=1&page_size=20` where list endpoints support pagination.

Prices and totals are stored and returned as integer cents. For example, `1299` means `$12.99`.

## ⚙️ Setup

You'll need to have [Docker](https://www.docker.com/) installed.

First, create a new `.env` file following `.env.example`.

Then, run the app with:

```bash
make docker/up
```

This starts the database, runs migrations, and starts the API on `http://localhost:4000`.

## 🔧 Makefile commands

| Command | Description |
|---|---|
| `make help` | Print the help message with all available commands |
| `make docker/up` | Build Docker images and start all services in detached mode |
| `make docker/up/attached` | Build Docker images and start all services while following logs |
| `make docker/down` | Stop all running services |
| `make docker/down/volumes` | Stop all services and remove Docker volumes, wiping DB data |
| `make docker/nuke` | Stop services, remove volumes, images, and orphan containers |
| `make docker/logs` | Follow logs for all services |
| `make docker/logs/api` | Follow logs for the API service only |
| `make docker/psql` | Open a `psql` shell connected to the containerized PostgreSQL database |
| `make db/migrate/new name=<migration_name>` | Create a new sequential SQL migration file |
| `make db/migrate/up` | Run database migrations against the containerized database |
| `make db/migrate/down` | Roll back all database migrations |
| `make db/schema` | Dump the current database schema to `schema.sql` |
| `make build/api` | Build the `cmd/api` application locally and generate a Linux AMD64 binary |
| `make test` | Start the test DB and run all tests |

## 🍕 Examples

Set the base URL once:

```bash
BASE_URL=http://localhost:4000
```

For protected examples, replace `$CUSTOMER_TOKEN`, `$ADMIN_TOKEN`, or `$STAFF_TOKEN` with a token returned by `POST /tokens/authentication`.

### Register a customer

```bash
curl --request POST \
  --url "$BASE_URL/users" \
  --header 'Content-Type: application/json' \
  --data '{
    "name": "Tomas",
    "email": "tomas@example.com",
    "password": "password123"
  }'
```

```json
{
  "activation_token": "<activation-jwt>",
  "user": {
    "id": 8,
    "created_at": "2026-06-06T12:00:00Z",
    "name": "Tomas",
    "email": "tomas@example.com",
    "activated": false,
    "role": "customer"
  }
}
```

### Activate and authenticate

```bash
curl --request PUT \
  --url "$BASE_URL/users/activate" \
  --header 'Content-Type: application/json' \
  --data '{
    "token": "<activation-jwt>"
  }'
```

```bash
curl --request POST \
  --url "$BASE_URL/tokens/authentication" \
  --header 'Content-Type: application/json' \
  --data '{
    "email": "tomas@example.com",
    "password": "password123"
  }'
```

```json
{
  "authentication_token": "<authentication-jwt>"
}
```

### Create a restaurant

Restaurants are no longer user accounts. An admin creates the restaurant record.

```bash
curl --request POST \
  --url "$BASE_URL/restaurants" \
  --header "Authorization: Bearer $ADMIN_TOKEN" \
  --header 'Content-Type: application/json' \
  --data '{
    "name": "Roma Pizza",
    "photo": "images/restaurants/roma.jpg",
    "address": "123 Market Street",
    "city": "Buenos Aires",
    "province": "Buenos Aires",
    "country": "Argentina",
    "latitude": -34.603722,
    "longitude": -58.381592
  }'
```

```json
{
  "restaurant": {
    "id": 7,
    "name": "Roma Pizza",
    "photo": "images/restaurants/roma.jpg",
    "address": "123 Market Street",
    "city": "Buenos Aires",
    "province": "Buenos Aires",
    "country": "Argentina",
    "latitude": -34.603722,
    "longitude": -58.381592,
    "created_at": "2026-06-06T12:05:00Z"
  }
}
```

### Add restaurant staff

An admin or restaurant owner can add staff. Use `owner` for the first manager of a restaurant, and `staff` for regular staff.

```bash
curl --request POST \
  --url "$BASE_URL/restaurants/7/staff" \
  --header "Authorization: Bearer $ADMIN_TOKEN" \
  --header 'Content-Type: application/json' \
  --data '{
    "email": "staff@example.com",
    "role": "owner"
  }'
```

```json
{
  "message": "staff member successfully added"
}
```

### List restaurants

```bash
curl --request GET \
  --url "$BASE_URL/restaurants" \
  --header "Authorization: Bearer $CUSTOMER_TOKEN"
```

```json
{
  "restaurants": [
    {
      "id": 7,
      "name": "Roma Pizza",
      "photo": "images/restaurants/roma.jpg",
      "address": "123 Market Street",
      "city": "Buenos Aires",
      "province": "Buenos Aires",
      "country": "Argentina",
      "latitude": -34.603722,
      "longitude": -58.381592,
      "created_at": "2026-06-06T12:05:00Z"
    }
  ]
}
```

### Add a dish

Prices are integer cents. This example creates a `$12.99` dish.

```bash
curl --request POST \
  --url "$BASE_URL/restaurants/7/dishes" \
  --header "Authorization: Bearer $STAFF_TOKEN" \
  --header 'Content-Type: application/json' \
  --data '{
    "name": "Neapolitan Pizza",
    "price": 1299,
    "description": "Tomato, mozzarella, basil, and olive oil",
    "categories": ["pizza", "vegetarian"]
  }'
```

```json
{
  "dish": {
    "id": 5,
    "restaurant_id": 7,
    "name": "Neapolitan Pizza",
    "price": 1299,
    "description": "Tomato, mozzarella, basil, and olive oil",
    "categories": ["pizza", "vegetarian"],
    "available": true,
    "updated_at": "2026-06-06T12:10:00Z"
  }
}
```

### List dishes for a restaurant

```bash
curl --request GET \
  --url "$BASE_URL/restaurants/7/dishes?available=true&categories=pizza&sort=price" \
  --header "Authorization: Bearer $CUSTOMER_TOKEN"
```

```json
{
  "dishes": [
    {
      "id": 5,
      "restaurant_id": 7,
      "name": "Neapolitan Pizza",
      "price": 1299,
      "description": "Tomato, mozzarella, basil, and olive oil",
      "categories": ["pizza", "vegetarian"],
      "available": true,
      "updated_at": "2026-06-06T12:10:00Z"
    }
  ],
  "metadata": {
    "current_page": 1,
    "page_size": 20,
    "first_page": 1,
    "last_page": 1,
    "total_records": 1
  }
}
```

### Upload a dish photo

```bash
curl --request POST \
  --url "$BASE_URL/restaurants/7/dishes/5/photo/" \
  --header "Authorization: Bearer $STAFF_TOKEN" \
  --form 'photo=@/path/to/neapolitan-pizza.jpg'
```

### Create an order

New orders start as `pending`.

```bash
curl --request POST \
  --url "$BASE_URL/restaurants/7/orders" \
  --header "Authorization: Bearer $CUSTOMER_TOKEN" \
  --header 'Content-Type: application/json' \
  --data '{
    "address": "Apartment 5D"
  }'
```

```json
{
  "order": {
    "id": 11,
    "user_id": 8,
    "restaurant_id": 7,
    "total": 0,
    "address": "Apartment 5D",
    "created_at": "2026-06-06T12:20:00Z",
    "updated_at": "2026-06-06T12:20:00Z",
    "status": "pending"
  }
}
```

### Add an item to an order

Order items snapshot the dish name and price at the time the item is added.

```bash
curl --request POST \
  --url "$BASE_URL/restaurants/7/orders/11/items" \
  --header "Authorization: Bearer $CUSTOMER_TOKEN" \
  --header 'Content-Type: application/json' \
  --data '{
    "dish_id": 5,
    "quantity": 2
  }'
```

```json
{
  "order_item": {
    "id": 18,
    "order_id": 11,
    "dish_id": 5,
    "dish_name": "Neapolitan Pizza",
    "unit_price": 1299,
    "quantity": 2,
    "subtotal": 2598
  }
}
```

### Get customer orders

```bash
curl --request GET \
  --url "$BASE_URL/users/me/orders?status=pending&sort=-id" \
  --header "Authorization: Bearer $CUSTOMER_TOKEN"
```

```json
{
  "orders": [
    {
      "order": {
        "id": 11,
        "user_id": 8,
        "restaurant_id": 7,
        "total": 2598,
        "address": "Apartment 5D",
        "created_at": "2026-06-06T12:20:00Z",
        "updated_at": "2026-06-06T12:22:00Z",
        "status": "pending"
      },
      "items": [
        {
          "dish": "Neapolitan Pizza",
          "quantity": 2,
          "subtotal": 2598
        }
      ]
    }
  ],
  "metadata": {
    "current_page": 1,
    "page_size": 50,
    "first_page": 1,
    "last_page": 1,
    "total_records": 1
  }
}
```

### Update an order status

Restaurant staff can move orders through these transitions:

`pending -> confirmed -> preparing -> ready -> delivered`

Staff can also cancel `pending` or `confirmed` orders.

```bash
curl --request PATCH \
  --url "$BASE_URL/restaurants/7/orders/11" \
  --header "Authorization: Bearer $STAFF_TOKEN" \
  --header 'Content-Type: application/json' \
  --data '{
    "status": "confirmed"
  }'
```

```json
{
  "order": {
    "id": 11,
    "user_id": 8,
    "restaurant_id": 7,
    "total": 2598,
    "address": "Apartment 5D",
    "created_at": "2026-06-06T12:20:00Z",
    "updated_at": "2026-06-06T12:30:00Z",
    "status": "confirmed"
  }
}
```

### Get logged-in user info

```bash
curl --request GET \
  --url "$BASE_URL/users/me" \
  --header "Authorization: Bearer $CUSTOMER_TOKEN"
```

```json
{
  "user": {
    "id": 8,
    "created_at": "2026-06-06T12:00:00Z",
    "name": "Tomas",
    "email": "tomas@example.com",
    "activated": true,
    "role": "customer"
  }
}
```
