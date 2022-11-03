package config

import (
	"github.com/ilyakaznacheev/cleanenv"
	"log"
	"sync"
)

type Config struct {
	SenderEmail    string `env:"sender_email" env-required:"true"`
	SenderPass     string `env:"sender_pass" env-required:"true"`
	RecipientEmail string `env:"recipient_email" env-required:"true"`
	Host           string `env:"host" env-required:"true"`
	Port           int    `env:"port" env-required:"true"`
	DataPath       string `env:"data_path" env-required:"true"`
}

const StructDateTimeFormat = "2006-01-02 15:04"
const StructDateFormat = "2006-01-02"

var instance *Config
var once sync.Once

func GetConfig(path string) *Config {
	once.Do(func() {
		log.Printf("read application config from path %s", path)

		instance = &Config{}

		if err := cleanenv.ReadConfig(path, instance); err != nil {
			log.Fatal(err)
		}
	})
	return instance
}
