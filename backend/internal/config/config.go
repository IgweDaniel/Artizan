package config

import (
	"os"
	"reflect"

	"github.com/caarlos0/env"
	"github.com/joho/godotenv"
)

type Config struct {
	Port      int    `env:"PORT,required"`
	JwtSecret string `env:"JWT_SECRET,required"`
	Env       string `env:"ENVIRONMENT,required"`
	DbUrl     string `env:"DATABASE_URL,required"`
}

func LoadConfig() (Config, error) {

	cfg := Config{} // 👈 new instance of `Config`
	if _, err := os.Stat(".env"); err == nil {
		err := godotenv.Load()
		if err != nil {
			return cfg, err
		}
	}

	err := parseEnv(&cfg) // 👈 Parse environment variables into `Config`
	if err != nil {
		return cfg, err
	}

	return cfg, nil
}

func parseEnv(cfg interface{}) error {
	v := reflect.ValueOf(cfg).Elem()

	for i := 0; i < v.NumField(); i++ {
		field := v.Field(i)

		if field.Kind() == reflect.Struct {
			err := parseEnv(field.Addr().Interface())
			if err != nil {
				return err
			}
		} else {
			err := env.Parse(cfg)
			if err != nil {
				return err
			}
		}
	}

	return nil
}
