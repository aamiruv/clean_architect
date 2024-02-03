package main

import (
	"context"
	"flag"
	"github.com/AmirMirzayi/clean_architecture/api/proto/auth"
	"github.com/AmirMirzayi/clean_architecture/api/router"
	"github.com/AmirMirzayi/clean_architecture/pkg/config"
	"github.com/AmirMirzayi/clean_architecture/pkg/logger"
	"github.com/AmirMirzayi/clean_architecture/pkg/logger/file"
	weblog "github.com/AmirMirzayi/clean_architecture/pkg/logger/web"
	"github.com/AmirMirzayi/clean_architecture/pkg/server/grpc"
	"github.com/AmirMirzayi/clean_architecture/pkg/server/web"
	"golang.org/x/sync/errgroup"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
)

const ShutdownTimeout = 5 * time.Second

func main() {
	var configPath string
	flag.StringVar(&configPath, "config", "config.json", "config file path, eg: -config=/path/to/file.json")
	flag.Parse()

	_, cancel := context.WithCancel(context.Background())
	defer cancel()

	cfg := config.LoadConfig(configPath)

	fileLogger := file.NewLogger(file.LogHourly, "log")
	theWebLog := weblog.NewLogger(cfg.GetLoggerURL())
	// log on file & http url same time
	complexLogger := logger.NewComplexLogger(fileLogger, theWebLog)

	log.SetFlags(log.Ltime | log.Lshortfile | log.LUTC)
	log.SetOutput(complexLogger)

	webLoggerFile := file.NewLogger(file.LogHourly, "weblog")
	webLogger := log.New(
		webLoggerFile, "",
		log.Ltime|log.Lshortfile|log.LUTC|log.Lmsgprefix,
	)

	webServer := web.NewServer(
		cfg.GetWeb().GetAddress(),
		webLogger,
		1<<11,
		cfg.GetWeb().GetIdleTimeout(),
		cfg.GetWeb().GetReadTimeOut(),
		cfg.GetWeb().GetWriteTimeout(),
		cfg.GetWeb().GetReadHeaderTimeout(),
	)
	router.RegisterHttpRoutes(webServer.GetMuxHandler())
	go func() {
		log.Panic(webServer.Run())
	}()
	log.Printf("initialize web server in address: %s", cfg.GetWeb().GetAddress())

	grpcServer := grpc.NewServer(cfg.GetGrpc().GetAddress())
	go func() {
		log.Panic(grpcServer.Run())
	}()
	log.Printf("initialize grpc server in address: %s", cfg.GetGrpc().GetAddress())
	auth.RegisterAuthServiceServer(grpcServer.GetServer(), auth.UnimplementedAuthServiceServer{})

	sigint := make(chan os.Signal, 1)
	signal.Notify(sigint, os.Interrupt, syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP, syscall.SIGQUIT, syscall.SIGKILL)
	<-sigint

	errGp := errgroup.Group{}
	errGp.Go(func() error {
		return webServer.GracefulShutdown(ShutdownTimeout)
	})
	errGp.Go(func() error {
		grpcServer.GracefulShutdown(ShutdownTimeout)
		return nil
	})
	if err := errGp.Wait(); err != nil {
		log.Fatal(err)
	}
}
