package main

import "fmt"
import "math/bits"
import "os"

//
const progStart = 0x200
const memSize = 4096

// Chip8 is our emulated processor state
type Chip8 struct {
	inst       uint16
	memory     []uint8
	v          [16]uint8      // register block
	index      uint16         // index reg
	pc         uint16         // program counter
	gfx        [64 * 32]uint8 // pixel array for graphics
	delayTimer uint8
	soundTimer uint8
	stack      [16]uint16
	sp         uint16
}

/*
0x000-0x1FF - Chip 8 interpreter (contains font set in emu)
0x050-0x0A0 - Used for the built in 4x5 pixel font set (0-F)
0x200-0xFFF - Program ROM and work RAM
*/

// LoadProgram loads the program from a file into the Chip8's memory.
func (c *Chip8) LoadProgram(prog string) {
	progFile, err := os.Open(prog)
	if err != nil {
		panic(err)
	}
	_, err = progFile.Read(c.memory[progStart:])
	if err != nil {
		panic(err)
	}
}

// Init initializes the chip8 instance.
func (c *Chip8) Init() {
	c.inst = 0
	c.index = 0
	c.pc = progStart
	c.sp = 0
	c.delayTimer = 0
	c.soundTimer = 0
	c.memory = make([]uint8, memSize)
}

// Decode decodes a single instruction.
func (c *Chip8) Decode() {
	topByte := bits.RotateLeft16(uint16(c.memory[c.pc]), 8) // shift the top byte up 8
	bottomByte := uint16(c.memory[c.pc+1])
	c.inst = topByte | bottomByte
}

func topNibble(i uint16) uint16 {
	return (i & 0xF000) >> 12
}
func bottomNibble(i uint16) uint16 {
	return (i & 0x000F)
}

func targetAddr(i uint16) uint16 {
	return (i & 0x0FFF)
}

// SetIndex sets the index register if current inst is ANNN
func (c *Chip8) SetIndex() {
	c.index = c.inst & 0x0FFF
}

// Execute executes a single instruction.
func (c *Chip8) Execute() {
	c.Decode()
	if c.inst == 0xFFFF {
		fmt.Println("End of program")
		os.Exit(0)
	}
	top := topNibble(c.inst)
	switch top {
	case 0xA:
		fmt.Println("0xA")
		c.SetIndex()
	case 0x0:
		switch bottomNibble(c.inst) {
		case 0x0:
			fmt.Println("clear screen")
			break
		case 0xE:
			fmt.Println("ret")
		}
	case 0x1:
		fmt.Printf("jump: %#x\n", targetAddr(c.inst))
	case 0x2:
		fmt.Printf("call sub at: %#x\n", targetAddr(c.inst))
	}

	c.pc += 2

}

func main() {
	chip := new(Chip8)
	chip.Init()
	// chip.LoadProgram("in.oct")
	chip.LoadProgram("../hapax8asm/test.bin")
	for {
		chip.Execute()
	}
}
