## KUTEmu

WIP: Trying to build a toy RiscV 32bit emulator for running Doom

## TODO:
* Remove OpenSBI and make a simple bootloader
* Combine all CSR registers which are shared.
* Fix ecall instruction
* Implement PLIC / CLIC
* Implement interrupt handling
* Implement VGA driver
* Implement WFI

## Others
Based on current execution privilege level, always adjust mstatus register. For example, all lower priv mode interrupts 
should be disabled and all higher mode interrupts enabled.
