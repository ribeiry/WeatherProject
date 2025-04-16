package db

import "WeatherProject/models"

type WeatherRepository interface {
	Insert(w models.WeatherResponse) error
	FindByDate(date string) (*models.WeatherResponse, error)
	FindAll() ([]models.WeatherResponse, error)
}
