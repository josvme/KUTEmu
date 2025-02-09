.section .text
.globl start

start:
    la      t0, supervisor
    csrw    mepc, t0
    la      t1, m_trap
    csrw    mtvec, t1
    li      t2, 0x1800
    csrc    mstatus, t2
    li      t3, 0x800
    csrs    mstatus, t3
    li      t4, 0x100
    csrs    medeleg, t4
    mret

m_trap:
    csrr    t0, mepc
    csrr    t1, mcause
    la      t2, supervisor
    csrw    mepc, t2
    addi    s0, s0, 1
    addi s1, x0, 10
    beq s0, s1, success
    mret

supervisor:
    la      t0, user
    csrw    sepc, t0
    la      t1, s_trap
    csrw    stvec, t1
    sret

s_trap:
    csrr    t0, sepc
    csrr    t1, scause
    ecall

user:
    csrr    t0, instret
    ecall

failure:
	li a0, 0
	li a7, 93
	ecall

success:
 	li a0, 42
 	li a7, 93
 	ecall