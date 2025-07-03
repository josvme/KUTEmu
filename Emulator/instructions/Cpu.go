package instructions

import (
	"fmt"
	"os"
)

type Cpu struct {
	PC             uint32
	Registers      [32]uint32
	Memory         Memory
	CSR            CSR
	AtomicReserved bool
	// 3 for machine, 1 for supervisor, 2 for hypervisor, 0 for user
	CurrentMode uint32
}

func (c *Cpu) ExecInst(i Inst) error {
	// Always reset register 0 to 0, to be sure
	c.Registers[0] = 0
	switch i.(type) {
	case RI:
		ins := i.(RI)
		executeR(ins, c)
	case II:
		ins := i.(II)
		executeI(ins, c)
	case SI:
		ins := i.(SI)
		executeS(ins, c)
	case BI:
		ins := i.(BI)
		executeB(ins, c)
	case JI:
		ins := i.(JI)
		executeJ(ins, c)
	case UI:
		ins := i.(UI)
		executeU(ins, c)
	case FI:
		ins := i.(FI)
		executeF(ins, c)
	}
	return nil
}

func (cpu *Cpu) HandleInterrupts(inst string) error {
	if inst == "wfi" && cpu.CSR.Registers[MIP] != 1 {
		cpu.CSR.Registers[MIP] = 1

		csr := cpu.CSR
		mode := cpu.CurrentMode
		// set mcause & scause register
		switch mode {
		case 1:
			csr.Registers[SCAUSE] = csr.Registers[SCAUSE] | 1<<31 | 1 // supervisor software interrupt
		default:
			csr.Registers[MCAUSE] = csr.Registers[MCAUSE] | 1<<31 | 3 // machine software interrupt
		}
		return nil
	}
	// Implement only direct mode now

	// Handle interrupts
	// Save pc to mepc / corresponding
	// move PC to interrupt / exception handler after checking mdeleg and meip / corresponding registers
	// mtvec has different address based on modes see page 24, https://people.eecs.berkeley.edu/~krste/papers/riscv-priv-spec-1.7.pdf
	// A trap in privilege level P causes a jump to the address mtvec + P ×0x40. Non-maskable interrupts
	//cause a jump to address mtvec + 0xFC.

	if cpu.CSR.Registers[MIP] == 1 {
		csr := cpu.CSR
		mode := cpu.CurrentMode

		// check if delegation is added for exceptions
		midelegReg := csr.Registers[MIDELEG]
		// This should be checked first
		if midelegReg&uint32(1<<3) == 1 {
			// go to supervisor trap
			csr.Registers[SEPC] = cpu.PC
			sstatus := ToMStatusReg(SSTATUS)
			sstatus.spp = mode
			csr.Registers[SSTATUS] = FromMStatusReg(sstatus)
			// This is very crappy, to accomodate +4
			stvec := ToMtvecReg(csr.Registers[STVEC])
			cpu.PC = stvec.base
			cpu.CurrentMode = 1
			// These are just wrong
			cpu.CSR.Registers[SIP] = 0
			cpu.CSR.Registers[MIP] = 0
		} else {
			// run machine trap
			csr.Registers[MEPC] = cpu.PC
			mstatus := ToMStatusReg(MSTATUS)
			mstatus.mpp = mode
			csr.Registers[MSTATUS] = FromMStatusReg(mstatus)
			// This is very crappy, to accomodate +4
			mtvec := ToMtvecReg(csr.Registers[MTVEC])
			cpu.PC = mtvec.base
			cpu.CurrentMode = 3
			// only for interrupts, software
			cpu.PC = (mtvec.base + 4*3) - 4
			cpu.CSR.Registers[MIP] = 0
		}

	}
	return nil
}

func executeF(inst FI, c *Cpu) {
	// Order memory/instruction access. We ignore them now.
	switch inst.Operation() {
	case "fence":
		c.PC += 4
	case "fence.i":
		c.PC += 4
	}
}

