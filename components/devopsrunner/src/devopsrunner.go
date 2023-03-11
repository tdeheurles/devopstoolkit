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

	stop, exitCode := d.SwitchBinaryIfNeeded()
	if stop {
		return exitCode
	}

	// GetCommonParameters()
	// TurnOnVerbose()

	// component.Assert()
	// command.Assert()
	if d.CommandFactory == nil {
		panic("DevopsRunner.CommandFactory is nil")
	}
	command := d.CommandFactory.Parse()
	if command == nil {
		panic("DevopsRunner.CommandFactory.Parse() returned nil")
	}

	return command.Execute()
}

func ParseCommand() {
	panic("ParseCommand is unimplemented")
}

func AssertCommand() {
	panic("AssertCommand is unimplemented")
}

func AssertComponent() {
	panic("AssertComponent is unimplemented")
}

func TurnOnVerbose() {
	panic("TurnOnVerbose is unimplemented")
}

func GetCommonParameters() {
	panic("GetCommonParameters is unimplemented")
}

func (d *DevopsRunner) SwitchBinaryIfNeeded() (bool, int) {
	expectedCommandVersion := d.Configuration.Data().Spec.Devopsrunner.Version
	if expectedCommandVersion == LocalDevelopmentVersionKey {
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
