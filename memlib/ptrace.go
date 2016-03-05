package memlib

import (
	"fmt"
	"encoding/binary"
	"golang.org/x/sys/unix"
)

func Wait(pid int) {
	var s unix.WaitStatus
	_, err := unix.Wait4(pid, &s, 0, nil)
	Check(err)
}

func Attach(pid int) {
	// try to attach
	err := unix.PtraceAttach(pid)
	Check(err)

	Wait(pid)
}

func Detach(pid int) {
	unix.PtraceDetach(pid)
}

func RecoverAndDetach(pid int) {
        if r := recover(); r != nil {
		fmt.Println("Recovered panic:", r)
        }
	Detach(pid)
}

func Peek(pid int, addr uintptr) (uint64, error) {
	data := make([]byte, LONGSIZE)
	_, err := unix.PtracePeekData(pid, addr, data)
	val := uint64(0)
	if err == nil {
		val = binary.LittleEndian.Uint64(data)
	}
	return val, err
}

func Poke(pid int, addr uintptr, val uint64) error {
	buf := make([]byte, LONGSIZE)
	binary.LittleEndian.PutUint64(buf, val)

	_, err := unix.PtracePokeData(pid, addr, buf)
	return err
}

