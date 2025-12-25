package main

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"

	"ziggy/cmd"
)

func main() {
	// Load .env file if it exists (ignore error if not found)
	_ = godotenv.Load()

	if err := cmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
