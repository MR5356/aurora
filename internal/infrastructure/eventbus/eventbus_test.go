package eventbus

import (
	"testing"
)

func TestGetEventBus(t *testing.T) {
	ins1 := GetEventBus()
	ins2 := GetEventBus()
	if ins1 != ins2 {
		t.Fail()
	}
}

func TestConcurrentGetEventBus(t *testing.T) {
	instances := make([]EventBus, 0)
	ch := make(chan EventBus, 100)

	for i := 0; i < 100; i++ {
		go func() {
			ch <- GetEventBus()
		}()
	}

	for i := 0; i < 100; i++ {
		instances = append(instances, <-ch)
	}

	for i := 1; i < len(instances); i++ {
		if instances[0] != instances[i] {
			t.Fail()
		}
	}
}
