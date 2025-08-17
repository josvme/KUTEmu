package instructions

import (
	"sync"
)

type Display struct {
	Screen map[uint32]uint32
	Mutex  sync.Mutex
}

func (d *Display) Write(v uint32, add uint32) error {
	d.Mutex.Lock()
	addr := add - VIRT_DISPLAY
	d.Screen[addr] = v
	//fmt.Printf("Writing to display %X, %X\n", addr, v)
	d.Mutex.Unlock()
	return nil
}
