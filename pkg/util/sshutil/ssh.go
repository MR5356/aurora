package sshutil

import (
	"bufio"
	"bytes"
	"context"
	"fmt"
	"github.com/go-ping/ping"
	"github.com/sirupsen/logrus"
	"golang.org/x/crypto/ssh"
	"strconv"
	"strings"
	"sync"
	"time"
)

type Client struct {
	host      HostInfo
	sshClient *ssh.Client
}

func NewSSHClient(host HostInfo) (client *Client, err error) {
	sshClient, err := ssh.Dial("tcp", fmt.Sprintf("%s:%d", host.Host, host.Port), &ssh.ClientConfig{
		User:            host.Username,
		Auth:            []ssh.AuthMethod{ssh.Password(host.Password)},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
		Timeout:         time.Second * 10,
	})

	if err != nil {
		return
	}

	client = &Client{
		host:      host,
		sshClient: sshClient,
	}
	return
}

func (c *Client) GetSession() (*ssh.Session, error) {
	return c.sshClient.NewSession()
}

func (c *Client) GetStats() (stats *Stats, err error) {
	stats = &Stats{}
	wg := &sync.WaitGroup{}

	wg.Add(1)
	go func() {
		defer wg.Done()
		memInfo, e := c.GetMemInfo()
		if e != nil {
			logrus.Warnf("get mem info error: %v", e)
			memInfo = &MemInfo{}
			memInfo.Default()
		}
		stats.MemInfo = memInfo
	}()
	wg.Add(1)
	go func() {
		defer wg.Done()
		cpuInfo, e := c.GetCpuInfo()
		if e != nil {
			logrus.Warnf("get cpu info error: %v", e)
			cpuInfo = &CpuInfo{}
			cpuInfo.Default()
		}
		stats.CpuInfo = cpuInfo
	}()
	wg.Add(1)
	go func() {
		defer wg.Done()
		uptime, e := c.GetUptime()
		if e != nil {
			logrus.Warnf("get uptime error: %v", e)
		}
		stats.Uptime = uptime
	}()
	wg.Add(1)
	go func() {
		defer wg.Done()
		rtt, e := c.GetPing()
		if e != nil {
			logrus.Warnf("get ping error: %v", e)
			err = e
		}
		stats.Ping = rtt
	}()
	wg.Add(1)
	go func() {
		defer wg.Done()
		netInfos, e := c.GetNetInfo()
		if e != nil {
			logrus.Warnf("get net info error: %v", e)
			ni := &NetInfo{}
			ni.Default()
			netInfos = []*NetInfo{
				ni,
			}
		}
		stats.NetInfo = netInfos
	}()
	wg.Wait()

	return
}

func (c *Client) GetNetInfo() (netInfos []*NetInfo, err error) {
	interfaceNames, err := c.getNetInterface()
	if err != nil {
		return
	}
	t1 := time.Now()
	preNet, err := c.Run("/bin/cat /proc/net/dev")
	t2 := time.Since(t1)
	if err != nil {
		return
	}
	time.Sleep(time.Second)
	t3 := time.Now()
	newNet, err := c.Run("/bin/cat /proc/net/dev")
	t4 := time.Since(t3)
	if err != nil {
		return
	}

	totalTime := (t4+t2)/2 + time.Second
	preScanner := bufio.NewScanner(strings.NewReader(preNet))
	nowScanner := bufio.NewScanner(strings.NewReader(newNet))
	for _, ifName := range interfaceNames {
		preNetInfo := &NetInfo{
			Name: ifName,
		}
		for preScanner.Scan() {
			parts := strings.Fields(preScanner.Text())
			if len(parts) == 17 && parts[0] == fmt.Sprintf("%s:", ifName) {
				preNetInfo.BytesRec, _ = strconv.ParseUint(parts[1], 10, 64)
				preNetInfo.BytesSent, _ = strconv.ParseUint(parts[9], 10, 64)
			}
		}

		netInfo := &NetInfo{
			Name: ifName,
		}

		for nowScanner.Scan() {
			parts := strings.Fields(nowScanner.Text())
			if len(parts) == 17 && parts[0] == fmt.Sprintf("%s:", ifName) {
				netInfo.BytesRec, _ = strconv.ParseUint(parts[1], 10, 64)
				netInfo.BytesSent, _ = strconv.ParseUint(parts[9], 10, 64)
				netInfo.SpeedRec = uint64(float64(netInfo.BytesRec-preNetInfo.BytesRec) / totalTime.Seconds())
				netInfo.SpeedSent = uint64(float64(netInfo.BytesSent-preNetInfo.BytesSent) / totalTime.Seconds())
			}
		}

		netInfos = append(netInfos, netInfo)
	}
	return
}

func (c *Client) getNetInterface() (interfaceNames []string, err error) {
	interfaceString, err := c.Run("/bin/ls -l /sys/class/net/")
	if err != nil {
		return
	}
	scanner := bufio.NewScanner(strings.NewReader(interfaceString))
	for scanner.Scan() {
		if !strings.Contains(scanner.Text(), "virtual") {
			parts := strings.Fields(scanner.Text())
			if len(parts) > 9 {
				interfaceNames = append(interfaceNames, parts[8])
			}
		}
	}
	return
}

