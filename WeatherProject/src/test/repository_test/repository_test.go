package repository_test

import (
	"WeatherProject/src/config"
	"WeatherProject/src/models"
	db "WeatherProject/src/repository"
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/patrickmn/go-cache"
	"github.com/stretchr/testify/assert"
)

func TestFindByDate(t *testing.T) {

	dbMock, mock, err := sqlmock.New()
	assert.NoError(t, err)

	defer dbMock.Close()

	repo := &db.MySqlWeatherRepository{DB: dbMock}
	config.CachingInMemory = cache.New(5*time.Minute, 10*time.Minute)

	expected := models.WeatherResponse{
		Country: "BR",
		Date:    "2024-04-10",
		Text:    "Sol com nuvens",
	}

	rows := sqlmock.NewRows([]string{"country", "date", "message"}).
		AddRow(expected.Country, expected.Date, expected.Text)

	mock.ExpectQuery(regexp.QuoteMeta("SELECT country, date,message FROM weather WHERE date = ? LIMIT 1")).
		WithArgs("2024-04-10").
		WillReturnRows(rows)

	result, err := repo.FindByDate("2024-04-10")

	if err != nil {
		t.Fatalf("Erro inesperado: %v", err)
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Houve falha nas expectativas do mock: %v", err)
	}
	assert.NoError(t, err)
	assert.Equal(t, &expected, result)

}

func TestInsertWeather(t *testing.T) {

	dbMock, mock, err := sqlmock.New()

	assert.NoError(t, err)

	defer dbMock.Close()

	repo := &db.MySqlWeatherRepository{DB: dbMock}

	expected := models.WeatherResponse{
		Country: "BR",
		Date:    "2024-04-10",
		Text:    "Sol com nuvens",
	}

	mock.ExpectPrepare(regexp.QuoteMeta("INSERT INTO weather (country, date, message) VALUES(?,?,?)")).
		ExpectExec().
		WithArgs(expected.Country, expected.Date, expected.Text).
		WillReturnResult(sqlmock.NewResult(1, 1))

	err = repo.Insert(expected)

	assert.NoError(t, err)

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestFindAll(t *testing.T) {
	dbMock, mock, err := sqlmock.New()
	assert.NoError(t, err)
	var expectedList []models.WeatherResponse

	defer dbMock.Close()

	repo := &db.MySqlWeatherRepository{DB: dbMock}
	config.CachingInMemory = cache.New(5*time.Minute, 10*time.Minute)

	expected := models.WeatherResponse{
		Country: "BR",
		Date:    "2024-04-10",
	}

	expectedList = append(expectedList, expected)

	rows := sqlmock.NewRows([]string{"country", "date"}).
		AddRow(expected.Country, expected.Date)

	mock.ExpectQuery(regexp.QuoteMeta("SELECT country, date FROM weather")).
		WillReturnRows(rows)

	result, err := repo.FindAll()

	if err != nil {
		t.Fatalf("Erro inesperado: %v", err)
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Houve falha nas expectativas do mock: %v", err)
	}
	assert.NoError(t, err)
	assert.Equal(t, expectedList, result)
}
