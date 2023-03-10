package main

import (
	"github.com/knadh/koanf"
	"github.com/knadh/koanf/parsers/yaml"
	"github.com/knadh/koanf/providers/file"
)

type Config struct {
	ApiVersion string `yaml:"apiVersion"`
	Kind       string `yaml:"kind"`
	Spec       struct {
		Devopsrunner struct {
			Version string `yaml:"version"`
		} `yaml:"devopsrunner"`
	} `yaml:"spec"`
}

type Configurationer interface {
	GetString(key string) string
}

type Configuration struct {
	store *koanf.Koanf
}

func NewConfiguration() *Configuration {
	configuration := &Configuration{
		store: koanf.New("."),
	}

	configuration.store.Load(file.Provider("./runner.config.yaml"), yaml.Parser())

	return configuration
}

func (c *Configuration) Store() *koanf.Koanf {
	return c.store
}

func (c *Configuration) GetString(key string) string {
	return c.store.String(key)
}
