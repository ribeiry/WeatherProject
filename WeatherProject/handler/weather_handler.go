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

	repo := &repository.MySqlWeatherRepository{DB: config.DB}
	service := &service.WeatherService{
		Repo:    repo,
		BaseURL: config.BaseURL,
		APIKey:  config.APIKey,
	}

	//Cria a camada do GIN e suas rotas
	router := gin.Default()
	router.GET("/temperatura", IndentedJSON(service))

	router.Run("localhost:8080")
	defer config.DB.Close()

}

func IndentedJSON(service service.WeatherServiceInterface) gin.HandlerFunc {

	return func(c *gin.Context) {
		response, err := service.GetTodayWeather()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao consultar temperatura"})
			return
		}
		c.IndentedJSON(http.StatusOK, response)
	}
}
