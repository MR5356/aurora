package config

import (
	"testing"
)

func TestConfig(t *testing.T) {
	_ = Current(
		WithPort(12345),
		WithDebug(true),
		WithGracePeriod(23),
		WithDatabase("sqlite", ":memory:"),
	)

	if Current().Server.Port != 12345 {
		t.Fail()
	}

	if !Current().Server.Debug {
		t.Fail()
	}

	if Current().Server.GracePeriod != 23 {
		t.Fail()
	}

	if Current().Database.DSN != ":memory:" {
		t.Fail()
	}

	if Current().Database.Driver != "sqlite" {
		t.Fail()
	}
}
