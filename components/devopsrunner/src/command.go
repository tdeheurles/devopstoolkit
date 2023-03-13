package main

type Commander interface {
	Assert()
	Execute() int
	Parse()
}

type Command struct {
}

func NewCommand() *Command {
	return &Command{}
}

func (c *Command) Assert() {
	panic("Command.Assert() is not implemented")
}

func (c *Command) Execute() int {
	panic("Command.Execute() is not implemented")
}

func (c *Command) Parse() {
	panic("Command.Parse() is not implemented")
}
