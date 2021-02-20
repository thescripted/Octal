package chip8

import (
	"fmt"
	"math/bits"
	"os"
)

//
const progStart = 0x200
const memSize = 4096

// Demonstration returns a test string.
func Demonstration() string {
	return "Hello, this is a demonstration."
}

// Chip8 is our emulated processor state
type Chip8 struct {
	inst       uint16
	memory     []byte        // memory of chip-8 VM
	register   [16]byte      // register block
	index      uint16        // index reg
	pc         uint16        // program counter
	gfx        [64 * 32]byte // pixel array for graphics
	stack      [16]uint16    // Call Stack
	sp         uint16        // Stack pointer
	delayTimer byte
	soundTimer byte
}

type opcode uint16

/*
0x000-0x1FF - Chip 8 interpreter (contains font set in emu)
0x050-0x0A0 - Used for the built in 4x5 pixel font set (0-F)
0x200-0xFFF - Program ROM and work RAM
*/

// New returns a new instance of a Chip8 VM.
func New() *Chip8 {
	chip := &Chip8{
		inst:       0,
		memory:     make([]byte, memSize),
		index:      0,
		pc:         progStart,
		sp:         0,
		delayTimer: 0,
		soundTimer: 0,
	}
	return chip
}

// Run will run an infinite game loop.
func (c *Chip8) Run() error {
	var err error
	for {
		c.Execute()
	}
	return err
}

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

// ToString GT Sucks
func (c *Chip8) ToString() string {
	return fmt.Sprintf("Chip State:\n\tinst: %#x\n\tindex: %#x\n\tpc: %#x\n\tsp: %d\n\t", c.inst, c.index, c.pc, c.sp)
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

// SetPC sets the PC register to the given address
func (c *Chip8) SetPC(newaddr uint16) {
	c.pc = newaddr
}

// IncPC increments the PC (adds 2 since the word is a short)
func (c *Chip8) IncPC() {
	c.pc += 2
}

// GetImm tech sucks
func (c *Chip8) GetImm(numDigs int) uint8 {
	switch numDigs {
	case 2:
		return uint8(c.inst & 0x00FF)
	case 3:
		return uint8(c.inst & 0x0FFF)
	default:
		panic("bad arg")
	}
}

// GetXReg tech sucks
func (c *Chip8) GetXReg() uint16 {
	return c.inst & 0x0F00
}

// GetYReg tech sucks
func (c *Chip8) GetYReg() uint16 {
	return c.inst & 0x00F0
}

// Math8 executes the correct math instruction based on the bottom nibble of an inst starting with 0x8.
func (c *Chip8) Math8() {
	x := c.GetXReg()
	y := c.GetYReg()
	xVal := c.v[x]
	yVal := c.v[y]
	switch bottomNibble(c.inst) {
	case 0x0:
		c.v[x] = yVal
	case 0x1:
		c.v[x] = xVal | yVal
	case 0x2:
		c.v[x] = xVal & yVal
	case 0x3:
		c.v[x] = xVal ^ yVal
	case 0x4:
		add := xVal + yVal
		if add > 255 {
			c.v[0xF] = 1
		}
		c.v[x] = uint8(add & 0xFF)
	case 0x5:
		c.v[x] = xVal - yVal
	case 0x6:
		c.v[x] = xVal >> 1
	case 0x7:
		c.v[x] = yVal - xVal
	case 0xE:
		c.v[x] = xVal << 1
	}
}

// Execute executes a single instruction.
func (c *Chip8) Execute() {
	c.Decode()
	fmt.Println(c.ToString())
	if c.inst == 0x0 {
		os.Exit(0)
	}
	top := topNibble(c.inst)
	switch top {
	case 0xA:
		c.SetIndex()
		c.IncPC()
	case 0x0:
		switch bottomNibble(c.inst) {
		case 0x0:
			fmt.Println("clear screen")
		case 0xE:
			fmt.Println("ret")
		}
		c.IncPC()
	// JUMP
	case 0x1:
		c.SetPC(targetAddr(c.inst))
	// CALL
	case 0x2:
		c.SetPC(targetAddr(c.inst))
	case 0x3:
		imm := c.GetImm(2)
		regNum := c.GetXReg()
		c.IncPC()
		if imm == c.v[regNum] {
			c.IncPC() // skip inst
		}
	case 0x4:
		imm := c.GetImm(2)
		regNum := c.GetXReg()
		c.IncPC()
		if imm != c.v[regNum] {
			c.IncPC()
		}
	case 0x5:
		x := c.GetXReg()
		y := c.GetYReg()
		c.IncPC()
		if c.v[x] == c.v[y] {
			c.IncPC()
		}
	case 0x6:
		x := c.GetXReg()
		imm := c.GetImm(2)
		c.v[x] = imm
		c.IncPC()
	case 0x7:
		x := c.GetXReg()
		imm := c.GetImm(2)
		c.v[x] += imm

	case 0x8:
		c.Math8()
		c.IncPC()
	case 0x9:
		x := c.GetXReg()
		y := c.GetYReg()
		c.IncPC()
		if c.v[x] != c.v[y] {
			c.IncPC()
		}

	}

}
