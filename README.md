# Chip-8 Emulator

This is a repo containing a CHIP-8 Emulator written in Go.

## Install

CHIP-8 Emulator relies on [SDL2](http://libsdl.org/index.php) for its GUI and keyboard inputs, so make sure that it is installed on your machine.


## Usage
To run this emulator, clone this repository into a local directory and `go build main.go -o chip8` in the root of that directory.
then, run `./chip8`, and follow the instructions to load a game that are currently available in the `rom/` directory.

An opcode test is also provided. To run the opcode tests, run `./chip8 --test-opcode`.
