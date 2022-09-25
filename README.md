# Go-Template

## Written in GoLang

This repository will serve as a base web application in Go.

- Built in Go version 1.19
- Uses the [chi](https://github.com/go-chi/chi/v5) router
- Uses [godotenv](https://github.com/joho/godotenv) for environmental variables
- Uses [Swaggo](https://github.com/swaggo/swag) to generate API documentation
- Uses [Go-Validator](https://github.com/asaskevich/govalidator) for validating incoming data

### Database (Object Relational Management)

- Uses [Ent](https://entgo.io/ent) for ORM (Postgres)

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

## Testing

To run tests use below command.

```
go test
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

## To update Database models

This will update Ent models and functions for ORM. Restarting the server will enact the changes on the DB.

```
go generate ./ent
```

---

## Role based access control (RBAC) settings

The authorization settings are found in the ./internal/auth/defaultPolicy.go file.
This data structure is used by the setupcasbin policy to implement policy in DB upon server start.

SetupCasbinPolicy functions in a way where it adds policies only if they're not found already.

Format of policy: Subject, Object, Action (ie. "Who" is accessing "DB object" to commit "CRUD action")
