package terminal

import (
	"github.com/creack/pty"
	"os"
	"os/exec"
	"runtime"
)

type LocalTerminal struct {
	config Config
	pty    *os.File
	cmd    *exec.Cmd
}

func NewLocalTerminal(config Config) *LocalTerminal {
	return &LocalTerminal{
		config: config,
	}
}

func (t *LocalTerminal) Start() error {
	shell := os.Getenv("SHELL")
	if shell == "" {
		shell = "/bin/bash"
	}
	if runtime.GOOS == "windows" {
		shell = "cmd.exe"
	}

	cmd := exec.Command(shell)
	f, err := pty.Start(cmd)
	if err != nil {
		return err
	}

	t.pty = f
	t.cmd = cmd

	return t.Resize(t.config.Cols, t.config.Rows)
}

func (t *LocalTerminal) Close() error {
	if t.pty != nil {
		t.pty.Close()
	}
	if t.cmd != nil && t.cmd.Process != nil {
		t.cmd.Process.Kill()
	}
	return nil
}

func (t *LocalTerminal) Write(p []byte) (n int, err error) {
	return t.pty.Write(p)
}

func (t *LocalTerminal) Read(p []byte) (n int, err error) {
	return t.pty.Read(p)
}

func (t *LocalTerminal) Resize(cols, rows uint32) error {
	return pty.Setsize(t.pty, &pty.Winsize{
		Rows: uint16(rows),
		Cols: uint16(cols),
	})
}
