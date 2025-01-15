# Cross Compiling 
https://wiki.nixos.org/wiki/Cross_Compiling

## Points to Note
* Once we are in `nix develop .#riscv`, $CC and others points to the correct binary.
* The cc(compiler) and as(assembler) like as defacto in linux is named spaced by the architecture. For example, riscv32-none-elf-cc (for riscv)

## OBJCPY elf to .img (flat file)
```shell
$OBJCOPY -O binary result/bin/hellomake hello.img
```

## OBJDUMP .img (disassemble)
```shell
$OBJDUMP -D hello.img -b binary -m riscv:rv32
```

## Running Qemu with flat file
```shell
qemu-system-riscv32 -machine virt -bios hello.img 
```