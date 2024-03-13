package main

import (
	"fmt"
	"go-template/pkg/calculator"
)

func main() {
	fmt.Println("Hello World")

	c := calculator.NewCalculator(5, 5)
	fmt.Println("5 + 5 = ", c.GetAnswer())
}
