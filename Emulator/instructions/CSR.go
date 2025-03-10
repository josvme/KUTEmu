package instructions

import (
	"fmt"
	"slices"
)

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

// More
const MSTATUSH = 0x310

type CSR struct {
	Registers [4096]uint32
}

type CSRPerm struct {
	mode      uint32
	privilege uint32
}

func getDetailsForCSR(csReg uint32) CSRPerm {
	// The top two bits (csr[11:10]) indicate whether the register is read/write (00, 01, or 10) or read-only (11).
	// The next two bits (csr[9:8]) encode the lowest privilege level that can access the CSR
	mode := csReg >> 10
	privilege := (csReg >> 8) & 0b0011
	return CSRPerm{
		mode:      mode,
		privilege: privilege,
	}
}

// 11 => read only
// 00, 01, 10 => read and write
// So 3 is read only, others are read write
// 3 for machine, 1 for supervisor, 2 for hypervisor, 0 for user

func getRWMode(csrReg uint32, currentCpuLevel uint32) uint32 {
	// Check if CSR is legal
	// Nothing bad should happen here as a CSR register is tied to its currentCpuLevel and privilege
	details := getDetailsForCSR(csrReg)
	if currentCpuLevel >= details.privilege {
		return details.mode
	}
	// This should be an exception
	return 4
}

func (csr *CSR) isInterruptEnabled(mode uint32) bool {
	status := ToMStatusReg(csr.Registers[MSTATUS])
	switch mode {
	case 0:
		return status.mie > 0
	case 1:
		return status.sie > 0
	case 3:
		return status.uie > 0
	}
	// This shouldn't happen
	fmt.Println("This shouldn't happen in interrupt enabled.......")
	return false
}

func (csr *CSR) setPreviousPrivilege(currentExecutionMode uint32, mode uint32, value uint32) {
	// First check if privilege is high enough to perform operation
	if currentExecutionMode < mode {
		panic("CurrentExecution mode should be higher")
	}
	status := ToMStatusReg(csr.Registers[MSTATUS])
	switch mode {
	case 0:
		status.mpp = value
		csr.Registers[MSTATUS] = FromMStatusReg(status)
	case 1:
		status.spp = value
		csr.Registers[MSTATUS] = FromMStatusReg(status)
	}
	// This shouldn't happen
	fmt.Println("This shouldn't happen in set previous privilege.......")
	return
}

func (csr *CSR) setPreviousInterruptEnabled(currentExecutionMode uint32, mode uint32, value uint32) {
	// First check if privilege is high enough to perform operation
	if currentExecutionMode < mode {
		panic("CurrentExecution mode should be higher")
	}
	status := ToMStatusReg(csr.Registers[MSTATUS])
	switch mode {
	case 0:
		status.mpie = value
		csr.Registers[MSTATUS] = FromMStatusReg(status)
	case 1:
		status.spie = value
		csr.Registers[MSTATUS] = FromMStatusReg(status)
	case 3:
		status.upie = value
		csr.Registers[MSTATUS] = FromMStatusReg(status)
	}
	// This shouldn't happen
	fmt.Println("This shouldn't happen in set previous interrupt .......")
	return
}

func (csr *CSR) setInterruptEnabled(currentExecutionMode uint32, mode uint32, value uint32) {
	// First check if privilege is high enough to perform operation
	if currentExecutionMode < mode {
		panic("CurrentExecution mode should be higher")
	}
	status := ToMStatusReg(csr.Registers[MSTATUS])
	switch mode {
	case 0:
		status.mie = value
		csr.Registers[MSTATUS] = FromMStatusReg(status)
	case 1:
		status.sie = value
		csr.Registers[MSTATUS] = FromMStatusReg(status)
	case 3:
		status.uie = value
		csr.Registers[MSTATUS] = FromMStatusReg(status)
	}

	// This shouldn't happen
	fmt.Println("This shouldn't happen in set interrupt .......")
	return
}

func (csr *CSR) GetValue(csrReg uint32, currentExecutionMode uint32, cpu *Cpu) uint32 {
	// Check if CSR itself is valid
	if !csr.isCSRValid(csrReg) {
		// illegal access
		csr.handleIllegalGet(csrReg, currentExecutionMode, uint32(3), cpu)
	}
	perm := getRWMode(csrReg, currentExecutionMode)
	if perm > 3 {
		// illegal access
		csr.handleIllegalGet(csrReg, currentExecutionMode, uint32(3), cpu)
	}
	switch {
	case csrReg == MISA:
		return (1 << 30) | 1 | (1 << 8) | (0 << 18) // 18 is supervisor mode
	case csrReg == MVENDORID:
		return 0
	case csrReg == MARCHID:
		return 0
	case csrReg == MIMPID:
		return 0
	case csrReg == MHARTID:
		return 0 // Single core system
	}
	return csr.Registers[csrReg]
}

