package instructions

const VIRT_UART0 = 0x10000000

type Memory struct {
	Map  map[uint32]byte
	Uart UART
}

func (m *Memory) LoadBytes(b []byte, location uint32) error {
	for i, bb := range b {
		m.Map[location+uint32(i)] = bb
	}
	return nil
}

func (m *Memory) WriteByte(b byte, location uint32) {
	if location >= VIRT_UART0 {
		_ = m.Uart.Write(b)
		return
	}
	m.Map[location] = b
}

func (m *Memory) WriteHalf(h uint16, location uint32) {
	m.WriteByte(byte(h&uint16(0xFF)), location)
	m.WriteByte(byte((h&uint16(0xFF00))>>8), location+1)
}

func (m *Memory) WriteWord(w uint32, location uint32) {
	m.WriteByte(byte(w&uint32(0xFF)), location)
	m.WriteByte(byte((w&uint32(0xFF00))>>8), location+1)
	m.WriteByte(byte((w&uint32(0xFF0000))>>16), location+2)
	m.WriteByte(byte((w&uint32(0xFF000000))>>24), location+3)
}

func (m *Memory) ReadByte(location uint32) byte {
	return m.Map[location]
}

func (m *Memory) ReadHalf(location uint32) uint16 {
	return uint16(m.ReadByte(location)) | (uint16(m.ReadByte(location+1)) << 8)
}

func (m *Memory) ReadWord(location uint32) uint32 {
	return uint32(m.ReadByte(location)) |
		(uint32(m.ReadByte(location+1)) << 8) |
		(uint32(m.ReadByte(location+2)) << 16) |
		(uint32(m.ReadByte(location+3)) << 24)
}
