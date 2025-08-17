package instructions

const VIRT_UART0 = 0x10000000
const VIRT_DISPLAY = 0x1D385000
const VIRT_DISPLAY_SIZE = 320 * 200

type Memory struct {
	Map     map[uint32]byte
	Uart    *UART
	Plic    *Plic
	Cpu     *Cpu
	Clint   *Clint
	Display *Display
}

// Hack
func (m *Memory) SetCpu(cpu *Cpu) {
	m.Cpu = cpu
}

func (m *Memory) LoadBytes(b []byte, location uint32) error {
	for i, bb := range b {
		m.Map[location+uint32(i)] = bb
	}
	return nil
}

func (m *Memory) WriteByte(b byte, location uint32) {
	if location >= VIRT_UART0 && location < VIRT_UART0+0x100 {
		_ = m.Uart.Write(b, location-VIRT_UART0)
		return
	}
	if location >= VIRT_DISPLAY && location < VIRT_DISPLAY+VIRT_DISPLAY_SIZE {
		//println("Writing to display byte")
	}

	m.Map[location] = b
}

func (m *Memory) WriteHalf(h uint16, location uint32) {
	m.WriteByte(byte(h&uint16(0xFF)), location)
	m.WriteByte(byte((h&uint16(0xFF00))>>8), location+1)
	if location >= VIRT_DISPLAY && location < VIRT_DISPLAY+VIRT_DISPLAY_SIZE {
		//println("Writing to display half")
	}
}

func (m *Memory) WriteWord(w uint32, location uint32) {
	if location >= PLIC_PRIORITY && location <= PLIC_INT_ENABLE+0x100 {
		_ = m.Plic.Write(w, location, m.Cpu)
		return
	}
	if location >= PLIC_THRESHOLD && location <= PLIC_THRESHOLD+0x100 {
		_ = m.Plic.Write(w, location, m.Cpu)
		return
	}

	if location >= BASE_CLINT && location <= CLINT_END {
		_ = m.Clint.Write(w, location, m.Cpu)
		return
	}

	if location >= VIRT_DISPLAY && location < VIRT_DISPLAY+VIRT_DISPLAY_SIZE {
		_ = m.Display.Write(w, location)
		return
	}

	m.WriteByte(byte(w&uint32(0xFF)), location)
	m.WriteByte(byte((w&uint32(0xFF00))>>8), location+1)
	m.WriteByte(byte((w&uint32(0xFF0000))>>16), location+2)
	m.WriteByte(byte((w&uint32(0xFF000000))>>24), location+3)

}

func (m *Memory) ReadByte(location uint32) byte {
	if location >= VIRT_UART0 && location <= VIRT_UART0+0x16 {
		b, _ := m.Uart.Read(location - VIRT_UART0)
		return b
	}
	return m.Map[location]
}

func (m *Memory) ReadHalf(location uint32) uint16 {
	return uint16(m.ReadByte(location)) | (uint16(m.ReadByte(location+1)) << 8)
}

func (m *Memory) ReadWord(location uint32) uint32 {
	if location >= PLIC_PRIORITY && location <= PLIC_INT_ENABLE+0x100 {
		return m.Plic.Read(location)
	}
	if location >= PLIC_THRESHOLD && location <= PLIC_THRESHOLD+0x100 {
		return m.Plic.Read(location)
	}

	if location >= VIRT_DISPLAY && location < VIRT_DISPLAY+VIRT_DISPLAY_SIZE {
		return m.Display.Screen[location-VIRT_DISPLAY]
	}

	if location >= BASE_CLINT && location == 0x200BFFC {
		return uint32(((m.Clint.Mtime << 32) >> 32) & 0xFFFFFFFF)
	}

	if location >= BASE_CLINT && location == 0x200BFF8 {
		return uint32((m.Clint.Mtime >> 32) & 0xFFFFFFFF)
	}

	return uint32(int32(uint32(m.ReadByte(location)) |
		(uint32(m.ReadByte(location+1)) << 8) |
		(uint32(m.ReadByte(location+2)) << 16) |
		(uint32(m.ReadByte(location+3)) << 24)))
}
