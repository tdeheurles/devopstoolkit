package main

import (
	"flag"
	"strings"

	"github.com/itzg/go-flagsfiller"
	"github.com/knadh/koanf"
	"github.com/knadh/koanf/parsers/yaml"
	"github.com/knadh/koanf/providers/file"
	"github.com/knadh/koanf/providers/posflag"
	"github.com/spf13/pflag"
)

var LocalDevelopmentVersionKey = "local"
var DevopsrunnerVersionKey = "spec.devopsrunner.version"

type Config struct {
	ApiVersion string `yaml:"apiVersion"`
	Kind       string `yaml:"kind"`
	Spec       struct {
		Devopsrunner struct {
			Version string `yaml:"version" default:"local" usage:"Version of the devopsrunner to use, can be 'local' or a version number"`
			Debug   bool   `yaml:"debug" default:"false" usage:"Enable debug mode"`
		} `yaml:"devopsrunner"`
		Command struct {
			Path []string `yaml:"path" usage:"Path to the command binaries"`
		} `yaml:"command"`
	} `yaml:"spec"`
}

type Configurationer interface {
	Data() Config
}

type Configuration struct {
	store      *koanf.Koanf
	data       Config
	args       []string
	configPath string
}

func NewConfiguration(args []string) *Configuration {

	pflagset := pflag.FlagSet{}
	configPath := pflagset.String("config", "runner.config.yaml", "Path to the configuration file (default: runner.config.yaml)")
	pflagset.ParseErrorsWhitelist.UnknownFlags = true
	pflagset.Parse(args)

	configuration := &Configuration{
		store:      koanf.New("."),
		data:       Config{},
		args:       args,
		configPath: *configPath,
	}

	parseConfigurationFiles(configuration)
	parseCommandLine(configuration)
	feelConfig(configuration)

	verbose(configuration)
	return configuration
}

func verbose(configuration *Configuration) {
	if configuration.data.Spec.Devopsrunner.Debug {
		println("--- configuration.store ---")
		configuration.store.Print()
		println("\n")
	}
}

func (c *Configuration) Data() Config {
	return c.data
}

func feelConfig(configuration *Configuration) {
	// We fill the configuration.Data structure
	// see: https://github.com/knadh/koanf#unmarshalling-and-marshalling
	err := configuration.store.Unmarshal("", &configuration.data)
	if err != nil {
		panic("error unmarshalling config --> " + err.Error())
	}
}

func parseConfigurationFiles(configuration *Configuration) {
	configuration.store.Load(file.Provider(configuration.configPath), yaml.Parser())
	// println("configuration after runner.config.yaml:")
	// configuration.store.Print()
	// println("")
}

func parseCommandLine(configuration *Configuration) {
	// We define all the flags from the Config structure
	// see: https://github.com/itzg/go-flagsfiller
	var config Config
	filler := flagsfiller.New()
	flagset := &flag.FlagSet{}
	err := filler.Fill(flagset, &config)
	if err != nil {
		panic(err)
	}

	// We transfert the flagset (go std flag) to a pflagset (spf13/pflag)
	// see: https://github.com/spf13/pflag
	pflagset := &pflag.FlagSet{}
	pflagset.SetNormalizeFunc(wordSepNormalizeFunc)
	pflagset.AddGoFlagSet(flagset)

	// We parse the pflagset
	pflagset.Parse(configuration.args)

	// We load the pflagset into the koanf store
	// see: https://github.com/knadh/koanf#reading-from-command-line
	if err := configuration.store.Load(posflag.Provider(pflagset, ".", configuration.store), nil); err != nil {
		panic("error loading config from CLI --> " + err.Error())
	}

	// We print the configuration
	// println("configuration after runner.config.yaml and command line:")
	// configuration.store.Print()
	// println("")
}

// wordSepNormalizeFunc changes all flags from camel case to dot notation
func wordSepNormalizeFunc(f *pflag.FlagSet, name string) pflag.NormalizedName {
	from := []string{"-", "_"}
	to := "."
	for _, sep := range from {
		name = strings.Replace(name, sep, to, -1)
	}
	return pflag.NormalizedName(name)
}
