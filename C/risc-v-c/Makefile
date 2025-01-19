# The arch is riscv32-none-elf-

helloc: hello.o hello.ld
	${LD} -T hello.ld -m elf32lriscv -nostdlib -o helloc helloc.o --no-relax
	${OBJCOPY} -O binary helloc hello.img


hello.o: hello.c
	${CC} -march=rv32i -mabi=ilp32 -fno-stack-protector -Wall -O2 -nostdlib -nostartfiles  -c hello.c -o helloc.o --verbose
