package emulator

import (
	"fmt"
	"os"
	"riscv/instructions"
)

type Emulator struct {
	cpu instructions.Cpu
}

const VIRT_DRAM = 0x00000000
const VIRT_OPENSBI_START = 0x80200000
const VIRT_VIRTIO = 0x10001000

func (e *Emulator) Run() {
	memory := instructions.Memory{Map: make(map[uint32]byte)}
	cpu := instructions.Cpu{
		PC:          0x00000000,
		Registers:   [32]uint32{},
		Memory:      memory,
		CurrentMode: 3,
	}
	//path := os.Getenv("OBJ_PATH")
	//body, _ := os.ReadFile(path)
	// body, _ := os.ReadFile("./../C/risc-v-c/result/bin/hello.img")
	path := os.Getenv("OBJ_PATH")
	body, _ := os.ReadFile(path)
	_ = memory.LoadBytes(body, VIRT_DRAM)

	var b [4]byte
	for {
		b[0] = memory.ReadByte(cpu.PC)
		b[1] = memory.ReadByte(cpu.PC + 1)
		b[2] = memory.ReadByte(cpu.PC + 2)
		b[3] = memory.ReadByte(cpu.PC + 3)

		// fail on empty instructions
		if b[0] == 0 && b[1] == 0 && b[2] == 0 && b[3] == 0 {
			fmt.Printf("\nFail: Empty instruction. PC: %x\n", cpu.PC)
			cpu.PC += 4
			break
		}
		inst := instructions.DecodeBytes(b)

		//fmt.Println(fmt.Sprintf("Executing Bytes: %x, Opcode: %s, Operation: #%v PC: %x ", instructions.TransformLittleToBig(b), inst.Operation(), inst, cpu.PC))

		_ = cpu.ExecInst(inst)
		_ = cpu.HandleInterrupts()
		// Handle interrupts / exceptions

	}
}
