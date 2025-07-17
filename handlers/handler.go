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

// WeatherHandler is responsible for handling weather-related requests and managing dependencies.
type WeatherHandler struct {
	LocationService      services.LocationService     // Service to retrieve location data
	WeatherService       services.WeatherService      // Service to retrieve weather data
	CepValidator         *shared.CepValidator         // Validator for validating CEP (Brazilian ZIP code)
	TemperatureConverter *shared.TemperatureConverter // Utility to convert temperatures between Celsius, Fahrenheit, and Kelvin
}

// NewWeatherHandler creates and returns a new WeatherHandler with everything initialized
// Nova instância do WeatherHandler é criada e retornada com todos os serviços e utilitários inicializados
func NewWeatherHandler(
	locationService services.LocationService,
	weatherService services.WeatherService,
	temperatureConverter *shared.TemperatureConverter,
	chBrasilAPI, chViaCEP chan models.Location,
) *WeatherHandler {
	// Initialize channels for fetching location data
	// Inicializa os canais para buscar dados de localização

	return &WeatherHandler{
		LocationService:      locationService,                   // Assign location service
		WeatherService:       weatherService,                    // Assign weather service
		CepValidator:         shared.NewCepValidator(`^\d{8}$`), // Assign CEP validator with a regex pattern
		TemperatureConverter: temperatureConverter,              // Assign temperature converter utility
	}
}

// WeatherHandlerFunc handles the HTTP requests for weather data
// Função que lida com as requisições HTTP para obter dados meteorológicos
func (h *WeatherHandler) WeatherHandlerFunc() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Retrieve the 'cep' query parameter from the URL
		// Obtém o parâmetro 'cep' da URL da requisição
		cep := strings.TrimSpace(r.URL.Query().Get("cep"))

		// Create channels for receiving location data from APIs
		// Cria canais para receber dados de localização das APIs
		chBrasilAPI := make(chan models.Location)
		chViaCEP := make(chan models.Location)

		// Validate the CEP input
		// Valida o CEP fornecido
		if !h.CepValidator.IsValidCep(cep) {
			// Respond with an error message if CEP is invalid
			// Retorna uma resposta de erro caso o CEP seja inválido
			response := models.ErrorResponse{
				Error: "invalid zipcode", // Error message in English
			}
			// Set the HTTP status code to 422 (Unprocessable Entity)
			// Define o código de status HTTP como 422 (Entidade não processável)
			w.WriteHeader(http.StatusUnprocessableEntity)

			// Encode the response into JSON and send it to the client
			// Codifica a resposta em JSON e envia para o cliente
			json.NewEncoder(w).Encode(response)
			return
		}

		// Fetch location data based on CEP, using channels to simulate multiple API responses
		// Busca dados de localização com base no CEP, utilizando canais para simular múltiplas respostas de APIs
		location, err := h.LocationService.GetLocationFromCEP(cep, chBrasilAPI, chViaCEP)
		if err != nil || location.City == nil {
			// Respond with an error message if the location cannot be found
			// Retorna uma resposta de erro caso não seja possível encontrar a localização
			response := models.ErrorResponse{
				Error: "can not find zipcode", // Error message in English
			}
			// Set the HTTP status code to 422 (Unprocessable Entity)
			// Define o código de status HTTP como 422 (Entidade não processável)
			w.WriteHeader(http.StatusNotFound)

			// Encode the response into JSON and send it to the client
			// Codifica a resposta em JSON e envia para o cliente
			json.NewEncoder(w).Encode(response)
			return
		}

		// Fetch temperature for the city
		// Busca a temperatura para a cidade
		tempC, err := h.WeatherService.GetTemperature(*location.City)
		if err != nil {
			// Respond with an error message if fetching the temperature fails
			// Retorna uma resposta de erro caso a busca pela temperatura falhe
			response := models.ErrorResponse{
				Error: "failed to get temperature", // Error message in English
			}
			// Set the HTTP status code to 422 (Unprocessable Entity)
			// Define o código de status HTTP como 500 (Erro interno do servidor)
			w.WriteHeader(http.StatusInternalServerError)

			// Encode the response into JSON and send it to the client
			// Codifica a resposta em JSON e envia para o cliente
			json.NewEncoder(w).Encode(response)
			return
		}

		// Convert temperature using the shared utility
		// Converte a temperatura utilizando a ferramenta compartilhada
		tempF := h.TemperatureConverter.CelsiusToFahrenheit(tempC)
		tempK := h.TemperatureConverter.CelsiusToKelvin(tempC)

		// Prepare the response with temperature data in Celsius, Fahrenheit, and Kelvin
		// Prepara a resposta com os dados de temperatura em Celsius, Fahrenheit e Kelvin
		response := models.TemperatureResponse{
			Celsius:    tempC, // Temperature in Celsius
			Fahrenheit: tempF, // Temperature in Fahrenheit
			Kelvin:     tempK, // Temperature in Kelvin
		}

		// Send the response as JSON
		// Envia a resposta como JSON
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	}
}
