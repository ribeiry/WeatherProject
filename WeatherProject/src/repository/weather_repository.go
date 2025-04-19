package db

import (
	"WeatherProject/src/config"
	"WeatherProject/src/models"
	"database/sql"
	"log"

	"github.com/patrickmn/go-cache"
)

type MySqlWeatherRepository struct {
	DB *sql.DB
}

func (r *MySqlWeatherRepository) Insert(weather models.WeatherResponse) error {

	stmtIns, err := r.DB.Prepare("INSERT INTO weather (country, date, message) VALUES(?,?,?)")

	if err != nil {
		return err
	}

	defer stmtIns.Close()

	_, err = stmtIns.Exec(weather.Country, weather.Date, weather.Text)
	if err != nil {
		return err
	}
	return err
}

func (r *MySqlWeatherRepository) FindByDate(date string) (*models.WeatherResponse, error) {

	var w models.WeatherResponse

	if cached, found := config.CachingInMemory.Get(date); found {
		if w, ok := cached.(*models.WeatherResponse); ok {
			return w, nil
		}

	}

	query := "SELECT country, date,message FROM weather WHERE date = ? LIMIT 1"

	rows := r.DB.QueryRow(query, date)

	err := rows.Scan(&w.Country, &w.Date, &w.Text)

	if err != nil {
		if err == sql.ErrNoRows {
			log.Println("Nenhum resultado encontrado.")
			return nil, nil
		} else {
			return nil, err
		}
	} else {

		//Setando o Cache
		config.CachingInMemory.Set(date, &w, cache.DefaultExpiration)
		return &w, nil
	}
}

func (r *MySqlWeatherRepository) FindAll() ([]models.WeatherResponse, error) {

	rows, err := r.DB.Query("SELECT country, date FROM weather")

	if err != nil {
		log.Println(err.Error())
	}

	var list []models.WeatherResponse
	for rows.Next() {

		var w models.WeatherResponse
		err = rows.Scan(&w.Country, &w.Date)

		if err != nil {
			log.Println(err)
		}
		list = append(list, w)
	}
	return list, nil
}
