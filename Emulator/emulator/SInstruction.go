package emulator

type SI struct {
	Opcode byte
	SIM1   byte
	F3     byte
	RS1    byte
	RS2    byte
	SIM2   byte
}

func (i SI) decode(inst uint32) Inst {
	op := decodeOpcode(inst)
	sim1 := decodeSIM1(inst)
	f3 := decodeF3(inst)
	rs1 := decodeRS1(inst)
	rs2 := decodeRS2(inst)
	sim2 := decodeSIM2(inst)
	return SI{
		Opcode: op,
		SIM1:   sim1,
		F3:     f3,
		RS1:    rs1,
		RS2:    rs2,
		SIM2:   sim2,
	}
}

func decodeSIM2(inst uint32) byte {
	return byte(getBitsAsUInt32(inst, 7, 11))
}

func decodeSIM1(inst uint32) byte {
	return byte(getBitsAsUInt32(inst, 25, 31))
}
