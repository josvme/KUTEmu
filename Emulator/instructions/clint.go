package instructions

import "errors"

const BASE_CLINT = 0x2000000
const CLINT_END = 0x200BFFF
const MTIMECMP_OFFSET = 0x4000
const MTIME_OFFSET = 0xBFF8

type Clint struct {
	Msip     uint32
	Mtimecmp uint64
	Mtime    uint64
}

func (c *Clint) Write(v uint32, addr uint32, cpu *Cpu) error {
	if addr > CLINT_END || addr < BASE_CLINT {
		return errors.New("Invalid address for clint")
	}

	if addr == BASE_CLINT {
		c.Msip = v
		// trigger MSIP (Software interrupt)
		// set MSIP (3) of MIP
		// We assume software interrupt
		c.TriggerInterrupt(3, cpu)
	}

	if addr == BASE_CLINT+MTIMECMP_OFFSET {
		// so clear everything except top 32 bits
		c.Mtimecmp = (c.Mtimecmp & uint64(0xFFFFFFFF<<32)) | uint64(v)
		// lower part of mtimecmp register
	}

	if addr == BASE_CLINT+MTIME_OFFSET+4 {
		// upper part of mtimecmp
		// so clear everything except bottom 32 bits
		c.Mtimecmp = (c.Mtimecmp & uint64(0xFFFFFFFF)) | uint64(v)<<32
	}

	if addr == BASE_CLINT+MTIME_OFFSET {
		// so clear everything except top 32 bits
		c.Mtime = (c.Mtime & uint64(0xFFFFFFFF<<32)) | uint64(v)
		// lower part of mtimecmp register
	}

	if addr == BASE_CLINT+MTIMECMP_OFFSET+4 {
		// upper part of mtimecmp
		// so clear everything except bottom 32 bits
		c.Mtime = (c.Mtime & uint64(0xFFFFFFFF)) | uint64(v)<<32
	}

	return nil
}

func (c *Clint) TriggerTimerInterrupt(cpu *Cpu) {
	if cpu.CSR.Registers[MIP] > 0 {
		return
	}
	// set mtip bit on position 7
	cpu.CSR.Registers[MIP] = cpu.CSR.Registers[MIP] | 1<<7
}

func (c *Clint) TriggerInterrupt(id uint32, cpu *Cpu) {
	if cpu.CSR.Registers[MIP] > 0 {
		return
	}
	cpu.CSR.Registers[MIP] = cpu.CSR.Registers[MIP] | 1<<id
}
