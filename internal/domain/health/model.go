package health

import (
	"github.com/MR5356/aurora/internal/infrastructure/database"
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

type Statistics struct {
	Total     int64                 `json:"total"`
	Up        int64                 `json:"up"`
	Down      int64                 `json:"down"`
	Unknown   int64                 `json:"unknown"`
	Error     int64                 `json:"error"`
	Ping      int64                 `json:"ping"`
	SSH       int64                 `json:"ssh"`
	HTTP      int64                 `json:"http"`
	Database  int64                 `json:"database"`
	ErrorList []*HealthListResponse `json:"errorList"`
	SlowList  []*HealthListResponse `json:"slowList"` // rtt > 460
}

type Count struct {
	Total   int64 `json:"total"`
	Up      int64 `json:"up"`
	Down    int64 `json:"down"`
	Unknown int64 `json:"unknown"`
	Error   int64 `json:"error"`
}

type HealthListResponse struct {
	ID      uuid.UUID `json:"id" gorm:"type:uuid;primaryKey" swaggerignore:"true" example:"00000000-0000-0000-0000-000000000000"`
	Title   string    `json:"title" gorm:"not null" validate:"required"`
	Desc    string    `json:"desc"`
	Type    string    `json:"type" gorm:"length:32" validate:"oneof=ping ssh http database"`
	Enabled bool      `json:"enabled"`
	Status  string    `json:"status"` // last result
	RTT     int64     `json:"rtt"`    // last result
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
