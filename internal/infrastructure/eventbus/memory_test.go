package eventbus

import (
	"testing"
)

func TestNewMemoryEventBus(t *testing.T) {
	bus := NewMemoryEventBus()
	if bus == nil {
		t.Log("New Memory EventBus not Created!")
		t.Fail()
	}
}

func TestMemoryEventBus_Subscribe(t *testing.T) {
	bus := NewMemoryEventBus()
	if bus.Subscribe("topic", func() {}) != nil {
		t.Fail()
	}

	if bus.Subscribe("topic", "String") == nil {
		t.Fail()
	}
}

func TestMemoryEventBus_UnSubscribe(t *testing.T) {
	bus := NewMemoryEventBus()
	handler := func() {}
	topic := "topic"
	bus.Subscribe(topic, handler)
	if bus.UnSubscribe(topic, handler) != nil {
		t.Fail()
	}

	if bus.UnSubscribe(topic, handler) == nil {
		t.Fail()
	}
}

func TestMemoryEventBus_Publish(t *testing.T) {
	bus := NewMemoryEventBus()
	_ = bus.Subscribe("topic", func(a int) {
		if a != 10 {
			t.Fail()
		}
	})

	_ = bus.Publish("topic", 10)
}

func TestMemoryEventBus_TryLock(t *testing.T) {
	bus := NewMemoryEventBus()
	lockKey := "lock"
	if err := bus.TryLock(lockKey); err != nil {
		t.Fail()
	}

	if err := bus.TryLock(lockKey); err != ErrAlreadyLocked {
		t.Fail()
	}

	if err := bus.UnLock(lockKey); err != nil {
		t.Fail()
	}

	if err := bus.TryLock(lockKey); err != nil {
		t.Fail()
	}

	bus.UnLock(lockKey)
}

func TestMemoryEventBus_UnLock(t *testing.T) {
	bus := NewMemoryEventBus()
	lockKey := "lock"
	bus.TryLock(lockKey)
	if err := bus.UnLock(lockKey); err != nil {
		t.Fail()
	}
}
