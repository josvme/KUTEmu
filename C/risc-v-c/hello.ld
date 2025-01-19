MEMORY
{
   BRAM (RWX) : ORIGIN = 0x80000000, LENGTH = 128
}
SECTIONS {
  .text : {
	helloc.o(.text)
  } >BRAM
}