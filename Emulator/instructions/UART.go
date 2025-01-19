package instructions

import "fmt"

type UART struct {
}

func (u *UART) Write(b byte) error {
	if b == 0 {
		return nil
	}
	fmt.Printf("%c", b)
	return nil
}

// TODO
func (u *UART) Read() (byte, error) {
	return 0, nil
}
