# Go-Template

## Written in GoLang

This repository will serve as a base web application in Go.

- Built in Go version 1.19
- Uses the [chi](https://github.com/go-chi/chi/v5) router
- Uses [godotenv](https://github.com/joho/godotenv) for environmental variables
- Uses [Swaggo](https://github.com/swaggo/swag) to generate API documentation
- Uses [Go-Validator](https://github.com/asaskevich/govalidator) for validating incoming data

### Environment

You will be required to setup a .env file in the root of the project folder. This will need to contain the database details, and the encryption (HMAC) / JSON web token session secrets.

```
DB_USER=postgres
DB_PASS=
DB_HOST=localhost
DB_PORT=5432
DB_NAME=
SESSIONS_SECRET_KEY=
HMAC_SECRET=
```

### Database (Object Relational Management)

- Uses [Gorm](https://gorm.io) for ORM (Postgres)

### Security / Authentication / Authorization

- Uses [bcrypt](https://golang.org/x/crypto) for password hashing/decrypting
- Uses [Golang-jwt](https://github.com/golang-jwt/jwt) for JSON Web Token authentication
- Uses [casbin](https://github.com/casbin/casbin/v2) for Role based access control authorization

## To run Go server

```
go run ./cmd
```

## To run using Docker

"container-name" is typically the github address of your project. (ie. dmawardi/go-template)

```
<!-- Builds docker image -->
docker build -t container-name .


<!-- runs docker image and matches port -->
docker run --publish 8080:8080 container-name
```

---

## How to use

Follow these steps to add a feature to the API. This template uses the clean architecture pattern

1. Build schema and auto migrate in ./internal/db using ORM instructions below.
2. Build repository in ./internal/repository which is the interaction between the DB and the application. This should use a struct with receiver functions.
3. Build incoming DTO models for Create and Update JSON requests in ./internal/models
4. Build service in ./internal/service that uses the repository and applies business logic.
5. Build the controller (handler) in ./internal/controller that accepts the request, performs data validation, then sends to the service to interact with database.
6. Add validation to handler using govalidator. This is done by adding `valid:""` key-value pairs to struct DTO definitions (/internal/models) that are being passed into the ValidateStruct function (used in controllers).
7. Add the new controller to the API struct in the ./internal/routes/routes.go file. This allows it to be used within the Routes function in the same file. Build routes to use the handlers that have been created in step 4 using the api struct.
8. Update the ApiSetup function in the ./cmd/main.go file to build the new repository, service, and controller.
9. Add the route to the RBAC authorization policy file (./internal/auth/defaultPolicy.go)
10. (Testing) For e2e testing, you will need to update the controllers_test.go file in ./internal/controller. Updates are required in the buildAPI, setupDatabase & setupDBAuthAppModels functions

---

## Testing

To run all tests use below command.

```
go test ./...
```

This will run all files that match the testing file naming convention (\*\_test.go).

Tests for repositories, services, and controllers should be in their respsective directories. The controllers folder consists of E2E tests.  
The setup for these tests is in the controllers_test.go file.

Upon adding a new module:
-Make sure to build a DB struct that contains the modules of: repo, service, & controller. This object should then be added to the testDbRepo which serves as the test connection that will be serving the requests for the DB and API.

#### Additional flags

- "-V" prints more detailed results
- "-cover" will provide a test coverage report

---

## API documentation

API documentation is auto generated using markdown within code. This is achieved using Swag.
You must navigate to folder with main.go to generate.

The below commands must be used upon making changes to the API in order to regenerate the API docs.

- "-d" directory flag allows custom directory to be used
- "--pd" flag parses dependecies as well

```
swag init -d ./cmd --pd
```

This will update API documentation generated in the ./docs folder. It is served on path /swagger

## To use Database ORM

To edit schemas: Go to ./internal/db/schemas.go

The schemas are Structs based off of gorm.Model.

After creating the schema in schemas.go, go to db.go and add to automigrate.

Note for creating and updating using GORM: Relationship data that does not yet exist will be created as a new entry. However, if you try to edit an existing record, it will not allow you to.

## Role based access control (RBAC) settings

The authorization settings are found in the ./internal/auth/defaultPolicy.go file.
This data structure is used by the setupcasbin policy to implement policy in DB upon server start.

SetupCasbinPolicy functions in a way where it adds policies only if they're not found already.

Format of policy: Subject, Object, Action (ie. "Who" is accessing "DB object" to commit "CRUD action")
