package app

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/mechta-market/gotemplate/internal/config"
	"google.golang.org/grpc"
)

type App struct {
	globalTracerCloser io.Closer

	pgpool *pgxpool.Pool

	grpcServer *GrpcServer
	httpServer *GrpcGwServer

	ctx       context.Context
	ctxCancel context.CancelFunc

	exitCode int
}

func (a *App) Init() {
	var err error

	a.ctx, a.ctxCancel = context.WithCancel(context.Background())

	// logger
	{
		slogOpts := &slog.HandlerOptions{
			AddSource: config.Conf.Debug,
			Level:     slog.LevelInfo,
		}
		switch strings.ToLower(config.Conf.LogLevel) {
		case "debug":
			slogOpts.Level = slog.LevelDebug
		case "info":
			slogOpts.Level = slog.LevelInfo
		case "warn":
			slogOpts.Level = slog.LevelWarn
		case "error":
			slogOpts.Level = slog.LevelError
		}
		logger := slog.New(slog.NewJSONHandler(os.Stdout, slogOpts))
		if config.Conf.Debug {
			logger = slog.New(slog.NewTextHandler(os.Stdout, slogOpts))
		}
		slog.SetDefault(logger)
	}

	// globalTracer
	// {
	// 	if config.Conf.WithTracing && config.Conf.JaegerAddress != "" {
	// 		slog.Info("tracing enabled")
	// 		_, a.globalTracerCloser, err = tracerInitGlobal(config.Conf.JaegerAddress, constant.ServiceName)
	// 		errCheck(err, "tracerInitGlobal")
	// 	}
	// }

	// pgpool
	{
		pgConf, err := pgxpool.ParseConfig(config.Conf.PgDsn)
		errCheck(err, "pgxpool.ParseConfig")

		pgConf.MaxConns = 10
		pgConf.MinConns = 2
		pgConf.MaxConnLifetime = 3 * time.Minute
		pgConf.MaxConnIdleTime = time.Minute
		pgConf.HealthCheckPeriod = 15 * time.Second

		a.pgpool, err = pgxpool.NewWithConfig(context.Background(), pgConf)
		errCheck(err, "pgxpool.NewWithConfig")
	}

	// grpc server
	{
		a.grpcServer = NewGrpcServer("main", func(server *grpc.Server) {
			// server registers
		})
	}

	// http-gw server
	{
		a.httpServer, err = NewGrpcGwServer(
			[]func(ctx context.Context, mux *runtime.ServeMux, conn *grpc.ClientConn) error{
				// server registers
			},
			[]HttpRoute{
				{
					Method: http.MethodGet, Path: "/tst", Handler: func(w http.ResponseWriter, r *http.Request, _ map[string]string) {
						slog.Error("test error", "error", errors.New("test error"))
					},
				},
			})
		errCheck(err, "NewGrpcGwServer")
	}
}

func (a *App) PreStartHook() {
	slog.Info("PreStartHook")
}

func (a *App) Start() {
	slog.Info("Starting")

	// grpc server
	{
		err := a.grpcServer.Start()
		errCheck(err, "grpcServer.Start")
	}

	// grpc-gw server
	{
		a.httpServer.Start()
	}
}

func (a *App) Listen() {
	signalCtx, signalCtxCancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer signalCtxCancel()

	// wait signal
	sig := <-signalCtx.Done()
	slog.Info("Shutdown signal received", "signal", sig)
}

func (a *App) Stop() {
	slog.Info("Shutting down...")

	// stop context
	a.ctxCancel()

	// grpc-gw server
	{
		err := a.httpServer.Stop(context.Background())
		if err != nil {
			slog.Error("httpServer.Stop", "error", err)
			a.exitCode = 1
		}
	}

	// grpc server
	a.grpcServer.Stop()
}

func (a *App) WaitJobs() {
	slog.Info("waiting jobs")
}

func (a *App) Exit() {
	slog.Info("Exit")

	if a.globalTracerCloser != nil {
		_ = a.globalTracerCloser.Close()
	}

	// flush stdout

	os.Exit(a.exitCode)
}

func errCheck(err error, msg string) {
	if err != nil {
		if msg != "" {
			err = fmt.Errorf("%s: %w", msg, err)
		}
		slog.Error(err.Error())
		os.Exit(1)
	}
}
