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
const OP_TOPLEVEL_ATOMIC_RI = 0b0101111

func (i RI) Operation() string {
	if i.Opcode != OP_TOPLEVEL_RI && i.Opcode != OP_TOPLEVEL_ATOMIC_RI {
		panic("This shouldn't happen")
	}

	//Rl := i.F7 & 0b00000001
	//Aq := i.F7 & 0b00000010
	F5 := (i.F7 & 0b01111100 << 1) >> 3
	if i.Opcode == OP_TOPLEVEL_RI {
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
		// Add atomic instructions
		case i.F3 == 0x2 && F5 == 0x02:
			return "lr.w"
		case i.F3 == 0x2 && F5 == 0x03:
			return "sc.w"

		default:
			panic("Unknown R Operation")
		}
	}

	if i.Opcode == OP_TOPLEVEL_ATOMIC_RI {
		switch {
		case i.F3 == 0x2 && F5 == 0x01:
			return "amoswap.w"
		case i.F3 == 0x2 && F5 == 0x00:
			return "amoadd.w"
		case i.F3 == 0x2 && F5 == 0x0C:
			return "amoand.w"
		case i.F3 == 0x2 && F5 == 0x08:
			return "amoor.w"
		case i.F3 == 0x2 && F5 == 0x04:
			return "amoxor.w"
		case i.F3 == 0x2 && F5 == 0x14:
			return "amomax.w"
		case i.F3 == 0x2 && F5 == 0x10:
			return "amomin.w"
		case i.F3 == 0x2 && F5 == 0x1c:
			return "amomaxu.w"
		case i.F3 == 0x2 && F5 == 0x18:
			return "amominu.w"
		// Add multiply instructions
		case i.F3 == 0x0 && i.F7 == 0x01:
			return "mul"
		case i.F3 == 0x1 && i.F7 == 0x01:
			return "mulh"
		case i.F3 == 0x2 && i.F7 == 0x01:
			return "mulhsu"
		case i.F3 == 0x3 && i.F7 == 0x01:
			return "mulhu"
		case i.F3 == 0x4 && i.F7 == 0x01:
			return "div"
		case i.F3 == 0x5 && i.F7 == 0x01:
			return "divu"
		case i.F3 == 0x6 && i.F7 == 0x01:
			return "rem"
		case i.F3 == 0x7 && i.F7 == 0x01:
			return "remu"
		default:
			panic("Unknown A Operation")
		}
	}
	panic("Unknown Operation")
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
