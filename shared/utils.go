package shared

import (
	"regexp"
)

// TemperatureConverter provides conversion methods for temperature.
type TemperatureConverter struct{}

// CelsiusToFahrenheit converts Celsius to Fahrenheit
func (tc *TemperatureConverter) CelsiusToFahrenheit(c float64) float64 {
	return c*1.8 + 32
}

// CelsiusToKelvin converts Celsius to Kelvin
func (tc *TemperatureConverter) CelsiusToKelvin(c float64) float64 {
	return c + 273
}

// CepValidator validates if a CEP is in the correct format.
type CepValidator struct {
	RegexPattern string
}

// NewCepValidator creates a new CepValidator with a given regex pattern.
func NewCepValidator(pattern string) *CepValidator {
	return &CepValidator{RegexPattern: pattern}
}

// IsValidCep checks if the provided CEP matches the regex pattern.
func (cv *CepValidator) IsValidCep(cep string) bool {
	match, _ := regexp.MatchString(cv.RegexPattern, cep)
	return match
}
