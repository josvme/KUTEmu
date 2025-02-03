package instructions

type II struct {
	Opcode byte
	RD     byte
	F3     byte
	RS1    byte
	IIM    uint16
}

const OP_TOPLEVEL_ARITH = 0b0010011
const OP_TOPLEVEL_LOAD = 0b0000011
const OP_TOPLEVEL_JUMP = 0b1101111
const OP_TOPLEVEL_JUMP_2 = 0b1100111
const OP_TOPLEVEL_ENVIRON = 0b1110011

func (i II) Operation() string {
	switch {
	case i.F3 == 0x0 && i.Opcode == OP_TOPLEVEL_ARITH:
		return "addi"
	case i.F3 == 0x4 && i.Opcode == OP_TOPLEVEL_ARITH:
		return "xori"
	case i.F3 == 0x6 && i.Opcode == OP_TOPLEVEL_ARITH:
		return "ori"
	case i.F3 == 0x7 && i.Opcode == OP_TOPLEVEL_ARITH:
		return "andi"
	case i.F3 == 0x1 && i.Opcode == OP_TOPLEVEL_ARITH && immValue(i.IIM) == 0x00:
		return "slli"
	case i.F3 == 0x5 && i.Opcode == OP_TOPLEVEL_ARITH && immValue(i.IIM) == 0x00:
		return "srli"
	case i.F3 == 0x5 && i.Opcode == OP_TOPLEVEL_ARITH && immValue(i.IIM) == 0x20:
		return "srai"
	case i.F3 == 0x2 && i.Opcode == OP_TOPLEVEL_ARITH:
		return "slti"
	case i.F3 == 0x3 && i.Opcode == OP_TOPLEVEL_ARITH:
		return "sltiu"
	case i.F3 == 0x0 && i.Opcode == OP_TOPLEVEL_LOAD:
		return "lb"
	case i.F3 == 0x1 && i.Opcode == OP_TOPLEVEL_LOAD:
		return "lh"
	case i.F3 == 0x2 && i.Opcode == OP_TOPLEVEL_LOAD:
		return "lw"
	case i.F3 == 0x4 && i.Opcode == OP_TOPLEVEL_LOAD:
		return "lbu"
	case i.F3 == 0x5 && i.Opcode == OP_TOPLEVEL_LOAD:
		return "lhu"
	case i.F3 == 0x0 && i.Opcode == OP_TOPLEVEL_JUMP_2:
		return "jalr"
	// We also Zicsr/Env instructions here.
	case i.F3 == 0x0 && i.Opcode == OP_TOPLEVEL_ENVIRON && i.IIM == 0x00:
		return "ecall"
	case i.F3 == 0x0 && i.Opcode == OP_TOPLEVEL_ENVIRON && i.IIM == 0x01:
		return "ebreak"
	case i.F3 == 0x1 && i.Opcode == OP_TOPLEVEL_ENVIRON:
		return "csrrw"
	case i.F3 == 0x2 && i.Opcode == OP_TOPLEVEL_ENVIRON:
		return "csrrs"

	case i.F3 == 0x3 && i.Opcode == OP_TOPLEVEL_ENVIRON:
		return "csrrc"

	case i.F3 == 0x5 && i.Opcode == OP_TOPLEVEL_ENVIRON:
		return "csrrwi"

	case i.F3 == 0x6 && i.Opcode == OP_TOPLEVEL_ENVIRON:
		return "csrrsi"

	case i.F3 == 0x7 && i.Opcode == OP_TOPLEVEL_ENVIRON:
		return "csrrci"
	default:
		panic("Unknown Operation")
	}
}

func immValue(im uint16) byte {
	start := 5
	end := 11
	big := uint16(0)
	small := uint16(0)
	for i := range end + 1 {
		big = big | (1 << i)
	}
	for i := range start {
		small = small | (1 << i)
	}
	mask := big ^ small
	val := mask & im
	val = val >> start
	return byte(val)
}

func (i II) Decode(inst uint32) Inst {
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
		IIM:    iim1,
	}
}

func decodeIIM1(inst uint32) uint16 {
	return uint16(getBitsAsUInt32(inst, 20, 31))
}
