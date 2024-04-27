package server

import (
	"aurora/pkg/config"
	"context"
	"testing"
	"time"
)

func TestNew(t *testing.T) {
	cfg := config.New(config.WithDatabase("sqlite", ":memory:"))

	svc, err := New(cfg)
	if err != nil {
		t.Fail()
	}

	if svc == nil {
		t.Fail()
	}
}

func TestServer_Run(t *testing.T) {
	cfg := config.New(config.WithDatabase("sqlite", ":memory:"))

	svc, err := New(cfg)
	if err != nil {
		t.Fail()
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	go func() {
		if err := svc.Run(); err != nil {
			t.Fail()
		}
	}()

	<-ctx.Done()
}
