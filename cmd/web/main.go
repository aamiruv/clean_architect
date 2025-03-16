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

	_ "github.com/go-sql-driver/mysql"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	_ "github.com/lib/pq"
	"golang.org/x/sync/errgroup"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/reflection"
	_ "modernc.org/sqlite"

	"github.com/amirzayi/clean_architec/api/handler"
	"github.com/amirzayi/clean_architec/internal/auth"
	"github.com/amirzayi/clean_architec/internal/auth/service"
	"github.com/amirzayi/clean_architec/internal/auth/usecase"
	"github.com/amirzayi/clean_architec/internal/user"
	"github.com/amirzayi/clean_architec/pkg/config"
	"github.com/amirzayi/clean_architec/pkg/httpmiddleware"
	"github.com/amirzayi/clean_architec/pkg/interceptor"
	"github.com/amirzayi/clean_architec/pkg/logger/filelog"
	"github.com/amirzayi/clean_architec/pkg/logger/remotelog"
	"github.com/amirzayi/clean_architec/pkg/server/grpcserver"
	"github.com/amirzayi/clean_architec/pkg/server/webserver"
)

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

	// todo: configurable log writer(ex: ELK, prometheus, web-service, etc.)
	// specific logger used for server metric
	serverMetricLogger := slog.NewLogLogger(slog.NewJSONHandler(os.Stdout, nil), slog.LevelInfo)
	// specific logger used for server(grpc&http) panic
	serverPanicLogger := slog.NewLogLogger(slog.NewJSONHandler(os.Stdout, nil), slog.LevelInfo)

	gwMux := runtime.NewServeMux()

	muxHandler := http.NewServeMux()
	muxHandler.Handle("/", gwMux)

	responseTimeMiddleware := func(handler http.Handler) http.Handler {
		return httpmiddleware.MeterResponseTime(handler, serverMetricLogger)
	}
	recoveryMiddleware := func(handler http.Handler) http.Handler {
		return httpmiddleware.Recovery(handler, serverPanicLogger)
	}

	// should metric response time even if panic occurred?
	apiHandler := httpmiddleware.Chain(muxHandler,
		responseTimeMiddleware,
		recoveryMiddleware,
		httpmiddleware.EnforceJSON)

	webServer := webserver.New(
		webserver.WithHandler(apiHandler),
		webserver.WithAddress(cfg.Web().Address()),
		webserver.WithLogger(webServerLogger),
		webserver.WithTimeouts(
			cfg.Web().IdleTimeout(),
			cfg.Web().ReadTimeOut(),
			cfg.Web().WriteTimeout(),
			cfg.Web().ReadHeaderTimeout(),
			cfg.Web().ShutdownTimeout(),
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
		cfg.GRPC().ShutdownTimeout(),
		grpc.MaxRecvMsgSize(cfg.GRPC().MaxReceiveMsgSize()),
		grpc.ReadBufferSize(cfg.GRPC().ReadBufferSize()),
		grpc.ChainUnaryInterceptor(responseTimeMeterInterceptor, recoveryInterceptor),
	)

	if cfg.GRPC().HasReflection() {
		reflection.Register(grpcServer.Server())
	}

	grpcDialOptions := []grpc.DialOption{
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	}
	if err = auth.RegisterGateway(ctx, gwMux, cfg.GRPC().Address(), grpcDialOptions...); err != nil {
		return err
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

	sigint := make(chan os.Signal, 1)
	signal.Notify(sigint, os.Interrupt, syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP, syscall.SIGQUIT, syscall.SIGKILL)

	userService := user.NewService(user.NewSQLRepository(db))

	authService := service.NewAuthService()
	authUseCase := usecase.NewAuthUseCase(authService, userService)

	handler.Register(muxHandler, webServerLogger, authUseCase)

	auth.InitializeAuthServer(grpcServer.Server(), db, authUseCase)

	select {
	case err = <-errCh:
		return err

	case sig := <-sigint:
		logger.Info(fmt.Sprintf("received signal %s", sig))

		errGp := errgroup.Group{}
		errGp.Go(func() error {
			return webServer.GracefulShutdown()
		})
		errGp.Go(func() error {
			grpcServer.GracefulShutdown()
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
