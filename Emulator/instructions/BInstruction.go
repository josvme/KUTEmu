package instructions

type BI struct {
	Opcode byte
	BIM1   byte
	F3     byte
	RS1    byte
	RS2    byte
	BIM2   byte
}

const OP_TOPLEVEL_BI = 0b1100011

func (i BI) Operation() string {
	if i.Opcode != OP_TOPLEVEL_BI {
		panic("This shouldn't happen")
	}
	switch {
	case i.F3 == 0x0:
		return "beq"
	case i.F3 == 0x1:
		return "bne"
	case i.F3 == 0x4:
		return "blt"
	case i.F3 == 0x5:
		return "bge"
	case i.F3 == 0x6:
		return "bltu"
	case i.F3 == 0x7:
		return "bgeu"
	default:
		panic("Unknown Operation")
	}
}

func (i BI) Decode(inst uint32) Inst {
	op := decodeOpcode(inst)
	bim1 := decodeBIM1(inst)
	f3 := decodeF3(inst)
	rs1 := decodeRS1(inst)
	rs2 := decodeRS2(inst)
	bims2 := decodeBIM2(inst)
	return BI{
		Opcode: op,
		BIM1:   bim1,
		F3:     f3,
		RS1:    rs1,
		RS2:    rs2,
		BIM2:   bims2,
	}
}

func decodeBIM1(inst uint32) byte {
	return byte(getBitsAsUInt32(inst, 7, 11))
}

func decodeBIM2(inst uint32) byte {
	return byte(getBitsAsUInt32(inst, 25, 31))
}
