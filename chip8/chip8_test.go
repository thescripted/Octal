package chip8

import "testing"

// TestDecode does test
func TestDecode(t *testing.T) {
	chip := new(Chip8)
	chip.Init()
	chip.memory[chip.pc] = 0xA2
	chip.memory[chip.pc+1] = 0xF0
	chip.Decode()
	if chip.inst != 0xA2F0 {
		t.Errorf("Got %#x, expected 0xA2F0", chip.inst)
	}
}

func TestExecute(t *testing.T) {
	chip := new(Chip8)
	chip.Init()
	chip.memory[chip.pc] = 0xA2
	chip.memory[chip.pc+1] = 0xF0
	chip.Execute()
	if chip.index != 0x02F0 {
		t.Errorf("Got %#x, expected 0x02F0", chip.index)
	}

}

func TestLoadProgram(t *testing.T) {
	chip := new(Chip8)
	chip.Init()
	chip.LoadProgram("in.oct")
	chip.Decode()
	if chip.inst != 0xA2F0 {
		t.Errorf("Got %#x, expected 0xA2F0", chip.inst)
	}
}
