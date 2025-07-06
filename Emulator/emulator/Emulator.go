package emulator

import (
	"fmt"
	"os"
	"riscv/instructions"
)

type Emulator struct {
	cpu instructions.Cpu
}

const VIRT_DRAM = 0x80000000
const VIRT_OPENSBI_START = 0x80200000
const VIRT_VIRTIO = 0x10001000
const DTB = 0x87e00000

func (e *Emulator) Run() {
	uart := instructions.NewUART()
	plic := &instructions.Plic{
		Priority:  1,
		Pending:   make([]uint32, 0),
		Enable:    0,
		Threshold: 1,
		Claim:     0,
		Complete:  0,
	}
	memory := &instructions.Memory{Map: make(map[uint32]byte), Uart: uart, Plic: plic}
	csr := &instructions.CSR{
		Registers: make([]uint32, 4096),
	}
	cpu := instructions.Cpu{
		PC:          0x80000000,
		Registers:   [32]uint32{},
		Memory:      memory,
		CSR:         csr,
		CurrentMode: 3,
	}
	memory.SetCpu(&cpu)
	cpu.Memory = memory
	//f, _ := os.Create("e.log")
	//defer f.Close()

	//body, _ := os.ReadFile("./../C/OS/fw_dynamic.bin")
	path := os.Getenv("OBJ_PATH")
	path = "/home/josv/Projects/RiscV/Tests/os.img"
	body, _ := os.ReadFile(path)
	_ = memory.LoadBytes(body, VIRT_DRAM)

	//body, _ = os.ReadFile("two.dtb")
	//_ = memory.LoadBytes(body, DTB)

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
		//mstatus := instructions.ToMStatusReg(cpu.CSR.GetValue(instructions.MSTATUS, cpu.CurrentMode, &cpu))
		//fmt.Println(fmt.Sprintf("before mstatus: %x", mstatus))

		// f.WriteString(fmt.Sprintf("0x%x\n", cpu.PC))
		_ = cpu.ExecInst(inst)

		//mstatus = instructions.ToMStatusReg(cpu.CSR.GetValue(instructions.MSTATUS, cpu.CurrentMode, &cpu))
		//fmt.Println(fmt.Sprintf("after mstatus: %x", mstatus))

		if uart.DataExistsToRead() && cpu.CSR.Registers[instructions.MIE] > 0 {
			// trigger a plic interrupt
			plic.TriggerInterrupt(10, &cpu)
		}
		_ = cpu.HandleInterrupts(inst.Operation())

		//mstatus = instructions.ToMStatusReg(cpu.CSR.GetValue(instructions.MSTATUS, cpu.CurrentMode, &cpu))
		//fmt.Println(fmt.Sprintf("post interrupt mstatus: %x", mstatus))
		// Handle interrupts / exceptions
	}
}
