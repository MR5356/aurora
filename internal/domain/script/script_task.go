package script

import (
	"context"
	"encoding/json"
	"github.com/MR5356/aurora/internal/domain/host"
	database2 "github.com/MR5356/aurora/internal/infrastructure/database"
	"github.com/MR5356/jietan/pkg/executor"
	"github.com/MR5356/jietan/pkg/executor/api"
	"github.com/sirupsen/logrus"
	"sync"
	"time"
)

const (
	taskStatusRunning  = "running"
	taskStatusFinished = "finished"
	taskStatusFailed   = "failed"
)

var jobMap = sync.Map{}

type JobInfo struct {
	exec      executor.Executor
	ctx       context.Context
	ctxCancel context.CancelFunc
}

type Task struct {
	params   *RunScriptParams
	recordDB *database2.BaseMapper[*Record]
	scriptDB *database2.BaseMapper[*Script]
	hostDB   *database2.BaseMapper[*host.Host]

	finished bool
}

func NewTask() *Task {
	return &Task{
		recordDB: database2.NewMapper(database2.GetDB(), &Record{}),
		scriptDB: database2.NewMapper(database2.GetDB(), &Script{}),
		hostDB:   database2.NewMapper(database2.GetDB(), &host.Host{}),
	}
}

type TaskParams struct {
	Script *Script
	Hosts  []*api.HostInfo
	Params string
}

func (t *Task) SetParams(params string) {
	ps := new(RunScriptParams)
	if err := json.Unmarshal([]byte(params), ps); err != nil {
		logrus.Errorf("set task params error: %v", err)
		return
	}

	t.params = ps
}

func (t *Task) Run() {
	if t.params == nil {
		logrus.Errorf("task params is empty")
		return
	}
	script, err := t.scriptDB.Detail(&Script{ID: t.params.ScriptId})
	if err != nil {
		logrus.Errorf("script %s not found", t.params.ScriptId)
		return
	}
	hosts := make([]*api.HostInfo, 0)
	for _, id := range t.params.HostIds {
		h, err := t.hostDB.Detail(&host.Host{ID: id})
		if err != nil {
			logrus.Errorf("host %s not found", id)
			return
		}
		hosts = append(hosts, &api.HostInfo{
			Host:       h.HostInfo.Host,
			Port:       int(h.HostInfo.Port),
			User:       h.HostInfo.Username,
			Passwd:     h.HostInfo.Password,
			PrivateKey: h.HostInfo.PrivateKey,
			Passphrase: h.HostInfo.Passphrase,
		})
	}

	hostsStr, _ := json.Marshal(hosts)
	record := &Record{
		ScriptTitle: script.Title,
		Script:      script.Content,
		Hosts:       string(hostsStr),
		Params:      t.params.Params,
		Status:      taskStatusRunning,
	}

	if err := t.recordDB.Insert(record); err != nil {
		logrus.Errorf("add record error: %v", err)
	}

	exec := executor.GetExecutor("remote")

	ctx, cancel := context.WithCancel(context.Background())
	jobMap.Store(record.ID.String(), &JobInfo{
		exec:      exec,
		ctx:       ctx,
		ctxCancel: cancel,
	})

	defer jobMap.Delete(record.ID.String())

	// refresh log
	go func() {
		for !t.finished {
			log := exec.GetResult(api.ResultFieldLog, nil).(map[string][]string)
			logStr, _ := json.Marshal(log)
			record.Result = string(logStr)
			if err := t.recordDB.DB.Updates(record).Error; err != nil {
				logrus.Errorf("update record error: %v", err)
			}
			time.Sleep(time.Second)
		}
	}()

	res := exec.Execute(ctx, &api.ExecuteParams{
		"hosts":  hosts,
		"script": script.Content,
		"params": t.params.Params,
	})

	errs, _ := json.Marshal(res.Data["error"])
	log := res.Data["log"].(map[string][]string)
	logStr, _ := json.Marshal(log)
	record.Result = string(logStr)
	if res.Status == api.Success {
		record.Status = taskStatusFinished
	} else {
		record.Status = taskStatusFailed
	}
	record.Message = res.Message
	record.Error = string(errs)

	t.finished = true
	if err := t.recordDB.DB.Updates(record).Error; err != nil {
		logrus.Errorf("update record error: %v", err)
	}
}
