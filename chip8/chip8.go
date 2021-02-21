package chip8

import (
	"fmt"
	"os"
)

const (
	progStart    = 2
	memSize      = 4096
	gfxSize      = 64 * 32
	registerSize = 16
	stackSize    = 16
)

// Chip8 is our emulated processor state
type Chip8 struct {
	memory     []byte   // memory of chip-8 VM
	registers  []byte   // register block
	index      uint16   // index reg
	pc         uint16   // program counter
	gfx        []byte   // pixel array for graphics
	stack      []uint16 // Call Stack
	sp         byte     // Stack pointer
	delayTimer byte     // delay timer
	soundTimer byte     // sound timer
}

type opcode struct {
	instruction uint16
	second      uint16
	third       uint16
	fourth      uint16
	lowerByte   uint16
	memAddress  uint16
}

// ChipRuntimeError catch programmer's generated errors at runtime.
type ChipRuntimeError struct {
	Err    error
	lineno int
}

// ChipLoaderError catch file loading errors.
type ChipLoaderError struct {
	Err  error
	file string
}

func (e *ChipRuntimeError) Error() string {
	return fmt.Sprintf("Error at line %d: %s", e.lineno, e.Err.Error())
}

func (e *ChipLoaderError) Error() string {
	return fmt.Sprintf("Error with %s: %s", e.file, e.Err.Error())
}

/*
0x000-0x1FF - Chip 8 interpreter (contains font set in emu)
0x050-0x0A0 - Used for the built in 4x5 pixel font set (0-F)
0x200-0xFFF - Program ROM and work RAM
*/

// New returns a new instance of a Chip8 VM.
func New() *Chip8 {
	// initialize the system.
	chip := &Chip8{}
	chip.initialize()

	// load the fontset.

	return chip
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
	c.gfx = make([]byte, gfxSize)
	c.stack = make([]uint16, stackSize)
}

// Run will run an infinite game loop.
func (c *Chip8) Run() error {
	for {
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

// decode reads and parses the first two memory address to use in execution.
// invalid inputs should be tested here.
func (c *Chip8) decode(current uint16) opcode {
	var parsedOpcode opcode

	parsedOpcode.instruction = current & 0xF000
	parsedOpcode.second = current & 0x0F00 >> 8
	parsedOpcode.third = current & 0x00F0 >> 4
	parsedOpcode.fourth = current & 0x000F
	parsedOpcode.lowerByte = current & 0x00FF
	parsedOpcode.memAddress = current & 0x0FFF

	return parsedOpcode
}

func (c *Chip8) emulateCycle() error {
	// fetch -- find a way to make this neater.
	currentOpcode := uint16(c.memory[c.pc])<<8 | uint16(c.memory[c.pc+1])
	opcode := c.decode(currentOpcode)
	c.pc += 2

	// decode
	switch opcode.instruction { // exposes the first half-byte in the opcode
	}

	// execute
	return nil
}

func interpret(op opcode) uint16 {
	return uint16(0)
}

func execute(operation uint16) {
}
