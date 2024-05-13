package runner

import (
	"context"
	"errors"
	"fmt"
	_ "github.com/MR5356/aurora/pkg/log"
	"github.com/MR5356/aurora/pkg/response"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/sirupsen/logrus"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

type Config struct {
	Host  string
	Port  int
	Token string
	Debug bool
}

func Run(cfg *Config) error {
	conn, _, err := websocket.DefaultDialer.Dial(fmt.Sprintf("ws://%s:%d/api/v1/runner/conn?token=%s", cfg.Host, cfg.Port, cfg.Token), nil)
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	// 发送消息
	err = conn.WriteMessage(websocket.TextMessage, []byte("Hello, world!"))
	if err != nil {
		log.Fatal(err)
	}

	// 读取消息
	for {
		messageType, p, err := conn.ReadMessage()
		if err != nil {
			log.Fatal(err)
		}
		logrus.Infof("Received message type: %d", messageType)
		log.Println("Received message:", string(p))
	}
}

func Run1(cfg *Config) error {
	if cfg.Debug {
		gin.SetMode(gin.DebugMode)
		logrus.SetLevel(logrus.DebugLevel)
	} else {
		gin.SetMode(gin.ReleaseMode)
	}

	engine := gin.Default()
	engine.MaxMultipartMemory = 8 << 20
	engine.GET("/ping", func(c *gin.Context) {
		response.Success(c, nil)
	})

	engine.GET("/metrics", func(handler http.Handler) gin.HandlerFunc {
		return func(context *gin.Context) {
			handler.ServeHTTP(context.Writer, context.Request)
		}
	}(promhttp.Handler()))

	server := &http.Server{
		Addr:    fmt.Sprintf(":%d", cfg.Port),
		Handler: engine,
	}

	go func() {
		logrus.Infof("server running on port %d", cfg.Port)
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
