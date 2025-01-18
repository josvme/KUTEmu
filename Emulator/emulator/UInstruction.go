package emulator

type UI struct {
	Opcode byte
	RD     byte
	UIM1   uint32
}

func (i UI) decode(inst uint32) Inst {
	op := decodeOpcode(inst)
	rd := decodeRD(inst)
	uim1 := decodeUIM1(inst)
	return UI{
		Opcode: op,
		RD:     rd,
		UIM1:   uim1,
	}
}

func decodeUIM1(inst uint32) uint32 {
	return getBitsAsUInt32(inst, 12, 31)
}
