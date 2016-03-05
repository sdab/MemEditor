package main

import (
	"bufio"
	"encoding/binary"
	"flag"
	"fmt"
	"golang.org/x/sys/unix"
	"os"
	"unsafe"
)

func readNextValue() int {
	fmt.Printf("Please insert a value to search:\n")
	var val int
	fmt.Scan(&val)
	return val
}

type VMARegion struct {
	Begin, End uint64
}

// Attempt to read proc memory map
func ScanMemory(pid int) []VMARegion {
	filename := fmt.Sprintf("/proc/%d/maps", pid)
	f, err := os.Open(filename)
	check(err)

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

var addresses map[uintptr]int

const LONGSIZE = uint64(unsafe.Sizeof(uint64(0)))

func check(err error) {
	if err != nil {
		panic(err)
	}
}

func wait(pid int) {
	var s unix.WaitStatus
	_, err := unix.Wait4(pid, &s, 0, nil)
	check(err)
}

func attach(pid int) {
	// try to attach
	err := unix.PtraceAttach(pid)
	check(err)

	wait(pid)
}

func detach(pid int) {
	err := unix.PtraceDetach(pid)
	check(err)
}

func peek(pid int, addr uintptr) (uint64, error) {
	data := make([]byte, LONGSIZE)
	_, err := unix.PtracePeekData(pid, addr, data)
	val := uint64(0)
	if err == nil {
		val = binary.LittleEndian.Uint64(data)
	}
	return val, err
}

func poke(pid int, addr uintptr, val uint64) error {
	buf := make([]byte, LONGSIZE)
	binary.LittleEndian.PutUint64(buf, val)

	_, err := unix.PtracePokeData(pid, addr, buf)
	return err
}

func recoverAndDetach(pid int) {
        if r := recover(); r != nil {
		fmt.Println("Recovered panic:", r)
        }
	detach(pid)
}

func searchAddresses(pid, val int) {
	attach(pid)
	defer recoverAndDetach(pid)

	if addresses == nil {
		// scan virtual memory of pid
		regions := ScanMemory(pid)
		fmt.Printf("Read %d memory regions.\n", len(regions))
		addresses = make(map[uintptr]int)
		// search for all addressess with val
		searchAllAddresses(pid, val, regions)
	} else {
		searchExistingAddresses(pid, val)
	}
}

func searchAllAddresses(pid, val int, regions []VMARegion) {
	count := 0
	countErrors := 0

	fmt.Println("Reading addresses.")
	for _, region := range regions {
		size := region.End - region.Begin

		fmt.Printf("Reading address region %x-%x with size %d.\n", region.Begin, region.End, size)

		for i := uint64(0); i < size/LONGSIZE; i++ {
			addr := uintptr(region.Begin+i*LONGSIZE)
			v, err := peek(pid, addr)
			count++
			if err != nil {
				countErrors++
				continue
			}

			if v == uint64(val) {
				addresses[addr] = val

			}
		}
	}
	fmt.Printf("Searched %d addresses with %d errors.\n", count, countErrors)
}

func searchExistingAddresses(pid, val int) {
	count := 0
	countErrors := 0

	fmt.Println("Reading addresses.")
	for addr, _ := range addresses {
		r_val, err := peek(pid, addr)
		count++
		if err != nil {
			countErrors++
			continue
		}
		
		if r_val != uint64(val) {
			delete(addresses, addr)
		} else {
			addresses[addr] = val
		}
	}
}

func WriteToAddress(pid int, addr uintptr, val int) {
	fmt.Printf("Writing value %d to address %x on pid %d\n", val, addr, pid)
	
	attach(pid)
	defer recoverAndDetach(pid)

	err := poke(pid, addr, uint64(val))
	if err != nil {
		fmt.Println(pid)
		panic(err)
	}	
}

func main() {
	// TODO: be able to search pids
	// read pid
	pid := flag.Int("pid", 0, "Pid of the process to inspect.")
	flag.Parse()

	if *pid == 0 {
		panic("Must input pid.")
	}

	for ;; {
		// read value
		val := readNextValue()
		fmt.Printf("Searching for %d:\n", val)

		// search pid memory for addresses with value
		searchAddresses(*pid, val)

		// report number, if < 10 print them otherwise narrow down
		// by looping on searching with value
		fmt.Printf("Found %d matching addresses.\n", len(addresses))
		if len(addresses) < 10 {
			for i, v := range addresses {
				fmt.Printf("Address %x : val %d.\n", i, v)
			}
		}

		if len(addresses) <= 1 {
			break
		}
	}

	if len(addresses) == 0 {
		fmt.Println("No addresses matched that value.")
		return
	}

	// change value
	for addr,_ := range addresses {
		fmt.Printf("Change address %d to value:\n", addr)
		var new_val int
		fmt.Scan(&new_val)
		WriteToAddress(*pid, addr, new_val)
	}
}
