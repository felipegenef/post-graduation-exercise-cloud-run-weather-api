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
}

// NewWeatherHandler creates and returns a new WeatherHandler with everything initialized
func NewWeatherHandler(
	locationService services.LocationService,
	weatherService services.WeatherService,
	temperatureConverter *shared.TemperatureConverter,
	chBrasilAPI, chViaCEP chan models.Location,

) *WeatherHandler {
	// Initialize channels for fetching location data

	return &WeatherHandler{
		LocationService:      locationService,
		WeatherService:       weatherService,
		CepValidator:         shared.NewCepValidator(`^\d{8}$`),
		TemperatureConverter: temperatureConverter,
	}
}

// WeatherHandlerFunc handles the HTTP requests for weather data
func (h *WeatherHandler) WeatherHandlerFunc() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		cep := strings.TrimSpace(r.URL.Query().Get("cep"))
		chBrasilAPI := make(chan models.Location)
		chViaCEP := make(chan models.Location)
		// Validate the CEP input
		if !h.CepValidator.IsValidCep(cep) {
			response := models.ErrorResponse{
				Error: "invalid zipcode",
			}
			// Set the HTTP status code to 422 (Unprocessable Entity)
			w.WriteHeader(http.StatusUnprocessableEntity)

			// Encode the response into JSON and send it to the client
			json.NewEncoder(w).Encode(response)
			return
		}

		// Fetch location data based on CEP, using channels to simulate multiple API responses
		location, err := h.LocationService.GetLocationFromCEP(cep, chBrasilAPI, chViaCEP)
		if err != nil || location.City == nil {
			response := models.ErrorResponse{
				Error: "can not find zipcode",
			}
			// Set the HTTP status code to 422 (Unprocessable Entity)
			w.WriteHeader(http.StatusNotFound)

			// Encode the response into JSON and send it to the client
			json.NewEncoder(w).Encode(response)
			return
		}

		// Fetch temperature for the city
		tempC, err := h.WeatherService.GetTemperature(*location.City)
		if err != nil {
			response := models.ErrorResponse{
				Error: "failed to get temperature",
			}
			// Set the HTTP status code to 422 (Unprocessable Entity)
			w.WriteHeader(http.StatusInternalServerError)

			// Encode the response into JSON and send it to the client
			json.NewEncoder(w).Encode(response)
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
		json.NewEncoder(w).Encode(response)
	}
}
