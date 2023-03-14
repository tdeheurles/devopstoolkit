package main

import (
	"fmt"
	"strings"
)

type DevopsRunnerParameter struct {
	Name string
	Type string
	Tag  string
}

func NewDevopsRunnerParameter(line string) (*DevopsRunnerParameter, error) {
	cleanedLine := line
	for _, elt := range []string{"# do:", "#do:", "#do :", "# do :", " "} {
		if strings.HasPrefix(cleanedLine, elt) {
			cleanedLine = strings.Replace(cleanedLine, elt, "", 1)
		}
	}

	split := strings.SplitN(cleanedLine, " ", 3)
	if len(split) != 3 {
		return nil, fmt.Errorf("line %s is not a valid devopsrunner parameter", line)
	}

	return &DevopsRunnerParameter{
		Name: split[0],
		Type: split[1],
		Tag:  split[2],
	}, nil
}
