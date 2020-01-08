package main

import (
	"fmt"
	"github.com/eiannone/keyboard"
)

func main() {
	char, _, err := keyboard.GetSingleKey()
	if err != nil {
		panic(err)
	}
	fmt.Printf("You pressed: %q\r\n", char)
}
