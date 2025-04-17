package handler_test

import (
	"WeatherProject/handler"
	"WeatherProject/models"
	mocksOA "WeatherProject/test/mocks"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestIndentedJSON(t *testing.T) {

	mockRepo := new(mocksOA.WeatherServiceInterface)

	expected := &models.WeatherResponse{
		Country: "BR",
		Date:    time.Now().Format("2006-01-02"),
		Text:    "A Zona de Convergência Intertropical (ZCIT) está ativa e espalha nuvens carregadas pelo RN e o norte do CE. Uma nova frente fria já começa a influenciar o tempo no RS. No começo da tarde deste domingo (31), pancadas de chuva eram observadas na fronteira com o Uruguai. Até à noite, as pancadas de chuva se espalham por mais áreas do estado. Uma área de baixa pressão no Paraguai, associado a um cavado nos níveis médios da atmosfera, estimula a formação de nuvens carregadas em MS. ",
	}

	mockRepo.On("GetTodayWeather").Return(expected, nil)

	mockRepo.On("FindByDate", expected.Date).Return(expected, nil)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	handler := handler.IndentedJSON(mockRepo)
	handler(c)

	assert.Equal(t, http.StatusOK, w.Code)
}
