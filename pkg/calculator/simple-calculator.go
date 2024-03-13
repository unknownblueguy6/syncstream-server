package calculator

import "go-template/pkg/operator"

type Calculator interface {
	GetAnswer() int
}

type simpleClr struct {
	adr operator.Adder
}

func NewCalculator(x, y int) Calculator {
	return &simpleClr{
		adr: operator.NewAdder(x, y),
	}
}

func (s *simpleClr) GetAnswer() int {
	return s.adr.Add()
}
