package main

import (
	"fmt"
	"log"
	"net/http"
)

const portNumber = ":8080"

func main() {
	fmt.Printf("Starting application on port: %s\n", portNumber)

	srv := &http.Server{
		Addr:    portNumber,
		Handler: routes(),
	}

	// Listen and serve using server settings above
	err := srv.ListenAndServe()
	if err != nil {

		log.Fatal(err)
	}
}