func (c *Client) GetPing() (rtt int64, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()
	stdoutChan := make(chan int64)
	errChan := make(chan error)

	rtt = 9999

	go func() {
		pinger, err := ping.NewPinger(c.host.Host)
		if err != nil {
			errChan <- err
		}
		pinger.Count = 1
		pinger.SetPrivileged(true)
		err = pinger.Run()
		if err != nil {
			errChan <- err
		}
		stats := pinger.Statistics()
		rtt = stats.AvgRtt.Milliseconds()
		stdoutChan <- rtt
	}()

	select {
	case err = <-errChan:
		return rtt, err
	case rtt = <-stdoutChan:
		return rtt, err
	case <-ctx.Done():
		return rtt, fmt.Errorf("ping timeout")
	}
}

func (c *Client) GetUptime() (uptime int64, err error) {
	ut, err := c.Run("/bin/cat /proc/uptime")
	if err != nil {
		return
	}
	parts := strings.Fields(ut)
	if len(parts) == 2 {
		upSecond, err := strconv.ParseFloat(parts[0], 64)
		if err != nil {
			return uptime, err
		}
		return int64(upSecond), nil
	}
	return
}

func (c *Client) GetCpuInfo() (cpuInfo *CpuInfo, err error) {
	pre, err := c.getCpuInfo()
	if err != nil {
		return
	}
	time.Sleep(time.Second / 2)
	now, err := c.getCpuInfo()
	if err != nil {
		return
	}
	total := now.total - pre.total
	cpuInfo = &CpuInfo{
		User:    (now.User - pre.User) / total * 100,
		Nice:    (now.Nice - pre.Nice) / total * 100,
		System:  (now.System - pre.System) / total * 100,
		Idle:    (now.Idle - pre.Idle) / total * 100,
		IOWait:  (now.IOWait - pre.IOWait) / total * 100,
		Irq:     (now.Irq - pre.Irq) / total * 100,
		SoftIrq: (now.SoftIrq - pre.SoftIrq) / total * 100,
		Steal:   (now.Steal - pre.Steal) / total * 100,
		Guest:   (now.Guest - pre.Guest) / total * 100,
		Cores:   now.Cores,
	}
	cpuInfo.Percent = 100 - cpuInfo.Idle
	return
}

func (c *Client) getCpuInfo() (cpuInfo *CpuInfo, err error) {
	cpuInfo = &CpuInfo{}
	cpuInfo.Default()
	cpu, err := c.Run("/bin/cat /proc/stat")
	if err != nil {
		return cpuInfo, err
	}
	scanner := bufio.NewScanner(strings.NewReader(cpu))
	for scanner.Scan() {
		fields := strings.Fields(scanner.Text())
		if strings.Contains(fields[0], "cpu") && fields[0] != "cpu" {
			cpuInfo.Cores += 1
		}
		if len(fields) > 0 && fields[0] == "cpu" {
			numFields := len(fields)
			for i := 1; i < numFields; i++ {
				val, err := strconv.ParseFloat(fields[i], 64)
				if err != nil {
					return cpuInfo, err
				}
				cpuInfo.total += val
				switch i {
				case 1:
					cpuInfo.User = val
				case 2:
					cpuInfo.Nice = val
				case 3:
					cpuInfo.System = val
				case 4:
					cpuInfo.Idle = val
				case 5:
					cpuInfo.IOWait = val
				case 6:
					cpuInfo.Irq = val
				case 7:
					cpuInfo.SoftIrq = val
				case 8:
					cpuInfo.Steal = val
				case 9:
					cpuInfo.Guest = val
				}
			}
		}
	}
	return
}

func (c *Client) GetMemInfo() (memInfo *MemInfo, err error) {
	memInfo = &MemInfo{}
	memInfo.Default()
	mem, err := c.Run("/bin/cat /proc/meminfo")
	if err != nil {
		return
	}
	scanner := bufio.NewScanner(strings.NewReader(mem))
	for scanner.Scan() {
		parts := strings.Fields(scanner.Text())
		if len(parts) == 3 {
			val, err := strconv.ParseUint(parts[1], 10, 64)
			if err != nil {
				return nil, err
			}
			val *= 1024
			switch parts[0] {
			case "MemTotal:":
				memInfo.Total = val
			case "MemFree:":
				memInfo.Free = val
			case "MemAvailable:":
				memInfo.Available = val
			case "Buffers:":
				memInfo.Buffers = val
			case "Cached:":
				memInfo.Cached = val
			case "SwapTotal:":
				memInfo.SwapTotal = val
			case "SwapFree:":
				memInfo.SwapFree = val
			}
		}
	}
	memInfo.Percent = (float64(memInfo.Total-memInfo.Available) / float64(memInfo.Total)) * 100
	return
}

func (c *Client) Close() {
	err := c.sshClient.Close()
	if err != nil {
		logrus.Warnf("SSH client close error: %v", err)
	}
}

func (c *Client) Run(cmd string) (stdout string, err error) {
	stdoutChan := make(chan string)
	errChan := make(chan error)

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	go func() {
		run, err := c.RunCmd(cmd)
		if err != nil {
			errChan <- err
		}
		stdoutChan <- run
	}()
	select {
	case err = <-errChan:
		return stdout, err
	case stdout = <-stdoutChan:
		return stdout, err
	case <-ctx.Done():
		return "", fmt.Errorf("timeout")
	}
}

func (c *Client) RunCmd(cmd string) (stdout string, err error) {
	session, err := c.sshClient.NewSession()
	if err != nil {
		return
	}
	defer session.Close()

	var buffer bytes.Buffer
	session.Stdout = &buffer

	err = session.Run(cmd)
	if err != nil {
		return
	}
	stdout = buffer.String()
	return
}
