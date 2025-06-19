#!/bin/bash

pushd ../Emulator
go build
cp riscv ../Tests/
popd

for i in *.dump;
do
  name="${i%.*}"
  riscv32-none-elf-objcopy -O binary $name $name.img
  export OBJ_PATH=`realpath $name.img`
  export MODE="test"
  echo "Running: $name"
  ./riscv
  rm $name.img
done