package terminal

import (
	"fmt"
	"golang.org/x/crypto/ssh"
	"io"
	"time"
)

type SSHTerminal struct {
	config    Config
	sc        SSHConfig
	sshClient *ssh.Client
	session   *ssh.Session
	stdin     io.WriteCloser
	stdout    io.Reader
}

type SSHConfig struct {
	Host       string
	Port       int
	Username   string
	Password   string
	PrivateKey string
	Passphrase string
}

func NewSSHTerminal(sc SSHConfig, config Config) *SSHTerminal {
	return &SSHTerminal{
		config: config,
		sc:     sc,
	}
}

func (t *SSHTerminal) Start() error {
	config := &ssh.ClientConfig{
		User:            t.sc.Username,
		Auth:            t.sc.GetAuthMethods(),
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
		Timeout:         time.Second * 30,
	}

	client, err := ssh.Dial("tcp", fmt.Sprintf("%s:%d", t.sc.Host, t.sc.Port), config)
	if err != nil {
		return err
	}

	t.sshClient = client

	session, err := client.NewSession()
	if err != nil {
		return err
	}
	t.session = session

	t.stdin, err = session.StdinPipe()
	if err != nil {
		return err
	}
	t.stdout, err = session.StdoutPipe()
	if err != nil {
		return err
	}

	modes := ssh.TerminalModes{
		ssh.ECHO:          1,
		ssh.TTY_OP_ISPEED: 14400,
		ssh.TTY_OP_OSPEED: 14400,
	}

	if err := session.RequestPty("xterm-256color", int(t.config.Rows), int(t.config.Cols), modes); err != nil {
		return err
	}

	return session.Shell()
}

func (t *SSHTerminal) Close() error {
	if t.session != nil {
		t.session.Close()
	}

	if t.sshClient != nil {
		t.sshClient.Close()
	}
	return nil
}

func (t *SSHTerminal) Write(p []byte) (n int, err error) {
	return t.stdin.Write(p)
}

func (t *SSHTerminal) Read(p []byte) (n int, err error) {
	return t.stdout.Read(p)
}

func (t *SSHTerminal) Resize(cols, rows uint32) error {
	return t.session.WindowChange(int(cols), int(rows))
}

func (sc *SSHConfig) GetAuthMethods() []ssh.AuthMethod {
	authMethods := make([]ssh.AuthMethod, 0)

	if len(sc.PrivateKey) > 0 {
		if len(sc.Passphrase) > 0 {
			signer, err := ssh.ParsePrivateKeyWithPassphrase([]byte(sc.PrivateKey), []byte(sc.Passphrase))
			if err == nil {
				authMethods = append(authMethods, ssh.PublicKeys(signer))
			}
		} else {
			signer, err := ssh.ParsePrivateKey([]byte(sc.PrivateKey))
			if err == nil {
				authMethods = append(authMethods, ssh.PublicKeys(signer))
			}
		}
	}

	if len(sc.Password) > 0 {
		authMethods = append(authMethods, ssh.Password(sc.Password))
	}

	return authMethods
}
