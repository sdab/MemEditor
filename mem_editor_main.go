package main

import (
	"flag"
	"fmt"
)

func readNextValue() int {
	fmt.Printf("Please insert a value to search:\n")
	var val int
	fmt.Scan(&val)
	return val
}

func main() {
	// TODO: be able to search pids
	// read pid
	pid := flag.Int("pid", 0, "Pid of the process to inspect.")
	flag.Parse()

	if *pid == 0 {
		panic("Must input pid.")
	}

	cmdHandler := NewCmdHandler(*pid)

	for ;; {
		cmdHandler.HandleNextCmd()
		// // read value
		// val := readNextValue()
		// fmt.Printf("Searching for %d:\n", val)

		// // search pid memory for addresses with value
		// addresses = scanner.SearchAddresses(val)

		// // report number, if < 10 print them otherwise narrow down
		// // by looping on searching with value
		// fmt.Printf("Found %d matching addresses.\n", len(addresses))
		// if len(addresses) < 10 {
		// 	for _, addr := range addresses {
		// 		fmt.Printf("Address %x : val %d.\n", addr, val)
		// 	}
		// }

		// if len(addresses) <= 1 {
		// 	break
		// }
	}

	// if len(addresses) == 0 {
	// 	fmt.Println("No addresses matched that value.")
	// 	return
	// }

	// // change value
	// for _, addr := range addresses {
	// 	fmt.Printf("Change address %d to value:\n", addr)
	// 	var new_val int
	// 	fmt.Scan(&new_val)
	// 	scanner.WriteToAddress(addr, new_val)
	// }
}
