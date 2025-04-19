package config_test

import (
	"WeatherProject/src/config"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

func TestInit(t *testing.T) {

	db, mock, err := sqlmock.New()
	assert.NoError(t, err)

	defer db.Close()

	mock.ExpectPing()

	config.InitDatabase(config.DBConnection)
	err = mock.ExpectationsWereMet()
	assert.NoError(t, err, "nem todas as expectativas do mock foram atendidas")

}
