package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	models "post-graduation-exercise-cloud-run-weather-api/models"
	service "post-graduation-exercise-cloud-run-weather-api/services"
	utils "post-graduation-exercise-cloud-run-weather-api/shared"
	"strings"
)

func WeatherHandler(w http.ResponseWriter, r *http.Request) {
	cep := strings.TrimSpace(r.URL.Query().Get("cep"))

	if !utils.IsValidCep(cep) {
		http.Error(w, "invalid zipcode", http.StatusUnprocessableEntity)
		return
	}

	location, err := service.GetLocationFromCEP(cep)
	if err != nil || location.City == nil {
		fmt.Println("Error looking for Zip Code: %v", err)
		http.Error(w, "can not find zipcode", http.StatusNotFound)
		return
	}

	tempC, err := service.GetTemperature(*location.City)
	if err != nil {
		fmt.Println("Error looking for Temperature: %v", err)
		http.Error(w, "failed to get temperature", http.StatusInternalServerError)
		return
	}

	tempF := utils.CelsiusToFahrenheit(tempC)
	tempK := utils.CelsiusToKelvin(tempC)

	response := models.TemperatureResponse{
		Celsius:    tempC,
		Fahrenheit: tempF,
		Kelvin:     tempK,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
