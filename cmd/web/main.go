package main

import (
	"flag"
	"github.com/AmirMirzayi/clean_architecture/pkg/config"
	"log"
	"net/http"
	"os"
)

var configPath string

func init() {
	flag.StringVar(&configPath, "config", "config.json", "config file path, eg: -config=/path/to/file.json")
	flag.Parse()
}

func main() {
	f, err := os.OpenFile("log.txt", os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		panic(err)
	}
	log.SetOutput(f)
	cfg := config.LoadConfig(configPath)
	_ = cfg
	mux := http.NewServeMux()
	srv := &http.Server{
		Handler: mux,
		Addr:    cfg.GetWeb().GetAddress(),
	}
	log.Fatalln(srv.ListenAndServe())
}
