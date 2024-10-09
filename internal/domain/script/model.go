package script

import (
	"database/sql/driver"
	"encoding/json"
	"github.com/MR5356/aurora/internal/infrastructure/database"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Script struct {
	ID      uuid.UUID `json:"id" gorm:"primary_key;type:uuid;" swaggerignore:"true"`
	Title   string    `json:"title" gorm:"not null" validate:"required"`
	Desc    string    `json:"desc"`
	Content string    `json:"content" gorm:"not null" validate:"required"`
	Type    string    `json:"type" validate:"oneof=shell python"`

	database.BaseModel
}

func (s *Script) TableName() string {
	return "script"
}

func (s *Script) BeforeCreate(tx *gorm.DB) error {
	if s.ID == uuid.Nil {
		s.ID = uuid.New()
	}
	return nil
}

type Params map[string]any

func (p *Params) Scan(val interface{}) error {
	s := val.(string)
	return json.Unmarshal([]byte(s), p)
}

func (p Params) Value() (driver.Value, error) {
	s, err := json.Marshal(p)
	return string(s), err
}

type Record struct {
	ID          uuid.UUID `json:"id" gorm:"primary_key;type:uuid;" swaggerignore:"true"`
	ScriptTitle string    `json:"scriptTitle"`
	Script      string    `json:"script"`
	Hosts       string    `json:"hosts"`
	Params      string    `json:"params"`
	Result      string    `json:"result"`
	Status      string    `json:"status"`
	Message     string    `json:"message"`
	Error       string    `json:"error"`

	database.BaseModel
}

func (r *Record) TableName() string {
	return "script_record"
}

func (r *Record) BeforeCreate(tx *gorm.DB) error {
	if r.ID == uuid.Nil {
		r.ID = uuid.New()
	}
	return nil
}

type RunScriptParams struct {
	ScriptId uuid.UUID
	HostIds  []uuid.UUID
	Params   string
}
