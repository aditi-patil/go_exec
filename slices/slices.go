package main

import "pic"

func Pic(dx, dy int) [][]uint8 {
	s := make([][]uint8, dy)
	for i := range s {
		s[i] = make([]uint8, dx)
	}
	for j, value := range s {
		value[j] = uint8(j * 2)
	}

	return s
}

func main() {
	pic.Show(Pic)
}
