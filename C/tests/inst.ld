MEMORY
{
   BRAM (RWX) : ORIGIN = 0x00000000, LENGTH = 8024
}
SECTIONS {
  .text : {
	*.o(.text)
  } >BRAM
}