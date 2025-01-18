package instructions

type JI struct {
	Opcode byte
	RD     byte
	JIM1   uint32
}

const OP_TOPLEVEL_JI = 0b1101111

func (i JI) Operation() string {
	if i.Opcode != OP_TOPLEVEL_JI {
		panic("This shouldn't happen")
	}
	return "jal"
}

func (i JI) Decode(inst uint32) Inst {
	op := decodeOpcode(inst)
	rd := decodeRD(inst)
	jim1 := decodeJIM1(inst)
	return JI{
		Opcode: op,
		RD:     rd,
		JIM1:   jim1,
	}
}

func decodeJIM1(inst uint32) uint32 {
	return getBitsAsUInt32(inst, 12, 31)
}
