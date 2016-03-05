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

// Returns a help prompt describing the commands that CmdHandler
// can take as input.
func (c *CmdHandler) Usage() string {
	help := "Prompt:\n" +
		"%d>: prompt where %d is the number of tracking " +
		"addresses.\n" +
		"Commands:\n" +
		"* scan val - scans the current values of all " +
		"tracked addresses and filters the tracked " +
		"addresses by value. Scans the whole of mapped " +
		"memory if there are no tracked addresses (such " +
		"as on startup or after a reset).\n" +
		"* list - lists all tracked addresses and their " +
		"last values.\n" +
		"* update - scans the current values of all " +
		"tracked addresses.\n" +
		"* set addr val - Writes val to address addr.\n" +
		"* setall val - Writes val to all tracked " +
		"addresses.\n" +
		"* reset - Removes all tracked addresses. The " +
		"next scan will read all of mapped memory.\n" +
		"* help - prints the commands\n"
        
	return help
}

// %d>: shell where %d is the number of tracking addresses
// commands:
//  scan val - scans the current values of all tracked addresses and
//   filters the tracked addresses by value. Scans the whole of 
//   mapped memory if there are no tracked addresses (such as on 
//   startup or after a reset).
//
//  list - lists all tracked addresses and their last values.
//
//  update - scans the current values of all tracked addresses.
//
//  set addr val - Writes val to address addr.
//
//  setall val - Writes val to all tracked addresses.
//
//  reset - Removes all tracked addresses. The next scan will read
//   all of mapped memory.
//
//  help - prints the commands
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
		fmt.Printf("Unknown command.\nHelp text:\n%s", c.Usage())
	}
}
