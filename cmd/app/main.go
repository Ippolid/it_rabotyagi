package main

import (
	"log"
)

func main() {
	if err := app.Run(); err != nil {
		log.Fatal("Failed to start application:", err)
	}
}
