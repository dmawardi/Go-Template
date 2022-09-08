# Go-Template

## Written in GoLang

This repository will serve as a base web application in Go.

- Built in Go version 1.19
- Uses the [chi](https://github.com/go-chi/chi/v5) router

## To run

```
go run ./cmd
```

## To run using Docker

```
<!-- Builds docker image -->
docker build -t container-name .


<!-- runs docker image and matches port -->
docker run --publish 8080:8080 container-name
```
