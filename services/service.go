package services

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"post-graduation-exercise-cloud-run-weather-api/models"
	"time"
)

// APIClient defines the behavior of an external API client.
type APIClient interface {
	Get(url string) (*http.Response, error)
}

// LocationService is an interface that defines the methods to interact with location services.
type LocationService interface {
	GetLocationFromCEP(cep string, chBrasilAPI, chViaCEP chan models.Location) (models.Location, error)
}

// WeatherService is an interface that defines the methods for interacting with weather services.
type WeatherService interface {
	GetTemperature(city string) (float64, error)
	GetClient() APIClient // Add this method to the interface
}

// WeatherServiceImpl is the concrete implementation of the WeatherService interface.
type WeatherServiceImpl struct {
	Client APIClient
}

// APIClientImpl is the concrete implementation of the APIClient interface.
type APIClientImpl struct {
	Client *http.Client
}

// LocationServiceImpl is the concrete implementation of the LocationService interface.
type LocationServiceImpl struct {
	WeatherService WeatherService
}

// NewWeatherService creates and returns a new instance of WeatherServiceImpl.
func NewWeatherService(client APIClient) WeatherService {
	return &WeatherServiceImpl{
		Client: client,
	}
}

// NewLocationService creates and returns a new LocationServiceImpl instance.
func NewLocationService(weatherService WeatherService) LocationService {
	return &LocationServiceImpl{
		WeatherService: weatherService,
	}
}

// NewAPIClient creates and returns a new instance of APIClientImpl.
func NewAPIClient(client *http.Client) *APIClientImpl {
	return &APIClientImpl{
		Client: client,
	}
}

// GetTemperature retrieves the current temperature for a given city.
func (ws *WeatherServiceImpl) GetTemperature(city string) (float64, error) {
	apiKey := os.Getenv("WEATHER_API_KEY")
	// Fix spaces on names
	encodedCity := url.QueryEscape(city)
	url := fmt.Sprintf("https://api.weatherapi.com/v1/current.json?key=%s&q=%s", apiKey, encodedCity)

	resp, err := ws.Client.Get(url)
	if err != nil {
		fmt.Printf("error fetching weather %v", err)
		return 0, err
	}
	defer resp.Body.Close()

	var weather models.WeatherResponse
	if err := json.NewDecoder(resp.Body).Decode(&weather); err != nil {
		fmt.Printf("error decoding weather %v", err)
		return 0, err
	}

	return weather.Current.TempC, nil
}

// GetClient returns the APIClient used in WeatherServiceImpl.
func (ws *WeatherServiceImpl) GetClient() APIClient {
	return ws.Client
}

// GetLocationFromCEP retrieves location data based on a given CEP.
func (ls *LocationServiceImpl) GetLocationFromCEP(cep string, chBrasilAPI, chViaCEP chan models.Location) (models.Location, error) {
	timeout := time.After(10 * time.Second)

	// Asynchronously fetch data from the APIs
	go ls.fetchFromBrasilAPI(cep, chBrasilAPI)
	go ls.fetchFromViaCEP(cep, chViaCEP)

	select {
	case res := <-chBrasilAPI:
		if res.Localidade != nil {
			return res, nil
		}
	case res := <-chViaCEP:
		if res.Localidade != nil {
			return res, nil
		}
	case <-timeout:
		return models.Location{}, errors.New("timeout after 10 seconds")
	}

	return models.Location{}, errors.New("error searching for CEP data")
}

// fetchFromBrasilAPI fetches location data from the BrasilAPI.
func (ls *LocationServiceImpl) fetchFromBrasilAPI(cep string, ch chan models.Location) {
	url := fmt.Sprintf("https://brasilapi.com.br/api/cep/v1/%s", cep)
	resp, err := ls.WeatherService.GetClient().Get(url) // Use GetClient to avoid casting
	if err != nil || resp.StatusCode != http.StatusOK {
		fmt.Printf("error fetching brasil api %v", err)
		ch <- models.Location{}
		return
	}
	defer resp.Body.Close()

	var address models.BrasilAPIResponse
	if err := json.NewDecoder(resp.Body).Decode(&address); err != nil {
		fmt.Printf("error decoding brasil api %v", err)
		ch <- models.Location{}
		return
	}

	ch <- models.Location{
		Cep:        &cep,
		Localidade: &address.Neighborhood,
		Uf:         &address.State,
		City:       &address.City,
	}
}

// fetchFromViaCEP fetches location data from the ViaCEP API.
func (ls *LocationServiceImpl) fetchFromViaCEP(cep string, ch chan models.Location) {
	url := fmt.Sprintf("http://viacep.com.br/ws/%s/json", cep)
	resp, err := ls.WeatherService.GetClient().Get(url) // Use GetClient to avoid casting
	if err != nil || resp.StatusCode != http.StatusOK {
		fmt.Printf("error fetching viacep api %v", err)
		ch <- models.Location{}
		return
	}
	defer resp.Body.Close()

	var address models.ViaCEPResponse
	if err := json.NewDecoder(resp.Body).Decode(&address); err != nil {
		fmt.Printf("error decoding viacep api %v", err)
		ch <- models.Location{}
		return
	}

	ch <- models.Location{
		Cep:        &cep,
		Localidade: &address.Localidade,
		Uf:         &address.UF,
		City:       &address.Localidade,
	}
}

// Get performs an HTTP GET request.
func (api *APIClientImpl) Get(url string) (*http.Response, error) {
	return api.Client.Get(url)
}
