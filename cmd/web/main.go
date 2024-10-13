// Package main runs a web server with that dependencies.
package main

import (
	"context"
	"database/sql"
	"flag"
	"fmt"
	"io"
	"log/slog"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"golang.org/x/sync/errgroup"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/reflection"

	"github.com/AmirMirzayi/clean_architecture/api/httprouter"
	"github.com/AmirMirzayi/clean_architecture/internal/auth"
	"github.com/AmirMirzayi/clean_architecture/pkg/config"
	"github.com/AmirMirzayi/clean_architecture/pkg/logger/filelog"
	"github.com/AmirMirzayi/clean_architecture/pkg/logger/remotelog"
	"github.com/AmirMirzayi/clean_architecture/pkg/server/grpcserver"
	"github.com/AmirMirzayi/clean_architecture/pkg/server/webserver"
)

const ShutdownTimeout = 5 * time.Second

func main() {
	if err := run(); err != nil {
		panic(err)
	}
}

func run() error {
	var configPath string
	flag.StringVar(&configPath, "config", "config.json", "config file path, eg: -config=/path/to/file.json")
	flag.Parse()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	cfg, err := config.LoadConfigOrDefault(configPath)
	if err != nil {
		return err
	}

	var logWriters []io.Writer
	if cfg.Logger().Console() {
		logWriters = append(logWriters, os.Stdout)
	}
	if cfg.Logger().Directory() != "" {
		fileLogger := filelog.New(filelog.LoggerType(cfg.Logger().FileCreationMode()), cfg.Logger().Directory())
		logWriters = append(logWriters, fileLogger)
	}
	if cfg.Logger().RemoteURL() != "" {
		remoteLogger := remotelog.New(cfg.Logger().RemoteURL())
		logWriters = append(logWriters, remoteLogger)
	}

	logWriter := io.MultiWriter(logWriters...)
	logger := slog.New(slog.NewJSONHandler(logWriter, &slog.HandlerOptions{AddSource: true, Level: slog.Level(cfg.Logger().Level())}))
	// no need to pass logger to another part of application
	slog.SetDefault(logger)

	webServerLogFile := filelog.New(filelog.LogHourly, "weblog")
	webServerLogWriter := io.MultiWriter(os.Stdout, webServerLogFile)
	webServerLogger := slog.NewLogLogger(slog.NewJSONHandler(webServerLogWriter, nil), slog.LevelInfo)

	gwMux := runtime.NewServeMux()

	muxHandler := httprouter.New()
	muxHandler.Handle("/", gwMux)

	webServer := webserver.New(
		webserver.WithHandler(muxHandler),
		webserver.WithAddress(cfg.Web().Address()),
		webserver.WithLogger(webServerLogger),
		webserver.WithTimeouts(
			cfg.Web().IdleTimeout(),
			cfg.Web().ReadTimeOut(),
			cfg.Web().WriteTimeout(),
			cfg.Web().ReadHeaderTimeout(),
		),
	)

	grpcServer := grpcserver.New(
		cfg.GRPC().Address(),
		grpc.MaxRecvMsgSize(cfg.GRPC().MaxReceiveMsgSize()),
		grpc.ReadBufferSize(cfg.GRPC().ReadBufferSize()),
	)

	if cfg.GRPC().HasReflection() {
		reflection.Register(grpcServer.Server())
	}

	errCh := make(chan error)

	var (
		db *sql.DB
		wg sync.WaitGroup
	)

	auth.InitializeAuthServer(grpcServer.Server(), db)

	go func() {
		wg.Add(1)
		defer wg.Done()
		db, err = sql.Open("mysql", cfg.DB().ConnectionString())
		if err != nil {
			errCh <- fmt.Errorf("failed to open database connection: %w", err)
		}
		if err = db.Ping(); err != nil {
			errCh <- fmt.Errorf("failed to ping database: %w", err)
		}
	}()

	go func() {
		wg.Add(1)
		defer wg.Done()
		if err = webServer.Run(); err != nil {
			errCh <- fmt.Errorf("failed to run web server: %w", err)
		}
	}()
	logger.Info("web server initialized", "address", cfg.Web().Address())

	go func() {
		wg.Add(1)
		defer wg.Done()
		if err = grpcServer.Run(); err != nil {
			errCh <- fmt.Errorf("failed to run grpc server: %w", err)
		}
	}()
	logger.Info("grpc server initialized", "address", cfg.GRPC().Address())

	go func() {
		wg.Wait()
		close(errCh)
	}()

	grpcDialOptions := []grpc.DialOption{
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	}
	if err = auth.RegisterGateway(ctx, gwMux, cfg.GRPC().Address(), grpcDialOptions...); err != nil {
		return err
	}

	sigint := make(chan os.Signal, 1)
	signal.Notify(sigint, os.Interrupt, syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP, syscall.SIGQUIT, syscall.SIGKILL)

	select {
	case err = <-errCh:
		return err

	case sig := <-sigint:
		logger.Info("received signal", "", sig)

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
