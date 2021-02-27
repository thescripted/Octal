package chip8

import (
	"fmt"
	"testing"
)

// TestJump calls the 1NNN opcode, checking for a valid pc result.
func TestJump(t *testing.T) {
	dummyDrawChannel := make(chan int)
	chip := New()
	chip.memory[chip.pc] = 0x1A
	chip.memory[chip.pc+1] = 0x00
	chip.emulateCycle(dummyDrawChannel)
	correct := uint16(0xA00)
	if chip.pc != correct {
		t.Errorf("Got %#x, expected %#x", chip.pc, correct)
	}

}

// TestFont calls opcode FX29, checks if the index is set to the correct font character.
func TestFont(t *testing.T) {
	dummyDrawChannel := make(chan int)
	var fontChar byte = 0x5
	chip := New()
	chip.V[0] = fontChar
	chip.memory[chip.pc], chip.memory[chip.pc+1] = 0xF0, 0x29
	chip.emulateCycle(dummyDrawChannel)
	fmt.Println("Chip Registers:", chip.V)
	fmt.Printf("Chip Index: %#x\n", chip.index)
	var fontLocation uint16 = fontStart + uint16(fontChar)*fontLength
	if chip.index != fontLocation {
		t.Errorf("Got %#x, expected %#x", chip.index, fontLocation)
	}

}

// TestDecode does test
func TestDecode(t *testing.T) {
	chip := New()
	chip.memory[chip.pc] = 0xA2
	chip.memory[chip.pc+1] = 0xF0
	// chip.Decode()
	// 	if chip.inst != 0xA2F0 {
	//		t.Errorf("Got %#x, expected 0xA2F0", chip.inst)
	//	}
}

func TestExecute(t *testing.T) {
	chip := New()
	chip.memory[chip.pc] = 0xA2
	chip.memory[chip.pc+1] = 0xF0
	// chip.Execute()
	//	if chip.index != 0x02F0 {
	//		t.Errorf("Got %#x, expected 0x02F0", chip.index)
	//	}

}

func TestLoadProgram(t *testing.T) {
	chip := New()
	chip.LoadProgram("in.oct")
	// chip.Decode()
	//	if chip.inst != 0xA2F0 {
	//		t.Errorf("Got %#x, expected 0xA2F0", chip.inst)
	//	}
}
