package server

import (
	"context"
	"embed"
	"errors"
	"fmt"
	"github.com/MR5356/aurora/docs"
	"github.com/MR5356/aurora/internal/config"
	"github.com/MR5356/aurora/internal/domain/health"
	"github.com/MR5356/aurora/internal/domain/host"
	"github.com/MR5356/aurora/internal/domain/notify"
	"github.com/MR5356/aurora/internal/domain/pipeline"
	"github.com/MR5356/aurora/internal/domain/plugin"
	"github.com/MR5356/aurora/internal/domain/schedule"
	"github.com/MR5356/aurora/internal/domain/script"
	"github.com/MR5356/aurora/internal/domain/system"
	"github.com/MR5356/aurora/internal/domain/user"
	"github.com/MR5356/aurora/internal/domain/user/oauth"
	"github.com/MR5356/aurora/internal/infrastructure/database"
	"github.com/MR5356/aurora/internal/infrastructure/eventbus"
	"github.com/MR5356/aurora/internal/response"
	ginmiddleware2 "github.com/MR5356/aurora/internal/server/ginmiddleware"
	_ "github.com/MR5356/aurora/pkg/log"
	"github.com/MR5356/aurora/pkg/util/structutil"
	"github.com/gin-contrib/gzip"
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/sirupsen/logrus"
	"os"
	"os/signal"
	"syscall"
	"time"

	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"net/http"
)

type Server struct {
	engine *gin.Engine
}

//go:embed static
var fs embed.FS

func New(cfg *config.Config) (server *Server, err error) {
	logrus.Infof("config: \n%+v", structutil.Struct2String(config.Current()))

	// init log level
	if cfg.Server.Debug {
		gin.SetMode(gin.DebugMode)
		logrus.SetLevel(logrus.DebugLevel)
	} else {
		gin.SetMode(gin.ReleaseMode)
	}

	os.MkdirAll(cfg.Server.PluginPath, os.ModePerm)

	// init infrastructure
	oauth.NewOAuthManager(cfg)
	user.NewJWTService(cfg)
	database.NewDatabase(cfg)
	eventbus.NewEventBus(cfg)

	engine := gin.Default()
	engine.MaxMultipartMemory = 8 << 20

	// init gin infrastructure
	engine.Use(
		ginmiddleware2.Record(),
		ginmiddleware2.Static("/", ginmiddleware2.NewStaticFileSystem(fs, "static")),
	)

	// 404
	engine.NoRoute(func(ctx *gin.Context) {
		response.Error(ctx, response.CodeNotFound)
	})

	api := engine.Group(cfg.Server.Prefix)
	api.Use(
		gzip.Gzip(gzip.DefaultCompression),
		ginmiddleware2.MustLogin(),
	)

	// metrics
	engine.GET("/api/v1/metrics", func(handler http.Handler) gin.HandlerFunc {
		return func(context *gin.Context) {
			handler.ServeHTTP(context.Writer, context.Request)
		}
	}(promhttp.Handler()))

	// swagger
	docs.SwaggerInfo.Title = "Aurora API"
	docs.SwaggerInfo.Description = "Aurora API"
	docs.SwaggerInfo.Version = "1.0"
	docs.SwaggerInfo.BasePath = cfg.Server.Prefix
	api.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// service
	services := []Service{
		script.GetService(),
		user.GetService(),
		system.GetService(),
		notify.GetService(),
		pipeline.GetService(),
		host.GetService(),
		health.GetService(),
		schedule.GetService(),
	}

	for _, svc := range services {
		if err := svc.Initialize(); err != nil {
			return nil, err
		}
	}

	// controller
	controllers := []Controller{
		schedule.NewController(),
		user.NewController(),
		system.NewController(),
		notify.NewController(),
		pipeline.NewController(),
		host.NewController(),
		health.NewController(),
		plugin.NewController(),
		script.NewController(),
	}

	for _, ctl := range controllers {
		ctl.RegisterRoute(api)
	}

	server = &Server{
		engine: engine,
	}
	return
}

func (s *Server) Run() error {
	server := &http.Server{
		Addr:    fmt.Sprintf(":%d", config.Current().Server.Port),
		Handler: s.engine,
	}

	go func() {
		logrus.Infof("server running on port %d", config.Current().Server.Port)
		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			logrus.Fatalf("server listening error: %v", err)
		}
	}()

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)

	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(config.Current().Server.GracePeriod)*time.Second)
	defer cancel()

	ch := <-sig
	logrus.Infof("server receive signal: %s", ch.String())
	return server.Shutdown(ctx)
}

func (s *Server) Shutdown() {

}
