package main

import (
	"fmt"
)

func main() {
	var x *int
	x = &3

	fmt.Println(*x)
}
