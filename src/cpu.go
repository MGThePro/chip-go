package main

import (
	"math/rand"
)

var opcode uint16
var I uint16
var pc uint16 = 0x200
var memory [4096]uint8
var V [16]uint8
var stack [16]uint16
var sp uint16
var delay_timer uint8
var sound_timer uint8
var key [16]uint8
var framebuffer [64][32]uint8
var fontset = [80]uint8 {
	0xF0, 0x90, 0x90, 0x90, 0xF0,
	0x20, 0x60, 0x20, 0x20, 0x70,
	0xF0, 0x10, 0xF0, 0x80, 0xF0,
	0xF0, 0x10, 0xF0, 0x10, 0xF0,
	0x90, 0x90, 0xF0, 0x10, 0x10,
	0xF0, 0x80, 0xF0, 0x10, 0xF0,
	0xF0, 0x80, 0xF0, 0x90, 0xF0,
	0xF0, 0x10, 0x20, 0x40, 0x40,
	0xF0, 0x90, 0xF0, 0x90, 0xF0,
	0xF0, 0x90, 0xF0, 0x10, 0xF0,
	0xF0, 0x90, 0xF0, 0x90, 0x90,
	0xE0, 0x90, 0xE0, 0x90, 0xE0,
	0xF0, 0x80, 0x80, 0x80, 0xF0,
	0xE0, 0x90, 0x90, 0x90, 0xE0,
	0xF0, 0x80, 0xF0, 0x80, 0xF0,
	0xF0, 0x80, 0xF0, 0x80, 0x80,
}

var funcmap = map[uint16] func(){
	0x0000 : op_0nnn,
	0x00E0 : op_00E0,
	0x00EE : op_00EE,
	0x1000 : op_1xxx,
	0x2000 : op_2xxx,
	0x3000 : op_3xkk,
	0x4000 : op_4xkk,
	0x5000 : op_5xy0,
	0x6000 : op_6xkk,
	0x7000 : op_7xkk,
	0x8000 : op_8xy0,
	0x8001 : op_8xy1,
	0x8002 : op_8xy2,
	0x8003 : op_8xy3,
	0x8004 : op_8xy4,
	0x8005 : op_8xy5,
	0x8006 : op_8xy6,
	0x8007 : op_8xy7,
	0x800E : op_8xyE,
	0x9000 : op_9xy0,
	0xA000 : op_Annn,
	0xB000 : op_Bnnn,
	0xC000 : op_Cxkk,
	0xD000 : op_Dxyn,
	0xE09E : op_Ex9E,
	0xE0A1 : op_ExA1,
	0xF007 : op_Fx07,
	0xF00A : op_Fx0A,
	0xF015 : op_Fx15,
	0xF018 : op_Fx18,
	0xF01E : op_Fx1E,
	0xF029 : op_Fx29,
	0xF033 : op_Fx33,
	0xF055 : op_Fx55,
	0xF065 : op_Fx65,
}

func fetch() {
	opcode = uint16(memory[pc])<<8 | uint16(memory[pc+1])
}

func execute() {
	parsedOpcode := opcode
	if opcode <= 0x0FFF && opcode != 0x00E0 && opcode != 0x00EE {
		parsedOpcode = opcode & 0x0000
	}
	if opcode >= 0x1000 && opcode <= 0x7FFF {
		parsedOpcode = opcode & 0xF000
	}
	if opcode >= 0x8000 && opcode <= 0x9FFF {
		parsedOpcode = opcode & 0xF00F
	}
	if opcode >= 0xA000 && opcode <= 0xDFFF {
		parsedOpcode = opcode & 0xF000
	}
	if opcode >= 0xE000 && opcode <= 0xFFFF {
		parsedOpcode = opcode & 0xF0FF
	}
	if funcmap[parsedOpcode] != nil {
		funcmap[parsedOpcode]()
	}
}

func op_0nnn() {
	nnn := opcode & 0x0FFF
	pc = nnn - 2
}

func op_00E0() {
	for x:= 0; 64 > x; x++{
		for y:= 0; 32 > y; y++{
			framebuffer[x][y] = 0
		}
	}
}

func op_00EE() {
	pc = stack[sp]
	sp--
}

func op_1xxx() {
	pc = opcode & 0x0FFF
	pc -= 2
}

func op_2xxx() {
	sp++
	stack[sp] = pc
	pc = opcode & 0x0FFF
}

func op_3xkk() {
	x := (opcode & 0x0F00) >> 8
	kk := opcode & 0x00FF
	if V[x] == uint8(kk) {
		pc += 2
	}
}

func op_4xkk() {
	x := (opcode & 0x0F00) >> 8
	kk := opcode & 0x00FF
	if V[x] != uint8(kk) {
		pc += 2
	}
}

func op_5xy0() {
	x := (opcode & 0x0F00) >> 8
	y := (opcode & 0x00F0) >> 4
	if V[x] == V[y] {
		pc += 2
	}
}

func op_6xkk() {
	x := (opcode & 0x0F00) >> 8
	kk := opcode & 0x00FF
	V[x] = uint8(kk)
}

