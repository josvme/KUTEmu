package emulator

import (
	"os"
	"riscv/instructions"
)

type Emulator struct {
	cpu instructions.Cpu
}

const VIRT_DRAM = 0x80000000
const VIRT_VIRTIO = 0x10001000

func (e *Emulator) Run() {
	memory := instructions.Memory{Map: make(map[uint32]byte)}
	cpu := instructions.Cpu{
		PC:        0x80000000,
		Registers: [32]uint32{},
		Memory:    memory,
	}
	body, _ := os.ReadFile("./../C/hello.img")

	_ = memory.LoadBytes(body, VIRT_DRAM)

	var b [4]byte
	for {
		b[0] = memory.Read(cpu.PC)
		b[1] = memory.Read(cpu.PC + 1)
		b[2] = memory.Read(cpu.PC + 2)
		b[3] = memory.Read(cpu.PC + 3)

		inst := instructions.DecodeBytes(b)

		//fmt.Println(fmt.Sprintf("Executing Bytes: %x, Opcode: %s, Operation: #%v PC: %x ", instructions.TransformLittleToBig(b), inst.Operation(), inst, cpu.PC))

		_ = cpu.ExecInst(inst)

	}
}
