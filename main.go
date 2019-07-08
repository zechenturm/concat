package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"os/exec"

	"github.com/go-yaml/yaml"
)

// BufferSize specifies the size of the read and write buffers used
const BufferSize = 4

type command struct {
	Files []string `yaml:"files"`
	Cmd   string   `yaml:"cmd"`
	Args  []string `yaml:"args"`
	cmnd  *exec.Cmd
}

type input struct {
	Files []string  `yaml:"files"`
	Cmds  []command `yaml:"cmds"`
}

// RelevantCmds returns the commands relevant to the file(name) given
func (in *input) RelevantCmds(file string) []*command {
	var commands []*command
	for _, c := range in.Cmds {
		if c.IsRelevant(file) {
			commands = append(commands, &c)
		}
	}
	return commands
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

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Usage: apply <configfile>")
		return
	}
	data, err := ioutil.ReadFile(os.Args[1])
	if err != nil {
		fmt.Println(err)
		return
	}

	var in input
	err = yaml.Unmarshal(data, &in)
	if err != nil {
		fmt.Println(err)
		return
	}

	for _, file := range in.Files {
		f, err := os.Open(file)
		if err != nil {
			fmt.Println(err)
			continue
		}
		rf := in.RelevantCmds(file)
		if rf != nil {
			for _, cmd := range rf {
				cmd.Init()
				stdin, stdout, err := cmd.GetStdPipes()
				if err != nil {
					fmt.Println(err)
					continue
				}
				cmd.Start()

				go func() {
					_, err := io.Copy(stdin, f)
					if err != nil {
						fmt.Println(err)
					}
					err = f.Close()
					if err != nil {
						fmt.Println(err)
					}
					err = stdin.Close()
					if err != nil {
						fmt.Println(err)
					}
				}()

				go func() {
					var err error
					var n int
					data := make([]byte, BufferSize)
					for err == nil {
						n, err = stdout.Read(data)
						fmt.Print(string(data[:n]))
					}
				}()
				err = cmd.Wait()
				if err != nil {
					fmt.Println("Error for", file, ":", err)
				}
			}
		} else {
			var n int
			data := make([]byte, BufferSize)
			for err == nil {
				n, err = f.Read(data)
				fmt.Print(string(data[:n]))
			}
		}
	}
}
