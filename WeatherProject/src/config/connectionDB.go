package config

import (
	"database/sql"

	_ "github.com/go-sql-driver/mysql"
)

var MySqlDB *sql.DB

func InitDatabase(conection string) error {
	var err error
	MySqlDB, err = sql.Open("mysql", conection)

	if err != nil {
		return err
	}
	err = MySqlDB.Ping()
	if err != nil {
		return err
	}
	return nil
}
