package pipeline

import (
	"github.com/MR5356/aurora/pkg/middleware/database"
	"github.com/MR5356/aurora/pkg/util/structutil"
	"github.com/MR5356/aurora/pkg/util/validate"
	"github.com/MR5356/go-workflow"
	"github.com/sirupsen/logrus"
	"sync"
)

var (
	onceService sync.Once
	service     *Service
)

type Service struct {
	wfDB *database.BaseMapper[*Workflow]
	nDB  *database.BaseMapper[*Nodes]
	eDB  *database.BaseMapper[*Edges]
}

func GetService() *Service {
	onceService.Do(func() {
		service = &Service{
			wfDB: database.NewMapper(database.GetDB(), &Workflow{}),
			nDB:  database.NewMapper(database.GetDB(), &Nodes{}),
			eDB:  database.NewMapper(database.GetDB(), &Edges{}),
		}
	})
	return service
}

func (s *Service) AddWorkflow(wf *WorkflowRequest) error {
	if err := validate.Validate(wf); err != nil {
		logrus.Errorf("validate workflow failed, error: %v", err)
		return err
	}

	tx := database.GetDB().Begin()
	defer tx.Rollback()

	if err := s.wfDB.Insert(wf.Workflow, tx); err != nil {
		logrus.Errorf("insert workflow failed, error: %v", err)
		return err
	}

	for _, n := range wf.Nodes {
		node := &Nodes{
			ID:         n.Id,
			WorkflowId: wf.Workflow.ID,
			Uses:       n.Uses,
			Label:      n.Label,
			Params:     n.Params,
		}
		if err := s.nDB.Insert(node, tx); err != nil {
			logrus.Errorf("insert node failed, error: %v", err)
			return err
		}
	}

	for _, e := range wf.Edges {
		edge := &Edges{
			WorkflowId: wf.Workflow.ID,
			Source:     e.Source,
			Target:     e.Target,
		}
		if err := s.eDB.Insert(edge, tx); err != nil {
			logrus.Errorf("insert edge failed, error: %v", err)
			return err
		}
	}

	tx.Commit()

	return nil
}

func (s *Service) UpdateWorkflow(wf *WorkflowRequest) error {
	if err := validate.Validate(wf); err != nil {
		logrus.Errorf("validate workflow failed, error: %v", err)
		return err
	}

	tx := database.GetDB().Begin()
	defer tx.Rollback()

	if err := s.wfDB.Update(&Workflow{ID: wf.Workflow.ID}, structutil.Struct2Map(wf.Workflow), tx); err != nil {
		logrus.Errorf("update workflow failed, error: %v", err)
		return err
	}

	if err := tx.Where("workflow_id = ?", wf.Workflow.ID).Delete(&Nodes{}).Error; err != nil {
		logrus.Errorf("delete node failed, error: %v", err)
		return err
	}
	for _, n := range wf.Nodes {
		node := &Nodes{
			ID:         n.Id,
			WorkflowId: wf.Workflow.ID,
			Uses:       n.Uses,
			Label:      n.Label,
			Params:     n.Params,
		}
		if err := s.nDB.Insert(node, tx); err != nil {
			logrus.Errorf("insert node failed, error: %v", err)
			return err
		}
	}

	if err := tx.Where("workflow_id = ?", wf.Workflow.ID).Delete(&Edges{}).Error; err != nil {
		logrus.Errorf("delete edge failed, error: %v", err)
		return err
	}
	for _, e := range wf.Edges {
		edge := &Edges{
			WorkflowId: wf.Workflow.ID,
			Source:     e.Source,
			Target:     e.Target,
		}
		if err := s.eDB.Insert(edge, tx); err != nil {
			logrus.Errorf("insert edge failed, error: %v", err)
			return err
		}
	}
	tx.Commit()
	return nil
}

func (s *Service) GetWorkflow(wf *Workflow) (*WorkflowRequest, error) {
	wf, err := s.wfDB.Detail(&Workflow{ID: wf.ID})
	if err != nil {
		return nil, err
	}
	nodes, err := s.nDB.List(&Nodes{WorkflowId: wf.ID})
	if err != nil {
		return nil, err
	}
	node := make([]*workflow.Node, 0)
	for _, n := range nodes {
		node = append(node, &workflow.Node{
			Id:     n.ID,
			Label:  n.Label,
			Uses:   n.Uses,
			Params: n.Params,
			Status: workflow.NodeStatusPending,
		})
	}
	edges, err := s.eDB.List(&Edges{WorkflowId: wf.ID})
	if err != nil {
		return nil, err
	}
	edge := make([]*workflow.Edge, 0)
	for _, e := range edges {
		edge = append(edge, &workflow.Edge{
			Source: e.Source,
			Target: e.Target,
			Status: workflow.NodeStatusPending,
		})
	}
	wfr := &WorkflowRequest{
		&Workflow{},
		&workflow.WorkflowDAG{},
	}
	wfr.ID = wf.ID
	wfr.Title = wf.Title
	wfr.Owner = wf.Owner
	wfr.CreatedAt = wf.CreatedAt
	wfr.UpdatedAt = wf.UpdatedAt
	wfr.Edges = edge
	wfr.Nodes = node
	return wfr, nil
}

func (s *Service) ListWorkflow(wf *Workflow) ([]*WorkflowRequest, error) {
	res := make([]*WorkflowRequest, 0)
	wfs, err := s.wfDB.List(wf)
	if err != nil {
		return nil, err
	}
	for _, w := range wfs {
		wfr, err := s.GetWorkflow(w)
		if err != nil {
			return nil, err
		}
		res = append(res, wfr)
	}
	return res, nil
}

func (s *Service) Initialize() error {
	if err := database.GetDB().AutoMigrate(&Workflow{}, &Nodes{}, &Edges{}); err != nil {
		return err
	}
	return nil
}
