package config

import (
	"encoding/json"
	"os"
)

type appConfig struct {
	db        db
	web       web
	grpc      grpc
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
		Grpc      grpc   `json:"grpc"`
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

func (app appConfig) GetDB() db {
	return app.db
}

func (app appConfig) GetWeb() web {
	return app.web
}

func (app appConfig) GetGrpc() grpc {
	return app.grpc
}

func (app appConfig) GetLoggerURL() string {
	return app.loggerURL
}
