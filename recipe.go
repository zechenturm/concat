package main

import (
	"fmt"
	"io"
)

type recipe struct {
	Files      []string  `yaml:"files"`
	Cmds       []command `yaml:"commands"`
	IgnoreFile bool      `yaml:"ignoreFile"`
}

// IsRelevant checks if the comand is relevant to the file(name) passed in
func (r *recipe) IsRelevant(file string) bool {
	for _, matchFile := range r.Files {
		if file == matchFile {
			return true
		}
	}
	return false
}

// CmdCount returns the total number of commands in the recipe
func (r *recipe) CmdCount() int {
	return len(r.Cmds)
}

// Init initialises the recipe
func (r *recipe) Init(file io.ReadCloser) {
	if len(r.Cmds) == 0 {
		return
	}

	for i := range r.Cmds {
		r.Cmds[i].Init()
	}

	cmdIn, err := r.Cmds[0].GetStdin()
	if err != nil {
		errorln(err)
		return
	}
	connect(file, cmdIn)

	if r.CmdCount() >= 2 {
		for i := 1; i < r.CmdCount(); i++ {
			cout0, err := r.Cmds[i-1].GetStdout()
			if err != nil {
				errorln(err)
			}
			cin1, err := r.Cmds[i].GetStdin()
			if err != nil {
				errorln(err)
			}
			connect(cout0, cin1)
		}
	}

	cout, err := r.Cmds[len(r.Cmds)-1].GetStdout()
	if err != nil {
		errorln(err)
		return
	}

	go func() {
		var n int
		buf := make([]byte, BufferSize)
		for err == nil {
			n, err = cout.Read(buf)
			fmt.Print(string(buf[:n]))

		}
	}()

}

// Execute executes the recipe, blocking until execution is done
func (r *recipe) Execute() {
	if r.CmdCount() == 0 {
		return
	}
	for _, c := range r.Cmds {
		err := c.Start()
		if err != nil {
			errorln("Error starting:", err)
			return
		}
	}

	err := r.Cmds[len(r.Cmds)-1].Wait()
	if err != nil {
		errorln("Error waiting:", err)
	}
}

func connect(src io.ReadCloser, dst io.WriteCloser) {
	go func() {
		_, err := io.Copy(dst, src)
		if err != nil {
			errorln(err)
		}
		err = src.Close()
		if err != nil {
			errorln(err)
		}
		err = dst.Close()
		if err != nil {
			errorln(err)
		}
	}()
}
