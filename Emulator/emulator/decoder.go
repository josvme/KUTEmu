package emulator

import "fmt"

type Inst interface {
	decode(inst uint32) Inst
}

func transformLittleToBig(inst [4]byte) uint32 {
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

func decodeBytes(by [4]byte) Inst {
	c := transformLittleToBig(by)
	op := decodeOpcode(c)
	switch op {
	// R
	case 0b0110011:
		return RI{}.decode(c)
	// I
	case 0b0010011:
		return II{}.decode(c)
	// I
	case 0b0000011:
		return II{}.decode(c)
	// S
	case 0b0100011:
		return SI{}.decode(c)
	// B
	case 0b1100011:
		return BI{}.decode(c)
	// J
	case 0b1101111:
		return JI{}.decode(c)
	// I
	case 0b1100111:
		return II{}.decode(c)
	// U
	case 0b0110111:
		return UI{}.decode(c)
	// U
	case 0b0010111:
		return UI{}.decode(c)
	// I
	case 0b1110011:
		return II{}.decode(c)
	// Fence
	case 0b0001111:
		return FI{}.decode(c)

	default:
		panic(fmt.Sprintf("Unknown instruction %v", c))
	}
}
