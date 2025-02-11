package config

import (
	"log"
	"os"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
	"github.com/joho/godotenv"
)

type Config struct {
	Env       string    `json:"env" env-default:"local"`
	Servers   Servers   `json:"servers"`
	PrimaryDB PrimaryDB `json:"database"`
}
type Servers struct {
	HTTPServer HTTPServer `json:"http"`
}
type PrimaryDB struct {
	Host        string `json:"host"`
	Port        string `json:"port"`
	User        string `json:"user"`
	Dbname      string `json:"dbname"`
	MaxAttempts int    `json:"max_attempts"`
}
type HTTPServer struct {
	Port        string        `json:"port"`
	Host        string        `json:"host"`
	Timeout     time.Duration `json:"timeout"`
	IdleTimeout time.Duration `json:"idle_Timeout"`
}

func MustLoadConfig() *Config {
	if err := godotenv.Load(".env"); err != nil {
		log.Fatal("CONFIG_PATH IS NOT SET")
		return nil
	}
	configpath := os.Getenv("CONFIG_PATH")
	if configpath == "" {
		log.Fatal("CONFIG_PATH IS NOT SET")
	}
	if _, err := os.Stat(configpath); os.IsNotExist(err) {
		log.Fatalf("config file does not exist: %s", configpath)
	}
	var cfg Config
	if err := cleanenv.ReadConfig(configpath, &cfg); err != nil {
		log.Fatalf("config file does not exist: %s", configpath)
	}
	return &cfg
}
