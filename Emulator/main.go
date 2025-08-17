package main

import "riscv/emulator"

func main() {
	emu := emulator.NewEmulator()
	emu.Run()
}
