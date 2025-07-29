package main

import (
	"GetCurrency/internal/app"
	"log"
)

func main() {
	if err := app.Run(); err != nil {
		log.Fatalf("failed to run application: %v", err)
	}
}
