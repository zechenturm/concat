package main

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/go-yaml/yaml"
)

// BufferSize specifies the size of the read and write buffers used
const BufferSize = 4

type input struct {
	Files []string `yaml:"files"`
	Cmds  []recipe `yaml:"recipes"`
}

// RelevantCmd returns the commands relevant to the file(name) given
func (in *input) RelevantCmd(file string) *recipe {
	for _, c := range in.Cmds {
		if c.IsRelevant(file) {
			return &c
		}
	}
	return nil
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
		rf := in.RelevantCmd(file)
		if rf != nil {
			// pipeCommands(f, &rf.Cmds[0])
			rf.Init(f)
			rf.Execute()
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
