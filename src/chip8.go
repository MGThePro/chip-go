package main

import (
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/imdraw"
	"github.com/faiface/pixel/pixelgl"
	"os"
)

func main() {
	loadGame()
	loadFont()
	loadFuncs()
	pixelgl.Run(run)
}

func loadGame() {
	gameArr, err := os.ReadFile("game.ch8")
	if err != nil {
		panic(err)
	}
	for i := 0; i < len(gameArr); i++ {
		memory[0x200+i] = gameArr[i]
	}
}

func loadFont() {
	for i := 0; i < len(fontSet); i++ {
		memory[i] = fontSet[i]
	}
}

func drawFramebuffer(imd *imdraw.IMDraw) {
	imd.Color = pixel.RGB(1, 1, 1)
	for y := 0; y < 32; y++ {
		for x := 0; x < 64; x++ {
			if framebuffer[x][y] > 0 {
				imd.Color = pixel.RGB(1, 1, 1)
			} else {
				imd.Color = pixel.RGB(0, 0, 0)
			}
			imd.Push(pixel.V(float64(x*10), float64(310-y*10)))
			imd.Push(pixel.V(float64((x*10)+10), float64((310-y*10)+10)))
			imd.Rectangle(0)
		}
	}
}

func pollInput(win *pixelgl.Window) {
	if win.Pressed(pixelgl.Key0) {
		key[0] = 1
	} else {
		key[0] = 0
	}
	if win.Pressed(pixelgl.Key1) {
		key[1] = 1
	} else {
		key[1] = 0
	}
	if win.Pressed(pixelgl.Key2) {
		key[2] = 1
	} else {
		key[2] = 0
	}
	if win.Pressed(pixelgl.Key3) {
		key[3] = 1
	} else {
		key[3] = 0
	}
	if win.Pressed(pixelgl.Key4) {
		key[4] = 1
	} else {
		key[4] = 0
	}
	if win.Pressed(pixelgl.Key5) {
		key[5] = 1
	} else {
		key[5] = 0
	}
	if win.Pressed(pixelgl.Key6) {
		key[6] = 1
	} else {
		key[6] = 0
	}
	if win.Pressed(pixelgl.Key7) {
		key[7] = 1
	} else {
		key[7] = 0
	}
	if win.Pressed(pixelgl.Key8) {
		key[8] = 1
	} else {
		key[8] = 0
	}
	if win.Pressed(pixelgl.Key9) {
		key[9] = 1
	} else {
		key[9] = 0
	}
	if win.Pressed(pixelgl.KeyA) {
		key[10] = 1
	} else {
		key[10] = 0
	}
	if win.Pressed(pixelgl.KeyB) {
		key[11] = 1
	} else {
		key[11] = 0
	}
	if win.Pressed(pixelgl.KeyC) {
		key[12] = 1
	} else {
		key[12] = 0
	}
	if win.Pressed(pixelgl.KeyD) {
		key[13] = 1
	} else {
		key[13] = 0
	}
	if win.Pressed(pixelgl.KeyE) {
		key[14] = 1
	} else {
		key[14] = 0
	}
	if win.Pressed(pixelgl.KeyF) {
		key[15] = 1
	} else {
		key[15] = 0
	}
}

func run() {
	cfg := pixelgl.WindowConfig{
		Title:  "Chip-go",
		Bounds: pixel.R(0, 0, 640, 320),
		VSync:  true,
	}
	win, err := pixelgl.NewWindow(cfg)
	if err != nil {
		panic(err)
	}
	cyclesPerFrame := 500 / 60
	for !win.Closed() {
		cyclesDelta := 0
		pollInput(win)
		for cyclesDelta < cyclesPerFrame {
			cyclesDelta++
			fetch()
			execute()
			pc += 2
		}
		if soundTimer > 0 {
			soundTimer--
		}
		if delayTimer > 0 {
			delayTimer--
		}
		imd := imdraw.New(nil)
		drawFramebuffer(imd)
		imd.Draw(win)
		win.Update()
	}
}
