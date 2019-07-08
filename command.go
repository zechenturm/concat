package main

import (
	"io"
	"os/exec"
)

type command struct {
	Cmd  string   `yaml:"cmd"`
	Args []string `yaml:"args"`
	cmnd *exec.Cmd
}

// Init builds an exec.Cmd to exectute
func (cmd *command) Init() {
	cmd.cmnd = exec.Command(cmd.Cmd, cmd.Args...)
}

func (cmd *command) Start() error {
	return cmd.cmnd.Start()
}

func (cmd *command) Wait() error {
	return cmd.cmnd.Wait()
}

func (cmd *command) GetStdPipes() (io.WriteCloser, io.ReadCloser, error) {
	stdin, err := cmd.cmnd.StdinPipe()
	stdout, err := cmd.cmnd.StdoutPipe()
	return stdin, stdout, err
}
