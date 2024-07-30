package main

import (
	"fmt"
	"EX_2/stringunpack" // Импортируем наш пакет
)

func main() {
	inputs := []string{
		"a4bc2d2e",
		"abcd",
		"45",
		"",
		"qwe\\4\\5",
		"qwe\\45",
		"qwe\\\\5",
		"abc\\",
	}

	for _, input := range inputs {
		result, err := stringunpack.UnpackString(input)
		if err != nil {
			fmt.Printf("Error unpacking string %q: %v", input, err)
		} else {
			fmt.Printf("String %q: %q\n", input, result)
		}
	}
}
