package host

import (
	"context"
	"encoding/json"
	"fmt"
	database2 "github.com/MR5356/aurora/internal/infrastructure/database"
	"github.com/MR5356/aurora/pkg/util/sshutil"
	"sync"
	"time"

	"github.com/MR5356/aurora/pkg/util/cacheutil"
	"github.com/MR5356/aurora/pkg/util/container"
	"github.com/MR5356/aurora/pkg/util/structutil"
	"github.com/MR5356/aurora/pkg/util/validate"
	"github.com/MR5356/jietan/pkg/executor"
	"github.com/MR5356/jietan/pkg/executor/api"
	"github.com/google/uuid"
	"golang.org/x/crypto/ssh"
)

var (
	onceService sync.Once
	service     *Service
)

type Service struct {
	hostDb               database2.Mapper[*Host]
	groupDb              *database2.BaseMapper[*Group]
	containerClientCache *cacheutil.CountdownCache[container.Client]
	hostClientCache      *cacheutil.CountdownCache[*sshutil.Client]

	statsCache *StatsCache
}

func GetService() *Service {
	onceService.Do(func() {
		service = &Service{
			hostDb:               database2.NewMapper(database2.GetDB(), &Host{}),
			groupDb:              database2.NewMapper(database2.GetDB(), &Group{}),
			containerClientCache: cacheutil.NewCountdownCache[container.Client](time.Minute * 30),
			hostClientCache:      cacheutil.NewCountdownCache[*sshutil.Client](time.Minute * 30),
			statsCache:           &StatsCache{},
		}
	})
	return service
}

// ListGroup list host group
func (s *Service) ListGroup(group *Group) ([]*Group, error) {
	res := make([]*Group, 0)
	if err := s.groupDb.DB.Preload("Hosts").Find(&res, group).Error; err != nil {
		return res, err
	}

	for _, g := range res {
		for _, h := range g.Hosts {
			h.HostInfo.Password = ""
		}
	}

	return res, nil
}

// AddGroup add host group
func (s *Service) AddGroup(group *Group) error {
	group.ID = uuid.Nil
	if err := validate.Validate(group); err != nil {
		return err
	}
	return s.groupDb.Insert(group)
}

// DeleteGroup delete host group
func (s *Service) DeleteGroup(id uuid.UUID) error {
	if err := s.groupDb.DB.Where("group_id = ?", id).Delete(&Host{GroupId: id}).Error; err != nil {
		return err
	}

	return s.groupDb.Delete(&Group{ID: id})
}

// UpdateGroup update host group
func (s *Service) UpdateGroup(group *Group) error {
	if err := validate.Validate(group); err != nil {
		return err
	}
	return s.groupDb.DB.Updates(group).Error
}

// AddHost add host
func (s *Service) AddHost(host *Host) error {
	host.ID = uuid.Nil
	if err := validate.Validate(host); err != nil {
		return err
	}

	if err := s.checkHost(host); err != nil {
		return err
	}

	return s.hostDb.Insert(host)
}

// UpdateHost update host
func (s *Service) UpdateHost(host *Host) error {
	if err := validate.Validate(host); err != nil {
		return err
	}

	if err := s.checkHost(host); err != nil {
		return err
	}

	return s.hostDb.Update(&Host{ID: host.ID}, structutil.Struct2Map(host))
}

// DeleteHost delete host
func (s *Service) DeleteHost(id uuid.UUID) error {
	return s.hostDb.Delete(&Host{ID: id})
}

// DetailHost detail host
func (s *Service) DetailHost(id uuid.UUID) (*Host, error) {
	return s.hostDb.Detail(&Host{ID: id})
}

// ListHost list host
func (s *Service) ListHost(host *Host) ([]*Host, error) {
	res := make([]*Host, 0)
	if err := s.hostDb.GetDB().Joins("Group").Find(&res, host).Error; err != nil {
		return res, err
	}

	for _, h := range res {
		h.HostInfo.Password = ""
	}

	return res, nil
}

func (s *Service) checkHost(machine *Host) error {
	sshClient, err := ssh.Dial("tcp", fmt.Sprintf("%s:%d", machine.HostInfo.Host, machine.HostInfo.Port), &ssh.ClientConfig{
		User:            machine.HostInfo.Username,
		Auth:            machine.HostInfo.GetAuthMethods(),
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
		Timeout:         time.Second * 20,
	})
	if err != nil {
		return err
	}
	defer sshClient.Close()
	exec := executor.GetExecutor("remote")
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*20)
	defer cancel()
	hostInfo := &api.HostInfo{
		Host:   machine.HostInfo.Host,
		Port:   int(machine.HostInfo.Port),
		User:   machine.HostInfo.Username,
		Passwd: machine.HostInfo.Password,
	}
	a := exec.Execute(ctx, &api.ExecuteParams{
		"hosts": []*api.HostInfo{
			hostInfo,
		},
		"script": `#!/bin/bash

# 系统信息
os=$(cat /etc/os-release | grep ^ID= | cut -d '=' -f2)
os=${os//\"/} # 去掉双引号
kernel=$(uname -r)
hostname=$(hostname)
arch=$(uname -m)

# 硬件信息
cpu_count=$(lscpu | grep "^CPU(s):" | cut -d ':' -f2 | awk '{$1=$1;print}')
mem_size=$(free -h | grep Mem | awk '{print $2}')

# containerd 版本信息
containerd_version=$(containerd --version 2>/dev/null || k3s crictl version 2>/dev/null | grep 'RuntimeVersion' | awk '{print $2}')

# docker 版本信息
docker_version=$(docker version --format '{{.Server.APIVersion}}' 2>/dev/null)

# 构建JSON
json=$(cat <<EOF
{
    "os": "$os",
    "kernel": "$kernel",
    "hostname": "$hostname",
    "arch": "$arch",
    "cpu": "$cpu_count",
    "mem": "$mem_size",
    "containerd": "$containerd_version",
    "docker": "$docker_version"
}
EOF
)

# 输出
echo $json`,
		"params": "",
	})
	metaInfo := new(MetaInfo)
	err = json.Unmarshal([]byte(a.Data["log"].(map[string][]string)[hostInfo.String()][0]), metaInfo)
	if err != nil {
		metaInfo = &MetaInfo{
			OS:       "unknown",
			Kernel:   "unknown",
			Hostname: "unknown",
			Arch:     "unknown",
			Cpu:      "unknown",
			Mem:      "unknown",
		}
	}

	// TODO 增加系统信息检测
	machine.MetaInfo = *metaInfo
	return nil
}

func (s *Service) Initialize() error {
	if err := database2.GetDB().AutoMigrate(&Host{}, &Group{}); err != nil {
		return err
	}

	if err := s.groupDb.DB.Where(&Group{ID: uuid.MustParse("b0ea5261-4185-44f3-b16b-ef7e6b681775")}).Attrs(&Group{Title: "default"}).FirstOrCreate(&Group{}).Error; err != nil {
		return err
	}

	return nil
}
