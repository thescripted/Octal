package main

import (
	"fmt"
	"github.com/veandco/go-sdl2/sdl"
	"time"
)

const (
	screenWidth  = 640
	screenHeight = 320
	pixelWidth   = screenWidth / 64
	pixelHeight  = screenHeight / 32
)

func main() {
	if err := sdl.Init(sdl.INIT_VIDEO); err != nil {
		panic(err)
	}

	window, err := sdl.CreateWindow("CHIP - 8.", sdl.WINDOWPOS_UNDEFINED, sdl.WINDOWPOS_UNDEFINED, screenWidth, screenHeight, 0)
	if err != nil {
		panic(err)
	}
	defer window.Destroy()

	renderer, err := sdl.CreateRenderer(window, -1, sdl.RENDERER_ACCELERATED)
	if err != nil {
		panic(err)
	}
	defer renderer.Destroy()

	sdl.SetHint(sdl.HINT_RENDER_SCALE_QUALITY, "linear")

	renderer.SetDrawColor(96, 128, 128, 255)
	renderer.Clear()

	pixel := sdl.Rect{
		X: 0,
		Y: 0,
		W: pixelWidth,
		H: pixelHeight,
	}

	var k uint8 = 0
	for i := 0; i < 64; i++ {
		pixel.X = int32(i * pixelWidth)
		for j := 0; j < 32; j++ {
			pixel.Y = int32(j * pixelHeight)
			renderer.SetDrawColor(255*k, 255*k, 255*k, 255)
			renderer.FillRect(&pixel)
			if k == 0 {
				k = uint8(1)
			} else {
				k = uint8(0)
			}
		}
	}
	renderer.Present()

	// Game Cycles
	clock := time.NewTicker(time.Millisecond)
	timer := time.NewTicker(time.Second / 60)
	video := time.NewTicker(time.Second / 60)
	fps := time.NewTicker(time.Second)
	frames := 0

	for processEvent() {
		select {
		case <-fps.C: // FPS Capture
			fmt.Println("Frames:", frames)
			frames = 0
		case <-clock.C: // Emulate Cycle.
		case <-timer.C: // SoundTimer and DelayTimer
		case <-video.C: // Draw
		default:
		}
		frames++
		sdl.Delay(16)
	}

	defer sdl.Quit()
}

func processEvent() bool {
	for event := sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
		switch event.(type) {
		case *sdl.QuitEvent:
			println("Quitting...")
			return false
		}
	}
	return true
}
