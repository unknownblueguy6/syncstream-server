package operator

type Adder interface {
	Add() int
}

type adder struct {
	x, y int
}

func NewAdder(x, y int) Adder {
	return &adder{
		x: x,
		y: y,
	}
}

func (a *adder) Add() int {
	return a.x + a.y
}
