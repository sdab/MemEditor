package memlib

import (
	"fmt"
	"os"
	"bufio"
)

// A type that performs memory scans on a pid
type Scanner struct {
	pid int
	addresses map[uintptr]int
}

// A virtual memory address region
type VMARegion struct {
	Begin, End uint64
}

// Returns an initialized scanner
func NewScanner(pid int) *Scanner {
	return &Scanner{pid, nil}
}

// Scans the /proc/pid/maps file for the process' mapped 
// virtual memory
func (s *Scanner) ScanVirtualMemory() []VMARegion {
	filename := fmt.Sprintf("/proc/%d/maps", s.pid)
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

// Returns the addresses that store val
func (s *Scanner) SearchAddresses(val int) []uintptr {
	Attach(s.pid)
	defer RecoverAndDetach(s.pid)

	if s.addresses == nil {
		// scan virtual memory of pid
		regions := s.ScanVirtualMemory()
		fmt.Printf("Read %d memory regions.\n", len(regions))
		s.addresses = make(map[uintptr]int)
		// search for all addressess with val
		return s.searchAllAddresses(val, regions)
	} else {
		return s.searchExistingAddresses(val)
	}
}

func (s *Scanner) searchAllAddresses(val int, regions []VMARegion) []uintptr {
	count := 0
	countErrors := 0
	out_ptrs := make([]uintptr, 0, 100)

	fmt.Println("Reading addresses.")
	for _, region := range regions {
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

			if v == uint64(val) {
				s.addresses[addr] = val
				out_ptrs = append(out_ptrs, addr)
			}
		}
	}
	fmt.Printf("Searched %d addresses with %d errors.\n", count, countErrors)
	return out_ptrs
}

func (s *Scanner) searchExistingAddresses(val int) []uintptr {
	count := 0
	countErrors := 0
	out_ptrs := make([]uintptr, 0, 100)

	fmt.Println("Reading addresses.")
	for addr, _ := range s.addresses {
		r_val, err := Peek(s.pid, addr)
		count++
		if err != nil {
			countErrors++
			continue
		}
		
		if r_val != uint64(val) {
			delete(s.addresses, addr)
		} else {
			s.addresses[addr] = val
			out_ptrs = append(out_ptrs, addr)
		}
	}
	return out_ptrs
}

// Writes the value to the address passed in.
// TODO: return error rather than panic
func (s *Scanner) WriteToAddress(addr uintptr, val int) {
	fmt.Printf("Writing value %d to address %x on pid %d\n", val, addr, s.pid)
	
	Attach(s.pid)
	defer RecoverAndDetach(s.pid)

	err := Poke(s.pid, addr, uint64(val))
	if err != nil {
		fmt.Println(s.pid)
		panic(err)
	}
}