func executeR(inst RI, c *Cpu) {
	switch inst.Operation() {
	case "add":
		c.Registers[inst.RD] = c.Registers[inst.RS1] + c.Registers[inst.RS2]
		c.PC += 4

	case "sub":
		c.Registers[inst.RD] = c.Registers[inst.RS1] - c.Registers[inst.RS2]
		c.PC += 4

	case "xor":
		c.Registers[inst.RD] = c.Registers[inst.RS1] ^ c.Registers[inst.RS2]
		c.PC += 4

	case "or":
		c.Registers[inst.RD] = c.Registers[inst.RS1] | c.Registers[inst.RS2]
		c.PC += 4

	case "and":
		c.Registers[inst.RD] = c.Registers[inst.RS1] & c.Registers[inst.RS2]
		c.PC += 4

	case "sll":
		c.Registers[inst.RD] = c.Registers[inst.RS1] << (c.Registers[inst.RS2] & 0x1F)
		c.PC += 4

	case "srl":
		c.Registers[inst.RD] = c.Registers[inst.RS1] >> (c.Registers[inst.RS2] & 0x1F)
		c.PC += 4

	// Arithmetic Left shift RS1 by lower 5 bits of RS2
	case "sra":
		c.Registers[inst.RD] = uint32(int32(c.Registers[inst.RS1]) >> byte(c.Registers[inst.RS2]&0x1F))
		c.PC += 4

	// Signed compare
	case "slt":
		if int32(c.Registers[inst.RS1]) < int32(c.Registers[inst.RS2]) {
			c.Registers[inst.RD] = 1
		} else {
			c.Registers[inst.RD] = 0
		}
		c.PC += 4

	// Unsigned compare
	case "sltu":
		if c.Registers[inst.RS1] < c.Registers[inst.RS2] {
			c.Registers[inst.RD] = 1
		} else {
			c.Registers[inst.RD] = 0
		}
		c.PC += 4
	// Atomic Instructions
	case "lr.w":
		c.Registers[inst.RD] = c.Memory.ReadWord(c.Registers[inst.RS1])
		c.AtomicReserved = true
		c.PC += 4
	case "sc.w":
		if c.AtomicReserved {
			c.Memory.WriteWord(c.Registers[inst.RS2], c.Registers[inst.RS1])
			c.Registers[inst.RD] = 0
		} else {
			c.Registers[inst.RD] = 1
		}
		c.AtomicReserved = false
		c.PC += 4
	case "amoswap.w":
		c.Registers[inst.RD] = c.Memory.ReadWord(c.Registers[inst.RS1])
		c.Memory.WriteWord(c.Registers[inst.RS2], c.Registers[inst.RS1])
		c.PC += 4
	case "amoadd.w":
		c.Registers[inst.RD] = c.Memory.ReadWord(c.Registers[inst.RS1])
		c.Memory.WriteWord(c.Registers[inst.RS2]+c.Registers[inst.RD], c.Registers[inst.RS1])
		c.PC += 4
	case "amoand.w":
		c.Registers[inst.RD] = c.Memory.ReadWord(c.Registers[inst.RS1])
		c.Memory.WriteWord(c.Registers[inst.RS2]&c.Registers[inst.RD], c.Registers[inst.RS1])
		c.PC += 4
	case "amoor.w":
		c.Registers[inst.RD] = c.Memory.ReadWord(c.Registers[inst.RS1])
		c.Memory.WriteWord(c.Registers[inst.RS2]|c.Registers[inst.RD], c.Registers[inst.RS1])
		c.PC += 4
	case "amoxor.w":
		c.Registers[inst.RD] = c.Memory.ReadWord(c.Registers[inst.RS1])
		c.Memory.WriteWord(c.Registers[inst.RS2]^c.Registers[inst.RD], c.Registers[inst.RS1])
		c.PC += 4
	case "amomax.w":
		c.Registers[inst.RD] = c.Memory.ReadWord(c.Registers[inst.RS1])
		c.Memory.WriteWord(uint32(max(int32(c.Registers[inst.RS2]), int32(c.Registers[inst.RD]))), c.Registers[inst.RS1])
		c.PC += 4
	case "amomin.w":
		c.Registers[inst.RD] = c.Memory.ReadWord(c.Registers[inst.RS1])
		c.Memory.WriteWord(uint32(min(int32(c.Registers[inst.RS2]), int32(c.Registers[inst.RD]))), c.Registers[inst.RS1])
		c.PC += 4
	case "amomaxu.w":
		c.Registers[inst.RD] = c.Memory.ReadWord(c.Registers[inst.RS1])
		c.Memory.WriteWord(max(c.Registers[inst.RS2], c.Registers[inst.RD]), c.Registers[inst.RS1])
		c.PC += 4
	case "amominu.w":
		c.Registers[inst.RD] = c.Memory.ReadWord(c.Registers[inst.RS1])
		c.Memory.WriteWord(min(c.Registers[inst.RS2], c.Registers[inst.RD]), c.Registers[inst.RS1])
		c.PC += 4
		// Multiply Instructions
	case "mul":
		c.Registers[inst.RD] = uint32(int32(c.Registers[inst.RS1]) * int32(c.Registers[inst.RS2]))
		c.PC += 4
	case "mulh":
		c.Registers[inst.RD] = uint32(int64(int32(c.Registers[inst.RS1])) * int64(int32(c.Registers[inst.RS2])) >> 32)
		c.PC += 4
	case "mulhsu":
		// RS2 is unsigned and RS1 is signed
		c.Registers[inst.RD] = uint32(uint64(int32(c.Registers[inst.RS1])) * uint64(c.Registers[inst.RS2]) >> 32)
		c.PC += 4
	case "mulhu":
		c.Registers[inst.RD] = uint32(uint64(c.Registers[inst.RS1]) * uint64(c.Registers[inst.RS2]) >> 32)
		c.PC += 4
	case "div":
		c.Registers[inst.RD] = uint32(int32(c.Registers[inst.RS1]) / int32(c.Registers[inst.RS2]))
		c.PC += 4
	case "divu":
		c.Registers[inst.RD] = c.Registers[inst.RS1] / c.Registers[inst.RS2]
		c.PC += 4
	case "rem":
		c.Registers[inst.RD] = uint32(int32(c.Registers[inst.RS1]) % int32(c.Registers[inst.RS2]))
		c.PC += 4
	case "remu":
		c.Registers[inst.RD] = c.Registers[inst.RS1] % c.Registers[inst.RS2]
		c.PC += 4
	}
}

