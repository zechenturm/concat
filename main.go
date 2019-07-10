package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"strings"

	"github.com/go-yaml/yaml"
)

// BufferSize specifies the size of the read and write buffers used
const BufferSize = 4

type stringRC struct {
	S   string
	src io.Reader
}

func (src *stringRC) Read(b []byte) (int, error) {
	if src.src == nil {
		src.src = strings.NewReader(src.S)
	}
	return src.src.Read(b)
}

func (src *stringRC) Close() error {
	return nil
}

func errorln(e ...interface{}) {
	fmt.Fprintln(os.Stderr, e...)
}

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
		errorln("Usage: apply <configfile>")
		return
	}
	data, err := ioutil.ReadFile(os.Args[1])
	if err != nil {
		errorln(err)
		return
	}

	var in input
	err = yaml.Unmarshal(data, &in)
	if err != nil {
		errorln(err)
		return
	}

	for _, file := range in.Files {
		rf := in.RelevantCmd(file)
		var f io.ReadCloser
		if rf != nil {
			// pipeCommands(f, &rf.Cmds[0])
			fmt.Println("relevant file:", rf)
			if rf.IgnoreFile == false {
				f, err = os.Open(file)
				if err != nil {
					errorln("Error opening", file, ":", err)
					continue
				}
			} else {
				f = &stringRC{S: file}
			}
			rf.Init(f)
			rf.Execute()
		} else {
			f, err = os.Open(file)
			if err != nil {
				errorln("Error opening", file, ":", err)
				continue
			}
			var n int
			data := make([]byte, BufferSize)
			for err == nil {
				n, err = f.Read(data)
				fmt.Print(string(data[:n]))
			}
		}
	}
}
