# The arch is riscv32-none-elf-

helloc: hello.o hello.ld
	${LD} -T hello.ld --no-dynamic-linker -m elf32lriscv -static -nostdlib  -s -o helloc helloc.o

hello.o: hello.c
	${CC} -march=rv32i -mabi=ilp32 -fno-stack-protector -r hello.c -o helloc.o
