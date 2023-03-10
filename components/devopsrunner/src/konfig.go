package main

import (
	"github.com/lalamove/konfig"
	"github.com/lalamove/konfig/loader/klenv"
	"github.com/lalamove/konfig/loader/klfile"
	"github.com/lalamove/konfig/loader/klflag"
	"github.com/lalamove/konfig/parser/kpyaml"
)

func loadAndWatchKonfig() {
	konfig.Init(konfig.DefaultConfig())
	bind()

	loadFiles()
	loadEnv()
	loadFlags()

	loadWatch()
}

func bind() {
	konfig.Bind(Config{})
}

func loadWatch() {
	if err := konfig.LoadWatch(); err != nil {
		println(err)
	}
}

func loadEnv() {
	konfig.RegisterLoader(
		klenv.New(&klenv.Config{
			Regexp: "^RUNNER_.*",
		}),
	)
}

func loadFiles() {
	konfig.RegisterLoaderWatcher(
		klfile.New(&klfile.Config{
			Files: []klfile.File{
				{
					Path:   "./runner.config.yaml",
					Parser: kpyaml.Parser,
				},
			},
			Watch: true,
		}),
		// optionally you can pass config hooks to run when a file is changed
		func(c konfig.Store) error {
			println(c.Name() + " have changed")
			return nil
		},
	)
}

func loadFlags() {
	konfig.RegisterLoader(
		klflag.New(&klflag.Config{}),
	)
}
