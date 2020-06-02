package exec

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	//	"strings"
	"syscall"
	"time"
)

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func printError(err error, app string) {
	if err != nil {
		//check(err)
		os.Stderr.WriteString(fmt.Sprintf("["+app+"] Error: %s\n", err.Error()))
	}
}

func printOutput(outs []byte) {
	if len(outs) > 0 {
		fmt.Printf(string(outs))
	}
}

func print(line string) {
	fmt.Println(line)
}

func Run(app string, args []string, dir string, interval int) string {
	print("Ejecutando " + app)

	cmd := exec.Command(app)
	if dir != "" {
		cmd.Dir = dir
	}
	outputBytes := &bytes.Buffer{}
	cmd.Stdout = outputBytes

	err := cmd.Start()
	printError(err, app)

	timer := time.NewTimer(time.Second * time.Duration(interval))
	go func(timer *time.Timer, cmd *exec.Cmd, app string) {
		for _ = range timer.C {
			err := cmd.Process.Signal(os.Kill)
			printError(err, app)
		}
	}(timer, cmd, args[0])

	// Only proceed once the process has finished
	cmd.Wait()
	return string(outputBytes.Bytes())
}

func Execute(app string, args []string, dir string) string {
	print("Ejecutando " + app)
	cmd := exec.Command(app, args...)
	if dir != "" {
		cmd.Dir = dir
	}
	cmd.SysProcAttr = &syscall.SysProcAttr{Setpgid: true}
	outputBytes := &bytes.Buffer{}
	cmd.Stdout = outputBytes

	err := cmd.Run()
	printError(err, app)

	syscall.Kill(-cmd.Process.Pid, syscall.SIGKILL)

	return string(outputBytes.Bytes())
}

func Start(app string, args []string, dir string, interval int) string {
	print("Ejecutando " + app)
	cmd := exec.Command(app, args...)
	cmd.Dir = dir
	cmd.SysProcAttr = &syscall.SysProcAttr{Setpgid: true}
	outputBytes := &bytes.Buffer{}
	cmd.Stdout = outputBytes

	err := cmd.Start()
	printError(err, app)

	time.Sleep(time.Second * time.Duration(interval))

	syscall.Kill(-cmd.Process.Pid, syscall.SIGKILL)

	return string(outputBytes.Bytes())
}