func executeI(inst II, c *Cpu) {
	switch inst.Operation() {
	case "addi":
		c.Registers[inst.RD] = c.Registers[inst.RS1] + uint32(int32(uint32(inst.IIM)<<20)>>20)
		c.PC += 4

	case "xori":
		c.Registers[inst.RD] = c.Registers[inst.RS1] ^ uint32(int16(inst.IIM<<4)>>4)
		c.PC += 4

	case "ori":
		c.Registers[inst.RD] = c.Registers[inst.RS1] | uint32(int16(inst.IIM<<4)>>4)
		c.PC += 4

	case "andi":
		c.Registers[inst.RD] = c.Registers[inst.RS1] & uint32(int32(uint32(inst.IIM)<<20)>>20)
		c.PC += 4

	// Shifts should use only last 6 bits
	case "slli":
		c.Registers[inst.RD] = c.Registers[inst.RS1] << (uint32(inst.IIM) & 0x1F)
		c.PC += 4

	case "srli":
		c.Registers[inst.RD] = c.Registers[inst.RS1] >> (uint32(inst.IIM) & 0x1F)
		c.PC += 4

	// Arithmetic Shift, Golang does arithmetic shifts(msb-ext) for signed and logical for unsigned(zero-ext)
	case "srai":
		// We need bottom 5 bits only
		c.Registers[inst.RD] = uint32(int32(c.Registers[inst.RS1]) >> ((inst.IIM << 11) >> 11))
		c.PC += 4

	case "slti":
		// Signed value
		if int32(c.Registers[inst.RS1]) < int32(int16(inst.IIM<<4)>>4) {
			c.Registers[inst.RD] = 1
		} else {
			c.Registers[inst.RD] = 0
		}
		c.PC += 4

	case "sltiu":
		if c.Registers[inst.RS1] < uint32(int16(inst.IIM<<4)>>4) {
			c.Registers[inst.RD] = 1
		} else {
			c.Registers[inst.RD] = 0
		}
		c.PC += 4

	// All Load ones are signed offsets
	case "lb":
		rdi := int32(c.Registers[inst.RS1]) + int32(int16(inst.IIM<<4)>>4)
		c.Registers[inst.RD] = uint32(int8(c.Memory.ReadByte(uint32(rdi))))
		c.PC += 4

	// All Load ones are signed offsets
	case "lh":
		rdi := int32(c.Registers[inst.RS1]) + int32(int16(inst.IIM<<4)>>4)
		c.Registers[inst.RD] = uint32(int16(uint16(c.Memory.ReadByte(uint32(rdi))) | (uint16(c.Memory.ReadByte(uint32(rdi+1))) << 8)))
		c.PC += 4

	// All Load ones are signed offsets
	case "lw":
		rdi := int32(c.Registers[inst.RS1]) + int32(int16(inst.IIM<<4)>>4)
		c.Registers[inst.RD] = uint32(int32(uint32(c.Memory.ReadByte(uint32(rdi))) |
			(uint32(c.Memory.ReadByte(uint32(rdi+1))) << 8) |
			(uint32(c.Memory.ReadByte(uint32(rdi+2))) << 16) |
			(uint32(c.Memory.ReadByte(uint32(rdi+3))) << 24)))
		c.PC += 4

	// All Load ones are signed offsets
	case "lbu":
		rdi := c.Registers[inst.RS1] + uint32(int16(inst.IIM<<4)>>4)
		c.Registers[inst.RD] = uint32(c.Memory.ReadByte(rdi))
		c.PC += 4

	// All Load ones are signed offsets
	case "lhu":
		rdi := c.Registers[inst.RS1] + uint32(int16(inst.IIM<<4)>>4)
		k1 := uint32(c.Memory.ReadByte(rdi))
		k2 := (uint32(c.Memory.ReadByte(rdi+1)) << 8)
		k := k1 | k2
		c.Registers[inst.RD] = k
		c.PC += 4

	case "jalr":
		// This is required because RS1 and RD can be same register
		oldV := c.Registers[inst.RS1]
		c.Registers[inst.RD] = c.PC + 4
		c.PC = oldV + uint32(int16(inst.IIM<<4)>>4)

	case "ecall":
		if os.Getenv("MODE") == "test" || true {
			if c.Registers[10] == 42 || c.Registers[10] == 0 {
				fmt.Fprintln(os.Stdout, fmt.Sprintf("Test Succeeded"))
			} else {
				fmt.Fprintln(os.Stdout, fmt.Sprintf("Ecall: testId: %d, Failed", c.Registers[3]))
			}
			os.Exit(0)
		}
		// If not in test mode Switch context to OS
		c.PC += 4

	case "ebreak":
		// Switch access to Debugger
		c.PC += 4

	case "csrrw":
		// Ignore reading values / registers twice
		xs := c.Registers[inst.RS1]
		c.Registers[inst.RD] = c.CSR.GetValue(uint32(inst.IIM), c.CurrentMode, c)
		c.CSR.SetValue(uint32(inst.IIM), xs, c.CurrentMode, c)
		c.PC += 4

	// For all i or immediate instructions for csr RD is a 5 bit field
	case "csrrwi":
		c.Registers[inst.RD] = c.CSR.GetValue(uint32(inst.IIM), c.CurrentMode, c)
		c.CSR.SetValue(uint32(inst.IIM), uint32(inst.RS1), c.CurrentMode, c)
		c.PC += 4

	case "csrrs":
		// We need more checks here to see if we can indeed modify the registers based on privilege level
		// at which processor is working
		kk := c.CSR.GetValue(uint32(inst.IIM), c.CurrentMode, c)
		csrExisting := c.CSR.GetValue(uint32(inst.IIM), c.CurrentMode, c)
		csrBitmask := c.Registers[inst.RS1]
		c.CSR.SetValue(uint32(inst.IIM), csrBitmask|csrExisting, c.CurrentMode, c)
		c.Registers[inst.RD] = kk
		c.PC += 4

	case "csrrsi":
		// We need more checks here to see if we can indeed modify the registers based on privilege level
		// at which processor is working
		kk := c.CSR.GetValue(uint32(inst.IIM), c.CurrentMode, c)
		// RS1 has the immediate values
		csrExisting := c.CSR.GetValue(uint32(inst.IIM), c.CurrentMode, c)
		csrBitmask := uint32(inst.RS1)
		c.CSR.SetValue(uint32(inst.IIM), csrBitmask|csrExisting, c.CurrentMode, c)
		c.Registers[inst.RD] = kk
		c.PC += 4

	case "csrrc":
		// We need more checks here to see if we can indeed modify the registers based on privilege level
		// at which processor is working
		kk := c.CSR.GetValue(uint32(inst.IIM), c.CurrentMode, c)
		csrExisting := c.CSR.GetValue(uint32(inst.IIM), c.CurrentMode, c)
		csrBitmask := c.Registers[inst.RS1]
		c.CSR.SetValue(uint32(inst.IIM), csrExisting & ^csrBitmask, c.CurrentMode, c)
		c.Registers[inst.RD] = kk
		c.PC += 4

	case "csrrci":
		// We need more checks here to see if we can indeed modify the registers based on privilege level
		// at which processor is working
		kk := c.CSR.GetValue(uint32(inst.IIM), c.CurrentMode, c)
		csrExisting := c.CSR.GetValue(uint32(inst.IIM), c.CurrentMode, c)
		csrBitmask := uint32(inst.RS1)
		c.CSR.SetValue(uint32(inst.IIM), csrExisting & ^csrBitmask, c.CurrentMode, c)
		c.Registers[inst.RD] = kk
		c.PC += 4
	case "sret":
		statusReg := ToMStatusReg(c.CSR.GetValue(SSTATUS, c.CurrentMode, c))
		statusReg.sie = statusReg.spie
		// restore mode of processor
		c.CurrentMode = statusReg.spp
		statusReg.spie = 1
		// Set to U-Mode
		statusReg.spp = 0
		c.CSR.Registers[SSTATUS] = FromMStatusReg(statusReg)
		// change return PC
		c.PC = c.CSR.Registers[SEPC]
	case "mret":
		statusReg := ToMStatusReg(c.CSR.GetValue(MSTATUS, c.CurrentMode, c))
		statusReg.mie = statusReg.mpie
		// restore mode of processor
		c.CurrentMode = statusReg.mpp
		statusReg.mpie = 1
		if statusReg.mpp != 3 {
			statusReg.mprv = 0
		}
		// Set to U-Mode
		statusReg.mpp = 0
		c.CSR.Registers[MSTATUS] = FromMStatusReg(statusReg)
		// change return PC
		c.PC = c.CSR.Registers[MEPC]
	case "sfence.vma":
		// nop operation
		c.PC += 4
	case "wfi":
		// Wait until interrupt comes
		mip := c.CSR.Registers[MIP]
		if mip > 0 {
			c.PC += 4
		}
		// If nothing then simply wait
	default:
		panic("running instruction failed")
	}
}

