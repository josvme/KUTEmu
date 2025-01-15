# The arch is riscv32-none-elf-

hello: hello.o hello.ld
	${LD} -T hello.ld --no-dynamic-linker -m elf32lriscv -static -nostdlib -s -o hellomake hello.o

hello.o: hello.s
	${AS} -march=rv32i -mabi=ilp32 -o hello.o -c hello.s
