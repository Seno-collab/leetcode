package main

import "fmt"

func main() {
	s := "AZaz09?"
	for i := range s {
		fmt.Println(checkCharacters(s[i]))
	}
	fmt.Printf("%T\n", s[1])
}

func checkCharacters(v uint8) bool {
	fmt.Println(v)
	if v >= 65 && v <= 90 {
		return true
	}
	if v >= 97 && v <= 122 {
		return true
	}
	if v >= 48 || v <= 57 {
		return true
	}
	return false
}
