package instructions

type SI struct {
	Opcode byte
	F3     byte
	RS1    byte
	RS2    byte
	SIM    uint16
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
	f3 := decodeF3(inst)
	rs1 := decodeRS1(inst)
	rs2 := decodeRS2(inst)
	sim := decodeSIM(inst)
	return SI{
		Opcode: op,
		F3:     f3,
		RS1:    rs1,
		RS2:    rs2,
		SIM:    sim,
	}
}

func decodeSIM(inst uint32) uint16 {
	lower := uint16(getBitsAsUInt32(inst, 7, 11))
	upper := uint16(getBitsAsUInt32(inst, 25, 31))

	return upper<<5 | lower
}
