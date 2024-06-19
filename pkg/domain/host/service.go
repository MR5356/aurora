package host

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/MR5356/aurora/pkg/middleware/database"
	"github.com/MR5356/aurora/pkg/util/structutil"
	"github.com/MR5356/aurora/pkg/util/validate"
	"github.com/MR5356/jietan/pkg/executor"
	"github.com/MR5356/jietan/pkg/executor/api"
	"github.com/google/uuid"
	"golang.org/x/crypto/ssh"
	"sync"
	"time"
)

var (
	onceService sync.Once
	service     *Service
)

type Service struct {
	hostDb  *database.BaseMapper[*Host]
	groupDb *database.BaseMapper[*Group]
}

func GetService() *Service {
	onceService.Do(func() {
		service = &Service{
			hostDb:  database.NewMapper(database.GetDB(), &Host{}),
			groupDb: database.NewMapper(database.GetDB(), &Group{}),
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
	if err := s.hostDb.DB.Joins("Group").Find(&res, host).Error; err != nil {
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
		Auth:            []ssh.AuthMethod{ssh.Password(machine.HostInfo.Password)},
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
		"script": "#!/bin/bash\n\n# 系统信息\nos=$(cat /etc/os-release | grep ^ID= | cut -d '=' -f2)\nos=${os//\\\"/} # 去掉双引号\nkernel=$(uname -r)\nhostname=$(hostname)\narch=$(uname -m)\n\n# 硬件信息\ncpu_count=$(lscpu | grep \"^CPU:\\|^CPU(s)\" | cut -d ':' -f2 | awk '{$1=$1;print}')\nmem_size=$(free -h | grep Mem | awk '{print $2}')\n\n# 构建JSON\njson=\"{\\\"os\\\": \\\"$os\\\",\n        \\\"kernel\\\": \\\"$kernel\\\",\n        \\\"hostname\\\": \\\"$hostname\\\",\n        \\\"arch\\\": \\\"$arch\\\",\n        \\\"cpu\\\": \\\"$cpu_count\\\",\n        \\\"mem\\\": \\\"$mem_size\\\"}\"\n\n# 输出\necho $json",
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
	if err := database.GetDB().AutoMigrate(&Host{}, &Group{}); err != nil {
		return err
	}

	if err := s.groupDb.DB.Where(&Group{ID: uuid.MustParse("b0ea5261-4185-44f3-b16b-ef7e6b681775")}).Attrs(&Group{Title: "default"}).FirstOrCreate(&Group{}).Error; err != nil {
		return err
	}
	return nil
}
