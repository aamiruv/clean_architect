package config

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"

	"github.com/BurntSushi/toml"
	"gopkg.in/yaml.v3"
)

type appConfig struct {
	db        db
	web       web
	grpc      grpc
	loggerURL string
}

type tmpConfig struct {
	DB struct {
		IP       string `json:"ip" yaml:"ip"`
		Port     uint   `json:"port" yaml:"port"`
		UserName string `json:"userName" yaml:"userName"`
		Password string `json:"password" yaml:"password"`
		DBName   string `json:"dbName" yaml:"dbName"`
	} `json:"db" yaml:"db"`
	Web struct {
		BindingIPAddress       string `json:"bindingIpAddress" yaml:"bindingIpAddress"`
		Port                   uint   `json:"port" yaml:"port"`
		ReadTimeOutInSec       uint   `json:"readTimeOutInSec" yaml:"readTimeOutInSec"`
		IdleTimeoutInSec       uint   `json:"idleTimeoutInSec" yaml:"idleTimeoutInSec"`
		WriteTimeoutInSec      uint   `json:"writeTimeoutInSec" yaml:"writeTimeoutInSec"`
		ReadHeaderTimeoutInSec uint   `json:"readHeaderTimeoutInSec" yaml:"readHeaderTimeoutInSec"`
	} `json:"web" yaml:"web"`
	GRPC struct {
		BindingIPAddress  string `json:"bindingIpAddress" yaml:"bindingIpAddress"`
		Port              uint   `json:"port" yaml:"port"`
		MaxReceiveMsgSize int    `json:"maxReceiveMsgSize" yaml:"maxReceiveMsgSize"`
		ReadBufferSize    int    `json:"readBufferSize" yaml:"readBufferSize"`
		HasReflection     bool   `json:"hasReflection" yaml:"hasReflection"`
	} `json:"grpc" yaml:"grpc"`
	LoggerURL string `json:"loggerURL" yaml:"loggerURL"`
}

func LoadConfig(fileAddress string) (appConfig, error) {
	bytes, err := os.ReadFile(fileAddress)
	if err != nil {
		return appConfig{}, err
	}

	tmpCfg := tmpConfig{}

	switch filepath.Ext(fileAddress) {
	case ".json":
		err = json.Unmarshal(bytes, &tmpCfg)
	case ".yml", ".yaml":
		err = yaml.Unmarshal(bytes, &tmpCfg)
	case ".toml":
		err = toml.Unmarshal(bytes, &tmpCfg)
	default:
		err = errors.New("Unsupported config's file type")
	}

	if err != nil {
		return appConfig{}, err
	}

	return appConfig{
		db: db{
			ip:       tmpCfg.DB.IP,
			port:     tmpCfg.DB.Port,
			userName: tmpCfg.DB.UserName,
			password: tmpCfg.DB.Password,
			dbName:   tmpCfg.DB.DBName,
		},
		web: web{
			bindingIpAddress:       tmpCfg.Web.BindingIPAddress,
			port:                   tmpCfg.Web.Port,
			readTimeOutInSec:       tmpCfg.Web.ReadTimeOutInSec,
			idleTimeoutInSec:       tmpCfg.Web.IdleTimeoutInSec,
			writeTimeoutInSec:      tmpCfg.Web.WriteTimeoutInSec,
			readHeaderTimeoutInSec: tmpCfg.Web.ReadHeaderTimeoutInSec,
		},
		grpc: grpc{
			bindingIpAddress:  tmpCfg.GRPC.BindingIPAddress,
			port:              tmpCfg.GRPC.Port,
			maxReceiveMsgSize: tmpCfg.GRPC.MaxReceiveMsgSize,
			readBufferSize:    tmpCfg.GRPC.ReadBufferSize,
			hasReflection:     tmpCfg.GRPC.HasReflection,
		},
		loggerURL: tmpCfg.LoggerURL,
	}, nil
}

func (app appConfig) DB() db {
	return app.db
}

func (app appConfig) Web() web {
	return app.web
}

func (app appConfig) GRPC() grpc {
	return app.grpc
}

func (app appConfig) LoggerURL() string {
	return app.loggerURL
}
