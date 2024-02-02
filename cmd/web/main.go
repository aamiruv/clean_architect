package main

import (
	"context"
	"flag"
	"github.com/AmirMirzayi/clean_architecture/pkg/config"
	"github.com/AmirMirzayi/clean_architecture/pkg/logger/file"
	"github.com/AmirMirzayi/clean_architecture/pkg/web"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
)

var configPath string

func init() {
	flag.StringVar(&configPath, "config", "config.json", "config file path, eg: -config=/path/to/file.json")
	flag.Parse()
}

func main() {
	_, cancel := context.WithCancel(context.Background())
	defer cancel()

	logger := file.NewLogger(file.LogDaily, "log")
	log.SetOutput(logger)
	log.SetFlags(log.Ltime | log.Lshortfile | log.LUTC)

	webLoggerFile := file.NewLogger(file.LogHourly, "weblog")
	webLogger := log.New(
		webLoggerFile, "",
		log.Ltime|log.Lshortfile|log.LUTC|log.Lmsgprefix,
	)
	cfg := config.LoadConfig(configPath)

	webServer := web.NewServer(
		cfg.GetWeb().GetAddress(),
		webLogger,
		1<<11,
		time.Duration(cfg.GetWeb().IdleTimeoutInSec)*time.Second,
		time.Duration(cfg.GetWeb().ReadTimeOutInSec)*time.Second,
		time.Duration(cfg.GetWeb().WriteTimeoutInSec)*time.Second,
		time.Duration(cfg.GetWeb().ReadHeaderTimeoutInSec)*time.Second,
	)
	log.Printf("initialize web server in address: %s", cfg.GetWeb().GetAddress())
	go func() {
		log.Panic(webServer.Run())
	}()

	sigint := make(chan os.Signal, 1)
	signal.Notify(sigint, os.Interrupt, syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP, syscall.SIGQUIT, syscall.SIGKILL)
	<-sigint
	if err := webServer.GracefulShutdown(5 * time.Second); err != nil {
		log.Println(err)
	}
}
