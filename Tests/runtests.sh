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
  echo $OBJ_PATH
  export MODE="test"
  echo "Running: $name"
  timeout 2 ./riscv
  rm $name.img
done