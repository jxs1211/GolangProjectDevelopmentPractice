package main

import (
	"fmt"
	"os"
)

type S struct {
	a int
	b string
	c bool
}

func main() {
	// s := S{a: 1, b: "S", c: true}
	// fmt.Printf("%+v\n, %-v\n", s, s)
	// fmt.Printf("sasas")
	for i, v := range os.Environ() {
		fmt.Println(i, v)
	}
	// var err nil
	// var perr *os.PathError
	// errors.Is
	// errors.As(err, perr)
}
