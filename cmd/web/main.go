package main

import (
	"context"
	"database/sql"
	"errors"
	"flag"
	"github.com/AmirMirzayi/clean_architecture/api/router"
	"github.com/AmirMirzayi/clean_architecture/internal/auth"
	"github.com/AmirMirzayi/clean_architecture/pkg/config"
	"github.com/AmirMirzayi/clean_architecture/pkg/logger"
	"github.com/AmirMirzayi/clean_architecture/pkg/logger/file"
	weblog "github.com/AmirMirzayi/clean_architecture/pkg/logger/web"
	"github.com/AmirMirzayi/clean_architecture/pkg/server/grpc"
	"github.com/AmirMirzayi/clean_architecture/pkg/server/web"
	_ "github.com/go-sql-driver/mysql"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"golang.org/x/sync/errgroup"
	grpc2 "google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log"
	"net/http"
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

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	cfg := config.LoadConfig(configPath)

	db, err := sql.Open("mysql", cfg.GetDB().GetConnectionString())
	if err != nil {
		panic(err)
	}
	if err = db.Ping(); err != nil {
		panic(err)
	}

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
		web.WithAddress(cfg.GetWeb().GetAddress()),
		web.WithLogger(webLogger),
		web.WithTimeout(
			cfg.GetWeb().GetIdleTimeout(),
			cfg.GetWeb().GetReadTimeOut(),
			cfg.GetWeb().GetWriteTimeout(),
			cfg.GetWeb().GetReadHeaderTimeout(),
		),
	)
	router.RegisterHttpRoutes(webServer.GetMuxHandler())
	go func() {
		if err := webServer.Run(); !errors.Is(err, http.ErrServerClosed) {
			log.Panic(err)
		}
	}()
	log.Printf("initialize web server in address: %s", cfg.GetWeb().GetAddress())

	grpcServer := grpc.NewServer(cfg.GetGrpc().GetAddress())
	go func() {
		log.Panic(grpcServer.Run())
	}()
	log.Printf("initialize grpc server in address: %s", cfg.GetGrpc().GetAddress())

	auth.InitializeAuthServer(grpcServer.GetServer(), db)

	mux := runtime.NewServeMux()
	webServer.GetMuxHandler().Handle("/", mux)

	dialOptions := []grpc2.DialOption{
		grpc2.WithTransportCredentials(insecure.NewCredentials()),
	}
	if err := auth.RegisterGateway(ctx, mux, cfg.GetGrpc().GetAddress(), dialOptions...); err != nil {
		log.Panic(err)
	}

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
	errGp.Go(func() error {
		return db.Close()
	})
	if err := errGp.Wait(); err != nil {
		log.Fatal(err)
	}
}
