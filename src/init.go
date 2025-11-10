package main

import (
	"database/sql"
	"log"
	"os"

	"gopkg.in/yaml.v3"

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

const configPath = "config.yaml"

var (
	db  *sql.DB
	cfg config
)

func loadConfig() (err error) {
	cfgBytes, err := os.ReadFile(configPath)
	if err != nil {
		log.Printf("[Config] Reading config file error: %s", err.Error())
		return
	}

	err = yaml.Unmarshal(cfgBytes, &cfg)
	if err != nil {
		log.Printf("[Config] Unmarshalling config error: %s", err.Error())
		return
	}
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
