package config

type Config struct {
	App
}

type App struct {
	GoEnv  string
	ApiUrl string
}

func DefaultConfig() *Config {
	var c Config

	c.GoEnv = "development"
	c.ApiUrl = "http://127.0.0.1:4000"

	return &c
}
