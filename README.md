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
3. Build service in ./internal/service that uses the repository and applies business logic.
4. Build the controller (handler) in ./internal/controller that accepts the request, performs data validation, then sends to the service to interact with database.
5. Add the new controller to the API struct in the ./cmd/routes.go file. This allows it to be used within the routes.
6. Update routes in ./cmd/routes.go to use the handler that has been created in step 5.
7. Add validation to handler using govalidator. This functions by adding `valid:""` key-value pairs to struct DTO definitions (/internal/models) that are being passed into the ValidateStruct function (used in controller).

---

## Testing

To run all tests use below command.

```
go test ./...
```

This will run all files that match the testing file naming convention (\*\_test.go).

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

## Role based access control (RBAC) settings

The authorization settings are found in the ./internal/auth/defaultPolicy.go file.
This data structure is used by the setupcasbin policy to implement policy in DB upon server start.

SetupCasbinPolicy functions in a way where it adds policies only if they're not found already.

Format of policy: Subject, Object, Action (ie. "Who" is accessing "DB object" to commit "CRUD action")
