# Helpful Commands

We can check the addresses of functions/variables using llvm-nm
```bash
llvm-nm kernel.elf
```

Run objdump on file

```bash
llvm-objdump -d kernel.elf
```

Look into `kernel.map` file to see the output of linker to see what is where 

## Compiling OpenSBI
```shell
cd ~/Projects/RiscV/SBI/opensbi
export PLATFORM_RISCV_ISA=rv32ima
export PLATFORM_RISCV_XLEN=32
make LLVM=1 PLATFORM=generic
```