func (csr *CSR) isCSRValid(reg uint32) bool {
	r := int(reg)
	v := []int{USTATUS, UIE, UTVEC, USCRATCH, UEPC, UCAUSE, UTVAL, UIP, CYCLE, TIME, INSTRET, CYCLEH, TIMEH,
		SSTATUS, SEDELEG, SIDELEG, SIE, STVEC, SCOUNTEREN, SSCRATCH, SEPC, SCAUSE, STVAL, SIP, SRW,
		MVENDORID, MARCHID, MIMPID, MHARTID, MSTATUS, MISA, MEDELEG, MIDELEG, MIE, MTVEC, MCOUNTEREN,
		MSCRATCH, MEPC, MCAUSE, MTVAL, MIP, MCYCLE, MINSTRET, MCYCLEH, MSTATUSH}
	t := slices.Contains(v, r)
	return t
}

func (csr *CSR) handleIllegalGet(reg uint32, cpulevel uint32, exception uint32, cpu *Cpu) {
	// set corresponding bits in CSR registers
	//fmt.Println("handling illegal access")
	csr.handleExceptions(cpulevel, exception, cpu)
}

func (csr *CSR) SetValue(csrReg uint32, value uint32, currentExecutionMode uint32, cpu *Cpu) {
	// Check if CSR itself is valid
	if !csr.isCSRValid(csrReg) {
		// illegal access
		csr.handleExceptions(currentExecutionMode, uint32(3), cpu)
	}
	// First check if privilege is high enough to perform operation
	rwmode := getRWMode(csrReg, currentExecutionMode)
	if rwmode >= 3 {
		// Attempts to access a non-existent CSR raise an illegal instruction exception. Attempts to access a
		//CSR without appropriate privilege level or to write a read-only register also raise illegal instruction
		//exceptions. A read/write register might also contain some bits that are read-only, in which case
		//writes to the read-only bits are ignored.
		csr.handleExceptions(currentExecutionMode, uint32(3), cpu)
		return
	}
	csr.Registers[csrReg] = value
}

func (csr *CSR) handleExceptions(mode uint32, exception uint32, cpu *Cpu) {
	// Save pc to mepc / corresponding
	// move PC to interrupt / exception handler after checking mdeleg and meip / corresponding registers
	// mtvec has different address based on modes see page 24, https://people.eecs.berkeley.edu/~krste/papers/riscv-priv-spec-1.7.pdf
	// A trap in privilege level P causes a jump to the address mtvec + P ×0x40. Non-maskable interrupts
	//cause a jump to address mtvec + 0xFC.

	// set mcause & scause register
	switch mode {
	case 1:
		csr.Registers[SCAUSE] = csr.Registers[SCAUSE] | exception
	default:
		csr.Registers[MCAUSE] = csr.Registers[MCAUSE] | exception
	}

	// check if delegation is added for exceptions
	medelegReg := csr.Registers[MEDELEG]
	if medelegReg&uint32(1<<3) == 1 {
		// go to supervisor trap
		csr.Registers[SEPC] = cpu.PC
		sstatus := ToMStatusReg(SSTATUS)
		sstatus.spp = mode
		csr.Registers[SSTATUS] = FromMStatusReg(sstatus)
		// This is very crappy, to accomodate +4
		stvec := ToMtvecReg(csr.Registers[STVEC])
		cpu.PC = stvec.base - 4
		cpu.CurrentMode = 1
	} else {
		// run machine trap
		csr.Registers[MEPC] = cpu.PC
		mstatus := ToMStatusReg(MSTATUS)
		mstatus.mpp = mode
		csr.Registers[MSTATUS] = FromMStatusReg(mstatus)
		// This is very crappy, to accomodate +4
		mtvec := ToMtvecReg(csr.Registers[MTVEC])
		cpu.PC = mtvec.base - 4
		cpu.CurrentMode = 3
		//{
		//	// only for interrupts
		//	cpu.PC = (mtvec.base + 4*exception) - 4
		//}
	}

	switch exception {
	case 3:
		//fmt.Println("handling illegal access of value")
	}
}

/*
When MODE=Direct, all traps into machine mode cause the pc to be set to the address in the BASE field.
When MODE=Vectored, all synchronous exceptions into machine mode cause the pc to be set to the address in the BASE field,
whereas interrupts cause the pc to be set to the address in the BASE field plus four times the interrupt cause number.
For example, a machine-mode timer interrupt (see Table [mcauses]) causes the pc to be set to BASE+0x1c.
*/
type MtvecReg struct {
	mode uint32
	base uint32
}

