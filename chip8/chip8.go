package chip8

import (
	"fmt"
	"math/rand"
	"os"
)

const (
	// memSize is the size of the memory in the Chip-8 emulator
	memSize = 0x1000

	// gfxSize is the size of the graphics card rendered by pixel.
	gfxSize = 0x800

	// fontStart is where the font set begins in memory.
	fontStart = 0x50

	// fontLength is the length of an individual font.
	fontLength = 0x5

	// programStart is where all Chip-8 ROMs instruction start at in memory.
	programStart = 0x200

	// registerSize is the size of the Vx register in Chip-8.
	registerSize = 0x10

	// stackSize is the size of the call-stack used in Chip-8 call routines.
	stackSize = 0x10
)

// fontSet is the default font for Chip-8. This is loaded on init.
var fontSet = [80]byte{
	0xF0, 0x90, 0x90, 0x90, 0xF0, // 0
	0x20, 0x60, 0x20, 0x20, 0x70, // 1
	0xF0, 0x10, 0xF0, 0x80, 0xF0, // 2
	0xF0, 0x10, 0xF0, 0x10, 0xF0, // 3
	0x90, 0x90, 0xF0, 0x10, 0x10, // 4
	0xF0, 0x80, 0xF0, 0x10, 0xF0, // 5
	0xF0, 0x80, 0xF0, 0x90, 0xF0, // 6
	0xF0, 0x10, 0x20, 0x40, 0x40, // 7
	0xF0, 0x90, 0xF0, 0x90, 0xF0, // 8
	0xF0, 0x90, 0xF0, 0x10, 0xF0, // 9
	0xF0, 0x90, 0xF0, 0x90, 0x90, // A
	0xE0, 0x90, 0xE0, 0x90, 0xE0, // B
	0xF0, 0x80, 0x80, 0x80, 0xF0, // C
	0xE0, 0x90, 0x90, 0x90, 0xE0, // D
	0xF0, 0x80, 0xF0, 0x80, 0xF0, // E
	0xF0, 0x80, 0xF0, 0x80, 0x80, // F
}

// Chip8 contains the emulated processor state.
type Chip8 struct {
	Video      [gfxSize]byte // pixel array for graphics
	SoundTimer byte          // sound timer
	Key        [16]byte      // Key Bindings

	index      uint16             // index reg
	pc         uint16             // program counter
	sp         byte               // Stack pointer
	delayTimer byte               // delay timer
	memory     [memSize]byte      // memory of chip-8 VM
	v          [registerSize]byte // register block
	stack      [stackSize]uint16  // Call Stack
}

// opcode is a data structure containing the instruction.
type opcode struct {
	instruction uint16
	x           uint16
	y           uint16
	n           uint16
	lowerByte   uint16
	address     uint16
}

// New returns a new instance of a Chip8 VM.
func New() *Chip8 {
	chip := &Chip8{}
	chip.init()

	return chip
}

// LoadProgram loads the program from a file into the Chip8's memory.
func (c *Chip8) LoadProgram(program string) error {
	c.init()
	programFile, err := os.Open(program)
	if err != nil {
		return err
	}
	_, err = programFile.Read(c.memory[programStart:])
	if err != nil {
		return err
	}

	return nil
}

// PressKey turns on a key flag.
func (c *Chip8) PressKey(key uint) {
	c.Key[key] = 1
}

// ReleaseKey turns off a key flag.
func (c *Chip8) ReleaseKey(key uint) {
	c.Key[key] = 0
}

