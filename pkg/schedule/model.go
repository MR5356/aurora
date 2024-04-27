package schedule

import (
	"github.com/MR5356/aurora/pkg/middleware/database"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"time"
)

type Schedule struct {
	ID         uuid.UUID `json:"id" gorm:"primary_key;type:uuid;"`
	Title      string    `json:"title" validate:"required"`
	Desc       string    `json:"desc"`
	CronString string    `json:"cronString" validate:"required"`
	TaskName   string    `json:"taskName" validate:"required"`
	NextTime   time.Time `json:"nextTime" gorm:"-"`
	Params     string    `json:"params"`
	Enabled    bool      `json:"enabled"`
	Status     string    `json:"status"`

	database.BaseModel
}

func (s *Schedule) TableName() string {
	return "schedule"
}

func (s *Schedule) BeforeCreate(tx *gorm.DB) error {
	s.ID = uuid.New()
	return nil
}

type Record struct {
	ID         uuid.UUID `json:"id" gorm:"primary_key;type:uuid;"`
	ScheduleID uuid.UUID `json:"scheduleID" gorm:"type:uuid;"`
	Title      string    `json:"title" validate:"required"`
	TaskName   string    `json:"taskName" validate:"required"`
	Params     string    `json:"params"`
	Status     string    `json:"status"`

	database.BaseModel
}

func (r *Record) TableName() string {
	return "schedule_record"
}

func (r *Record) BeforeCreate(tx *gorm.DB) error {
	r.ID = uuid.New()
	return nil
}

type Executor struct {
	Name        string `json:"name"`
	DisplayName string `json:"displayName"`
}
