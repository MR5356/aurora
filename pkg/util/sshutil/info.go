package sshutil

import (
	"database/sql/driver"
	"encoding/json"
)

// Stats 基础信息
type Stats struct {
	MemInfo *MemInfo   `json:"memInfo"`
	CpuInfo *CpuInfo   `json:"cpuInfo"`
	NetInfo []*NetInfo `json:"netInfo"`
	Uptime  int64      `json:"uptime"`
	Ping    int64      `json:"ping"`
}

type HostInfo struct {
	Host     string `json:"host"`
	Port     uint16 `json:"port"`
	Username string `json:"username"`
	Password string `json:"password"`
}

func (h *HostInfo) Scan(val interface{}) error {
	s := val.(string)
	err := json.Unmarshal([]byte(s), &h)
	return err
}

func (h HostInfo) Value() (driver.Value, error) {
	s, err := json.Marshal(h)
	return string(s), err
}

func (s *Stats) Default() {
	s.MemInfo = &MemInfo{}
	s.MemInfo.Default()
	s.CpuInfo = &CpuInfo{}
	s.CpuInfo.Default()
	s.Ping = 9999
	s.ProcessNil()
}

func (s *Stats) ProcessNil() {
	if s.NetInfo == nil {
		s.NetInfo = make([]*NetInfo, 0)
		s.NetInfo = append(s.NetInfo, &NetInfo{
			Name:      "",
			BytesRec:  0,
			BytesSent: 0,
			SpeedRec:  0,
			SpeedSent: 0,
		})
		s.NetInfo[0].Default()
	}
}

// MemInfo 内存信息
type MemInfo struct {
	Percent   float64 `json:"percent"`
	Total     uint64  `json:"total"`
	Free      uint64  `json:"free"`
	Available uint64  `json:"available"`
	Buffers   uint64  `json:"buffers"`
	Cached    uint64  `json:"cached"`
	SwapTotal uint64  `json:"swapTotal"`
	SwapFree  uint64  `json:"swapFree"`
}

func (i *MemInfo) Default() {
	i.SwapFree = 0
	i.SwapTotal = 0
	i.Percent = 0
	i.Total = 0
	i.Free = 0
	i.Available = 0
	i.Buffers = 0
	i.Cached = 0
}

// CpuInfo cpu信息
type CpuInfo struct {
	Percent float64 `json:"percent"`
	User    float64 `json:"user"`
	Nice    float64 `json:"nice"`
	System  float64 `json:"system"`
	Idle    float64 `json:"idle"`
	IOWait  float64 `json:"iowait"`
	Irq     float64 `json:"irq"`
	SoftIrq float64 `json:"softIrq"`
	Steal   float64 `json:"steal"`
	Guest   float64 `json:"guest"`
	Cores   uint    `json:"cores"`

	total float64
}

func (i *CpuInfo) Default() {
	i.Percent = 0
	i.User = 0
	i.Nice = 0
	i.System = 0
	i.Idle = 0
	i.IOWait = 0
	i.Irq = 0
	i.SoftIrq = 0
	i.Steal = 0
	i.Guest = 0
	i.Cores = 0
}

// NetInfo 网络信息
type NetInfo struct {
	Name      string `json:"name"`
	BytesRec  uint64 `json:"bytesRec"`
	BytesSent uint64 `json:"bytesSent"`
	SpeedRec  uint64 `json:"speedRec"`
	SpeedSent uint64 `json:"speedSent"`
}

func (i *NetInfo) Default() {
	i.Name = "Unknown"
	i.BytesRec = 0
	i.BytesSent = 0
	i.SpeedRec = 0
	i.SpeedSent = 0
}
