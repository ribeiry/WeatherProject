package handler

import (
	"WeatherProject/config"
	repository "WeatherProject/repository"
	service "WeatherProject/services"
	"net/http"

	"github.com/gin-gonic/gin"
)

func Run() {

	//Inicia a configuracao do Banco
	config.Init()

	//Inicia a configuracao do Cache
	config.InitCache()

	//Cria a camada do GIN e suas rotas
	router := gin.Default()
	router.GET("/temperatura", IndentedJSON)

	router.Run("localhost:8080")
	defer config.DB.Close()

}

func IndentedJSON(c *gin.Context) {

	repo := &repository.MySqlWeatherRepository{DB: config.DB}
	service := &service.WeatherService{Repo: repo, BaseURL: config.BaseURL, APIKey: config.APIKey}

	response, err := service.GetTodayWeather()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao consultar temperatura"})
		return
	}
	c.IndentedJSON(http.StatusOK, response)
}
