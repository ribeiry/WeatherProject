package services_test

import (
	"WeatherProject/src/config"
	"WeatherProject/src/models"
	service "WeatherProject/src/services"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	mocksOA "WeatherProject/src/test/mocks"

	"github.com/stretchr/testify/assert"
)

func TestGetTodayWeather(t *testing.T) {

	mockRepo := new(mocksOA.MockWeatherRepository)
	service := &service.WeatherService{Repo: mockRepo, BaseURL: config.BaseURL, APIKey: config.APIKey}

	expected := &models.WeatherResponse{
		Country: "BR",
		Date:    time.Now().Format("2006-01-02"),
		Text:    "A Zona de Convergência Intertropical (ZCIT) está ativa e espalha nuvens carregadas pelo RN e o norte do CE. Uma nova frente fria já começa a influenciar o tempo no RS. No começo da tarde deste domingo (31), pancadas de chuva eram observadas na fronteira com o Uruguai. Até à noite, as pancadas de chuva se espalham por mais áreas do estado. Uma área de baixa pressão no Paraguai, associado a um cavado nos níveis médios da atmosfera, estimula a formação de nuvens carregadas em MS. ",
	}

	mockRepo.On("FindByDate", expected.Date).Return(expected, nil)

	result, err := service.GetTodayWeather()

	assert.NoError(t, err)
	assert.Equal(t, expected, result)
	mockRepo.AssertExpectations(t)
}

func TestGetTemperature(t *testing.T) {

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode([]models.WeatherResponse{
			{
				Country: "BR",
				Date:    "2024-04-19",
				Text:    "Sol com Nuvens",
			},
		})
	})

	server := httptest.NewServer(handler)

	defer server.Close()

	service := &service.WeatherService{
		BaseURL: server.URL,
		APIKey:  "fake-token",
		Client:  server.Client(),
	}

	resp, err := service.GetTemperature()

	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, "BR", resp.Country)
	assert.Equal(t, "2024-04-19", resp.Date)
	assert.Equal(t, "Sol com Nuvens", resp.Text)
}

func TestGetAllWeather(t *testing.T) {

	mockRepo := new(mocksOA.MockWeatherRepository)
	var expectedList []models.WeatherResponse
	service := &service.WeatherService{
		Repo:    mockRepo,
		BaseURL: config.BaseURL,
		APIKey:  config.APIKey,
	}

	expected := models.WeatherResponse{
		Country: "BR",
		Date:    "2024-04-10",
	}

	expectedList = append(expectedList, expected)

	mockRepo.On("FindAll").Return(expectedList, nil)

	result, err := service.GetAllWeather()

	assert.NoError(t, err)
	assert.Equal(t, expectedList, result)
	mockRepo.AssertExpectations(t)
}

func TestSaveWeather(t *testing.T) {

	mockRepo := new(mocksOA.MockWeatherRepository)
	service := &service.WeatherService{
		Repo:    mockRepo,
		BaseURL: config.BaseURL,
		APIKey:  config.APIKey,
	}

	WeatherInn := models.WeatherResponse{
		Country: "BR",
		Date:    time.Now().Format("2006-01-02"),
		Text:    "Sol com Nuvens",
	}

	Country := "BR"
	Text := "Sol com Nuvens"

	mockRepo.On("Insert", WeatherInn).Return(nil, nil)

	err := service.CreateWeatherEntry(Country, Text)

	assert.NoError(t, err)
}
