package main

type CommandFactorier interface {
	Parse() Commander
}

type CommandFactory struct {
}

func NewCommandFactory() *CommandFactory {
	return &CommandFactory{}
}

func (c *CommandFactory) Parse() Commander {
	panic("CommandFactory.Parse() is not implemented")
}
