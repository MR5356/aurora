package notify

import (
	"context"
	"github.com/MR5356/aurora/pkg/config"
	"sync"
)

const (
	TypeEmail = "email"
)

var (
	notifierManager     *NotifierManager
	onceNotifierManager sync.Once
)

type Notifier interface {
	Send(ctx context.Context, msg *MessageTemplate) error
}

type NotifierManager struct {
	notifiers map[string]Notifier
}

func GetNotifierManager() *NotifierManager {
	onceNotifierManager.Do(func() {
		notifierManager = &NotifierManager{
			notifiers: make(map[string]Notifier),
		}
		notifierManager.AddNotifier(TypeEmail, NewEmailNotifier(config.Current().Email))
	})
	return notifierManager
}

func (nm *NotifierManager) AddNotifier(name string, notifier Notifier) {
	nm.notifiers[name] = notifier
}

func (nm *NotifierManager) GetNotifier(name string) Notifier {
	return nm.notifiers[name]
}
