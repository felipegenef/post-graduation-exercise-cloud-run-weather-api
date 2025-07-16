package services

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"os"
	models "post-graduation-exercise-cloud-run-weather-api/models"
	"time"
)

func GetLocationFromCEP(cep string) (models.Location, error) {
	// Canais para as respostas das duas APIs
	// Channels to receive responses from both APIs
	chBrasilAPI := make(chan models.Location)
	chViaCEP := make(chan models.Location)

	// Timeout para a resposta (1 segundo)
	// Timeout for the response (1 second)
	timeout := time.After(10 * time.Second)

	// Inicia as goroutines para as duas APIs
	// Start the goroutines for both APIs
	go fetchFromBrasilAPI(cep, chBrasilAPI)
	go fetchFromViaCEP(cep, chViaCEP)

	// Usando select para esperar a resposta mais rápida
	// Using select to wait for the fastest response
	select {
	case res := <-chBrasilAPI:
		if res.Localidade != nil {
			fmt.Println("Found Data from BrasilAPI first")
			return res, nil
		}
	case res := <-chViaCEP:
		if res.Localidade != nil {
			fmt.Println("Found Data from ViaCEP first")
			return res, nil
		}
	case <-timeout:
		// Caso o tempo de resposta tenha excedido 10 segundos
		// If the response time exceeds 10 seconds
		fmt.Println("Erro: Timeout. Nenhuma resposta recebida dentro do tempo limite.")
		return models.Location{}, errors.New("timeout após 10 segundos")
	}
	return models.Location{}, errors.New("erro searching for cep data")
}

func GetTemperature(city string) (float64, error) {
	apiKey := os.Getenv("WEATHER_API_KEY")
	// Fix spaces on names
	encodedCity := url.QueryEscape(city)
	url := fmt.Sprintf("https://api.weatherapi.com/v1/current.json?key=%s&q=%s", apiKey, encodedCity)

	resp, err := http.Get(url)
	if err != nil {
		fmt.Println("error fetching weather API: %v", err)
		return 0, err
	}
	defer resp.Body.Close()

	var weather models.WeatherResponse
	if err := json.NewDecoder(resp.Body).Decode(&weather); err != nil {
		fmt.Println("error decoding weather API: %v", err)
		return 0, err
	}

	return weather.Current.TempC, nil
}

// Função para buscar os dados usando a API BrasilAPI
// Function to fetch data using the BrasilAPI
func fetchFromBrasilAPI(cep string, ch chan models.Location) {
	url := fmt.Sprintf("https://brasilapi.com.br/api/cep/v1/%s", cep)
	resp, err := http.Get(url)
	if err != nil {
		// Caso haja erro na requisição, envia resposta com erro
		// If there is an error in the request, send an error response
		fmt.Println("Error from BrasilAPI : %v", err)
		ch <- models.Location{}
		return
	}
	defer resp.Body.Close()

	// Verifica se a resposta foi OK (status 200)
	// Check if the response status is OK (status 200)
	if resp.StatusCode != http.StatusOK {
		fmt.Println("Error from BrasilAPI : %v", err)
		ch <- models.Location{}
		return
	}

	var address models.BrasilAPIResponse
	if err := json.NewDecoder(resp.Body).Decode(&address); err != nil {
		// Se falhar ao decodificar os dados, envia resposta com erro
		// If decoding the data fails, send an error response
		fmt.Println("Error from BrasilAPI : %v", err)
		ch <- models.Location{}
		return
	}

	// Envia os dados recebidos pela API
	// Send the data received from the API
	ch <- models.Location{Cep: &cep, Localidade: &address.Neighborhood, Uf: &address.State, City: &address.City}
}

// Função para buscar os dados usando a API ViaCEP
// Function to fetch data using the ViaCEP API
func fetchFromViaCEP(cep string, ch chan models.Location) {
	url := fmt.Sprintf("http://viacep.com.br/ws/%s/json", cep)
	resp, err := http.Get(url)
	if err != nil {
		// Caso haja erro na requisição, envia resposta com erro
		// If there is an error in the request, send an error response
		fmt.Println("Error from ViaCepAPI : %v", err)
		ch <- models.Location{}
		return
	}
	defer resp.Body.Close()

	// Verifica se a resposta foi OK (status 200)
	// Check if the response status is OK (status 200)
	if resp.StatusCode != http.StatusOK {
		fmt.Println("Error from ViaCepAPI : %v", err)
		ch <- models.Location{}
		return
	}

	var address models.ViaCEPResponse
	if err := json.NewDecoder(resp.Body).Decode(&address); err != nil {
		// Se falhar ao decodificar os dados, envia resposta com erro
		// If decoding the data fails, send an error response
		fmt.Println("Error from decoding ViaCepAPI response: %v", err)
		ch <- models.Location{}
		return
	}

	// Envia os dados recebidos pela API
	// Send the data received from the API
	ch <- models.Location{Cep: &cep, Localidade: &address.Localidade, Uf: &address.UF, City: &address.Localidade}
}
