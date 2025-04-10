// Package config provides configuration management for various services
// such as database, grpc server, http server, etc.
// It reads configuration from file and make available to other parts of application.
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

// AppConfig holds the configurations for the entire application, including
// db, web server, and grpc server configurations.
// It has no Exported fields to encapsulate configurations.
type AppConfig struct {
	db     db
	web    web
	grpc   grpc
	logger logger
	auth   auth
}

func (app AppConfig) DB() db {
	return app.db
}

func (app AppConfig) Web() web {
	return app.web
}

func (app AppConfig) GRPC() grpc {
	return app.grpc
}

func (app AppConfig) Logger() logger {
	return app.logger
}

func (app AppConfig) Auth() auth {
	return app.auth
}

// tmpConfig holds the configurations for the entire application, including
// db, web server, and grpc server configurations.
// It should have Exported fields to work with tags.
type tmpConfig struct {
	DB struct {
		Driver   string `default:"sqlite" json:"driver" yaml:"driver" toml:"driver"`
		IP       string `default:"127.0.0.1" json:"ip" yaml:"ip" toml:"ip"`
		Port     uint   `default:"3306" json:"port" yaml:"port" toml:"port"`
		UserName string `default:"amir" json:"userName" yaml:"userName" toml:"userName"`
		Password string `default:"mirzaei" json:"password" yaml:"password" toml:"password"`
		Name     string `default:"clean-architect" json:"name" yaml:"name" toml:"name"`
		Path     string `default:"." json:"path" yaml:"path" toml:"path"`
	} `json:"db" yaml:"db" toml:"db"`
	Web struct {
		BindingIPAddress       string `default:"0.0.0.0" json:"bindingIpAddress" yaml:"bindingIpAddress" toml:"bindingIpAddress"`
		Port                   uint   `default:"8071" json:"port" yaml:"port" toml:"port"`
		ReadTimeOutInSec       uint   `default:"7" json:"readTimeOutInSec" yaml:"readTimeOutInSec" toml:"readTimeOutInSec"`
		IdleTimeoutInSec       uint   `default:"10" json:"idleTimeoutInSec" yaml:"idleTimeoutInSec" toml:"idleTimeoutInSec"`
		WriteTimeoutInSec      uint   `default:"20" json:"writeTimeoutInSec" yaml:"writeTimeoutInSec" toml:"writeTimeoutInSec"`
		ReadHeaderTimeoutInSec uint   `default:"1" json:"readHeaderTimeoutInSec" yaml:"readHeaderTimeoutInSec" toml:"readHeaderTimeoutInSec"`
		ShutdownTimeoutInSec   uint   `default:"1" json:"shutdownTimeoutInSec" yaml:"shutdownTimeoutInSec" toml:"shutdownTimeoutInSec"`
	} `json:"web" yaml:"web" toml:"web"`
	GRPC struct {
		BindingIPAddress     string `default:"127.0.0.1" json:"bindingIpAddress" yaml:"bindingIpAddress" toml:"bindingIpAddress"`
		Port                 uint   `default:"8070" json:"port" yaml:"port" toml:"port"`
		MaxReceiveMsgSize    int    `default:"5120" json:"maxReceiveMsgSize" yaml:"maxReceiveMsgSize" toml:"maxReceiveMsgSize"`
		ReadBufferSize       int    `default:"5120" json:"readBufferSize" yaml:"readBufferSize" toml:"readBufferSize"`
		HasReflection        bool   `default:"true" json:"hasReflection" yaml:"hasReflection" toml:"hasReflection"`
		ShutdownTimeoutInSec uint   `default:"1" json:"shutdownTimeoutInSec" yaml:"shutdownTimeoutInSec" toml:"shutdownTimeoutInSec"`
	} `json:"grpc" yaml:"grpc" toml:"grpc"`
	Logger struct {
		Level            int    `default:"0" json:"level" yaml:"level" toml:"level"`
		Directory        string `default:"log" json:"directory" yaml:"directory" toml:"directory"`
		FileCreationMode int    `default:"0" json:"fileCreationMode" yaml:"fileCreationMode" toml:"fileCreationMode"`
		RemoteURL        string `default:"" json:"remoteURL" yaml:"remoteURL" toml:"remoteURL"`
		Console          bool   `default:"true" json:"console" yaml:"console" toml:"console"`
	} `json:"logger" yaml:"logger" toml:"logger"`
	Auth struct {
		Secret   string `default:"some_secret" json:"secret" yaml:"secret" toml:"secret"`
		LifeTime int    `default:"1" json:"lifeTime" yaml:"lifeTime" toml:"lifeTime"`
	} `json:"auth" yaml:"auth" toml:"auth"`
}

func (cfg tmpConfig) ToAppConfig() AppConfig {
	return AppConfig{
		db: db{
			driver:   cfg.DB.Driver,
			ip:       cfg.DB.IP,
			port:     cfg.DB.Port,
			userName: cfg.DB.UserName,
			password: cfg.DB.Password,
			name:     cfg.DB.Name,
			path:     cfg.DB.Path,
		},
		web: web{
			bindingIpAddress:       cfg.Web.BindingIPAddress,
			port:                   cfg.Web.Port,
			readTimeOutInSec:       cfg.Web.ReadTimeOutInSec,
			idleTimeoutInSec:       cfg.Web.IdleTimeoutInSec,
			writeTimeoutInSec:      cfg.Web.WriteTimeoutInSec,
			readHeaderTimeoutInSec: cfg.Web.ReadHeaderTimeoutInSec,
			shutdownTimeout:        cfg.Web.ShutdownTimeoutInSec,
		},
		grpc: grpc{
			bindingIpAddress:  cfg.GRPC.BindingIPAddress,
			port:              cfg.GRPC.Port,
			maxReceiveMsgSize: cfg.GRPC.MaxReceiveMsgSize,
			readBufferSize:    cfg.GRPC.ReadBufferSize,
			hasReflection:     cfg.GRPC.HasReflection,
			shutdownTimeout:   cfg.GRPC.ShutdownTimeoutInSec,
		},
		logger: logger{
			level:            cfg.Logger.Level,
			directory:        cfg.Logger.Directory,
			fileCreationMode: cfg.Logger.FileCreationMode,
			remoteURL:        cfg.Logger.RemoteURL,
			console:          cfg.Logger.Console,
		},
		auth: auth{
			secret:   cfg.Auth.Secret,
			lifeTime: cfg.Auth.LifeTime,
		},
	}
}

// LoadConfig will return AppConfig which that values are filled by given config file's address.
func LoadConfig(fileAddress string) (AppConfig, error) {
	bytes, err := os.ReadFile(fileAddress)
	if err != nil {
		return AppConfig{}, err
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
		return AppConfig{}, err
	}

	return tmpCfg.ToAppConfig(), nil
}

// LoadConfigOrDefault will do LoadConfig. if loading had problem, then returns default values config.
func LoadConfigOrDefault(fileAddress string) (AppConfig, error) {
	cfg, err := LoadConfig(fileAddress)
	if err == nil {
		return cfg, nil
	}

	var tmpCfg tmpConfig
	err = defaults.Set(&tmpCfg)
	return tmpCfg.ToAppConfig(), err
}
