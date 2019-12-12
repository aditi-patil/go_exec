package main

import (
	"fmt"
	"strings"
)

var s = "This is a text string a text which is having too long string"

func main() {
	m := make(map[string]int)
	a := strings.Fields(s)
	for i := 0; i < len(a); i++ {
		v, ok := m[a[i]]
		if ok == true {
			m[a[i]] = v + 1
		} else {
			m[a[i]] = 1
		}
	}
	fmt.Println(m)
}
