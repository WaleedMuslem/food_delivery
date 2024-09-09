package config

import (
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	Port                   string
	AccessSecret           string
	AccessLifetimeminutes  int
	RefreshSecret          string
	RefreshLifetimeminutes int
	DbUsername             string
	DbPassword             string
	DbName                 string
}

func NewConfig() *Config {

	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	accessMin, err := strconv.Atoi(os.Getenv("ACCESS_LIFTIME_MINUTES"))
	if err != nil {
		log.Fatal("Error Passing ACCESS_LIFTIME_MINUTES")
	}

	refreshMin, err := strconv.Atoi(os.Getenv("REFRESH_LIFTIME_MINUTES"))
	if err != nil {
		log.Fatal("Error Passing ACCESS_LIFTIME_MINUTES")
	}

	return &Config{
		Port:                   os.Getenv("PORT"),
		AccessSecret:           os.Getenv("ACCESS_SECERT"),
		AccessLifetimeminutes:  accessMin,
		RefreshSecret:          os.Getenv("REFRESH_SECERT"),
		RefreshLifetimeminutes: refreshMin,
		DbUsername:             os.Getenv("DB_USERNAME"),
		DbPassword:             os.Getenv("DB_PASSWORD"),
		DbName:                 os.Getenv("DB_NAME"),
	}
}
