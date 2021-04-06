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
var delayTimer uint8
var soundTimer uint8
var key [16]uint8
var framebuffer [64][32]uint8
var fontSet = [80]uint8{
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

var funcArray [61542]func()

func loadFuncs() {
	funcArray[0x0000] = op_0xxx
	funcArray[0x00E0] = op_00E0
	funcArray[0x00EE] = op_00EE
	funcArray[0x1000] = op_1xxx
	funcArray[0x2000] = op_2xxx
	funcArray[0x3000] = op_3xkk
	funcArray[0x4000] = op_4xkk
	funcArray[0x5000] = op_5xy0
	funcArray[0x6000] = op_6xkk
	funcArray[0x7000] = op_7xkk
	funcArray[0x8000] = op_8xy0
	funcArray[0x8001] = op_8xy1
	funcArray[0x8002] = op_8xy2
	funcArray[0x8003] = op_8xy3
	funcArray[0x8004] = op_8xy4
	funcArray[0x8005] = op_8xy5
	funcArray[0x8006] = op_8xy6
	funcArray[0x8007] = op_8xy7
	funcArray[0x800E] = op_8xyE
	funcArray[0x9000] = op_9xy0
	funcArray[0xA000] = op_Annn
	funcArray[0xB000] = op_Bnnn
	funcArray[0xC000] = op_Cxkk
	funcArray[0xD000] = op_Dxyn
	funcArray[0xE09E] = op_Ex9E
	funcArray[0xE0A1] = op_ExA1
	funcArray[0xF007] = op_Fx07
	funcArray[0xF00A] = op_Fx0A
	funcArray[0xF015] = op_Fx15
	funcArray[0xF018] = op_Fx18
	funcArray[0xF01E] = op_Fx1E
	funcArray[0xF029] = op_Fx29
	funcArray[0xF033] = op_Fx33
	funcArray[0xF055] = op_Fx55
	funcArray[0xF065] = op_Fx65
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
	if opcode >= 0xE000 {
		parsedOpcode = opcode & 0xF0FF
	}
	if parsedOpcode <= 0xF065 && funcArray[parsedOpcode] != nil {
		funcArray[parsedOpcode]()
	}
}

func op_0xxx() {
	nnn := opcode & 0x0FFF
	sp++
	stack[sp] = pc
	pc = nnn - 2
}

func op_00E0() {
	for x := 0; 64 > x; x++ {
		for y := 0; 32 > y; y++ {
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
	pc = (opcode & 0x0FFF) - 2
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
	if uint16(V[x])+uint16(V[y]) > 255 {
		V[0xF] = 1
	} else {
		V[0xF] = 0
	}
	V[x] = V[x] + V[y]
}

func op_8xy5() {
	x := (opcode & 0x0F00) >> 8
	y := (opcode & 0x00F0) >> 4
	if V[x] >= V[y] {
		V[0xF] = 1
	} else {
		V[0xF] = 0
	}
	V[x] -= V[y]
}

func op_8xy6() {
	x := (opcode & 0x0F00) >> 8
	V[0xF] = V[x] & 0x01
	V[x] = V[x] >> 1
}

func op_8xy7() {
	x := (opcode & 0x0F00) >> 8
	y := (opcode & 0x00F0) >> 4
	if V[x] <= V[y] {
		V[0xF] = 1
	} else {
		V[0xF] = 0
	}
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
	if V[x] != V[y] {
		pc += 2
	}
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
	n := opcode & 0x000F
	V[0xF] = 0
	startX := V[x] % 63
	startY := V[y] % 31
	for row := 0; uint16(row) < n; row++ {
		if uint8(row)+startY > 31 {
			break
		}
		spriteByte := memory[I+uint16(row)]
		for col := 0; col < 8; col++ {
			if uint8(col)+startX > 63 {
				break
			}
			spritePixel := spriteByte & (0x80 >> col)
			screenPixel := framebuffer[startX+uint8(col)][startY+uint8(row)]
			if spritePixel > 0 && screenPixel > 0 {
				V[0xF] = 1
			}
			framebuffer[startX+uint8(col)][startY+uint8(row)] = screenPixel ^ spritePixel
		}
	}
}

func op_Ex9E() {
	x := (opcode & 0x0F00) >> 8
	if key[V[x]] == 1 {
		pc += 2
	}
}

func op_ExA1() {
	x := (opcode & 0x0F00) >> 8
	if key[V[x]] != 1 {
		pc += 2
	}
}

func op_Fx07() {
	x := (opcode & 0x0F00) >> 8
	V[x] = delayTimer
}

func op_Fx0A() {
	x := (opcode & 0x0F00) >> 8
	for i := 0; i < 16; i++ {
		if key[i] == 1 {
			V[x] = uint8(i)
			return
		}
	}
	pc -= 2
}

func op_Fx15() {
	x := (opcode & 0x0F00) >> 8
	delayTimer = V[x]
}

func op_Fx18() {
	x := (opcode & 0x0F00) >> 8
	soundTimer = V[x]
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
	memory[I+1] = (V[x] - (V[x] % 10) - (memory[I] * 100)) / 10
	memory[I+2] = V[x] - (memory[I] * 100) - (memory[I+1] * 10)
}

func op_Fx55() {
	x := (opcode & 0x0F00) >> 8
	for i := 0; uint16(i) <= x; i++ {
		memory[I] = V[i]
		I++
	}
}

func op_Fx65() {
	x := (opcode & 0x0F00) >> 8
	for i := 0; uint16(i) <= x; i++ {
		V[i] = memory[I]
		I++
	}
}
