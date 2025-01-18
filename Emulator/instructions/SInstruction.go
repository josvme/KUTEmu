package instructions

type SI struct {
	Opcode byte
	SIM1   byte
	F3     byte
	RS1    byte
	RS2    byte
	SIM2   byte
}

const OP_TOPLEVEL_SI = 0b0100011

func (i SI) Operation() string {
	if i.Opcode != OP_TOPLEVEL_SI {
		panic("This shouldn't happen")
	}
	switch {
	case i.F3 == 0x00:
		return "sb"
	case i.F3 == 0x01:
		return "sh"
	case i.F3 == 0x02:
		return "sw"
	default:
		panic("Unknown Operation")
	}
}

func (i SI) Decode(inst uint32) Inst {
	op := decodeOpcode(inst)
	sim1 := decodeSIM1(inst)
	f3 := decodeF3(inst)
	rs1 := decodeRS1(inst)
	rs2 := decodeRS2(inst)
	sim2 := decodeSIM2(inst)
	return SI{
		Opcode: op,
		SIM1:   sim1,
		F3:     f3,
		RS1:    rs1,
		RS2:    rs2,
		SIM2:   sim2,
	}
}

func decodeSIM2(inst uint32) byte {
	return byte(getBitsAsUInt32(inst, 7, 11))
}

func decodeSIM1(inst uint32) byte {
	return byte(getBitsAsUInt32(inst, 25, 31))
}
