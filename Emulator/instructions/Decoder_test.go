package instructions

import (
	"fmt"
	"reflect"
	"testing"
)

func TestTransform(t *testing.T) {

	expected := uint32(0x06800513)
	got := TransformLittleToBig([4]byte{0x13, 0x05, 0x80, 0x06})

	if expected != got {
		t.Errorf("Expected %v, Got %v", expected, got)
	}
}

func TestDecodeOpcode(t *testing.T) {
	expected := byte(0b0010011)
	// addi x10, x0, 104
	got := decodeOpcode(uint32(0x06800513))

	if expected != got {
		t.Errorf("Expected %v, Got %v", expected, got)
	}
}

func PrettyPrintStructAsBinary(s interface{}) {
	v := reflect.ValueOf(s)
	typeOfS := v.Type()

	for i := 0; i < v.NumField(); i++ {
		fmt.Printf("Field: %s\tValue: %b\n", typeOfS.Field(i).Name, v.Field(i).Interface())
	}
}

func TestDecoderInstructions(t *testing.T) {
	decodeInst := func(inst uint32, want Inst) Inst {
		return want.Decode(inst)
	}
	type TestInst struct {
		inst uint32
		want Inst
	}

	tests := []TestInst{
		{
			inst: 0x007403B3,
			want: RI{
				Opcode: 0b0110011,
				RD:     0b00111,
				F3:     0b000,
				RS1:    0b01000,
				RS2:    0b00111,
				F7:     0b0000000,
			},
		},
		{
			// jalr
			inst: 0x00008067,
			want: II{
				Opcode: 0b1100111,
				RD:     0b00000,
				F3:     0b000,
				RS1:    0b00001,
				IIM:    0b000000000000,
			},
		},
		{
			inst: 0x06800513,
			want: II{
				Opcode: 0b0010011,
				RD:     0b01010,
				F3:     0b000,
				RS1:    0b00000,
				IIM:    0b000001101000,
			},
		},
		{
			// csrrw
			inst: 0x30401073,
			want: II{
				Opcode: 0b1110011,
				RD:     0b00000,
				F3:     0b001,
				RS1:    0b00000,
				IIM:    0b001100000100,
			},
		},
		{
			inst: 0xF99FF06F,
			want: JI{
				Opcode: 0b1101111,
				RD:     0b00000,
				JIM:    0b111111111111110011000,
			},
		},
		{
			inst: 0x00a58023,
			want: SI{
				Opcode: 0b0100011,
				SIM:    0b0000000,
				F3:     0b000,
				RS1:    0b01011,
				RS2:    0b01010,
			},
		},
		{
			inst: 0x06705063,
			want: BI{
				Opcode: 0b1100011,
				F3:     0b101,
				RS1:    0b00000,
				RS2:    0b00111,
				BIM:    0b0000001100000,
			},
		},
		{
			inst: 0x01643037,
			want: UI{
				Opcode: 0b0110111,
				RD:     0b00000,
				UIM1:   0b00000001011001000011,
			},
		},
		{
			inst: 0x0ff0000f,
			want: FI{
				Opcode: 0b0001111,
				RD:     0b00000,
				F3:     0b000,
				RS1:    0b00000,
				Succ:   0b1111,
				Pred:   0b1111,
				FM:     0b0000,
			},
		},
	}

	for _, tt := range tests {
		got := decodeInst(tt.inst, tt.want)
		if !reflect.DeepEqual(tt.want, got) {
			t.Errorf("Expected: %+v, Got: %+v, Instruction: %x", tt.want, got, tt.inst)
			fmt.Println("Want")
			PrettyPrintStructAsBinary(tt.want)
			fmt.Println("Got")
			PrettyPrintStructAsBinary(got)
		}

	}
}
