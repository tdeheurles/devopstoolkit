package main

type DevopsRunnerer interface {
	Run()
}

type DevopsRunner struct {
	binaryVersion  string
	binaryExecutor BinaryExecutorer
	commandFactory CommandFactorier
	configuration  Configurationer
}

func NewDevopsRunner(version string, binaryExecutor BinaryExecutorer, commandFactory CommandFactorier, configuration Configurationer) *DevopsRunner {
	if binaryExecutor == nil {
		panic("DevopsRunner.BinaryExecutor is nil")
	}

	if configuration == nil {
		panic("DevopsRunner.Configuration is nil")
	}

	if commandFactory == nil {
		panic("DevopsRunner.CommandFactory is nil")
	}

	return &DevopsRunner{
		binaryVersion:  version,
		binaryExecutor: binaryExecutor,
		configuration:  configuration,
		commandFactory: commandFactory,
	}
}

func (d *DevopsRunner) Run() int {

	stop, exitCode := d.SwitchBinaryIfNeeded()
	if stop {
		return exitCode
	}

	command, err := d.commandFactory.Parse(d.configuration)
	if err != nil {
		panic(err)
	}

	return command.Execute()
}

func (d *DevopsRunner) SwitchBinaryIfNeeded() (bool, int) {
	expectedCommandVersion := d.configuration.Data().Spec.Devopsrunner.Version
	if expectedCommandVersion == LocalDevelopmentVersionKey {
		return false, 0
	}

	if expectedCommandVersion == "" {
		return false, 0
	}

	if d.binaryVersion == expectedCommandVersion {
		return false, 0
	}

	return d.binaryExecutor.Execute(expectedCommandVersion)
}
