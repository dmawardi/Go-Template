# Go-Template

## Written in GoLang

This repository will serve as a base web application in Go.

- Built in Go version 1.19
- Uses the [chi](https://github.com/go-chi/chi/v5) router

## To run Go server

```
go run ./cmd
```

## API documentation

API documentation is auto generated using markdown within code. This is achieved using Swag.
You must navigate to folder with main.go to generate.

The below commands must be used upon making changes to the API in order to regenerate the API docs.

- "-d" directory flag allows custom directory to be used
- "--pd" flag parses dependecies as well

```
swag init -d ./cmd --pd
```

## To update Ent models

This will update Ent models and functions for ORM

```
go generate ./ent
```

## To run using Docker

```
<!-- Builds docker image -->
docker build -t container-name .


<!-- runs docker image and matches port -->
docker run --publish 8080:8080 container-name
```
