package service

import "WeatherProject/models"

type WeatherServiceInterface interface {
	CreateWeatherEntry(locale, msg string) error
	GetTodayWeather() (*models.WeatherResponse, error)
	GetAllWeather() ([]models.WeatherResponse, error)
	GetTemperature() (*models.WeatherResponse, error)
	GetTemperatureDay(date string) (*models.WeatherResponse, error)
}
