package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewDevopsRunnerParameterParsing(t *testing.T) {
	// given
	line := "# do: Name string `json:\"name\"`"

	// when
	devopsRunnerParameter, err := NewDevopsRunnerParameter(line)
	if err != nil {
		t.Fatal(err)
	}

	// then
	assert.Equal(t, "Name", devopsRunnerParameter.Name)
	assert.Equal(t, "string", devopsRunnerParameter.Type)
	assert.Equal(t, "`json:\"name\"`", devopsRunnerParameter.Tag)
}
