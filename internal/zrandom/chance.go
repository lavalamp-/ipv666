package zrandom

import (
	"math/rand"
)

func GenerateHostBits(bitCount int) []byte {
	var toReturn []byte
	var curByte byte = 0x00
	curPos := 0
	for i := 0; i < bitCount; i++ {
		curByte = (curByte << 1) | (byte)(uint8(rand.Intn(2)))
		curPos ++
		if curPos == 8 {
			toReturn = append([]byte{curByte}, toReturn...)
			curByte = 0x00
			curPos = 0
		}
	}
	if curPos != 0 {
		toReturn = append([]byte{curByte}, toReturn...)
	}
	for len(toReturn) < 16 {
		toReturn = append([]byte{0x00}, toReturn...)
	}
	return toReturn
}

func GenerateRandomBits(bitCount uint8) []byte {
	var toReturn []byte
	var curByte byte = 0x00
	curPos := 0
	var i uint8
	for i = 0; i < bitCount; i++ {
		curByte = (curByte << 1) | (byte)(uint8(rand.Intn(2)))
		curPos ++
		if curPos == 8 {
			toReturn = append([]byte{curByte}, toReturn...)
			curByte = 0x00
			curPos = 0
		}
	}
	if curPos != 0 {
		toReturn = append([]byte{curByte}, toReturn...)
	}
	return toReturn
}