func op_7xkk() {
	x := (opcode & 0x0F00) >> 8
	kk := opcode & 0x00FF
	V[x] += uint8(kk)
}

func op_8xy0() {
	x := (opcode & 0x0F00) >> 8
	y := (opcode & 0x00F0) >> 4
	V[x] = V[y]
}

func op_8xy1() {
	x := (opcode & 0x0F00) >> 8
	y := (opcode & 0x00F0) >> 4
	V[x] = V[x] | V[y]
}

func op_8xy2() {
	x := (opcode & 0x0F00) >> 8
	y := (opcode & 0x00F0) >> 4
	V[x] = V[x] & V[y]
}

func op_8xy3() {
	x := (opcode & 0x0F00) >> 8
	y := (opcode & 0x00F0) >> 4
	V[x] = V[x] ^ V[y]
}

func op_8xy4() {
	x := (opcode & 0x0F00) >> 8
	y := (opcode & 0x00F0) >> 4
	if uint16(V[x]) + uint16(V[y]) > 255 {
		V[0xF] = 1
	} else { V[0xF] = 0 }
	V[x] = V[x] + V[y]
}

func op_8xy5() {
	x := (opcode & 0x0F00) >> 8
	y := (opcode & 0x00F0) >> 4
	if V[x] > V[y] {
		V[0xF] = 1
	} else { V[0xF] = 0 }
	V[x] = V[x] - V[y]
}

func op_8xy6() {
	x := (opcode & 0x0F00) >> 8
	V[0xF] = uint8(x & 0x01)
	V[x] = V[x] >> 1
}

func op_8xy7() {
	x := (opcode & 0x0F00) >> 8
	y := (opcode & 0x00F0) >> 4
	if V[x] < V[y] {
		V[0xF] = 1
	} else { V[0xF] = 0 }
	V[x] = V[y] - V[x]
}

func op_8xyE() {
	x := (opcode & 0x0F00) >> 8
	V[0xF] = (V[x] & 0x80) >> 7
	V[x] = V[x] << 1
}

func op_9xy0() {
	x := (opcode & 0x0F00) >> 8
	y := (opcode & 0x00F0) >> 4
	if V[x] > V[y] { pc += 2 }
}

func op_Annn() {
	I = opcode & 0x0FFF
}

func op_Bnnn() {
	nnn := opcode & 0x0FFF
	pc = nnn + uint16(V[0])
	pc -= 2
}

func op_Cxkk() {
	x := (opcode & 0x0F00) >> 8
	kk := opcode & 0x00FF
	V[x] = uint8(rand.Int()) & uint8(kk)

}

func op_Dxyn() {
	x := (opcode & 0x0F00) >> 8
	y := (opcode & 0x00F0) >> 4
	n := (opcode & 0x000F)
	V[0xF] = 0
	startx := V[x] % 63
	starty := V[y] % 31
	for row := 0; uint16(row) < n; row++ {
		if uint8(row) + starty > 31 { break }
		spriteByte := memory[I + uint16(row)]
		for col := 0; col < 8; col++ {
			if uint8(col) + startx > 63 { break }
			spritePixel := spriteByte & (0x80 >> col)
			screenPixel := framebuffer[startx+uint8(col)][starty+uint8(row)]
			if spritePixel > 0 && screenPixel > 0 {
				V[0xF] = 1
			}
			framebuffer[startx+uint8(col)][starty+uint8(row)] = screenPixel ^ spritePixel
		}
	}
}

func op_Ex9E() {
	x := (opcode & 0x0F00) >> 8
	if key[x] == 1 { pc += 2 }
}

func op_ExA1() {
	x := (opcode & 0x0F00) >> 8
	if key[x] != 1 { pc += 2 }
}

func op_Fx07() {
	x := (opcode & 0x0F00) >> 8
	V[x] = delay_timer
}

func op_Fx0A() {
	x := (opcode & 0x0F00) >> 8
	for {
		for i := 0 ; i < 16 ; i++ {
			if key[i] == 1 {
				V[x] = uint8(i)
				return
			}
		}
	}
}

func op_Fx15() {
	x := (opcode & 0x0F00) >> 8
	delay_timer = V[x]
}

func op_Fx18() {
	x := (opcode & 0x0F00) >> 8
	sound_timer = V[x]
}

func op_Fx1E() {
	x := (opcode & 0x0F00) >> 8
	I = I + x
}

func op_Fx29() {
	x := (opcode & 0x0F00) >> 8
	I = x * 5
}

func op_Fx33() {
	x := (opcode & 0x0F00) >> 8
	memory[I] = (V[x] - (V[x] % 100)) / 100
	memory[I+1] = (V[x] - (V[x] % 10) - (memory[I]*100)) / 10
	memory[I+2] = V[x] - (memory[I]*100) - (memory[I+1]*10)
}

func op_Fx55() {
	x := (opcode & 0x0F00) >> 8
	for i:= 0; uint16(i) <= x; i++ {
		memory[I+uint16(i)] = V[i]
	}
}

func op_Fx65() {
	x := (opcode & 0x0F00) >> 8
	for i:= 0; uint16(i) <= x; i++ {
		V[i] = memory[I+uint16(i)]
	}
}