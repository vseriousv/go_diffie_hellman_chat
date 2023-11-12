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
	c.ApiUrl = "https://dh-chat.swapline.io"

	return &c
}