func ToMtvecReg(r uint32) MtvecReg {
	mode := getBitsAsUInt32(r, 0, 1)
	base := getBitsAsUInt32(r, 2, 31)
	return MtvecReg{
		mode: mode,
		base: base * 4,
	}
}

func FromMtvecReg(r MtvecReg) uint32 {
	return r.mode | (r.base>>2)<<2
}

type MscratchReg struct {
	value uint32
}

type MepcReg struct {
	// lowest 1 bit is always zero
	value uint32
}

// Has fault address can be page fault/ bad instruction
type MtvalReg struct {
}

// Matches layout of bits in mip register
// Restricted views of the mip and mie registers appear as the sip/sie, and uip/uie registers in
// S-mode and U-mode respectively. If an interrupt is delegated to privilege mode x by setting a bit in
// the mideleg register, it becomes visible in the x ip register and is maskable using the x ie register.
// Otherwise, the corresponding bits in x ip and x ie appear to be hardwired to zero.
type MedelegReg struct {
}

type MidelegReg struct {
}

// Used to identify the exception / interrupt type and jump to corresponding handler
type McauseReg struct {
	exceptionCode uint32
	isInterrupt   uint32
}

// Memory mapped, 64 bit
type Mtime struct {
}

// Memory mapped, 64 bit
type Mtimecmp struct {
}

type MipReg struct {
	usip uint32
	ssip uint32
	msip uint32
	utip uint32
	stip uint32
	mtip uint32
	ueip uint32
	seip uint32
	meip uint32
}

type MieReg struct {
	usie uint32
	ssie uint32
	msie uint32
	utie uint32
	stie uint32
	mtie uint32
	ueie uint32
	seie uint32
	meie uint32
}

type MStatusReg struct {
	uie  uint32
	sie  uint32
	mie  uint32
	upie uint32
	spie uint32
	mpie uint32
	spp  uint32
	mpp  uint32
	fs   uint32
	xs   uint32
	mprv uint32
	sum  uint32
	mxr  uint32
	tvm  uint32
	tw   uint32
	tsr  uint32
	sd   uint32
}

func maskBitSet(start uint32, end uint32) uint32 {
	result := uint32(0)
	for i := start; i <= end; i++ {
		result = result | 1
		result = result << 1
	}

	for i := uint32(0); i < start; i++ {
		result = result << 1
	}
	return result
}

func ToMStatusReg(r uint32) MStatusReg {
	uie := getBitsAsUInt32(r, 0, 0)
	sie := getBitsAsUInt32(r, 1, 1)
	mie := getBitsAsUInt32(r, 3, 3)
	upie := getBitsAsUInt32(r, 4, 4)
	spie := getBitsAsUInt32(r, 5, 5)
	mpie := getBitsAsUInt32(r, 7, 7)
	spp := getBitsAsUInt32(r, 8, 8)
	mpp := getBitsAsUInt32(r, 11, 12)
	fs := getBitsAsUInt32(r, 13, 14)
	xs := getBitsAsUInt32(r, 15, 16)
	mprv := getBitsAsUInt32(r, 17, 17)
	sum := getBitsAsUInt32(r, 18, 18)
	mxr := getBitsAsUInt32(r, 19, 19)
	tvm := getBitsAsUInt32(r, 20, 20)
	tw := getBitsAsUInt32(r, 21, 21)
	tsr := getBitsAsUInt32(r, 22, 22)
	sd := getBitsAsUInt32(r, 31, 31)

	return MStatusReg{
		uie:  uie,
		sie:  sie,
		mie:  mie,
		upie: upie,
		spie: spie,
		mpie: mpie,
		spp:  spp,
		mpp:  mpp,
		fs:   fs,
		xs:   xs,
		mprv: mprv,
		sum:  sum,
		mxr:  mxr,
		tvm:  tvm,
		tw:   tw,
		tsr:  tsr,
		sd:   sd,
	}
}

func FromMStatusReg(r MStatusReg) uint32 {
	return r.uie |
		r.sie<<1 |
		r.mie<<3 |
		r.upie<<4 |
		r.spie<<5 |
		r.mpie<<7 |
		r.spp<<8 |
		r.mpp<<11 |
		r.fs<<13 |
		r.xs<<15 |
		r.mprv<<17 |
		r.sum<<18 |
		r.mxr<<19 |
		r.tvm<<20 |
		r.tw<<21 |
		r.tsr<<22 |
		r.sd<<31
}
