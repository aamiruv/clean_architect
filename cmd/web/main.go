// Package main runs a web server with that dependencies.
package main

import (
	"context"
	"database/sql"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/AmirMirzayi/clean_architecture/pkg/logger/remotelog"
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
	flag.StringVar(&configPath, "config", "config.json", "config file path, eg: -config=/path/to/file.json")
	flag.Parse()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	cfg, err := config.LoadConfigOrDefault(configPath)
	if err != nil {
		return err
	}

	fileLogger := filelog.New(filelog.LogHourly, "log")
	remoteLogger := remotelog.New(cfg.LoggerURL())
	// log on console & file & http url at same time
	logWriter := io.MultiWriter(os.Stdout, fileLogger, remoteLogger)
	log.SetFlags(log.Ltime | log.Lshortfile | log.LUTC)
	log.SetOutput(logWriter)

	webLoggerFile := filelog.New(filelog.LogHourly, "weblog")
	webServerLogger := log.New(
		webLoggerFile, "",
		log.Ltime|log.Lshortfile|log.LUTC|log.Lmsgprefix,
	)

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

	errCh := make(chan error)

	var (
		db *sql.DB
		wg sync.WaitGroup
	)

	go func() {
		wg.Add(1)
		defer wg.Done()
		db, err = sql.Open("mysql", cfg.DB().ConnectionString())
		if err != nil {
			errCh <- err
		}
		if err = db.Ping(); err != nil {
			errCh <- err
		}
	}()

	go func() {
		wg.Add(1)
		defer wg.Done()
		if err = webServer.Run(); err != nil {
			errCh <- fmt.Errorf("failed to initialize web server: %w", err)
		}
	}()
	fmt.Printf("web server initialized in address: %s\n\r", cfg.Web().Address())

	go func() {
		wg.Add(1)
		defer wg.Done()
		if err = grpcServer.Run(); err != nil {
			errCh <- fmt.Errorf("failed to initialize grpc server: %w", err)
		}
	}()
	fmt.Printf("grpc server initialized in address: %s\n\r", cfg.GRPC().Address())

	if cfg.GRPC().HasReflection() {
		reflection.Register(grpcServer.Server())
	}

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
		fmt.Printf("received signal %s", sig)

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
