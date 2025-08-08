package main

import (
	"errors"
	"fmt"
)

type temporary interface {
	Temporary() bool
}

type temp struct {
	Err error
}

func (t temp) Temporary() bool {
	return true
}

func (t temp) Error() string {
	return t.Err.Error()
}

func foo() error {
	return fmt.Errorf("foo: %w", temp{Err: fmt.Errorf("temp error")})
}

func main() {
	err := foo()
	var t temporary
	if errors.As(err, &t) {
		if t.Temporary() {
			fmt.Println("this is temporary")
		}
	}
}
