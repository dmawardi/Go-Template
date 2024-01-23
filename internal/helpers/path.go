package helpers

import (
	"log"
	"os"
	"strings"
)

// Build path from working directory
func BuildPathFromWorkingDirectory(urlFromWD string) string {
	// generate path
	dirPath, err := os.Getwd()
	if err != nil {
		log.Fatal("Could not get working directory")
	}

	// Split path to remove excess path when running tests
	splitPath := strings.Split(dirPath, "internal")

	// Grab initial part of path and join with path from project root directory
	urlPath := splitPath[0] + urlFromWD
	return urlPath
}
