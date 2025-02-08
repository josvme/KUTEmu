#!/bin/bash

set -xue

QEMU=qemu-system-riscv32
CC=clang

CFLAGS="-fuse-ld=lld -std=c11 -O2 -g3 -Wall -Wextra --target=riscv32-unknown-elf -fno-stack-protector -ffreestanding -nostdlib"

# Build the kernel
$CC $CFLAGS -Wl,-Tkernel.ld -Wl,-Map=kernel.map -o kernel.elf kernel.c -Wl,--no-dynamic-linker

# Start qemu
$QEMU -machine virt -bios default -nographic -serial mon:stdio --no-reboot -kernel kernel.elf