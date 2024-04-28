package eventbus

import (
	"errors"
	"github.com/MR5356/aurora/pkg/config"
	"sync"
)

var (
	eb   EventBus
	once sync.Once
)

var (
	ErrAlreadyLocked = errors.New("already locked")
)

type EventBus interface {
	// Subscribe subscribe event
	Subscribe(topic string, handler interface{}) error

	// UnSubscribe unsubscribe event
	UnSubscribe(topic string, handler interface{}) error

	// Publish publish event
	Publish(topic string, data interface{}) error

	// TryLock try lock, if lock success return nil
	TryLock(key string) error

	// UnLock unLock
	UnLock(key string) error
}

func NewEventBus(cfg *config.Config) EventBus {
	once.Do(func() {
		eb = NewMemoryEventBus()
	})
	return eb
}

func GetEventBus() EventBus {
	return eb
}
