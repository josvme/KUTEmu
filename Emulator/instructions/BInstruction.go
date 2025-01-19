package instructions

type BI struct {
	Opcode byte
	F3     byte
	RS1    byte
	RS2    byte
	BIM    uint16
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
	bim := decodeBIM(inst)
	f3 := decodeF3(inst)
	rs1 := decodeRS1(inst)
	rs2 := decodeRS2(inst)
	return BI{
		Opcode: op,
		F3:     f3,
		RS1:    rs1,
		RS2:    rs2,
		BIM:    bim,
	}
}

func decodeBIM(inst uint32) uint16 {
	b14 := uint16(getBitsAsUInt32(inst, 8, 11))
	b11 := uint16(getBitsAsUInt32(inst, 7, 7))
	b105 := uint16(getBitsAsUInt32(inst, 25, 30))
	b12 := uint16(getBitsAsUInt32(inst, 31, 31))
	a := b14 << 1
	b := b11 << 11
	c := b105 << 5
	d := b12 << 12
	e := a | b | c | d
	return e
}
