package pluginhub

import (
	"context"
	"errors"
	"fmt"
	"github.com/MR5356/aurora/pkg/response"
	"github.com/MR5356/aurora/pkg/util/fileutil"
	"github.com/MR5356/aurora/pkg/util/structutil"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

type Server struct {
	engine *gin.Engine
}

func New() *Server {
	engine := gin.Default()
	engine.MaxMultipartMemory = 8 << 20

	engine.NoRoute(func(context *gin.Context) {
		response.Error(context, response.CodeNotFound)
	})

	engine.Static("/plugin", "./_plugins")

	cfg := new(TaskConfig)
	fileutil.NewStructFromFile("./_plugins/checkout/task.yml", cfg)
	logrus.Infof("%+v", structutil.Struct2String(cfg))

	return &Server{engine: engine}
}

func (s *Server) Run(port int) error {
	server := &http.Server{
		Addr:    fmt.Sprintf(":%d", port),
		Handler: s.engine,
	}

	go func() {
		logrus.Infof("server running on port %d", port)
		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			logrus.Fatalf("server run failed, err: %v", err)
		}
	}()

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
	defer cancel()

	ch := <-sig
	logrus.Infof("server receive signal: %s", ch.String())
	return server.Shutdown(ctx)
}
