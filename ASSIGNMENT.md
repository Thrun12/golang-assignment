# Golang assignment 💻

This assignment will form the basis for the job interview. You will show-case your implementation and talk about the choices you made and obstacles if you encountered any.

## Prequisites ✔️

1. Git
2. Golang
3. Docker (optional)

## Assignment 📝

Fork this repository and commit your code to your own repository.

1. Create a simple Golang application that will serve a RESTful API with a resource of your choice (pets, books, memes, etc). The endpoints should support Create, Read, Update and Delete operations.
2. The data should be persisted in a database like SQLite, MySQL, PostgreSQL, etc.

If you are stuck or want some inspiration see the tips below.

### Optionals

1. Add tests for you REST API.
2. Create a dockerfile for your application.

## Tips 🧞

### Application example

_Endpoints:_

```sh
# Get all pets
GET /pets

# Get a specific pet
GET /pets/:id

# Create a pet
POST /pets

# Update a pet
PUT /pets/:id

# Delete a pet
DELETE /pets/:id
```

_Model:_

```json
{
    "id": 1,
    "name": "Fluffy",
    "age": 3,
    "breed": "cat"
}
```

### Database

This package can be used to connect to a local SQLite database:

-   https://github.com/mattn/go-sqlite3

If you want an ORM (Object Relational Mapping) library, you can use GORM:

-   https://gorm.io/docs/connecting_to_the_database.html#SQLite

### REST API

For creating a HTTP server, you can use echo:

-   https://github.com/labstack/echo

### REST API testing

Documentation for testing echo endpoints:

-   https://echo.labstack.com/guide/testing/
