package main

import (
	"fmt"
	"os"
	"strconv"
	"sync"

	"github.com/dwladdimiroc/stats-storm/exec"
	files "github.com/dwladdimiroc/stats-storm/read_file"
)

func main() {
	nameApp := os.Args[1]

	duration, errParser := strconv.Atoi(os.Args[2])
	if errParser != nil {
		fmt.Println("Error... No es un número válido el ingresado como parámetro...")
		fmt.Println(errParser)
		os.Exit(2)
	}
	var wg sync.WaitGroup

	// (1) Execute Stats CPU
	appCmdCPU := "sar"
	argsCmdCPU := []string{"-u", "1"}
	dirCmdCPU := ""
	var outputCPU string

	wg.Add(1)
	go func(appCmd string, argsCmd []string, output *string, dirCmdCPU string, duration int) {
		defer wg.Done()
		*output = exec.Start(appCmd, argsCmd, dirCmdCPU, duration)
	}(appCmdCPU, argsCmdCPU, &outputCPU, dirCmdCPU, duration)

	// (2) Execute Stats Memory
	appCmdMemory := "vmstat"
	argsCmdMemory := []string{"-n", "-S", "M", "1"}
	dirCmdMemory := ""
	var outputMemory string

	wg.Add(1)
	go func(appCmd string, argsCmd []string, output *string, dirCmdMemory string, duration int) {
		defer wg.Done()
		*output = exec.Start(appCmd, argsCmd, dirCmdMemory, duration)
	}(appCmdMemory, argsCmdMemory, &outputMemory, dirCmdMemory, duration)

	// (3) Execute App
	appCmdStormApp := "sh"
	argsCmdStormApp := []string{"startApp.sh"}
	dirCmdStormApp := "apps"
	exec.Execute(appCmdStormApp, argsCmdStormApp, dirCmdStormApp)

	//	// (4) Execute Monitor
	appCmdMonitor := "java"
	argsCmdMonitor := []string{"-jar", "storm-monitor.jar", nameApp, "true"}
	dirCmdMonitor := "monitor"

	wg.Add(1)
	go func(appCmd string, argsCmd []string, dirCmdMonitor string, duration int) {
		defer wg.Done()
		exec.Start(appCmd, argsCmd, dirCmdMonitor, duration)
	}(appCmdMonitor, argsCmdMonitor, dirCmdMonitor, duration)

	wg.Wait()

	files.ParseCPU(outputCPU, nameApp)
	files.ParseMemory(outputMemory, nameApp)

	appCmdStormApp = "storm"
	argsCmdStormApp = []string{"kill", nameApp, "-w", "0"}
	dirCmdStormApp = ""
	exec.Execute(appCmdStormApp, argsCmdStormApp, dirCmdStormApp)
}
