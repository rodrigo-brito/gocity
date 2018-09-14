package test

import "fmt"

var (
	a = 1
	b = 2
)

const batata = 3

type AliasExample int

type Example struct {
	Value string
}

func (e *Example) Print() {
	fmt.Println(e.Value)
}

func (e *Example) Print2() {
	fmt.Println(e.Value)
}

func Print(values ...string) {
	for value := range values {
		fmt.Println(" ", value)
	}
}

func Print2(values ...string) {
	for value := range values {
		fmt.Println(" ", value)
	}
}
