package main

import (
	"fmt"
	"github.com/thescripted/Octal/chip8"
	"github.com/veandco/go-sdl2/sdl"
	"time"
)

const (
	screenWidth  = 640
	screenHeight = 320
	pixelWidth   = screenWidth / 64
	pixelHeight  = screenHeight / 32
)

var (
	Chip     *chip8.Chip8
	window   *sdl.Window
	renderer *sdl.Renderer
	pixel    sdl.Rect
	err      error
	KeyMap   = map[sdl.Scancode]uint{
		sdl.SCANCODE_X: 0x0,
		sdl.SCANCODE_1: 0x1,
		sdl.SCANCODE_2: 0x2,
		sdl.SCANCODE_3: 0x3,
		sdl.SCANCODE_Q: 0x4,
		sdl.SCANCODE_W: 0x5,
		sdl.SCANCODE_E: 0x6,
		sdl.SCANCODE_A: 0x7,
		sdl.SCANCODE_S: 0x8,
		sdl.SCANCODE_D: 0x9,
		sdl.SCANCODE_Z: 0xA,
		sdl.SCANCODE_C: 0xB,
		sdl.SCANCODE_4: 0xC,
		sdl.SCANCODE_R: 0xD,
		sdl.SCANCODE_F: 0xE,
		sdl.SCANCODE_V: 0xF,
	}
)

func main() {
	// Initialize Chip
	Chip = chip8.New()
	if err = Chip.LoadProgram("./rom/IBM Logo.ch8"); err != nil {
		panic(err)
	}

	// Initialize Graphics
	if err = sdl.Init(sdl.INIT_VIDEO); err != nil {
		panic(err)
	}
	window, err = sdl.CreateWindow("CHIP - 8", sdl.WINDOWPOS_UNDEFINED, sdl.WINDOWPOS_UNDEFINED, screenWidth, screenHeight, 0)
	if err != nil {
		panic(err)
	}
	defer window.Destroy()
	renderer, err = sdl.CreateRenderer(window, -1, sdl.RENDERER_ACCELERATED)
	if err != nil {
		panic(err)
	}
	defer renderer.Destroy()
	sdl.SetHint(sdl.HINT_RENDER_SCALE_QUALITY, "linear")

	// Set background
	renderer.SetDrawColor(96, 128, 128, 255)
	renderer.Clear()
	renderer.Present()

	// Initialize Pixel
	pixel = sdl.Rect{
		X: 0,
		Y: 0,
		W: pixelWidth,
		H: pixelHeight,
	}

	// Game Cycles
	clock := time.NewTicker(time.Millisecond)
	timer := time.NewTicker(time.Second / 60)
	video := time.NewTicker(time.Second / 60)
	fps := time.NewTicker(time.Second)
	frames := 0

	// Game Loop
	for processEvent() {
		select {
		case <-fps.C: // FPS Capture
			fmt.Println("Frames:", frames)

			frames = 0
		case <-clock.C: // Emulate Cycle.
			Chip.Tick()
		case <-timer.C: // SoundTimer and DelayTimer
		case <-video.C: // Draw
			draw(Chip.Video)
		default:
		}
		frames++
		sdl.Delay(1) // 1ms Delay
	}

	defer sdl.Quit()
}

func draw(video [0x800]byte) {
	var k byte
	for i := 0; i < 64; i++ {
		pixel.X = int32(i * pixelWidth)
		for j := 0; j < 32; j++ {
			k = video[i+64*j]
			pixel.Y = int32(j * pixelHeight)
			renderer.SetDrawColor(255*k, 255*k, 255*k, 255)
			renderer.FillRect(&pixel)
		}
	}
	renderer.Present()
}

func processEvent() bool {
	for event := sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
		switch ev := event.(type) {
		case *sdl.QuitEvent:
			return false
		case *sdl.KeyboardEvent:
			if ev.Type == sdl.KEYUP {
				if key, ok := KeyMap[ev.Keysym.Scancode]; ok {
					fmt.Println("Releasing:", key, ev.Keysym.Scancode)
					Chip.ReleaseKey(key)
				}
			} else {
				if key, ok := KeyMap[ev.Keysym.Scancode]; ok {
					fmt.Println("Pressing:", key, ev.Keysym.Scancode)
					Chip.PressKey(key)
				}
				// Should listen to other Keyboard Events here...
			}
		}
	}
	return true
}
