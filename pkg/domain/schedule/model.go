package schedule

import (
	"github.com/MR5356/aurora/pkg/middleware/database"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"time"
)

type Schedule struct {
	ID         uuid.UUID `json:"id" gorm:"primary_key;type:uuid;" example:"00000000-0000-0000-0000-000000000000"`
	Title      string    `json:"title" validate:"required"`
	Desc       string    `json:"desc"`
	CronString string    `json:"cronString" validate:"required" example:"*/5 * * * * *"`
	Executor   string    `json:"executor" validate:"required" example:"test"`
	NextTime   time.Time `json:"nextTime" gorm:"-" swaggerignore:"true"`
	Params     string    `json:"params"`
	Enabled    bool      `json:"enabled" example:"true"`
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
	ID         uuid.UUID `json:"id" gorm:"primary_key;type:uuid;" swaggerignore:"true"`
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
	task        func() Task
}
