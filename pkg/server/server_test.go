package server

import (
	"github.com/MR5356/aurora/pkg/config"
	"testing"
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

// can not pass on github actions
//func TestServer_Run(t *testing.T) {
//	cfg := config.New(config.WithDatabase("sqlite", ":memory:"))
//
//	svc, err := New(cfg)
//	if err != nil {
//		t.Fail()
//	}
//
//	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
//	defer cancel()
//
//	go func() {
//		if err := svc.Run(); err != nil {
//			t.Fail()
//		}
//	}()
//
//	<-ctx.Done()
//}
