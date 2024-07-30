package host

import (
	"encoding/json"
	"github.com/MR5356/aurora/pkg/util/sshutil"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"sync"
	"time"
)

const (
	getCPUInfo = `#!/bin/bash

# 读取 /proc/stat 文件内容
stat_content=$(cat /proc/stat)

# 提取 CPU 总体信息
cpu_line=$(echo "$stat_content" | grep "^cpu ")

# 提取开机时间
btime=$(echo "$stat_content" | grep "^btime" | awk '{print$2}')

read -r _ user nice system idle iowait irq softirq steal guest guest_nice <<< $cpu_line

# 计算总时间
total=$(($user + $nice + $system + $idle + $iowait + $irq + $softirq + $steal + $guest + $guest_nice))

# 计算各项百分比
user_percent=$(awk "BEGIN {print ($user / $total)}")
nice_percent=$(awk "BEGIN {print ($nice / $total)}")
system_percent=$(awk "BEGIN {print ($system / $total)}")
idle_percent=$(awk "BEGIN {print ($idle / $total)}")
iowait_percent=$(awk "BEGIN {print ($iowait / $total)}")
irq_percent=$(awk "BEGIN {print ($irq / $total)}")
softirq_percent=$(awk "BEGIN {print ($softirq / $total)}")
steal_percent=$(awk "BEGIN {print ($steal / $total)}")
guest_percent=$(awk "BEGIN {print ($guest / $total)}")
guest_nice_percent=$(awk "BEGIN {print ($guest_nice / $total)}")

# 计算总体 CPU 使用率
total_usage_percent=$(awk "BEGIN {print 1 - $idle_percent}")

# 提取每个 CPU 核心的信息
cpu_cores=$(echo "$stat_content" | grep "^cpu[0-9]" | while read -r line; do
    read -r core user nice system idle iowait irq softirq steal guest guest_nice <<< $line
    total=$(($user + $nice + $system + $idle + $iowait + $irq + $softirq + $steal + $guest + $guest_nice))
    user_percent=$(awk "BEGIN {print ($user / $total)}")
    nice_percent=$(awk "BEGIN {print ($nice / $total)}")
    system_percent=$(awk "BEGIN {print ($system / $total)}")
    idle_percent=$(awk "BEGIN {print ($idle / $total)}")
    iowait_percent=$(awk "BEGIN {print ($iowait / $total)}")
    irq_percent=$(awk "BEGIN {print ($irq / $total)}")
    softirq_percent=$(awk "BEGIN {print ($softirq / $total)}")
    steal_percent=$(awk "BEGIN {print ($steal / $total)}")
    guest_percent=$(awk "BEGIN {print ($guest / $total)}")
    guest_nice_percent=$(awk "BEGIN {print ($guest_nice / $total)}")
    total_usage_percent=$(awk "BEGIN {print 1 - $idle_percent}")
    echo "{\"percent\": $total_usage_percent, \"system\": $system_percent, \"user\": $user_percent, \"iowait\": $iowait_percent, \"steal\": $steal_percent, \"idle\": $idle_percent}"
done | paste -sd,)

# 获取系统运行时间和负载
# uptime_info=$(uptime -p | sed 's/up //')
load=$(uptime | awk -F'load average: ' '{print $2}')

# 创建 JSON 格式
json=$(cat <<EOF
{
  "percent": $total_usage_percent,
  "system": $system_percent,
  "user": $user_percent,
  "iowait": $iowait_percent,
  "steal": $steal_percent,
  "idle": $idle_percent,
  "uptime": $btime,
  "load": "$load",
  "children": [$cpu_cores]
}
EOF
)

# 输出 JSON
echo "$json"
`
	getDiskInfo = `#!/bin/bash

# 确保我们有足够权限来访问系统文件
# if [ "$(id -u)" -ne "0" ]; then
#     echo "This script needs to be run as root."
#     exit 1
# fi

# 获取所有磁盘分区
partitions=$(lsblk -ln -o NAME,TYPE,MOUNTPOINT | grep -v "^$" | grep "/" | awk '{print $1}')

# 临时文件用于存储初始和最终读写字节数
temp_file_initial=$(mktemp)
temp_file_final=$(mktemp)
temp_file_sizes=$(mktemp)
temp_file_mounts=$(mktemp)
temp_file_fstype=$(mktemp)
temp_file_usage=$(mktemp)

# 记录初始读写字节数
# echo "Collecting initial disk stats..."
for partition in $partitions; do
    grep "$partition" /proc/diskstats | awk '{print $3, $4, $8, $6, $10, $11, $12, $13}' >> $temp_file_initial
done

# 记录分区总大小、挂载点、格式和使用情况
# echo "Collecting disk sizes, mount points, and file system types..."
lsblk -ln -o NAME,MOUNTPOINT,FSTYPE | grep -E "($(echo $partitions | sed 's/ /|/g'))" | awk '{print $1, $2, $3}' >> $temp_file_mounts
df -B1 | awk 'NR>1 {print $1, $2, $3, $6}' >> $temp_file_usage

# 获取分区的格式信息
lsblk -ln -o NAME,FSTYPE | grep -E "($(echo $partitions | sed 's/ /|/g'))" | awk '{print $1, $2}' >> $temp_file_fstype

# 等待1秒钟来获取新的数据
sleep 1

# 记录最终读写字节数
# echo "Collecting final disk stats..."
for partition in $partitions; do
    grep "$partition" /proc/diskstats | awk '{print $3, $4, $8, $6, $10, $11, $12, $13}' >> $temp_file_final
done

# 输出 JSON 格式
echo "["

first=true

while IFS=' ' read -r device mount_point fstype; do
    # 从最终数据中提取相关数据
    read -r _ read_ios_final write_ios_final read_sectors_final write_sectors_final read_time_final write_time_final discard_time_final < <(grep "$device" $temp_file_final)
    
    # 从初始数据中提取相关数据
    read -r _ read_ios_initial write_ios_initial read_sectors_initial write_sectors_initial read_time_initial write_time_initial discard_time_initial < <(grep "$device" $temp_file_initial)
    
    # 从分区的格式信息文件中提取文件系统类型
    fs_type=$(grep "$device" $temp_file_fstype | awk '{print $2}')
    
    # 从已使用空间文件中提取已使用大小
    used_size_str=$(grep "$device" $temp_file_usage | awk '{print $3}')
    used_size_bytes=$used_size_str

    # 从 df 命令输出中提取总大小
    total_size_str=$(grep "$device" $temp_file_usage | awk '{print $2}')
    total_size_bytes=$total_size_str
    
    # 计算读写差异
    read_ios_diff=$((read_ios_final - read_ios_initial))
    write_ios_diff=$((write_ios_final - write_ios_initial))
    read_bytes_diff=$((read_sectors_final * 512 - read_sectors_initial * 512))
    write_bytes_diff=$((write_sectors_final * 512 - write_sectors_initial * 512))

    # 计算读写速度 (B/s)
    read_speed=$((read_bytes_diff / 1))  # 计算每秒的速度
    write_speed=$((write_bytes_diff / 1)) # 计算每秒的速度

    # 计算读写延迟 (ms)
    read_time_diff=$((read_time_final - read_time_initial))
    write_time_diff=$((write_time_final - write_time_initial))

    # 避免除以0的情况
    if [ "$read_ios_diff" -ne "0" ]; then
        read_latency=$(awk "BEGIN {print ($read_time_diff / $read_ios_diff) * 1000}")
    else
        read_latency=0
    fi

    if [ "$write_ios_diff" -ne "0" ]; then
        write_latency=$(awk "BEGIN {print ($write_time_diff / $write_ios_diff) * 1000}")
    else
        write_latency=0
    fi

    # 计算读写总字节数
    total_read_bytes=$((read_sectors_final * 512))
    total_write_bytes=$((write_sectors_final * 512))

    # 构建 JSON 对象
    if [ "$first" = true ]; then
        first=false
    else
        echo ","
    fi

    fs_type=$(echo -e "$fs_type" | tr -d '\n')

    echo "  {"
    echo "    \"partition\": \"/dev/$device\","
    echo "    \"size_bytes\": $total_size_bytes,"
    echo "    \"used_size_bytes\": $used_size_bytes,"
    echo "    \"mount_point\": \"$mount_point\","
    echo "    \"fs_type\": \"$fs_type\","
    echo "    \"read_speed_Bps\": $read_speed,"
    echo "    \"write_speed_Bps\": $write_speed,"
    echo "    \"read_iops\": $read_ios_diff,"
    echo "    \"write_iops\": $write_ios_diff,"
    echo "    \"read_latency_ms\": $read_latency,"
    echo "    \"write_latency_ms\": $write_latency,"
    echo "    \"total_read_bytes\": $total_read_bytes,"
    echo "    \"total_write_bytes\": $total_write_bytes"
    echo "  }"

done < $temp_file_mounts

echo "]"
rm $temp_file_initial $temp_file_final $temp_file_sizes $temp_file_mounts $temp_file_fstype $temp_file_usage

exit 0`
	getMemInfo = `#!/bin/bash

# 使用 free -b 命令获取内存信息
mem_info=$(free -b)

# 提取内存信息行
mem_line=$(echo "$mem_info" | grep -E "^Mem:")

# 解析内存信息
total=$(echo $mem_line | awk '{print $2}')
used=$(echo $mem_line | awk '{print $3}')
free=$(echo $mem_line | awk '{print $4}')
shared=$(echo $mem_line | awk '{print $5}')
buffcache=$(echo $mem_line | awk '{print $6}')
available=$(echo $mem_line | awk '{print $7}')

# 创建 JSON 格式
json=$(cat <<EOF
{
  "total": $total,
  "used": $used,
  "free": $free,
  "shared": $shared,
  "buffcache": $buffcache,
  "available": $available
}
EOF
)

# 输出 JSON
echo "$json"`
	getNetworkInfo = `#!/bin/bash

# 获取所有网络接口名称
interfaces=$(ls /sys/class/net)

# 初始化 JSON 数组
json="["

# 获取第一次读取的接收和发送字节数
declare -A initial_rx_bytes
declare -A initial_tx_bytes

for interface in $interfaces; do
  initial_rx_bytes[$interface]=$(cat /sys/class/net/$interface/statistics/rx_bytes)
  initial_tx_bytes[$interface]=$(cat /sys/class/net/$interface/statistics/tx_bytes)
done

# 等待一秒钟
sleep 1

# 遍历每个网络接口
for interface in $interfaces; do
  # 获取接口类型
  type=$(cat /sys/class/net/$interface/type)
  if [[ $type -eq 1 ]]; then
    type="Ethernet"
  elif [[ $type -eq 512 ]]; then
    type="Loopback"
  else
    type="Unknown"
  fi

  # 获取第二次读取的接收和发送字节数
  final_rx_bytes=$(cat /sys/class/net/$interface/statistics/rx_bytes)
  final_tx_bytes=$(cat /sys/class/net/$interface/statistics/tx_bytes)

  # 计算每秒接收和发送的字节数
  speed_rec=$((final_rx_bytes - ${initial_rx_bytes[$interface]}))
  speed_sent=$((final_tx_bytes - ${initial_tx_bytes[$interface]}))

  # 获取总的接收和发送的字节数
  bytes_rec=$final_rx_bytes
  bytes_sent=$final_tx_bytes

  # 追加接口信息到 JSON 数组
  json="$json{\"name\": \"$interface\", \"type\": \"$type\", \"bytes_rec\": $bytes_rec, \"bytes_sent\": $bytes_sent, \"speed_rec\": $speed_rec, \"speed_sent\": $speed_sent},"
done

# 去掉最后一个逗号并关闭 JSON 数组
json="${json%,}]"

# 输出 JSON
echo "$json"
`
	getProcessInfo = `#!/bin/bash

# 使用 ps 命令获取进程信息
process_info=$(ps -eo pcpu,pmem,comm --sort=-pcpu | awk 'NR>1 {print $1, $2, $3}')

# 初始化 JSON 数组
json="["

# 遍历进程信息
while read -r cpu mem command; do
  # 格式化为 JSON 字符串
  json="$json{\"cpu\": $cpu, \"mem\": $mem, \"command\": \"$command\"},"
done <<< "$process_info"

# 去掉最后一个逗号并关闭 JSON 数组
json="${json%,}]"

# 输出 JSON
echo "$json"
`
)

