package main

import (
	"fmt"
	"log"
	"os"

	"github.com/knadh/koanf"
	"github.com/knadh/koanf/parsers/yaml"
	"github.com/knadh/koanf/providers/file"
	"github.com/knadh/koanf/providers/posflag"

	flag "github.com/spf13/pflag"
)

func loadAndWatchKoanf() {

	var k = koanf.New(".")
	f := flag.NewFlagSet("foo", flag.ContinueOnError)
	f.Usage = func() {
		fmt.Println(f.FlagUsages())
		os.Exit(0)
	}
	f.Parse(os.Args[1:])

	k.Load(file.Provider("./runner.config.yaml"), yaml.Parser())

	if err := k.Load(posflag.Provider(f, ".", k), nil); err != nil {
		log.Fatalf("error loading config: %v", err)
	}

	k.Print()
}
