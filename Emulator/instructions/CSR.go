package instructions

// User Trap Setup
const USTATUS = 0x000
const UIE = 0x004
const UTVEC = 0x005

// User Trap Handling
const USCRATCH = 0x040
const UEPC = 0x041
const UCAUSE = 0x042
const UTVAL = 0x043
const UIP = 0x044

// User Counters / Timers
const CYCLE = 0xC00
const TIME = 0xC01
const INSTRET = 0xC02
const CYCLEH = 0xC80
const TIMEH = 0xC81

// Supervisor Trap Setup
const SSTATUS = 0x100
const SEDELEG = 0x102
const SIDELEG = 0x103
const SIE = 0x104
const STVEC = 0x105
const SCOUNTEREN = 0x106

// Supervisor Trap Handling
const SSCRATCH = 0x140
const SEPC = 0x141
const SCAUSE = 0x142
const STVAL = 0x143
const SIP = 0x144
const SRW = 0x180

// Machine Information Registers
const MVENDORID = 0xF11
const MARCHID = 0xF12
const MIMPID = 0xF13
const MHARTID = 0xF14

// Machine Trap Setup
const MSTATUS = 0x300
const MISA = 0x301
const MEDELEG = 0x302
const MIDELEG = 0x303
const MIE = 0x404
const MTVEC = 0x305
const MCOUNTEREN = 0x306

// Machine Trap Handling
const MSCRATCH = 0x340
const MEPC = 0x341
const MCAUSE = 0x342
const MTVAL = 0x343
const MIP = 0x344

// Machine Counters / Timers
const MCYCLE = 0xB00
const MINSTRET = 0xB02
const MCYCLEH = 0xCB80

type CSR struct {
	Registers [4096]uint32
}

type CSRPerm struct {
	mode      uint32
	privilege uint32
}

func GetWritableBits(mode uint8, csrReg uint32) uint32 {
	return 0
}

func GetDetailsForCSR(csReg uint32) CSRPerm {
	mode := csReg << 30
	privilege := (csReg << 28) >> 30
	return CSRPerm{
		mode:      mode,
		privilege: privilege,
	}
}

// 11 => read only
// 00, 01, 10 => read and write
// So 3 is read only, others are read write
// 3 for machine, 1 for supervisor, 2 for hypervisor, 0 for user

func GetPermission(mode uint8, csrReg uint32) uint32 {
	// Nothing bad should happen here as a CSR register is tied to its mode and privilege
	details := GetDetailsForCSR(csrReg)
	if uint32(mode) >= details.mode {
		return details.privilege
	}
	// This should be an exception
	return 3
}

func CanWrite(mode uint8, csrReg uint32) bool {
	return false
}

func CanRead(mode uint8, csrReg uint32) bool {
	return false
}

func GetValue(csrReg uint32) uint32 {
	return 0
}

func SetValue(csrReg uint32) uint32 {
	return 0
}
