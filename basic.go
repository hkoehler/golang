package main

import "fmt"

type T struct {
	i int
	s string
}

func (t *T) String() string {
	return fmt.Sprintf("(i=%v, s=%q)", t.i, t.s)	
}

func main() {
	var t T
	fmt.Printf("%v\n", t)
}
