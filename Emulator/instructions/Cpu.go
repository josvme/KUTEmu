package instructions

type Cpu struct {
	PC        uint32
	Registers [32]uint32
	Memory    Memory
}

func (c *Cpu) ExecInst(i Inst) error {
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

func executeF(inst FI, c *Cpu) {
	// Order memory access
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
		c.Registers[inst.RD] = c.Registers[inst.RS1] << c.Registers[inst.RS2]
		c.PC += 4

	case "srl":
		c.Registers[inst.RD] = c.Registers[inst.RS1] >> c.Registers[inst.RS2]
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
	}
}

func executeI(inst II, c *Cpu) {
	switch inst.Operation() {
	case "addi":
		c.Registers[inst.RD] = c.Registers[inst.RS1] + uint32(inst.IIM)
		c.PC += 4

	case "xori":
		c.Registers[inst.RD] = c.Registers[inst.RS1] ^ uint32(inst.IIM)
		c.PC += 4

	case "ori":
		c.Registers[inst.RD] = c.Registers[inst.RS1] | uint32(inst.IIM)
		c.PC += 4

	case "andi":
		c.Registers[inst.RD] = c.Registers[inst.RS1] & uint32(inst.IIM)
		c.PC += 4

	case "slli":
		c.Registers[inst.RD] = c.Registers[inst.RS1] << uint32(inst.IIM) & 0x1F
		c.PC += 4

	case "srli":
		c.Registers[inst.RD] = c.Registers[inst.RS1] >> uint32(inst.IIM) & 0x1F
		c.PC += 4

	// Arithmetic Shift, Golang does arithmetic shifts(msb-ext) for signed and logical for unsigned(zero-ext)
	case "srai":
		c.Registers[inst.RD] = uint32(int32(c.Registers[inst.RS1]) >> inst.IIM)
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
		if c.Registers[inst.RS1] < uint32(inst.IIM) {
			c.Registers[inst.RD] = 1
		} else {
			c.Registers[inst.RD] = 0
		}
		c.PC += 4

	// All Load ones are signed offsets
	case "lb":
		rdi := int32(c.Registers[inst.RS1]) + int32(int16(inst.IIM<<3)>>3)
		existingRD := c.Registers[inst.RD] & 0xFFFFFF00
		c.Registers[inst.RD] = existingRD | uint32(c.Memory.Read(uint32(rdi)))
		c.PC += 4

	// All Load ones are signed offsets
	case "lh":
		rdi := int32(c.Registers[inst.RS1]) + int32(int16(inst.IIM<<3)>>3)
		existingRD := c.Registers[inst.RD] & 0xFFFF0000
		c.Registers[inst.RD] = existingRD | uint32(c.Memory.Read(uint32(rdi))) | (uint32(c.Memory.Read(uint32(rdi+1))) << 8)
		c.PC += 4

	// All Load ones are signed offsets
	case "lw":
		rdi := int32(c.Registers[inst.RS1]) + int32(int16(inst.IIM<<3)>>3)
		c.Registers[inst.RD] = uint32(c.Memory.Read(uint32(rdi))) |
			(uint32(c.Memory.Read(uint32(rdi+1))) << 8) |
			(uint32(c.Memory.Read(uint32(rdi+2))) << 16) |
			(uint32(c.Memory.Read(uint32(rdi+3))) << 24)
		c.PC += 4

	// All Load ones are signed offsets
	case "lbu":
		rdi := c.Registers[inst.RS1] + uint32(int16(inst.IIM<<3)>>3)
		existingRD := c.Registers[inst.RD] & 0xFFFFFF00
		c.Registers[inst.RD] = existingRD | uint32(c.Memory.Read(rdi))
		c.PC += 4

	// All Load ones are signed offsets
	case "lhu":
		rdi := c.Registers[inst.RS1] + uint32(int16(inst.IIM<<3)>>3)
		existingRD := c.Registers[inst.RD] & 0xFFFF0000
		c.Registers[inst.RD] = existingRD | uint32(c.Memory.Read(rdi)) | (uint32(c.Memory.Read(rdi+1)) << 8)
		c.PC += 4

	case "jalr":
		c.Registers[inst.RD] = c.PC + 4
		c.PC += uint32(inst.IIM)

	case "ecall":
		// Switch context to OS
		c.PC += 4

	case "ebreak":
		// Switch access to DB
		c.PC += 4
	}
}

func executeS(inst SI, c *Cpu) {
	switch inst.Operation() {
	// All Store ones are signed offsets
	case "sb":
		c.Memory.Write(byte(c.Registers[inst.RS2]&uint32(0xFF)), c.Registers[inst.RS1]+uint32(int16(inst.SIM<<4)>>4))
		c.PC += 4

	// All Store ones are signed offsets
	case "sh":
		c.Memory.Write(byte(c.Registers[inst.RS2]&uint32(0xFF)), c.Registers[inst.RS1]+uint32(int16(inst.SIM<<4)>>4))
		c.Memory.Write(byte((c.Registers[inst.RS2]&uint32(0xFF00))>>8), c.Registers[inst.RS1]+uint32(int16(inst.SIM<<4)>>4)+1)
		c.PC += 4

	case "sw":
		c.Memory.Write(byte(c.Registers[inst.RS2]&uint32(0xFF)), c.Registers[inst.RS1]+uint32(int16(inst.SIM<<4)>>4))
		c.Memory.Write(byte((c.Registers[inst.RS2]&uint32(0xFF00))>>8), c.Registers[inst.RS1]+uint32(int16(inst.SIM<<4)>>4)+1)
		c.Memory.Write(byte((c.Registers[inst.RS2]&uint32(0xFF0000))>>16), c.Registers[inst.RS1]+uint32(int16(inst.SIM<<4)>>4)+2)
		c.Memory.Write(byte((c.Registers[inst.RS2]&uint32(0xFF000000))>>24), c.Registers[inst.RS1]+uint32(int16(inst.SIM<<4)>>4)+3)
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
		if c.Registers[inst.RS1] < c.Registers[inst.RS2] {
			// We ignore Overflows here
			c.PC = uint32(int32(c.PC) + int32(inst.BIM))
		} else {
			c.PC += 4
		}
	case "bge":
		if c.Registers[inst.RS1] >= c.Registers[inst.RS2] {
			// We ignore Overflows here
			c.PC = uint32(int32(c.PC) + int32(inst.BIM))
		} else {
			c.PC += 4
		}
	case "bltu":
		if c.Registers[inst.RS1] < c.Registers[inst.RS2] {
			c.PC = c.PC + uint32(inst.BIM)
		} else {
			c.PC += 4
		}
	case "bgeu":
		if c.Registers[inst.RS1] >= c.Registers[inst.RS2] {
			c.PC = c.PC + uint32(inst.BIM)
		} else {
			c.PC += 4
		}
	}
}

func executeJ(inst JI, c *Cpu) {
	switch inst.Operation() {
	case "jal":
		c.Registers[inst.RD] = c.PC + 4
		c.PC += uint32(int32(inst.JIM<<11) >> 11)
	}
}

// Immediate value not sign extended, all others are sign extended
func executeU(inst UI, c *Cpu) {
	switch inst.Operation() {
	case "lui":
		c.Registers[inst.RD] = inst.UIM1 << 12
		c.PC += 4
	case "auipc":
		c.Registers[inst.RD] = c.PC + inst.UIM1<<12
		c.PC += 4
	}
}
