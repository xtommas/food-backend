<div align="center">
  <h1>Food delivery backend</h1>

![Go](https://img.shields.io/badge/go-%2300ADD8.svg?style=for-the-badge&logo=go&logoColor=white)
![Postgres](https://img.shields.io/badge/postgres-%23316192.svg?style=for-the-badge&logo=postgresql&logoColor=white)
![JWT](https://img.shields.io/badge/JWT-black?style=for-the-badge&logo=JSON%20web%20tokens)
![Postman](https://img.shields.io/badge/Postman-FF6C37?style=for-the-badge&logo=postman&logoColor=white)

<p align="left">Backend for a food delivery application, where restaurants can register, add dishes and handle orders; and users can make orders to the restaurants.</p>

</div>

The API supports the following endpoints and actions:

<div align="center">

| Method | URL Pattern                                        | Action                                                                              | Filtering                                                                                                                                            |
| :----- | :------------------------------------------------- | :---------------------------------------------------------------------------------- | :--------------------------------------------------------------------------------------------------------------------------------------------------- |
| GET    | /healthcheck                                       | Show application health and version information                                     |                                                                                                                                                      |
| GET    | /restaurants                                       | Show details of all the restaurants                                                 |                                                                                                                                                      |
| GET    | /restaurants/:restaurant_id/dishes                 | Get all the dishes for a restaurant                                                 | Filter by `?available=true/false`, `?name=something`, and `?categories=something`. Also, `?sort=id/-id/name/-name/price/-price/available/-available` |
| POST   | /restaurants/:restaurant_id/dishes                 | Lets restaurants add dishes to their menu                                           |                                                                                                                                                      |
| GET    | /restaurants/:restaurant_id/dishes/:id             | Get a specific dish from a restaurant                                               |                                                                                                                                                      |
| PATCH  | /restaurants/:restaurant_id/dishes/:id             | Update the details of a dish                                                        |                                                                                                                                                      |
| DELETE | /restaurants/:restaurant_id/dishes/:id             | Delete a dish                                                                       |                                                                                                                                                      |
| POST   | /restaurants/:restaurant_id/dishes/:id/photo/      | Add a picture for a dish                                                            |                                                                                                                                                      |
| GET    | /restaurants/:restaurant_id/dishes/:id/photo/      | Get the picture for a dish                                                          |                                                                                                                                                      |
| POST   | /restaurants/:restaurant_id/orders                 | Make and order to a restaurant                                                      |                                                                                                                                                      |
| GET    | /restaurants/:restaurant_id/orders                 | Get the orders for a restaurant                                                     | Filter by `?status=created/in%20progress/ready/delivered/cancelled` and `?sort=id/-id/total/-total/status/-status`                                   |
| GET    | /restaurants/:restaurant_id/orders/:order_id       | Get details for a specific order, including items ordered                           | Filter by `?status=created/in%20progress/ready/delivered/cancelled` and `?sort=id/-id/total/-total/status/-status`                                   |
| PATCH  | /restaurants/:restaurant_id/orders/:order_id       | Update the status of an order                                                       |                                                                                                                                                      |
| POST   | /restaurants/:restaurant_id/orders/:order_id/items | Add items to an order                                                               |                                                                                                                                                      |
| GET    | /restaurants/:restaurant_id/orders/:order_id/items | Get the items for a specific order                                                  |                                                                                                                                                      |
| POST   | /users                                             | Register a new user (customer or restaurant)                                        |                                                                                                                                                      |
| PUT    | /users/activate                                    | Activate a user's account                                                           |                                                                                                                                                      |
| PUT    | /users/password                                    | Change user password                                                                |                                                                                                                                                      |
| GET    | /users/me                                          | Get details of the authenticated user                                               |                                                                                                                                                      |
| POST   | /users/me/photo                                    | Set a profile picture for the authenticated user                                    |                                                                                                                                                      |
| GET    | /users/me/photo                                    | Get the picture for the authenticated user                                          |                                                                                                                                                      |
| PATCH  | /users/me                                          | Update email and name for the currently authenticated user                          |                                                                                                                                                      |
| GET    | /users/me/orders                                   | Get orders made by currently authenticated user                                     | Filter by `?status=created/in%20progress/ready/delivered/cancelled` and `?sort=id/-id/total/-total/status/-status`                                   |
| GET    | /users/me/orders/:order_id                         | Get a specific order made by the authenticated user, including ordered items        |                                                                                                                                                      |
| GET    | /users/me/orders/:order_id/items                   | Get the items of a specific order made by the authenticated user                    |                                                                                                                                                      |
| POST   | /tokens/authentication                             | Get an authentication token for a user, with email and password on the request body |                                                                                                                                                      |
| POST   | /tokens/password-reset                             | Get a password reset token for a user, with email on the request body               |                                                                                                                                                      |
| POST   | /tokens/activation                                 | Get an activation token for a user, with email on the request body                  |                                                                                                                                                      |
| GET    | /debug/vars                                        | Display application metrics                                                         |                                                                                                                                                      |

</div>

## ⚙️ Setup

You'll need to set up the PostgreSQL database using the `databse.sql` file and then running

```bash
$ make db/migrations/up
```

Also, you should put the database Data Source Name and the JWT secret in a `.envrc` file, it should look something like this:

```bash
export DB_DSN=postgres://food:password@localhost/food?sslmode=disable
export JWT_SECRET=e7X29mLufqQNGGyEFl5rpSHSs_RZUtt69Maur82U_iSc4PIjFT2Dtt9r2U4VAO5odfp7OPDeg5TN4o0-wuNRZA
```

Now, you can run the application with

```bash
$ make run/api
```

Additionally, you can build the application using

```bash
$ make build/api
```

and run it using:

```bash
./bin/api -db-dsn=postgres://food:yourpassword@localhost/food?sslmode=disable
```

## 🍕 Examples

### Creating a user

Request

```bash
curl --request POST \
  --url http://localhost:4000/users \
  --header 'Content-Type: application/json' \
  --data '{
	"name": "Tomas",
	"email": "tomas@example.com",
	"password": "password",
	"role": "customer"
}'
```

Response

```json
{
  "activation_token": "eyJhbGciOiJIUzI1NiJ9.eyJhdWQiOlsiZ2l0aHViLmNvbS94dG9tbWFzL2Zvb2QtYmFja2VuZCJdLCJleHAiOjE2OTE4Njg2ODUuNjk2OTY0LCJpYXQiOjE2OTE3ODIyODUuNjk2OTY0LCJpc3MiOiJnaXRodWIuY29tL3h0b21tYXMvZm9vZC1iYWNrZW5kIiwibmJmIjoxNjkxNzgyMjg1LjY5Njk2NCwic2NvcGUiOiJhY3RpdmF0aW9uIiwic3ViIjoiOCJ9.bGalg2RSaSpOTLP7E0msNhhp4G3eDLTxNE52jotPXDQ",
  "user": {
    "id": 8,
    "created_at": "2023-08-11T16:31:26-03:00",
    "name": "Tomas",
    "email": "tomas@example.com",
    "activated": false,
    "role": "customer"
  }
}
```

### Activating a user

Request

```bash
curl --request PUT \
  --url http://localhost:4000/users/activate \
  --header 'Content-Type: application/json' \
  --data '{
	"token": "eyJhbGciOiJIUzI1NiJ9.eyJhdWQiOlsiZ2l0aHViLmNvbS94dG9tbWFzL2Zvb2QtYmFja2VuZCJdLCJleHAiOjE2OTE4Njg2ODUuNjk2OTY0LCJpYXQiOjE2OTE3ODIyODUuNjk2OTY0LCJpc3MiOiJnaXRodWIuY29tL3h0b21tYXMvZm9vZC1iYWNrZW5kIiwibmJmIjoxNjkxNzgyMjg1LjY5Njk2NCwic2NvcGUiOiJhY3RpdmF0aW9uIiwic3ViIjoiOCJ9.bGalg2RSaSpOTLP7E0msNhhp4G3eDLTxNE52jotPXDQ"
}'
```

Response

```json
{
  "user": {
    "id": 8,
    "created_at": "2023-08-11T16:31:26-03:00",
    "name": "Tomas",
    "email": "tomas@example.com",
    "activated": true,
    "role": "customer"
  }
}
```

### Authenticating a user

Request

```bash
curl --request POST \
  --url http://localhost:4000/tokens/authentication \
  --header 'Content-Type: application/json' \
  --data '{
	"email": "tomas@example.com",
	"password": "password"
}'
```

Response

```json
{
  "authentication_token": "eyJhbGciOiJIUzI1NiJ9.eyJhdWQiOlsiZ2l0aHViLmNvbS94dG9tbWFzL2Zvb2QtYmFja2VuZCJdLCJleHAiOjE2OTIxMjU5MDIuNDI2NzU0MiwiaWF0IjoxNjkyMDM5NTAyLjQyNjc1NDIsImlzcyI6ImdpdGh1Yi5jb20veHRvbW1hcy9mb29kLWJhY2tlbmQiLCJuYmYiOjE2OTIwMzk1MDIuNDI2NzU0Miwic2NvcGUiOiJhdXRoZW50aWNhdGlvbiIsInN1YiI6IjcifQ.TAhf_qF6uhyRdZ-IbDsT4JSwUCk5JfHg9jNXzeHNTgo"
}
```

> **Note**
> You can follow the same workflow to create a restaurant, just change the role to "restaurant" in the request body

### Making an order

Request

```bash
curl --request POST \
  --url http://localhost:4000/restaurants/7/orders \
  --header 'Authorization: Bearer eyJhbGciOiJIUzI1NiJ9.eyJhdWQiOlsiZ2l0aHViLmNvbS94dG9tbWFzL2Zvb2QtYmFja2VuZCJdLCJleHAiOjE2OTIxMTk2MTYuNTA3NzIzMywiaWF0IjoxNjkyMDMzMjE2LjUwNzcyMzMsImlzcyI6ImdpdGh1Yi5jb20veHRvbW1hcy9mb29kLWJhY2tlbmQiLCJuYmYiOjE2OTIwMzMyMTYuNTA3NzIzMywic2NvcGUiOiJhdXRoZW50aWNhdGlvbiIsInN1YiI6IjgifQ.c1RVx6S5r5iEyTWnZQCbxIYRz-uaaUcU1zUrCkiiPGo' \
  --header 'Content-Type: application/json' \
  --data '{
	"address": "Apartment 5D"
}'
```

Response

```json
{
  "order": {
    "id": 5,
    "user_id": 8,
    "restaurant_id": 7,
    "total": "$0",
    "address": "Apartment 5D",
    "created_at": "0001-01-01T00:00:00Z",
    "status": "created"
  }
}
```

### Adding items to the order

Request

```bash
curl --request POST \
  --url http://localhost:4000/restaurants/7/orders/5/items \
  --header 'Authorization: Bearer eyJhbGciOiJIUzI1NiJ9.eyJhdWQiOlsiZ2l0aHViLmNvbS94dG9tbWFzL2Zvb2QtYmFja2VuZCJdLCJleHAiOjE2OTIxMTk2MTYuNTA3NzIzMywiaWF0IjoxNjkyMDMzMjE2LjUwNzcyMzMsImlzcyI6ImdpdGh1Yi5jb20veHRvbW1hcy9mb29kLWJhY2tlbmQiLCJuYmYiOjE2OTIwMzMyMTYuNTA3NzIzMywic2NvcGUiOiJhdXRoZW50aWNhdGlvbiIsInN1YiI6IjgifQ.c1RVx6S5r5iEyTWnZQCbxIYRz-uaaUcU1zUrCkiiPGo' \
  --header 'Content-Type: application/json' \
  --data '{
	"dish_id": 5,
	"quantity": 2
}'
```

Response

```json
{
  "order_item": {
    "id": 6,
    "order_id": 5,
    "dish_id": 5,
    "quantity": 2,
    "subtotal": "$8000"
  }
}
```

### Getting orders for a user

Request

```bash
curl --request GET \
  --url http://localhost:4000/users/me/orders \
  --header 'Authorization: Bearer eyJhbGciOiJIUzI1NiJ9.eyJhdWQiOlsiZ2l0aHViLmNvbS94dG9tbWFzL2Zvb2QtYmFja2VuZCJdLCJleHAiOjE2OTIxMTk2MTYuNTA3NzIzMywiaWF0IjoxNjkyMDMzMjE2LjUwNzcyMzMsImlzcyI6ImdpdGh1Yi5jb20veHRvbW1hcy9mb29kLWJhY2tlbmQiLCJuYmYiOjE2OTIwMzMyMTYuNTA3NzIzMywic2NvcGUiOiJhdXRoZW50aWNhdGlvbiIsInN1YiI6IjgifQ.c1RVx6S5r5iEyTWnZQCbxIYRz-uaaUcU1zUrCkiiPGo'
```

Response

```json
{
  "metadata": {
    "current_page": 1,
    "page_size": 50,
    "first_page": 1,
    "last_page": 1,
    "total_records": 3
  },
  "orders": [
    {
      "order": {
        "id": 3,
        "user_id": 8,
        "restaurant_id": 7,
        "total": "$13000",
        "address": "Apartment 5D",
        "created_at": "2023-08-14T12:09:39-03:00",
        "status": "delivered"
      },
      "items": [
        {
          "dish": "Cheese Pizza",
          "quantity": 3,
          "subtotal": "$9000"
        },
        {
          "dish": "Neapolitan Pizza",
          "quantity": 1,
          "subtotal": "$4000"
        }
      ]
    },
    {
      "order": {
        "id": 4,
        "user_id": 8,
        "restaurant_id": 7,
        "total": "$13000",
        "address": "Apartment 5D",
        "created_at": "2023-08-14T15:54:25-03:00",
        "status": "delivered"
      },
      "items": [
        {
          "dish": "Neapolitan Pizza",
          "quantity": 1,
          "subtotal": "$4000"
        },
        {
          "dish": "Cheese Pizza",
          "quantity": 3,
          "subtotal": "$9000"
        }
      ]
    },
    {
      "order": {
        "id": 5,
        "user_id": 8,
        "restaurant_id": 7,
        "total": "$17000",
        "address": "Apartment 5D",
        "created_at": "2023-08-14T16:05:37-03:00",
        "status": "in progress"
      },
      "items": [
        {
          "dish": "Cheese Pizza",
          "quantity": 3,
          "subtotal": "$9000"
        },
        {
          "dish": "Neapolitan Pizza",
          "quantity": 2,
          "subtotal": "$8000"
        }
      ]
    }
  ]
}
```

### Getting all available restaurants

Request

```shell
curl --request GET \
  --url http://localhost:4000/restaurants \
  --header 'Authorization: Bearer eyJhbGciOiJIUzI1NiJ9.eyJhdWQiOlsiZ2l0aHViLmNvbS94dG9tbWFzL2Zvb2QtYmFja2VuZCJdLCJleHAiOjE2OTIxMzkwOTEuNDcwMjM0MiwiaWF0IjoxNjkyMDUyNjkxLjQ3MDIzNDIsImlzcyI6ImdpdGh1Yi5jb20veHRvbW1hcy9mb29kLWJhY2tlbmQiLCJuYmYiOjE2OTIwNTI2OTEuNDcwMjM0Miwic2NvcGUiOiJhdXRoZW50aWNhdGlvbiIsInN1YiI6IjgifQ.KQJjOi6xjsx3TFS51Cm4mMvCn1_ViQUacK1W83vb8vY'
```

Response

```json
{
  "restaurants": [
    {
      "id": 7,
      "photo": "images/users/7.jpg",
      "created_at": "2023-08-11T16:29:06-03:00",
      "name": "Domino's",
      "email": "dominos@example.com",
      "activated": true,
      "role": "restaurant"
    },
    {
      "id": 5,
      "photo": "images/users/5.jpg",
      "created_at": "2023-08-11T16:26:09-03:00",
      "name": "McDonald's",
      "email": "mcdonalds@example.com",
      "activated": true,
      "role": "restaurant"
    }
  ]
}
```

### Adding a new dish

Request

```shell
curl --request POST \
  --url http://localhost:4000/restaurants/7/dishes \
  --header 'Authorization: Bearer eyJhbGciOiJIUzI1NiJ9.eyJhdWQiOlsiZ2l0aHViLmNvbS94dG9tbWFzL2Zvb2QtYmFja2VuZCJdLCJleHAiOjE2OTIxMTk1NjcuNjMxNTQ5MSwiaWF0IjoxNjkyMDMzMTY3LjYzMTU0OTEsImlzcyI6ImdpdGh1Yi5jb20veHRvbW1hcy9mb29kLWJhY2tlbmQiLCJuYmYiOjE2OTIwMzMxNjcuNjMxNTQ5MSwic2NvcGUiOiJhdXRoZW50aWNhdGlvbiIsInN1YiI6IjcifQ.UZymbl2dOg8zm2IBlYgmYyCtmgU_3fmHnHtZbIbb23E' \
  --header 'Content-Type: application/json' \
  --data '{
	"name": "Neapolitan Pizza",
	"price": "4000",
	"description": "Neapolitan pizza with tomatoes and mozzarella cheese",
	"categories": ["Pizza"]
}'
```

Response

```json
{
  "dish": {
    "id": 5,
    "restaurant_id": 7,
    "name": "Neapolitan Pizza",
    "price": "$4000",
    "description": "Neapolitan pizza with tomatoes and mozzarella cheese",
    "category": ["Pizza"],
    "available": true
  }
}
```

### Adding an image to a dish

Request

```shell
curl --request POST \
  --url http://localhost:4000/restaurants/7/dishes/5/photo \
  --header 'Authorization: Bearer eyJhbGciOiJIUzI1NiJ9.eyJhdWQiOlsiZ2l0aHViLmNvbS94dG9tbWFzL2Zvb2QtYmFja2VuZCJdLCJleHAiOjE2OTE4NzQ5NTMuNzEyMjIwMiwiaWF0IjoxNjkxNzg4NTUzLjcxMjIyMDIsImlzcyI6ImdpdGh1Yi5jb20veHRvbW1hcy9mb29kLWJhY2tlbmQiLCJuYmYiOjE2OTE3ODg1NTMuNzEyMjIwMiwic2NvcGUiOiJhdXRoZW50aWNhdGlvbiIsInN1YiI6IjcifQ.HByX7enm41hH8isboADOSPmVSbfuHzyZj50mSAPK188' \
  --header 'Content-Type: multipart/form-data' \
  --form 'photo=@C:\Users\tomas\Downloads\neapolitan_pizza.jpg'
```

### Getting logged in user/restaurant info

Request

```shell
curl --request GET \
  --url http://localhost:4000/users/me \
  --header 'Authorization: Bearer eyJhbGciOiJIUzI1NiJ9.eyJhdWQiOlsiZ2l0aHViLmNvbS94dG9tbWFzL2Zvb2QtYmFja2VuZCJdLCJleHAiOjE2OTIwMzIwMTIuMTgyMTQyLCJpYXQiOjE2OTE5NDU2MTIuMTgyMTQyLCJpc3MiOiJnaXRodWIuY29tL3h0b21tYXMvZm9vZC1iYWNrZW5kIiwibmJmIjoxNjkxOTQ1NjEyLjE4MjE0Miwic2NvcGUiOiJhdXRoZW50aWNhdGlvbiIsInN1YiI6IjUifQ.KHDeD8jD5dcGYKsVauJO95S1hk0CR7qmv43hqYB5RhU'
```

Response

```json
{
  "user": {
    "id": 5,
    "photo": "images/users/5.jpg",
    "created_at": "2023-08-11T16:26:09-03:00",
    "name": "McDonald's",
    "email": "mcdonalds@example.com",
    "activated": true,
    "role": "restaurant"
  }
}
```
