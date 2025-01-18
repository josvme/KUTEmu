package emulator

type RI struct {
	Opcode byte
	RD     byte
	F3     byte
	RS1    byte
	RS2    byte
	F7     byte
}

func (i RI) decode(inst uint32) Inst {
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
