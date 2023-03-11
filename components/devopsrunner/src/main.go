package main

import "os"

var version = "0.0.0"

func main() {
	binaryExecutor := NewBinaryExecutor()
	commandFactory := NewCommandFactory()
	configuration := NewConfiguration(os.Args)
	devopsRunner := NewDevopsRunner(version, binaryExecutor, commandFactory, configuration)
	devopsRunner.Run()
}
