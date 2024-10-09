package notify

import (
	"github.com/MR5356/aurora/internal/infrastructure/database"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type MessageTemplate struct {
	ID             uuid.UUID `json:"id" gorm:"primary_key;type:uuid;"`
	Event          string    `json:"event" gorm:"not null;"`
	Subject        string    `json:"subject" gorm:"not null;"`
	Body           string    `json:"body" gorm:"not null;"`
	Level          string    `json:"level" gorm:"not null;"`
	DefaultSubject string    `json:"defaultSubject" gorm:"not null;"`
	DefaultBody    string    `json:"defaultBody" gorm:"not null;"`

	Receivers MessageReceiver `json:"-" gorm:"-"`

	database.BaseModel
}

type MessageReceiver struct {
	Receivers []string `json:"-" gorm:"-"`
	Type      string   `json:"-" gorm:"-"`
}

func (m *MessageTemplate) TableName() string {
	return "notify_message"
}

func (m *MessageTemplate) BeforeCreate(tx *gorm.DB) error {
	m.ID = uuid.New()
	return nil
}
