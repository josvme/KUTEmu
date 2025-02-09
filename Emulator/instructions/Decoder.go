package instructions

import "fmt"

type Inst interface {
	Decode(inst uint32) Inst
	Operation() string
}

func TransformLittleToBig(inst [4]byte) uint32 {
	return uint32(inst[3])<<24 | uint32(inst[2])<<16 | uint32(inst[1])<<8 | uint32(inst[0])
}

func getBitsAsUInt32(inst uint32, start int, end int) uint32 {
	big := uint32(0)
	small := uint32(0)
	for i := range end + 1 {
		big = big | (1 << i)
	}
	for i := range start {
		small = small | (1 << i)
	}
	mask := big ^ small
	val := mask & inst
	val = val >> start
	return val
}

func DecodeBytes(by [4]byte) Inst {
	c := TransformLittleToBig(by)
	op := decodeOpcode(c)
	switch op {
	// R
	case 0b0110011:
		return RI{}.Decode(c)
	// R
	case 0b0101111:
		return RI{}.Decode(c)
	// I
	// I
	case 0b0010011:
		return II{}.Decode(c)
	// I
	case 0b0000011:
		return II{}.Decode(c)
	// S
	case 0b0100011:
		return SI{}.Decode(c)
	// B
	case 0b1100011:
		return BI{}.Decode(c)
	// J
	case 0b1101111:
		return JI{}.Decode(c)
	// I
	case 0b1100111:
		return II{}.Decode(c)
	// U
	case 0b0110111:
		return UI{}.Decode(c)
	// U
	case 0b0010111:
		return UI{}.Decode(c)
	// I
	case 0b1110011:
		return II{}.Decode(c)
	// Fence
	case 0b0001111:
		return FI{}.Decode(c)

	default:
		panic(fmt.Sprintf("Unknown instruction 0x%X", c))
	}
}
