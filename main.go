package main

import (
	"log"
	"net/http"
	"os"

	handlers "post-graduation-exercise-cloud-run-weather-api/handlers"

	"github.com/joho/godotenv"
)

func main() {
	// Carrega as variáveis de ambiente do .env
	godotenv.Load()

	http.HandleFunc("/weather", handlers.WeatherHandler)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	log.Printf("Server running on port %s", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
