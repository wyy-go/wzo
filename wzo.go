package wzo

import (
	"io"
	"os"
	"os/signal"
	"syscall"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/jaeger"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.12.0"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"

	"github.com/wyy-go/wzo/core/config"
	"github.com/wyy-go/wzo/core/log"
	"github.com/wyy-go/wzo/core/pkg/env"
	"github.com/wyy-go/wzo/core/registry"
	"github.com/wyy-go/wzo/core/registry/etcd"
	httpserver "github.com/wyy-go/wzo/core/transport/http/server"
	"github.com/wyy-go/wzo/core/transport/rpc/server"
)

type App struct {
	opts       Options
	zc         *zconfig
	rpcServer  *server.Server
	httpServer *httpserver.Server
}

func (a *App) GetHttpServer() *httpserver.Server {
	return a.httpServer
}

func (a *App) GetRpcServer() *server.Server {
	return a.rpcServer
}

type zconfig struct {
	App struct {
		Mode string
		Name string
	}
	Logger struct {
		Level      string `json:"level"`
		Filename   string `json:"filename"`
		MaxSize    int    `json:"maxSize"`
		MaxBackups int    `json:"maxBackups"`
		MaxAge     int    `json:"maxAge"`
		Compress   bool   `json:"compress"`
	}
	Http struct {
		Addr string
	}
	Rpc struct {
		Addr string
	}
	Tracer struct {
		Addr string
	}
	Registry struct {
		BasePath       string
		EtcdAddr       []string
		UpdateInterval int
	}
}

func New(opts ...Option) *App {
	options := newOptions(opts...)
	zc := &zconfig{}
	if err := config.Unmarshal(zc); err != nil {
		log.Fatal(err)
	}

	if zc.App.Name == "" {
		log.Fatal("config item app.name can't be empty")
	}

	env.Set(zc.App.Mode)

	level, err := zapcore.ParseLevel(zc.Logger.Level)
	if err != nil {
		level = log.InfoLevel
	}
	if env.IsDev() {
		w := &lumberjack.Logger{
			Filename:   zc.Logger.Filename,
			MaxSize:    zc.Logger.MaxSize,
			MaxBackups: zc.Logger.MaxBackups,
			MaxAge:     zc.Logger.MaxAge,
			Compress:   zc.Logger.Compress,
		}
		l := log.NewTee([]io.Writer{os.Stderr, w}, level, log.WithCaller(true))
		log.ResetDefault(l)
	} else {
		w := &lumberjack.Logger{
			Filename:   zc.Logger.Filename,
			MaxSize:    zc.Logger.MaxSize,
			MaxBackups: zc.Logger.MaxBackups,
			MaxAge:     zc.Logger.MaxAge,
			Compress:   zc.Logger.Compress,
		}
		l := log.New(w, level, log.WithCaller(true))
		log.ResetDefault(l)
	}

	app := &App{
		opts: options,
		zc:   zc,
	}

	tracing := false
	if zc.Tracer.Addr != "" {
		setTracerProvider(zc.Tracer.Addr, zc.App.Name)
		tracing = true
	} else {
		setNoExporterTracerProvider(zc.App.Name)
	}

	if app.opts.InitRpcServer != nil {
		app.rpcServer = server.NewServer(
			server.Name(zc.App.Name),
			server.Addr(zc.Rpc.Addr),
			server.BasePath(zc.Registry.BasePath),
			server.UpdateInterval(zc.Registry.UpdateInterval),
			server.EtcdAddr(zc.Registry.EtcdAddr),
			server.Tracing(tracing),
			server.InitRpcServer(app.opts.InitRpcServer),
		)
	}
	mode := "debug"
	if env.IsRelease() {
		mode = "release"
	}
	if app.opts.InitHttpServer != nil {
		var r registry.Registry
		var opts []etcd.Option
		if zc.Registry.BasePath != "" {
			opts = append(opts, etcd.BasePath(zc.Registry.BasePath))
		}
		app.httpServer = httpserver.NewServer(
			httpserver.Name(zc.App.Name),
			httpserver.Addr(zc.Http.Addr),
			httpserver.Mode(mode),
			httpserver.Tracing(tracing),
			httpserver.Registry(r),
			httpserver.InitHttpServer(app.opts.InitHttpServer),
		)
	}

	return app
}

func setTracerProvider(endpoint string, name string) *trace.TracerProvider {
	exp, err := jaeger.New(jaeger.WithCollectorEndpoint(jaeger.WithEndpoint(endpoint)))
	if err != nil {
		log.Fatal(err.Error())
	}
	tp := trace.NewTracerProvider(
		trace.WithSampler(trace.AlwaysSample()),
		trace.WithBatcher(exp),
		trace.WithResource(resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceNameKey.String(name),
		)),
	)

	otel.SetTracerProvider(tp)
	otel.SetTextMapPropagator(
		propagation.NewCompositeTextMapPropagator(
			propagation.TraceContext{},
			propagation.Baggage{},
		),
	)

	return tp
}

func setNoExporterTracerProvider(name string) *trace.TracerProvider {
	tp := trace.NewTracerProvider(
		trace.WithSampler(trace.AlwaysSample()),
		trace.WithResource(resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceNameKey.String(name),
		)),
	)

	otel.SetTracerProvider(tp)
	otel.SetTextMapPropagator(
		propagation.NewCompositeTextMapPropagator(
			propagation.TraceContext{},
			propagation.Baggage{},
		),
	)

	return tp
}

func (a *App) Run() error {
	for _, f := range a.opts.BeforeStart {
		if err := f(); err != nil {
			return err
		}
	}

	if a.rpcServer != nil {
		if err := a.rpcServer.Start(); err != nil {
			return err
		}
	}

	if a.httpServer != nil {
		if err := a.httpServer.Start(); err != nil {
			return err
		}
	}

	for _, f := range a.opts.AfterStart {
		if err := f(); err != nil {
			return err
		}
	}

	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGTERM, syscall.SIGINT, syscall.SIGQUIT)
	log.Infof("received signal %s", <-ch)

	for _, f := range a.opts.BeforeStop {
		if err := f(); err != nil {
			return err
		}
	}

	if a.rpcServer != nil {
		_ = a.rpcServer.Stop()
	}

	if a.httpServer != nil {
		_ = a.httpServer.Stop()
	}

	for _, f := range a.opts.AfterStop {
		if err := f(); err != nil {
			return err
		}
	}

	return nil
}
