package db

import "WeatherProject/src/models"

type WeatherRepository interface {
	Insert(w models.WeatherResponse) error
	FindByDate(date string) (*models.WeatherResponse, error)
	FindAll() ([]models.WeatherResponse, error)
}
