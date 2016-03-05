package main
// A simple toy program with a value variable that continously
// reads in values from the user and stores it

import (
	"fmt"
)

var kValue int

func readNextValue() int {
	fmt.Printf("Please insert a value to store:\n")
	var val int
	fmt.Scan(&val)
	return val
}

func usage() string {
	return "Commands:\nSet 'n'\nRead"
}

func readNextCmd() {
	fmt.Printf(">:")
	var cmd string
	fmt.Scan(&cmd)

	if cmd == "Read" {
		fmt.Printf("Stored value is %d.\n", kValue)
	} else if cmd == "Set" {
		var val int
		fmt.Scan(&val)
		kValue = val
	} else {
		fmt.Println(usage())
	}
}

func main() {
	for {
		readNextCmd()
	}
}
