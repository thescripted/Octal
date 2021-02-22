package chip8

import (
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
