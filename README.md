## KUTEmu

WIP: Trying to build a toy RiscV 32bit emulator for running Doom

## TODO:
* Combine all CSR registers which are shared.
* Fix ecall instruction
* Implement framebuffer driver
* Implement Supervisor mode
* Implement CLIC

## Others
Based on current execution privilege level, always adjust mstatus register. For example, all lower priv mode interrupts 
should be disabled and all higher mode interrupts enabled.
