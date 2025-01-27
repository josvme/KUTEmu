#!/bin/bash

for i in *.s;
do
  $AS -march=rv32ia -mabi=ilp32 -o $i.o -c $i > /dev/null 2>&1
  $LD -T inst.ld --no-dynamic-linker -m elf32lriscv -static -s -o bins/$i.bin $i.o > /dev/null 2>&1
  retVal=$?
  if [ $retVal -ne 0 ]; then
      echo "Loader Failed for $i"
      continue
  fi
  rm *.o
  $OBJCOPY -O binary bins/$i.bin bins/$i.img
  rm bins/*.bin
  export OBJ_PATH=`realpath bins/$i.img`
  export MODE="test"
  echo "Running: $i"
  ./riscv
  rm bins/$i.img
done