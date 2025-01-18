package emulator

type FI struct {
	Opcode byte
	RD     byte
	F3     byte
	RS1    byte
	Succ   byte
	Pred   byte
	FM     byte
}

func (i FI) decode(inst uint32) Inst {
	op := decodeOpcode(inst)
	rd := decodeRD(inst)
	f3 := decodeF3(inst)
	rs1 := decodeRS1(inst)
	succ := decodeSucc(inst)
	pred := decodePred(inst)
	fm := decodeFM(inst)

	return FI{
		Opcode: op,
		RD:     rd,
		F3:     f3,
		RS1:    rs1,
		Succ:   succ,
		Pred:   pred,
		FM:     fm,
	}
}

func decodeSucc(inst uint32) byte {
	return byte(getBitsAsUInt32(inst, 20, 23))
}

func decodePred(inst uint32) byte {
	return byte(getBitsAsUInt32(inst, 24, 27))
}

func decodeFM(inst uint32) byte {
	return byte(getBitsAsUInt32(inst, 28, 31))
}