// Tick is the Fetch-Decode-Execute routine. It will process one `tick` of instruction.
func (c *Chip8) Tick() error {
	currentOpcode := uint16(c.memory[c.pc])<<8 | uint16(c.memory[c.pc+1])
	opcode := decode(currentOpcode)
	c.pc += 2

	registerX := &c.v[opcode.x]
	registerY := &c.v[opcode.y]
	flagRegister := &c.v[0xF]

	switch opcode.instruction {
	case 0x0000:
		switch opcode.lowerByte {
		case 0xE0: // Clear
			for i := range c.Video {
				c.Video[i] = 0
			}
		case 0xEE: // return from call routine
			c.pc = c.pop()
		}

	case 0x1000:
		c.pc = opcode.address

	case 0x2000:
		c.push(c.pc)
		c.pc = opcode.address

	case 0x3000:
		if *registerX == byte(opcode.lowerByte) {
			c.pc += 2
		}

	case 0x4000:
		if *registerX != byte(opcode.lowerByte) {
			c.pc += 2
		}

	case 0x5000:
		if *registerX == *registerY {
			c.pc += 2
		}

	case 0x6000:
		*registerX = byte(opcode.lowerByte)

	case 0x7000:
		*registerX += byte(opcode.lowerByte)

	case 0x8000:
		switch opcode.n {
		case 0x0:
			*registerX = *registerY

		case 0x1:
			*registerX |= *registerY

		case 0x2:
			*registerX &= *registerY

		case 0x3:
			*registerX ^= *registerY

		case 0x4:
			sum := *registerX + *registerY
			*flagRegister = 0
			if sum > 0xFF {
				*flagRegister = 1
			}
			*registerX = sum

		case 0x5:
			minuend := *registerX
			subtrahend := *registerY
			*flagRegister = 0
			if minuend > subtrahend {
				*flagRegister = 1
			}
			*registerX = minuend - subtrahend

		case 0x6:
			lsb := *registerY & 1
			*flagRegister = lsb
			*registerX = *registerY >> 1

		case 0x7:
			minuend := *registerY
			subtrahend := *registerX
			*flagRegister = 0
			if minuend > subtrahend {
				*flagRegister = 1
			}
			*registerX = minuend - subtrahend

		case 0xE:
			msb := (*registerY & (1 << 7)) >> 7 // sizeof(byte) = 8
			*flagRegister = msb
			*registerX = *registerY << 1
		}

	case 0x9000:
		if *registerX != *registerY {
			c.pc += 2
		}

	case 0xA000:
		c.index = opcode.address

	case 0xB000: // AMBIGUOUS: Should provide Configuration
		c.pc = opcode.address
		c.v[0x0] = byte(opcode.address)

	case 0xC000:
		*registerX = byte(rand.Intn(0x100)) & byte(opcode.lowerByte)

	case 0xD000: // Draw
		height := int(opcode.n)
		xCoordinate := int(*registerX) % 64
		yCoordinate := int(*registerY) % 32

		*flagRegister = 0 // set flag to zero
		for y := 0; y < height; y++ {
			pixel := c.memory[c.index+uint16(y)]
			for x := 0; x < 8; x++ {
				if (pixel & (0x80 >> x)) != 0 { // if pixel_item is on
					if c.Video[xCoordinate+x+((yCoordinate+y)*64)] == 1 { // and Video is also on
						*flagRegister = 1
					}
					c.Video[xCoordinate+x+((yCoordinate+y)*64)] ^= 1
				}
			}
		}

	case 0xE000:
		switch opcode.lowerByte {
		case 0x9E:
			if c.Key[*registerX] == 1 {
				c.pc += 2
			}

		case 0xA1:
			if c.Key[*registerX] == 0 {
				c.pc += 2
			}
		}

	case 0xF000:
		switch opcode.lowerByte {
		case 0x07:
			*registerX = c.delayTimer

		case 0x15:
			c.delayTimer = *registerX

		case 0x18:
			c.SoundTimer = *registerX

		case 0x1E:
			sum := c.index + uint16(*registerX)
			*flagRegister = 0
			if sum > 0xFFFF {
				*flagRegister = 1
			}
			c.index = sum

		case 0x0A:
			fmt.Println("I am being called.")
			keyPressed := false
			for i := range c.Key {
				if c.Key[i] == 1 {
					*registerX = byte(i)
					keyPressed = true
					break
				}
			}
			if !keyPressed { // wait for key input.
				c.pc -= 2
			}

		case 0x29:
			c.index = uint16(fontStart + *registerX*fontLength)

		case 0x33:
			c.memory[c.index] = *registerX / 100
			c.memory[c.index+1] = (*registerX / 10) % 10
			c.memory[c.index+2] = *registerX % 10

		case 0x55:
			// if 0 is provided, just add registerX to memory.
			if opcode.x == 0 {
				c.memory[c.index] = *registerX
			} else {
				var i uint16
				for i = 0; i <= opcode.x; i++ {
					c.memory[c.index+i] = c.v[i]
				}
			}
			c.index += opcode.x + 1 // AMBIGUOUS !!!

		case 0x65:
			// if 0 is provided, just read into registerX from memory.
			if opcode.x == 0x0 {
				*registerX = c.memory[c.index]
			} else {
				var i uint16
				for i = 0; i <= opcode.x; i++ {
					c.v[i] = c.memory[c.index+i]
				}
			}
			c.index += opcode.x + 1 // AMBIGUOUS !!!
		}
	}

	return nil
}

// EmulateTimer decrements the sound and delay timer.
func (c *Chip8) EmulateTimer() {
	if c.delayTimer > 0 {
		c.delayTimer--
	}
	if c.SoundTimer > 0 {
		c.SoundTimer--
		fmt.Println("Beeping!")
	}
}

// init will clear display, stack, registers (v), and memory.
func (c *Chip8) init() {
	c.pc = programStart // program counter starts at 0x200
	c.sp = 0
	c.index = 0
	c.delayTimer = 0
	c.SoundTimer = 0

	// clear Register.
	for i := range c.v {
		c.v[i] = 0
	}

	// clear call stack.
	for i := range c.stack {
		c.stack[i] = 0
	}

	// clear graphics.
	for i := range c.Video {
		c.Video[i] = 0
	}

	// init Key.
	for i := range c.Key {
		c.Key[i] = 0
	}

	// load the fontset into memory. By convention, fontset occupies 0x50-0x9F.
	for i, font := range fontSet {
		c.memory[fontStart+i] = font
	}
}

// decode reads and parses the first two memory address to use in execution.
func decode(current uint16) opcode {
	parsedOpcode := opcode{
		instruction: current & 0xF000,
		x:           current & 0x0F00 >> 8,
		y:           current & 0x00F0 >> 4,
		n:           current & 0x000F,
		lowerByte:   current & 0x00FF,
		address:     current & 0x0FFF,
	}

	return parsedOpcode
}

// push will push a reference to the previous pc, as long as the stack isn't full.
func (c *Chip8) push(v uint16) {
	if c.sp >= stackSize-1 {
		panic("Stack is full.")
	}
	c.stack[c.sp] = v
	c.sp++
}

// pop will pop the recently pushed pc, as long as the stack isn't empty.
func (c *Chip8) pop() uint16 {
	if c.sp <= 0 {
		panic("Stack is empty.")
	}
	c.sp--
	pc := c.stack[c.sp]
	return pc
}
