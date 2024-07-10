package containerd

import (
	"fmt"
	"github.com/MR5356/aurora/pkg/util/sshutil"
	"github.com/containerd/containerd"
	"github.com/sirupsen/logrus"
	"io"
	"net"
	"os"
	"strings"
)

const (
	localSocketFmt  = "/tmp/%s-containerd.sock"
	remoteSocketCmd = "find /run -name containerd.sock"
)

var remoteSocket = ""

type Client struct {
	client *containerd.Client

	socket string
	tunnel net.Conn
}

func NewClientWithSSH(sshInfo *sshutil.HostInfo) (*Client, error) {
	c := &Client{}

	sshClient, err := sshutil.NewSSHClient(*sshInfo)
	if err != nil {
		return nil, err
	}

	session, err := sshClient.GetSession()
	if err != nil {
		return nil, err
	}

	// get remote containerd.sock path
	if output, err := session.Output(remoteSocketCmd); err != nil {
		return nil, err
	} else {
		remoteSocket = strings.Trim(string(output), "\n")
	}
	logrus.Debugf("remote socket: %s", remoteSocket)

	// connect to containerd via ssh
	// FIXME: BUG: ssh: rejected: administratively prohibited (open failed)
	tunnel, err := sshClient.GetClient().Dial("unix", remoteSocket)
	if err != nil {
		logrus.Errorf("Failed to connect to containerd: %v", err)
		return nil, err
	}
	c.tunnel = tunnel

	localSocket := fmt.Sprintf(localSocketFmt, sshInfo.Host)
	_ = os.Remove(localSocket)
	c.socket = localSocket

	// create local listener
	localListener, err := net.Listen("unix", localSocket)
	if err != nil {
		logrus.Errorf("Failed to listen on local socket: %v", err)
		return nil, err
	}

	go func() {
		for {
			localConn, err := localListener.Accept()
			if err != nil {
				logrus.Errorf("Failed to accept connection: %v", err)
				return
			}
			go func() {
				defer localConn.Close()
				defer tunnel.Close()

				go io.Copy(localConn, tunnel)
				io.Copy(tunnel, localConn)
			}()
		}
	}()

	// create containerd client
	containerdClient, err := containerd.New(localSocket)
	if err != nil {
		logrus.Errorf("Failed to create containerd client: %v", err)
		return nil, err
	}
	c.client = containerdClient

	return c, nil
}

func (c *Client) Close() {
	_ = c.client.Close()
	_ = c.tunnel.Close()
	_ = os.Remove(c.socket)
}
