package host

import (
	"database/sql/driver"
	"encoding/json"
	"time"

	"github.com/MR5356/aurora/pkg/middleware/database"
	"github.com/MR5356/aurora/pkg/util/sshutil"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Host struct {
	ID       uuid.UUID        `json:"id" gorm:"type:uuid;primaryKey" swaggerignore:"true" example:"00000000-0000-0000-0000-000000000000"`
	Title    string           `json:"title" gorm:"not null" validate:"required"`
	Desc     string           `json:"desc"`
	HostInfo sshutil.HostInfo `json:"hostInfo" gorm:"not null"`
	MetaInfo MetaInfo         `json:"metaInfo" swaggerignore:"true"`
	Group    Group            `json:"group" swaggerignore:"true" validate:"omitempty"`
	GroupId  uuid.UUID        `json:"groupId"`

	database.BaseModel
}

type Group struct {
	ID    uuid.UUID `json:"id" gorm:"type:uuid;primaryKey" swaggerignore:"true" example:"00000000-0000-0000-0000-000000000000"`
	Title string    `json:"title" gorm:"unique;not null" validate:"required"`
	Hosts []*Host   `json:"hosts" gorm:"foreignkey:GroupId" swaggerignore:"true"`

	CreatedAt time.Time `json:"createdAt" swaggerignore:"true"`
	UpdatedAt time.Time `json:"updatedAt" swaggerignore:"true"`
}

func (g *Group) TableName() string {
	return "host_group"
}

func (g *Group) BeforeCreate(tx *gorm.DB) error {
	if g.ID == uuid.Nil {
		g.ID = uuid.New()
	}
	return nil
}

type MetaInfo struct {
	OS         string `json:"os"`
	Kernel     string `json:"kernel"`
	Hostname   string `json:"hostname"`
	Arch       string `json:"arch"`
	Cpu        string `json:"cpu"`
	Mem        string `json:"mem"`
	Containerd string `json:"containerd"`
	Docker     string `json:"docker"`
}

func (m *Host) TableName() string {
	return "host"
}

func (m *Host) BeforeCreate(tx *gorm.DB) error {
	if m.ID == uuid.Nil {
		m.ID = uuid.New()
	}
	return nil
}

func (m *MetaInfo) Scan(val interface{}) error {
	s := val.(string)
	err := json.Unmarshal([]byte(s), &m)
	return err
}

func (m MetaInfo) Value() (driver.Value, error) {
	s, err := json.Marshal(m)
	return string(s), err
}
