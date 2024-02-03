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
	"github.com/AmirMirzayi/clean_architecture/pkg/web"
	"google.golang.org/grpc"
	"log"
	"net"
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

type tmp struct {
	auth.UnimplementedAuthServiceServer
}

func (t tmp) Register(context.Context, *auth.RegisterRequest) (*auth.RegisterResponse, error) {
	return &auth.RegisterResponse{UserId: "amir"}, nil
}

func main() {
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

	lis, err := net.Listen("tcp", "localhost:8070")
	if err != nil {
		panic(err)
	}
	grpcServer := grpc.NewServer()

	auth.RegisterAuthServiceServer(grpcServer, &tmp{})
	go grpcServer.Serve(lis)

	sigint := make(chan os.Signal, 1)
	signal.Notify(sigint, os.Interrupt, syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP, syscall.SIGQUIT, syscall.SIGKILL)
	<-sigint
	if err := webServer.GracefulShutdown(5 * time.Second); err != nil {
		log.Println(err)
	}
	grpcServer.GracefulStop()
}