func executeS(inst SI, c *Cpu) {
	switch inst.Operation() {
	// All Store ones are signed offsets
	case "sb":
		c.Memory.WriteByte(byte(c.Registers[int(inst.RS2)]&uint32(0xFF)), c.Registers[int(inst.RS1)]+uint32(int16(inst.SIM<<4)>>4))
		c.PC += 4

	// All Store ones are signed offsets
	case "sh":
		c.Memory.WriteByte(byte(c.Registers[int(inst.RS2)]&uint32(0xFF)), c.Registers[int(inst.RS1)]+uint32(int16(inst.SIM<<4)>>4))
		c.Memory.WriteByte(byte((c.Registers[int(inst.RS2)]&uint32(0xFF00))>>8), c.Registers[int(inst.RS1)]+uint32(int16(inst.SIM<<4)>>4)+1)
		c.PC += 4

	case "sw":
		c.Memory.WriteByte(byte(c.Registers[int(inst.RS2)]&uint32(0xFF)), c.Registers[int(inst.RS1)]+uint32(int16(inst.SIM<<4)>>4))
		c.Memory.WriteByte(byte((c.Registers[int(inst.RS2)]&uint32(0xFF00))>>8), c.Registers[int(inst.RS1)]+uint32(int16(inst.SIM<<4)>>4)+1)
		c.Memory.WriteByte(byte((c.Registers[int(inst.RS2)]&uint32(0xFF0000))>>16), c.Registers[int(inst.RS1)]+uint32(int16(inst.SIM<<4)>>4)+2)
		c.Memory.WriteByte(byte((c.Registers[int(inst.RS2)]&uint32(0xFF000000))>>24), c.Registers[int(inst.RS1)]+uint32(int16(inst.SIM<<4)>>4)+3)
		c.PC += 4
	}
}

