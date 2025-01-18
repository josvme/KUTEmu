package emulator

type JI struct {
	Opcode byte
	RD     byte
	JIM1   uint32
}

func (i JI) decode(inst uint32) Inst {
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
