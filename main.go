package main

import (
	"flag"
	"fmt"
	"github.com/manifoldco/promptui"
	"github.com/pkg/errors"
	"github.com/thescripted/Octal/chip8"
	"github.com/veandco/go-sdl2/sdl"
	"time"

	"io/ioutil"
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
)

var KeyMap = map[sdl.Scancode]uint{
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

// TODO(ben): Support debug flag

func main() {
	testOpCode := flag.Bool("test-op", false, "Run the opcode test on the emulator.")

	window, renderer, err := initializeSDL()
	if err != nil {
		panic(err)
	}
	defer window.Destroy()
	defer renderer.Destroy()
	defer sdl.Quit()
	renderer.SetDrawColor(96, 128, 128, 255)
	renderer.Clear()
	renderer.Present()
	pixel = sdl.Rect{
		X: 0,
		Y: 0,
		W: pixelWidth,
		H: pixelHeight,
	}

	Chip = chip8.New()
	program, err := getProgram(testOpCode)
	if err != nil {
		panic(err)
	}
	err = Chip.LoadProgram(program)
	if err != nil {
		panic(err)
	}

	clock := time.NewTicker(time.Millisecond)
	timer := time.NewTicker(time.Second / 60)
	video := time.NewTicker(time.Second / 60)
	fps := time.NewTicker(time.Second)
	frames := 0
	for processEvent() {
		select {
		case <-fps.C:
			fmt.Println("Frames:", frames)
			frames = 0
		case <-clock.C:
			Chip.Tick()
		case <-timer.C:
		case <-video.C:
			draw(Chip.Video)
		default:
		}
		frames++
		sdl.Delay(1)
	}
}

func getProgram(testOpCode *bool) (string, error) {
	romDir := "./rom"
	files, err := ioutil.ReadDir(romDir)
	if err != nil {
		return "", errors.Wrap(err, "Unable to read files from Directory.")
	}
	if *testOpCode {
		return "test_opcode.ch8", nil
	}
	var selection []string
	for _, f := range files {
		selection = append(selection, f.Name())
	}
	prompt := promptui.Select{
		Label: "Select a game",
		Items: selection,
	}
	_, result, err := prompt.Run()
	if err != nil {
		return "", errors.Wrap(err, "Unable to run Prompt.")
	}
	return romDir + "/" + result, nil

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

				if ev.Keysym.Scancode == sdl.SCANCODE_RIGHT {
					// TODO(ben): Debug Mode
				}
			}
		}
	}
	return true
}
func initializeSDL() (*sdl.Window, *sdl.Renderer, error) {
	if err = sdl.Init(sdl.INIT_VIDEO); err != nil {
		panic(err)
	}
	sdl.SetHint(sdl.HINT_RENDER_SCALE_QUALITY, "linear")
	window, err = sdl.CreateWindow("CHIP - 8", sdl.WINDOWPOS_UNDEFINED, sdl.WINDOWPOS_UNDEFINED, screenWidth, screenHeight, 0)
	if err != nil {
		return window, renderer, errors.Wrap(err, "Create window failed.")
	}
	renderer, err = sdl.CreateRenderer(window, -1, sdl.RENDERER_ACCELERATED)
	if err != nil {
		return window, renderer, errors.Wrap(err, "Create renderer failed.")
	}
	return window, renderer, nil
}
