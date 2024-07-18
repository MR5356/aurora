package containerd

import (
	"context"
	"net"
	"os"
	"strings"

	"github.com/MR5356/aurora/pkg/util/sshutil"
	"github.com/containerd/containerd"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
)

const (
	localSocketFmt  = "./tmp/%s-containerd.sock"
	remoteSocketCmd = "find /run -name containerd.sock"
)

var remoteSocket = ""

func init() {
	os.MkdirAll("./tmp", os.ModePerm)
}

type Client struct {
	client *containerd.Client

	socket   string
	tunnel   net.Conn
	grocConn *grpc.ClientConn
}

func NewClientWithSSH(sshInfo *sshutil.HostInfo) (*Client, error) {
	c := &Client{}

	sshClient, err := sshutil.NewSSHClient(*sshInfo)
	if err != nil {
		logrus.Debugf("new ssh client error: %+v", err)
		return nil, err
	}

	session, err := sshClient.GetSession()
	if err != nil {
		logrus.Debugf("get ssh session error: %+v", err)
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

	conn, err := grpc.Dial("", grpc.WithContextDialer(func(ctx context.Context, _ string) (net.Conn, error) {
		return tunnel, nil
	}), grpc.WithInsecure())
	if err != nil {
		logrus.Debugf("create grpc client error: %+v", err)
		return nil, err
	}
	c.grocConn = conn

	// create containerd client
	containerdClient, err := containerd.NewWithConn(conn)
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
	_ = c.grocConn.Close()
	_ = os.Remove(c.socket)
}
