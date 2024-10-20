// Package main runs a web server with that dependencies.
package main

import (
	"context"
	"database/sql"
	"flag"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	_ "github.com/lib/pq"
	"golang.org/x/sync/errgroup"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/reflection"
	_ "modernc.org/sqlite"

	"github.com/AmirMirzayi/clean_architecture/api/httprouter"
	"github.com/AmirMirzayi/clean_architecture/internal/auth"
	"github.com/AmirMirzayi/clean_architecture/pkg/config"
	"github.com/AmirMirzayi/clean_architecture/pkg/interceptor"
	"github.com/AmirMirzayi/clean_architecture/pkg/logger/filelog"
	"github.com/AmirMirzayi/clean_architecture/pkg/logger/remotelog"
	"github.com/AmirMirzayi/clean_architecture/pkg/middleware"
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

	// todo: configurable log writer(ex: elastic, prometheus, web-service, etc.)
	// specific logger used for server metric
	serverMetricLogger := slog.NewLogLogger(slog.NewJSONHandler(os.Stdout, nil), slog.LevelInfo)
	// specific logger used for server(grpc&http) panic
	serverPanicLogger := slog.NewLogLogger(slog.NewJSONHandler(os.Stdout, nil), slog.LevelInfo)

	gwMux := runtime.NewServeMux()

	muxHandler := httprouter.New()
	muxHandler.Handle("/", gwMux)

	responseTimeMiddleware := func(handler http.Handler) http.Handler {
		return middleware.MeterResponseTime(handler, serverMetricLogger)
	}
	recoveryMiddleware := func(handler http.Handler) http.Handler {
		return middleware.Recovery(handler, serverPanicLogger)
	}

	// should metric response time even if panic occurred?
	handler := middleware.Chain(muxHandler, responseTimeMiddleware, recoveryMiddleware)

	webServer := webserver.New(
		webserver.WithHandler(handler),
		webserver.WithAddress(cfg.Web().Address()),
		webserver.WithLogger(webServerLogger),
		webserver.WithTimeouts(
			cfg.Web().IdleTimeout(),
			cfg.Web().ReadTimeOut(),
			cfg.Web().WriteTimeout(),
			cfg.Web().ReadHeaderTimeout(),
		),
	)

	responseTimeMeterInterceptor := func(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {
		return interceptor.ResponseTimeMeter(ctx, req, info, handler, serverMetricLogger)
	}
	recoveryInterceptor := func(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {
		return interceptor.Recovery(ctx, req, info, handler, serverPanicLogger)
	}

	grpcServer := grpcserver.New(
		cfg.GRPC().Address(),
		grpc.MaxRecvMsgSize(cfg.GRPC().MaxReceiveMsgSize()),
		grpc.ReadBufferSize(cfg.GRPC().ReadBufferSize()),
		grpc.ChainUnaryInterceptor(responseTimeMeterInterceptor, recoveryInterceptor),
	)

	if cfg.GRPC().HasReflection() {
		reflection.Register(grpcServer.Server())
	}

	errCh := make(chan error)

	var (
		db *sql.DB
		wg sync.WaitGroup
	)

	wg.Add(3)

	go func() {
		defer wg.Done()
		db, err = sql.Open(cfg.DB().Driver(), cfg.DB().ConnectionString())
		if err != nil {
			errCh <- fmt.Errorf("failed to open database connection: %w", err)
		}
		if err = db.Ping(); err != nil {
			errCh <- fmt.Errorf("failed to ping database: %w", err)
		}
	}()

	go func() {
		defer wg.Done()
		if err = webServer.Run(); err != nil {
			errCh <- fmt.Errorf("failed to run web server: %w", err)
		}
	}()
	logger.Info(fmt.Sprintf("web server initialized on %s", cfg.Web().Address()))

	go func() {
		defer wg.Done()
		if err = grpcServer.Run(); err != nil {
			errCh <- fmt.Errorf("failed to run grpc server: %w", err)
		}
	}()
	logger.Info(fmt.Sprintf("grpc server initialized on %s", cfg.GRPC().Address()))

	go func() {
		wg.Wait()
		close(errCh)
	}()

	auth.InitializeAuthServer(grpcServer.Server(), db)

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
		logger.Info(fmt.Sprintf("received signal %s", sig))

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
