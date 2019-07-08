package main

import (
	"io"
	"os/exec"
)

type command struct {
	Files []string `yaml:"files"`
	Cmd   string   `yaml:"cmd"`
	Args  []string `yaml:"args"`
	cmnd  *exec.Cmd
}

// IsRelevant checks if the comand is relevant to the file(name) passed in
func (cmd *command) IsRelevant(file string) bool {
	for _, matchFile := range cmd.Files {
		if file == matchFile {
			return true
		}
	}
	return false
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
