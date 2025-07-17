package handlers

import (
	"encoding/json"
	"net/http"
	"post-graduation-exercise-cloud-run-weather-api/models"
	"post-graduation-exercise-cloud-run-weather-api/services"
	"post-graduation-exercise-cloud-run-weather-api/shared"
	"strings"
)

// Define interfaces for services that can be injected

type WeatherHandler struct {
	LocationService      services.LocationService
	WeatherService       services.WeatherService
	CepValidator         *shared.CepValidator
	TemperatureConverter *shared.TemperatureConverter
	ChBrasilAPI          chan models.Location
	ChViaCEP             chan models.Location
}

// NewWeatherHandler creates and returns a new WeatherHandler with everything initialized
func NewWeatherHandler() *WeatherHandler {
	// Initialize channels for fetching location data
	chBrasilAPI := make(chan models.Location)
	chViaCEP := make(chan models.Location)
	client := &http.Client{}
	cepValidator := shared.NewCepValidator(`^\d{8}$`) // Example CEP regex pattern
	temperatureConverter := &shared.TemperatureConverter{}
	apiClient := &services.APIClientImpl{Client: client}
	weatherService := services.NewWeatherService(apiClient)
	// Initialize LocationService (depends on WeatherService)
	locationService := services.NewLocationService(weatherService)

	return &WeatherHandler{
		LocationService:      locationService,
		WeatherService:       weatherService,
		CepValidator:         cepValidator,
		TemperatureConverter: temperatureConverter,
		ChBrasilAPI:          chBrasilAPI,
		ChViaCEP:             chViaCEP,
	}
}

// WeatherHandlerFunc handles the HTTP requests for weather data
func (h *WeatherHandler) WeatherHandlerFunc() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		cep := strings.TrimSpace(r.URL.Query().Get("cep"))

		// Validate the CEP input
		if !h.CepValidator.IsValidCep(cep) {
			http.Error(w, "invalid zipcode", http.StatusUnprocessableEntity)
			return
		}

		// Fetch location data based on CEP, using channels to simulate multiple API responses
		location, err := h.LocationService.GetLocationFromCEP(cep, h.ChBrasilAPI, h.ChViaCEP)
		if err != nil || location.City == nil {
			http.Error(w, "cannot find zipcode", http.StatusNotFound)
			return
		}

		// Fetch temperature for the city
		tempC, err := h.WeatherService.GetTemperature(*location.City)
		if err != nil {
			http.Error(w, "failed to get temperature", http.StatusInternalServerError)
			return
		}

		// Convert temperature using the shared utility
		tempF := h.TemperatureConverter.CelsiusToFahrenheit(tempC)
		tempK := h.TemperatureConverter.CelsiusToKelvin(tempC)

		// Prepare the response
		response := models.TemperatureResponse{
			Celsius:    tempC,
			Fahrenheit: tempF,
			Kelvin:     tempK,
		}

		// Send the response as JSON
		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(response); err != nil {
			http.Error(w, "failed to encode response", http.StatusInternalServerError)
		}
	}
}
