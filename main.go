package main

import (
	"log"

	"github.com/thescripted/hapax8/chip8"
)

// main: A CLI tool to initialize the chip-8 emulator and a menu-selection.
func main() {
	chip := chip8.New()
	chip.LoadProgram("./out.bin")
	err := chip.Run()
	if err != nil {
		log.Fatalf("emulator failed to start: %s", err)
	}
}
