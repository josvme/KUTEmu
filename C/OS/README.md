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