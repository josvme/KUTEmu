#!/bin/bash

pushd ../Emulator
go build
cp riscv ../Tests/
popd

filename="os.elf"
name="${filename%.*}"
riscv32-none-elf-objcopy -O binary $filename $name.img
export OBJ_PATH=`realpath $name.img`
echo $OBJ_PATH
./riscv