func executeB(inst BI, c *Cpu) {
	switch inst.Operation() {
	case "beq":
		if c.Registers[inst.RS1] == c.Registers[inst.RS2] {
			c.PC += uint32(int16(inst.BIM<<4) >> 4)
		} else {
			c.PC += 4
		}
	case "bne":
		if c.Registers[inst.RS1] != c.Registers[inst.RS2] {
			c.PC += uint32(int16(inst.BIM<<4) >> 4)
		} else {
			c.PC += 4
		}
	case "blt":
		if int32(c.Registers[inst.RS1]) < int32(c.Registers[inst.RS2]) {
			c.PC = c.PC + uint32(int16(inst.BIM<<4)>>4)
		} else {
			c.PC += 4
		}
	case "bge":
		if int32(c.Registers[inst.RS1]) >= int32(c.Registers[inst.RS2]) {
			c.PC = c.PC + uint32(int16(inst.BIM<<4)>>4)
		} else {
			c.PC += 4
		}
	case "bltu":
		if c.Registers[inst.RS1] < c.Registers[inst.RS2] {
			c.PC = c.PC + uint32(int16(inst.BIM<<4)>>4)
		} else {
			c.PC += 4
		}
	case "bgeu":
		if c.Registers[inst.RS1] >= c.Registers[inst.RS2] {
			c.PC = c.PC + uint32(int16(inst.BIM<<4)>>4)
		} else {
			c.PC += 4
		}
	}
}

func executeJ(inst JI, c *Cpu) {
	switch inst.Operation() {
	case "jal":
		c.Registers[inst.RD] = c.PC + 4
		c.PC += uint32(int32(inst.JIM<<12) >> 12)
	}
}

// Immediate value not sign extended, all others are sign extended
func executeU(inst UI, c *Cpu) {
	switch inst.Operation() {
	case "lui":
		c.Registers[inst.RD] = inst.UIM1 << 12
		c.PC += 4
	case "auipc":
		c.Registers[inst.RD] = c.PC + uint32(int32(inst.UIM1<<12))
		c.PC += 4
	}
}
