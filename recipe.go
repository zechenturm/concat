package main

import (
	"fmt"
	"io"
	"os"
)

type recipe struct {
	Files []string  `yaml:"files"`
	Cmds  []command `yaml:"commands"`
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

func (r *recipe) CmdCount() int {
	return len(r.Cmds)
}

func (r *recipe) Init(file io.ReadCloser) {
	if len(r.Cmds) == 0 {
		return
	}

	for i := range r.Cmds {
		r.Cmds[i].Init()
	}

	cmdIn, err := r.Cmds[0].GetStdin()
	if err != nil {
		fmt.Println(err)
		return
	}
	connect(file, cmdIn)

	if r.CmdCount() >= 2 {
		for i := 1; i < r.CmdCount(); i++ {
			fmt.Println(i)
			cout0, err := r.Cmds[i-1].GetStdout()
			if err != nil {
				fmt.Println(err)
			}
			cin1, err := r.Cmds[i].GetStdin()
			if err != nil {
				fmt.Println(err)
			}
			connect(cout0, cin1)
		}
	}

	cout, err := r.Cmds[len(r.Cmds)-1].GetStdout()
	if err != nil {
		fmt.Println(err)
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

func (r *recipe) Execute() {
	if r.CmdCount() == 0 {
		return
	}
	for _, c := range r.Cmds {
		err := c.Start()
		if err != nil {
			fmt.Println("Error starting:", err)
			return
		}
	}

	err := r.Cmds[len(r.Cmds)-1].Wait()
	if err != nil {
		fmt.Println("Error wating:", err)
	}
}

func connect(src io.ReadCloser, dst io.WriteCloser) {
	go func() {
		_, err := io.Copy(dst, src)
		if err != nil {
			fmt.Println(err)
		}
		err = src.Close()
		if err != nil {
			fmt.Println(err)
		}
		err = dst.Close()
		if err != nil {
			fmt.Println(err)
		}
	}()
}

func pipeCommands(f *os.File, cmd *command) {
	cmd.Init()
	stdin, stdout, err := cmd.GetStdPipes()
	if err != nil {
		fmt.Println(err)
		return
	}
	cmd.Start()

	connect(f, stdin)

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
		fmt.Println("Error for", f.Name(), ":", err)
	}
}
