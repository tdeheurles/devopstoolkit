package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"path"
	"sync"
)

// BinaryExecutorer is an interface for executing a binary
type BinaryExecutorer interface {
	Execute(version string) (stop bool, exitCode int)
}

type BinaryExecutor struct {
}

func NewBinaryExecutor() *BinaryExecutor {
	return &BinaryExecutor{}
}

func (b BinaryExecutor) Execute(version string) (stop bool, exitCode int) {
	// TODO: install the binary if it doesn't exist
	home := os.Getenv("HOME")
	commandPath := path.Join(home, ".devopstoolkit", "bin", "devopsrunner-"+version)
	cmd := exec.Command(commandPath, os.Args[1:]...)

	stdoutPipe, err := cmd.StdoutPipe()
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	stderrPipe, err := cmd.StderrPipe()
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	err = cmd.Start()
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	var wg sync.WaitGroup
	wg.Add(2)
	go func() {
		defer wg.Done()
		scanner := bufio.NewScanner(stdoutPipe)
		for scanner.Scan() {
			fmt.Println(scanner.Text())
		}
	}()
	go func() {
		defer wg.Done()
		scanner := bufio.NewScanner(stderrPipe)
		for scanner.Scan() {
			fmt.Println(scanner.Text())
		}
	}()
	err = cmd.Wait()
	wg.Wait()
	if err != nil {
		fmt.Println("Error:", err)
		cmd.Process.Kill()
		return
	}

	return true, cmd.ProcessState.ExitCode()
}
