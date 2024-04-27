package main

import (
	"aurora/pkg/config"
	"aurora/pkg/server"
	"github.com/sirupsen/logrus"
)

func main() {
	cfg := config.New(
		config.WithDebug(true),
	)

	svc, err := server.New(cfg)
	if err != nil {
		logrus.Fatalf("server.New failed, err: %v", err)
	}

	if err := svc.Run(); err != nil {
		logrus.Fatalf("server.Run failed, err: %v", err)
	}
}
