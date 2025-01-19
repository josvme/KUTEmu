package instructions

type JI struct {
	Opcode byte
	RD     byte
	JIM    uint32
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
		JIM:    jim1,
	}
}

func decodeJIM1(inst uint32) uint32 {
	b1219 := getBitsAsUInt32(inst, 12, 19)
	b11 := getBitsAsUInt32(inst, 20, 20)
	b110 := getBitsAsUInt32(inst, 21, 30)
	b20 := getBitsAsUInt32(inst, 31, 31)
	a := (b1219 << 12)
	b := (b11 << 11)
	c := (b110 << 1)
	d := (b20 << 20)
	e := a | b | c | d
	return e
}
