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

## Connect gdb to Qemu
```shell
riscv32-none-elf-gdb bins/csr.s.bin -ex "target remote :1234"
```
Note: Install gdb dashboard and make below changes to .gdbinit file
```shell
define show_vars
    display /x $mstatus
    display /x $mepc
    display /x $sstatus
    display /x $sepc
end
```
And call show_vars once program is started, as display works once program is started.

Note: Above we use .bin for gdb (elf file with symbols) and .img (flat file) for qemu

## Resources
https://riscv-programming.org/book/riscv-book.html has a lot of details
https://luplab.gitlab.io/rvcodecjs/