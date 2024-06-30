package health

import (
	"github.com/MR5356/aurora/pkg/middleware/database"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Health struct {
	ID      uuid.UUID `json:"id" gorm:"type:uuid;primaryKey" swaggerignore:"true" example:"00000000-0000-0000-0000-000000000000"`
	Title   string    `json:"title" gorm:"not null" validate:"required"`
	Desc    string    `json:"desc"`
	Type    string    `json:"type" gorm:"length:32" validate:"oneof=ping ssh http database"`
	Enabled bool      `json:"enabled"`
	Params  string    `json:"params" validate:"required"`
	Status  string    `json:"status"` // last result
	RTT     int64     `json:"rtt"`    // last result

	database.BaseModel
}

func (h *Health) TableName() string {
	return "health_check"
}

func (h *Health) BeforeCreate(tx *gorm.DB) error {
	if h.ID == uuid.Nil {
		h.ID = uuid.New()
	}
	return nil
}

type Param struct {
	Key      string `json:"key"`
	Value    any    `json:"value"`
	Title    string `json:"title"`
	Type     string `json:"type"`
	Required bool   `json:"required"`
	Desc     string `json:"desc"`
}

type CheckType struct {
	Title  string  `json:"title"`
	Type   string  `json:"type"`
	Desc   string  `json:"desc"`
	Params []Param `json:"params"`
}

type Params []Param

func (ps *Params) GetKey(key string) any {
	for _, param := range *ps {
		if param.Key == key {
			return param.Value
		}
	}
	return ""
}

type Record struct {
	ID       uuid.UUID `json:"id" gorm:"type:uuid;primaryKey" swaggerignore:"true" example:"00000000-0000-0000-0000-000000000000"`
	ParentId uuid.UUID `json:"parentId" gorm:"type:uuid;" swaggerignore:"true" example:"00000000-0000-0000-0000-000000000000"`

	Status string `json:"status"`
	Rtt    int64  `json:"rtt"`
	Result string `json:"result"` // json string

	database.BaseModel
}

func (h *Record) TableName() string {
	return "health_check_record"
}

func (h *Record) BeforeCreate(tx *gorm.DB) error {
	if h.ID == uuid.Nil {
		h.ID = uuid.New()
	}
	return nil
}
