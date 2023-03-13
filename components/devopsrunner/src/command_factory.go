package main

type CommandFactorier interface {
	Parse() (Commander, error)
}

type CommandFactory struct {
}

func NewCommandFactory() *CommandFactory {
	return &CommandFactory{}
}

func (c *CommandFactory) Parse() (Commander, error) {
	panic("CommandFactory.Parse() is not implemented")
}
