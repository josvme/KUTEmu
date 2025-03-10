# RISC-V bare metal example

To build the simply run `make hello`

Run with QEMU: `qemu-system-riscv64 -machine virt -bios hellomake`

Note: We convert this into flatfile format before loading to Qemu. See [README](../README.md)

## Resources
https://popovicu.com/posts/bare-metal-programming-risc-v/

Run with compiled OpenSBI and qemu.
```shell
make hello
make qemu
```