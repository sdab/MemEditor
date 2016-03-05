package main

import(
	"fmt"
	mem "github.com/sdab/MemEditor/memlib"
)

type CmdHandler struct {
	scanner *mem.Scanner
	addresses []mem.ScanResult
}

func NewCmdHandler(pid int) *CmdHandler {
	scanner := mem.NewScanner(pid)
	return &CmdHandler{scanner, nil}
}

func (c *CmdHandler) Usage() string {
	// TODO: implement
	return "Usage: "
}

// %d>: shell where %d is the number of tracking addresses
// commands:
// scan val - read the current values of all tracked addresses and filter the tracked addresses by value
// list - list all tracked addresses and their last values
// update - read the current values of all tracked addresses
// set addr val
// setall val
// reset
// help
func (c *CmdHandler) HandleNextCmd() {
	fmt.Printf("%d>:", len(c.addresses))
	var cmd string
	fmt.Scan(&cmd)

	if cmd == "scan" {
		var val uint64
		fmt.Scan(&val)

		if c.addresses == nil {
			c.addresses = c.scanner.ScanAll(val)
		} else {
			c.addresses = c.scanner.ScanAndFilter(val, c.addresses)
		}
	} else if cmd == "list" {
		fmt.Println("\nAddress: value")
		for _, result := range c.addresses {
			fmt.Printf("%x: %d\n", result.Address, result.Value)
		}
	} else if cmd == "update" {
		c.addresses = c.scanner.Scan(c.addresses)
	} else if cmd == "set" {
		var addr uintptr
		fmt.Scan(&addr)
		var val uint64
		fmt.Scan(&val)
		
		c.scanner.WriteToAddress(addr, val)
		// update scan so we can see updated value
		c.addresses = c.scanner.Scan(c.addresses)
	} else if cmd == "setall" {
		var val uint64
		fmt.Scan(&val)

		for _, result := range c.addresses {
			c.scanner.WriteToAddress(result.Address, val)
		}
		// update scan so we can see updated value
		c.addresses = c.scanner.Scan(c.addresses)
	} else if cmd == "reset" {
		c.addresses = nil
	} else if cmd == "help" {
		fmt.Println(c.Usage())
	} else {
		fmt.Printf("Unknown command. %s\n", c.Usage())
	}


}
