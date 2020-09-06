package config

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/lib/pq" // for including postgre driver

	"github.com/svakode/svachan/utils"
)

func InitDB() (dbConn *sql.DB) {
	var dsn string

	dbProvider := utils.FatalGetString("DATABASE_PROVIDER")
	dbHost := utils.FatalGetString("DATABASE_HOST")
	dbPort := utils.FatalGetString("DATABASE_PORT")
	dbUser := utils.FatalGetString("DATABASE_USER")
	dbPass := utils.FatalGetString("DATABASE_PASS")
	dbName := utils.FatalGetString("DATABASE_NAME")

	if dbProvider == "postgres" {
		dsn = fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", dbHost, dbPort, dbUser, dbPass, dbName)
	}

	dbConn, err := sql.Open("postgres", dsn)
	if err != nil {
		fmt.Println(err)
	}

	dbConn.SetMaxIdleConns(0)

	err = dbConn.Ping()
	if err != nil {
		fmt.Println(err.Error())
		log.Fatal("Error connecting to database")
		os.Exit(1)
	}

	return
}
