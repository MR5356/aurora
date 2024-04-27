package eventbus

import (
	"errors"
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

func GetEventBus() EventBus {
	// var err error
	once.Do(func() {
		eb = NewMemoryEventBus()
	})
	return eb
}
