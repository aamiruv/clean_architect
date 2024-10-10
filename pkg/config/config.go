package config

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"

	"github.com/BurntSushi/toml"
	"github.com/creasty/defaults"
	"gopkg.in/yaml.v3"
)

type appConfig struct {
	db        db
	web       web
	grpc      grpc
	loggerURL string
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

type tmpConfig struct {
	DB struct {
		IP       string `default:"127.0.0.1" json:"ip" yaml:"ip" toml:"ip"`
		Port     uint   `default:"3306" json:"port" yaml:"port" toml:"port"`
		UserName string `default:"amir" json:"userName" yaml:"userName" toml:"userName"`
		Password string `default:"mirzaei" json:"password" yaml:"password" toml:"password"`
		DBName   string `default:"clean-architect" json:"dbName" yaml:"dbName" toml:"dbName"`
	} `json:"db" yaml:"db" toml:"db"`
	Web struct {
		BindingIPAddress       string `default:"0.0.0.0" json:"bindingIpAddress" yaml:"bindingIpAddress" toml:"bindingIpAddress"`
		Port                   uint   `default:"8071" json:"port" yaml:"port" toml:"port"`
		ReadTimeOutInSec       uint   `default:"7" json:"readTimeOutInSec" yaml:"readTimeOutInSec" toml:"readTimeOutInSec"`
		IdleTimeoutInSec       uint   `default:"10" json:"idleTimeoutInSec" yaml:"idleTimeoutInSec" toml:"idleTimeoutInSec"`
		WriteTimeoutInSec      uint   `default:"20" json:"writeTimeoutInSec" yaml:"writeTimeoutInSec" toml:"writeTimeoutInSec"`
		ReadHeaderTimeoutInSec uint   `default:"1" json:"readHeaderTimeoutInSec" yaml:"readHeaderTimeoutInSec" toml:"readHeaderTimeoutInSec"`
	} `json:"web" yaml:"web" toml:"web"`
	GRPC struct {
		BindingIPAddress  string `default:"127.0.0.1" json:"bindingIpAddress" yaml:"bindingIpAddress" toml:"bindingIpAddress"`
		Port              uint   `default:"8070" json:"port" yaml:"port" toml:"port"`
		MaxReceiveMsgSize int    `default:"5120" json:"maxReceiveMsgSize" yaml:"maxReceiveMsgSize" toml:"maxReceiveMsgSize"`
		ReadBufferSize    int    `default:"5120" json:"readBufferSize" yaml:"readBufferSize" toml:"readBufferSize"`
		HasReflection     bool   `default:"true" json:"hasReflection" yaml:"hasReflection" toml:"hasReflection"`
	} `json:"grpc" yaml:"grpc" toml:"grpc"`
	LoggerURL string `default:"http://127.0.0.1:8080/ping" json:"loggerURL" yaml:"loggerURL" toml:"loggerURL"`
}

func (cfg tmpConfig) ToAppConfig() appConfig {
	return appConfig{
		db: db{
			ip:       cfg.DB.IP,
			port:     cfg.DB.Port,
			userName: cfg.DB.UserName,
			password: cfg.DB.Password,
			dbName:   cfg.DB.DBName,
		},
		web: web{
			bindingIpAddress:       cfg.Web.BindingIPAddress,
			port:                   cfg.Web.Port,
			readTimeOutInSec:       cfg.Web.ReadTimeOutInSec,
			idleTimeoutInSec:       cfg.Web.IdleTimeoutInSec,
			writeTimeoutInSec:      cfg.Web.WriteTimeoutInSec,
			readHeaderTimeoutInSec: cfg.Web.ReadHeaderTimeoutInSec,
		},
		grpc: grpc{
			bindingIpAddress:  cfg.GRPC.BindingIPAddress,
			port:              cfg.GRPC.Port,
			maxReceiveMsgSize: cfg.GRPC.MaxReceiveMsgSize,
			readBufferSize:    cfg.GRPC.ReadBufferSize,
			hasReflection:     cfg.GRPC.HasReflection,
		},
		loggerURL: cfg.LoggerURL,
	}
}

func LoadConfig(fileAddress string) (appConfig, error) {
	bytes, err := os.ReadFile(fileAddress)
	if err != nil {
		return appConfig{}, err
	}

	var tmpCfg tmpConfig

	switch filepath.Ext(fileAddress) {
	case ".json":
		err = json.Unmarshal(bytes, &tmpCfg)
	case ".yml", ".yaml":
		err = yaml.Unmarshal(bytes, &tmpCfg)
	case ".toml":
		err = toml.Unmarshal(bytes, &tmpCfg)
	default:
		err = errors.New("unsupported config's file type")
	}

	if err != nil {
		return appConfig{}, err
	}

	return tmpCfg.ToAppConfig(), nil
}

func LoadConfigOrDefault(fileAddress string) (appConfig, error) {
	cfg, err := LoadConfig(fileAddress)
	if err == nil {
		return cfg, nil
	}

	var tmpCfg tmpConfig
	err = defaults.Set(&tmpCfg)
	return tmpCfg.ToAppConfig(), err
}
