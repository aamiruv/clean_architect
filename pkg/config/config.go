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
	tmpCfg := struct {
		DB        db     `json:"db"`
		Web       web    `json:"web"`
		Grpc      grpc   `json:"grpc"`
		LoggerURL string `json:"loggerURL"`
	}{}
	if err = json.Unmarshal(bytes, &tmpCfg); err != nil {
		panic(err)
	}
	return appConfig{
		db:        tmpCfg.DB,
		web:       tmpCfg.Web,
		grpc:      tmpCfg.Grpc,
		loggerURL: tmpCfg.LoggerURL,
	}
}

func (app appConfig) DB() db {
	return app.db
}

func (app appConfig) Web() web {
	return app.web
}

func (app appConfig) Grpc() grpc {
	return app.grpc
}

func (app appConfig) LoggerURL() string {
	return app.loggerURL
}
