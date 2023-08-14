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
