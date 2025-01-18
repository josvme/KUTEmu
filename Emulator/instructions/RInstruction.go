package instructions

type RI struct {
	Opcode byte
	RD     byte
	F3     byte
	RS1    byte
	RS2    byte
	F7     byte
}

const OP_TOPLEVEL_RI = 0b0110011

func (i RI) Operation() string {
	if i.Opcode != OP_TOPLEVEL_RI {
		panic("This shouldn't happen")
	}
	switch {
	case i.F3 == 0x0 && i.F7 == 0x00:
		return "add"
	case i.F3 == 0x0 && i.F7 == 0x20:
		return "sub"
	case i.F3 == 0x4 && i.F7 == 0x00:
		return "xor"
	case i.F3 == 0x6 && i.F7 == 0x00:
		return "or"
	case i.F3 == 0x7 && i.F7 == 0x00:
		return "and"
	case i.F3 == 0x1 && i.F7 == 0x00:
		return "sll"
	case i.F3 == 0x5 && i.F7 == 0x00:
		return "srl"
	case i.F3 == 0x5 && i.F7 == 0x20:
		return "sra"
	case i.F3 == 0x2 && i.F7 == 0x00:
		return "slt"
	case i.F3 == 0x3 && i.F7 == 0x00:
		return "sltu"
	default:
		panic("Unknown Operation")
	}
}

func (i RI) Decode(inst uint32) Inst {
	op := decodeOpcode(inst)
	rd := decodeRD(inst)
	f3 := decodeF3(inst)
	rs1 := decodeRS1(inst)
	rs2 := decodeRS2(inst)
	f7 := decodeF7(inst)
	return RI{
		Opcode: op,
		RD:     rd,
		F3:     f3,
		RS1:    rs1,
		RS2:    rs2,
		F7:     f7,
	}
}
