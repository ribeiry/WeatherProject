package handler

import (
	"WeatherProject/src/config"
	repository "WeatherProject/src/repository"
	service "WeatherProject/src/services"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

func Run() {

	//Inicia a configuracao do Banco
	error := config.InitDatabase(config.DBConnection)

	if error != nil {
		log.Println("Error to connection database %v", error)
	}
	//Inicia a configuracao do Cache
	config.InitCache()

	service := SetUpService()

	router := ConfigRoutes(service)

	// Roda o servidor
	router.Run("localhost:8080")
	defer config.MySqlDB.Close()

}

func ConfigRoutes(service *service.WeatherService) *gin.Engine {
	//Cria a camada do GIN e suas rotas
	router := gin.Default()
	router.GET("/temperatura", IndentedJSON(service))

	return router
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

func SetUpService() *service.WeatherService {

	repo := &repository.MySqlWeatherRepository{DB: config.MySqlDB}
	return &service.WeatherService{
		Repo:    repo,
		BaseURL: config.BaseURL,
		APIKey:  config.APIKey,
	}
}
