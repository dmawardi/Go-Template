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
# Database Settings
DB_USER=postgres
DB_PASS=
DB_HOST=localhost
DB_PORT=5432
DB_NAME=
SESSIONS_SECRET_KEY=
HMAC_SECRET=
# SMTP Settings
SMTP_HOST=
SMTP_PORT=
SMTP_USERNAME=
SMTP_PASSWORD=
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

---

## How to use

Follow these steps to add a feature to the API. This template uses the clean architecture pattern

1. Build schema and auto migrate in ./internal/db using ORM instructions below. (Be sure to build the receiver functions getID() and ObtainValue() for the admin panel)
2. Build repository in ./internal/repository which is the interaction between the DB and the application. This should use a struct with receiver functions.
3. Build incoming DTO models for Create and Update JSON requests in ./internal/models
4. Build service in ./internal/service that uses the repository and applies business logic.
5. Build the controller (handler) in ./internal/controller that accepts the request, performs data validation, then sends to the service to interact with database.
6. Add validation to handler using govalidator. This is done by adding `valid:""` key-value pairs to struct DTO definitions (/internal/models) that are being passed into the ValidateStruct function (used in controllers).
7. Add the new controller to the API struct in the ./internal/routes/routes.go file. This allows it to be used within the Routes function in the same file. Build routes to use the handlers that have been created in step 4 using the api struct.
8. Update the ApiSetup function in the ./cmd/main.go file to build the new repository, service, and controller.
9. Add the route to the RBAC authorization policy file (./internal/auth/rbac_policy.go)
   ADMIN
10. Add a new file in ./internal/admin-panel with the title of the schema being built. Inside will need to contain:
11. The db schema will need to fit the specs of the db.AdminPanelSchema, so add two receiver functions to your schema struct.

- The first will be a funciton that returns the ID of the schema (GetId())
- The second will be a function that returns the value of the field given a key (ObtainValue())
  These functions will be used in the admin panel to display the data.

11. Add the routes in the AddAdminRoutes function in ./internal/routes/routes.go
12. Add the route to the RBAC authorization policy file (./internal/auth/rbac_policy.go)

13. (Testing) For e2e testing, you will need to update the controllers_test.go file in ./internal/controller. Updates are required in the testDbRepo struct, buildAPI, setupDatabase & setupDBAuthAppModels functions

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

The below commands must be used upon making changes to the API in order to regenerate the API docs.

- "-d" directory flag allows custom directory to be used
- "-g" flag allows direct pointing to the main.go file for generation of swagger annotations from files that are imported (ie. controllers, services, repositories, etc.)
- "--pd" flag parses dependecies as well
- "--parseInternal" flag parses internal packages

```
<!-- Generate docs from home folder -->
<!-- Remove old API docs folder in static -->
<!-- To move the generated folder to be accessible to users -->
swag init -d ./internal/controller -g ../../cmd/main.go --pd --parseInternal
rm -rf static/docs
mv docs static/
```

This will update API documentation generated in the ./docs folder. It is served on path /swagger

## To use Database ORM

To edit schemas: Go to ./internal/db/schemas.go

The schemas are Structs based off of gorm.Model.

After creating the schema in schemas.go, go to db.go and add to automigrate.

For the admin panel, you will need to add two receiver functions to your schema struct in order for it to adhere to the db.AdminPanelSchema interface.

Note for creating and updating using GORM: Relationship data that does not yet exist will be created as a new entry. However, if you try to edit an existing record, it will not allow you to.

## Role based access control (RBAC) settings

The authorization model is found in the ./internal/auth/rbac_model.conf file.
This data structure is used by the setupcasbin policy to implement policy in DB upon server start.

The default policy is found in the ./internal/auth/rbac_policy.csv file.

Format of policy: Subject, Object, Action ("Who" is accessing "DB object" to commit "CRUD action")

Eg. admin, /api/v1/users, POST

In the policy implementation above:
p = Used to assign permissions to roles
eg. Assigning read permission to user role for /api/me endpoint
| p type | v0 | v1 | v2 |
| ------ | ---- | ------- | ---- |
| p | user | /api/me | read |

g = Used to assign roles to users
eg. Assigning a moderator role to user with id 2
| p type | v0 | v1 |
| ------ | ---- | ------- |
| g | 2 | moderator |

g2 = Used to assign roles to roles to create an inheritance heirarchy
eg. Allocating all permissions for moderator role to admin role
| p type | v0 | v1 |
| ------ | ---- | ------- |
| g2 | admin | moderator |

What about record level control?
eg. User can only edit their own profile

For sake of flexibility, reducing casbin policy model complexity, and to avoid having to create a new policy for each new record, we will use custom application logic within handlers to check if the user is the owner of the record to allow passing of the request to the service.

Request -> Authorization (does user have permission through role?) -> Validation (does user own record?) -> Service (perform CRUD operation)

## To run using Docker

To run the application within a Docker container, you will need to build the image and run the container.

When running Docker on a Mac with an ARM processor, you will need to use the buildx command to build the image for amd64. This is where the --platform option comes in handy.

"container-name" is typically the github address of your project. (ie. dmawardi/go-template)

```
<!-- Builds docker image -->
docker build -t container-name .
<!-- Builds Docker image for amd64 (if on arm64) -->
docker buildx build --platform linux/amd64 -t container-name .


<!-- runs docker image and matches port -->
docker run --publish 8080:8080 container-name
```

In order to run the Docker image on a server, you will need to push the image to a Docker registry (Docker Hub). This can be done using Docker Desktop

## To deploy container on a server

First, ensure that Docker is installed on the server.

Then, pull the Docker image on the server using the container name (dmawardi/go-template:latest) and the docker pull command:

```

docker pull container-name:version

```

Then, run the Docker image on the server using the following command:

```

docker run -d -p 8080:8080 container-name:version

```
