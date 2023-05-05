package config

import (
	"os"
	"strings"

	"github.com/joho/godotenv"
)

var Default = map[string]Config{
	"APP_URL":     "localhost",
	"APP_ENV":     "local",
	"APP_PORT":    "3000",
	"DB_DRIVER":   "mysql",
	"DB_HOST":     "127.0.0.1",
	"DB_PORT":     "3306",
	"DB_NAME":     "crud-golang",
	"DB_USER":     "root",
	"DB_PASSWORD": "",
}

type Config string

func init() {
	godotenv.Load()
}

func Get(key string) Config {
	value := Config(os.Getenv(key))
	if value == "" {
		value = Default[key]
	}
	return value
}

func (c Config) String() string {
	return string(c)
}

func (c Config) Bool() bool {
	if strings.ToLower(c.String()) == "true" {
		return true
	}
	return false
}