func (s *Service) GetHostStats(id uuid.UUID) (*Stats, error) {
	client, err := s.getHostClient(id)
	if err != nil {
		return nil, err
	}
	wg := sync.WaitGroup{}
	wg.Add(1)
	go func() {
		defer wg.Done()
		info, err := client.Run(getCPUInfo)
		if err != nil {
			logrus.Errorf("get cpu info error: %v", err)
			return
		}
		cpu := new(CPUInfo)
		if err := json.Unmarshal([]byte(info), cpu); err != nil {
			return
		}
		s.statsCache.Set(id.String(), cpu)
	}()
	wg.Add(1)
	go func() {
		defer wg.Done()
		info, err := client.Run(getDiskInfo)
		if err != nil {
			logrus.Errorf("get disk info error: %v", err)
			return
		}
		disk := new(DiskInfo)
		if err := json.Unmarshal([]byte(info), disk); err != nil {
			logrus.Errorf("unmarshal disk info error: %v, info: %s", err, info)
			return
		}
		s.statsCache.Set(id.String(), disk)
	}()
	wg.Add(1)
	go func() {
		defer wg.Done()
		st := time.Now()
		info, err := client.Run(getMemInfo)
		rtt := time.Since(st).Milliseconds()
		s.statsCache.Set(id.String(), RTT(rtt))
		if err != nil {
			logrus.Errorf("get mem info error: %v", err)
			return
		}
		mem := new(MemInfo)
		if err := json.Unmarshal([]byte(info), mem); err != nil {
			return
		}
		s.statsCache.Set(id.String(), mem)
	}()
	wg.Add(1)
	go func() {
		defer wg.Done()
		info, err := client.Run(getNetworkInfo)
		if err != nil {
			logrus.Errorf("get network info error: %v", err)
			return
		}
		network := new(NetworkInfo)
		if err := json.Unmarshal([]byte(info), network); err != nil {
			return
		}
		s.statsCache.Set(id.String(), network)
	}()
	wg.Add(1)
	go func() {
		defer wg.Done()
		info, err := client.Run(getProcessInfo)
		if err != nil {
			logrus.Errorf("get process info error: %v", err)
			return
		}
		process := new(ProcessInfo)
		if err := json.Unmarshal([]byte(info), process); err != nil {
			return
		}
		s.statsCache.Set(id.String(), process)
	}()
	wg.Wait()
	return s.statsCache.Get(id.String()), nil
}

func (s *Service) getHostClient(id uuid.UUID) (*sshutil.Client, error) {
	if client, ok := s.hostClientCache.Get(id.String()); ok {
		return client, nil
	}

	if host, err := s.DetailHost(id); err != nil {
		return nil, err
	} else {
		if client, err := sshutil.NewSSHClient(host.HostInfo); err != nil {
			return nil, err
		} else {
			s.hostClientCache.Set(id.String(), client)
			return client, nil
		}
	}
}
