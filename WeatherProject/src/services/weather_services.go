package service

import (
	"WeatherProject/config"
	"WeatherProject/models"
	db "WeatherProject/repository"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"sync"
	"time"
)

func WeatherServiceConstructor(repo db.WeatherRepository,
	baseUrl string, apiKey string,
	semaphore chan struct{}) *WeatherService {

	return &WeatherService{Repo: repo, BaseURL: baseUrl, APIKey: apiKey, Semaphore: semaphore}
}

type WeatherService struct {
	Repo        db.WeatherRepository
	BaseURL     string
	APIKey      string
	Client      *http.Client
	Semaphore   chan struct{}
	SucessCount int
	FailCount   int
	Mutex       sync.Mutex
}

func (repository *WeatherService) CreateWeatherEntry(locale, msg string) error {
	w := models.WeatherResponse{
		Country: locale,
		Date:    time.Now().Format("2006-01-02"),
		Text:    msg,
	}

	return repository.Repo.Insert(w)
}

func (repository *WeatherService) GetTemperatureDay(date string) (*models.WeatherResponse, error) {

	resp := repository.findByDateAsync(date)

	repository.Mutex.Lock()
	defer repository.Mutex.Unlock()

	if resp.Err != nil {
		repository.FailCount++
		return nil, resp.Err
	}

	repository.SucessCount++

	return &models.WeatherResponse{
		Country: resp.DataResponse.Country,
		Date:    resp.DataResponse.Date,
		Text:    resp.DataResponse.Text,
	}, nil
}

func (repository *WeatherService) findByDateAsync(date string) models.WeatherFindResult {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)

	defer cancel()

	repository.Semaphore <- struct{}{}
	defer func() {
		<-repository.Semaphore
	}()

	findRespWeather := make(chan models.WeatherFindResult, 1)

	go func() {
		data, err := repository.Repo.FindByDate(date)
		findRespWeather <- models.WeatherFindResult{DataResponse: data, Err: err}
	}()

	select {
	case resp := <-findRespWeather:
		if resp.Err != nil {
			repository.FailCount++
			log.Println("Error ao buscar no banco: ", resp.Err)
			return models.WeatherFindResult{
				Err: fmt.Errorf("Erro ao buscar ao banco de dados"),
			}
		} else {
			repository.SucessCount++
			log.Println("Busca efetuada com sucesso")
			return resp
		}
	case <-ctx.Done():
		repository.FailCount++
		log.Println(" timeout para esperar o Banco de dados")
		return models.WeatherFindResult{
			Err: fmt.Errorf("Erro ao buscar ao banco de dados"),
		}
	}

}

func (repository *WeatherService) GetTodayWeather() (*models.WeatherResponse, error) {
	today := time.Now().Format("2006-01-02")
	weather, err := repository.Repo.FindByDate(today)

	if err != nil || weather != nil {
		log.Println("Encontrado no banco de dados")
		return weather, nil
	}
	log.Println("Not found in database, call api ...")

	return repository.fetchAndCacheWeather()
}

func (repository *WeatherService) fetchAndCacheWeather() (*models.WeatherResponse, error) {

	weather, err := repository.GetTemperature()

	if err != nil {
		log.Println("Error ao buscar na API")
		return nil, err
	}

	go repository.saveWeatherAsync(weather.Country, weather.Text)

	return weather, nil
}

func (repository *WeatherService) saveWeatherAsync(country, text string) {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	repository.Semaphore <- struct{}{}
	defer func() {
		<-repository.Semaphore
	}()

	saveErrChan := make(chan error, 1)

	go func() {
		saveErrChan <- repository.CreateWeatherEntry(country, text)

	}()

	select {
	case err := <-saveErrChan:
		if err != nil {
			repository.FailCount++
			log.Println("Error ao salvar no banco: ", err)
		} else {
			repository.SucessCount++
			log.Println("Salvamento feito com sucesso")
		}
	case <-ctx.Done():
		repository.Mutex.Lock()
		defer repository.Mutex.Unlock()

		repository.FailCount++
		log.Println(" timeout para esperar o Banco de dados")
	}
}

func (repo *WeatherService) GetAllWeather() ([]models.WeatherResponse, error) {
	return repo.Repo.FindAll()
}

func (repo *WeatherService) GetTemperature() (*models.WeatherResponse, error) {

	client := repo.Client
	url := fmt.Sprintf("%s%s?token=%s", repo.BaseURL, config.SynopticPath, repo.APIKey)

	if client == nil {
		client = http.DefaultClient
	}
	resp, err := client.Get(url)

	if err != nil {
		log.Println(err)
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)

	if err != nil {
		log.Println(err)
		return nil, err

	}
	var weather []models.WeatherResponse
	if err := json.Unmarshal(body, &weather); err != nil {
		log.Println(err)
		return nil, err

	}
	if len(weather) == 0 {
		log.Println(err)

		return nil, err
	}
	return &weather[0], nil
}
