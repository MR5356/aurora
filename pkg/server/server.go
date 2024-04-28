package server

import (
	"context"
	"embed"
	"errors"
	"fmt"
	"github.com/MR5356/aurora/docs"
	"github.com/MR5356/aurora/pkg/config"
	"github.com/MR5356/aurora/pkg/domain/schedule"
	"github.com/MR5356/aurora/pkg/domain/user"
	"github.com/MR5356/aurora/pkg/domain/user/oauth"
	_ "github.com/MR5356/aurora/pkg/log"
	"github.com/MR5356/aurora/pkg/middleware/database"
	"github.com/MR5356/aurora/pkg/middleware/eventbus"
	"github.com/MR5356/aurora/pkg/response"
	"github.com/MR5356/aurora/pkg/server/ginmiddleware"
	"github.com/MR5356/aurora/pkg/util/structutil"
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

	// init middleware
	oauth.NewOAuthManager(cfg)
	user.NewJWTService(cfg)
	database.NewDatabase(cfg)
	eventbus.NewEventBus(cfg)

	engine := gin.Default()
	engine.MaxMultipartMemory = 8 << 20

	// init gin middleware
	engine.Use(
		ginmiddleware.Record(),
		ginmiddleware.Static("/", ginmiddleware.NewStaticFileSystem(fs, "static")),
	)

	// 404
	engine.NoRoute(func(ctx *gin.Context) {
		response.Error(ctx, response.CodeNotFound)
	})

	api := engine.Group(cfg.Server.Prefix)

	// metrics
	api.GET("/metrics", func(handler http.Handler) gin.HandlerFunc {
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
		schedule.GetService(),
		user.GetService(),
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
