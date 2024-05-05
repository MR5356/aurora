package notify

import (
	"context"
	"github.com/MR5356/aurora/pkg/middleware/database"
	"github.com/MR5356/aurora/pkg/middleware/eventbus"
	"github.com/sirupsen/logrus"
	"sync"
)

const (
	TopicSendMessage = "topic.notify.send_message"
)

var (
	once    sync.Once
	service *Service
)

type Service struct {
	msgTemplateDB *database.BaseMapper[*MessageTemplate]
}

func GetService() *Service {
	once.Do(func() {
		service = &Service{
			msgTemplateDB: database.NewMapper(database.GetDB(), &MessageTemplate{}),
		}
	})
	return service
}

func (s *Service) sendMessage(msg *MessageTemplate) error {
	logrus.Infof("send message: %+v", msg)
	return GetNotifierManager().GetNotifier(msg.Receivers.Type).Send(context.Background(), msg)
}

func (s *Service) Initialize() error {
	if err := database.GetDB().AutoMigrate(&MessageTemplate{}); err != nil {
		return err
	}

	if err := eventbus.GetEventBus().Subscribe(TopicSendMessage, s.sendMessage); err != nil {
		return err
	}
	return nil
}
