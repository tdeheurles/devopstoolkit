package main

import "strings"

type CommandFactorier interface {
	Parse(configuration Configurationer) (Commander, error)
}

type CommandFactory struct {
}

func NewCommandFactory() *CommandFactory {
	return &CommandFactory{}
}

func (c *CommandFactory) Parse(Configuration Configurationer) (Commander, error) {
	panic("CommandFactory.Parse() is not implemented")
}

func (c *CommandFactory) GetDevopsRunnerParameters(content string) ([]DevopsRunnerParameter, error) {

	lines := strings.Split(content, "\n")
	devopsRunnerParameters := []DevopsRunnerParameter{}
	prefixes := []string{"# do:", "#do:", "#do :", "# do :"}
	for _, line := range lines {
		for _, prefix := range prefixes {
			if strings.HasPrefix(line, prefix) {
				devopsRunnerParameter, err := NewDevopsRunnerParameter(line)
				if err != nil {
					return nil, err
				}
				devopsRunnerParameters = append(devopsRunnerParameters, *devopsRunnerParameter)
			}
		}
	}

	return devopsRunnerParameters, nil
}
