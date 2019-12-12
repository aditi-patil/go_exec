package main

import (
	"fmt"
)

func Substract(a, b int) (d int) {
	d = a - b
	return d
}

func main() {
	fmt.Println(Substract(35, 12))
}
