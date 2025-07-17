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
	// Load env vars from .env file
	godotenv.Load()

	weatherHandler := handlers.NewWeatherHandler()

	http.HandleFunc("/weather", weatherHandler.WeatherHandlerFunc())

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	log.Printf("Server running on port %s", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
