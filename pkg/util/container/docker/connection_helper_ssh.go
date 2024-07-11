package docker

import (
	"context"
	"github.com/MR5356/aurora/pkg/util/sshutil"
	"github.com/pkg/errors"
	"golang.org/x/crypto/ssh"
	"io"
	"net"
	"os"
	"strings"
	"sync/atomic"
	"time"
)

const (
	dockerDialStdio = "docker system dial-stdio"
	dummyHost       = "https://docker.ac.cn"
)

type SSHConnectionHelper struct {
	Dialer func(ctx context.Context, network, addr string) (net.Conn, error)
	Host   string
}

func GetSSHConnectionHelper(info *sshutil.HostInfo) (*SSHConnectionHelper, error) {
	return &SSHConnectionHelper{
		Dialer: func(ctx context.Context, network, addr string) (net.Conn, error) {
			var conn sshConn

			sshClient, err := sshutil.NewSSHClient(*info)
			if err != nil {
				return nil, err
			}

			session, err := sshClient.GetSession()
			if err != nil {
				return nil, err
			}

			conn.session = session
			if conn.stdin, err = session.StdinPipe(); err != nil {
				return nil, err
			}

			if conn.stdout, err = session.StdoutPipe(); err != nil {
				return nil, err
			}

			if conn.stderr, err = session.StderrPipe(); err != nil {
				return nil, err
			}

			conn.localAddr = dummyAddr{network: "dummy", s: "dummy-local"}
			conn.remoteAddr = dummyAddr{network: "dummy", s: "dummy-remote"}

			return &conn, session.Start(dockerDialStdio)
		},
		Host: dummyHost,
	}, nil
}

func (c *SSHConnectionHelper) GetDialer() func(ctx context.Context, network, addr string) (net.Conn, error) {
	return c.Dialer
}

func (c *SSHConnectionHelper) GetHost() string {
	return c.Host
}

type dummyAddr struct {
	network string
	s       string
}

func (d dummyAddr) Network() string {
	return d.network
}

func (d dummyAddr) String() string {
	return d.s
}

type sshConn struct {
	session    *ssh.Session
	stdin      io.WriteCloser
	stdout     io.Reader
	stderr     io.Reader
	closing    atomic.Bool
	localAddr  net.Addr
	remoteAddr net.Addr
}

func (c *sshConn) Read(p []byte) (int, error) {
	n, err := c.stdout.Read(p)
	if c.closing.Load() {
		return n, err
	}
	return n, c.handleEOF(err)
}

func (c *sshConn) Write(p []byte) (int, error) {
	n, err := c.stdin.Write(p)
	if c.closing.Load() {
		return n, err
	}
	return n, c.handleEOF(err)
}

func (c *sshConn) handleEOF(err error) error {
	if err != io.EOF {
		return err
	}

	return errors.Errorf("execution failed with %v, make sure the URL is valid, and Docker 18.09 or later is installed on the remote host", err)
}

func (c *sshConn) Close() error {
	c.closing.Store(true)
	defer c.closing.Store(false)

	if err := c.stdin.Close(); err != nil && strings.Contains(err.Error(), os.ErrClosed.Error()) {
		return err
	}

	return c.session.Close()
}

func (c *sshConn) LocalAddr() net.Addr {
	return c.localAddr
}

func (c *sshConn) RemoteAddr() net.Addr {
	return c.remoteAddr
}

func (c *sshConn) SetDeadline(t time.Time) error {
	return nil
}

func (c *sshConn) SetReadDeadline(t time.Time) error {
	return nil
}

func (c *sshConn) SetWriteDeadline(t time.Time) error {
	return nil
}
