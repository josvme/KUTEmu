## KUTEmu

A Toy RiscV 32bit emulator

## Compiling OpenSBI
```shell
cd ~/Projects/RiscV/SBI/opensbi
export PLATFORM_RISCV_ISA=rv32ima
export PLATFORM_RISCV_XLEN=32
make LLVM=1 PLATFORM=generic
```