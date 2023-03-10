package main

type BinaryExecutorer interface {
	Execute(version string) (bool, int)
}

type BinaryExecutor struct {
}

func NewBinaryExecutor() *BinaryExecutor {
	return &BinaryExecutor{}
}

func (b BinaryExecutor) Execute(version string) (bool, int) {
	panic("BinaryExecutor.Execute() is not implemented")
}
