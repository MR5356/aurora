package pipeline

import (
	"database/sql/driver"
	"encoding/json"
	"github.com/MR5356/aurora/pkg/middleware/database"
	"github.com/MR5356/go-workflow"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"time"
)

type Edges struct {
	ID         string `json:"id" gorm:"primary_key;" example:"00000000-0000-0000-0000-000000000000"`
	WorkflowId string `json:"workflowId" example:"00000000-0000-0000-0000-000000000000"`
	Source     string `json:"source" example:"00000000-0000-0000-0000-000000000000"`
	Target     string `json:"target" example:"00000000-0000-0000-0000-000000000000"`

	database.BaseModel
}

func (e *Edges) TableName() string {
	return "workflow_edge"
}

func (e *Edges) BeforeCreate(tx *gorm.DB) error {
	if len(e.ID) == 0 {
		e.ID = uuid.NewString()
	}
	return nil
}

type Nodes struct {
	ID         string `json:"id" gorm:"primary_key;" example:"00000000-0000-0000-0000-000000000000"`
	WorkflowId string `json:"workflowId" example:"00000000-0000-0000-0000-000000000000"`
	Uses       string `json:"uses" example:"test"`
	Label      string `json:"label" example:"test"`
	Params     Params `json:"params" example:"test"`

	CreatedAt time.Time `json:"createdAt" swaggerignore:"true"`
	UpdatedAt time.Time `json:"updatedAt" swaggerignore:"true"`
}

type Params []*workflow.TaskParam

func (p *Params) Scan(value interface{}) error {
	s := value.(string)
	return json.Unmarshal([]byte(s), p)
}

func (p Params) Value() (driver.Value, error) {
	s, err := json.Marshal(p)
	return string(s), err
}

func (n *Nodes) TableName() string {
	return "workflow_node"
}

func (n *Nodes) BeforeCreate(tx *gorm.DB) error {
	if len(n.ID) == 0 {
		n.ID = uuid.NewString()
	}
	return nil
}

type Workflow struct {
	ID    string `json:"id" gorm:"primary_key;" example:"00000000-0000-0000-0000-000000000000"`
	Title string `json:"title" example:"test" validate:"required"`
	Owner string `json:"owner" example:"test" validate:"required"`

	database.BaseModel
}

func (w *Workflow) TableName() string {
	return "workflow"
}

func (w *Workflow) BeforeCreate(tx *gorm.DB) error {
	if len(w.ID) == 0 {
		w.ID = uuid.NewString()
	}
	return nil
}

type WorkflowRequest struct {
	*Workflow
	*workflow.WorkflowDAG
}
