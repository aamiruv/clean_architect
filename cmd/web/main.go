package main

import (
	"context"
	"database/sql"
	"errors"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"golang.org/x/sync/errgroup"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"github.com/AmirMirzayi/clean_architecture/api/router"
	"github.com/AmirMirzayi/clean_architecture/internal/auth"
	"github.com/AmirMirzayi/clean_architecture/pkg/config"
	"github.com/AmirMirzayi/clean_architecture/pkg/logger/filelog"
	"github.com/AmirMirzayi/clean_architecture/pkg/logger/weblog"
	"github.com/AmirMirzayi/clean_architecture/pkg/server/grpcserver"
	"github.com/AmirMirzayi/clean_architecture/pkg/server/webserver"
)

const ShutdownTimeout = 5 * time.Second

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {
	var configPath string
	flag.StringVar(&configPath, "config", "config.json", "config filelog path, eg: -config=/path/to/filelog.json")
	flag.Parse()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	cfg := config.LoadConfig(configPath)

	db, err := sql.Open("mysql", cfg.DB().ConnectionString())
	if err != nil {
		return err
	}
	if err = db.Ping(); err != nil {
		return err
	}

	fileLogger := filelog.New(filelog.LogHourly, "log")
	theWebLog := weblog.New(cfg.LoggerURL())
	// log on console & file & http url at same time
	logWriter := io.MultiWriter(os.Stdout, fileLogger, theWebLog)
	log.SetFlags(log.Ltime | log.Lshortfile | log.LUTC)
	log.SetOutput(logWriter)

	webLoggerFile := filelog.New(filelog.LogHourly, "weblog")
	webLogger := log.New(
		webLoggerFile, "",
		log.Ltime|log.Lshortfile|log.LUTC|log.Lmsgprefix,
	)

	gwMux := runtime.NewServeMux()

	muxHandler := router.New()
	muxHandler.Handle("/", gwMux)

	webServer := webserver.New(
		webserver.WithMuxHandler(muxHandler),
		webserver.WithAddress(cfg.Web().Address()),
		webserver.WithLogger(webLogger),
		webserver.WithTimeout(
			cfg.Web().IdleTimeout(),
			cfg.Web().ReadTimeOut(),
			cfg.Web().WriteTimeout(),
			cfg.Web().ReadHeaderTimeout(),
		),
	)

	errCh := make(chan error)
	defer close(errCh)

	go func() {
		if err = webServer.Run(); !errors.Is(err, http.ErrServerClosed) {
			errCh <- fmt.Errorf("failed to initialize web server: %w", err)
		}
	}()
	log.Printf("web server initialized in address: %s", cfg.Web().Address())

	grpcServer := grpcserver.New(cfg.Grpc().Address())
	auth.InitializeAuthServer(grpcServer.Server(), db)
	go func() {
		if err = grpcServer.Run(); err != nil {
			errCh <- fmt.Errorf("failed to initialize grpc server: %w", err)
		}
	}()
	log.Printf("grpc server initialized in address: %s", cfg.Grpc().Address())

	dialOptions := []grpc.DialOption{
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	}
	if err = auth.RegisterGateway(ctx, gwMux, cfg.Grpc().Address(), dialOptions...); err != nil {
		return err
	}

	sigint := make(chan os.Signal, 1)
	signal.Notify(sigint, os.Interrupt, syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP, syscall.SIGQUIT, syscall.SIGKILL)

	select {
	case err = <-errCh:
		return err

	case sig := <-sigint:
		log.Printf("received signal %s", sig)

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
		if err = errGp.Wait(); err != nil {
			return err
		}
	}

	return nil
}
