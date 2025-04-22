package config

import (
	"database/sql"
	"log"

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
		log.Println("Erro connectionDB")
		return err
	}
	return nil
}
