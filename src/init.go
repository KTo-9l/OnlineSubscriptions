package main

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

type config struct {
	Server struct {
		Port string `yaml:"port"`
	} `yaml:"server"`
	DB struct {
		DriverName       string `yaml:"driverName"`
		ConnectionString string `yaml:"connectionString"`
	} `yaml:"db"`
	Logs struct {
		Path string `yaml:"path"`
	} `yaml:"logs"`
}

const configPath = "config.env"

var (
	db  *sql.DB
	cfg config
)

func loadConfig() (err error) {
	err = godotenv.Load(configPath)
	if err != nil {
		fmt.Printf("[Config] Loading .env file error: %s", err.Error())
		return
	}

	cfg.Server.Port = fmt.Sprintf(":%s", os.Getenv("SERVER_PORT"))
	if cfg.Server.Port == ":" {
		return errors.New("empty server port was received from config.env")
	}

	cfg.DB.DriverName = os.Getenv("DB_DRIVER_NAME")
	cfg.DB.ConnectionString = fmt.Sprintf("postgresql://%s:%s@%s:%s/%s?sslmode=%s",
		os.Getenv("DB_USER"), os.Getenv("DB_PASSWORD"), os.Getenv("DB_HOST"), os.Getenv("DB_PORT"), os.Getenv("DB_NAME"), os.Getenv("DB_SSL"))

	cfg.Logs.Path = os.Getenv("LOGS_PATH")

	return
}

func initLog() {
	f, err := os.OpenFile(cfg.Logs.Path, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0600)
	if err != nil {
		panic(err)
	}
	log.SetOutput(f)
}

func initDB() (err error) {
	db, err = sql.Open(cfg.DB.DriverName, cfg.DB.ConnectionString)
	if err != nil {
		return err
	}
	return db.Ping()
}
