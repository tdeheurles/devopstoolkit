package main

var version = "0.0.0"

func main() {
	binaryExecutor := NewBinaryExecutor()
	commandFactory := NewCommandFactory()
	configuration := NewConfiguration()
	devopsRunner := NewDevopsRunner(version, binaryExecutor, commandFactory, configuration)
	devopsRunner.Run()
}
