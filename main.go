package main

import (
	"fmt"
	"os"
	"time"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/imdraw"
	"github.com/faiface/pixel/pixelgl"
	"github.com/thescripted/hapax8/chip8"
	"golang.org/x/image/colornames"
)

func run() { // callback function to our "main" routine
	chip := chip8.New()
	if err := chip.LoadProgram("./IBM Logo.ch8"); err != nil { // load from CLI
		fmt.Println("Error:", err)
		os.Exit(1)
	}

	drawSig := make(chan int)
	errSig := make(chan error)
	go chip.Run(drawSig, errSig)

	graphics := chip.Gfx

	cfg := pixelgl.WindowConfig{
		Title:  "GT Sucks!",
		Bounds: pixel.R(0, 0, 640, 320),
		VSync:  true,
	}

	win, err := pixelgl.NewWindow(cfg)
	if err != nil {
		panic(err)
	}

	imd := imdraw.New(nil)

	// FPS Calcuation
	var (
		frames = 0
		second = time.Tick(time.Second)
	)

	win.Clear(colornames.Gold)
	for !win.Closed() {
		select {
		case <-second: // FPS Capture
			win.SetTitle(fmt.Sprintf("%s | FPS: %d", cfg.Title, frames))
			frames = 0
		case <-drawSig: // Draw Capture
			imd.Clear()
			idx := 0
			for y := win.Bounds().Max.Y; y > win.Bounds().Min.Y; y -= 320 / 32 {
				for x := win.Bounds().Min.X; x < win.Bounds().Max.X; x += 640 / 64 {
					if graphics[idx]%2 == 0 { // memory pixels are binary. this needs to be more efficient.
						imd.Color = pixel.RGB(0, 0, 0)
					} else {
						imd.Color = pixel.RGB(1, 1, 1)
					}
					imd.Push(pixel.V(x, y), pixel.V(x+10, y+10))
					imd.Rectangle(0)
					idx++
				}
			}
			imd.Draw(win)
		default:
		}

		if win.JustPressed(pixelgl.KeyK) {
			fmt.Println("Pressed key K")
		}
		if win.JustPressed(pixelgl.KeyJ) {
			fmt.Println("Pressed key J")
		}
		frames++
		win.Update()
	}
}

func main() {
	pixelgl.Run(run) // enable pixelgl to capture main function
}
