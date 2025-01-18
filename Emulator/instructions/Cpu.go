package instructions

type Cpu struct {
	PC        uint32
	Registers [32]uint32
	Memory    Memory
}

func (c *Cpu) ExecInst(i Inst) error {
	switch i.(type) {
	case II:
		ins := i.(II)
		executeI(ins, c)
	case UI:
		ins := i.(UI)
		executeU(ins, c)
	case SI:
		ins := i.(SI)
		executeS(ins, c)
	case JI:
		ins := i.(JI)
		executeJ(ins, c)
	}
	return nil
}

func executeJ(inst JI, c *Cpu) {
	switch inst.Operation() {
	case "jal":
		c.Registers[inst.RD] = c.PC + 4
		c.PC += inst.JIM1
	}
}

func executeS(inst SI, c *Cpu) {
	switch inst.Operation() {
	case "sb":
		c.Memory.Write(byte(c.Registers[inst.RS2]&uint32(0xFF)), c.Registers[inst.RS1]+uint32(inst.SIM2<<5)+uint32(inst.SIM1))
		c.PC += 4
	}
}

func executeU(inst UI, c *Cpu) {
	switch inst.Operation() {
	case "lui":
		c.Registers[inst.RD] = inst.UIM1 << 12
		c.PC += 4
	}
}

func executeI(inst II, c *Cpu) {
	switch inst.Operation() {
	case "addi":
		c.Registers[inst.RD] = c.Registers[inst.RS1] + uint32(inst.IIM1)
		c.PC += 4
	}
}
