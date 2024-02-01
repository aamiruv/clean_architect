package main

import (
	"flag"
	"github.com/AmirMirzayi/clean_architecture/pkg/config"
	"github.com/AmirMirzayi/clean_architecture/pkg/logger/file"
	"log"
	"net/http"
)

var configPath string

func init() {
	flag.StringVar(&configPath, "config", "config.json", "config file path, eg: -config=/path/to/file.json")
	flag.Parse()
}

func main() {
	logger := file.NewLogger(file.LogDaily, "log")
	log.SetOutput(logger)
	log.SetFlags(log.Ltime | log.Lshortfile | log.LUTC)

	webLoggerFile := file.NewLogger(file.LogHourly, "weblog")
	webLogger := log.New(
		webLoggerFile, "",
		log.Ltime|log.Lshortfile|log.LUTC|log.Lmsgprefix,
	)
	cfg := config.LoadConfig(configPath)
	mux := http.NewServeMux()
	srv := &http.Server{
		Handler:  mux,
		Addr:     cfg.GetWeb().GetAddress(),
		ErrorLog: webLogger,
	}

	log.Printf("initialize web server in address: %s", cfg.GetWeb().GetAddress())
	log.Fatalln(srv.ListenAndServe())
}
