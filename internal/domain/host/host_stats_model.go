package host

import "sync"

type StatsCache map[string]*Stats

func (sc *StatsCache) Set(key string, value any) {
	stats, ok := (*sc)[key]
	if !ok {
		(*sc)[key] = (&Stats{}).Default()
		stats = (*sc)[key]
	}
	stats.mutex.Lock()
	defer stats.mutex.Unlock()
	switch value.(type) {
	case *CPUInfo:
		stats.CpuInfo = value.(*CPUInfo)
	case *MemInfo:
		stats.MemInfo = value.(*MemInfo)
	case *DiskInfo:
		stats.DiskInfo = value.(*DiskInfo)
	case *NetworkInfo:
		stats.NetworkInfo = value.(*NetworkInfo)
	case *ProcessInfo:
		stats.ProcessInfo = value.(*ProcessInfo)
	case RTT:
		stats.Rtt = value.(RTT)
	}
}

func (sc *StatsCache) Get(key string) *Stats {
	if stats, ok := (*sc)[key]; ok {
		return stats
	} else {
		return (&Stats{}).Default()
	}
}

type Stats struct {
	mutex       sync.Mutex
	CpuInfo     *CPUInfo     `json:"cpuInfo"`
	MemInfo     *MemInfo     `json:"memInfo"`
	DiskInfo    *DiskInfo    `json:"diskInfo"`
	NetworkInfo *NetworkInfo `json:"networkInfo"`
	ProcessInfo *ProcessInfo `json:"processInfo"`
	Rtt         RTT          `json:"rtt"`
}

func (s *Stats) Default() *Stats {
	return &Stats{
		CpuInfo:     &CPUInfo{},
		MemInfo:     &MemInfo{},
		DiskInfo:    &DiskInfo{},
		NetworkInfo: &NetworkInfo{},
		ProcessInfo: &ProcessInfo{},
		Rtt:         0,
	}
}

type RTT int64

type CPUInfo struct {
	Percent  float64   `json:"percent"`
	System   float64   `json:"system"`
	User     float64   `json:"user"`
	IOWait   float64   `json:"iowait"`
	Steal    float64   `json:"steal"`
	Idle     float64   `json:"idle"`
	Children []CPUInfo `json:"children"`
	Uptime   int64     `json:"uptime"`
	Load     string    `json:"load"`
}

type MemInfo struct {
	Total     int64 `json:"total"`
	Used      int64 `json:"used"`
	Free      int64 `json:"free"`
	Shared    int   `json:"shared"`
	BuffCache int64 `json:"buffcache"`
	Available int64 `json:"available"`
}

type DiskInfo []DiskItem
type DiskItem struct {
	Partition       string `json:"partition"`
	SizeBytes       int    `json:"size_bytes"`
	UsedSizeBytes   int    `json:"used_size_bytes"`
	MountPoint      string `json:"mount_point"`
	FsType          string `json:"fs_type"`
	ReadSpeedBps    int    `json:"read_speed_Bps"`
	WriteSpeedBps   int    `json:"write_speed_Bps"`
	ReadIops        int    `json:"read_iops"`
	WriteIops       int    `json:"write_iops"`
	ReadLatencyMs   int    `json:"read_latency_ms"`
	WriteLatencyMs  int    `json:"write_latency_ms"`
	TotalReadBytes  int    `json:"total_read_bytes"`
	TotalWriteBytes int    `json:"total_write_bytes"`
}

type NetworkInfo []NetworkItem
type NetworkItem struct {
	Name      string `json:"name"`
	Type      string `json:"type"`
	BytesRec  int    `json:"bytes_rec"`
	BytesSent int    `json:"bytes_sent"`
	SpeedRec  int64  `json:"speed_rec"`
	SpeedSent int    `json:"speed_sent"`
}

type ProcessInfo []ProcessItem
type ProcessItem struct {
	CPU     float64 `json:"cpu"`
	Mem     float64 `json:"mem"`
	Command string  `json:"command"`
}
