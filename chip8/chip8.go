package chip8

import (
	"fmt"
	"os"
	"time"
)

const (
	memSize      = 0x1000
	gfxSize      = 0x800
	progStart    = 0x200
	registerSize = 0x10
	stackSize    = 0x10
)

// Chip8 is our emulated processor state
type Chip8 struct {
	memory     []byte // memory of chip-8 VM
	registers  []byte // register block
	index      uint16 // index reg
	pc         uint16 // program counter
	Gfx        []byte // pixel array for graphics
	stack      stack  // Call Stack
	sp         byte   // Stack pointer
	delayTimer byte   // delay timer
	soundTimer byte   // sound timer
}

type opcode struct {
	instruction uint16
	second      uint16
	third       uint16
	fourth      uint16
	lowerByte   uint16
	memAddress  uint16
}

type stack []uint16

// ChipRuntimeError catch programmer's generated errors at runtime.
type ChipRuntimeError struct {
	lineno int
	Err    error
}

// ChipLoaderError catch file loading errors.
type ChipLoaderError struct {
	file string
	Err  error
}

// Should consider handling runtime error for stack limit here.
func (s stack) push(v uint16) stack {
	return append(s, v)
}

func (s stack) pop() (stack, uint16) {
	l := len(s)
	if l <= 0 {
		// handle runtime error if programmer pops from empty stack.
	}
	return s[:l-1], s[l-1]

}

func (e *ChipRuntimeError) Error() string {
	return fmt.Sprintf("Error at line %d: %s", e.lineno, e.Err.Error())
}

func (e *ChipLoaderError) Error() string {
	return fmt.Sprintf("Error with %s: %s", e.file, e.Err.Error())
}

// New returns a new instance of a Chip8 VM.
func New() *Chip8 {
	chip := &Chip8{}
	chip.initialize()
	// load the fontset.

	return chip
}

// Run will run an infinite game loop.
func (c *Chip8) Run(drawSig chan int, errSig chan error) {
	ticker := time.NewTicker(time.Second / 60) // 60Hz

	for {
		select {
		case <-ticker.C:
			c.emulateCycle(drawSig)
		default:
		}
	}
}

// LoadProgram loads the program from a file into the Chip8's memory.
func (c *Chip8) LoadProgram(prog string) error {
	c.initialize()
	progFile, err := os.Open(prog)
	if err != nil {
		return err
	}
	_, err = progFile.Read(c.memory[progStart:])
	if err != nil {
		return err
	}

	return nil
}

// initialize will clear display, stack, registers, and memory.
// this is used when the emulator begins or is reset with another game.
func (c *Chip8) initialize() {
	c.pc = progStart // program counter starts at 0x200
	c.sp = 0
	c.index = 0
	c.delayTimer = 0
	c.soundTimer = 0
	c.memory = make([]byte, memSize)
	c.registers = make([]byte, registerSize)
	c.Gfx = make([]byte, gfxSize)
	c.stack = make(stack, 0) // will need to be considerate of stacksize
}

// decode reads and parses the first two memory address to use in execution.
func decode(current uint16) opcode {
	parsedOpcode := opcode{
		instruction: current & 0xF000,
		second:      current & 0x0F00 >> 8,
		third:       current & 0x00F0 >> 4,
		fourth:      current & 0x000F,
		lowerByte:   current & 0x00FF,
		memAddress:  current & 0x0FFF,
	}

	return parsedOpcode
}

func (c *Chip8) emulateCycle(drawSig chan int) error {
	// fetch -- find a way to make this neater.
	currentOpcode := uint16(c.memory[c.pc])<<8 | uint16(c.memory[c.pc+1])
	opcode := decode(currentOpcode)
	c.pc += 2

	switch opcode.instruction { // exposes the first half-byte in the opcode
	case 0x0000:
		switch opcode.lowerByte {
		case 0xE0: // Clear
			for i := range c.Gfx {
				c.Gfx[i] = 0
			}
			drawSig <- 1 // Draw to Canvas.
		case 0xEE: // return from call routine
			c.stack, c.pc = c.stack.pop()
		}
	case 0x1000:
		c.pc = opcode.memAddress
	case 0x2000:
		c.stack = c.stack.push(c.pc)
		c.pc = opcode.memAddress
	case 0x3000:
		if c.registers[opcode.second] == byte(opcode.lowerByte) {
			c.pc += 2
		}
	case 0x4000:
		if c.registers[opcode.second] != byte(opcode.lowerByte) {
			c.pc += 2
		}
	case 0x5000:
		if c.registers[opcode.second] == c.registers[opcode.third] {
			c.pc += 2
		}
	case 0x6000:
		c.registers[opcode.second] = byte(opcode.lowerByte)
	case 0x7000:
		c.registers[opcode.second] += byte(opcode.lowerByte)
	case 0x8000:
		switch opcode.fourth {
		case 0x0:
			c.registers[opcode.second] = c.registers[opcode.third]
		case 0x1:
			c.registers[opcode.second] |= c.registers[opcode.third]
		case 0x2:
			c.registers[opcode.second] &= c.registers[opcode.third]
		case 0x3:
			c.registers[opcode.second] ^= c.registers[opcode.third]
		case 0x4:
			sum := c.registers[opcode.second] + c.registers[opcode.third]
			c.registers[0xF] = 0
			if sum > 0xFF {
				c.registers[0xF] = 1
			}
			c.registers[opcode.second] = sum
		case 0x5:
			minuend := c.registers[opcode.second]
			subtrahend := c.registers[opcode.third]
			c.registers[0xF] = 0
			if minuend > subtrahend {
				c.registers[0xF] = 1
			}
			c.registers[opcode.second] = minuend - subtrahend
		case 0x7:
			minuend := c.registers[opcode.third]
			subtrahend := c.registers[opcode.second]
			c.registers[0xF] = 0
			if minuend > subtrahend {
				c.registers[0xF] = 1
			}
			c.registers[opcode.second] = minuend - subtrahend
		}
	case 0x9000:
		if c.registers[opcode.second] != c.registers[opcode.third] {
			c.pc += 2
		}
	case 0xA000:
		c.index = opcode.memAddress
	case 0xD000: // Draw
		c.registers[0xF] = 0 // set flag to zero
		height := int(opcode.fourth)
		xCoord := int(c.registers[opcode.second]) % 64
		yCoord := int(c.registers[opcode.third]) % 32

		for y := 0; y < height; y++ {
			pixel := c.memory[c.index+uint16(y)]
			for x := 0; x < 8; x++ {
				if (pixel & (0x80 >> x)) != 0 { // if pixel_item is on
					if c.Gfx[xCoord+x+((yCoord+y)*64)] == 1 { // and Gfx is also on
						c.registers[0xF] = 1 // then turn the flag on.
					}
					c.Gfx[xCoord+x+((yCoord+y)*64)] ^= 1
				}
			}
		}
		drawSig <- 1 // Draw to Canvas.
	}
	return nil
}
