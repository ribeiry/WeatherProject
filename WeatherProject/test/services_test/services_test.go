package services_test

import (
	"WeatherProject/config"
	"WeatherProject/models"
	service "WeatherProject/services"
	"testing"
	"time"

	mocksOA "WeatherProject/test/mocks"

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
