package instructions

type MtvecReg struct {
	mode uint32
	base uint32
}

func ToMtvecReg(r uint32) MtvecReg {
	mode := getBitsAsUInt32(r, 0, 1)
	base := getBitsAsUInt32(r, 2, 31)
	return MtvecReg{
		mode: mode,
		base: base,
	}
}

func FromMtvecReg(r MtvecReg) uint32 {
	return r.mode | r.base<<2
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
