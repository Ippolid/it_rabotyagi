package main

import (
	"itpath/internal/app"
	"log"
)

func main() {
	if err := app.Run(); err != nil {
		log.Fatal("Failed to start application:", err)
	}
}
