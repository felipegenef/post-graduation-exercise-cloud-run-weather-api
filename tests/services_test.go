package tests

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"os"
	"post-graduation-exercise-cloud-run-weather-api/models"
	"post-graduation-exercise-cloud-run-weather-api/services"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockWeatherService com retorno fixo
type MockWeatherService struct {
	mock.Mock
}

func (m *MockWeatherService) GetTemperature(city string) (float64, error) {
	args := m.Called(city)
	return args.Get(0).(float64), args.Error(1)
}

func (m *MockWeatherService) GetClient() services.APIClient {
	args := m.Called()
	return args.Get(0).(services.APIClient)
}

// MockApiClient implementando o método Get
type MockApiClient struct {
	mock.Mock
}

// Mock do método Get com switch para URL
func (m *MockApiClient) Get(url string) (*http.Response, error) {
	args := m.Called(url)
	return args.Get(0).(*http.Response), args.Error(1)
}

// MockLocationService simula a obtenção de localização
type MockLocationService struct {
	mock.Mock
}

// Simula a obtenção de localização a partir do CEP
func (m *MockLocationService) GetLocationFromCEP(cep string, chBrasilAPI, chViaCEP chan models.Location) (models.Location, error) {
	args := m.Called(cep, chBrasilAPI, chViaCEP)
	return args.Get(0).(models.Location), args.Error(1)
}

func TestRaceFetch(t *testing.T) {
	mockLocationService := new(MockLocationService)

	// Canais
	chBrasilAPI := make(chan models.Location)
	chViaCEP := make(chan models.Location)

	cep := "12345678"
	localidadeBrasilAPI := "São Paulo"
	ufBrasilAPI := "SP"
	locationBrasilAPI := models.Location{
		Cep:        &cep,
		Localidade: &localidadeBrasilAPI,
		Uf:         &ufBrasilAPI,
	}

	localidadeViaCEP := "Rio de Janeiro"
	ufViaCEP := "RJ"
	locationViaCEP := models.Location{
		Cep:        &cep,
		Localidade: &localidadeViaCEP,
		Uf:         &ufViaCEP,
	}

	// Expect behavior on MockLocationService
	mockLocationService.On("GetLocationFromCEP", cep, chBrasilAPI, chViaCEP).Return(locationBrasilAPI, nil)

	// Testando a resposta com chBrasilAPI primeiro
	go func() {
		chBrasilAPI <- locationBrasilAPI
	}()

	// Testando a resposta com chViaCEP primeiro
	go func() {
		time.Sleep(1 * time.Second)
		chViaCEP <- locationViaCEP
	}()

	// Esperado o resultado
	resultado, err := mockLocationService.GetLocationFromCEP(cep, chBrasilAPI, chViaCEP)
	assert.NoError(t, err)
	assert.Equal(t, locationBrasilAPI, resultado)

	// Verificando chamadas
	mockLocationService.AssertExpectations(t)
}

func TestHttpFetch(t *testing.T) {
	mockApiClient := new(MockApiClient)
	apiKey := os.Getenv("WEATHER_API_KEY")
	locationService := services.NewWeatherService(mockApiClient)

	// Mock do retorno do método Get
	mockApiClient.On("Get", fmt.Sprintf("https://api.weatherapi.com/v1/current.json?key=%v&q=SP", apiKey)).
		Return(&http.Response{
			StatusCode: 200,
			Body:       io.NopCloser(bytes.NewReader([]byte(`{"current": {"temp_c":13.14}}`))),
		}, nil).Once()

	mockApiClient.On("Get", mock.Anything).Return(&http.Response{
		StatusCode: 200,
		Body:       io.NopCloser(bytes.NewReader([]byte(`{"current": {"temp_c":13.0}}`))),
	}, nil).Once()

	// Teste para a URL "SP"
	response, err := locationService.GetTemperature("SP")
	assert.NoError(t, err)
	assert.Equal(t, 13.14, response)

	// Teste para a URL "other-city" (outro valor)
	response, err = locationService.GetTemperature("other-city")
	assert.NoError(t, err)
	assert.Equal(t, 13.0, response)

	// Verificando as expectativas dos mocks
	mockApiClient.AssertExpectations(t)
}
