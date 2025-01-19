MEMORY
{
   BRAM (RWX) : ORIGIN = 0x80000000, LENGTH = 0x200
}
SECTIONS
{
    everything :
    {
	. = ALIGN(4);
	helloc.o (.text)
        *(.*)
    } >BRAM
}