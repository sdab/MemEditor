package memlib

import (
	"fmt"
	"os"
	"bufio"
)

// A type that performs memory scans on a pid
type Scanner struct {
	pid int
	regions []VMARegion
}

// A virtual memory address region
type VMARegion struct {
	Begin, End uint64
}

type ScanResult struct {
	Address uintptr
	Value uint64
}

// Returns an initialized scanner
func NewScanner(pid int) *Scanner {
	return &Scanner{pid, ScanVirtualMemory(pid)}
}

// Scans the /proc/pid/maps file for the process' mapped 
// virtual memory
func ScanVirtualMemory(pid int) []VMARegion {
	filename := fmt.Sprintf("/proc/%d/maps", pid)
	f, err := os.Open(filename)
	Check(err)

	reader := bufio.NewReader(f)
	var regions []VMARegion

	fmt.Println("Reading mapped mem file.")
	for {
		// TODO: handle IsPrefix
		line, _, err := reader.ReadLine()
		if err != nil {
			break
		}
		var begin, end uint64
		fmt.Sscanf(string(line), "%x-%x", &begin, &end)
		regions = append(regions, VMARegion{begin, end})
	}
	return regions
}

// Scans all addresses and filters them by value.
func (s *Scanner) ScanAll(val uint64) []ScanResult {
	Attach(s.pid)
	defer RecoverAndDetach(s.pid)

	count := 0
	countErrors := 0
	results := make([]ScanResult, 0, 100)

	fmt.Println("Reading addresses.")
	for _, region := range s.regions {
		size := region.End - region.Begin

		fmt.Printf("Reading address region %x-%x with size %d.\n", region.Begin, region.End, size)

		for i := uint64(0); i < size/LONGSIZE; i++ {
			addr := uintptr(region.Begin+i*LONGSIZE)
			v, err := Peek(s.pid, addr)
			count++
			if err != nil {
				countErrors++
				continue
			}

			if v == val {
				results = append(results, ScanResult{addr, v})
			}
		}
	}
	fmt.Printf("Searched %d addresses with %d errors.\n", count, countErrors)
	return results
}

// Scans addresses in ScanResult array and returns an updated
// ScanResult array.
func (s *Scanner) Scan(addresses []ScanResult) []ScanResult {
	Attach(s.pid)
	defer RecoverAndDetach(s.pid)

	count := 0
	countErrors := 0
	results := make([]ScanResult, 0, 100)

	fmt.Println("Reading addresses.")
	for _, result := range addresses {
		r_val, err := Peek(s.pid, result.Address)
		count++
		if err != nil {
			countErrors++
			continue
		}
		
		results = append(results, ScanResult{result.Address, r_val})
	}
	return results
}

// Scans addresses in ScanResult array and returns those that are
// storing val.
func (s *Scanner) ScanAndFilter(val uint64, addresses []ScanResult) []ScanResult {
	results := s.Scan(addresses)
	out := make([]ScanResult, 0, 100)
	// filter results by value
	for _, result := range results {
		if result.Value == val {
			out = append(out, result)
		}
	}
	return out
}

// Writes the value to the address passed in.
// TODO: return error rather than panic
func (s *Scanner) WriteToAddress(addr uintptr, val uint64) {
	fmt.Printf("Writing value %d to address %x on pid %d\n", val, addr, s.pid)
	
	Attach(s.pid)
	defer RecoverAndDetach(s.pid)

	err := Poke(s.pid, addr, val)
	if err != nil {
		fmt.Println(s.pid)
		panic(err)
	}
}
