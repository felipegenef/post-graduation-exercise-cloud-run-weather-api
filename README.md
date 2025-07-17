# Weather API - Golang (Eng version below)

![Test Coverage](https://codecov.io/gh/felipegenef/post-graduation-exercise-cloud-run-weather-api/branch/main/graph/badge.svg)
![Test Status](https://github.com/felipegenef/post-graduation-exercise-cloud-run-weather-api/actions/workflows/go.yaml/badge.svg)

## Descrição

Este projeto é uma API desenvolvida como parte de um exercício de pós-graduação em Golang. A API consulta duas fontes externas para obter informações de localização a partir de um **CEP (Código de Endereçamento Postal)** e retorna a temperatura atual da cidade correspondente. A consulta é feita simultaneamente usando duas APIs externas: **BrasilAPI** e **ViaCEP**, e, após validar o CEP, a temperatura é obtida de uma API de clima. A temperatura é convertida para **Celsius**, **Fahrenheit** e **Kelvin**.

### Requisitos do exercício

- O sistema deve receber um **CEP válido de 8 dígitos**.
- O sistema deve realizar a pesquisa do CEP e encontrar o nome da localização, retornando as temperaturas nas escalas: Celsius, Fahrenheit e Kelvin.
- O sistema deve responder adequadamente nos seguintes cenários:
  - **Em caso de sucesso:**
    - Código HTTP: `200`
    - Response Body: 
    ```json
    { "temp_C": 28.5, "temp_F": 83.3, "temp_K": 301.65 }
    ```
  - **Em caso de falha (CEP inválido):**
    - Código HTTP: `422`
    - Mensagem: `"invalid zipcode"`
  - **Em caso de falha (CEP não encontrado):**
    - Código HTTP: `404`
    - Mensagem: `"can not find zipcode"`

- O sistema deve ser deployado no **Google Cloud Run**.

## Funcionalidades

- Consulta a localização a partir de um CEP utilizando as APIs **BrasilAPI** e **ViaCEP**.
- Validação do formato do CEP antes de realizar a consulta.
- Consulta à temperatura atual da cidade usando uma API externa de clima.
- Conversão da temperatura para **Celsius**, **Fahrenheit** e **Kelvin**.
- Resposta estruturada em formato **JSON** com a temperatura nas três escalas.
- Tratamento de erros para respostas inválidas ou falhas de API.

## Requisitos

- Go 1.23.3 ou superior
- Variáveis de ambiente configuradas, incluindo a chave de API para o serviço de clima.
- Acesso à internet para consumir as APIs externas de localização e clima.

## Como Acessar a API

A API está hospedada no **Google Cloud Run** e pode ser acessada através do endpoint:

```bash
curl https://weather-api-76fmx4exrq-uc.a.run.app/weather?cep=01025020
```

# Weather API - Golang (Versão em Português acima)

![Test Coverage](https://codecov.io/gh/felipegenef/post-graduation-exercise-cloud-run-weather-api/branch/main/graph/badge.svg)
![Test Status](https://github.com/felipegenef/post-graduation-exercise-cloud-run-weather-api/actions/workflows/go.yaml/badge.svg)

## Description

This project is an API developed as part of a postgraduate exercise in Golang. The API queries two external sources to obtain location information based on a **ZIP code (CEP)** and returns the current temperature of the corresponding city. The request is made simultaneously using two external APIs: **BrasilAPI** and **ViaCEP**, and after validating the ZIP code, the temperature is fetched from a weather API. The temperature is converted into **Celsius**, **Fahrenheit**, and **Kelvin**.

### Exercise Requirements

- The system must accept a **valid 8-digit ZIP code**.
- The system must query the ZIP code and find the location's name, then return the temperatures formatted as: Celsius, Fahrenheit, and Kelvin.
- The system should respond appropriately in the following scenarios:
  - **On success:**
    - HTTP Status Code: `200`
    - Response Body: 
    ```json
    { "temp_C": 28.5, "temp_F": 83.3, "temp_K": 301.65 }
    ```
  - **On failure (invalid ZIP code):**
    - HTTP Status Code: `422`
    - Message: `"invalid zipcode"`
  - **On failure (ZIP code not found):**
    - HTTP Status Code: `404`
    - Message: `"can not find zipcode"`

- The system must be deployed on **Google Cloud Run**.

## Features

- Queries location from a ZIP code using **BrasilAPI** and **ViaCEP**.
- Validates the ZIP code format before making the request.
- Fetches the current temperature of the city using an external weather API.
- Converts the temperature to **Celsius**, **Fahrenheit**, and **Kelvin**.
- Responds with a structured **JSON** response containing the temperature in the three scales.
- Error handling for invalid responses or API failures.

## Requirements

- Go 1.23.3 or higher
- Environment variables configured, including the API key for the weather service.
- Internet access to consume the external location and weather APIs.

## How to Access the API

The API is hosted on **Google Cloud Run** and can be accessed via the endpoint:

```bash
curl https://weather-api-76fmx4exrq-uc.a.run.app/weather?cep=01025020
```