package instructions

func decodeRS1(inst uint32) byte {
	return byte(getBitsAsUInt32(inst, 15, 19))
}

func decodeRS2(inst uint32) byte {
	return byte(getBitsAsUInt32(inst, 20, 24))
}

func decodeF3(inst uint32) byte {
	return byte(getBitsAsUInt32(inst, 12, 14))
}

func decodeF7(inst uint32) byte {
	return byte(getBitsAsUInt32(inst, 25, 31))
}

func decodeOpcode(inst uint32) byte {
	return byte(getBitsAsUInt32(inst, 0, 6))
}

func decodeRD(inst uint32) byte {
	return byte(getBitsAsUInt32(inst, 7, 11))
}
