package main

import (
	"log"
	"os"

	// Blank-import the function package so the init() runs
	_ "geps.dev/progress"
	"github.com/GoogleCloudPlatform/functions-framework-go/funcframework"
)

func main() {
	// Set FUNCTION_TARGET if not already set
	if os.Getenv("FUNCTION_TARGET") == "" {
		os.Setenv("FUNCTION_TARGET", "Progress")
	}
	
	// Use PORT environment variable, or default to 8080.
	port := "8080"
	if envPort := os.Getenv("PORT"); envPort != "" {
		port = envPort
	}
	if err := funcframework.Start(port); err != nil {
		log.Fatalf("funcframework.Start: %v\n", err)
	}
}
