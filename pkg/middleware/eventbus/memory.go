package eventbus

import (
	evbus "github.com/asaskevich/EventBus"
	"sync"
)

type MemoryEventBus struct {
	mutexLockMap sync.Map
	evbus        evbus.Bus
}

func (eb *MemoryEventBus) Unlock(key string) error {
	//TODO implement me
	panic("implement me")
}

func NewMemoryEventBus() *MemoryEventBus {
	return &MemoryEventBus{
		evbus:        evbus.New(),
		mutexLockMap: sync.Map{},
	}
}

func (eb *MemoryEventBus) TryLock(key string) error {
	if _, locked := eb.mutexLockMap.LoadOrStore(key, true); locked {
		return ErrAlreadyLocked
	}
	return nil
}

func (eb *MemoryEventBus) UnLock(key string) error {
	eb.mutexLockMap.Delete(key)
	return nil
}

func (eb *MemoryEventBus) Subscribe(topic string, handler interface{}) error {
	return eb.evbus.Subscribe(topic, handler)
}

func (eb *MemoryEventBus) UnSubscribe(topic string, handler interface{}) error {
	return eb.evbus.Unsubscribe(topic, handler)
}

func (eb *MemoryEventBus) Publish(topic string, data interface{}) error {
	eb.evbus.Publish(topic, data)
	return nil
}
