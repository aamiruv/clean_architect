package main

import (
	"flag"
	"github.com/AmirMirzayi/clean_architecture/pkg/config"
)

var configPath string

func init() {
	configPath = *flag.String("config", "config.json", "config file path, eg: -config=/path/to/file.json")

	cfg := config.LoadConfig(configPath)
	_ = cfg
}
