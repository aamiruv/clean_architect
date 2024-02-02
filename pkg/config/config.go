package config

import (
	"encoding/json"
	"os"
)

type appConfig struct {
	db        db
	web       web
	loggerURL string
}

func LoadConfig(fileAddress string) appConfig {
	bytes, err := os.ReadFile(fileAddress)
	if err != nil {
		panic(err)
	}
	type tmpConfig struct {
		DB        db     `json:"db"`
		Web       web    `json:"web"`
		LoggerURL string `json:"loggerURL"`
	}
	tmpCfg := tmpConfig{}
	if err = json.Unmarshal(bytes, &tmpCfg); err != nil {
		panic(err)
	}
	return appConfig{
		db:        tmpCfg.DB,
		web:       tmpCfg.Web,
		loggerURL: tmpCfg.LoggerURL,
	}
}

func (cfg appConfig) GetDB() db {
	return cfg.db
}

func (app appConfig) GetWeb() web {
	return app.web
}

func (app appConfig) GetLoggerURL() string {
	return app.loggerURL
}
