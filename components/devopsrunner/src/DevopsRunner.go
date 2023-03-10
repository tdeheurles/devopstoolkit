package main

type DevopsRunnerer interface {
	Run()
}

type DevopsRunner struct {
	BinaryVersion  string
	BinaryExecutor BinaryExecutorer
	CommandFactory CommandFactorier
	Configuration  Configurationer
}

func NewDevopsRunner(version string, binaryExecutor BinaryExecutorer, commandFactory CommandFactorier, configuration Configurationer) *DevopsRunner {
	return &DevopsRunner{
		BinaryVersion:  version,
		BinaryExecutor: binaryExecutor,
		Configuration:  configuration,
		CommandFactory: commandFactory,
	}
}

func (d *DevopsRunner) Run() int {

	// d.ConfigurationStore

	stop, exitCode := d.SwitchBinaryIfNeeded()
	if stop {
		return exitCode
	}

	// GetCommonParameters()
	// TurnOnVerbose()

	// component.Assert()
	// command.Assert()
	if d.CommandFactory == nil {
		panic("commandFactory is nil")
	}
	command := d.CommandFactory.Parse()
	if command == nil {
		panic("command is nil")
	}

	return command.Execute()
}

func ParseCommand() {
	panic("unimplemented")
}

func AssertCommand() {
	panic("unimplemented")
}

func AssertComponent() {
	panic("unimplemented")
}

func TurnOnVerbose() {
	panic("unimplemented")
}

func GetCommonParameters() {
	panic("unimplemented")
}

func (d *DevopsRunner) SwitchBinaryIfNeeded() (bool, int) {
	expectedCommandVersion := d.Configuration.GetString("spec.devopsrunner.version")
	if expectedCommandVersion == "local" {
		return false, 0
	}

	if expectedCommandVersion == "" {
		return false, 0
	}

	if d.BinaryVersion == expectedCommandVersion {
		return false, 0
	}

	return d.BinaryExecutor.Execute(expectedCommandVersion)
}
