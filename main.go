package main

import (
	"log"
	"net/http"
	"os"

	handlers "post-graduation-exercise-cloud-run-weather-api/handlers"
	"post-graduation-exercise-cloud-run-weather-api/models"
	"post-graduation-exercise-cloud-run-weather-api/services"
	"post-graduation-exercise-cloud-run-weather-api/shared"

	"github.com/joho/godotenv"
)

func getHandler() *handlers.WeatherHandler {
	chBrasilAPI := make(chan models.Location)
	chViaCEP := make(chan models.Location)
	client := &http.Client{}
	temperatureConverter := &shared.TemperatureConverter{}
	apiClient := &services.APIClientImpl{Client: client}
	weatherService := services.NewWeatherService(apiClient)
	// Initialize LocationService (depends on WeatherService)
	locationService := services.NewLocationService(weatherService)
	handler := handlers.NewWeatherHandler(
		locationService,
		weatherService,
		temperatureConverter,
		chBrasilAPI,
		chViaCEP,
	)
	return handler
}

func main() {
	// Carrega as variáveis de ambiente do .env
	// Load env vars from .env file
	godotenv.Load()

	weatherHandler := getHandler()

	http.HandleFunc("/weather", weatherHandler.WeatherHandlerFunc())

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	log.Printf("Server running on port %s", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
