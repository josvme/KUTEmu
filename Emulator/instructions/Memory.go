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

func (m *Memory) Write(b byte, location uint32) {
	if location >= VIRT_UART0 {
		_ = m.Uart.Write(b)
		return
	}
	m.Map[location] = b
}

func (m *Memory) Read(location uint32) byte {
	return m.Map[location]
}
