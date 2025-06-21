package instructions

import (
	"fmt"
	"os"
)

const THR = 0
const RBR = 0
const DLL = 0
const IER = 1
const DLM = 1
const IIR = 2
const FCR = 2
const LCR = 3
const MCR = 4
const LSR = 5
const MSR = 6
const SCR = 7

const DLAB_FLAG = 1 << 7

type UART struct {
	registerRT [8]byte
	registerDL [3]byte
}

func (u *UART) getDLabFlag() byte {
	return u.registerRT[LCR] & DLAB_FLAG
}
func (u *UART) interruptsEnabled() byte {
	return u.registerRT[IER] & 0x1
}

func NewUART() UART {
	return UART{
		// We are ready to transmit and receive
		registerRT: [8]byte{0, 0, 0, 0, 0, 32, 0, 0},
		registerDL: [3]byte{0, 0, 0},
	}
}

func (u *UART) Write(b byte, location uint32) error {
	// This is receive mode. This is receive because we receive it and write to stdout
	if location == THR && u.getDLabFlag() == 0 {
		fmt.Printf("%c", b)
		// show you can write more data
		// we write to the register, so the register can be completely used by read buffer
		u.registerRT[5] = 32
		return nil
	}

	if (u.getDLabFlag() == 0) || (location >= 3) {
		u.registerRT[location] = b
	}

	if u.getDLabFlag() == 1 && location <= 2 {
		u.registerDL[location] = b
	}

	return nil
}

// TODO
func (u *UART) Read(location uint32) (byte, error) {
	if location == LSR {
		// check stdin if there is some bytes to read, if there is nothing to read simply return nil
		// Check if there's data available to read from stdin
		stat, err := os.Stdin.Stat()
		if err != nil {
			return u.registerRT[LSR], err
		}

		// Check if there's data in the stdin buffer
		if (stat.Mode()&os.ModeCharDevice) == 0 && stat.Size() > 0 {
			// Data is available to read
			// Set bit 0 (DR - Data Ready) to 1
			u.registerRT[LSR] |= 0x1
		}

		return u.registerRT[LSR], nil
	}

	if location == RBR && u.getDLabFlag() == 0 {
		b := make([]byte, 1)
		_, err := os.Stdin.Read(b)
		if err != nil {
			return 0, err
		}
		u.registerRT[RBR] = b[0]
		// unsets data is there bit
		u.registerRT[5] = u.registerRT[5] & 0xFE
		return b[0], nil
	}

	if u.getDLabFlag() == 1 && location <= 2 {
		return u.registerDL[location], nil
	}

	return u.registerRT[location], nil
}
