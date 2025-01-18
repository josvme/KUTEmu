package emulator

type II struct {
	Opcode byte
	RD     byte
	F3     byte
	RS1    byte
	IIM1   uint16
}

func (i II) decode(inst uint32) Inst {
	op := decodeOpcode(inst)
	rd := decodeRD(inst)
	f3 := decodeF3(inst)
	rs1 := decodeRS1(inst)
	iim1 := decodeIIM1(inst)

	return II{
		Opcode: op,
		RD:     rd,
		F3:     f3,
		RS1:    rs1,
		IIM1:   iim1,
	}
}

func decodeIIM1(inst uint32) uint16 {
	return uint16(getBitsAsUInt32(inst, 20, 31))
}
