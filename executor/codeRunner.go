package executor

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

//Language represents the detected language of the snippet (currently supported: go, java, c)
type Language int

const (
	//GOLANG is the go language
	GOLANG Language = iota
	//CLANG is the c language
	CLANG
	//JAVALANG is the java language
	JAVALANG
	//UNKNOWN is a not supported language
	UNKNOWN
)

//CodePacket represents all relevant data for the code executution of 1 file
type CodePacket struct {
	Title    string
	Path     string
	Filename string
	Inputs   []string
	Timeout  int
	RunData  *RunInfo
}

//RunInfo contains all data from the execution of the code
type RunInfo struct {
	Lang    Language
	OutStd  string
	OutErr  string
	Runtime float64
}

//Run executes the code of the code packet and stores the outputs
func (p CodePacket) Run() {

	lang := detectFileType(p.Filename)

	fmt.Println("running file: " + p.Path + p.Filename + " type=" + fmt.Sprintf("%d", lang))
	tBegin := time.Now()

	cmdString := "go"
	args := []string{"run", p.Path + p.Filename}

	if lang == CLANG {
		filename := compileC(p)
		println("executing compiled program")
		cmdString = filepath.FromSlash("./" + filename)
		args = []string{}

		defer os.Remove(cmdString)
	}

	if lang == UNKNOWN {
		fmt.Println("got unknown language!")
		return
	}

	cmd := exec.Command(cmdString, args...)
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	stdin, err := cmd.StdinPipe()
	if err != nil {
		fmt.Println(err)
		p.RunData.OutErr = p.RunData.OutErr + err.Error()
		return
	}
	defer stdin.Close()

	err = cmd.Start()
	if err != nil {
		fmt.Println("code execution error: " + err.Error())
		fmt.Println(string(stderr.Bytes()))
		fmt.Println("---------------------------------")

		p.RunData.OutErr = p.RunData.OutErr + err.Error() + "\n" + string(stderr.Bytes())
		return
	}

	fmt.Println("setting inputs...")
	done := make([]chan error, len(p.Inputs))
	for i, elem := range p.Inputs {
		io.WriteString(stdin, elem+"\n")

		index := i

		//write the answer into the done channel
		go func() { done[index] <- cmd.Wait() }()
		addTime := 0
		if i == 0 {
			addTime = 200
		}

		//select whichever comes first: done chan or timer
		select {
		case _ = <-done[index]:
		case <-time.After(time.Duration(p.Timeout+addTime) * time.Millisecond):
		}

	}

	response := string(stdout.Bytes())
	e := string(stderr.Bytes())
	duration := time.Now().Sub(tBegin)
	//fmt.Println("runtime: " + fmt.Sprintf("%.4f", duration.Seconds()))

	//fmt.Print(response)

	p.RunData.OutStd = response
	p.RunData.OutErr = e
	p.RunData.Runtime = duration.Seconds()
	fmt.Println("code execution done")

}

func detectFileType(name string) Language {
	if strings.HasSuffix(name, ".go") {
		return GOLANG
	} else if strings.HasSuffix(name, ".c") {
		return CLANG
	} else if strings.HasSuffix(name, ".java") {
		return JAVALANG
	}

	return UNKNOWN

}

func compileC(p CodePacket) string {
	name := strings.Split(p.Filename, ".")[0]

	cmdString := "gcc"
	args := []string{p.Path + p.Filename, "-o", name + p.Title, "-lm"}
	err := exec.Command(cmdString, args...).Run()
	if err != nil {
		fmt.Println("got error while compiling c program; " + err.Error())
	}

	fmt.Println("compiled C program!")
	return name + p.Title
}
