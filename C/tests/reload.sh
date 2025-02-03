#!/bin/bash

pushd ../../Emulator
go build
cp riscv ../C/tests/
popd
bash opcode.sh