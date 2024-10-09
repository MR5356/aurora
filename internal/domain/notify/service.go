package notify

import (
	"context"
	database2 "github.com/MR5356/aurora/internal/infrastructure/database"
	"github.com/MR5356/aurora/internal/infrastructure/eventbus"
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
	msgTemplateDB *database2.BaseMapper[*MessageTemplate]
}

func GetService() *Service {
	once.Do(func() {
		service = &Service{
			msgTemplateDB: database2.NewMapper(database2.GetDB(), &MessageTemplate{}),
		}
	})
	return service
}

func (s *Service) sendMessage(msg *MessageTemplate) error {
	logrus.Infof("send message: %+v", msg)
	return GetNotifierManager().GetNotifier(msg.Receivers.Type).Send(context.Background(), msg)
}

func (s *Service) Initialize() error {
	if err := database2.GetDB().AutoMigrate(&MessageTemplate{}); err != nil {
		return err
	}

	if err := eventbus.GetEventBus().Subscribe(TopicSendMessage, s.sendMessage); err != nil {
		return err
	}
	return nil
}
