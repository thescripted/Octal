package main

import (
	"fmt"
	"os"
	"time"

	"github.com/thescripted/hapax8/chip8"
	"github.com/veandco/go-sdl2/sdl"
)

var (
	// Chip is the CHIP_8 Virtual Machine
	Chip *chip8.Chip8
)

func main() { // callback function to our "main" routine

	Chip = chip8.New()

	// Move this over to a CLI tool.
	if err := Chip.LoadProgram("./rom/octojam2title.ch8"); err != nil { // load from CLI
		fmt.Println("Error:", err)
		os.Exit(1)
	}

	// initialize window & renderer.
	window, renderer, err := sdl.CreateWindowAndRenderer(800, 600, sdl.WINDOWEVENT_SHOWN)
	if err != nil {
		panic(err)
	}
	defer window.Destroy()
	defer renderer.Destroy()

	window.SetTitle("Testing SDL2")

	format := sdl.PIXELFORMAT_ABGR8888
	access := sdl.TEXTUREACCESS_STREAMING
	tex, err := renderer.CreateTexture(uint32(format), access, 128, 64)
	if err != nil {
		panic(err)
	}
	defer tex.Destroy()

	// Game Cycles
	clock := time.NewTicker(time.Millisecond)
	video := time.NewTicker(time.Second / 60)
	timer := time.NewTicker(time.Second / 60)
	fps := time.NewTicker(time.Second)
	frames := 0

	// Game Loop.
	for processEvent() {
		select {
		case <-fps.C: // FPS Capture
			fmt.Println("Frames:", frames)
			frames = 0
		case <-clock.C: // Emulate Cycle.
			Chip.EmulateCycle()
		case <-video.C: // Draw
			draw()
		case <-timer.C: // SoundTimer and DelayTimer
			Chip.EmulateTimer()
		default:
		}
		frames++
	}
}

// draw draws the graphics onto the screen..
func draw() {
}

// Process event register keys & other external game information.
func processEvent() bool {
	for event := sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
		switch event.(type) {
		case *sdl.QuitEvent:
			return false
		}
	}
	return true
}
