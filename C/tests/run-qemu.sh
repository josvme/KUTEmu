#!/bin/bash

set -xue

QEMU=qemu-system-riscv32
# Start qemu
$QEMU -machine virt -cpu rv32,pmp=false -smp 1 -s -S -nographic -bios bins/csr.s.img