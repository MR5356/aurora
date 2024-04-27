package server

import (
	"context"
	"errors"
	"fmt"
	"github.com/MR5356/aurora/docs"
	"github.com/MR5356/aurora/pkg/config"
	"github.com/MR5356/aurora/pkg/middleware/database"
	"github.com/MR5356/aurora/pkg/response"
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

func New(cfg *config.Config) (server *Server, err error) {
	// init logger level
	if cfg.Server.Debug {
		gin.SetMode(gin.DebugMode)
		logrus.SetLevel(logrus.DebugLevel)
	} else {
		gin.SetMode(gin.ReleaseMode)
	}

	// init database
	database.GetDB()

	engine := gin.Default()
	engine.MaxMultipartMemory = 8 << 20

	// init gin middleware
	engine.Use()

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
	services := []Service{}

	for _, svc := range services {
		if err := svc.Initialize(); err != nil {
			return nil, err
		}
	}

	// controller
	controllers := []Controller{}

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
