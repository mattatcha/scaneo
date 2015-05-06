package testdata

import "time"

type Amet struct{}

func (a Amet) AmetFunc() int {
	return 0
}

type Fizz struct {
	Buzz *time.Time
}

type foo struct {
	bar time.Time
}
